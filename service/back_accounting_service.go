package service

import (
	"bytes"
	"cybergame-api/helper"
	"cybergame-api/model"
	"cybergame-api/repository"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type AccountingService interface {
	CheckCurrentAdminId(input any) (*int64, error)
	CheckCurrentUsername(input any) (*string, error)
	CheckConfirmationPassword(req model.ConfirmRequest) (*bool, error)

	GetBanks(req model.BankListRequest) (*model.SuccessWithPagination, error)
	GetAccountTypes(req model.AccountTypeListRequest) (*model.SuccessWithPagination, error)

	GetBankAccountById(req model.BankGetByIdRequest) (*model.BankAccount, error)
	GetBankAccounts(req model.BankAccountListRequest) (*model.SuccessWithPagination, error)
	GetBankAccountPriorities() (*model.SuccessWithPagination, error)
	CreateBankAccount(body model.BankAccountCreateBody) error
	UpdateBankAccount(id int64, req model.BankAccountUpdateRequest) error
	DeleteBankAccount(id int64) error

	GetTransactionById(req model.BankGetByIdRequest) (*model.BankAccountTransaction, error)
	GetTransactions(req model.BankAccountTransactionListRequest) (*model.SuccessWithPagination, error)
	CreateTransaction(body model.BankAccountTransactionBody) error
	UpdateTransaction(id int64, body model.BankAccountTransactionBody) error
	DeleteTransaction(id int64) error

	GetTransferById(req model.BankGetByIdRequest) (*model.BankAccountTransfer, error)
	GetTransfers(req model.BankAccountTransferListRequest) (*model.SuccessWithPagination, error)
	CreateTransfer(body model.BankAccountTransferBody) error
	ConfirmTransfer(id int64, actorId int64) error
	DeleteTransfer(id int64) error

	GetAccountStatements(req model.BankAccountStatementListRequest) (*model.SuccessWithPagination, error)
	GetAccountStatementById(req model.BankGetByIdRequest) (*model.BankStatement, error)
	AddAccountStatementToWebhook(req model.RecheckWebhookRequest) error

	GetExternalSettings() (*model.ExternalSettings, error)
	GetCustomerAccountsInfo(req model.CustomerAccountInfoRequest) (*model.CustomerAccountInfo, error)
	GetExternalAccounts() (*model.SuccessWithPagination, error)
	GetExternalAccountBalance(req model.ExternalAccountStatusRequest) (*model.ExternalAccountBalance, error)
	GetExternalAccountStatus(req model.ExternalAccountStatusRequest) (*model.ExternalAccountStatus, error)
	CreateExternalAccount(body model.ExternalAccountCreateBody) (*model.ExternalAccountCreateResponse, error)
	UpdateExternalAccount(body model.ExternalAccountUpdateBody) (*model.ExternalAccountCreateResponse, error)
	EnableExternalAccount(req model.ExternalAccountEnableRequest) (*model.ExternalAccountStatus, error)
	DeleteExternalAccount(req model.ExternalAccountStatusRequest) error
	TransferExternalAccount(req model.ExternalAccountTransferRequest) error
	CreateBankStatementFromWebhook(data model.WebhookStatement) error
	CreateBotaccountConfig(body model.BotAccountConfigCreateBody) error

	GetExternalAccountLogs(req model.ExternalStatementListRequest) (*model.SuccessWithPagination, error)
	GetExternalAccountStatements(req model.ExternalStatementListRequest) (*model.SuccessWithPagination, error)
	GetExternalAccountStatementByTimestamp(req model.ExternalStatementListRequest) (*model.SuccessWithPagination, error)
	CreateWebhookLog(logType string, jsonRequest string) (*int64, error)
	SetSuccessWebhookLog(id int64, jsonPayload string) error
	SetFailedWebhookLog(id int64, logStatus string) error
}

type accountingService struct {
	repo repository.AccountingRepository
}

var invalidConfirmation = "Invalid confirmation password"
var invalidCurrentAdminId = "Invalid current user id"

var recordNotFound = "record not found"
var bankNotFound = "Bank not found"
var bankAccountNotFound = "Account not found"
var transactionNotFound = "Transsaction not found"
var transferNotFound = "Transfer not found"

func NewAccountingService(
	repo repository.AccountingRepository,
) AccountingService {
	return &accountingService{repo}
}

func (s *accountingService) CheckCurrentAdminId(input any) (*int64, error) {

	// input := c.MustGet("adminId")
	if input == nil {
		return nil, badRequest(invalidCurrentAdminId)
	}
	var adminId = int64(input.(float64))
	if adminId <= 0 {
		return nil, badRequest(invalidCurrentAdminId)
	}
	return &adminId, nil
}

func (s *accountingService) CheckCurrentUsername(input any) (*string, error) {

	// input := c.MustGet("username")
	if input == nil {
		return nil, badRequest(invalidCurrentAdminId)
	}
	var username = input.(string)
	// if username == "" {
	// 	return nil, badRequest(invalidCurrentAdminId)
	// }
	return &username, nil
}

func (s *accountingService) CheckConfirmationPassword(req model.ConfirmRequest) (*bool, error) {

	user, err := s.repo.GetAdminById(req.UserId)
	if err != nil {
		fmt.Println(req)
		return nil, notFound(invalidConfirmation)
	}
	if user == nil {
		return nil, badRequest(invalidConfirmation)
	}
	if err := helper.CompareAdminPassword(req.Password, user.Password); err != nil {
		return nil, badRequest(invalidConfirmation)
	}
	token := true
	return &token, nil
}

func (s *accountingService) GetBanks(req model.BankListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.UnlimitPagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}

	records, err := s.repo.GetBanks(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *accountingService) GetAccountTypes(req model.AccountTypeListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.UnlimitPagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}

	records, err := s.repo.GetAccountTypes(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *accountingService) GetBankAccountById(req model.BankGetByIdRequest) (*model.BankAccount, error) {

	err := s.UpdateBankAccountBalanceById(req.Id)
	log.Panicln(err)
	// if err != nil {
	// return nil, internalServerError(err.Error())
	// }

	record, err := s.repo.GetBankAccountById(req.Id)
	if err != nil {
		if err.Error() == recordNotFound {
			return nil, notFound(bankAccountNotFound)
		}
		return nil, internalServerError(err.Error())
	}
	return record, nil
}

func (s *accountingService) GetBankAccounts(req model.BankAccountListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.UnlimitPagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}

	err := s.UpdateAllBankAccountBotStatus()
	log.Println(err)

	list, err := s.repo.GetBankAccounts(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return list, nil
}

func (s *accountingService) GetBankAccountPriorities() (*model.SuccessWithPagination, error) {

	list, err := s.repo.GetBankAccountPriorities()
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return list, nil
}

func (s *accountingService) CreateBankAccount(body model.BankAccountCreateBody) error {

	bank, err := s.repo.GetBankById(body.BankId)
	if err != nil {
		fmt.Println(err)
		if err.Error() == recordNotFound {
			return notFound(bankNotFound)
		}
		return badRequest("Invalid Bank")
	}

	accountType, err := s.repo.GetAccounTypeById(body.AccountTypeId)
	if err != nil {
		fmt.Println(err)
		return badRequest("Invalid Account Type")
	}

	acNo := helper.StripAllButNumbers(body.AccountNumber)
	exist, err := s.repo.HasBankAccount(acNo)
	if err != nil {
		fmt.Println(err)
		return internalServerError(err.Error())
	}
	if exist {
		return badRequest("Account already exist")
	}

	var createBody model.BankAccountCreateBody
	createBody.BankId = bank.Id
	createBody.AccountTypeId = accountType.Id
	createBody.AccountName = body.AccountName
	createBody.AccountNumber = acNo
	createBody.AccountPriorityId = body.AccountPriorityId
	createBody.AutoCreditFlag = body.AutoCreditFlag
	createBody.IsMainWithdraw = body.IsMainWithdraw
	createBody.AutoWithdrawFlag = body.AutoWithdrawFlag
	createBody.AutoWithdrawCreditFlag = body.AutoWithdrawCreditFlag
	createBody.AutoWithdrawConfirmFlag = body.AutoWithdrawConfirmFlag
	createBody.AutoTransferMaxAmount = body.AutoTransferMaxAmount
	createBody.AutoWithdrawMaxAmount = body.AutoWithdrawMaxAmount
	createBody.DeviceUid = body.DeviceUid
	// อัพเดทหลังจากเรียกบอท createBody.PinCode = data.PinCode
	createBody.QrWalletStatus = body.QrWalletStatus
	createBody.AccountStatus = body.AccountStatus
	createBody.AccountBalance = 0
	createBody.ConnectionStatus = "disconnected"
	if err := s.repo.CreateBankAccount(createBody); err != nil {
		return internalServerError(err.Error())
	}

	// allowCreateExternalAccount := false
	// config, _ := s.GetExternalAccountConfig("allow_create_external_account")
	// if config != nil {
	// 	if config.ConfigVal == "list" {
	// 		accountConfig, errConfig := s.HasExternalAccountConfig("allow_external_account_number", acNo)
	// 		if errConfig != nil {
	// 			return nil
	// 		}
	// 		if accountConfig.ConfigVal == acNo {
	// 			allowCreateExternalAccount = true
	// 		}
	// 	} else if config.ConfigVal == "all" {
	// 		allowCreateExternalAccount = true
	// 	}
	// }

	if s.IsAllowCreateExternalAccount(acNo) && body.DeviceUid != "" && body.PinCode != "" && !s.HasExternalAccount(acNo) {
		// if _, err := s.HasExternalAccountConfig("allow_external_account_number", acNo); err != nil {
		// 	return nil
		// }

		// FASTBANK
		var createExternalBody model.ExternalAccountCreateBody
		createExternalBody.AccountNo = acNo
		createExternalBody.BankCode = bank.Code
		createExternalBody.DeviceId = body.DeviceUid
		// ไม่ได้ใช้ createExternalBody.Password = data.Password
		createExternalBody.Pin = body.PinCode
		// ไม่ได้ใช้ createExternalBody.Username = data.Username
		// ไม่ได้ใช้ createExternalBody.WebhookNotifyUrl = os.Getenv("ACCOUNTING_LOCAL_WEBHOOK_ENDPOINT") + "/api/accounting/webhooks/noti"
		createExternalBody.WebhookUrl = os.Getenv("ACCOUNTING_LOCAL_WEBHOOK_ENDPOINT") + "/api/accounting/webhooks/action"
		if createResp, err := s.CreateExternalAccount(createExternalBody); err != nil {
			_, webhookErr := s.CreateWebhookLog("CreateBankAccount.CreateExternalAccount, ERROR", helper.StructJson(struct {
				data model.BankAccountCreateBody
				err  error
			}{body, err}))
			log.Panicln(webhookErr)
			return internalServerError(err.Error())
		} else {
			// Update EncryptionPin
			account, err := s.repo.GetBankAccountByAccountNumber(acNo)
			if err != nil {
				_, webhookErr := s.CreateWebhookLog("CreateBankAccount.GetBankAccountByAccountNumber, ERROR", helper.StructJson(struct {
					data model.BankAccountCreateBody
					err  error
				}{body, err}))
				log.Panicln(webhookErr)
				return internalServerError(err.Error())
			}
			var updateBody model.BankAccountUpdateBody
			updateBody.PinCode = &createResp.Pin
			updateBody.ExternalId = &createResp.Id
			if err := s.repo.UpdateBankAccount(account.Id, updateBody); err != nil {
				_, webhookErr := s.CreateWebhookLog("CreateBankAccount.UpdateBankAccount, ERROR", helper.StructJson(struct {
					data model.BankAccountUpdateBody
					err  error
				}{updateBody, err}))
				log.Panicln(webhookErr)
				return internalServerError(err.Error())
			}
		}
	}

	return nil
}

func (s *accountingService) UpdateBankAccount(id int64, req model.BankAccountUpdateRequest) error {

	account, err := s.repo.GetBankAccountById(id)
	if err != nil {
		return internalServerError(err.Error())
	}

	var updateBody model.BankAccountUpdateBody
	var updateExBody model.ExternalAccountUpdateBody
	onExternalChange := false
	isAllowNoMainWithdraw := true
	onMainWithdrawChange := false

	// Validate
	if req.BankId != nil && account.BankId != *req.BankId {
		bank, err := s.repo.GetBankById(*req.BankId)
		if err != nil {
			fmt.Println(err)
			if err.Error() == recordNotFound {
				return notFound(bankNotFound)
			}
			return badRequest("Invalid Bank")
		}
		updateBody.BankId = &bank.Id
		// onExternalChange = true
		updateExBody.BankCode = bank.Code
	}
	if req.AccountTypeId != nil && account.AccountTypeId != *req.AccountTypeId {
		accountType, err := s.repo.GetAccounTypeById(*req.AccountTypeId)
		if err != nil {
			fmt.Println(err)
			return badRequest("Invalid Account Type")
		}
		updateBody.AccountTypeId = &accountType.Id
	}
	if req.AccountName != nil && account.AccountName != *req.AccountName {
		updateBody.AccountName = req.AccountName
	}
	if req.AccountNumber != nil && account.AccountNumber != *req.AccountNumber {
		acNo := helper.StripAllButNumbers(*req.AccountNumber)
		if acNo != "" {
			check, err := s.repo.HasBankAccount(acNo)
			if err != nil {
				return internalServerError(err.Error())
			}
			if !check {
				fmt.Println(acNo)
				return notFound("Account already exist")
			}
			updateBody.AccountNumber = &acNo
			// onExternalChange = true
			account.AccountNumber = acNo
		} else {
			updateBody.AccountNumber = &account.AccountNumber
		}
	}
	if req.DeviceUid != nil && account.DeviceUid != *req.DeviceUid {
		updateBody.DeviceUid = req.DeviceUid
		// onExternalChange = true
		updateExBody.DeviceId = &account.DeviceUid
	}
	if req.PinCode != nil {
		// updateBody.PinCode = req.PinCode
		onExternalChange = true
		updateExBody.Pin = req.PinCode
	}
	if req.AutoCreditFlag != nil && account.AutoCreditFlag != *req.AutoCreditFlag {
		updateBody.AutoCreditFlag = req.AutoCreditFlag
	}
	if req.IsMainWithdraw != nil && account.IsMainWithdraw != *req.IsMainWithdraw {
		if isAllowNoMainWithdraw {
			updateBody.IsMainWithdraw = req.IsMainWithdraw
			onMainWithdrawChange = true
		} else {
			if *req.IsMainWithdraw {
				// reset all
				updateBody.IsMainWithdraw = req.IsMainWithdraw
				onMainWithdrawChange = true
			} else {
				// cant set to false if no other main account
				onMainWithdrawChange = false
			}
		}
	}
	if req.AutoWithdrawFlag != nil && account.AutoWithdrawFlag != *req.AutoWithdrawFlag {
		updateBody.AutoWithdrawFlag = req.AutoWithdrawFlag
	}
	if req.AutoWithdrawCreditFlag != nil && account.AutoWithdrawCreditFlag != *req.AutoWithdrawCreditFlag {
		updateBody.AutoWithdrawCreditFlag = req.AutoWithdrawCreditFlag
	}
	if req.AutoWithdrawConfirmFlag != nil && account.AutoWithdrawConfirmFlag != *req.AutoWithdrawConfirmFlag {
		updateBody.AutoWithdrawConfirmFlag = req.AutoWithdrawConfirmFlag
	}
	if req.AutoWithdrawMaxAmount != nil && account.AutoWithdrawMaxAmount != *req.AutoWithdrawMaxAmount {
		updateBody.AutoWithdrawMaxAmount = req.AutoWithdrawMaxAmount
	}
	if req.AutoTransferMaxAmount != nil && account.AutoTransferMaxAmount != *req.AutoTransferMaxAmount {
		updateBody.AutoTransferMaxAmount = req.AutoTransferMaxAmount
	}
	if req.AccountPriorityId != nil && account.AccountPriorityId != *req.AccountPriorityId {
		updateBody.AccountPriorityId = req.AccountPriorityId
	}
	if req.QrWalletStatus != nil && account.QrWalletStatus != *req.QrWalletStatus {
		updateBody.QrWalletStatus = req.QrWalletStatus
	}
	if req.AccountStatus != nil && account.AccountStatus != *req.AccountStatus {
		updateBody.AccountStatus = req.AccountStatus
	}

	if onExternalChange && s.IsAllowCreateExternalAccount(account.AccountNumber) && updateExBody.DeviceId != nil && updateExBody.Pin != nil {
		// if _, err := s.HasExternalAccountConfig("allow_external_account_number", acNo); err != nil {
		// 	return nil
		// }

		if updateExBody.DeviceId == nil {
			updateExBody.DeviceId = &account.DeviceUid
		}
		// if updateExBody.Pin == nil {
		// 	updateExBody.Pin = &account.PinCode
		// }

		// Create if not exist
		if !s.HasExternalAccount(account.AccountNumber) {
			var createExternalBody model.ExternalAccountCreateBody
			createExternalBody.AccountNo = account.AccountNumber
			createExternalBody.BankCode = account.BankCode
			createExternalBody.DeviceId = *updateExBody.DeviceId
			// ไม่ได้ใช้ createExternalBody.Password = data.Password
			createExternalBody.Pin = *updateExBody.Pin
			// ไม่ได้ใช้ createExternalBody.Username = data.Username
			// ไม่ได้ใช้ createExternalBody.WebhookNotifyUrl = os.Getenv("ACCOUNTING_LOCAL_WEBHOOK_ENDPOINT") + "/api/accounting/webhooks/noti"
			createExternalBody.WebhookUrl = os.Getenv("ACCOUNTING_LOCAL_WEBHOOK_ENDPOINT") + "/api/accounting/webhooks/action"
			if createResp, err := s.CreateExternalAccount(createExternalBody); err != nil {
				_, webhookErr := s.CreateWebhookLog("UpdateBankAccount.CreateExternalAccount, ERROR", helper.StructJson(struct {
					req model.BankAccountUpdateRequest
					err error
				}{req, err}))
				log.Panicln(webhookErr)
				return internalServerError(err.Error())
			} else {
				// Update EncryptionPin
				updateBody.PinCode = &createResp.Pin
				updateBody.ExternalId = &createResp.Id
			}
		} else {
			updateExBody.AccountNo = account.AccountNumber
			// ไม่ได้ใช้ updateExBody.WebhookNotifyUrl = os.Getenv("ACCOUNTING_LOCAL_WEBHOOK_ENDPOINT") + "/api/accounting/webhooks/noti"
			updateExBody.WebhookUrl = os.Getenv("ACCOUNTING_LOCAL_WEBHOOK_ENDPOINT") + "/api/accounting/webhooks/action"
			if externalCreateResp, err := s.UpdateExternalAccount(updateExBody); err != nil {
				_, webhookErr := s.CreateWebhookLog("UpdateBankAccount, ERROR", helper.StructJson(struct {
					id  int64
					req model.BankAccountUpdateRequest
					err error
				}{id, req, err}))
				log.Panicln(webhookErr)
				return internalServerError(err.Error())
			} else {
				// Update EncryptionPin
				updateBody.PinCode = &externalCreateResp.Pin
				updateBody.ExternalId = &externalCreateResp.Id
			}
		}
	}

	if onMainWithdrawChange {
		if err := s.repo.ResetMainWithdrawBankAccount(); err != nil {
			return internalServerError(err.Error())
		}
	}

	if err := s.repo.UpdateBankAccount(account.Id, updateBody); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) UpdateBankAccountBalanceById(id int64) error {

	account, err := s.repo.GetBankAccountById(id)
	if err != nil {
		return internalServerError(err.Error())
	}
	now := time.Now()
	if account.LastConnUpdateAt != nil {
		if now.Sub(*account.LastConnUpdateAt).Seconds() < 30 {
			return nil
		}
	}

	statusActiveFlag := "active"
	statusDisconnectedFlag := "disconnected"
	var updateData model.BankAccountUpdateBody
	updateData.LastConnUpdateAt = &now
	updateData.ConnectionStatus = &statusDisconnectedFlag

	if account.DeviceUid != "" && account.PinCode != "" {
		var query model.ExternalAccountStatusRequest
		query.AccountNumber = account.AccountNumber
		statusResp, err := s.GetExternalAccountStatus(query)
		if err != nil {
			return internalServerError(err.Error())
		}
		if statusResp.Status == "online" {
			updateData.ConnectionStatus = &statusActiveFlag
		} else {
			fmt.Println("statusResp", statusResp)
			updateData.ConnectionStatus = &statusDisconnectedFlag
		}
		balaceResp, err := s.GetExternalAccountBalance(query)
		if err != nil {
			return internalServerError(err.Error())
		}
		if balaceResp.AccountNo == account.AccountNumber {
			balance, _ := strconv.ParseFloat(strings.TrimSpace(balaceResp.AccountBalance), 64)
			updateData.AccountBalance = &balance
		} else {
			fmt.Println("ERROR, balaceResp: ", balaceResp)
			return internalServerError(err.Error())
		}
	}

	if err := s.repo.UpdateBankAccount(id, updateData); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) UpdateAllBankAccountBotStatus() error {

	var query model.BankAccountListRequest
	query.Limit = 100
	query.Page = 0
	accounts, err := s.repo.GetBotBankAccounts(query)
	if err != nil {
		return internalServerError(err.Error())
	}
	now := time.Now()
	errorMinDelay := time.Now().Add(time.Minute * 5)
	statusActiveFlag := "active"
	statusDisconnectedFlag := "disconnected"
	for _, account := range accounts.List.([]model.BankAccountResponse) {

		if account.LastConnUpdateAt != nil {
			if now.Sub(*account.LastConnUpdateAt).Seconds() < 30 {
				continue
			}
		}
		var data model.BankAccountUpdateBody
		data.LastConnUpdateAt = &now
		data.ConnectionStatus = &statusDisconnectedFlag
		// data.AccountBalance = 0

		// FASTBANK
		var query model.ExternalAccountStatusRequest
		query.AccountNumber = account.AccountNumber
		statusResp, err := s.GetExternalAccountStatus(query)
		if err != nil {
			data.LastConnUpdateAt = &errorMinDelay
			// fmt.Println("ERROR", err.Error())
		} else {
			if statusResp.Status == "online" {
				data.ConnectionStatus = &statusActiveFlag
			} else {
				fmt.Println("statusResp", statusResp)
				data.ConnectionStatus = &statusDisconnectedFlag
			}
		}

		balaceResp, err := s.GetExternalAccountBalance(query)
		if err != nil {
			data.LastConnUpdateAt = &errorMinDelay
			// fmt.Println("ERROR", err.Error())
		} else {
			if balaceResp.AccountNo == account.AccountNumber {
				balance, _ := strconv.ParseFloat(strings.TrimSpace(balaceResp.AccountBalance), 64)
				data.AccountBalance = &balance
			} else {
				data.LastConnUpdateAt = &errorMinDelay
				// fmt.Println("ERROR, balaceResp: ", balaceResp)
			}
		}

		if err := s.repo.UpdateBankAccount(account.Id, data); err != nil {
			fmt.Println("ERROR, UPDATE ", err.Error())
		}
	}

	return nil
}

func (s *accountingService) DeleteBankAccount(id int64) error {

	account, err := s.repo.GetBankAccountById(id)
	if err != nil {
		return internalServerError(err.Error())
	}

	if account.ExternalId != "" && s.HasExternalAccount(account.AccountNumber) {
		var query model.ExternalAccountStatusRequest
		query.AccountNumber = account.AccountNumber
		if err := s.DeleteExternalAccount(query); err != nil {
			return internalServerError(err.Error())
		}
	}

	var updateBody model.BankAccountDeleteBody
	updateBody.AccountNumber = fmt.Sprintf("%s_del%d", account.AccountNumber, account.Id)
	updateBody.DeletedAt = time.Now()
	if err := s.repo.DeleteBankAccount(id, updateBody); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) GetTransactionById(req model.BankGetByIdRequest) (*model.BankAccountTransaction, error) {

	record, err := s.repo.GetTransactionById(req.Id)
	if err != nil {
		if err.Error() == recordNotFound {
			return nil, notFound(transactionNotFound)
		}
		return nil, internalServerError(err.Error())
	}
	return record, nil
}

func (s *accountingService) GetTransactions(req model.BankAccountTransactionListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.UnlimitPagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	list, err := s.repo.GetTransactions(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return list, nil
}

func (s *accountingService) CreateTransaction(body model.BankAccountTransactionBody) error {

	account, err := s.repo.GetBankAccountById(body.AccountId)
	if err != nil {
		fmt.Println(err)
		return badRequest("Invalid Bank Account")
	}

	var transaction model.BankAccountTransactionBody
	transaction.AccountId = account.Id
	transaction.Description = body.Description
	transaction.TransferType = body.TransferType
	transaction.Amount = body.Amount
	transaction.TransferAt = body.TransferAt
	transaction.CreatedByUsername = body.CreatedByUsername

	if err := s.repo.CreateTransaction(transaction); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) UpdateTransaction(id int64, body model.BankAccountTransactionBody) error {

	_, err := s.repo.GetTransactionById(id)
	if err != nil {
		return internalServerError(err.Error())
	}

	// no Update
	return notFound("Function not found")
}

func (s *accountingService) DeleteTransaction(id int64) error {

	_, err := s.repo.GetTransactionById(id)
	if err != nil {
		return internalServerError(err.Error())
	}

	if err := s.repo.DeleteTransaction(id); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) GetTransferById(req model.BankGetByIdRequest) (*model.BankAccountTransfer, error) {

	record, err := s.repo.GetTransferById(req.Id)
	if err != nil {
		if err.Error() == recordNotFound {
			return nil, notFound(transferNotFound)
		}
		return nil, internalServerError(err.Error())
	}
	return record, nil
}

func (s *accountingService) GetTransfers(req model.BankAccountTransferListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.UnlimitPagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	list, err := s.repo.GetTransfers(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return list, nil
}

func (s *accountingService) CreateTransfer(body model.BankAccountTransferBody) error {

	fromAccount, err := s.repo.GetBankAccountById(body.FromAccountId)
	if err != nil {
		fmt.Println(err)
		return badRequest("Invalid source Bank Account")
	}

	toAccount, err := s.repo.GetBankAccountById(body.ToAccountId)
	if err != nil {
		fmt.Println(err)
		return badRequest("Invalid destination Bank Account")
	}

	var createBody model.BankAccountTransferBody
	createBody.FromAccountId = fromAccount.Id
	createBody.FromBankId = fromAccount.BankId
	createBody.FromAccountName = fromAccount.AccountName
	createBody.FromAccountNumber = fromAccount.AccountNumber
	createBody.ToAccountId = toAccount.Id
	createBody.ToBankId = toAccount.BankId
	createBody.ToAccountName = toAccount.AccountName
	createBody.ToAccountNumber = toAccount.AccountNumber
	createBody.Amount = body.Amount
	createBody.TransferAt = body.TransferAt
	createBody.CreatedByUsername = body.CreatedByUsername
	createBody.Status = "pending"
	if err := s.repo.CreateTransfer(createBody); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) ConfirmTransfer(id int64, actorId int64) error {

	transfer, err := s.repo.GetTransferById(id)
	if err != nil {
		return internalServerError(err.Error())
	}

	if transfer.Status == "pending" {
		var body model.BankAccountTransferConfirmBody
		body.Status = "confirmed"
		body.ConfirmedAt = time.Now()
		body.ConfirmedByUserId = actorId
		if err := s.repo.ConfirmTransfer(id, body); err != nil {
			return internalServerError(err.Error())
		}
	} else {
		return badRequest("Transfer not in pending status")
	}
	return nil
}

func (s *accountingService) DeleteTransfer(id int64) error {

	_, err := s.repo.GetTransferById(id)
	if err != nil {
		return internalServerError(err.Error())
	}

	if err := s.repo.DeleteTransfer(id); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) GetAccountStatements(req model.BankAccountStatementListRequest) (*model.SuccessWithPagination, error) {

	systemAccount, err := s.repo.GetBankAccountById(req.AccountId)
	if err != nil {
		if err.Error() == recordNotFound {
			return nil, notFound(bankAccountNotFound)
		}
		return nil, internalServerError(err.Error())
	}

	var query model.ExternalStatementListRequest
	query.AccountNumber = systemAccount.AccountNumber
	query.Page = req.Page
	query.Limit = req.Limit
	records, err := s.GetExternalAccountStatements(query)
	if err != nil {
		return &model.SuccessWithPagination{}, nil
		// onError return Empty | return nil, internalServerError(err.Error())
	}

	// todo : Sync ?
	var externalIds []int64
	for _, record := range records.List.([]model.ExternalStatement) {
		externalIds = append(externalIds, record.Id)
	}
	fmt.Println("externalIds", externalIds)

	dbStatements, err := s.repo.GetBankExternalStatements(externalIds)
	if err != nil {
		return records, nil
		// onError return Empty | return nil, internalServerError(err.Error())
	}

	for i, record := range records.List.([]model.ExternalStatement) {
		var exist *model.BankStatementResponse = nil
		for _, dbRecord := range dbStatements.List.([]model.BankStatementResponse) {
			if record.Id == dbRecord.ExternalId {
				exist = &dbRecord
				break
			}
		}
		tempRecord := records.List.([]model.ExternalStatement)[i]
		tempRecord.AccountDetail = fmt.Sprintf("%s: %s", systemAccount.BankName, systemAccount.AccountName)
		// "createdDate": "2023-05-08 11:17:05",
		// "dateTime": "2023-05-08 11:15:37",
		// tempRecord.DateTime = strings.Replace(tempRecord.DateTime, " ", "T", 1) + "Z"
		// "updatedDate": "2023-05-08 11:17:05"
		if record.TxnCode == "X1" || record.TxnCode == "CR" {
			tempRecord.StatementType = "transfer_in"
			tempRecord.Amount = record.Amount
		} else if record.TxnCode == "X2" || record.TxnCode == "DR" {
			tempRecord.StatementType = "transfer_out"
			tempRecord.Amount = record.Amount * -1
		}
		if exist != nil {
			tempRecord.Status = exist.Status
		} else {
			tempRecord.Status = "new"
		}
		records.List.([]model.ExternalStatement)[i] = tempRecord
	}
	return records, nil
}

func (s *accountingService) GetAccountStatementById(req model.BankGetByIdRequest) (*model.BankStatement, error) {

	record, err := s.repo.GetBankStatementById(req.Id)
	if err != nil {
		if err.Error() == recordNotFound {
			return nil, notFound(err.Error())
		}
		return nil, internalServerError(err.Error())
	}
	return record, nil
}

func (s *accountingService) AddAccountStatementToWebhook(req model.RecheckWebhookRequest) error {

	systemAccount, err := s.repo.GetBankAccountById(req.AccountId)
	if err != nil {
		if err.Error() == recordNotFound {
			return notFound(bankAccountNotFound)
		}
		return internalServerError(err.Error())
	}

	if err := s.repo.HasBankExternalStatements(req.ExternalId); err != nil {
		if err.Error() == recordNotFound {
			var query model.ExternalStatementListRequest
			query.AccountNumber = systemAccount.AccountNumber
			query.OfDateTime = req.OfDateTime
			query.Page = 1
			query.Limit = 100
			records, err := s.GetExternalAccountStatements(query)
			if err != nil {
				return internalServerError(err.Error())
			}
			for _, record := range records.List.([]model.ExternalStatement) {
				if record.Id == req.ExternalId {
					err := s.CreateBankStatementFromExternalStatement(record)
					if err == nil {
						continue
					}
				}
			}
		}
	}
	return nil
}

func (s *accountingService) GetExternalSettings() (*model.ExternalSettings, error) {

	var body model.ExternalSettings
	body.ApiEndpoint = os.Getenv("ACCOUNTING_API_ENDPOINT")
	body.ApiKey = os.Getenv("ACCOUNTING_API_KEY")
	body.LocalWebhookEndpoint = os.Getenv("ACCOUNTING_LOCAL_WEBHOOK_ENDPOINT")

	return &body, nil
}

func (s *accountingService) HasExternalAccount(accountNumber string) bool {

	data, err := s.GetExternalAccounts()
	if err != nil {
		return true
	}

	for _, account := range data.List.([]model.ExternalAccount) {
		if account.AccountNo == accountNumber {
			return true
		}
	}
	return false
}

func (s *accountingService) GetExternalAccountConfig(key string) (*model.BotAccountConfig, error) {

	var query model.BotAccountConfigListRequest
	query.SearchKey = &key

	data, err := s.repo.GetBotaccountConfigs(query)
	if err != nil {
		return nil, internalServerError(err.Error())
	}

	for _, record := range data.List.([]model.BotAccountConfig) {
		if record.ConfigKey == key {
			return &record, nil
		}
	}
	return nil, notFound("Config not found")
}

func (s *accountingService) IsAllowCreateExternalAccount(accountNumber string) bool {

	allowCreateExternalAccount := false
	config, _ := s.GetExternalAccountConfig("allow_create_external_account")
	if config != nil {
		if config.ConfigVal == "list" {
			accountConfig, errConfig := s.HasExternalAccountConfig("allow_external_account_number", accountNumber)
			if errConfig != nil {
				return false
			}
			if accountConfig.ConfigVal == accountNumber {
				allowCreateExternalAccount = true
			}
		} else if config.ConfigVal == "all" {
			allowCreateExternalAccount = true
		}
	}
	return allowCreateExternalAccount
}

func (s *accountingService) HasExternalAccountConfig(key string, value string) (*model.BotAccountConfig, error) {

	var query model.BotAccountConfigListRequest
	query.SearchKey = &key
	query.SearchValue = &value

	data, err := s.repo.GetBotaccountConfigs(query)
	if err != nil {
		return nil, internalServerError(err.Error())
	}

	for _, record := range data.List.([]model.BotAccountConfig) {
		if record.ConfigKey == key {
			return &record, nil
		}
	}
	return nil, notFound("Config not found")
}

func (s *accountingService) GetCustomerAccountsInfo(req model.CustomerAccountInfoRequest) (*model.CustomerAccountInfo, error) {

	botAccount, err := s.repo.GetActiveExternalAccount()
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	req.AccountFrom = botAccount.AccountNumber
	// b, err := json.Marshal(req)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil, internalServerError("Error from JSON")
	// }
	// fmt.Println(string(b))

	client := &http.Client{}
	// curl -X POST "https://api.fastbankapi.com/api/v2/statement/verifyTransfer" -H "accept: */*" -H "apiKey: aa.bb" -H "Content-Type: application/json" -d "{ \"accountFrom\": \"cccc\", \"accountTo\": \"dddd\", \"bankCode\": \"bay\"}"
	data, _ := json.Marshal(req)
	reqExternal, _ := http.NewRequest("POST", os.Getenv("ACCOUNTING_API_ENDPOINT")+"/api/v2/statement/verifyTransfer", bytes.NewBuffer(data))
	reqExternal.Header.Set("apiKey", os.Getenv("ACCOUNTING_API_KEY"))
	reqExternal.Header.Set("Content-Type", "application/json")
	response, err := client.Do(reqExternal)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	if response.StatusCode != 200 {
		fmt.Println(response)
		return nil, internalServerError("Error from external API")
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var result model.CustomerAccountInfoReponse
	errJson := json.Unmarshal(responseData, &result)
	if errJson != nil {
		return nil, internalServerError("Error from JSON response")
	}
	return &result.Data, nil
}

func (s *accountingService) GetExternalAccounts() (*model.SuccessWithPagination, error) {

	client := &http.Client{}
	reqExternal, _ := http.NewRequest("GET", os.Getenv("ACCOUNTING_API_ENDPOINT")+"/api/v2/site/bankAccount", nil)
	reqExternal.Header.Set("apiKey", os.Getenv("ACCOUNTING_API_KEY"))
	response, err := client.Do(reqExternal)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	if response.StatusCode != 200 {
		fmt.Println(response)
		return nil, internalServerError("Error from external API")
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var list []model.ExternalAccount
	errJson := json.Unmarshal(responseData, &list)
	if errJson != nil {
		return nil, internalServerError("Error from JSON response")
	}

	// End count total records for pagination purposes (without limit and offset) //
	var result model.SuccessWithPagination
	result.List = list
	result.Total = int64(len(list))
	return &result, nil
}

func (s *accountingService) GetExternalAccountBalance(query model.ExternalAccountStatusRequest) (*model.ExternalAccountBalance, error) {

	client := &http.Client{}
	// curl -X GET "https://api.fastbankapi.com/api/v2/statement/balance?accountNo=hhhhhh" -H "accept: */*" -H "apiKey: aaaaaa.bbbbbb"
	reqExternal, _ := http.NewRequest("GET", os.Getenv("ACCOUNTING_API_ENDPOINT")+"/api/v2/statement/balance?accountNo="+query.AccountNumber, nil)
	reqExternal.Header.Set("apiKey", os.Getenv("ACCOUNTING_API_KEY"))
	response, err := client.Do(reqExternal)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	if response.StatusCode != 200 {
		fmt.Println(response)
		return nil, internalServerError("Error from external API")
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var result model.ExternalAccountBalance
	errJson := json.Unmarshal(responseData, &result)
	if errJson != nil {
		return nil, internalServerError("Error from JSON response")
	}
	if result.AccountNo != query.AccountNumber {
		_, err := s.CreateWebhookLog("GetExternalAccountBalance, ERROR:", string(responseData))
		log.Panicln(err)
		return nil, notFound("Bank account not found")
	}
	return &result, nil
}

func (s *accountingService) GetExternalAccountStatus(query model.ExternalAccountStatusRequest) (*model.ExternalAccountStatus, error) {

	client := &http.Client{}
	reqExternal, _ := http.NewRequest("GET", os.Getenv("ACCOUNTING_API_ENDPOINT")+"/api/v2/site/bank-status?accountNo="+query.AccountNumber, nil)
	reqExternal.Header.Set("apiKey", os.Getenv("ACCOUNTING_API_KEY"))
	response, err := client.Do(reqExternal)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != 200 {
		_, err := s.CreateWebhookLog("GetExternalAccountStatus, ERROR", helper.StructJson(struct {
			query        model.ExternalAccountStatusRequest
			responseJson string
		}{query, string(responseData)}))
		log.Panicln(err)
		return nil, notFound("External account not found")
	}

	var result model.ExternalAccountStatus
	errJson := json.Unmarshal(responseData, &result)
	if errJson != nil {
		return nil, internalServerError("Error from JSON response")
	}
	return &result, nil
}

func (s *accountingService) CreateExternalAccount(body model.ExternalAccountCreateBody) (*model.ExternalAccountCreateResponse, error) {

	client := &http.Client{}
	data, _ := json.Marshal(body)
	reqExternal, _ := http.NewRequest("POST", os.Getenv("ACCOUNTING_API_ENDPOINT")+"/api/v2/site/bankAccount", bytes.NewBuffer(data))
	reqExternal.Header.Set("apiKey", os.Getenv("ACCOUNTING_API_KEY"))
	reqExternal.Header.Set("Content-Type", "application/json")
	response, err := client.Do(reqExternal)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != 200 {
		_, err := s.CreateWebhookLog("CreateExternalAccount, ERROR:", string(responseData))
		log.Panicln(err)
		return nil, internalServerError("Error from external API")
	}

	var result model.ExternalAccountCreateResponse
	errJson := json.Unmarshal(responseData, &result)
	if errJson != nil {
		return nil, internalServerError("Error from JSON response")
	}
	jsonResult, err := json.Marshal(result)
	if err == nil {
		_, webhookErr := s.CreateWebhookLog("CreateExternalAccount, SUCCESS", string(jsonResult))
		log.Panicln(webhookErr)
	}
	return &result, nil
}

func (s *accountingService) UpdateExternalAccount(body model.ExternalAccountUpdateBody) (*model.ExternalAccountCreateResponse, error) {

	client := &http.Client{}
	data, _ := json.Marshal(body)
	reqExternal, _ := http.NewRequest("PUT", os.Getenv("ACCOUNTING_API_ENDPOINT")+"/api/v2/site/bankAccount", bytes.NewBuffer(data))
	reqExternal.Header.Set("apiKey", os.Getenv("ACCOUNTING_API_KEY"))
	reqExternal.Header.Set("Content-Type", "application/json")
	response, err := client.Do(reqExternal)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != 200 {
		_, err := s.CreateWebhookLog("UpdateExternalAccount, ERROR:", string(responseData))
		log.Panicln(err)
		return nil, internalServerError("Error from external API")
	}
	var result model.ExternalAccountCreateResponse
	errJson := json.Unmarshal(responseData, &result)
	if errJson != nil {
		return nil, internalServerError("Error from JSON response")
	}
	jsonResult, err := json.Marshal(result)
	if err == nil {
		_, webhookErr := s.CreateWebhookLog("UpdateExternalAccount, SUCCESS", string(jsonResult))
		log.Panicln(webhookErr)
	}
	return &result, nil
}

func (s *accountingService) DeleteExternalAccount(query model.ExternalAccountStatusRequest) error {

	client := &http.Client{}
	reqExternal, _ := http.NewRequest("DELETE", os.Getenv("ACCOUNTING_API_ENDPOINT")+"/api/v2/site/bankAccount/"+query.AccountNumber, nil)
	reqExternal.Header.Set("apiKey", os.Getenv("ACCOUNTING_API_KEY"))
	response, err := client.Do(reqExternal)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	if response.StatusCode != 200 {
		fmt.Println(response)
		return internalServerError("Error from external API")
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	_, webhookErr := s.CreateWebhookLog("DeleteExternalAccount, responseData:", string(responseData))
	log.Panicln(webhookErr)
	return nil
}

func (s *accountingService) EnableExternalAccount(req model.ExternalAccountEnableRequest) (*model.ExternalAccountStatus, error) {

	client := &http.Client{}
	// curl -X POST "https://api.fastbankapi.com/api/v2/site/enable-bank" -H "accept: */*" -H "apiKey: 123" -H "Content-Type: application/json" -d "{ \"accountNo\": \"string\", \"enable\": true}"
	data, _ := json.Marshal(req)
	reqExternal, _ := http.NewRequest("POST", os.Getenv("ACCOUNTING_API_ENDPOINT")+"/api/v2/site/enable-bank", bytes.NewBuffer(data))
	reqExternal.Header.Set("apiKey", os.Getenv("ACCOUNTING_API_KEY"))
	reqExternal.Header.Set("Content-Type", "application/json")
	response, err := client.Do(reqExternal)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	if response.StatusCode != 200 {
		fmt.Println(response)
		return nil, internalServerError("Error from external API")
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("EnableExternalAccount:", string(responseData))
	// {"success":true,"enable":true,"status":"online"}
	// {"success":true,"enable":false,"status":"offline"}
	var result model.ExternalAccountStatus
	errJson := json.Unmarshal(responseData, &result)
	if errJson != nil {
		return nil, internalServerError("Error from JSON response")
	}
	return &result, nil
}

func (s *accountingService) GetExternalAccountLogs(req model.ExternalStatementListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.UnlimitPagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}

	client := &http.Client{}
	// curl -X GET "https://api.fastbankapi.com/api/v2/site/bankAccount/logs?accountNo=aaaaaaaaaaaaaa&page=0&size=10" -H "accept: */*" -H "apiKey: xxxxxxxxxx.yyyyyyyyyyy"
	queryString := fmt.Sprintf("&page=%d&size=%d", req.Page, req.Limit)
	fullPath := os.Getenv("ACCOUNTING_API_ENDPOINT") + "/api/v2/site/bankAccount/logs?accountNo=" + req.AccountNumber + queryString
	reqExternal, _ := http.NewRequest("GET", fullPath, nil)
	reqExternal.Header.Set("apiKey", os.Getenv("ACCOUNTING_API_KEY"))
	response, err := client.Do(reqExternal)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	if response.StatusCode != 200 {
		fmt.Println(response)
		return nil, internalServerError("Error from external API")
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var externalList model.ExternalListWithPagination
	errJson := json.Unmarshal(responseData, &externalList)
	if errJson != nil {
		return nil, internalServerError("Error from JSON response")
	}
	// fmt.Println("response", string(responseData))

	// End count total records for pagination purposes (without limit and offset) //
	var result model.SuccessWithPagination
	result.List = externalList.Content
	result.Total = externalList.TotalElements
	return &result, nil
}

func (s *accountingService) GetExternalAccountStatements(req model.ExternalStatementListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.UnlimitPagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}

	client := &http.Client{}
	// https://api.fastbankapi.com/api/v2/statement?accountNo=aaaaa&page=0&size=10&txnCode=all
	// curl -X GET "https://api.fastbankapi.com/api/v2/statement?accountNo=hhhhhh&page=0&size=10&txnCode=all" -H "accept: */*" -H "apiKey: aaaaaaaa.bbbbbbbbbbb"
	queryString := fmt.Sprintf("&page=%d&size=%d&txnCode=all", req.Page, req.Limit)
	if req.OfDateTime != "" {
		queryString += fmt.Sprintf("&date=%s", req.OfDateTime)
	}
	fullPath := os.Getenv("ACCOUNTING_API_ENDPOINT") + "/api/v2/statement?accountNo=" + req.AccountNumber + queryString
	reqExternal, _ := http.NewRequest("GET", fullPath, nil)
	reqExternal.Header.Set("apiKey", os.Getenv("ACCOUNTING_API_KEY"))
	response, err := client.Do(reqExternal)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	if response.StatusCode != 200 {
		fmt.Println(fullPath, response)
		return nil, internalServerError("Error from external API")
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var externalList model.ExternalStatementListWithPagination
	errJson := json.Unmarshal(responseData, &externalList)
	if errJson != nil {
		fmt.Println(errJson)
		return nil, internalServerError("Error from JSON response")
	}
	// fmt.Println("response", string(responseData))

	// End count total records for pagination purposes (without limit and offset) //
	var result model.SuccessWithPagination
	result.List = externalList.Content
	result.Total = externalList.TotalElements
	return &result, nil
}

func (s *accountingService) GetExternalAccountStatementByTimestamp(req model.ExternalStatementListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.UnlimitPagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}

	client := &http.Client{}
	// https://api.fastbankapi.com/api/v2/statement?accountNo=aaaaa&date=2023-05-08%2011%3A15%3A37&page=0&size=10&txnCode=all
	// curl -X GET "https://api.fastbankapi.com/api/v2/statement?accountNo=aaaaa&date=2023-05-08%2011%3A15%3A37&page=0&size=10&txnCode=all" -H "accept: */*" -H "apiKey: bbbbb.cccccc"
	queryString := fmt.Sprintf("&page=%d&size=%d&txnCode=all", req.Page, req.Limit)
	if req.OfDateTime != "" {
		queryString += fmt.Sprintf("&date=%s", req.OfDateTime)
	}
	fullPath := os.Getenv("ACCOUNTING_API_ENDPOINT") + "/api/v2/statement?accountNo=" + req.AccountNumber + queryString
	reqExternal, _ := http.NewRequest("GET", fullPath, nil)
	reqExternal.Header.Set("apiKey", os.Getenv("ACCOUNTING_API_KEY"))
	response, err := client.Do(reqExternal)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	if response.StatusCode != 200 {
		fmt.Println(fullPath, response)
		return nil, internalServerError("Error from external API")
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var externalList model.ExternalStatementListWithPagination
	errJson := json.Unmarshal(responseData, &externalList)
	if errJson != nil {
		fmt.Println(errJson)
		return nil, internalServerError("Error from JSON response")
	}
	// fmt.Println("response", string(responseData))

	// End count total records for pagination purposes (without limit and offset) //
	var result model.SuccessWithPagination
	result.List = externalList.Content
	result.Total = externalList.TotalElements
	return &result, nil
}

func (s *accountingService) TransferExternalAccount(req model.ExternalAccountTransferRequest) error {

	var body model.ExternalAccountTransferBody
	systemAccount, err := s.repo.GetBankAccountById(req.SystemAccountId)
	if err != nil {
		if err.Error() == recordNotFound {
			return notFound(bankAccountNotFound)
		}
		return internalServerError(err.Error())
	}
	body.AccountForm = systemAccount.AccountNumber
	body.AccountTo = req.AccountNumber
	body.Amount = req.Amount
	body.BankCode = req.BankCode
	body.Pin = systemAccount.PinCode

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
			return internalServerError("Error from JSON response")
		}
		fmt.Println("errorModel", errorModel)
		if errorModel.Error != "" {
			return internalServerError(errorModel.Error)
		}
		return internalServerError("Error from external API")
	}
	return nil
}

func (s *accountingService) CreateBankStatementFromWebhook(data model.WebhookStatement) error {

	systemAccount, err := s.repo.GetBankAccountByExternalId(data.BankAccountId)
	if err != nil {
		fmt.Println(err)
		return badRequest("Invalid Bank Account")
	}

	// Check Exist statement
	_, errOldStatement := s.repo.GetWebhookStatementByExternalId(data.Id)
	if errOldStatement != nil && errOldStatement.Error() == recordNotFound {
		var bodyCreateState model.BankStatementCreateBody
		bodyCreateState.AccountId = systemAccount.Id
		bodyCreateState.ExternalId = data.Id
		if data.TxnCode == "X1" || data.TxnCode == "CR" {
			bodyCreateState.StatementType = "transfer_in"
			bodyCreateState.Amount = data.Amount
		} else if data.TxnCode == "X2" || data.TxnCode == "DR" {
			bodyCreateState.StatementType = "transfer_out"
			bodyCreateState.Amount = data.Amount * -1
		} else {
			_, err := s.CreateWebhookLog("unsupport TxnCode found, WebhookStatement:", helper.StructJson(struct{ data model.WebhookStatement }{data}))
			log.Panicln(err)
			return badRequest("Invalid TxnCode")
		}

		bank, err := s.GetBankFromWebhook(data.Info)
		if err != nil {
			return err
		}
		bodyCreateState.FromBankId = bank.Id
		accountNumber, _ := s.GetAccountNoFromWebhook(bank.Code, data.Info)
		bodyCreateState.FromAccountNumber = accountNumber

		bodyCreateState.Detail = data.TxnDescription + " " + data.Info
		bodyCreateState.TransferAt = data.DateTime
		bodyCreateState.Status = "pending"

		insertId, err := s.repo.CreateWebhookStatement(bodyCreateState)
		if err != nil {
			return internalServerError(err.Error())
		}

		// Auto Match if == 1
		var reqPosibleList model.MemberPossibleListRequest
		statement, err := s.repo.GetBankStatementById(*insertId)
		if err != nil {
			return err
		}
		reqPosibleList.UnknownStatementId = statement.Id
		reqPosibleList.UserBankCode = &statement.FromBankCode
		reqPosibleList.UserAccountNumber = &statement.FromAccountNumber

		records, err := s.repo.GetPossibleStatementOwners(reqPosibleList)
		if err != nil {
			return err
		}
		if records.Total == 1 {
			for _, possibleOwner := range records.List.([]model.Member) {
				// Auto create transaction
				if bodyCreateState.StatementType == "transfer_in" {
					// AutoDeposit
					if systemAccount.AutoCreditFlag == "auto" {
						if err := s.CreateAutoDepositTransaction(possibleOwner, bodyCreateState); err != nil {
							return internalServerError(err.Error())
						}
						var statementMatchRequest model.BankStatementMatchRequest
						statementMatchRequest.UserId = possibleOwner.Id
						if err := s.SetStatementOwnerMatched(statement.Id, statementMatchRequest); err != nil {
							return internalServerError(err.Error())
						}
					} else {
						if err := s.CreateDepositTransaction(possibleOwner, bodyCreateState); err != nil {
							return internalServerError(err.Error())
						}
					}
				} else if bodyCreateState.StatementType == "transfer_out" {
					// Auto ignore, no need to match
					var statementMatchRequest model.BankStatementMatchRequest
					statementMatchRequest.ConfirmedAt = time.Now()
					statementMatchRequest.ConfirmedByUserId = 0
					statementMatchRequest.ConfirmedByUsername = "อัตโนมัติ"
					if err := s.IgnoreStatementOwner(statement.Id, statementMatchRequest); err != nil {
						return internalServerError(err.Error())
					}
				}
			}
		} else if records.Total > 1 {
			return errors.New("too many possible owners")
		} else {
			return errors.New("no possible owners")
		}
	} else {
		return errors.New("statement already exists")
	}
	return nil
}

func (s *accountingService) CreateBankStatementFromExternalStatement(data model.ExternalStatement) error {

	systemAccount, err := s.repo.GetBankAccountByExternalId(data.BankAccountId)
	if err != nil {
		fmt.Println(err)
		return badRequest("Invalid Bank Account")
	}

	_, errOldStatement := s.repo.GetWebhookStatementByExternalId(data.Id)
	if errOldStatement != nil && errOldStatement.Error() == recordNotFound {
		var bodyCreateState model.BankStatementCreateBody
		bodyCreateState.AccountId = systemAccount.Id
		bodyCreateState.ExternalId = data.Id
		if data.TxnCode == "X1" || data.TxnCode == "CR" {
			bodyCreateState.StatementType = "transfer_in"
			bodyCreateState.Amount = data.Amount
		} else if data.TxnCode == "X2" || data.TxnCode == "DR" {
			bodyCreateState.StatementType = "transfer_out"
			bodyCreateState.Amount = data.Amount * -1
		} else {
			// s.CreateWebhookLog("unsupport TxnCode found, WebhookStatement:", helper.StructJson(struct{ data model.WebhookStatement }{data}))
			return badRequest("Invalid TxnCode")
		}

		bank, _ := s.GetBankFromWebhook(data.Info)
		bodyCreateState.FromBankId = bank.Id
		accountNumber, _ := s.GetAccountNoFromWebhook(bank.Code, data.Info)
		bodyCreateState.FromAccountNumber = accountNumber

		bodyCreateState.Detail = data.TxnDescription + " " + data.Info
		// tempDateTime := strings.Replace(data.DateTime, " ", "T", 1) + "Z"
		// unixTimeUTC := time.Unix(1405544146, 0)               //gives unix time stamp in utc
		// unitTimeInRFC3339 := unixTimeUTC.Format(time.RFC3339) // converts utc time to RFC3339 format
		// timeParseLayout := "DateTime"
		tempDateTime, _ := time.Parse("2006-01-02 15:04:05", data.DateTime)
		fmt.Println(data.DateTime, tempDateTime)
		bodyCreateState.TransferAt = tempDateTime
		bodyCreateState.Status = "pending"

		insertId, err := s.repo.CreateWebhookStatement(bodyCreateState)
		if err != nil {
			return internalServerError(err.Error())
		}

		// Auto Match if == 1
		var reqPosibleList model.MemberPossibleListRequest
		statement, err := s.repo.GetBankStatementById(*insertId)
		if err != nil {
			return nil
		}
		reqPosibleList.UnknownStatementId = statement.Id
		reqPosibleList.UserBankCode = &statement.FromBankCode
		reqPosibleList.UserAccountNumber = &statement.FromAccountNumber

		records, err := s.repo.GetPossibleStatementOwners(reqPosibleList)
		if err != nil {
			return nil
		}
		if records.Total == 1 {
			for _, possibleOwner := range records.List.([]model.Member) {
				// Auto create transaction
				if bodyCreateState.StatementType == "transfer_in" {
					// AutoDeposit
					if systemAccount.AutoCreditFlag == "auto" {
						if err := s.CreateAutoDepositTransaction(possibleOwner, bodyCreateState); err != nil {
							// return internalServerError(err.Error())
							return nil
						}
						var statementMatchRequest model.BankStatementMatchRequest
						statementMatchRequest.UserId = possibleOwner.Id
						if err := s.SetStatementOwnerMatched(statement.Id, statementMatchRequest); err != nil {
							// return internalServerError(err.Error())
							return nil
						}
					} else {
						if err := s.CreateDepositTransaction(possibleOwner, bodyCreateState); err != nil {
							// return internalServerError(err.Error())
							return nil
						}
					}
				} else if bodyCreateState.StatementType == "transfer_out" {
					// Auto ignore, no need to match
					var statementMatchRequest model.BankStatementMatchRequest
					statementMatchRequest.ConfirmedAt = time.Now()
					statementMatchRequest.ConfirmedByUserId = 0
					statementMatchRequest.ConfirmedByUsername = "อัตโนมัติ"
					if err := s.IgnoreStatementOwner(statement.Id, statementMatchRequest); err != nil {
						// return internalServerError(err.Error())
						return nil
					}
				}
			}
		}
	}

	return nil
}

func (s *accountingService) CreateAutoDepositTransaction(possibleOwner model.Member, bodyCreateState model.BankStatementCreateBody) error {

	if bodyCreateState.StatementType == "transfer_in" {
		var createDepositBody model.BankTransactionCreateBody
		createDepositBody.MemberCode = possibleOwner.MemberCode
		createDepositBody.TransferType = "deposit"
		createDepositBody.CreditAmount = bodyCreateState.Amount
		createDepositBody.TransferAt = &bodyCreateState.TransferAt
		createDepositBody.CreatedByUserId = 0
		createDepositBody.CreatedByUsername = "อัตโนมัติ"
		createDepositBody.IsAutoCredit = true
		// later: promotionId bonusAmount
		// bodyCreateState.AccountId = systemAccount.Id == createDepositBody.ToAccountId = &systemAccount.Id
		createDepositBody.ToAccountId = &bodyCreateState.AccountId
		transId, err := s.CreateBankTransaction(createDepositBody)
		if err != nil {
			fmt.Print(err.Error())
			return internalServerError(err.Error())
		}
		var confirmTransReq model.BankConfirmDepositRequest
		confirmTransReq.ConfirmedAt = time.Now()
		confirmTransReq.ConfirmedByUserId = 0
		confirmTransReq.ConfirmedByUsername = "อัตโนมัติ"
		confirmTransReq.TransferAt = &bodyCreateState.TransferAt
		// confirmTransReq.BonusAmount = 0
		if err := s.ConfirmDepositTransaction(*transId, confirmTransReq); err != nil {
			fmt.Print(err.Error())
			return internalServerError(err.Error())
		}
		// no need if IsAutoCredit = true err := s.ConfirmDepositCredit(*transId, confirmTransReq); err != nil {
		// 	fmt.Print(err.Error())
		// 	return internalServerError(err.Error())
		// }
		// var statementMatchRequest model.BankStatementMatchRequest
		// statementMatchRequest.UserId = possibleOwner.Id
		// if err := s.MatchStatementOwner(statement.Id, statementMatchRequest); err != nil {
		// 	// return internalServerError(err.Error())
		// 	return nil
		// }
	}

	return nil
}

func (s *accountingService) CreateDepositTransaction(possibleOwner model.Member, bodyCreateState model.BankStatementCreateBody) error {

	if bodyCreateState.StatementType == "transfer_in" {
		var createDepositBody model.BankTransactionCreateBody
		createDepositBody.MemberCode = possibleOwner.MemberCode
		createDepositBody.TransferType = "deposit"
		createDepositBody.CreditAmount = bodyCreateState.Amount
		createDepositBody.TransferAt = &bodyCreateState.TransferAt
		createDepositBody.CreatedByUserId = 0
		createDepositBody.CreatedByUsername = "อัตโนมัติ"
		createDepositBody.IsAutoCredit = false
		// later: promotionId bonusAmount
		// bodyCreateState.AccountId = systemAccount.Id == createDepositBody.ToAccountId = &systemAccount.Id
		createDepositBody.ToAccountId = &bodyCreateState.AccountId
		_, err := s.CreateBankTransaction(createDepositBody)
		if err != nil {
			fmt.Print(err.Error())
			return internalServerError(err.Error())
		}
		// var confirmTransReq model.BankConfirmDepositRequest
		// confirmTransReq.ConfirmedAt = time.Now()
		// confirmTransReq.ConfirmedByUserId = 0
		// confirmTransReq.ConfirmedByUsername = "อัตโนมัติ"
		// confirmTransReq.TransferAt = &bodyCreateState.TransferAt
		// // confirmTransReq.BonusAmount = 0
		// if err := s.ConfirmDepositTransaction(*transId, confirmTransReq); err != nil {
		// 	fmt.Print(err.Error())
		// 	return internalServerError(err.Error())
		// }
		// if err := s.ConfirmDepositCredit(*transId, confirmTransReq); err != nil {
		// 	fmt.Print(err.Error())
		// 	return internalServerError(err.Error())
		// }
		// var statementMatchRequest model.BankStatementMatchRequest
		// statementMatchRequest.UserId = possibleOwner.Id
		// if err := s.MatchStatementOwner(statement.Id, statementMatchRequest); err != nil {
		// 	// return internalServerError(err.Error())
		// 	return nil
		// }
	}

	return nil
}

func (s *accountingService) CreateBankTransaction(data model.BankTransactionCreateBody) (*int64, error) {

	var transId *int64
	var body model.BankTransactionCreateBody
	body.TransferAt = data.TransferAt
	body.CreatedByUserId = data.CreatedByUserId
	body.CreatedByUsername = data.CreatedByUsername
	body.Status = "pending"

	if data.TransferType == "deposit" {
		member, err := s.repo.GetUserByMemberCode(data.MemberCode)
		if err != nil {
			fmt.Println(err)
			return nil, badRequest("Invalid Member code")
		}
		bank, err := s.repo.GetBankByCode(member.BankCode)
		if err != nil {
			fmt.Println(err)
			return nil, badRequest("Invalid User Bank")
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
			return nil, badRequest("Input Bank Account")
		}
		toAccount, err := s.repo.GetDepositAccountById(*data.ToAccountId)
		if err != nil {
			fmt.Println(err)
			return nil, badRequest("Invalid Bank Account")
		}
		body.ToAccountId = &toAccount.Id
		body.ToBankId = &toAccount.BankId
		body.ToAccountName = &toAccount.AccountName
		body.ToAccountNumber = &toAccount.AccountNumber
		// later: createBonus + refDeposit
		body.PromotionId = data.PromotionId

		if insertId, err := s.repo.CreateBankDepositTransaction(body); err == nil {
			transId = insertId
		} else {
			return nil, internalServerError(err.Error())
		}
	} else if data.TransferType == "withdraw" {
		member, err := s.repo.GetUserByMemberCode(data.MemberCode)
		if err != nil {
			fmt.Println(err)
			return nil, badRequest("Invalid Member code")
		}
		bank, err := s.repo.GetBankByCode(member.BankCode)
		if err != nil {
			fmt.Println(err)
			return nil, badRequest("Invalid User Bank")
		}
		body.MemberCode = *member.MemberCode
		body.UserId = member.Id
		body.CreditAmount = data.CreditAmount
		body.TransferType = data.TransferType

		// Withdraw SystemAccount is not requried
		if data.FromAccountId != nil {
			fromAccount, err := s.repo.GetWithdrawAccountById(*data.FromAccountId)
			if err != nil {
				fmt.Println(err)
				return nil, badRequest("Invalid Bank Account")
			}
			body.FromAccountId = &fromAccount.Id
			body.FromBankId = &fromAccount.BankId
			body.FromAccountName = &fromAccount.AccountName
			body.FromAccountNumber = &fromAccount.AccountNumber
		}

		body.ToBankId = &bank.Id
		body.ToAccountName = &member.Fullname
		body.ToAccountNumber = &member.BankAccount
		body.Status = "pending" // == pending_credit

		if insertId, err := s.repo.CreateBankWithdrawTransaction(body); err == nil {
			transId = insertId
		} else {
			return nil, internalServerError(err.Error())
		}
	} else {
		return nil, badRequest("Invalid Transfer Type")
	}

	return transId, nil
}

func (s *accountingService) SetStatementOwnerMatched(id int64, req model.BankStatementMatchRequest) error {

	statement, err := s.repo.GetBankStatementById(id)
	if err != nil {
		return internalServerError(err.Error())
	}
	if statement.Status != "pending" {
		return badRequest("Statement is not pending")
	}
	member, err := s.repo.GetMemberById(req.UserId)
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
	if err := s.repo.CreateStatementAction(createBody); err == nil {
		var body model.BankStatementUpdateBody
		body.Status = "confirmed"
		if err := s.repo.UpdateBankStatement(id, body); err != nil {
			return internalServerError(err.Error())
		}
	} else {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) IgnoreStatementOwner(id int64, req model.BankStatementMatchRequest) error {

	statement, err := s.repo.GetBankStatementById(id)
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
	if err := s.repo.CreateStatementAction(createBody); err == nil {
		if err := s.repo.IgnoreStatementOwner(id, body); err != nil {
			return internalServerError(err.Error())
		}
	} else {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) ConfirmDepositTransaction(id int64, req model.BankConfirmDepositRequest) error {

	record, err := s.repo.GetBankTransactionById(id)
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
		createBody.BonusAmount = *req.BonusAmount
	}
	createBody.ConfirmedAt = req.ConfirmedAt
	createBody.ConfirmedByUserId = req.ConfirmedByUserId
	createBody.ConfirmedByUsername = req.ConfirmedByUsername
	if actionId, err := s.repo.CreateTransactionAction(createBody); err == nil {
		// do nothing ?
		if err := s.repo.ConfirmPendingDepositTransaction(id, updateData); err != nil {
			if err := s.repo.RollbackTransactionAction(*actionId); err == nil {
				return internalServerError(err.Error())
			}
			return internalServerError(err.Error())
		}
	} else {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) ConfirmDepositCredit(id int64, req model.BankConfirmDepositRequest) error {

	record, err := s.repo.GetBankTransactionById(id)
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
	if _, err := s.repo.CreateTransactionAction(createBody); err == nil {
		fmt.Println("ConfirmPendingTransaction updateData:", helper.StructJson(updateData))
		if err := s.increaseMemberCreditFromDeposit(record.UserId, record.CreditAmount, "ฝากเงิน"); err != nil {
			return internalServerError(err.Error())
		}
		if record.BonusAmount > 0 {
			if err := s.increaseMemberCreditFromBonus(record.UserId, record.CreditAmount, "ได้รับโบนัสจากการฝากเงิน"); err != nil {
				return internalServerError(err.Error())
			}
		}
		if err := s.repo.ConfirmPendingCreditDepositTransaction(id, updateData); err != nil {
			return internalServerError(err.Error())
		}
	} else {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) increaseMemberCreditFromDeposit(userId int64, creditAmount float64, info string) error {

	statementType, err := s.repo.GetMemberStatementTypeByCode("deposit")
	if err != nil {
		return badRequest("Invalid Type")
	}

	var body model.MemberStatementCreateBody
	body.UserId = userId
	body.StatementTypeId = statementType.Id
	body.Info = info
	body.Amount = creditAmount
	if err := s.repo.IncreaseMemberCredit(body); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) increaseMemberCreditFromBonus(userId int64, creditAmount float64, info string) error {

	statementType, err := s.repo.GetMemberStatementTypeByCode("bonus")
	if err != nil {
		return badRequest("Invalid Type")
	}

	var body model.MemberStatementCreateBody
	body.UserId = userId
	body.StatementTypeId = statementType.Id
	body.Info = info
	body.Amount = creditAmount
	if err := s.repo.IncreaseMemberCredit(body); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) GetBankFromWebhookInfo(info string) (*model.BankResponse, error) {

	// sample : "กรุงศรีอยุธยา (BAY) /X123456",
	// fmt.Print("หาธนาคาร getBankIdFromStatementInformation.kt info:", info)
	infoStr := strings.ToLower(info)
	// fmt.Print("หาธนาคาร getBankIdFromStatementInformation.kt infoStr:", infoStr)

	var req model.BankListRequest
	banks, err := s.repo.GetBanks(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	for _, bank := range banks.List.([]model.BankResponse) {
		if strings.Contains(infoStr, bank.Code) {
			return &bank, nil
		} else if strings.Contains(infoStr, bank.Name) {
			return &bank, nil
		}
	}
	return nil, badRequest("Bank not found")
}

func (s *accountingService) GetBankFromWebhook(info string) (*model.BankResponse, error) {

	bank, err := s.GetBankFromWebhookInfo(info)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return bank, nil
}

func (s *accountingService) GetAccountNoFromWebhook(bankCode string, info string) (string, error) {

	// fmt.Print("GetAccountNoFromWebhook data:", helper.StructJson(data))
	infoStr := strings.ToLower(info)
	if bankCode == "scb" {
		// BankConstant.SCB -> statement.info.lowercase(Locale.getDefault()).split(" x")[1].take(4)
		infoStrings := strings.Split(infoStr, " x")
		if len(infoStrings) > 1 && len(infoStrings[1]) >= 4 {
			return infoStrings[1][:4], nil
		}
	} else {
		// else -> statement.info.lowercase(Locale.getDefault()).split("/x")[1].take(6)
		infoStrings := strings.Split(infoStr, "/x")
		if len(infoStrings) > 1 && len(infoStrings[1]) >= 6 {
			return infoStrings[1][:6], nil
		}
	}
	// later :
	// scb จะมีชื่อ ให้เช็คชื่อต่อ
	// KBANK xxx-x-x0000
	// เบื้องต้น support 2 ธนาคารก่อน
	return "", badRequest("AccountNo not found")
}

func (s *accountingService) CreateWebhookLog(logType string, jsonRequest string) (*int64, error) {

	var body model.WebhookLogCreateBody
	body.JsonRequest = jsonRequest
	body.JsonPayload = "{}"
	body.LogType = logType
	body.Status = "pending"

	insertId, err := s.repo.CreateWebhookLog(body)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return insertId, nil
}

func (s *accountingService) SetSuccessWebhookLog(id int64, jsonPayload string) error {

	var body model.WebhookLogUpdateBody
	body.JsonPayload = jsonPayload
	body.Status = "success"
	if err := s.repo.UpdateWebhookLog(id, body); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) SetFailedWebhookLog(id int64, jsonPayload string) error {

	var body model.WebhookLogUpdateBody
	body.JsonPayload = jsonPayload
	body.Status = "failed"
	if err := s.repo.UpdateWebhookLog(id, body); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) CreateBotaccountConfig(data model.BotAccountConfigCreateBody) error {

	if err := s.repo.CreateBotaccountConfig(data); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}
