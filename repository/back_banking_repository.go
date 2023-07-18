package repository

import (
	"bytes"
	"cybergame-api/model"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
)

func NewBankingRepository(db *gorm.DB) BankingRepository {
	return &repo{db}
}

type BankingRepository interface {
	GetBankStatementById(id int64) (*model.BankStatement, error)
	GetBankStatements(req model.BankStatementListRequest) (*model.SuccessWithPagination, error)
	GetBankExternalStatements(externalIds []int64) (*model.SuccessWithPagination, error)
	HasBankExternalStatements(externalId int64) error
	GetBankStatementSummary(req model.BankStatementListRequest) (*model.BankStatementSummary, error)
	CreateBankStatement(data model.BankStatementCreateBody) error
	UpdateBankStatement(id int64, data model.BankStatementUpdateBody) error
	MatchStatementOwner(id int64, data model.BankStatementUpdateBody) error
	IgnoreStatementOwner(id int64, data model.BankStatementUpdateBody) error
	DeleteBankStatement(id int64) error

	GetBankTransactionById(id int64) (*model.BankTransaction, error)
	GetBankTransactions(req model.BankTransactionListRequest) (*model.SuccessWithPagination, error)
	GetBankTransactionStatusCount(req model.BankTransactionListRequest) (*model.SuccessWithPagination, error)
	CreateBankDepositTransaction(data model.BankTransactionCreateBody) (*int64, error)
	CreateBankWithdrawTransaction(data model.BankTransactionCreateBody) (*int64, error)
	CreateBonusTransaction(data model.BonusTransactionCreateBody) error
	UpdateBankTransaction(id int64, data interface{}) error
	DeleteBankTransaction(id int64) error

	GetPendingDepositTransactions(req model.PendingDepositTransactionListRequest) (*model.SuccessWithPagination, error)
	GetPendingWithdrawTransactions(req model.PendingWithdrawTransactionListRequest) (*model.SuccessWithPagination, error)
	CreateTransactionAction(data model.CreateBankTransactionActionBody) (*int64, error)
	RollbackTransactionAction(actionId int64) error
	CreateStatementAction(data model.CreateBankStatementActionBody) error
	ConfirmPendingDepositTransaction(id int64, data model.BankDepositTransactionConfirmBody) error
	ConfirmPendingCreditDepositTransaction(id int64, data model.BankDepositTransactionConfirmBody) error
	CheckMemeberHasEnoughtCredit(memberId int64, creditAmount float64) error
	ConfirmPendingWithdrawTransaction(id int64, data model.BankWithdrawTransactionConfirmBody) error
	ConfirmPendingWithdrawTransfer(id int64, data model.BankWithdrawTransactionConfirmBody) error
	CancelPendingTransaction(id int64, data model.BankTransactionCancelBody) error
	GetFinishedTransactions(req model.FinishedTransactionListRequest) (*model.SuccessWithPagination, error)
	RemoveFinishedTransaction(id int64, data model.BankTransactionRemoveBody) error
	GetRemovedTransactions(req model.RemovedTransactionListRequest) (*model.SuccessWithPagination, error)

	GetMemberById(id int64) (*model.Member, error)
	GetMemberByCode(code string) (*model.Member, error)
	GetMembers(req model.MemberListRequest) (*model.SuccessWithPagination, error)
	GetPossibleStatementOwners(req model.MemberPossibleListRequest) (*model.SuccessWithPagination, error)
	GetMemberTransactions(req model.MemberTransactionListRequest) (*model.SuccessWithPagination, error)
	GetMemberTransactionSummary(req model.MemberTransactionListRequest) (*model.MemberTransactionSummary, error)
	IncreaseMemberCredit(body model.MemberStatementCreateBody) error
	DecreaseMemberCredit(body model.MemberStatementCreateBody) error

	TransferExternalAccount(body model.ExternalAccountTransferBody) error

	GetMemberStatementById(id int64) (*model.MemberStatementResponse, error)
	GetMemberStatements(req model.MemberStatementListRequest) (*model.SuccessWithPagination, error)
	GetMemberStatementTypeByCode(code string) (*model.MemberStatementType, error)
	GetMemberStatementTypes(req model.SimpleListRequest) (*model.SuccessWithPagination, error)
	CreateMemberStatement(data model.MemberStatementCreateBody) (*int64, error)
}

func (r repo) GetBankStatementById(id int64) (*model.BankStatement, error) {
	var record model.BankStatement
	selectedFields := "statements.id, statements.account_id, statements.external_id, statements.detail, statements.statement_type, statements.transfer_at, statements.from_bank_id, statements.from_account_number, statements.amount, statements.status, statements.created_at, statements.updated_at"
	selectedFields += ",accounts.account_name, accounts.account_number, accounts.account_type_id, accounts.bank_id"
	selectedFields += ",banks.name as bank_name, banks.code as bank_code, banks.icon_url as bank_icon_url, banks.type_flag as bank_type_flag"
	selectedFields += ",from_banks.name as from_bank_name, from_banks.code as from_bank_code, from_banks.icon_url as from_bank_icon_url"
	if err := r.db.Table("Bank_statements as statements").
		Select(selectedFields).
		Joins("LEFT JOIN Bank_accounts AS accounts ON accounts.id = statements.account_id").
		Joins("LEFT JOIN Banks AS banks ON banks.id = accounts.bank_id").
		Joins("LEFT JOIN Banks AS from_banks ON from_banks.id = statements.from_bank_id").
		Where("statements.id = ?", id).
		Where("statements.deleted_at IS NULL").
		First(&record).
		Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r repo) GetBankStatements(req model.BankStatementListRequest) (*model.SuccessWithPagination, error) {

	var list []model.BankStatementResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset)
	count := r.db.Table("Bank_statements as statements")
	count = count.Joins("LEFT JOIN Bank_accounts AS accounts ON accounts.id = statements.account_id")
	count = count.Joins("LEFT JOIN Banks AS banks ON banks.id = accounts.bank_id")
	count = count.Joins("LEFT JOIN Banks AS from_banks ON from_banks.id = statements.from_bank_id")
	count = count.Select("statements.id")
	if req.AccountId != "" {
		count = count.Where("statements.account_id = ?", req.AccountId)
	}
	if req.FromTransferDate != "" {
		count = count.Where("statements.transfer_at >= ?", req.FromTransferDate)
	}
	if req.ToTransferDate != "" {
		count = count.Where("statements.transfer_at <= ?", req.ToTransferDate)
	}
	if req.StatementType != "" {
		count = count.Where("statements.statement_type = ?", req.StatementType)
	}
	if req.Status != "" {
		count = count.Where("statements.status = ?", req.Status)
	}
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where(r.db.Where("accounts.account_name LIKE ?", search_like).Or("accounts.account_number LIKE ?", search_like))
	}

	if err = count.
		Where("statements.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}
	if total > 0 {
		// SELECT //
		selectedFields := "statements.id, statements.account_id, statements.external_id, statements.detail, statements.statement_type, statements.transfer_at, statements.from_bank_id, statements.from_account_number, statements.amount, statements.status, statements.created_at, statements.updated_at"
		selectedFields += ",accounts.account_name, accounts.account_number, accounts.account_type_id, accounts.bank_id"
		selectedFields += ",banks.name as bank_name, banks.code as bank_code, banks.icon_url as bank_icon_url, banks.type_flag as bank_type_flag"
		selectedFields += ",from_banks.name as from_bank_name, from_banks.code as from_bank_code, from_banks.icon_url as from_bank_icon_url"
		query := r.db.Table("Bank_statements as statements")
		query = query.Select(selectedFields)
		query = query.Joins("LEFT JOIN Bank_accounts AS accounts ON accounts.id = statements.account_id")
		query = query.Joins("LEFT JOIN Banks AS banks ON banks.id = accounts.bank_id")
		query = query.Joins("LEFT JOIN Banks AS from_banks ON from_banks.id = statements.from_bank_id")
		if req.AccountId != "" {
			query = query.Where("statements.account_id = ?", req.AccountId)
		}
		if req.FromTransferDate != "" {
			query = query.Where("statements.transfer_at >= ?", req.FromTransferDate)
		}
		if req.ToTransferDate != "" {
			query = query.Where("statements.transfer_at <= ?", req.ToTransferDate)
		}
		if req.StatementType != "" {
			query = query.Where("statements.statement_type = ?", req.StatementType)
		}
		if req.Status != "" {
			query = query.Where("statements.status = ?", req.Status)
		}
		if req.Search != "" {
			search_like := fmt.Sprintf("%%%s%%", req.Search)
			query = query.Where(r.db.Where("accounts.account_name LIKE ?", search_like).Or("accounts.account_number LIKE ?", search_like))
		}

		// Sort by ANY //
		req.SortCol = strings.TrimSpace(req.SortCol)
		if req.SortCol != "" {
			if strings.ToLower(strings.TrimSpace(req.SortAsc)) == "desc" {
				req.SortAsc = "DESC"
			} else {
				req.SortAsc = "ASC"
			}
			query = query.Order(req.SortCol + " " + req.SortAsc)
		}
		if req.Limit > 0 {
			query = query.Limit(req.Limit)
		}
		if err = query.
			Where("statements.deleted_at IS NULL").
			Offset(req.Page * req.Limit).
			Scan(&list).
			Error; err != nil {
			return nil, err
		}
	}

	// End count total records for pagination purposes (without limit and offset)
	var result model.SuccessWithPagination
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) GetBankExternalStatements(externalIds []int64) (*model.SuccessWithPagination, error) {

	var list []model.BankStatementResponse
	var total int64 = 1
	var err error

	if total > 0 {
		// SELECT //
		selectedFields := "statements.id, statements.account_id, statements.external_id, statements.detail, statements.statement_type, statements.transfer_at, statements.from_bank_id, statements.from_account_number, statements.amount, statements.status, statements.created_at, statements.updated_at"

		query := r.db.Table("Bank_statements as statements")
		query = query.Select(selectedFields)

		if externalIds != nil {
			query = query.Where("statements.external_id IN ?", externalIds)
		}

		if err = query.
			Where("statements.deleted_at IS NULL").
			Scan(&list).
			Error; err != nil {
			return nil, err
		}
	}

	// End count total records for pagination purposes (without limit and offset)
	var result model.SuccessWithPagination
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) HasBankExternalStatements(externalId int64) error {

	var total int64 = 1
	var err error

	// Count total records for pagination purposes (without limit and offset)
	count := r.db.Table("Bank_statements as statements")
	count = count.Select("statements.id")
	count = count.Where("statements.external_id = ?", externalId)
	if err = count.
		Where("statements.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return err
	}
	if total > 0 {
		return nil
	}
	return errors.New("record not found")
}

