package handler

import (
	"cybergame-api/middleware"
	"cybergame-api/model"
	"cybergame-api/repository"
	"cybergame-api/service"
	"encoding/json"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type accountingController struct {
	accountingService service.AccountingService
}

func newAccountingController(
	accountingService service.AccountingService,
) accountingController {
	return accountingController{accountingService}
}

func AccountingController(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewAccountingRepository(db)
	service := service.NewAccountingService(repo)
	handler := newAccountingController(service)

	root := r.Group("/accounting")
	root.GET("/autocreditflags/list", middleware.Authorize, handler.getAutoCreditFlags)
	root.GET("/autowithdrawflags/list", middleware.Authorize, handler.getAutoWithdrawFlags)
	root.GET("/qrwalletstatuses/list", middleware.Authorize, handler.getQrWalletStatuses)
	root.GET("/accountpriorities/list", middleware.Authorize, handler.getAccountPriorities)
	root.GET("/accountstatuses/list", middleware.Authorize, handler.getAccountStatuses)
	root.GET("/accountbotstatuses/list", middleware.Authorize, handler.getAccountBotStatuses)
	root.GET("/transfertypes/list", middleware.Authorize, handler.getTransferTypes)

	bankRoute := root.Group("/banks")
	bankRoute.GET("/list", middleware.Authorize, handler.getBanks)

	accountTypeRoute := root.Group("/accounttypes")
	accountTypeRoute.GET("/list", middleware.Authorize, handler.getAccountTypes)

	accountRoute := root.Group("/bankaccounts")
	accountRoute.GET("/list", middleware.Authorize, handler.getBankAccounts)
	accountRoute.GET("/detail/:id", middleware.Authorize, handler.getBankAccountById)
	accountRoute.POST("", middleware.Authorize, handler.createBankAccount)
	accountRoute.PATCH("/:id", middleware.Authorize, handler.updateBankAccount)
	accountRoute.DELETE("/:id", middleware.Authorize, handler.deleteBankAccount)

	account2Route := root.Group("/bankaccounts2")
	account2Route.GET("/settings", middleware.Authorize, handler.getExternalSettings)
	account2Route.GET("/customeraccount", middleware.Authorize, handler.getCustomerAccountsInfo)
	account2Route.GET("/list", middleware.Authorize, handler.getExternalAccounts)
	account2Route.GET("/status/:account", middleware.Authorize, handler.getExternalAccountStatus)
	account2Route.GET("/balance/:account", middleware.Authorize, handler.getExternalAccountBalance)
	account2Route.POST("", middleware.Authorize, handler.createExternalAccount)
	account2Route.PUT("", middleware.Authorize, handler.updateExternalAccount)
	account2Route.PUT("/status", middleware.Authorize, handler.EnableExternalAccount)
	account2Route.DELETE("/:account", middleware.Authorize, handler.deleteExternalAccount)
	account2Route.POST("/transfer", middleware.Authorize, handler.transferExternalAccount)
	account2Route.GET("/logs", middleware.Authorize, handler.getExternalAccountLogs)
	account2Route.GET("/statements", middleware.Authorize, handler.getExternalAccountStatements)
	account2Route.POST("/config", middleware.Authorize, handler.createBotaccountConfig)

	webhookRoute := root.Group("/webhooks")
	webhookRoute.POST("/action", handler.webhookAction)
	webhookRoute.POST("/noti", handler.webhookNoti)

	transactionRoute := root.Group("/transactions")
	transactionRoute.GET("/list", middleware.Authorize, handler.getTransactions)
	transactionRoute.GET("/detail/:id", middleware.Authorize, handler.getTransactionById)
	transactionRoute.POST("", middleware.Authorize, handler.createTransaction)
	transactionRoute.DELETE("/:id", middleware.Authorize, handler.deleteTransaction)

	transferRoute := root.Group("/transfers")
	transferRoute.GET("/list", middleware.Authorize, handler.getTransfers)
	transferRoute.GET("/detail/:id", middleware.Authorize, handler.getTransferById)
	transferRoute.POST("", middleware.Authorize, handler.createTransfer)
	transferRoute.POST("/confirm/:id", middleware.Authorize, handler.confirmTransfer)
	transferRoute.DELETE("/:id", middleware.Authorize, handler.deleteTransfer)

	statementRoute := root.Group("/statements")
	statementRoute.GET("/list", middleware.Authorize, handler.getAccountStatements)
	statementRoute.GET("/detail/:id", middleware.Authorize, handler.getAccountStatementById)
	statementRoute.POST("/webhook", middleware.Authorize, handler.addAccountStatementToWebhook)

}

// @Summary getBanks get Bank List
// @Description ดึงข้อมูลตัวเลือก รายชื่อธนาคารทั้งหมด
// @Tags Accounting - Options
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Param search query string false "search"
// @Param sortCol query string false "sortCol"
// @Param sortAsc query string false "sortAsc"
// @Success 200 {object} model.SuccessWithPagination
// @Router /accounting/banks/list [get]
func (h accountingController) getBanks(c *gin.Context) {

	var query model.BankListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.accountingService.GetBanks(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary getAccountTypes
// @Description ดึงข้อมูลตัวเลือก ประเภทบัญชีธนาคารทั้งหมด
// @Tags Accounting - Options
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} model.SuccessWithPagination
// @Router /accounting/accounttypes/list [get]
func (h accountingController) getAccountTypes(c *gin.Context) {

	var query model.AccountTypeListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.accountingService.GetAccountTypes(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary getAutoCreditFlags
// @Description ดึงข้อมูลตัวเลือก การตั้งค่าปรับเครดิตอัตโนมัติ
// @Tags Accounting - Options
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} model.SuccessWithPagination
// @Router /accounting/autocreditflags/list [get]
func (h accountingController) getAutoCreditFlags(c *gin.Context) {
	var data = []model.SimpleOption{
		{Key: "manual", Name: "สร้างใบงานและปรับเครดิตเอง"},
		{Key: "auto", Name: "ปรับเครดิตออโต้ (Bot)"},
	}
	c.JSON(200, model.SuccessWithPagination{List: data, Total: 2})
}

// @Summary getAutoWithdrawFlags
// @Description ดึงข้อมูลตัวเลือก การตั้งค่าถอนโอนเงินอัตโนมัติ
// @Tags Accounting - Options
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} model.SuccessWithPagination
// @Router /accounting/autowithdrawflags/list [get]
func (h accountingController) getAutoWithdrawFlags(c *gin.Context) {
	var data = []model.SimpleOption{
		{Key: "manual", Name: "สร้างใบงานและปรับเครดิตเอง"},
		{Key: "auto_backoffice", Name: "บัญชีถอนหลัก ปรับเครดิตออโต้ คลิกผ่านระบบหลังบ้าน"},
		{Key: "auto_bot", Name: "บัญชีถอนหลัก ปรับเครดิตออโต้ โอนเงินออโต้ (Bot)"},
	}
	c.JSON(200, model.SuccessWithPagination{List: data, Total: 3})
}

// @Summary getQrWalletStatuses
// @Description ดึงข้อมูลตัวเลือก การเปิดใช้งาน QR Wallet
// @Tags Accounting - Options
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} model.SuccessWithPagination
// @Router /accounting/qrwalletstatuses/list [get]
func (h accountingController) getQrWalletStatuses(c *gin.Context) {
	var data = []model.SimpleOption{
		{Key: "use_qr", Name: "เปิด"},
		{Key: "disabled", Name: "ปิด"},
	}
	c.JSON(200, model.SuccessWithPagination{List: data, Total: 2})
}

// @Summary getAccountStatuses
// @Description ดึงข้อมูลตัวเลือก สถานะบัญชีธนาคาร
// @Tags Accounting - Options
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} model.SuccessWithPagination
// @Router /accounting/accountstatuses/list [get]
func (h accountingController) getAccountStatuses(c *gin.Context) {
	var data = []model.SimpleOption{
		{Key: "active", Name: "ใช้งาน"},
		{Key: "deactive", Name: "ระงับการใช้งาน"},
	}
	c.JSON(200, model.SuccessWithPagination{List: data, Total: 2})
}

