package service

import (
	"cybergame-api/helper"
	"cybergame-api/model"
	"cybergame-api/repository"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type BankingService interface {
	GetBankStatementById(req model.GetByIdRequest) (*model.BankStatement, error)
	GetBankStatements(req model.BankStatementListRequest) (*model.SuccessWithPagination, error)
	GetBankStatementSummary(req model.BankStatementListRequest) (*model.BankStatementSummary, error)
	CreateBankStatement(data model.BankStatementCreateBody) error
	MatchStatementOwner(id int64, req model.BankStatementMatchRequest) error
	IgnoreStatementOwner(id int64, req model.BankStatementMatchRequest) error
	DeleteBankStatement(id int64) error

	GetBankTransactionById(req model.BankTransactionGetRequest) (*model.BankTransaction, error)
	GetBankTransactions(req model.BankTransactionListRequest) (*model.SuccessWithPagination, error)
	GetBankDepositTransactions(req model.BankTransactionListRequest) (*model.SuccessWithPagination, error)
	GetBankWithdrawTransactions(req model.BankTransactionListRequest) (*model.SuccessWithPagination, error)
	GetBankDepositTransStatusCounts(req model.BankTransactionListRequest) (*model.BankDepositTransStatusCounts, error)
	GetBankWithdrawTransStatusCounts(req model.BankTransactionListRequest) (*model.BankWithdrawTransStatusCounts, error)
	CreateBankTransaction(data model.BankTransactionCreateBody) error
	CreateBonusTransaction(data model.BonusTransactionCreateBody) error
	UpdateBankTransaction(id int64, req model.BankTransactionUpdateRequest) error
	DeleteBankTransaction(id int64) error

	GetPendingDepositTransactions(req model.PendingDepositTransactionListRequest) (*model.SuccessWithPagination, error)
	GetPendingWithdrawTransactions(req model.PendingWithdrawTransactionListRequest) (*model.SuccessWithPagination, error)
	ConfirmDepositTransaction(id int64, req model.BankConfirmDepositRequest) error
	ConfirmDepositCredit(id int64, req model.BankConfirmDepositRequest) error
	ContinueAutoWithdrawTransaction(id int64) error
	ConfirmWithdrawTransaction(id int64, req model.BankConfirmCreditWithdrawRequest) error
	ConfirmWithdrawTransfer(id int64, req model.BankConfirmTransferWithdrawRequest) error
	CancelPendingTransaction(id int64, data model.BankTransactionCancelBody) error
	GetFinishedTransactions(req model.FinishedTransactionListRequest) (*model.SuccessWithPagination, error)
	RemoveFinishedTransaction(id int64, data model.BankTransactionRemoveBody) error
	GetRemovedTransactions(req model.RemovedTransactionListRequest) (*model.SuccessWithPagination, error)

	GetMemberByCode(code string) (*model.Member, error)
	GetMembers(req model.MemberListRequest) (*model.SuccessWithPagination, error)
	GetPossibleStatementOwners(req model.MemberPossibleListRequest) (*model.SuccessWithPagination, error)
	GetMemberTransactions(req model.MemberTransactionListRequest) (*model.SuccessWithPagination, error)
	GetMemberTransactionSummary(req model.MemberTransactionListRequest) (*model.MemberTransactionSummary, error)
	GetMemberStatementById(req model.GetByIdRequest) (*model.MemberStatementResponse, error)
	GetMemberStatements(req model.MemberStatementListRequest) (*model.SuccessWithPagination, error)
	ProcessMemberDepositCredit(userId int64, amount float64) error
	ProcessMemberWithdrawCredit(userId int64, amount float64) error
	ProcessMemberBonusCredit(userId int64, amount float64) error
	ProcessMemberGetbackCredit(userId int64, amount float64) error
}

var memberNotFound = "Member not found"
var bankStatementferNotFound = "Statement not found"
var bankTransactionferNotFound = "Transaction not found"

type bankingService struct {
	repoBanking      repository.BankingRepository
	repoAccounting   repository.AccountingRepository
	repoAgentConnect repository.AgentConnectRepository
}

func NewBankingService(
	repoBanking repository.BankingRepository,
	repoAccounting repository.AccountingRepository,
	repoAgentConnect repository.AgentConnectRepository,
) BankingService {
	return &bankingService{repoBanking, repoAccounting, repoAgentConnect}
}

func (s *bankingService) GetBankStatementById(req model.GetByIdRequest) (*model.BankStatement, error) {

	record, err := s.repoBanking.GetBankStatementById(req.Id)
	if err != nil {
		if err.Error() == recordNotFound {
			return nil, notFound(bankStatementferNotFound)
		}
		return nil, internalServerError(err.Error())
	}
	return record, nil
}

func (s *bankingService) GetBankStatements(req model.BankStatementListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.UnlimitPagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	records, err := s.repoBanking.GetBankStatements(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) GetBankStatementSummary(req model.BankStatementListRequest) (*model.BankStatementSummary, error) {

	records, err := s.repoBanking.GetBankStatementSummary(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) CreateBankStatement(data model.BankStatementCreateBody) error {

	toAccount, err := s.repoAccounting.GetBankAccountById(data.AccountId)
	if err != nil {
		fmt.Println(err)
		return badRequest("Invalid Bank Account")
	}
	var body model.BankStatementCreateBody
	body.AccountId = toAccount.Id
	if data.StatementType == "transfer_in" {
		body.Amount = data.Amount
	} else if data.StatementType == "transfer_out" {
		body.Amount = data.Amount * -1
	} else {
		return badRequest("Invalid Transfer Type")
	}
	body.Detail = data.Detail
	body.StatementType = data.StatementType
	body.TransferAt = data.TransferAt
	body.Status = "pending"

	if err := s.repoBanking.CreateBankStatement(body); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *bankingService) MatchStatementOwner(id int64, req model.BankStatementMatchRequest) error {

	statement, err := s.repoBanking.GetBankStatementById(id)
	if err != nil {
		return internalServerError(err.Error())
	}
	if statement.Status != "pending" {
		return badRequest("Statement is not pending")
	}
	member, err := s.repoBanking.GetMemberById(req.UserId)
	if err != nil {
		return badRequest("Invalid Member")
	}
	jsonBefore, _ := json.Marshal(statement)

	// TransAction
	var createBody model.CreateBankStatementActionBody
	createBody.StatementId = statement.Id
	createBody.UserId = member.Id
	createBody.ActionType = "confirmed"
	createBody.AccountId = statement.AccountId
	createBody.JsonBefore = string(jsonBefore)
	createBody.ConfirmedAt = req.ConfirmedAt
	createBody.ConfirmedByUserId = req.ConfirmedByUserId
	createBody.ConfirmedByUsername = req.ConfirmedByUsername
	if err := s.repoBanking.CreateStatementAction(createBody); err == nil {
		var body model.BankStatementUpdateBody
		body.Status = "confirmed"
		if err := s.repoBanking.MatchStatementOwner(id, body); err != nil {
			return internalServerError(err.Error())
		}
	} else {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *bankingService) IgnoreStatementOwner(id int64, req model.BankStatementMatchRequest) error {

	statement, err := s.repoBanking.GetBankStatementById(id)
	if err != nil {
		return internalServerError(err.Error())
	}
	if statement.Status != "pending" {
		return badRequest("Statement is not pending")
	}

	var body model.BankStatementUpdateBody
	body.Status = "ignored"
	jsonBefore, _ := json.Marshal(statement)

	// TransAction
	var createBody model.CreateBankStatementActionBody
	createBody.StatementId = statement.Id
	createBody.ActionType = "ignored"
	createBody.AccountId = statement.AccountId
	createBody.JsonBefore = string(jsonBefore)
	createBody.ConfirmedAt = req.ConfirmedAt
	createBody.ConfirmedByUserId = req.ConfirmedByUserId
	createBody.ConfirmedByUsername = req.ConfirmedByUsername
	if err := s.repoBanking.CreateStatementAction(createBody); err == nil {
		if err := s.repoBanking.IgnoreStatementOwner(id, body); err != nil {
			return internalServerError(err.Error())
		}
	} else {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *bankingService) DeleteBankStatement(id int64) error {

	_, err := s.repoBanking.GetBankStatementById(id)
	if err != nil {
		return internalServerError(err.Error())
	}

	if err := s.repoBanking.DeleteBankStatement(id); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *bankingService) GetBankTransactionById(req model.BankTransactionGetRequest) (*model.BankTransaction, error) {

	record, err := s.repoBanking.GetBankTransactionById(req.Id)
	if err != nil {
		if err.Error() == recordNotFound {
			return nil, notFound(bankTransactionferNotFound)
		}
		return nil, internalServerError(err.Error())
	}
	return record, nil
}

func (s *bankingService) GetBankTransactions(req model.BankTransactionListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.UnlimitPagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	records, err := s.repoBanking.GetBankTransactions(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) GetBankDepositTransactions(req model.BankTransactionListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.UnlimitPagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}

	if req.TransferType == "" || req.TransferType == "all" {
		req.TransferType = "all_deposit"
	}
	records, err := s.repoBanking.GetBankTransactions(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) GetBankWithdrawTransactions(req model.BankTransactionListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.UnlimitPagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}

	if req.TransferType == "" || req.TransferType == "all" {
		req.TransferType = "all_withdraw"
	}
	records, err := s.repoBanking.GetBankTransactions(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) GetBankDepositTransStatusCounts(req model.BankTransactionListRequest) (*model.BankDepositTransStatusCounts, error) {

	req.TransferType = "all_deposit"
	records, err := s.repoBanking.GetBankTransactionStatusCount(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}

	fmt.Println("GetBankDepositTransStatusCounts", records)

	var result model.BankDepositTransStatusCounts
	result.AllCount = 0
	for _, record := range records.List.([]model.BankTransactionStatusCount) {
		result.AllCount += record.Count
		if record.Status == "pending" {
			result.PendingCount = record.Count
		} else if record.Status == "pending_credit" {
			result.PendingCreditCount = record.Count
		} else if record.Status == "finished" {
			result.FinishedCount = record.Count
		} else if record.Status == "canceled" {
			result.FailedCount = record.Count
		}
	}
	return &result, nil
}

func (s *bankingService) GetBankWithdrawTransStatusCounts(req model.BankTransactionListRequest) (*model.BankWithdrawTransStatusCounts, error) {

	req.TransferType = "all_withdraw"
	records, err := s.repoBanking.GetBankTransactionStatusCount(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}

	fmt.Println("GetBankWithdrawTransStatusCounts", records)

	var result model.BankWithdrawTransStatusCounts
	result.AllCount = 0
	for _, record := range records.List.([]model.BankTransactionStatusCount) {
		result.AllCount += record.Count
		if record.Status == "pending_credit" {
			result.PendingCreditCount = record.Count
		} else if record.Status == "pending_transfer" {
			result.PendingTransferCount = record.Count
		} else if record.Status == "finished" {
			result.FinishedCount = record.Count
		} else if record.Status == "canceled" {
			result.FailedCount = record.Count
		}
	}
	return &result, nil
}

func (s *bankingService) CreateBankTransaction(data model.BankTransactionCreateBody) error {

	var body model.BankTransactionCreateBody
	body.TransferAt = data.TransferAt
	body.CreatedByUserId = data.CreatedByUserId
	body.CreatedByUsername = data.CreatedByUsername
	body.Status = "pending"

	if data.TransferType == "deposit" {
		member, err := s.repoAccounting.GetUserByMemberCode(data.MemberCode)
		if err != nil {
			fmt.Println(err)
			return badRequest("Invalid Member code")
		}
		bank, err := s.repoAccounting.GetBankByCode(member.BankCode)
		if err != nil {
			fmt.Println(err)
			return badRequest("Invalid User Bank")
		}
		body.MemberCode = *member.MemberCode
		body.UserId = member.Id
		body.CreditAmount = data.CreditAmount
		body.TransferType = data.TransferType
		body.DepositChannel = data.DepositChannel
		body.OverAmount = data.OverAmount
		body.IsAutoCredit = data.IsAutoCredit

		body.FromBankId = &bank.Id
		body.FromAccountName = &member.Fullname
		body.FromAccountNumber = &member.BankAccount
		if data.ToAccountId == nil {
			return badRequest("Input Bank Account")
		}
		toAccount, err := s.repoAccounting.GetDepositAccountById(*data.ToAccountId)
		if err != nil {
			fmt.Println(err)
			return badRequest("Invalid Bank Account")
		}
		body.ToAccountId = &toAccount.Id
		body.ToBankId = &toAccount.BankId
		body.ToAccountName = &toAccount.AccountName
		body.ToAccountNumber = &toAccount.AccountNumber
		// later: createBonus + refDeposit
		body.BonusAmount = data.BonusAmount
		body.PromotionId = data.PromotionId

		transactionId, err := s.repoBanking.CreateBankDepositTransaction(body)
		if err != nil {
			return internalServerError(err.Error())
		}

		agentName := os.Getenv("AGENT_NAME")
		sign := agentName + *member.Username
		timeNow := time.Now()
		agentData := model.AGCDeposit{
			Agentname:     agentName,
			PlayerName:    *member.Username,
			Amount:        body.CreditAmount,
			Timestamp:     timeNow.Unix(),
			Sign:          helper.CreateSign(sign, timeNow),
			TransactionId: strconv.FormatInt(*transactionId, 10),
		}

		if err := s.repoAgentConnect.Deposit(agentData); err != nil {
			return internalServerError(err.Error())
		}
	} else if data.TransferType == "withdraw" {
		var autoWithdrawCondition *model.BankAutoWithdrawCondition

		member, err := s.repoAccounting.GetUserByMemberCode(data.MemberCode)
		if err != nil {
			fmt.Println(err)
			return badRequest("Invalid Member code")
		}
		bank, err := s.repoAccounting.GetBankByCode(member.BankCode)
		if err != nil {
			fmt.Println(err)
			return badRequest("Invalid User Bank")
		}
		body.MemberCode = *member.MemberCode
		body.UserId = member.Id
		body.CreditAmount = data.CreditAmount
		body.TransferType = data.TransferType

		// Withdraw SystemAccount is no more requried
		if data.FromAccountId != nil {
			fromAccount, err := s.repoAccounting.GetWithdrawAccountById(*data.FromAccountId)
			if err != nil {
				fmt.Println(err)
				return badRequest("Invalid Bank Account")
			}
			body.FromAccountId = &fromAccount.Id
			body.FromBankId = &fromAccount.BankId
			body.FromAccountName = &fromAccount.AccountName
			body.FromAccountNumber = &fromAccount.AccountNumber
			if condition, err := s.GetNewAutoWithdrawCondition(body, *fromAccount); err == nil {
				autoWithdrawCondition = condition
			}
		}

		body.ToBankId = &bank.Id
		body.ToAccountName = &member.Fullname
		body.ToAccountNumber = &member.BankAccount
		body.Status = "pending_credit"

		// later : reserved/check current user credit ?
		// for now : no need to cut user Credit !
		if insertId, err := s.repoBanking.CreateBankWithdrawTransaction(body); err == nil {
			if insertId != nil && autoWithdrawCondition != nil {

				agentName := os.Getenv("AGENT_NAME")
				sign := agentName + *member.Username
				timeNow := time.Now()
				agentData := model.AGCDeposit{
					Agentname:     agentName,
					PlayerName:    *member.Username,
					Amount:        body.CreditAmount,
					Timestamp:     timeNow.Unix(),
					Sign:          helper.CreateSign(sign, timeNow),
					TransactionId: strconv.FormatInt(*insertId, 10),
				}

				if err := s.repoAgentConnect.Deposit(agentData); err != nil {
					return internalServerError(err.Error())
				}

				autoWithdrawCondition.TransId = *insertId
				if err := s.ProcessAutoWithdrawCondition(*autoWithdrawCondition); err != nil {
					return internalServerError(err.Error())
				}
			}
		} else {
			return internalServerError(err.Error())
		}
	} else if data.TransferType == "getcreditback" {
		// จะดึงยอดสลายไปเลย
		member, err := s.repoAccounting.GetUserByMemberCode(data.MemberCode)
		if err != nil {
			fmt.Println(err)
			return badRequest("Invalid Member code")
		}
		bank, err := s.repoAccounting.GetBankByCode(member.BankCode)
		if err != nil {
			fmt.Println(err)
			return badRequest("Invalid User Bank")
		}
		body.MemberCode = *member.MemberCode
		body.UserId = member.Id
		body.CreditAmount = data.CreditAmount
		body.TransferType = data.TransferType
		body.FromBankId = &bank.Id
		body.FromAccountName = &member.Fullname
		body.FromAccountNumber = &member.BankAccount
		body.Status = "pending_credit"
		// later : reserved/check current user credit ?
		// for now : not cut user Credit !
		if _, err := s.repoBanking.CreateBankWithdrawTransaction(body); err != nil {
			return internalServerError(err.Error())
		}
	} else {
		return badRequest("Invalid Transfer Type")
	}
	return nil
}

func (s *bankingService) GetNewAutoWithdrawCondition(body model.BankTransactionCreateBody, fromAccount model.BankAccount) (*model.BankAutoWithdrawCondition, error) {

	var autoWithdrawCondition model.BankAutoWithdrawCondition
	autoWithdrawCondition.UserId = body.UserId
	autoWithdrawCondition.TransStatus = body.Status
	autoWithdrawCondition.CreditAmount = body.CreditAmount
	autoWithdrawCondition.FromAccountId = fromAccount.Id

	if fromAccount.AutoWithdrawCreditFlag == "auto" {
		autoWithdrawCondition.AutoWithdrawConfirmFlag = "auto"
	}
	if fromAccount.AutoWithdrawCreditFlag == "auto" {
		autoWithdrawCondition.AutoWithdrawCreditFlag = "auto"
	}

	// later : More Conditions
	return &autoWithdrawCondition, nil
}

func (s *bankingService) GetAutoWithdrawCondition(body model.BankTransaction, fromAccount model.BankAccount) (*model.BankAutoWithdrawCondition, error) {

	var autoWithdrawCondition model.BankAutoWithdrawCondition
	autoWithdrawCondition.UserId = body.UserId
	autoWithdrawCondition.TransStatus = body.Status
	autoWithdrawCondition.CreditAmount = body.CreditAmount
	autoWithdrawCondition.FromAccountId = fromAccount.Id

	if fromAccount.AutoWithdrawCreditFlag == "auto" {
		autoWithdrawCondition.AutoWithdrawConfirmFlag = "auto"
	}
	if fromAccount.AutoWithdrawCreditFlag == "auto" {
		autoWithdrawCondition.AutoWithdrawCreditFlag = "auto"
	}

	// later : More Conditions
	return &autoWithdrawCondition, nil
}

func (s *bankingService) SetAutoWithdrawCondition(transId int64, curCondition *model.BankAutoWithdrawCondition) error {

	// Fix Same request on auto+auto
	transaction, err := s.repoBanking.GetBankTransactionById(transId)
	if err != nil {
		return internalServerError(err.Error())
	}
	if transaction.TransferType != "withdraw" && transaction.TransferType != "getcreditback" {
		return badRequest("Transaction is not withdraw")
	}
	curCondition.TransId = transaction.Id
	curCondition.UserId = transaction.UserId
	curCondition.TransStatus = transaction.Status
	curCondition.CreditAmount = transaction.CreditAmount
	curCondition.FromAccountId = transaction.FromAccountId

	// if fromAccount.AutoWithdrawCreditFlag == "auto" {
	// 	autoWithdrawCondition.AutoWithdrawConfirmFlag = "auto"
	// }
	// if fromAccount.AutoWithdrawCreditFlag == "auto" {
	// 	autoWithdrawCondition.AutoWithdrawCreditFlag = "auto"
	// }

	// later : More Conditions
	return nil
}

func (s *bankingService) ProcessAutoWithdrawCondition(req model.BankAutoWithdrawCondition) error {

	if req.TransStatus == "pending" && req.AutoWithdrawCreditFlag == "auto" {
		var confirmReq model.BankConfirmCreditWithdrawRequest
		confirmReq.FromAccountId = &req.FromAccountId
		confirmReq.CreditAmount = &req.CreditAmount
		confirmReq.BankChargeAmount = &req.BankChargeAmount
		confirmReq.ConfirmedAt = time.Now()
		confirmReq.ConfirmedByUserId = 0
		confirmReq.ConfirmedByUsername = "อัตโนมัติ"
		actionErr := s.ConfirmWithdrawTransaction(req.TransId, confirmReq)
		if actionErr != nil {
			return internalServerError(actionErr.Error())
		}
	} else if req.TransStatus == "pending_credit" && req.AutoWithdrawCreditFlag == "auto" {
		var confirmReq model.BankConfirmCreditWithdrawRequest
		confirmReq.FromAccountId = &req.FromAccountId
		confirmReq.CreditAmount = &req.CreditAmount
		confirmReq.BankChargeAmount = &req.BankChargeAmount
		confirmReq.ConfirmedAt = time.Now()
		confirmReq.ConfirmedByUserId = 0
		confirmReq.ConfirmedByUsername = "อัตโนมัติ"
		actionErr := s.ConfirmWithdrawTransaction(req.TransId, confirmReq)
		if actionErr != nil {
			return internalServerError(actionErr.Error())
		}
	} else if req.TransStatus == "pending_transfer" && req.AutoWithdrawConfirmFlag == "auto" {
		var confirmReq model.BankConfirmTransferWithdrawRequest
		confirmReq.FromAccountId = &req.FromAccountId
		// confirmReq.BankChargeAmount = &req.BankChargeAmount
		confirmReq.ConfirmedAt = time.Now()
		confirmReq.ConfirmedByUserId = 0
		confirmReq.ConfirmedByUsername = "อัตโนมัติ"
		actionErr := s.ConfirmWithdrawTransfer(req.TransId, confirmReq)
		if actionErr != nil {
			return internalServerError(actionErr.Error())
		}
	}
	return nil
}

func (s *bankingService) CreateBonusTransaction(data model.BonusTransactionCreateBody) error {

	member, err := s.repoAccounting.GetUserByMemberCode(data.MemberCode)
	if err != nil {
		fmt.Println(err)
		return badRequest("Invalid Member code")
	}
	bank, err := s.repoAccounting.GetBankByCode(member.BankCode)
	if err != nil {
		fmt.Println(err)
		return badRequest("Invalid User Bank")
	}

	var body model.BonusTransactionCreateBody
	body.MemberCode = *member.MemberCode
	body.UserId = member.Id
	body.TransferType = "bonus"
	body.ToAccountId = 0
	body.ToBankId = bank.Id
	body.ToAccountName = member.Fullname
	body.ToAccountNumber = member.BankAccount
	// body.BeforeAmount = data.BeforeAmount
	// body.AfterAmount = data.AfterAmount
	body.BonusAmount = data.BonusAmount
	body.BonusReason = data.BonusReason
	body.TransferAt = data.TransferAt
	body.CreatedByUserId = data.CreatedByUserId
	body.CreatedByUsername = data.CreatedByUsername
	body.Status = "pending_credit" // no need to confirm

	if err := s.repoBanking.CreateBonusTransaction(body); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *bankingService) UpdateBankTransaction(id int64, req model.BankTransactionUpdateRequest) error {

	trans, err := s.repoAccounting.GetBankTransactionById(id)
	if err != nil {
		return badRequest("Transaction not found")
	}

	fmt.Println("trans", trans)
	// Work by Data
	isTrue := true
	if req.IsAutoCredit != nil {
		if *req.IsAutoCredit == isTrue {
			fmt.Println("req.IsAutoCredit is true", req.IsAutoCredit)
		} else {
			fmt.Println("req.IsAutoCredit is false", req.IsAutoCredit)
		}
		if trans.Status == "pending" || trans.Status == "pending_credit" {
			// var body model.BankTransactionUpdateBody
			// body.IsAutoCredit = *req.IsAutoCredit
			// fmt.Println("body", body)
			data := map[string]interface{}{
				"is_auto_credit": *req.IsAutoCredit,
			}
			if err := s.repoBanking.UpdateBankTransaction(trans.Id, data); err != nil {
				return internalServerError(err.Error())
			}
		} else {
			return badRequest("Transaction is not pending")
		}
	}
	return nil
}

func (s *bankingService) DeleteBankTransaction(id int64) error {

	_, err := s.repoBanking.GetBankTransactionById(id)
	if err != nil {
		return internalServerError(err.Error())
	}

	if err := s.repoBanking.DeleteBankTransaction(id); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *bankingService) GetPendingDepositTransactions(req model.PendingDepositTransactionListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.UnlimitPagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	records, err := s.repoBanking.GetPendingDepositTransactions(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) GetPendingWithdrawTransactions(req model.PendingWithdrawTransactionListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.UnlimitPagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	records, err := s.repoBanking.GetPendingWithdrawTransactions(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) CancelPendingTransaction(id int64, data model.BankTransactionCancelBody) error {

	transaction, err := s.repoBanking.GetBankTransactionById(id)
	if err != nil {
		return internalServerError(err.Error())
	}
	if transaction.Status != "pending" && transaction.Status != "pending_credit" && transaction.Status != "pending_transfer" {
		return badRequest("Transaction is not pending")
	}
	jsonBefore, _ := json.Marshal(transaction)

	var createBody model.CreateBankTransactionActionBody
	if transaction.TransferType == "deposit" {
		createBody.ActionKey = fmt.Sprintf("CANCEL#%d", transaction.Id)
	} else if transaction.TransferType == "withdraw" {
		createBody.ActionKey = fmt.Sprintf("CANCEL#%d", transaction.Id)
	} else {
		// BONUS + GETBACK
		createBody.ActionKey = fmt.Sprintf("CANCEL#%d", transaction.Id)
	}
	createBody.TransactionId = transaction.Id
	createBody.UserId = transaction.UserId
	createBody.TransferType = transaction.TransferType
	createBody.FromAccountId = transaction.FromAccountId
	createBody.ToAccountId = transaction.ToAccountId
	createBody.JsonBefore = string(jsonBefore)
	createBody.TransferAt = transaction.TransferAt
	createBody.ConfirmedAt = data.CanceledAt
	createBody.ConfirmedByUserId = data.CanceledByUserId
	createBody.ConfirmedByUsername = data.CanceledByUsername
	if actionId, err := s.repoBanking.CreateTransactionAction(createBody); err == nil {

		if transaction.TransferType == "deposit" {
			// DO_NOTHING
		} else if transaction.TransferType == "withdraw" {
			// RETURN_CREDIT
			if err := s.increaseMemberCredit(transaction.UserId, transaction.CreditAmount, "withdraw", "คืนเครดิตจากการถอนไม่สำเร็จ"); err != nil {
				if err := s.repoBanking.RollbackTransactionAction(*actionId); err != nil {
					return internalServerError(err.Error())
				}
				return internalServerError(err.Error())
			}
		} else if transaction.TransferType == "bonus" {
			// DO_NOTHING
		} else if transaction.TransferType == "getcreditback" {
			// RETURN_CREDIT
			if err := s.increaseMemberCredit(transaction.UserId, transaction.CreditAmount, "getcreditback", "คืนเครดิตจากรายการที่ไม่สำเร็จ"); err != nil {
				if err := s.repoBanking.RollbackTransactionAction(*actionId); err != nil {
					return internalServerError(err.Error())
				}
				return internalServerError(err.Error())
			}
		}

		if err := s.repoBanking.CancelPendingTransaction(id, data); err != nil {
			if err := s.repoBanking.RollbackTransactionAction(*actionId); err != nil {
				return internalServerError(err.Error())
			}
			return internalServerError(err.Error())
		}

	} else {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *bankingService) ConfirmDepositTransaction(id int64, req model.BankConfirmDepositRequest) error {

	record, err := s.repoBanking.GetBankTransactionById(id)
	if err != nil {
		return internalServerError(err.Error())
	}
	if record.Status != "pending" {
		return badRequest("Transaction is not pending")
	}
	if record.TransferType != "deposit" && record.TransferType != "bonus" {
		return badRequest("Transaction is not deposit")
	}
	jsonBefore, _ := json.Marshal(record)

	var updateData model.BankDepositTransactionConfirmBody
	updateData.Status = "pending_credit"
	updateData.ConfirmedAt = req.ConfirmedAt
	updateData.ConfirmedByUserId = req.ConfirmedByUserId
	updateData.ConfirmedByUsername = req.ConfirmedByUsername
	// updateData.BonusAmount = req.BonusAmount

	var createBody model.CreateBankTransactionActionBody
	createBody.ActionKey = fmt.Sprintf("DCF_STATE#%d", record.Id)
	createBody.TransactionId = record.Id
	createBody.UserId = record.UserId
	createBody.TransferType = record.TransferType
	createBody.FromAccountId = record.FromAccountId
	createBody.ToAccountId = record.ToAccountId
	createBody.JsonBefore = string(jsonBefore)
	if req.TransferAt == nil {
		createBody.TransferAt = record.TransferAt
	} else {
		TransferAt := req.TransferAt
		createBody.TransferAt = TransferAt
		updateData.TransferAt = *TransferAt
	}
	if req.SlipUrl != nil {
		createBody.SlipUrl = *req.SlipUrl
	}
	createBody.CreditAmount = record.CreditAmount
	if req.BonusAmount != nil {
		createBody.BonusAmount = *req.BonusAmount // later: Bonus
	}
	createBody.ConfirmedAt = req.ConfirmedAt
	createBody.ConfirmedByUserId = req.ConfirmedByUserId
	createBody.ConfirmedByUsername = req.ConfirmedByUsername
	if actionId, err := s.repoBanking.CreateTransactionAction(createBody); err == nil {
		// do nothing ?
		if err := s.repoBanking.ConfirmPendingDepositTransaction(id, updateData); err != nil {
			if err := s.repoBanking.RollbackTransactionAction(*actionId); err != nil {
				return internalServerError(err.Error())
			}
			return internalServerError(err.Error())
		}
	} else {
		return internalServerError(err.Error())
	}

	if record.IsAutoCredit {
		// isAUtoCredit with same request
		if err := s.ConfirmDepositCredit(id, req); err != nil {
			return internalServerError(err.Error())
		}
	}
	return nil
}

func (s *bankingService) ConfirmDepositCredit(id int64, req model.BankConfirmDepositRequest) error {

	record, err := s.repoBanking.GetBankTransactionById(id)
	if err != nil {
		return internalServerError(err.Error())
	}
	if record.TransferType != "deposit" && record.TransferType != "bonus" {
		return badRequest("Transaction is not deposit")
	}
	if record.Status != "pending_credit" {
		return badRequest("Transaction is not pending")
	}
	jsonBefore, _ := json.Marshal(record)

	var updateData model.BankDepositTransactionConfirmBody
	updateData.Status = "finished"
	updateData.ConfirmedAt = req.ConfirmedAt
	updateData.ConfirmedByUserId = req.ConfirmedByUserId
	updateData.ConfirmedByUsername = req.ConfirmedByUsername
	if req.BonusAmount != nil {
		updateData.BonusAmount = *req.BonusAmount
		record.BonusAmount = *req.BonusAmount
	}

	var createBody model.CreateBankTransactionActionBody
	createBody.ActionKey = fmt.Sprintf("DCF_CREDIT#%d", record.Id)
	createBody.TransactionId = record.Id
	createBody.UserId = record.UserId
	createBody.TransferType = record.TransferType
	createBody.FromAccountId = record.FromAccountId
	createBody.ToAccountId = record.ToAccountId
	createBody.JsonBefore = string(jsonBefore)
	if req.TransferAt == nil {
		createBody.TransferAt = record.TransferAt
	} else {
		TransferAt := req.TransferAt
		createBody.TransferAt = TransferAt
		updateData.TransferAt = *TransferAt
	}
	if req.SlipUrl != nil {
		createBody.SlipUrl = *req.SlipUrl
	}
	createBody.CreditAmount = record.CreditAmount
	if req.BonusAmount != nil {
		createBody.BonusAmount = *req.BonusAmount
	}
	createBody.ConfirmedAt = req.ConfirmedAt
	createBody.ConfirmedByUserId = req.ConfirmedByUserId
	createBody.ConfirmedByUsername = req.ConfirmedByUsername
	// todo : transaction
	if _, err := s.repoBanking.CreateTransactionAction(createBody); err == nil {
		fmt.Println("ConfirmPendingTransaction updateData:", helper.StructJson(updateData))
		if err := s.increaseMemberCredit(record.UserId, record.CreditAmount, "deposit", "ฝากเครดิต"); err != nil {
			// if err := s.repoBanking.RollbackTransactionAction(*actionId); err != nil {
			// 	return internalServerError(err.Error())
			// }
			return internalServerError(err.Error())
		}
		if record.BonusAmount > 0 {
			if err := s.increaseMemberCredit(record.UserId, record.BonusAmount, "deposit", "ได้รับโบนัสจากการฝากเครดิต"); err != nil {
				// if err := s.repoBanking.RollbackTransactionAction(*actionId); err != nil {
				// 	return internalServerError(err.Error())
				// }
				return internalServerError(err.Error())
			}
		}
		if err := s.repoBanking.ConfirmPendingCreditDepositTransaction(id, updateData); err != nil {
			// if err := s.repoBanking.RollbackTransactionAction(*actionId); err != nil {
			// 	return internalServerError(err.Error())
			// }
			return internalServerError(err.Error())
		}
	} else {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *bankingService) increaseMemberCredit(userId int64, creditAmount float64, statementTypeName string, info string) error {

	statementType, err := s.repoBanking.GetMemberStatementTypeByCode(statementTypeName)
	if err != nil {
		return badRequest("Invalid Type")
	}

	var body model.MemberStatementCreateBody
	body.UserId = userId
	body.StatementTypeId = statementType.Id
	body.Info = info
	body.Amount = creditAmount
	if err := s.repoBanking.IncreaseMemberCredit(body); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *bankingService) decreaseMemberCredit(userId int64, creditAmount float64, statementTypeName string, info string) error {

	statementType, err := s.repoBanking.GetMemberStatementTypeByCode(statementTypeName)
	if err != nil {
		return badRequest("Invalid Type")
	}

	var body model.MemberStatementCreateBody
	body.UserId = userId
	body.StatementTypeId = statementType.Id
	body.Info = info
	body.Amount = creditAmount
	if err := s.repoBanking.DecreaseMemberCredit(body); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *bankingService) ContinueAutoWithdrawTransaction(id int64) error {

	record, err := s.repoBanking.GetBankTransactionById(id)
	if err != nil {
		return internalServerError(err.Error())
	}
	if record.TransferType != "withdraw" && record.TransferType != "getcreditback" {
		return badRequest("Transaction is not withdraw")
	}
	var fromAccountId = record.FromAccountId

	// later : AutoSelect ? checkBalance ?
	var autoWithdrawCondition *model.BankAutoWithdrawCondition
	systemAccount, err := s.repoAccounting.GetWithdrawAccountById(fromAccountId)
	if err == nil {
		if condition, err := s.GetAutoWithdrawCondition(*record, *systemAccount); err == nil {
			autoWithdrawCondition = condition
		}
	}

	if autoWithdrawCondition != nil {
		err := s.SetAutoWithdrawCondition(record.Id, autoWithdrawCondition)
		log.Println(err)
		if err := s.ProcessAutoWithdrawCondition(*autoWithdrawCondition); err != nil {
			return internalServerError(err.Error())
		}
	}
	return nil
}

func (s *bankingService) ConfirmWithdrawTransaction(id int64, req model.BankConfirmCreditWithdrawRequest) error {

	record, err := s.repoBanking.GetBankTransactionById(id)
	if err != nil {
		return internalServerError(err.Error())
	}
	if record.TransferType != "withdraw" && record.TransferType != "getcreditback" {
		return badRequest("Transaction is not withdraw")
	}
	var fromAccountId = record.FromAccountId
	var updateData model.BankWithdrawTransactionConfirmBody
	if record.TransferType == "getcreditback" {
		updateData.Status = "finished"
	} else {
		updateData.Status = "pending_transfer"
	}
	updateData.ConfirmedAt = req.ConfirmedAt
	updateData.ConfirmedByUserId = req.ConfirmedByUserId
	updateData.ConfirmedByUsername = req.ConfirmedByUsername
	if record.Status == "pending" || record.Status == "pending_credit" {
		if req.FromAccountId != nil {
			fromAccount, err := s.repoAccounting.GetBankAccountById(*req.FromAccountId)
			if err != nil {
				return badRequest("Invalid Bank Account")
			}
			fromAccountId = fromAccount.Id
			updateData.FromAccountId = &fromAccount.Id
		}
		if req.CreditAmount != nil {
			updateData.CreditAmount = *req.CreditAmount
		} else {
			updateData.CreditAmount = record.CreditAmount
		}
		if req.BankChargeAmount != nil {
			updateData.BankChargeAmount = *req.BankChargeAmount
		}
	} else {
		return badRequest("Transaction is not pending")
	}
	jsonBefore, _ := json.Marshal(record)

	// Check Credit/Balance
	if err := s.repoBanking.CheckMemeberHasEnoughtCredit(record.UserId, record.CreditAmount); err != nil {
		return internalServerError(err.Error())
	}

	// later : AutoSelect ? checkBalance ?
	var autoWithdrawCondition *model.BankAutoWithdrawCondition
	systemAccount, err := s.repoAccounting.GetWithdrawAccountById(fromAccountId)
	if err == nil {
		if condition, err := s.GetAutoWithdrawCondition(*record, *systemAccount); err == nil {
			autoWithdrawCondition = condition
		}
	}

	var createBody model.CreateBankTransactionActionBody
	createBody.ActionKey = fmt.Sprintf("CFW_CREDIT#%d", record.Id)
	createBody.TransactionId = record.Id
	createBody.UserId = record.UserId
	createBody.TransferType = record.TransferType
	createBody.FromAccountId = record.FromAccountId
	createBody.ToAccountId = record.ToAccountId
	createBody.JsonBefore = string(jsonBefore)
	createBody.TransferAt = record.TransferAt
	createBody.CreditAmount = updateData.CreditAmount
	createBody.BankChargeAmount = updateData.BankChargeAmount
	createBody.ConfirmedAt = req.ConfirmedAt
	createBody.ConfirmedByUserId = req.ConfirmedByUserId
	createBody.ConfirmedByUsername = req.ConfirmedByUsername
	// BOF : transaction
	if actionId, err := s.repoBanking.CreateTransactionAction(createBody); err == nil {
		if err := s.decreaseMemberCredit(record.UserId, record.CreditAmount, "withdraw", "ถอนเครดิต"); err != nil {
			if err := s.repoBanking.RollbackTransactionAction(*actionId); err != nil {
				return internalServerError(err.Error())
			}
			return internalServerError(err.Error())
		}
		if err := s.repoBanking.ConfirmPendingWithdrawTransaction(id, updateData); err != nil {
			if err := s.repoBanking.RollbackTransactionAction(*actionId); err != nil {
				return internalServerError(err.Error())
			}
			return internalServerError(err.Error())
		}
		// EOF : transaction
		if autoWithdrawCondition != nil {
			err := s.SetAutoWithdrawCondition(record.Id, autoWithdrawCondition)
			log.Println(err)
			if err := s.ProcessAutoWithdrawCondition(*autoWithdrawCondition); err != nil {
				return internalServerError(err.Error())
			}
		}
	} else {
		return internalServerError(err.Error())
	}

	return nil
}

func (s *bankingService) ConfirmWithdrawTransfer(id int64, req model.BankConfirmTransferWithdrawRequest) error {

	record, err := s.repoBanking.GetBankTransactionById(id)
	if err != nil {
		return internalServerError(err.Error())
	}
	if record.TransferType != "withdraw" {
		return badRequest("Transaction is not withdraw")
	}
	var fromAccountId = record.FromAccountId
	var updateData model.BankWithdrawTransactionConfirmBody
	updateData.Status = "finished"
	updateData.ConfirmedAt = req.ConfirmedAt
	updateData.ConfirmedByUserId = req.ConfirmedByUserId
	updateData.ConfirmedByUsername = req.ConfirmedByUsername
	if record.Status == "pending_transfer" {
		if req.FromAccountId != nil {
			fromAccount, err := s.repoAccounting.GetBankAccountById(*req.FromAccountId)
			if err != nil {
				return badRequest("Invalid Bank Account")
			}
			fromAccountId = fromAccount.Id
			updateData.FromAccountId = &fromAccountId
		}
		if req.BankChargeAmount != nil {
			updateData.BankChargeAmount = *req.BankChargeAmount
		}
	} else {
		return badRequest("Transaction is not transfer pending")
	}
	jsonBefore, _ := json.Marshal(record)

	// later : AutoSelectAccount ? checkBalance ?
	systemAccount, err := s.repoAccounting.GetWithdrawAccountById(fromAccountId)
	if err != nil {
		if err.Error() == recordNotFound {
			return notFound("Bank Account not found")
		}
		return internalServerError(err.Error())
	}

	var createBody model.CreateBankTransactionActionBody
	createBody.ActionKey = fmt.Sprintf("CFW_TRASFER#%d", record.Id)
	createBody.TransactionId = record.Id
	createBody.UserId = record.UserId
	createBody.TransferType = record.TransferType
	createBody.FromAccountId = record.FromAccountId
	createBody.ToAccountId = record.ToAccountId
	createBody.JsonBefore = string(jsonBefore)
	createBody.TransferAt = record.TransferAt
	createBody.CreditAmount = record.CreditAmount
	if req.BankChargeAmount != nil {
		createBody.BankChargeAmount = *req.BankChargeAmount
	}
	createBody.ConfirmedAt = req.ConfirmedAt
	createBody.ConfirmedByUserId = req.ConfirmedByUserId
	createBody.ConfirmedByUsername = req.ConfirmedByUsername
	// BOF : transaction
	if actionId, err := s.repoBanking.CreateTransactionAction(createBody); err == nil {
		allow_withdraw_from_account := false
		allow_withdraw_from_account_key := "allow_withdraw_from_account"
		var query model.BotAccountConfigListRequest
		query.SearchKey = &allow_withdraw_from_account_key
		config, _ := s.repoAccounting.GetBotaccountConfigByKey(query)
		// later : more condition
		if config != nil && config.ConfigVal == "all" {
			allow_withdraw_from_account = true
		}
		if allow_withdraw_from_account {
			// ExternalTransfer
			var body model.ExternalAccountTransferBody
			body.AccountForm = systemAccount.AccountNumber
			body.AccountTo = record.ToAccountNumber
			body.Amount = strconv.FormatFloat(record.CreditAmount, 'E', -1, 64)
			body.BankCode = record.ToBankCode
			body.Pin = systemAccount.PinCode
			if err := s.repoBanking.TransferExternalAccount(body); err != nil {
				if err := s.repoBanking.RollbackTransactionAction(*actionId); err != nil {
					return internalServerError(err.Error())
				}
				return internalServerError(err.Error())
			}
			// later : read from FASTBANK reponse, updateData.TransferAt = time.Now()
		}
		updateData.TransferAt = time.Now()
		if err := s.repoBanking.ConfirmPendingWithdrawTransfer(id, updateData); err != nil {
			if err := s.repoBanking.RollbackTransactionAction(*actionId); err != nil {
				return internalServerError(err.Error())
			}
			return internalServerError(err.Error())
		}
		// EOF : commit transaction
	} else {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *bankingService) GetFinishedTransactions(req model.FinishedTransactionListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.UnlimitPagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	records, err := s.repoBanking.GetFinishedTransactions(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) RemoveFinishedTransaction(id int64, data model.BankTransactionRemoveBody) error {

	record, err := s.repoBanking.GetBankTransactionById(id)
	if err != nil {
		return internalServerError(err.Error())
	}
	if record.Status != "finished" {
		return badRequest("Transaction is not finished")
	}

	if err := s.repoBanking.RemoveFinishedTransaction(id, data); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *bankingService) GetRemovedTransactions(req model.RemovedTransactionListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.UnlimitPagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	records, err := s.repoBanking.GetRemovedTransactions(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) GetMemberByCode(code string) (*model.Member, error) {

	if code == "" {
		return nil, badRequest("Code is required")
	}

	records, err := s.repoBanking.GetMemberByCode(code)
	if err != nil {
		if err.Error() == recordNotFound {
			return nil, notFound(memberNotFound)
		}
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) GetMembers(req model.MemberListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.UnlimitPagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	records, err := s.repoBanking.GetMembers(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) GetPossibleStatementOwners(req model.MemberPossibleListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.UnlimitPagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}

	statement, err := s.repoBanking.GetBankStatementById(req.UnknownStatementId)
	if err != nil {
		if err.Error() == recordNotFound {
			return nil, notFound(bankStatementferNotFound)
		}
		return nil, internalServerError(err.Error())
	}
	req.UserBankCode = &statement.FromBankCode
	req.UserAccountNumber = &statement.FromAccountNumber

	records, err := s.repoBanking.GetPossibleStatementOwners(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) GetMemberTransactions(req model.MemberTransactionListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.UnlimitPagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	records, err := s.repoBanking.GetMemberTransactions(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) GetMemberTransactionSummary(req model.MemberTransactionListRequest) (*model.MemberTransactionSummary, error) {

	result, err := s.repoBanking.GetMemberTransactionSummary(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return result, nil
}

func (s *bankingService) MatchDepositTransaction(id int64, req model.BankConfirmDepositRequest) error {

	record, err := s.repoBanking.GetBankTransactionById(id)
	if err != nil {
		return internalServerError(err.Error())
	}
	if record.Status != "pending" {
		return badRequest("Transaction is not pending")
	}
	if record.TransferType != "deposit" && record.TransferType != "bonus" {
		return badRequest("Transaction is not deposit")
	}
	if err := s.ConfirmDepositTransaction(record.UserId, req); err != nil {
		return internalServerError(err.Error())
	}
	// // no need if IsAutoCredit = true if err := s.ConfirmDepositCredit(record.UserId, req); err != nil {
	// 	return internalServerError(err.Error())
	// }
	return nil
}

func (s *bankingService) GetMemberStatementTypeByCode(name string) (*model.MemberStatementType, error) {

	record, err := s.repoBanking.GetMemberStatementTypeByCode(name)
	if err != nil {
		if err.Error() == recordNotFound {
			return nil, notFound(bankStatementferNotFound)
		}
		return nil, internalServerError(err.Error())
	}
	return record, nil
}

func (s *bankingService) GetMemberStatementTypes(req model.SimpleListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.UnlimitPagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	records, err := s.repoBanking.GetMemberStatementTypes(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) GetMemberStatementById(req model.GetByIdRequest) (*model.MemberStatementResponse, error) {

	record, err := s.repoBanking.GetMemberStatementById(req.Id)
	if err != nil {
		if err.Error() == recordNotFound {
			return nil, notFound(bankStatementferNotFound)
		}
		return nil, internalServerError(err.Error())
	}
	return record, nil
}

func (s *bankingService) GetMemberStatements(req model.MemberStatementListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.UnlimitPagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	records, err := s.repoBanking.GetMemberStatements(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) ProcessMemberDepositCredit(userId int64, amount float64) error {

	statementCode := "deposit"
	var body model.MemberStatementCreateBody

	member, err := s.repoBanking.GetMemberById(userId)
	if err != nil {
		return badRequest("Invalid Member")
	}

	statementType, err := s.repoBanking.GetMemberStatementTypeByCode(statementCode)
	if err != nil {
		return badRequest("Invalid Type")
	}

	// todo : more validate

	body.UserId = member.Id
	body.StatementTypeId = statementType.Id
	body.Info = "ฝากเครดิต"
	body.Amount = amount
	if err := s.repoBanking.IncreaseMemberCredit(body); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *bankingService) ProcessMemberWithdrawCredit(userId int64, amount float64) error {

	statementCode := "withdraw"
	var body model.MemberStatementCreateBody

	member, err := s.repoBanking.GetMemberById(userId)
	if err != nil {
		return badRequest("Invalid Member")
	}

	statementType, err := s.repoBanking.GetMemberStatementTypeByCode(statementCode)
	if err != nil {
		return badRequest("Invalid Type")
	}

	// todo : more validate

	body.UserId = member.Id
	body.StatementTypeId = statementType.Id
	body.Info = "ถอนเครดิต"
	body.Amount = amount
	if err := s.repoBanking.DecreaseMemberCredit(body); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *bankingService) ProcessMemberBonusCredit(userId int64, amount float64) error {

	statementCode := "bonus"
	var body model.MemberStatementCreateBody

	member, err := s.repoBanking.GetMemberById(userId)
	if err != nil {
		return badRequest("Invalid Member")
	}

	statementType, err := s.repoBanking.GetMemberStatementTypeByCode(statementCode)
	if err != nil {
		return badRequest("Invalid Type")
	}

	// todo : more validate

	body.UserId = member.Id
	body.StatementTypeId = statementType.Id
	body.Info = "ได้รับโบนัส"
	body.Amount = amount
	if err := s.repoBanking.IncreaseMemberCredit(body); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *bankingService) ProcessMemberGetbackCredit(userId int64, amount float64) error {

	statementCode := "getcreditback"
	var body model.MemberStatementCreateBody

	member, err := s.repoBanking.GetMemberById(userId)
	if err != nil {
		return badRequest("Invalid Member")
	}

	statementType, err := s.repoBanking.GetMemberStatementTypeByCode(statementCode)
	if err != nil {
		return badRequest("Invalid Type")
	}

	// todo : more validate

	body.UserId = member.Id
	body.StatementTypeId = statementType.Id
	body.Info = "ถูกดึงเครดิตคืน"
	body.Amount = amount
	if err := s.repoBanking.DecreaseMemberCredit(body); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}