func (r repo) GetBankStatementSummary(req model.BankStatementListRequest) (*model.BankStatementSummary, error) {

	var result model.BankStatementSummary
	var totalPendingStatementCount int64
	var totalPendingDepositCount int64
	var totalPendingWithdrawCount int64
	var err error

	// Count total records for pagination purposes (without limit and offset)
	count := r.db.Table("Bank_statements as statements")
	count = count.Joins("LEFT JOIN Bank_accounts AS accounts ON accounts.id = statements.account_id")
	count = count.Select("statements.id")
	count = count.Where("statements.status = ?", "pending")
	if req.AccountId != "" {
		count = count.Where("statements.account_id = ?", req.AccountId)
	}
	if req.StatementType != "" {
		count = count.Where("statements.statement_type = ?", req.StatementType)
	}
	if req.FromTransferDate != "" {
		count = count.Where("statements.transfer_at >= ?", req.FromTransferDate)
	}
	if req.ToTransferDate != "" {
		count = count.Where("statements.transfer_at <= ?", req.ToTransferDate)
	}
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where(r.db.Where("accounts.account_name LIKE ?", search_like).Or("accounts.account_number LIKE ?", search_like))
	}
	if err = count.
		Where("statements.deleted_at IS NULL").
		Count(&totalPendingStatementCount).
		Error; err != nil {
		return nil, err
	}

	// Count total records for pagination purposes (without limit and offset)
	countDeposit := r.db.Table("Bank_transactions as transactions")
	countDeposit = countDeposit.Select("transactions.id")
	countDeposit = countDeposit.Where("transactions.status = ?", "pending")
	countDeposit = countDeposit.Where("transactions.transfer_type = ?", "deposit")
	if req.AccountId != "" {
		countDeposit = countDeposit.Where("transactions.to_account_id = ?", req.AccountId)
	}
	if req.FromTransferDate != "" {
		countDeposit = countDeposit.Where("transactions.transfer_at >= ?", req.FromTransferDate)
	}
	if req.ToTransferDate != "" {
		countDeposit = countDeposit.Where("transactions.transfer_at <= ?", req.ToTransferDate)
	}
	if err = countDeposit.
		Where("transactions.deleted_at IS NULL").
		Count(&totalPendingDepositCount).
		Error; err != nil {
		return nil, err
	}

	// Count total records for pagination purposes (without limit and offset)
	countWithdraw := r.db.Table("Bank_transactions as transactions")
	countWithdraw = countWithdraw.Select("transactions.id")
	countWithdraw = countWithdraw.Where("transactions.status = ?", "pending")
	countWithdraw = countWithdraw.Where("transactions.transfer_type = ?", "withdraw")
	if req.AccountId != "" {
		countWithdraw = countWithdraw.Where("transactions.from_account_id = ?", req.AccountId)
	}
	if req.FromTransferDate != "" {
		countWithdraw = countWithdraw.Where("transactions.transfer_at >= ?", req.FromTransferDate)
	}
	if req.ToTransferDate != "" {
		countWithdraw = countWithdraw.Where("transactions.transfer_at <= ?", req.ToTransferDate)
	}
	if err = countWithdraw.
		Where("transactions.deleted_at IS NULL").
		Count(&totalPendingWithdrawCount).
		Error; err != nil {
		return nil, err
	}

	result.TotalPendingStatementCount = totalPendingStatementCount
	result.TotalPendingDepositCount = totalPendingDepositCount
	result.TotalPendingWithdrawCount = totalPendingWithdrawCount

	return &result, nil
}