// @Summary getAccountPriorities
// @Description ดึงข้อมูลตัวเลือก ลำดับความสำคัญบัญชีธนาคาร
// @Tags Accounting - Options
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} model.SuccessWithPagination
// @Router /accounting/accountpriorities/list [get]
func (h accountingController) getAccountPriorities(c *gin.Context) {

	// var data = []model.SimpleOption{
	// 	{Key: "new", Name: "ระดับ NEW ทั่วไป"},
	// 	{Key: "gold", Name: "ระดับ Gold ฝากมากกว่า 10 ครั้ง"},
	// 	{Key: "platinum", Name: "ระดับ Platinum ฝากมากกว่า 20 ครั้ง"},
	// 	{Key: "vip", Name: "ระดับ VIP ฝากมากกว่า 20 ครั้ง"},
	// 	{Key: "classic", Name: "ระดับ CLASSIC ฝากสะสมมากกว่า 1,000 บาท"},
	// 	{Key: "superior", Name: "ระดับ SUPERIOR ฝากสะสมมากกว่า 10,000 บาท"},
	// 	{Key: "deluxe", Name: "ระดับ DELUXE ฝากสะสมมากกว่า 100,000 บาท"},
	// 	{Key: "wisdom", Name: "ระดับ WISDOM ฝากสะสมมากกว่า 500,000 บาท"},
	// }
	// c.JSON(200, model.SuccessWithPagination{List: data, Total: 8})

	data, err := h.accountingService.GetBankAccountPriorities()
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary getAccountBotStatuses
// @Description ดึงข้อมูลตัวเลือก สถานะบอท
// @Tags Accounting - Options
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} model.SuccessWithPagination
// @Router /accounting/accountbotstatuses/list [get]
func (h accountingController) getAccountBotStatuses(c *gin.Context) {
	var data = []model.SimpleOption{
		{Key: "active", Name: "เชื่อมต่อ"},
		{Key: "disconnected", Name: "ไม่ได้เชื่อมต่อ"},
	}
	c.JSON(200, model.SuccessWithPagination{List: data, Total: 2})
}

