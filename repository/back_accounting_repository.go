package repository

import (
	"cybergame-api/model"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func NewAccountingRepository(db *gorm.DB) AccountingRepository {
	return &repo{db}
}

type AccountingRepository interface {
	GetAdminById(id int64) (*model.Admin, error)

	GetBanks(req model.BankListRequest) (*model.SuccessWithPagination, error)
	GetBankById(id int64) (*model.Bank, error)
	GetBankByCode(code string) (*model.Bank, error)

	GetAccountTypes(req model.AccountTypeListRequest) (*model.SuccessWithPagination, error)
	GetAccounTypeById(id int64) (*model.AccountType, error)

	GetUserByMemberCode(memberCode string) (*model.User, error)

	HasBankAccount(accountNumber string) (bool, error)
	GetBankAccountById(id int64) (*model.BankAccount, error)
	GetBankAccountByAccountNumber(accountNumber string) (*model.BankAccount, error)
	GetBankAccountByExternalId(id int64) (*model.BankAccount, error)
	GetActiveExternalAccount() (*model.BankAccount, error)
	GetDepositAccountById(id int64) (*model.BankAccount, error)
	GetWithdrawAccountById(id int64) (*model.BankAccount, error)
	GetBankAccounts(data model.BankAccountListRequest) (*model.SuccessWithPagination, error)
	GetBankAccountPriorities() (*model.SuccessWithPagination, error)
	GetBotBankAccounts(data model.BankAccountListRequest) (*model.SuccessWithPagination, error)
	CreateBankAccount(data model.BankAccountCreateBody) error
	ResetMainWithdrawBankAccount() error
	UpdateBankAccount(id int64, data model.BankAccountUpdateBody) error
	DeleteBankAccount(id int64, data model.BankAccountDeleteBody) error

	GetTransactionById(id int64) (*model.BankAccountTransaction, error)
	GetTransactions(data model.BankAccountTransactionListRequest) (*model.SuccessWithPagination, error)
	CreateTransaction(data model.BankAccountTransactionBody) error
	UpdateTransaction(id int64, data model.BankAccountTransactionBody) error
	DeleteTransaction(id int64) error

	GetTransferById(id int64) (*model.BankAccountTransfer, error)
	GetTransfers(data model.BankAccountTransferListRequest) (*model.SuccessWithPagination, error)
	CreateTransfer(data model.BankAccountTransferBody) error
	ConfirmTransfer(id int64, data model.BankAccountTransferConfirmBody) error
	DeleteTransfer(id int64) error

	CreateWebhookLog(body model.WebhookLogCreateBody) (*int64, error)
	UpdateWebhookLog(id int64, body model.WebhookLogUpdateBody) error
	GetWebhookStatementByExternalId(id int64) (*model.BankStatement, error)
	CreateWebhookStatement(body model.BankStatementCreateBody) (*int64, error)

	GetBotaccountConfigByKey(req model.BotAccountConfigListRequest) (*model.BotAccountConfig, error)
	GetBotaccountConfigs(req model.BotAccountConfigListRequest) (*model.SuccessWithPagination, error)
	CreateBotaccountConfig(data model.BotAccountConfigCreateBody) error
	DeleteBotaccountConfigByKey(key string) error
	DeleteBotaccountConfigById(id int64) error

	// Banking REPO
	GetBankStatements(req model.BankStatementListRequest) (*model.SuccessWithPagination, error)
	GetBankExternalStatements(externalIds []int64) (*model.SuccessWithPagination, error)
	HasBankExternalStatements(externalId int64) error
	GetMemberById(id int64) (*model.Member, error)
	IncreaseMemberCredit(body model.MemberStatementCreateBody) error
	GetMemberStatementTypeByCode(code string) (*model.MemberStatementType, error)
	GetPossibleStatementOwners(req model.MemberPossibleListRequest) (*model.SuccessWithPagination, error)
	GetBankStatementById(id int64) (*model.BankStatement, error)
	CreateBankDepositTransaction(data model.BankTransactionCreateBody) (*int64, error)
	CreateBankWithdrawTransaction(data model.BankTransactionCreateBody) (*int64, error)
	UpdateBankTransaction(id int64, data interface{}) error
	GetBankTransactionById(id int64) (*model.BankTransaction, error)
	CreateTransactionAction(data model.CreateBankTransactionActionBody) (*int64, error)
	RollbackTransactionAction(id int64) error
	ConfirmPendingDepositTransaction(id int64, data model.BankDepositTransactionConfirmBody) error
	ConfirmPendingCreditDepositTransaction(id int64, data model.BankDepositTransactionConfirmBody) error
	ConfirmPendingWithdrawTransaction(id int64, data model.BankWithdrawTransactionConfirmBody) error
	CreateStatementAction(data model.CreateBankStatementActionBody) error
	UpdateBankStatement(id int64, data model.BankStatementUpdateBody) error
	// MatchStatementOwner(id int64, data model.BankStatementUpdateBody) error
	IgnoreStatementOwner(id int64, data model.BankStatementUpdateBody) error
}

func (r repo) GetAdminById(id int64) (*model.Admin, error) {
	var admin model.Admin

	if err := r.db.Table("Admins").
		Select("id, username, phone, password, email, role").
		Where("id = ?", id).
		First(&admin).
		Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r repo) GetBanks(req model.BankListRequest) (*model.SuccessWithPagination, error) {

	var list []model.BankResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("Banks")
	count = count.Select("id")
	if req.Search != "" {
		count = count.Where("code = ?", req.Search)
	}
	if err = count.
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	if total > 0 {
		// SELECT //
		query := r.db.Table("Banks")
		query = query.Select("id, name, code, icon_url, icon_url, type_flag")
		if req.Search != "" {
			query = query.Where("code = ?", req.Search)
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
			Offset(req.Page * req.Limit).
			Scan(&list).
			Error; err != nil {
			return nil, err
		}
	}

	// End count total records for pagination purposes (without limit and offset) //
	var result model.SuccessWithPagination
	if list == nil {
		list = []model.BankResponse{}
	}
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) GetBankById(id int64) (*model.Bank, error) {

	var result *model.Bank
	if err := r.db.Table("Banks").
		Select("id, name, code, icon_url, type_flag").
		Where("id = ?", id).
		First(&result).
		Error; err != nil {
		return nil, err
	}

	// if result.Id == 0 {
	// 	return nil, errors.New(bankNotFound)
	// }
	return result, nil
}

func (r repo) GetBankByCode(code string) (*model.Bank, error) {

	var result *model.Bank
	if err := r.db.Table("Banks").
		Select("id, name, code, icon_url, type_flag").
		Where("code = ?", code).
		First(&result).
		Error; err != nil {
		return nil, err
	}

	// if result.Id == 0 {
	// 	return nil, errors.New(bankNotFound)
	// }
	return result, nil
}

func (r repo) GetAccountTypes(req model.AccountTypeListRequest) (*model.SuccessWithPagination, error) {

	var list []model.AccountTypeResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("Bank_account_types")
	count = count.Select("id")
	if req.Search != "" {
		count = count.Where("name = ?", req.Search)
	}
	if err = count.
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	if total > 0 {
		// SELECT //
		query := r.db.Table("Bank_account_types")
		query = query.Select("id, name, limit_flag")
		if req.Search != "" {
			query = query.Where("name = ?", req.Search)
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
			Offset(req.Page * req.Limit).
			Scan(&list).
			Error; err != nil {
			return nil, err
		}
	}

	// End count total records for pagination purposes (without limit and offset) //
	var result model.SuccessWithPagination
	if list == nil {
		list = []model.AccountTypeResponse{}
	}
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) GetAccounTypeById(id int64) (*model.AccountType, error) {

	var result model.AccountType
	if err := r.db.Table("Bank_account_types").
		Select("id, name, limit_flag").
		Where("id = ?", id).
		First(&result).
		Error; err != nil {
		return nil, err
	}

	// if result.Id == 0 {
	// 	return nil, errors.New("Account type not found")
	// }
	return &result, nil
}

func (r repo) GetUserByMemberCode(memberCode string) (*model.User, error) {

	var result model.User
	if err := r.db.Table("Users").
		Select("id, member_code, fullname, username, bankname, bank_code, bank_account, credit").
		Where("member_code = ?", memberCode).
		First(&result).
		Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (r repo) HasBankAccount(accountNumber string) (bool, error) {
	var count int64
	if err := r.db.Table("Bank_accounts").
		Select("id").
		Where("account_number = ?", accountNumber).
		Where("deleted_at IS NULL").
		Limit(1).
		Count(&count).
		Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r repo) GetBankAccountById(id int64) (*model.BankAccount, error) {

	var accounting model.BankAccount
	selectedFields := "accounts.id, accounts.bank_id, accounts.account_type_id, accounts.account_name, accounts.account_number, accounts.account_balance, accounts.account_priority_id, accounts.account_status, accounts.device_uid, accounts.pin_code, accounts.connection_status"
	selectedFields += ", accounts.auto_credit_flag, accounts.is_main_withdraw, accounts.auto_withdraw_flag, accounts.auto_withdraw_credit_flag, accounts.auto_withdraw_confirm_flag, accounts.auto_withdraw_max_amount, accounts.auto_transfer_max_amount, accounts.qr_wallet_status"
	selectedFields += ", accounts.last_conn_update_at, accounts.created_at, accounts.updated_at"
	selectedFields += ", banks.name as bank_name, banks.code as bank_code, banks.icon_url as bank_icon_url, banks.type_flag"
	selectedFields += ", account_types.name as account_type_name, account_types.limit_flag"
	if err := r.db.Table("Bank_accounts as accounts").
		Select(selectedFields).
		Joins("LEFT JOIN Banks AS banks ON banks.id = accounts.bank_id").
		Joins("LEFT JOIN Bank_account_types AS account_types ON account_types.id = accounts.account_type_id").
		Where("accounts.id = ?", id).
		Where("accounts.deleted_at IS NULL").
		First(&accounting).
		Error; err != nil {
		return nil, err
	}
	return &accounting, nil
}

func (r repo) GetDepositAccountById(id int64) (*model.BankAccount, error) {

	var accounting model.BankAccount
	selectedFields := "accounts.id, accounts.bank_id, accounts.account_type_id, accounts.account_name, accounts.account_number, accounts.account_balance, accounts.account_priority_id, accounts.account_status, accounts.device_uid, accounts.pin_code, accounts.connection_status"
	selectedFields += ", accounts.auto_credit_flag, accounts.is_main_withdraw, accounts.auto_withdraw_flag, accounts.auto_withdraw_credit_flag, accounts.auto_withdraw_confirm_flag, accounts.auto_withdraw_max_amount, accounts.auto_transfer_max_amount, accounts.qr_wallet_status"
	selectedFields += ", accounts.last_conn_update_at, accounts.created_at, accounts.updated_at"
	selectedFields += ", banks.name as bank_name, banks.code as bank_code, banks.icon_url as bank_icon_url, banks.type_flag"
	selectedFields += ", account_types.name as account_type_name, account_types.limit_flag"
	if err := r.db.Table("Bank_accounts as accounts").
		Select(selectedFields).
		Joins("LEFT JOIN Banks AS banks ON banks.id = accounts.bank_id").
		Joins("LEFT JOIN Bank_account_types AS account_types ON account_types.id = accounts.account_type_id").
		Where("accounts.id = ?", id).
		Where("account_types.allow_deposit = 1").
		Where("accounts.deleted_at IS NULL").
		First(&accounting).
		Error; err != nil {
		return nil, err
	}
	return &accounting, nil
}

func (r repo) GetWithdrawAccountById(id int64) (*model.BankAccount, error) {

	var accounting model.BankAccount
	selectedFields := "accounts.id, accounts.bank_id, accounts.account_type_id, accounts.account_name, accounts.account_number, accounts.account_balance, accounts.account_priority_id, accounts.account_status, accounts.device_uid, accounts.pin_code, accounts.connection_status"
	selectedFields += ", accounts.auto_credit_flag, accounts.is_main_withdraw, accounts.auto_withdraw_flag, accounts.auto_withdraw_credit_flag, accounts.auto_withdraw_confirm_flag, accounts.auto_withdraw_max_amount, accounts.auto_transfer_max_amount, accounts.qr_wallet_status"
	selectedFields += ", accounts.last_conn_update_at, accounts.created_at, accounts.updated_at"
	selectedFields += ", banks.name as bank_name, banks.code as bank_code, banks.icon_url as bank_icon_url, banks.type_flag"
	selectedFields += ", account_types.name as account_type_name, account_types.limit_flag"
	if err := r.db.Table("Bank_accounts as accounts").
		Select(selectedFields).
		Joins("LEFT JOIN Banks AS banks ON banks.id = accounts.bank_id").
		Joins("LEFT JOIN Bank_account_types AS account_types ON account_types.id = accounts.account_type_id").
		Where("accounts.id = ?", id).
		Where("account_types.allow_withdraw = 1").
		Where("accounts.deleted_at IS NULL").
		First(&accounting).
		Error; err != nil {
		return nil, err
	}
	return &accounting, nil
}

func (r repo) GetBankAccountByAccountNumber(accountNumber string) (*model.BankAccount, error) {

	var accounting model.BankAccount
	selectedFields := "accounts.id, accounts.bank_id, accounts.account_type_id, accounts.account_name, accounts.account_number, accounts.account_balance, accounts.account_priority_id, accounts.account_status, accounts.device_uid, accounts.pin_code, accounts.connection_status"
	selectedFields += ", accounts.auto_credit_flag, accounts.is_main_withdraw, accounts.auto_withdraw_flag, accounts.auto_withdraw_credit_flag, accounts.auto_withdraw_confirm_flag, accounts.auto_withdraw_max_amount, accounts.auto_transfer_max_amount, accounts.qr_wallet_status"
	selectedFields += ", accounts.last_conn_update_at, accounts.created_at, accounts.updated_at"
	selectedFields += ", banks.name as bank_name, banks.code as bank_code, banks.icon_url as bank_icon_url, banks.type_flag"
	selectedFields += ", account_types.name as account_type_name, account_types.limit_flag"
	if err := r.db.Table("Bank_accounts as accounts").
		Select(selectedFields).
		Joins("LEFT JOIN Banks AS banks ON banks.id = accounts.bank_id").
		Joins("LEFT JOIN Bank_account_types AS account_types ON account_types.id = accounts.account_type_id").
		Where("accounts.account_number = ?", accountNumber).
		Where("accounts.deleted_at IS NULL").
		First(&accounting).
		Error; err != nil {
		return nil, err
	}
	return &accounting, nil
}

func (r repo) GetBankAccountByExternalId(external_id int64) (*model.BankAccount, error) {

	var accounting model.BankAccount
	selectedFields := "accounts.id, accounts.bank_id, accounts.account_type_id, accounts.account_name, accounts.account_number, accounts.account_balance, accounts.account_priority_id, accounts.account_status, accounts.device_uid, accounts.pin_code, accounts.connection_status"
	selectedFields += ", accounts.auto_credit_flag, accounts.is_main_withdraw, accounts.auto_withdraw_flag, accounts.auto_withdraw_credit_flag, accounts.auto_withdraw_confirm_flag, accounts.auto_withdraw_max_amount, accounts.auto_transfer_max_amount, accounts.qr_wallet_status"
	selectedFields += ", accounts.last_conn_update_at, accounts.created_at, accounts.updated_at"
	if err := r.db.Table("Bank_accounts as accounts").
		Select(selectedFields).
		Where("accounts.external_id = ?", external_id).
		Where("accounts.deleted_at IS NULL").
		First(&accounting).
		Error; err != nil {
		return nil, err
	}
	return &accounting, nil
}

func (r repo) GetActiveExternalAccount() (*model.BankAccount, error) {

	var accounting model.BankAccount
	selectedFields := "accounts.id, accounts.bank_id, accounts.account_type_id, accounts.account_name, accounts.account_number, accounts.account_balance, accounts.account_priority_id, accounts.account_status, accounts.device_uid, accounts.pin_code, accounts.connection_status"
	selectedFields += ", accounts.auto_credit_flag, accounts.is_main_withdraw, accounts.auto_withdraw_flag, accounts.auto_withdraw_credit_flag, accounts.auto_withdraw_confirm_flag, accounts.auto_withdraw_max_amount, accounts.auto_transfer_max_amount, accounts.qr_wallet_status"
	selectedFields += ", accounts.last_conn_update_at, accounts.created_at, accounts.updated_at"
	if err := r.db.Table("Bank_accounts as accounts").
		Select(selectedFields).
		Where("accounts.connection_status = 'active'").
		Where("accounts.external_id IS NOT NULL").
		Where("accounts.deleted_at IS NULL").
		First(&accounting).
		Error; err != nil {
		return nil, err
	}
	return &accounting, nil
}

func (r repo) GetBankAccounts(req model.BankAccountListRequest) (*model.SuccessWithPagination, error) {

	var list []model.BankAccountResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("Bank_accounts AS accounts")
	count = count.Select("accounts.id")
	count = count.Joins("LEFT JOIN Banks AS banks ON banks.id = accounts.bank_id")
	count = count.Joins("LEFT JOIN Bank_account_types AS account_types ON account_types.id = accounts.account_type_id")
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where(r.db.Where("account_name LIKE ?", search_like).Or("account_number LIKE ?", search_like))
	}
	if req.AccountNumber != "" {
		count = count.Where("account_number = ?", req.AccountNumber)
	}
	if req.AccountType == "deposit" {
		count = count.Where("account_types.allow_deposit = 1")
	} else if req.AccountType == "withdraw" {
		count = count.Where("account_types.allow_withdraw = 1")
	}
	if err = count.
		Where("accounts.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	if total > 0 {
		// SELECT //
		query := r.db.Table("Bank_accounts AS accounts")
		selectedFields := "accounts.id, accounts.bank_id, accounts.account_type_id, accounts.account_name, accounts.account_number, accounts.account_balance, accounts.account_priority_id, accounts.account_status, accounts.device_uid, accounts.pin_code, accounts.connection_status"
		selectedFields += ", accounts.auto_credit_flag, accounts.is_main_withdraw, accounts.auto_withdraw_flag, accounts.auto_withdraw_credit_flag, accounts.auto_withdraw_confirm_flag, accounts.auto_withdraw_max_amount, accounts.auto_transfer_max_amount, accounts.qr_wallet_status"
		selectedFields += ", accounts.last_conn_update_at, accounts.created_at, accounts.updated_at"
		selectedFields += ", banks.name as bank_name, banks.code as bank_code, banks.icon_url as bank_icon_url, banks.type_flag"
		selectedFields += ", account_types.name as account_type_name, account_types.limit_flag"
		query = query.Select(selectedFields)
		query = query.Joins("LEFT JOIN Banks AS banks ON banks.id = accounts.bank_id")
		query = query.Joins("LEFT JOIN Bank_account_types AS account_types ON account_types.id = accounts.account_type_id")
		if req.Search != "" {
			search_like := fmt.Sprintf("%%%s%%", req.Search)
			query = query.Where(r.db.Where("accounts.account_name LIKE ?", search_like).Or("accounts.account_number LIKE ?", search_like))
		}
		if req.AccountNumber != "" {
			query = query.Where("accounts.account_number = ?", req.AccountNumber)
		}
		if req.AccountType == "deposit" {
			query = query.Where("account_types.allow_deposit = 1")
		} else if req.AccountType == "withdraw" {
			query = query.Where("account_types.allow_withdraw = 1")
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
			Where("accounts.deleted_at IS NULL").
			Offset(req.Page * req.Limit).
			Scan(&list).
			Error; err != nil {
			return nil, err
		}
	}

	// End count total records for pagination purposes (without limit and offset) //
	var result model.SuccessWithPagination
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) GetBankAccountPriorities() (*model.SuccessWithPagination, error) {

	var list []model.BankAccountPriority
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("Bank_account_priorities AS priorities")
	count = count.Select("priorities.id")
	if err = count.
		Where("priorities.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	if total > 0 {
		// SELECT //
		query := r.db.Table("Bank_account_priorities AS priorities")
		selectedFields := "priorities.id, priorities.name, priorities.condition_type, priorities.min_deposit_count, priorities.min_deposit_total, priorities.created_at, priorities.updated_at"
		query = query.Select(selectedFields)

		query = query.Order("id ASC")
		if err = query.
			Where("priorities.deleted_at IS NULL").
			Scan(&list).
			Error; err != nil {
			return nil, err
		}
	}

	// End count total records for pagination purposes (without limit and offset) //
	var result model.SuccessWithPagination
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) GetBotBankAccounts(req model.BankAccountListRequest) (*model.SuccessWithPagination, error) {

	var list []model.BankAccountResponse
	var total int64
	var err error
	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("Bank_accounts")
	count = count.Select("id")
	count = count.Where("device_uid != ?", "")
	count = count.Where("pin_code != ?", "")
	count = count.Where("external_id IS NOT NULL")
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where(r.db.Where("account_name LIKE ?", search_like).Or("account_number LIKE ?", search_like))
	}
	if err = count.
		Where("deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	if total > 0 {
		// SELECT //
		query := r.db.Table("Bank_accounts AS accounts")
		selectedFields := "accounts.id, accounts.bank_id, accounts.account_type_id, accounts.account_name, accounts.account_number, accounts.account_balance, accounts.account_priority_id, accounts.account_status, accounts.device_uid, accounts.pin_code, accounts.connection_status"
		selectedFields += ", accounts.auto_credit_flag, accounts.is_main_withdraw, accounts.auto_withdraw_flag, accounts.auto_withdraw_credit_flag, accounts.auto_withdraw_confirm_flag, accounts.auto_withdraw_max_amount, accounts.auto_transfer_max_amount, accounts.qr_wallet_status"
		selectedFields += ", accounts.last_conn_update_at, accounts.created_at, accounts.updated_at"
		selectedFields += ", banks.name as bank_name, banks.code as bank_code, banks.icon_url as bank_icon_url, banks.type_flag"
		selectedFields += ", account_types.name as account_type_name, account_types.limit_flag"
		query = query.Select(selectedFields)
		query = query.Joins("LEFT JOIN Banks AS banks ON banks.id = accounts.bank_id")
		query = query.Joins("LEFT JOIN Bank_account_types AS account_types ON account_types.id = accounts.account_type_id")
		query = query.Where("accounts.device_uid != ?", "")
		query = query.Where("accounts.pin_code != ?", "")
		query = query.Where("accounts.external_id IS NOT NULL")
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
			Where("accounts.deleted_at IS NULL").
			Offset(req.Page * req.Limit).
			Scan(&list).
			Error; err != nil {
			return nil, err
		}
	}

	// End count total records for pagination purposes (without limit and offset) //
	var result model.SuccessWithPagination
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) ResetMainWithdrawBankAccount() error {

	cleanUpData := map[string]interface{}{
		"is_main_withdraw": 0,
	}
	if err := r.db.Table("Bank_accounts").Where("is_main_withdraw != 0").Updates(cleanUpData).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) CreateBankAccount(data model.BankAccountCreateBody) error {

	if data.IsMainWithdraw {
		if err := r.ResetMainWithdrawBankAccount(); err != nil {
			return err
		}
	}

	if err := r.db.Table("Bank_accounts").Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) UpdateBankAccount(id int64, data model.BankAccountUpdateBody) error {
	if err := r.db.Table("Bank_accounts").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) DeleteBankAccount(id int64, data model.BankAccountDeleteBody) error {
	if err := r.db.Table("Bank_accounts").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) GetTransactionById(id int64) (*model.BankAccountTransaction, error) {
	var record model.BankAccountTransaction
	selectedFields := "transactions.id, transactions.account_id, transactions.description, transactions.transfer_type, transactions.amount, transactions.transfer_at, transactions.created_by_username, transactions.created_at, transactions.updated_at"
	selectedFields += ",accounts.bank_id, accounts.account_type_id, accounts.account_name, accounts.account_number, accounts.account_balance, accounts.account_priority_id, accounts.account_status, accounts.created_at, accounts.updated_at"
	selectedFields += ",banks.name as bank_name, banks.code as bank_code, banks.icon_url as bank_icon_url, banks.type_flag"
	if err := r.db.Table("Bank_account_transactions as transactions").
		Select(selectedFields).
		Joins("LEFT JOIN Bank_accounts AS accounts ON accounts.id = transactions.account_id").
		Joins("LEFT JOIN Banks AS banks ON banks.id = accounts.bank_id").
		Where("transactions.id = ?", id).
		Where("transactions.deleted_at IS NULL").
		First(&record).
		Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r repo) GetTransactions(req model.BankAccountTransactionListRequest) (*model.SuccessWithPagination, error) {

	var list []model.BankAccountTransactionResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("Bank_account_transactions as transactions")
	count = count.Select("transactions.id")
	count = count.Joins("LEFT JOIN Bank_accounts AS accounts ON accounts.id = transactions.account_id")
	count = count.Joins("LEFT JOIN Banks AS banks ON banks.id = accounts.bank_id")
	if req.AccountId != 0 {
		count = count.Where("transactions.account_id = ?", req.AccountId)
	}
	if req.FromCreatedDate != "" {
		count = count.Where("transactions.created_at >= ?", req.FromCreatedDate)
	}
	if req.ToCreatedDate != "" {
		count = count.Where("transactions.created_at <= ?", req.ToCreatedDate)
	}
	if req.TransferType != "" {
		count = count.Where("transactions.transfer_type = ?", req.TransferType)
	}
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where(r.db.Where("transactions.description LIKE ?", search_like).Or("accounts.account_name LIKE ?", search_like).Or("accounts.account_number LIKE ?", search_like))
	}
	if err = count.
		Where("transactions.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	if total > 0 {
		// SELECT //
		selectedFields := "transactions.id, transactions.account_id, transactions.description, transactions.transfer_type, transactions.amount, transactions.transfer_at, transactions.created_by_username, transactions.created_at, transactions.updated_at"
		selectedFields += ",accounts.bank_id, accounts.account_type_id, accounts.account_name, accounts.account_number, accounts.account_balance, accounts.account_priority_id, accounts.account_status, accounts.created_at, accounts.updated_at"
		selectedFields += ",banks.name as bank_name, banks.code as bank_code, banks.icon_url as bank_icon_url, banks.type_flag"
		query := r.db.Table("Bank_account_transactions as transactions")
		query = query.Select(selectedFields)
		query = query.Joins("LEFT JOIN Bank_accounts AS accounts ON accounts.id = transactions.account_id")
		query = query.Joins("LEFT JOIN Banks AS banks ON banks.id = accounts.bank_id")
		if req.AccountId != 0 {
			query = query.Where("transactions.account_id = ?", req.AccountId)
		}
		if req.FromCreatedDate != "" {
			query = query.Where("transactions.created_at >= ?", req.FromCreatedDate)
		}
		if req.ToCreatedDate != "" {
			query = query.Where("transactions.created_at <= ?", req.ToCreatedDate)
		}
		if req.TransferType != "" {
			query = query.Where("transactions.transfer_type = ?", req.TransferType)
		}
		if req.Search != "" {
			search_like := fmt.Sprintf("%%%s%%", req.Search)
			query = query.Where(r.db.Where("transactions.description LIKE ?", search_like).Or("accounts.account_name LIKE ?", search_like).Or("accounts.account_number LIKE ?", search_like))
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

	// End count total records for pagination purposes (without limit and offset) //
	var result model.SuccessWithPagination
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) CreateTransaction(data model.BankAccountTransactionBody) error {
	if err := r.db.Table("Bank_account_transactions").Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) UpdateTransaction(id int64, data model.BankAccountTransactionBody) error {
	if err := r.db.Table("Bank_account_transactions").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) DeleteTransaction(id int64) error {
	if err := r.db.Table("Bank_account_transactions").Where("id = ?", id).Delete(&model.BankAccountTransaction{}).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) GetTransferById(id int64) (*model.BankAccountTransfer, error) {
	var record model.BankAccountTransfer
	selectedFields := "transfers.id, transfers.from_account_id, transfers.from_bank_id, transfers.from_account_name, transfers.from_account_number"
	selectedFields += ",transfers.to_account_id, transfers.to_bank_id, transfers.to_account_name, transfers.to_account_number"
	selectedFields += ",transfers.amount, transfers.transfer_at, transfers.created_by_username, transfers.status, transfers.confirmed_at, transfers.confirmed_by_user_id, transfers.created_at, transfers.updated_at"
	selectedFields += ",from_banks.name as from_bank_name, from_banks.code as from_bank_code, from_banks.icon_url as from_bank_icon_url, from_banks.type_flag as from_bank_type_flag"
	selectedFields += ",to_banks.name as to_bank_name, to_banks.code as to_bank_code, to_banks.icon_url as to_bank_icon_url, to_banks.type_flag as to_bank_type_flag"
	if err := r.db.Table("Bank_account_transfers as transfers").
		Select(selectedFields).
		Joins("LEFT JOIN Banks AS from_banks ON from_banks.id = transfers.from_bank_id").
		Joins("LEFT JOIN Banks AS to_banks ON to_banks.id = transfers.to_bank_id").
		Where("transfers.id = ?", id).
		Where("transfers.deleted_at IS NULL").
		First(&record).
		Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r repo) GetTransfers(req model.BankAccountTransferListRequest) (*model.SuccessWithPagination, error) {

	var list []model.BankAccountTransferResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("Bank_account_transfers as transfers")
	count = count.Joins("LEFT JOIN Banks AS from_banks ON from_banks.id = transfers.from_bank_id")
	count = count.Joins("LEFT JOIN Banks AS to_banks ON to_banks.id = transfers.to_bank_id")
	count = count.Select("transfers.id")
	if req.AccountId != 0 {
		count = count.Where("transfers.from_account_id = ?", req.AccountId)
	}
	if req.FromCreatedDate != "" {
		count = count.Where("transfers.created_at >= ?", req.FromCreatedDate)
	}
	if req.ToCreatedDate != "" {
		count = count.Where("transfers.created_at <= ?", req.ToCreatedDate)
	}
	if req.ToAccountId != 0 {
		count = count.Where("transfers.to_account_id = ?", req.ToAccountId)
	}
	if req.Status != "" {
		count = count.Where("transfers.status = ?", req.Status)
	}
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where(r.db.Where("transfers.from_account_name LIKE ?", search_like).Or("transfers.from_account_number LIKE ?", search_like).Or("transfers.to_account_name LIKE ?", search_like).Or("transfers.to_account_number LIKE ?", search_like))
	}
	if err = count.
		Where("transfers.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	if total > 0 {
		// SELECT //
		selectedFields := "transfers.id, transfers.from_account_id, transfers.from_bank_id, transfers.from_account_name, transfers.from_account_number"
		selectedFields += ",transfers.to_account_id, transfers.to_bank_id, transfers.to_account_name, transfers.to_account_number"
		selectedFields += ",transfers.amount, transfers.transfer_at, transfers.created_by_username, transfers.status, transfers.confirmed_at, transfers.confirmed_by_user_id, transfers.created_at, transfers.updated_at"
		selectedFields += ",from_banks.name as from_bank_name, to_banks.name as to_bank_name"
		query := r.db.Table("Bank_account_transfers as transfers")
		query = query.Select(selectedFields)
		query = query.Joins("LEFT JOIN Banks AS from_banks ON from_banks.id = transfers.from_bank_id")
		query = query.Joins("LEFT JOIN Banks AS to_banks ON to_banks.id = transfers.to_bank_id")
		if req.AccountId != 0 {
			query = query.Where("transfers.from_account_id = ?", req.AccountId)
		}
		if req.FromCreatedDate != "" {
			query = query.Where("transfers.created_at >= ?", req.FromCreatedDate)
		}
		if req.ToCreatedDate != "" {
			query = query.Where("transfers.created_at <= ?", req.ToCreatedDate)
		}
		if req.ToAccountId != 0 {
			query = query.Where("transfers.to_account_id = ?", req.ToAccountId)
		}
		if req.Status != "" {
			query = query.Where("transfers.status = ?", req.Status)
		}
		if req.Search != "" {
			search_like := fmt.Sprintf("%%%s%%", req.Search)
			query = query.Where(r.db.Where("transfers.from_account_name LIKE ?", search_like).Or("transfers.from_account_number LIKE ?", search_like).Or("transfers.to_account_name LIKE ?", search_like).Or("transfers.to_account_number LIKE ?", search_like))
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
			Where("transfers.deleted_at IS NULL").
			Offset(req.Page * req.Limit).
			Scan(&list).
			Error; err != nil {
			return nil, err
		}
	}

	// End count total records for pagination purposes (without limit and offset) //
	var result model.SuccessWithPagination
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) CreateTransfer(data model.BankAccountTransferBody) error {
	if err := r.db.Table("Bank_account_transfers").Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) ConfirmTransfer(id int64, data model.BankAccountTransferConfirmBody) error {
	if err := r.db.Table("Bank_account_transfers").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) DeleteTransfer(id int64) error {
	if err := r.db.Table("Bank_account_transfers").Where("id = ?", id).Delete(&model.BankAccountTransfer{}).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) CreateWebhookLog(data model.WebhookLogCreateBody) (*int64, error) {
	if err := r.db.Table("Webhook_logs").Create(&data).Error; err != nil {
		return nil, err
	}
	return &data.Id, nil
}

func (r repo) UpdateWebhookLog(id int64, data model.WebhookLogUpdateBody) error {
	if err := r.db.Table("Webhook_logs").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) GetWebhookStatementByExternalId(external_id int64) (*model.BankStatement, error) {

	var record model.BankStatement
	selectedFields := "statements.id, statements.external_id, statements.account_id, statements.detail, statements.statement_type, statements.transfer_at, statements.amount, statements.status, statements.created_at, statements.updated_at"
	if err := r.db.Table("Bank_statements as statements").
		Select(selectedFields).
		Where("statements.external_id = ?", external_id).
		Where("statements.deleted_at IS NULL").
		First(&record).
		Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r repo) CreateWebhookStatement(data model.BankStatementCreateBody) (*int64, error) {
	if err := r.db.Table("Bank_statements").Create(&data).Error; err != nil {
		return nil, err
	}
	return &data.Id, nil
}

func (r repo) GetBotaccountConfigs(req model.BotAccountConfigListRequest) (*model.SuccessWithPagination, error) {

	var list []model.BotAccountConfig
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("Botaccount_config as configs")
	count = count.Select("configs.id")
	if req.SearchKey != nil {
		count = count.Where("configs.config_key = ?", req.SearchKey)
	}
	if req.SearchValue != nil {
		count = count.Where("configs.config_val = ?", req.SearchValue)
	}
	if err = count.
		Where("configs.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	if total > 0 {
		// SELECT //
		selectedFields := "configs.id, configs.config_key, configs.config_val"
		query := r.db.Table("Botaccount_config as configs")
		query = query.Select(selectedFields)
		if req.SearchKey != nil {
			query = query.Where("configs.config_key = ?", req.SearchKey)
		}
		if req.SearchValue != nil {
			query = query.Where("configs.config_val = ?", req.SearchValue)
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
			Where("configs.deleted_at IS NULL").
			Offset(req.Page * req.Limit).
			Scan(&list).
			Error; err != nil {
			return nil, err
		}
	}

	// End count total records for pagination purposes (without limit and offset) //
	var result model.SuccessWithPagination
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) GetBotaccountConfigByKey(req model.BotAccountConfigListRequest) (*model.BotAccountConfig, error) {

	var record model.BotAccountConfig
	var err error

	// SELECT //
	selectedFields := "configs.id, configs.config_key, configs.config_val"
	query := r.db.Table("Botaccount_config as configs")
	query = query.Select(selectedFields)
	if req.SearchKey != nil {
		query = query.Where("configs.config_key = ?", req.SearchKey)
	}
	if req.SearchValue != nil {
		query = query.Where("configs.config_val = ?", req.SearchValue)
	}
	if err = query.
		Where("configs.deleted_at IS NULL").
		First(&record).
		Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r repo) CreateBotaccountConfig(data model.BotAccountConfigCreateBody) error {
	if err := r.db.Table("Botaccount_config").Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) DeleteBotaccountConfigByKey(key string) error {
	if err := r.db.Table("Botaccount_config").Where("config_key = ?", key).Delete(&model.BankAccountTransfer{}).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) DeleteBotaccountConfigById(id int64) error {
	if err := r.db.Table("Botaccount_config").Where("id = ?", id).Delete(&model.BankAccountTransfer{}).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) MatchAutoStatementOwner(id int64, data model.BankStatementUpdateBody) error {
	if err := r.db.Table("Bank_statements").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}