func (r repo) CreateBankStatement(data model.BankStatementCreateBody) error {
	if err := r.db.Table("Bank_statements").Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) UpdateBankStatement(id int64, data model.BankStatementUpdateBody) error {
	if err := r.db.Table("Bank_statements").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) MatchStatementOwner(id int64, data model.BankStatementUpdateBody) error {
	if err := r.db.Table("Bank_statements").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) IgnoreStatementOwner(id int64, data model.BankStatementUpdateBody) error {
	if err := r.db.Table("Bank_statements").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) DeleteBankStatement(id int64) error {
	if err := r.db.Table("Bank_statements").Where("id = ?", id).Delete(&model.BankStatement{}).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) GetBankTransactionById(id int64) (*model.BankTransaction, error) {
	var record model.BankTransaction
	selectedFields := "transactions.id, transactions.user_id, transactions.transfer_type, transactions.promotion_id, transactions.from_account_id, transactions.from_bank_id, transactions.from_account_name, transactions.from_account_number, transactions.to_account_id, transactions.to_bank_id, transactions.to_account_name, transactions.to_account_number"
	selectedFields += ", transactions.credit_amount, transactions.paid_amount, transactions.over_amount, transactions.deposit_channel, transactions.bonus_amount, transactions.bonus_reason, transactions.before_amount, transactions.after_amount, transactions.bank_charge_amount"
	selectedFields += ", transactions.transfer_at, transactions.created_by_user_id, transactions.created_by_username, transactions.removed_at, transactions.removed_by_user_id, transactions.removed_by_username, transactions.status, transactions.status_detail, transactions.is_auto_credit"
	selectedFields += ", transactions.created_at, transactions.updated_at"
	selectedFields += ", from_banks.name as from_bank_name, from_banks.code as from_bank_code, from_banks.icon_url as from_bank_icon_url, from_banks.type_flag as from_bank_type_flag"
	selectedFields += ", to_banks.name as to_bank_name, to_banks.code as to_bank_code, to_banks.icon_url as to_bank_icon_url, to_banks.type_flag as to_bank_type_flag"
	selectedFields += ", users.member_code as member_code, users.username as user_username, users.firstname as user_firstname, users.lastname as user_lastname, users.fullname as user_fullname, users.phone as user_phone"
	if err := r.db.Table("Bank_transactions as transactions").
		Select(selectedFields).
		Joins("LEFT JOIN Banks AS from_banks ON from_banks.id = transactions.from_bank_id").
		Joins("LEFT JOIN Banks AS to_banks ON to_banks.id = transactions.to_bank_id").
		Joins("LEFT JOIN Users AS users ON users.id = transactions.user_id").
		Where("transactions.id = ?", id).
		Where("transactions.deleted_at IS NULL").
		First(&record).
		Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r repo) GetBankTransactionStatusCount(req model.BankTransactionListRequest) (*model.SuccessWithPagination, error) {

	var list []model.BankTransactionStatusCount
	var total int64 = 6
	var err error

	// SELECT //
	selectedFields := "transactions.status, count(*) as count"
	query := r.db.Table("Bank_transactions as transactions")
	query = query.Select(selectedFields)
	query = query.Joins("LEFT JOIN Banks AS from_banks ON from_banks.id = transactions.from_bank_id")
	query = query.Joins("LEFT JOIN Banks AS to_banks ON to_banks.id = transactions.to_bank_id")
	query = query.Joins("LEFT JOIN Users AS users ON users.id = transactions.user_id")
	query = query.Where("transactions.removed_at IS NULL")
	if req.MemberCode != "" {
		query = query.Where("transactions.member_code = ?", req.MemberCode)
	}
	if req.UserId != "" {
		query = query.Where("transactions.user_id = ?", req.UserId)
	}
	if req.FromTransferDate != "" {
		query = query.Where("transactions.transfer_at >= ?", req.FromTransferDate)
	}
	if req.ToTransferDate != "" {
		query = query.Where("transactions.transfer_at <= ?", req.ToTransferDate)
	}
	if req.TransferType != "" {
		if req.TransferType == "all_deposit" {
			query = query.Where(r.db.Where("transactions.transfer_type = ?", "deposit").Or("transactions.transfer_type = ?", "bonus"))
		} else if req.TransferType == "all_withdraw" {
			query = query.Where(r.db.Where("transactions.transfer_type = ?", "withdraw").Or("transactions.transfer_type = ?", "getcreditback"))
		} else {
			query = query.Where("transactions.transfer_type = ?", req.TransferType)
		}
	}
	if req.TransferStatus != "" {
		if req.TransferType == "failed" {
			query = query.Where(r.db.Where("transactions.transfer_type = ?", "canceled"))
		} else {
			query = query.Where("transactions.status = ?", req.TransferStatus)
		}
	}
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		query = query.Where(r.db.Where("transactions.from_account_name LIKE ?", search_like).Or("transactions.from_account_number LIKE ?", search_like).Or("transactions.to_account_name LIKE ?", search_like).Or("transactions.to_account_number LIKE ?", search_like))
	}

	if err = query.
		Where("transactions.deleted_at IS NULL").
		Group("transactions.status").
		Scan(&list).
		Error; err != nil {
		return nil, err
	}

	// End count total records for pagination purposes (without limit and offset)
	var result model.SuccessWithPagination
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) GetBankTransactions(req model.BankTransactionListRequest) (*model.SuccessWithPagination, error) {

	var list []model.BankTransactionResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset)
	count := r.db.Table("Bank_transactions as transactions")
	count = count.Select("transactions.id")
	count = count.Where("transactions.removed_at IS NULL")
	if req.MemberCode != "" {
		count = count.Where("transactions.member_code = ?", req.MemberCode)
	}
	if req.UserId != "" {
		count = count.Where("transactions.user_id = ?", req.UserId)
	}
	if req.FromTransferDate != "" {
		count = count.Where("transactions.transfer_at >= ?", req.FromTransferDate)
	}
	if req.ToTransferDate != "" {
		count = count.Where("transactions.transfer_at <= ?", req.ToTransferDate)
	}
	if req.TransferType != "" {
		if req.TransferType == "all_deposit" {
			count = count.Where(r.db.Where("transactions.transfer_type = ?", "deposit").Or("transactions.transfer_type = ?", "bonus"))
		} else if req.TransferType == "all_withdraw" {
			count = count.Where(r.db.Where("transactions.transfer_type = ?", "withdraw").Or("transactions.transfer_type = ?", "getcreditback"))
		} else {
			count = count.Where("transactions.transfer_type = ?", req.TransferType)
		}
	}
	if req.TransferStatus != "" {
		if req.TransferStatus == "failed" {
			count = count.Where(r.db.Where("transactions.status = ?", "canceled"))
		} else {
			count = count.Where("transactions.status = ?", req.TransferStatus)
		}
	}
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where(r.db.Where("transactions.from_account_name LIKE ?", search_like).Or("transactions.from_account_number LIKE ?", search_like).Or("transactions.to_account_name LIKE ?", search_like).Or("transactions.to_account_number LIKE ?", search_like))
	}

	if err = count.
		Where("transactions.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}
	if total > 0 {
		// SELECT //
		selectedFields := "transactions.id, transactions.user_id, transactions.transfer_type, transactions.promotion_id, transactions.from_account_id, transactions.from_bank_id, transactions.from_account_name, transactions.from_account_number, transactions.to_account_id, transactions.to_bank_id, transactions.to_account_name, transactions.to_account_number"
		selectedFields += ", transactions.credit_amount, transactions.paid_amount, transactions.over_amount, transactions.deposit_channel, transactions.bonus_amount, transactions.bonus_reason, transactions.before_amount, transactions.after_amount, transactions.bank_charge_amount"
		selectedFields += ", transactions.transfer_at, transactions.created_by_user_id, transactions.created_by_username, transactions.removed_at, transactions.removed_by_user_id, transactions.removed_by_username, transactions.status, transactions.status_detail, transactions.is_auto_credit"
		selectedFields += ", transactions.created_at, transactions.updated_at"
		selectedFields += ", from_banks.name as from_bank_name, from_banks.code as from_bank_code, from_banks.icon_url as from_bank_icon_url, from_banks.type_flag as from_bank_type_flag"
		selectedFields += ", to_banks.name as to_bank_name, to_banks.code as to_bank_code, to_banks.icon_url as to_bank_icon_url, to_banks.type_flag as to_bank_type_flag"
		selectedFields += ", users.member_code as member_code, users.username as user_username, users.firstname as user_firstname, users.lastname as user_lastname, users.fullname as user_fullname, users.phone as user_phone"
		query := r.db.Table("Bank_transactions as transactions")
		query = query.Select(selectedFields)
		query = query.Joins("LEFT JOIN Banks AS from_banks ON from_banks.id = transactions.from_bank_id")
		query = query.Joins("LEFT JOIN Banks AS to_banks ON to_banks.id = transactions.to_bank_id")
		query = query.Joins("LEFT JOIN Users AS users ON users.id = transactions.user_id")
		query = query.Where("transactions.removed_at IS NULL")
		if req.MemberCode != "" {
			query = query.Where("transactions.member_code = ?", req.MemberCode)
		}
		if req.UserId != "" {
			query = query.Where("transactions.user_id = ?", req.UserId)
		}
		if req.FromTransferDate != "" {
			query = query.Where("transactions.transfer_at >= ?", req.FromTransferDate)
		}
		if req.ToTransferDate != "" {
			query = query.Where("transactions.transfer_at <= ?", req.ToTransferDate)
		}
		if req.TransferType != "" {
			if req.TransferType == "all_deposit" {
				query = query.Where(r.db.Where("transactions.transfer_type = ?", "deposit").Or("transactions.transfer_type = ?", "bonus"))
			} else if req.TransferType == "all_withdraw" {
				query = query.Where(r.db.Where("transactions.transfer_type = ?", "withdraw").Or("transactions.transfer_type = ?", "getcreditback"))
			} else {
				query = query.Where("transactions.transfer_type = ?", req.TransferType)
			}
		}
		if req.TransferStatus != "" {
			if req.TransferStatus == "failed" {
				query = query.Where(r.db.Where("transactions.status = ?", "canceled"))
			} else {
				query = query.Where("transactions.status = ?", req.TransferStatus)
			}
		}

		if req.Search != "" {
			search_like := fmt.Sprintf("%%%s%%", req.Search)
			query = query.Where(r.db.Where("transactions.from_account_name LIKE ?", search_like).Or("transactions.from_account_number LIKE ?", search_like).Or("transactions.to_account_name LIKE ?", search_like).Or("transactions.to_account_number LIKE ?", search_like))
		}

		// Sort by ANY //
		req.SortCol = strings.TrimSpace(req.SortCol)
		if req.SortCol != "" {
			if strings.ToLower(strings.TrimSpace(req.SortAsc)) == "desc" {
				req.SortAsc = "DESC"
			} else {
				req.SortAsc = "ASC"
			}
			query = query.Order(req.SortCol + " " + req.SortAsc)
		}

		if req.Limit > 0 {
			query = query.Limit(req.Limit)
		}
		if err = query.
			Where("transactions.deleted_at IS NULL").
			Offset(req.Page * req.Limit).
			Scan(&list).
			Error; err != nil {
			return nil, err
		}
	}

	// End count total records for pagination purposes (without limit and offset)
	var result model.SuccessWithPagination
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) CreateBankDepositTransaction(data model.BankTransactionCreateBody) (*int64, error) {

	if err := r.db.Table("Bank_transactions").Create(&data).Error; err != nil {
		return nil, err
	}
	return &data.Id, nil
}

func (r repo) CreateBankWithdrawTransactionWithCut(data model.BankTransactionCreateBody) (*int64, error) {

	member, err := r.GetMemberById(data.UserId)
	if err != nil {
		return nil, err
	}

	if data.CreditAmount > 0 && member.Credit >= data.CreditAmount {
		if err := r.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Table("Users").Where("id = ?", data.UserId).UpdateColumn("credit", gorm.Expr("credit - ?", data.CreditAmount)).Error; err != nil {
				return err
			}
			if err := tx.Table("Bank_transactions").Create(&data).Error; err != nil {
				return err
			}
			// if creditBalance, err := r.GetMemberCredit(data.UserId); err != nil {
			// 	return err
			// } else if creditBalance <= 0 {
			// 	return fmt.Errorf("ZERO_CREDIT")
			// }
			return nil
		}); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("INSUFFICIENT_CREDIT")
	}
	return &data.Id, nil
}

func (r repo) CreateBankWithdrawTransaction(data model.BankTransactionCreateBody) (*int64, error) {

	if err := r.db.Table("Bank_transactions").Create(&data).Error; err != nil {
		return nil, err
	}
	return &data.Id, nil
}