// @Summary get Transfer Types
// @Description ดึงข้อมูลตัวเลือก ประเภทการทำธุรกรรม (ฝาก/ถอน)
// @Tags Accounting - Options
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} model.SuccessWithPagination
// @Router /accounting/transfertypes/list [get]
func (h accountingController) getTransferTypes(c *gin.Context) {
	var data = []model.SimpleOption{
		{Key: "deposit", Name: "ฝากเงิน"},
		{Key: "withdraw", Name: "ถอนเงิน"},
	}
	c.JSON(200, model.SuccessWithPagination{List: data, Total: 2})
}

// @Summary GetBankAccountList
// @Description ดึงข้อมูลลิสบัญชีธนาคาร ใช้แสดงในหน้า จัดการธนาคาร
// @Tags Accounting - Bank Accounts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.BankAccountListRequest true "BankAccountListRequest"
// @Success 200 {object} model.SuccessWithPagination
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/bankaccounts/list [get]
func (h accountingController) getBankAccounts(c *gin.Context) {

	var query model.BankAccountListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.accountingService.GetBankAccounts(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary GetBankAccountById
// @Description ดึงข้อมูลบัญชีธนาคาร ด้วย id
// @Tags Accounting - Bank Accounts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/bankaccounts/detail/{id} [get]
func (h accountingController) getBankAccountById(c *gin.Context) {

	var accounting model.BankGetByIdRequest
	if err := c.ShouldBindUri(&accounting); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.accountingService.GetBankAccountById(accounting)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary CreateBankAccount
// @Description สร้าง บัญชีธนาคาร ใหม่ ในหน้า จัดการธนาคาร
// @Tags Accounting - Bank Accounts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.BankAccountCreateBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/bankaccounts [post]
func (h accountingController) createBankAccount(c *gin.Context) {

	var accounting model.BankAccountCreateBody
	if err := c.ShouldBindJSON(&accounting); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(accounting); err != nil {
		HandleError(c, err)
		return
	}

	err := h.accountingService.CreateBankAccount(accounting)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(201, model.Success{Message: "Created success"})
}

// @Summary UpdateBankAccount
// @Description แก้ไข บัญชีธนาคาร ในหน้า จัดการธนาคาร
// @Tags Accounting - Bank Accounts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param body body model.BankAccountUpdateRequest true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/bankaccounts/{id} [patch]
func (h accountingController) updateBankAccount(c *gin.Context) {

	id := c.Param("id")
	identifier, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		HandleError(c, err)
		return
	}

	body := model.BankAccountUpdateRequest{}

	if err := c.ShouldBindJSON(&body); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(body); err != nil {
		HandleError(c, err)
		return
	}

	if err := h.accountingService.UpdateBankAccount(identifier, body); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Updated success"})
}

// @Summary DeleteBankAccount
// @Description ลบข้อมูลบัญชีธนาคาร ด้วย id
// @Tags Accounting - Bank Accounts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/bankaccounts/{id} [delete]
func (h accountingController) deleteBankAccount(c *gin.Context) {

	id := c.Param("id")
	identifier, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		HandleError(c, err)
		return
	}

	delErr := h.accountingService.DeleteBankAccount(identifier)
	if delErr != nil {
		HandleError(c, delErr)
		return
	}
	c.JSON(201, model.Success{Message: "Deleted success"})
}

// @Summary GetTransactionList
// @Description ดึงข้อมูลลิสธุรกรรม ใช้แสดงในหน้า จัดการธนาคาร - ธุรกรรม และ รายการฝากถอนเงินสด
// @Tags Accounting - Bank Account Transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.BankAccountTransactionListRequest true "BankAccountTransactionListRequest"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/transactions/list [get]
func (h accountingController) getTransactions(c *gin.Context) {

	var query model.BankAccountTransactionListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.accountingService.GetTransactions(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary GetTransactionById
// @Description ดึงข้อมูลธุรกรรมด้วย id *ยังไม่ได้ใช้งาน*
// @Tags Accounting - Bank Account Transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/transactions/detail/{id} [get]
func (h accountingController) getTransactionById(c *gin.Context) {

	var accounting model.BankGetByIdRequest

	if err := c.ShouldBindUri(&accounting); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.accountingService.GetTransactionById(accounting)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary CreateTransaction
// @Description สร้าง ธุรกรรม ในหน้า จัดการธนาคาร - ธุรกรรม ส่ง AccountId มาด้วย
// @Tags Accounting - Bank Account Transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.BankAccountTransactionBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/transactions [post]
func (h accountingController) createTransaction(c *gin.Context) {

	username, err := h.accountingService.CheckCurrentUsername(c.MustGet("username"))
	if err != nil {
		HandleError(c, err)
		return
	}

	var accounting model.BankAccountTransactionBody
	accounting.CreatedByUsername = *username
	if err := c.ShouldBindJSON(&accounting); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(accounting); err != nil {
		HandleError(c, err)
		return
	}

	if err := h.accountingService.CreateTransaction(accounting); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(201, model.Success{Message: "Created success"})
}

// @Summary DeleteTransaction
// @Description ลบข้อมูลธุรกรรมด้วย id ใช้ในหน้า จัดการธนาคาร - ธุรกรรม ส่งรหัสผ่านมาเพื่อยืนยันด้วย
// @Tags Accounting - Bank Account Transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param body body model.ConfirmRequest true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/transactions/{id} [delete]
func (h accountingController) deleteTransaction(c *gin.Context) {

	adminId, err := h.accountingService.CheckCurrentAdminId(c.MustGet("adminId"))
	if err != nil {
		HandleError(c, err)
		return
	}

	id := c.Param("id")
	identifier, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		HandleError(c, err)
		return
	}

	var confirmation model.ConfirmRequest
	confirmation.UserId = *adminId
	if err := c.ShouldBindJSON(&confirmation); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(confirmation); err != nil {
		HandleError(c, err)
		return
	}
	if _, err := h.accountingService.CheckConfirmationPassword(confirmation); err != nil {
		HandleError(c, err)
		return
	}

	delErr := h.accountingService.DeleteTransaction(identifier)
	if delErr != nil {
		HandleError(c, delErr)
		return
	}
	c.JSON(201, model.Success{Message: "Deleted success"})
}

// @Summary GetTransferList
// @Description ดึงข้อมูลลิสการโอนเงิน ใช้แสดงในหน้า จัดการธนาคาร - ธุรกรรม
// @Tags Accounting - Bank Account Transfers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.BankAccountTransferListRequest true "BankAccountTransferListRequest"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/transfers/list [get]
func (h accountingController) getTransfers(c *gin.Context) {

	var query model.BankAccountTransferListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.accountingService.GetTransfers(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary GetTransferByID
// @Description ดึงข้อมูลการโอนด้วย id *ยังไม่ได้ใช้งาน*
// @Tags Accounting - Bank Account Transfers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/transfers/detail/{id} [get]
func (h accountingController) getTransferById(c *gin.Context) {

	var accounting model.BankGetByIdRequest
	if err := c.ShouldBindUri(&accounting); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.accountingService.GetTransferById(accounting)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary CreateTransfer
// @Description สร้างข้อมูลการโอน
// @Tags Accounting - Bank Account Transfers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.BankAccountTransferBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/transfers [post]
func (h accountingController) createTransfer(c *gin.Context) {

	username, err := h.accountingService.CheckCurrentUsername(c.MustGet("username"))
	if err != nil {
		HandleError(c, err)
		return
	}

	var accounting model.BankAccountTransferBody
	accounting.CreatedByUsername = *username
	if err := c.ShouldBindJSON(&accounting); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(accounting); err != nil {
		HandleError(c, err)
		return
	}

	if err := h.accountingService.CreateTransfer(accounting); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(201, model.Success{Message: "Created success"})
}

// @Summary ConfirmTransfer
// @Description ยืนยันการโอน
// @Tags Accounting - Bank Account Transfers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/transfers/confirm/{id} [post]
func (h accountingController) confirmTransfer(c *gin.Context) {

	adminId, err := h.accountingService.CheckCurrentAdminId(c.MustGet("adminId"))
	if err != nil {
		HandleError(c, err)
		return
	}

	id := c.Param("id")
	identifier, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		HandleError(c, err)
		return
	}

	if err := h.accountingService.ConfirmTransfer(identifier, *adminId); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(201, model.Success{Message: "Updated success"})
}

// @Summary DeleteTransfer
// @Description ลบข้อมูลการโอนด้วย id ใช้ในหน้า จัดการธนาคาร - ธุรกรรม ส่งรหัสผ่านมาเพื่อยืนยันด้วย
// @Tags Accounting - Bank Account Transfers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param body body model.ConfirmRequest true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/transfers/{id} [delete]
func (h accountingController) deleteTransfer(c *gin.Context) {

	adminId, err := h.accountingService.CheckCurrentAdminId(c.MustGet("adminId"))
	if err != nil {
		HandleError(c, err)
		return
	}

	id := c.Param("id")
	identifier, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		HandleError(c, err)
		return
	}

	var confirmation model.ConfirmRequest
	confirmation.UserId = *adminId
	if err := c.ShouldBindJSON(&confirmation); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(confirmation); err != nil {
		HandleError(c, err)
		return
	}
	if _, err := h.accountingService.CheckConfirmationPassword(confirmation); err != nil {
		HandleError(c, err)
		return
	}

	delErr := h.accountingService.DeleteTransfer(identifier)
	if delErr != nil {
		HandleError(c, delErr)
		return
	}
	c.JSON(201, model.Success{Message: "Deleted success"})
}

// @Summary getAccountStatements รายการเดินบัญชีธนาคาร
// @Description ดึงข้อมูล Statement รายการเดินบัญชีธนาคาร จาก FASTBANK ตรงๆ
// @Tags Accounting - Bank Account Statements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.BankAccountStatementListRequest true "BankAccountStatementListRequest"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/statements/list [get]
func (h accountingController) getAccountStatements(c *gin.Context) {

	var query model.BankAccountStatementListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.accountingService.GetAccountStatements(query)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary GetAccountStatementById
// @Description ดึงข้อมูลการโอนด้วย id *ยังไม่ได้ใช้งาน*
// @Tags Accounting - Bank Account Statements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/statements/detail/{id} [get]
func (h accountingController) getAccountStatementById(c *gin.Context) {

	var req model.BankGetByIdRequest
	if err := c.ShouldBindUri(&req); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.accountingService.GetAccountStatementById(req)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary AddAccountStatementToWebhook
// @Description สั่งให้ยิง Statement ใหม่ ไปที่ Webhook ที่ตั้งไว้
// @Tags Accounting - Bank Account Statements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.RecheckWebhookRequest true "body"
// @Success 200 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/statements/webhook [post]
func (h accountingController) addAccountStatementToWebhook(c *gin.Context) {

	var reqCreate model.RecheckWebhookRequest
	if err := c.ShouldBindJSON(&reqCreate); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(reqCreate); err != nil {
		HandleError(c, err)
		return
	}

	if err := h.accountingService.AddAccountStatementToWebhook(reqCreate); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(200, model.Success{Message: "success"})
}

// @Summary GetCustomerAccountsInfo เช็คชื่อบัญชีธนาคารลูกค้า
// @Description ดึงข้อมูลบัญชีธนาคารของลูกค้า เพื่อเช็คชื่อบัญชีธนาคาร
// @Tags Accounting - FASTBANK
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.CustomerAccountInfoRequest true "CustomerAccountInfoRequest"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/bankaccounts2/customeraccount [post]
func (h accountingController) getCustomerAccountsInfo(c *gin.Context) {

	var query model.CustomerAccountInfoRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.accountingService.GetCustomerAccountsInfo(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary GetExternalSettings
// @Description อัพเดทข้อมูล บัญชีธนาคารบอท ด้วยเลขบัญชี
// @Tags Accounting - FASTBANK
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/bankaccounts2/settings [get]
func (h accountingController) getExternalSettings(c *gin.Context) {

	data, err := h.accountingService.GetExternalSettings()
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary GetExternalAccounts
// @Description ดึงข้อมูลลิสบัญชีธนาคาร บอท
// @Tags Accounting - FASTBANK
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.BankAccountListRequest true "BankAccountListRequest"
// @Success 200 {object} model.SuccessWithPagination
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/bankaccounts2/list [get]
func (h accountingController) getExternalAccounts(c *gin.Context) {

	// var query model.BankAccountListRequest
	// if err := c.ShouldBind(&query); err != nil {
	// 	HandleError(c, err)
	// 	return
	// }
	// if err := validator.New().Struct(query); err != nil {
	// 	HandleError(c, err)
	// 	return
	// }

	data, err := h.accountingService.GetExternalAccounts()
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary GetExternalAccountStatus
// @Description ดึงข้อมูล บัญชีธนาคารบอท ด้วยเลขบัญชี
// @Tags Accounting - FASTBANK
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param account path string true "accountNumber"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/bankaccounts2/status/{account} [get]
func (h accountingController) getExternalAccountStatus(c *gin.Context) {

	var query model.ExternalAccountStatusRequest
	query.AccountNumber = c.Param("account")

	data, err := h.accountingService.GetExternalAccountStatus(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary GetExternalAccountBalance
// @Description ดึงข้อมูล บัญชีธนาคารบอท ด้วยเลขบัญชี
// @Tags Accounting - FASTBANK
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param account path string true "accountNumber"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/bankaccounts2/balance/{account} [get]
func (h accountingController) getExternalAccountBalance(c *gin.Context) {

	var query model.ExternalAccountStatusRequest
	query.AccountNumber = c.Param("account")

	data, err := h.accountingService.GetExternalAccountBalance(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary CreateExternalAccount
// @Description สร้าง บัญชีธนาคารภายนอก ใหม่
// @Tags Accounting - FASTBANK
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.ExternalAccountCreateBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/bankaccounts2 [post]
func (h accountingController) createExternalAccount(c *gin.Context) {

	var accounting model.ExternalAccountCreateBody
	if err := c.ShouldBindJSON(&accounting); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(accounting); err != nil {
		HandleError(c, err)
		return
	}

	_, err := h.accountingService.CreateExternalAccount(accounting)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(201, model.Success{Message: "Created success"})
}

// @Summary UpdateExternalAccount
// @Description อัพเดทข้อมูล บัญชีธนาคารบอท ด้วยเลขบัญชี
// @Tags Accounting - FASTBANK
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.ExternalAccountCreateBody true "body"
// @Success 200 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/bankaccounts2/ [put]
func (h accountingController) updateExternalAccount(c *gin.Context) {

	var query model.ExternalAccountUpdateBody
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	_, err := h.accountingService.UpdateExternalAccount(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.Success{Message: "Update success"})
}

// @Summary EnableExternalAccount
// @Description เปิด ปิด สถานะบัญชีธนาคารบอท ด้วยเลขบัญชี
// @Tags Accounting - FASTBANK
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.ExternalAccountEnableRequest true "body"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/bankaccounts2/status [put]
func (h accountingController) EnableExternalAccount(c *gin.Context) {

	var query model.ExternalAccountEnableRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	result, err := h.accountingService.EnableExternalAccount(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "Update success", Data: result})
}

// @Summary DeleteExternalAccount
// @Description ลบข้อมูล บัญชีธนาคารบอท ด้วยเลขบัญชี
// @Tags Accounting - FASTBANK
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param account path string true "accountNumber"
// @Success 200 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/bankaccounts2/{account} [delete]
func (h accountingController) deleteExternalAccount(c *gin.Context) {

	var query model.ExternalAccountStatusRequest
	query.AccountNumber = c.Param("account")

	err := h.accountingService.DeleteExternalAccount(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.Success{Message: "Delete success"})
}

// @Summary TransferExternalAccount
// @Description โอนเงิน บัญชีธนาคารบอท
// @Tags Accounting - FASTBANK
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.ExternalAccountTransferRequest true "body"
// @Success 200 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/bankaccounts2/transfer [post]
func (h accountingController) transferExternalAccount(c *gin.Context) {

	var query model.ExternalAccountTransferRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	err := h.accountingService.TransferExternalAccount(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.Success{Message: "Transfer success"})
}

// @Summary GetExternalAccountLogs
// @Description ดึงข้อมูล Logs บัญชีธนาคารบอท
// @Tags Accounting - FASTBANK
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.ExternalStatementListRequest true "ExternalStatementListRequest"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/bankaccounts2/logs [get]
func (h accountingController) getExternalAccountLogs(c *gin.Context) {

	var query model.ExternalStatementListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.accountingService.GetExternalAccountLogs(query)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary GetExternalAccountStatements
// @Description ดึงข้อมูล Statement บัญชีธนาคารบอท
// @Tags Accounting - FASTBANK
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.ExternalStatementListRequest true "ExternalStatementListRequest"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/bankaccounts2/statements [get]
func (h accountingController) getExternalAccountStatements(c *gin.Context) {

	var query model.ExternalStatementListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.accountingService.GetExternalAccountStatements(query)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary WebhookAction
// @Description เว็บฮุคแบบ GET
// @Tags Accounting - FASTBANK
// @Accept json
// @Produce json
// @Param body body model.ExternalAccountEnableRequest true "body"
// @Success 200 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/webhooks/action [post]
func (h accountingController) webhookAction(c *gin.Context) {

	jsonData, errValidate := io.ReadAll(c.Request.Body)
	if errValidate != nil {
		HandleError(c, errValidate)
		return
	}

	// IsNewStateMentList
	var resp model.WebhookStatementResponse
	errJson := json.Unmarshal(jsonData, &resp)
	if errJson != nil {
		HandleError(c, errJson)
		return
	}
	jsonString := string(jsonData)
	insertId, err := h.accountingService.CreateWebhookLog("ACTION", jsonString)
	if err != nil {
		HandleError(c, err)
		return
	}
	var updateReq model.WebhookLogUpdateBody
	updateReq.Status = "success"
	updateReq.JsonPayload = "{}"

	// Do Work After
	if resp.NewStatementList != nil {
		for _, v := range resp.NewStatementList {
			err := h.accountingService.CreateBankStatementFromWebhook(v)
			if err != nil {
				updateReq.Status = err.Error()
				// later : support many error and payload info
				// updateReq.Status = "error"
				// updateReq.JsonPayload = err.Error()
			}
		}
	}
	updateReq.JsonPayload = "{}"

	// Update WebhookLog
	if updateReq.Status == "success" {
		if err = h.accountingService.SetSuccessWebhookLog(*insertId, "{}"); err != nil {
			HandleError(c, err)
			return
		}
	} else {
		if err = h.accountingService.SetFailedWebhookLog(*insertId, updateReq.Status); err != nil {
			HandleError(c, err)
			return
		}
	}
	c.JSON(200, model.Success{Message: "success"})
}

// @Summary WebhookNoti
// @Description เว็บฮุคแบบ POST
// @Tags Accounting - FASTBANK
// @Accept json
// @Produce json
// @Param body body model.ExternalAccountEnableRequest true "body"
// @Success 200 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/webhooks/noti [post]
func (h accountingController) webhookNoti(c *gin.Context) {

	jsonData, errValidate := io.ReadAll(c.Request.Body)
	if errValidate != nil {
		HandleError(c, errValidate)
		return
	}

	jsonString := string(jsonData)
	_, err := h.accountingService.CreateWebhookLog("NOTI", jsonString)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(200, model.Success{Message: "success"})
}

// @Summary CreateBotaccountConfig
// @Description เพิ่ม การตั้งค่าบัญชีธนาคาร ใหม่
// @Tags Accounting - FASTBANK
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.BotAccountConfigCreateBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/bankaccounts2/config [post]
func (h accountingController) createBotaccountConfig(c *gin.Context) {

	var reqCreate model.BotAccountConfigCreateBody
	if err := c.ShouldBindJSON(&reqCreate); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(reqCreate); err != nil {
		HandleError(c, err)
		return
	}

	err := h.accountingService.CreateBotaccountConfig(reqCreate)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(201, model.Success{Message: "Created success"})
}