func (r repo) CreateBonusTransaction(data model.BonusTransactionCreateBody) error {

	if err := r.db.Table("Bank_transactions").Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) UpdateBankTransaction(id int64, data interface{}) error {

	if err := r.db.Table("Bank_transactions").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) DeleteBankTransaction(id int64) error {

	if err := r.db.Table("Bank_transactions").Where("id = ?", id).Delete(&model.BankTransaction{}).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) GetPendingDepositTransactions(req model.PendingDepositTransactionListRequest) (*model.SuccessWithPagination, error) {

	var list []model.BankTransactionResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset)
	count := r.db.Table("Bank_transactions as transactions")
	count = count.Select("transactions.id")
	count = count.Where("transactions.transfer_type = ?", "deposit")
	count = count.Where("transactions.status = ?", "pending")
	count = count.Where("transactions.removed_at IS NULL")
	if req.FromTransferDate != "" {
		count = count.Where("transactions.transfer_at >= ?", req.FromTransferDate)
	}
	if req.ToTransferDate != "" {
		count = count.Where("transactions.transfer_at <= ?", req.ToTransferDate)
	}
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where("transactions.from_account_name LIKE ?", search_like)
		count = count.Or("transactions.from_account_number LIKE ?", search_like)
		count = count.Or("transactions.to_account_name LIKE ?", search_like)
		count = count.Or("transactions.to_account_number LIKE ?", search_like)
	}

	if err = count.
		Where("transactions.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}
	if total > 0 {
		// SELECT //
		selectedFields := "transactions.id, transactions.user_id, transactions.transfer_type, transactions.promotion_id, transactions.from_account_id, transactions.from_bank_id, transactions.from_account_name, transactions.from_account_number, transactions.to_account_id, transactions.to_bank_id, transactions.to_account_name, transactions.to_account_number"
		selectedFields += ", transactions.credit_amount, transactions.paid_amount, transactions.over_amount, transactions.deposit_channel, transactions.bonus_amount, transactions.bonus_reason, transactions.before_amount, transactions.after_amount, transactions.bank_charge_amount"
		selectedFields += ", transactions.transfer_at, transactions.created_by_user_id, transactions.created_by_username, transactions.removed_at, transactions.removed_by_user_id, transactions.removed_by_username, transactions.status, transactions.status_detail, transactions.is_auto_credit"
		selectedFields += ", transactions.created_at, transactions.updated_at"
		selectedFields += ", from_banks.name as from_bank_name, from_banks.code as from_bank_code, from_banks.icon_url as from_bank_icon_url, from_banks.type_flag as from_bank_type_flag"
		selectedFields += ", to_banks.name as to_bank_name, to_banks.code as to_bank_code, to_banks.icon_url as to_bank_icon_url, to_banks.type_flag as to_bank_type_flag"
		selectedFields += ", users.member_code as member_code, users.username as user_username, users.firstname as user_firstname, users.lastname as user_lastname, users.fullname as user_fullname, users.phone as user_phone"
		query := r.db.Table("Bank_transactions as transactions")
		query = query.Select(selectedFields)
		query = query.Joins("LEFT JOIN Banks as from_banks ON from_banks.id = transactions.from_bank_id")
		query = query.Joins("LEFT JOIN Banks as to_banks ON to_banks.id = transactions.to_bank_id")
		query = query.Joins("LEFT JOIN Users as users ON users.id = transactions.user_id")
		query = query.Where("transactions.transfer_type = ?", "deposit")
		query = query.Where("transactions.status = ?", "pending")
		query = query.Where("transactions.removed_at IS NULL")
		if req.FromTransferDate != "" {
			query = query.Where("transactions.transfer_at >= ?", req.FromTransferDate)
		}
		if req.ToTransferDate != "" {
			query = query.Where("transactions.transfer_at <= ?", req.ToTransferDate)
		}
		if req.Search != "" {
			search_like := fmt.Sprintf("%%%s%%", req.Search)
			query = query.Where("transactions.from_account_name LIKE ?", search_like)
			query = query.Or("transactions.from_account_number LIKE ?", search_like)
			query = query.Or("transactions.to_account_name LIKE ?", search_like)
			query = query.Or("transactions.to_account_number LIKE ?", search_like)
		}

		// Sort by ANY //
		req.SortCol = strings.TrimSpace(req.SortCol)
		if req.SortCol != "" {
			if strings.ToLower(strings.TrimSpace(req.SortAsc)) == "desc" {
				req.SortAsc = "DESC"
			} else {
				req.SortAsc = "ASC"
			}
			query = query.Order(req.SortCol + " " + req.SortAsc)
		}
		if req.Limit > 0 {
			query = query.Limit(req.Limit)
		}
		if err = query.
			Where("transactions.deleted_at IS NULL").
			Offset(req.Page * req.Limit).
			Scan(&list).
			Error; err != nil {
			return nil, err
		}
	}

	// End count total records for pagination purposes (without limit and offset)
	var result model.SuccessWithPagination
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) GetPendingWithdrawTransactions(req model.PendingWithdrawTransactionListRequest) (*model.SuccessWithPagination, error) {

	var list []model.BankTransactionResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset)
	count := r.db.Table("Bank_transactions as transactions")
	count = count.Select("transactions.id")
	count = count.Where("transactions.transfer_type = ?", "withdraw")
	count = count.Where("transactions.status = ?", "pending")
	count = count.Where("transactions.removed_at IS NULL")
	if req.FromTransferDate != "" {
		count = count.Where("transactions.transfer_at >= ?", req.FromTransferDate)
	}
	if req.ToTransferDate != "" {
		count = count.Where("transactions.transfer_at <= ?", req.ToTransferDate)
	}
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where("transactions.from_account_name LIKE ?", search_like)
		count = count.Or("transactions.from_account_number LIKE ?", search_like)
		count = count.Or("transactions.to_account_name LIKE ?", search_like)
		count = count.Or("transactions.to_account_number LIKE ?", search_like)
	}

	if err = count.
		Where("transactions.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}
	if total > 0 {
		// SELECT //
		selectedFields := "transactions.id, transactions.user_id, transactions.transfer_type, transactions.promotion_id, transactions.from_account_id, transactions.from_bank_id, transactions.from_account_name, transactions.from_account_number, transactions.to_account_id, transactions.to_bank_id, transactions.to_account_name, transactions.to_account_number"
		selectedFields += ", transactions.credit_amount, transactions.paid_amount, transactions.over_amount, transactions.deposit_channel, transactions.bonus_amount, transactions.bonus_reason, transactions.before_amount, transactions.after_amount, transactions.bank_charge_amount"
		selectedFields += ", transactions.transfer_at, transactions.created_by_user_id, transactions.created_by_username, transactions.removed_at, transactions.removed_by_user_id, transactions.removed_by_username, transactions.status, transactions.status_detail, transactions.is_auto_credit"
		selectedFields += ", transactions.created_at, transactions.updated_at"
		selectedFields += ", from_banks.name as from_bank_name, from_banks.code as from_bank_code, from_banks.icon_url as from_bank_icon_url, from_banks.type_flag as from_bank_type_flag"
		selectedFields += ", to_banks.name as to_bank_name, to_banks.code as to_bank_code, to_banks.icon_url as to_bank_icon_url, to_banks.type_flag as to_bank_type_flag"
		selectedFields += ", users.member_code as member_code, users.username as user_username, users.firstname as user_firstname, users.lastname as user_lastname, users.fullname as user_fullname, users.phone as user_phone"
		query := r.db.Table("Bank_transactions as transactions")
		query = query.Select(selectedFields)
		query = query.Joins("LEFT JOIN Banks as from_banks ON from_banks.id = transactions.from_bank_id")
		query = query.Joins("LEFT JOIN Banks as to_banks ON to_banks.id = transactions.to_bank_id")
		query = query.Joins("LEFT JOIN Users as users ON users.id = transactions.user_id")
		query = query.Where("transactions.transfer_type = ?", "withdraw")
		query = query.Where("transactions.status = ?", "pending")
		query = query.Where("transactions.removed_at IS NULL")
		if req.FromTransferDate != "" {
			query = query.Where("transactions.transfer_at >= ?", req.FromTransferDate)
		}
		if req.ToTransferDate != "" {
			query = query.Where("transactions.transfer_at <= ?", req.ToTransferDate)
		}
		if req.Search != "" {
			search_like := fmt.Sprintf("%%%s%%", req.Search)
			query = query.Where("transactions.from_account_name LIKE ?", search_like)
			query = query.Or("transactions.from_account_number LIKE ?", search_like)
			query = query.Or("transactions.to_account_name LIKE ?", search_like)
			query = query.Or("transactions.to_account_number LIKE ?", search_like)
		}

		// Sort by ANY //
		req.SortCol = strings.TrimSpace(req.SortCol)
		if req.SortCol != "" {
			if strings.ToLower(strings.TrimSpace(req.SortAsc)) == "desc" {
				req.SortAsc = "DESC"
			} else {
				req.SortAsc = "ASC"
			}
			query = query.Order(req.SortCol + " " + req.SortAsc)
		}

		if req.Limit > 0 {
			query = query.Limit(req.Limit)
		}
		if err = query.
			Where("transactions.deleted_at IS NULL").
			Offset(req.Page * req.Limit).
			Scan(&list).
			Error; err != nil {
			return nil, err
		}
	}

	// End count total records for pagination purposes (without limit and offset)
	var result model.SuccessWithPagination
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) CancelPendingTransaction(id int64, data model.BankTransactionCancelBody) error {
	//.Where("status = ?", "pending")
	if err := r.db.Table("Bank_transactions").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) CreateTransactionAction(data model.CreateBankTransactionActionBody) (*int64, error) {
	if err := r.db.Table("Bank_confirm_transactions").Create(&data).Error; err != nil {
		return nil, err
	}
	return &data.Id, nil
}

func (r repo) RollbackTransactionAction(actionId int64) error {
	data := map[string]interface{}{
		"action_key": fmt.Sprintf("ROLLBACK#%d", actionId),
		"deleted_at": time.Now(),
	}
	if err := r.db.Table("Bank_confirm_transactions").Where("id = ?", actionId).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) CreateStatementAction(data model.CreateBankStatementActionBody) error {
	if err := r.db.Table("Bank_confirm_statements").Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) ConfirmPendingDepositTransaction(id int64, body model.BankDepositTransactionConfirmBody) error {
	// todo :
	// data := map[string]interface{}{
	// 	"transfer_at":           body.TransferAt,
	// 	"bonus_amount":          body.BonusAmount,
	// 	"status":                body.Status,
	// 	"confirmed_at":          body.ConfirmedAt,
	// 	"confirmed_by_user_id":  body.ConfirmedByUserId,
	// 	"confirmed_by_username": body.ConfirmedByUsername,
	// }
	// .Where("status = ?", "pending")
	if err := r.db.Table("Bank_transactions").Where("id = ?", id).Updates(body).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) ConfirmPendingCreditDepositTransaction(id int64, body model.BankDepositTransactionConfirmBody) error {
	// todo :
	// data := map[string]interface{}{
	// 	"transfer_at":           body.TransferAt,
	// 	"bonus_amount":          body.BonusAmount,
	// 	"status":                body.Status,
	// 	"confirmed_at":          body.ConfirmedAt,
	// 	"confirmed_by_user_id":  body.ConfirmedByUserId,
	// 	"confirmed_by_username": body.ConfirmedByUsername,
	// }
	if err := r.db.Table("Bank_transactions").Where("id = ?", id).Where("status = ?", "pending_credit").Updates(body).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) CheckMemeberHasEnoughtCredit(memberId int64, creditAmount float64) error {

	member, err := r.GetMemberById(memberId)
	if err != nil {
		return err
	}
	if creditAmount <= 0 {
		return fmt.Errorf("INVALID_CREDIT_AMOUNT")
	}
	if member.Credit < creditAmount {
		return fmt.Errorf("INSUFFICIENT_CREDIT")
	}
	return nil
}

func (r repo) ConfirmPendingWithdrawTransaction(id int64, body model.BankWithdrawTransactionConfirmBody) error {
	// todo :
	//  data := map[string]interface{}{
	// 	"transfer_at":           body.TransferAt,
	// 	"credit_amount":         body.CreditAmount,
	// 	"bank_charge_amount":    body.BankChargeAmount,
	// 	"status":                body.Status,
	// 	"confirmed_at":          body.ConfirmedAt,
	// 	"confirmed_by_user_id":  body.ConfirmedByUserId,
	// 	"confirmed_by_username": body.ConfirmedByUsername,
	// }
	//.Where("status = ?", "pending")
	if err := r.db.Table("Bank_transactions").Where("id = ?", id).Updates(&body).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) ConfirmPendingWithdrawTransfer(id int64, body model.BankWithdrawTransactionConfirmBody) error {
	// todo :
	// data := map[string]interface{}{
	// 	"transfer_at":           body.TransferAt,
	// 	"bank_charge_amount":    body.BankChargeAmount,
	// 	"status":                body.Status,
	// 	"confirmed_at":          body.ConfirmedAt,
	// 	"confirmed_by_user_id":  body.ConfirmedByUserId,
	// 	"confirmed_by_username": body.ConfirmedByUsername,
	// }
	if err := r.db.Table("Bank_transactions").Where("id = ?", id).Where("status = ?", "pending_transfer").Updates(&body).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) GetFinishedTransactions(req model.FinishedTransactionListRequest) (*model.SuccessWithPagination, error) {

	var list []model.BankTransactionResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset)
	count := r.db.Table("Bank_transactions as transactions")
	count = count.Select("transactions.id")
	count = count.Where("transactions.status = ?", "finished")
	count = count.Where("transactions.removed_at IS NULL")
	if req.AccountId != "" {
		count = count.Where(r.db.Where("transactions.from_account_id = ?", req.AccountId).Or("transactions.to_account_id = ?", req.AccountId))
	}
	if req.FromTransferDate != "" {
		count = count.Where("transactions.transfer_at >= ?", req.FromTransferDate)
	}
	if req.ToTransferDate != "" {
		count = count.Where("transactions.transfer_at <= ?", req.ToTransferDate)
	}
	if req.TransferType != "" {
		count = count.Where("transactions.transfer_type = ?", req.TransferType)
	}
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where("transactions.from_account_name LIKE ?", search_like)
		count = count.Or("transactions.from_account_number LIKE ?", search_like)
		count = count.Or("transactions.to_account_name LIKE ?", search_like)
		count = count.Or("transactions.to_account_number LIKE ?", search_like)
	}

	if err = count.
		Where("transactions.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}
	if total > 0 {
		// SELECT //
		selectedFields := "transactions.id, transactions.user_id, transactions.transfer_type, transactions.promotion_id, transactions.from_account_id, transactions.from_bank_id, transactions.from_account_name, transactions.from_account_number, transactions.to_account_id, transactions.to_bank_id, transactions.to_account_name, transactions.to_account_number"
		selectedFields += ", transactions.credit_amount, transactions.paid_amount, transactions.over_amount, transactions.deposit_channel, transactions.bonus_amount, transactions.bonus_reason, transactions.before_amount, transactions.after_amount, transactions.bank_charge_amount"
		selectedFields += ", transactions.transfer_at, transactions.created_by_user_id, transactions.created_by_username, transactions.removed_at, transactions.removed_by_user_id, transactions.removed_by_username, transactions.status, transactions.status_detail, transactions.is_auto_credit"
		selectedFields += ", transactions.created_at, transactions.updated_at"
		selectedFields += ", from_banks.name as from_bank_name, from_banks.code as from_bank_code, from_banks.icon_url as from_bank_icon_url, from_banks.type_flag as from_bank_type_flag"
		selectedFields += ", to_banks.name as to_bank_name, to_banks.code as to_bank_code, to_banks.icon_url as to_bank_icon_url, to_banks.type_flag as to_bank_type_flag"
		selectedFields += ", users.member_code as member_code, users.username as user_username, users.firstname as user_firstname, users.lastname as user_lastname, users.fullname as user_fullname, users.phone as user_phone"
		query := r.db.Table("Bank_transactions as transactions")
		query = query.Select(selectedFields)
		query = query.Joins("LEFT JOIN Banks as from_banks ON from_banks.id = transactions.from_bank_id")
		query = query.Joins("LEFT JOIN Banks as to_banks ON to_banks.id = transactions.to_bank_id")
		query = query.Joins("LEFT JOIN Users as users ON users.id = transactions.user_id")
		query = query.Where("transactions.status = ?", "finished")
		query = query.Where("transactions.removed_at IS NULL")
		if req.AccountId != "" {
			query = query.Where(r.db.Where("transactions.from_account_id = ?", req.AccountId).Or("transactions.to_account_id = ?", req.AccountId))
		}
		if req.FromTransferDate != "" {
			query = query.Where("transactions.transfer_at >= ?", req.FromTransferDate)
		}
		if req.ToTransferDate != "" {
			query = query.Where("transactions.transfer_at <= ?", req.ToTransferDate)
		}
		if req.TransferType != "" {
			query = query.Where("transactions.transfer_type = ?", req.TransferType)
		}
		if req.Search != "" {
			search_like := fmt.Sprintf("%%%s%%", req.Search)
			query = query.Where("transactions.from_account_name LIKE ?", search_like)
			query = query.Or("transactions.from_account_number LIKE ?", search_like)
			query = query.Or("transactions.to_account_name LIKE ?", search_like)
			query = query.Or("transactions.to_account_number LIKE ?", search_like)
		}

		// Sort by ANY //
		req.SortCol = strings.TrimSpace(req.SortCol)
		if req.SortCol != "" {
			if strings.ToLower(strings.TrimSpace(req.SortAsc)) == "desc" {
				req.SortAsc = "DESC"
			} else {
				req.SortAsc = "ASC"
			}
			query = query.Order(req.SortCol + " " + req.SortAsc)
		}

		if req.Limit > 0 {
			query = query.Limit(req.Limit)
		}
		if err = query.
			Where("transactions.deleted_at IS NULL").
			Offset(req.Page * req.Limit).
			Scan(&list).
			Error; err != nil {
			return nil, err
		}
	}

	// End count total records for pagination purposes (without limit and offset)
	var result model.SuccessWithPagination
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) RemoveFinishedTransaction(id int64, data model.BankTransactionRemoveBody) error {
	if err := r.db.Table("Bank_transactions").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) GetRemovedTransactions(req model.RemovedTransactionListRequest) (*model.SuccessWithPagination, error) {

	var list []model.BankTransactionResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset)
	count := r.db.Table("Bank_transactions as transactions")
	count = count.Select("transactions.id")
	count = count.Where("transactions.removed_at IS NOT NULL")
	if req.FromTransferDate != "" {
		count = count.Where("transactions.transfer_at >= ?", req.FromTransferDate)
	}
	if req.ToTransferDate != "" {
		count = count.Where("transactions.transfer_at <= ?", req.ToTransferDate)
	}
	if req.TransferType != "" {
		count = count.Where("transactions.transfer_type = ?", req.TransferType)
	}
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where("transactions.from_account_name LIKE ?", search_like)
		count = count.Or("transactions.from_account_number LIKE ?", search_like)
		count = count.Or("transactions.to_account_name LIKE ?", search_like)
		count = count.Or("transactions.to_account_number LIKE ?", search_like)
	}

	if err = count.
		Where("transactions.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}
	if total > 0 {
		// SELECT //
		selectedFields := "transactions.id, transactions.user_id, transactions.transfer_type, transactions.promotion_id, transactions.from_account_id, transactions.from_bank_id, transactions.from_account_name, transactions.from_account_number, transactions.to_account_id, transactions.to_bank_id, transactions.to_account_name, transactions.to_account_number"
		selectedFields += ", transactions.credit_amount, transactions.paid_amount, transactions.over_amount, transactions.deposit_channel, transactions.bonus_amount, transactions.bonus_reason, transactions.before_amount, transactions.after_amount, transactions.bank_charge_amount"
		selectedFields += ", transactions.transfer_at, transactions.created_by_user_id, transactions.created_by_username, transactions.removed_at, transactions.removed_by_user_id, transactions.removed_by_username, transactions.status, transactions.status_detail, transactions.is_auto_credit"
		selectedFields += ", transactions.created_at, transactions.updated_at"
		selectedFields += ", from_banks.name as from_bank_name, from_banks.code as from_bank_code, from_banks.icon_url as from_bank_icon_url, from_banks.type_flag as from_bank_type_flag"
		selectedFields += ", to_banks.name as to_bank_name, to_banks.code as to_bank_code, to_banks.icon_url as to_bank_icon_url, to_banks.type_flag as to_bank_type_flag"
		selectedFields += ", users.member_code as member_code, users.username as user_username, users.firstname as user_firstname, users.lastname as user_lastname, users.fullname as user_fullname, users.phone as user_phone"
		query := r.db.Table("Bank_transactions as transactions")
		query = query.Select(selectedFields)
		query = query.Joins("LEFT JOIN Banks as from_banks ON from_banks.id = transactions.from_bank_id")
		query = query.Joins("LEFT JOIN Banks as to_banks ON to_banks.id = transactions.to_bank_id")
		query = query.Joins("LEFT JOIN Users as users ON users.id = transactions.user_id")
		query = query.Where("transactions.removed_at IS NOT NULL")
		if req.FromTransferDate != "" {
			query = query.Where("transactions.transfer_at >= ?", req.FromTransferDate)
		}
		if req.ToTransferDate != "" {
			query = query.Where("transactions.transfer_at <= ?", req.ToTransferDate)
		}
		if req.TransferType != "" {
			query = query.Where("transactions.transfer_type = ?", req.TransferType)
		}
		if req.Search != "" {
			search_like := fmt.Sprintf("%%%s%%", req.Search)
			query = query.Where("transactions.from_account_name LIKE ?", search_like)
			query = query.Or("transactions.from_account_number LIKE ?", search_like)
			query = query.Or("transactions.to_account_name LIKE ?", search_like)
			query = query.Or("transactions.to_account_number LIKE ?", search_like)
		}

		// Sort by ANY //
		req.SortCol = strings.TrimSpace(req.SortCol)
		if req.SortCol != "" {
			if strings.ToLower(strings.TrimSpace(req.SortAsc)) == "desc" {
				req.SortAsc = "DESC"
			} else {
				req.SortAsc = "ASC"
			}
			query = query.Order(req.SortCol + " " + req.SortAsc)
		}

		if req.Limit > 0 {
			query = query.Limit(req.Limit)
		}
		if err = query.
			Where("transactions.deleted_at IS NULL").
			Offset(req.Page * req.Limit).
			Scan(&list).
			Error; err != nil {
			return nil, err
		}
	}

	// End count total records for pagination purposes (without limit and offset)
	var result model.SuccessWithPagination
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) GetMemberById(id int64) (*model.Member, error) {
	var record model.Member

	selectedFields := "users.id, users.member_code, users.username, users.phone, users.firstname, users.lastname, users.fullname, users.credit, users.bankname, users.bank_account, users.promotion, users.status, users.channel, users.true_wallet, users.note, users.turnover_limit, users.created_at"
	if err := r.db.Table("Users as users").
		Select(selectedFields).
		Where("users.id = ?", id).
		Where("users.deleted_at IS NULL").
		First(&record).
		Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r repo) GetMemberCredit(id int64) (float64, error) {
	var record model.Member

	selectedFields := "users.id, users.credit"
	if err := r.db.Table("Users as users").
		Select(selectedFields).
		Where("users.id = ?", id).
		Where("users.deleted_at IS NULL").
		First(&record).
		Error; err != nil {
		return 0, err
	}
	return record.Credit, nil
}

func (r repo) GetMemberByCode(memberCode string) (*model.Member, error) {
	var record model.Member

	selectedFields := "users.id, users.member_code, users.username, users.phone, users.firstname, users.lastname, users.fullname, users.credit, users.bankname, users.bank_account, users.promotion, users.status, users.channel, users.true_wallet, users.note, users.turnover_limit, users.created_at"
	if err := r.db.Table("Users as users").
		Select(selectedFields).
		Where("users.member_code = ?", memberCode).
		Where("users.deleted_at IS NULL").
		First(&record).
		Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r repo) GetMembers(req model.MemberListRequest) (*model.SuccessWithPagination, error) {

	var list []model.Member
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset)
	count := r.db.Table("Users as users")
	count = count.Select("users.id")
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where(r.db.Where("users.username LIKE ?", search_like).Or("users.phone LIKE ?", search_like).Or("users.fullname LIKE ?", search_like).Or("users.bankname LIKE ?", search_like).Or("users.bank_account LIKE ?", search_like))
	}

	if err = count.
		Where("users.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}
	if total > 0 {
		// SELECT //
		selectedFields := "users.id, users.member_code, users.username, users.phone, users.firstname, users.lastname, users.fullname, users.credit, users.bankname, users.bank_account, users.promotion, users.status, users.channel, users.true_wallet, users.note, users.turnover_limit, users.created_at"
		query := r.db.Table("Users as users")
		query = query.Select(selectedFields)
		if req.Search != "" {
			search_like := fmt.Sprintf("%%%s%%", req.Search)
			query = query.Where(r.db.Where("users.username LIKE ?", search_like).Or("users.phone LIKE ?", search_like).Or("users.fullname LIKE ?", search_like).Or("users.bankname LIKE ?", search_like).Or("users.bank_account LIKE ?", search_like))
		}

		// Sort by ANY //
		req.SortCol = strings.TrimSpace(req.SortCol)
		if req.SortCol != "" {
			if strings.ToLower(strings.TrimSpace(req.SortAsc)) == "desc" {
				req.SortAsc = "DESC"
			} else {
				req.SortAsc = "ASC"
			}
			query = query.Order(req.SortCol + " " + req.SortAsc)
		}

		if req.Limit > 0 {
			query = query.Limit(req.Limit)
		}
		if err = query.
			Where("users.deleted_at IS NULL").
			Offset(req.Page * req.Limit).
			Scan(&list).
			Error; err != nil {
			return nil, err
		}
	}

	// End count total records for pagination purposes (without limit and offset)
	var result model.SuccessWithPagination
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) GetPossibleStatementOwners(req model.MemberPossibleListRequest) (*model.SuccessWithPagination, error) {

	var list []model.Member
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset)
	count := r.db.Table("Users as users")
	count = count.Select("users.id")
	if req.UserBankCode != nil {
		count = count.Where("users.bank_code = ?", *req.UserBankCode)
	}
	if req.UserAccountNumber != nil && *req.UserAccountNumber != "" {
		search_like := fmt.Sprintf("%%%s%%", *req.UserAccountNumber)
		count = count.Where("users.bank_account LIKE ?", search_like)
	}

	if err = count.
		Where("users.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}
	if total > 0 {
		// SELECT //
		selectedFields := "users.id, users.member_code, users.username, users.phone, users.firstname, users.lastname, users.fullname, users.credit, users.bankname, users.bank_account, users.promotion, users.status, users.channel, users.true_wallet, users.note, users.turnover_limit, users.created_at"
		query := r.db.Table("Users as users")
		query = query.Select(selectedFields)
		if req.UserBankCode != nil {
			query = query.Where("users.bank_code = ?", req.UserBankCode)
		}
		if req.UserAccountNumber != nil && *req.UserAccountNumber != "" {
			search_like := fmt.Sprintf("%%%s%%", *req.UserAccountNumber)
			query = query.Where("users.bank_account LIKE ?", search_like)
		}

		// Sort by ANY //
		req.SortCol = strings.TrimSpace(req.SortCol)
		if req.SortCol != "" {
			if strings.ToLower(strings.TrimSpace(req.SortAsc)) == "desc" {
				req.SortAsc = "DESC"
			} else {
				req.SortAsc = "ASC"
			}
			query = query.Order(req.SortCol + " " + req.SortAsc)
		}

		if req.Limit > 0 {
			query = query.Limit(req.Limit)
		}
		if err = query.
			Where("users.deleted_at IS NULL").
			Offset(req.Page * req.Limit).
			Scan(&list).
			Error; err != nil {
			return nil, err
		}
	}

	// End count total records for pagination purposes (without limit and offset)
	var result model.SuccessWithPagination
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) GetMemberTransactions(req model.MemberTransactionListRequest) (*model.SuccessWithPagination, error) {

	var list []model.MemberTransaction
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset)
	count := r.db.Table("Bank_transactions as transactions")
	count = count.Select("transactions.id")
	count = count.Where("transactions.removed_at IS NULL")
	if req.UserId != "" {
		count = count.Where("transactions.user_id = ?", req.UserId)
	}
	if req.FromTransferDate != "" {
		count = count.Where("transactions.transfer_at >= ?", req.FromTransferDate)
	}
	if req.ToTransferDate != "" {
		count = count.Where("transactions.transfer_at <= ?", req.ToTransferDate)
	}
	if req.TransferType != "" {
		count = count.Where("transactions.transfer_type = ?", req.TransferType)
	}
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where(r.db.Where("transactions.member_code LIKE ?", search_like).Or("transactions.from_account_name LIKE ?", search_like).Or("transactions.from_account_number LIKE ?", search_like).Or("transactions.to_account_name LIKE ?", search_like).Or("transactions.to_account_number LIKE ?", search_like))
	}

	if err = count.
		Where("transactions.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}
	if total > 0 {
		// SELECT //
		selectedFields := "transactions.id, transactions.user_id, transactions.transfer_type, transactions.promotion_id, transactions.from_account_id, transactions.from_bank_id, transactions.from_account_name, transactions.from_account_number, transactions.to_account_id, transactions.to_bank_id, transactions.to_account_name, transactions.to_account_number"
		selectedFields += ", transactions.credit_amount, transactions.paid_amount, transactions.over_amount, transactions.deposit_channel, transactions.bonus_amount, transactions.bonus_reason, transactions.before_amount, transactions.after_amount, transactions.bank_charge_amount"
		selectedFields += ", transactions.transfer_at, transactions.created_by_user_id, transactions.created_by_username, transactions.removed_at, transactions.removed_by_user_id, transactions.removed_by_username, transactions.status, transactions.status_detail, transactions.is_auto_credit"
		selectedFields += ", transactions.created_at, transactions.updated_at"
		selectedFields += ", from_banks.name as from_bank_name, from_banks.code as from_bank_code, from_banks.icon_url as from_bank_icon_url, from_banks.type_flag as from_bank_type_flag"
		selectedFields += ", to_banks.name as to_bank_name, to_banks.code as to_bank_code, to_banks.icon_url as to_bank_icon_url, to_banks.type_flag as to_bank_type_flag"
		selectedFields += ", users.member_code as member_code, users.username as user_username, users.firstname as user_firstname, users.lastname as user_lastname, users.fullname as user_fullname, users.phone as user_phone"
		query := r.db.Table("Bank_transactions as transactions")
		query = query.Select(selectedFields)
		query = query.Joins("LEFT JOIN Banks as from_banks ON from_banks.id = transactions.from_bank_id")
		query = query.Joins("LEFT JOIN Banks as to_banks ON to_banks.id = transactions.to_bank_id")
		query = query.Joins("LEFT JOIN Users as users ON users.id = transactions.user_id")
		query = query.Where("transactions.removed_at IS NULL")
		if req.UserId != "" {
			query = query.Where("transactions.user_id = ?", req.UserId)
		}
		if req.FromTransferDate != "" {
			query = query.Where("transactions.transfer_at >= ?", req.FromTransferDate)
		}
		if req.ToTransferDate != "" {
			query = query.Where("transactions.transfer_at <= ?", req.ToTransferDate)
		}
		if req.TransferType != "" {
			query = query.Where("transactions.transfer_type = ?", req.TransferType)
		}
		if req.Search != "" {
			search_like := fmt.Sprintf("%%%s%%", req.Search)
			query = query.Where(r.db.Where("transactions.member_code LIKE ?", search_like).Or("transactions.from_account_name LIKE ?", search_like).Or("transactions.from_account_number LIKE ?", search_like).Or("transactions.to_account_name LIKE ?", search_like).Or("transactions.to_account_number LIKE ?", search_like))
		}

		// Sort by ANY //
		req.SortCol = strings.TrimSpace(req.SortCol)
		if req.SortCol != "" {
			if strings.ToLower(strings.TrimSpace(req.SortAsc)) == "desc" {
				req.SortAsc = "DESC"
			} else {
				req.SortAsc = "ASC"
			}
			query = query.Order(req.SortCol + " " + req.SortAsc)
		}

		if req.Limit > 0 {
			query = query.Limit(req.Limit)
		}
		if err = query.
			Where("transactions.deleted_at IS NULL").
			Offset(req.Page * req.Limit).
			Scan(&list).
			Error; err != nil {
			return nil, err
		}
	}

	// End count total records for pagination purposes (without limit and offset)
	var result model.SuccessWithPagination
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) GetMemberTransactionSummary(req model.MemberTransactionListRequest) (*model.MemberTransactionSummary, error) {

	var result model.MemberTransactionSummary
	var err error

	// SELECT //
	selectedFields := "SUM(case when transfer_type = 'deposit' then credit_amount else 0 end) as total_deposit_amount, SUM(case when transfer_type = 'withdraw' then credit_amount else 0 end) as total_withdraw_amount, SUM(bonus_amount) as total_bom_amount"
	query := r.db.Table("Bank_transactions as transactions")
	query = query.Select(selectedFields)
	query = query.Joins("LEFT JOIN Banks as from_banks ON from_banks.id = transactions.from_bank_id")
	query = query.Joins("LEFT JOIN Banks as to_banks ON to_banks.id = transactions.to_bank_id")
	query = query.Where("transactions.removed_at IS NULL")
	if req.UserId != "" {
		query = query.Where("transactions.user_id = ?", req.UserId)
	}
	if req.FromTransferDate != "" {
		query = query.Where("transactions.transfer_at >= ?", req.FromTransferDate)
	}
	if req.ToTransferDate != "" {
		query = query.Where("transactions.transfer_at <= ?", req.ToTransferDate)
	}

	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		query = query.Where(r.db.Where("transactions.member_code LIKE ?", search_like).Or("transactions.from_account_name LIKE ?", search_like).Or("transactions.from_account_number LIKE ?", search_like).Or("transactions.to_account_name LIKE ?", search_like).Or("transactions.to_account_number LIKE ?", search_like))
	}

	if err = query.
		Where("transactions.deleted_at IS NULL").
		Scan(&result).
		Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r repo) IncreaseMemberCredit(body model.MemberStatementCreateBody) error {

	member, err := r.GetMemberById(body.UserId)
	if err != nil {
		return err
	}

	// todo : use agent credit
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		data := map[string]interface{}{
			"user_id":           member.Id,
			"statement_type_id": body.StatementTypeId,
			"transfer_at":       time.Now(),
			"info":              body.Info,
			"before_balance":    member.Credit,
			"amount":            body.Amount,
			"after_balance":     member.Credit + body.Amount,
		}
		if err := r.db.Table("User_statements").Create(&data).Error; err != nil {
			return err
		}
		if err := tx.Table("Users").Where("id = ?", member.Id).UpdateColumn("credit", gorm.Expr("credit + ?", body.Amount)).Error; err != nil {
			return err
		}
		// todo : sync with agent credit
		// var agentCreditBalance = from agent
		// if err := tx.Table("Users").Where("id = ?", member.Id).UpdateColumn("credit", agentCreditBalance).Error; err != nil {
		// 	return err
		// }
		return nil // COMMIT
	}); err != nil {
		return err
	}
	return nil
}

func (r repo) DecreaseMemberCredit(body model.MemberStatementCreateBody) error {

	member, err := r.GetMemberById(body.UserId)
	if err != nil {
		return err
	}
	// todo : check with agent credit
	if body.Amount > 0 && member.Credit >= body.Amount {
		// todo : use agent credit
		if err := r.db.Transaction(func(tx *gorm.DB) error {
			data := map[string]interface{}{
				"user_id":           member.Id,
				"statement_type_id": body.StatementTypeId,
				"transfer_at":       time.Now(),
				"info":              body.Info,
				"before_balance":    member.Credit,
				"amount":            body.Amount * -1,
				"after_balance":     member.Credit + body.Amount,
			}
			if err := r.db.Table("User_statements").Create(&data).Error; err != nil {
				return err
			}
			if err := tx.Table("Users").Where("id = ?", member.Id).UpdateColumn("credit", gorm.Expr("credit - ?", body.Amount)).Error; err != nil {
				return err
			}
			// todo : sync with agent credit
			// var agentCreditBalance = from agent
			// if err := tx.Table("Users").Where("id = ?", member.Id).UpdateColumn("credit", agentCreditBalance).Error; err != nil {
			// 	return err
			// }
			return nil // COMMIT
		}); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("NOT_ENOUGH_CREDIT")
	}
	return nil
}

func (r repo) TransferExternalAccount(body model.ExternalAccountTransferBody) error {

	client := &http.Client{}
	// curl -X POST "https://api.fastbankapi.com/api/v2/statement/transfer" -H "accept: */*" -H "apiKey: xxxxxxxxxx.yyyyyyyyyyy"
	//-H "Content-Type: application/json" -d "{ \"accountFrom\": \"aaaaaaaaaaaaaaaa\", \"accountTo\": \"bbbbbbbbbbbbbb\", \"amount\": \"8\", \"bankCode\": \"bay\", \"pin\": \"ccccc\"}"
	data, _ := json.Marshal(body)
	reqHttp, _ := http.NewRequest("POST", os.Getenv("ACCOUNTING_API_ENDPOINT")+"/api/v2/statement/transfer", bytes.NewBuffer(data))
	reqHttp.Header.Set("apiKey", os.Getenv("ACCOUNTING_API_KEY"))
	reqHttp.Header.Set("Content-Type", "application/json")
	response, err := client.Do(reqHttp)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("response", string(responseData))

	if response.StatusCode != 200 {
		var errorModel model.ExternalAccountError
		errJson := json.Unmarshal(responseData, &errorModel)
		if errJson != nil {
			return fmt.Errorf("EXTERNAL_REPONSE_JSON_ERROR")
		}
		fmt.Println("errorModel", errorModel)
		if errorModel.Error != "" {
			return fmt.Errorf(errorModel.Error)
		}
		return fmt.Errorf("EXTERNAL_API_ERROR")
	}
	return nil
}

func (r repo) GetMemberStatementTypeByCode(code string) (*model.MemberStatementType, error) {
	var record model.MemberStatementType
	selectedFields := "types.id, types.code, types.name, types.created_at, types.updated_at"
	if err := r.db.Table("User_statement_types as types").
		Select(selectedFields).
		Where("types.code = ?", code).
		Where("types.deleted_at IS NULL").
		First(&record).
		Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r repo) GetMemberStatementTypes(req model.SimpleListRequest) (*model.SuccessWithPagination, error) {

	var list []model.MemberStatementType
	var total int64
	var err error

	// COUNT total records for pagination purposes (without limit and offset)
	count := r.db.Table("User_statement_types as types")
	count = count.Select("types.id")
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where(r.db.Where("types.name LIKE ?", search_like))
	}
	if err = count.
		Where("types.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	if total > 0 {
		// SELECT //
		selectedFields := "types.id, types.code, types.name, types.created_at, types.updated_at"
		query := r.db.Table("User_statement_types as types")
		query = query.Select(selectedFields)
		if req.Search != "" {
			search_like := fmt.Sprintf("%%%s%%", req.Search)
			query = query.Where(r.db.Where("types.name LIKE ?", search_like))
		}

		// Sort by ANY //
		req.SortCol = strings.TrimSpace(req.SortCol)
		if req.SortCol != "" {
			if strings.ToLower(strings.TrimSpace(req.SortAsc)) == "desc" {
				req.SortAsc = "DESC"
			} else {
				req.SortAsc = "ASC"
			}
			query = query.Order(req.SortCol + " " + req.SortAsc)
		}

		if req.Limit > 0 {
			query = query.Limit(req.Limit)
		}
		if err = query.
			Where("types.deleted_at IS NULL").
			Offset(req.Page * req.Limit).
			Scan(&list).
			Error; err != nil {
			return nil, err
		}
	}

	// End count total records for pagination purposes (without limit and offset)
	var result model.SuccessWithPagination
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) GetMemberStatementById(id int64) (*model.MemberStatementResponse, error) {
	var record model.MemberStatementResponse
	selectedFields := "statements.id, statements.user_id, statements.statement_type_id, statements.transfer_at, statements.Info, statements.before_balance, statements.amount, statements.after_balance, statements.created_at, statements.updated_at"
	selectedFields += ",statement_types.name as statement_type_name"
	selectedFields += ",users.member_code as member_code, users.username as user_username, users.fullname as user_fullname"
	if err := r.db.Table("User_statements as statements").
		Select(selectedFields).
		Joins("LEFT JOIN User_statement_types AS statement_types ON statement_types.id = statements.statement_type_id").
		Joins("LEFT JOIN Users AS users ON users.id = statements.user_id").
		Where("statements.id = ?", id).
		Where("statements.deleted_at IS NULL").
		First(&record).
		Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r repo) GetMemberStatements(req model.MemberStatementListRequest) (*model.SuccessWithPagination, error) {

	var list []model.MemberStatementResponse
	var total int64
	var err error

	// COUNT total records for pagination purposes (without limit and offset)
	count := r.db.Table("User_statements as statements")
	count = count.Joins("LEFT JOIN User_statement_types AS statement_types ON statement_types.id = statements.statement_type_id")
	count = count.Joins("LEFT JOIN Users AS users ON users.id = statements.user_id")
	count = count.Select("statements.id")
	if req.UserId != "" {
		count = count.Where("statements.user_id = ?", req.UserId)
	}
	if req.FromTransferDate != "" {
		count = count.Where("statements.transfer_at >= ?", req.FromTransferDate)
	}
	if req.ToTransferDate != "" {
		count = count.Where("statements.transfer_at <= ?", req.ToTransferDate)
	}
	if req.StatementTypeId != "" {
		count = count.Where("statements.statement_type_id = ?", req.StatementTypeId)
	}
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where(r.db.Where("users.member_code LIKE ?", search_like).Or("users.username LIKE ?", search_like).Or("users.fullname LIKE ?", search_like))
	}
	if err = count.
		Where("statements.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	if total > 0 {
		// SELECT //
		selectedFields := "statements.id, statements.user_id, statements.statement_type_id, statements.transfer_at, statements.Info, statements.before_balance, statements.amount, statements.after_balance, statements.created_at, statements.updated_at"
		selectedFields += ",statement_types.name as statement_type_name"
		selectedFields += ",users.member_code as member_code, users.username as user_username, users.fullname as user_fullname"
		query := r.db.Table("User_statements as statements")
		query = query.Select(selectedFields)
		query = query.Joins("LEFT JOIN User_statement_types AS statement_types ON statement_types.id = statements.statement_type_id")
		query = query.Joins("LEFT JOIN Users AS users ON users.id = statements.user_id")
		if req.UserId != "" {
			query = query.Where("statements.user_id = ?", req.UserId)
		}
		if req.FromTransferDate != "" {
			query = query.Where("statements.transfer_at >= ?", req.FromTransferDate)
		}
		if req.ToTransferDate != "" {
			query = query.Where("statements.transfer_at <= ?", req.ToTransferDate)
		}
		if req.StatementTypeId != "" {
			query = query.Where("statements.statement_type_id = ?", req.StatementTypeId)
		}
		if req.Search != "" {
			search_like := fmt.Sprintf("%%%s%%", req.Search)
			query = query.Where(r.db.Where("users.member_code LIKE ?", search_like).Or("users.username LIKE ?", search_like).Or("users.fullname LIKE ?", search_like))
		}

		// Sort by ANY //
		req.SortCol = strings.TrimSpace(req.SortCol)
		if req.SortCol != "" {
			if strings.ToLower(strings.TrimSpace(req.SortAsc)) == "desc" {
				req.SortAsc = "DESC"
			} else {
				req.SortAsc = "ASC"
			}
			query = query.Order(req.SortCol + " " + req.SortAsc)
		}
		if req.Limit > 0 {
			query = query.Limit(req.Limit)
		}
		if err = query.
			Where("statements.deleted_at IS NULL").
			Offset(req.Page * req.Limit).
			Scan(&list).
			Error; err != nil {
			return nil, err
		}
	}

	// End count total records for pagination purposes (without limit and offset)
	var result model.SuccessWithPagination
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) CreateMemberStatement(data model.MemberStatementCreateBody) (*int64, error) {
	if err := r.db.Table("User_statements").Create(&data).Error; err != nil {
		return nil, err
	}
	return &data.Id, nil
}
