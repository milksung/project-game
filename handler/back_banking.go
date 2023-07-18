package handler

import (
	"cybergame-api/middleware"
	"cybergame-api/model"
	"cybergame-api/repository"
	"cybergame-api/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type bankingController struct {
	bankingService    service.BankingService
	accountingService service.AccountingService
}

func newBankingController(
	bankingService service.BankingService,
	accountingService service.AccountingService,
) bankingController {
	return bankingController{bankingService, accountingService}
}

func BankingController(r *gin.RouterGroup, db *gorm.DB) {

	repoBanking := repository.NewBankingRepository(db)
	repoAccounting := repository.NewAccountingRepository(db)
	repoAgentConnect := repository.NewAgentConnectRepository(db)
	service1 := service.NewBankingService(repoBanking, repoAccounting, repoAgentConnect)
	service2 := service.NewAccountingService(repoAccounting)
	handler := newBankingController(service1, service2)

	root := r.Group("/banking")
	root.GET("/transactiontypes/list", middleware.Authorize, handler.getTransactionTypes)
	root.GET("/transactionstatuses/list", middleware.Authorize, handler.getTransactionStatuses)
	root.GET("/statementtypes/list", middleware.Authorize, handler.getStatementTypes)

	statementRoute := root.Group("/statements")
	statementRoute.GET("/list", middleware.Authorize, handler.getBankStatements)
	statementRoute.GET("/summary", middleware.Authorize, handler.getBankStatementSummary)
	statementRoute.GET("/detail/:id", middleware.Authorize, handler.getBankStatementById)
	statementRoute.POST("", middleware.Authorize, handler.createBankStatement)
	statementRoute.POST("/matchowner/:id", middleware.Authorize, handler.matchStatementOwner)
	statementRoute.POST("/ignoreowner/:id", middleware.Authorize, handler.ignoreStatementOwner)
	statementRoute.DELETE("/:id", middleware.Authorize, handler.deleteBankStatement)

	transactionRoute := root.Group("/transactions")
	transactionRoute.GET("/list", middleware.Authorize, handler.getBankTransactions)
	transactionRoute.GET("/deposit_list", middleware.Authorize, handler.getBankDepositTransactions)
	transactionRoute.GET("/withdraw_list", middleware.Authorize, handler.getBankWithdrawTransactions)
	transactionRoute.GET("/count_deposit_statuses", middleware.Authorize, handler.getBankDepositTransStatusCounts)
	transactionRoute.GET("/count_withdraw_statuses", middleware.Authorize, handler.getBankWithdrawTransStatusCounts)
	transactionRoute.GET("/detail/:id", middleware.Authorize, handler.getBankTransactionById)
	transactionRoute.POST("", middleware.Authorize, handler.createBankTransaction)
	transactionRoute.PATCH("/:id", middleware.Authorize, handler.updateBankTransaction)
	transactionRoute.POST("/bonus", middleware.Authorize, handler.createBonusTransaction)

	transactionRoute.GET("/pendingdepositlist", middleware.Authorize, handler.getPendingDepositTransactions)
	transactionRoute.GET("/pendingwithdrawlist", middleware.Authorize, handler.getPendingWithdrawTransactions)
	transactionRoute.POST("/cancel/:id", middleware.Authorize, handler.cancelPendingTransaction)
	transactionRoute.POST("/confirmdeposit/:id", middleware.Authorize, handler.confirmDepositTransaction)
	transactionRoute.POST("/confirmdepositcredit/:id", middleware.Authorize, handler.confirmDepositCreditTransaction)
	transactionRoute.POST("/confirmcreditwithdraw/:id", middleware.Authorize, handler.confirmCreditWithdrawTransaction)
	transactionRoute.POST("/confirmtransferwithdraw/:id", middleware.Authorize, handler.confirmTransferWithdrawTransaction)
	transactionRoute.POST("/continueautowithdraw/:id", middleware.Authorize, handler.continueAutoWithdrawTransaction)
	transactionRoute.GET("/finishedlist", middleware.Authorize, handler.getFinishedTransactions)
	transactionRoute.POST("/remove/:id", middleware.Authorize, handler.removeFinishedTransaction)
	transactionRoute.GET("/removedlist", middleware.Authorize, handler.getRemovedTransactions)

	memberRoute := root.Group("/member")
	memberRoute.GET("/info/:code", middleware.Authorize, handler.getMemberByCode)
	memberRoute.GET("/list", middleware.Authorize, handler.getMembers)
	memberRoute.GET("/possibleowner", middleware.Authorize, handler.getPossibleStatementOwners)
	memberRoute.GET("/transactionsummary", middleware.Authorize, handler.getMemberTransactionSummary)
	memberRoute.GET("/transactions", middleware.Authorize, handler.getMemberTransactions)
	memberRoute.GET("/statements", middleware.Authorize, handler.getMemberStatements)
	memberRoute.GET("/statements/detail/:id", middleware.Authorize, handler.getMemberStatementById)
	// TEST
	memberRoute.POST("/statements1", middleware.Authorize, handler.processMemberDepositCredit)
	memberRoute.POST("/statements2", middleware.Authorize, handler.processMemberWithdrawCredit)
	memberRoute.POST("/statements3", middleware.Authorize, handler.processMemberBonusCredit)
	memberRoute.POST("/statements4", middleware.Authorize, handler.processMemberGetbackCredit)

}

// @Summary get Transaction Type List
// @Description ดึงข้อมูลตัวเลือก ประเภทการทำรายการ
// @Tags Banking - Options
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} model.SuccessWithPagination
// @Router /banking/transactiontypes/list [get]
func (h bankingController) getTransactionTypes(c *gin.Context) {
	var data = []model.SimpleOption{
		{Key: "deposit", Name: "ฝาก"},
		{Key: "withdraw", Name: "ถอน"},
		{Key: "getcreditback", Name: "ดึงเครดิตกลับ"},
	}
	c.JSON(200, model.SuccessWithPagination{List: data, Total: 2})
}

// @Summary get Transaction Status List
// @Description ดึงข้อมูลตัวเลือก สถานะรายการฝากถอน
// @Tags Banking - Options
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} model.SuccessWithPagination
// @Router /banking/transactionstatuses/list [get]
func (h bankingController) getTransactionStatuses(c *gin.Context) {
	var data = []model.SimpleOption{
		{Key: "pending", Name: "รอดำเนินการ"},
		{Key: "canceled", Name: "ยกเลิกแล้ว"},
		{Key: "finished", Name: "อนุมัติแล้ว"},
	}
	c.JSON(200, model.SuccessWithPagination{List: data, Total: 2})
}

// @Summary get Statement type List
// @Description ดึงข้อมูลตัวเลือก ประเภทรายการเดินบัญชี
// @Tags Banking - Options
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} model.SuccessWithPagination
// @Router /banking/statementtypes/list [get]
func (h bankingController) getStatementTypes(c *gin.Context) {
	var data = []model.SimpleOption{
		{Key: "transfer_in", Name: "โอนเข้า"},
		{Key: "transfer_out", Name: "โอนออก"},
	}
	c.JSON(200, model.SuccessWithPagination{List: data, Total: 2})
}

// @Summary GetStatementList รายการเดินบัญชี รายการโอนรอดำเนินการ
// @Description ดึงข้อมูลลิสการโอนเงิน ใช้แสดงในหน้า จัดการธนาคาร - ธุรกรรม
// @Tags Banking - Bank Account Statements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.BankStatementListRequest true "BankStatementListRequest"
// @Success 200 {object} model.SuccessWithPagination
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/statements/list [get]
func (h bankingController) getBankStatements(c *gin.Context) {

	var query model.BankStatementListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetBankStatements(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary GetBankStatementSummary
// @Description ดึงข้อมูลจำนวนรายการสรุป เช่น รายการโอนรอดำเนินการ รายการฝากถอนรอดำเนินการ
// @Tags Banking - Bank Account Statements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/statements/summary [get]
func (h bankingController) getBankStatementSummary(c *gin.Context) {

	var query model.BankStatementListRequest

	data, err := h.bankingService.GetBankStatementSummary(query)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary GetStatementByID
// @Description ดึงข้อมูลการโอนด้วย id *ยังไม่ได้ใช้งาน*
// @Tags Banking - Bank Account Statements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/statements/detail/{id} [get]
func (h bankingController) getBankStatementById(c *gin.Context) {

	var req model.GetByIdRequest

	if err := c.ShouldBindUri(&req); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetBankStatementById(req)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary CreateStatement
// @Description สร้างข้อมูลการโอน
// @Tags Banking - Bank Account Statements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.BankStatementCreateBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/statements [post]
func (h bankingController) createBankStatement(c *gin.Context) {

	var banking model.BankStatementCreateBody
	if err := c.ShouldBindJSON(&banking); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(banking); err != nil {
		HandleError(c, err)
		return
	}

	if err := h.bankingService.CreateBankStatement(banking); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(201, model.Success{Message: "Created success"})
}

// @Summary MatchStatementOwner ดำเนินการข้อมูลการเดินบัญชี
// @Description ดำเนินการข้อมูลการเดินบัญชี
// @Tags Banking - Bank Account Statements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param body body model.BankStatementMatchRequest true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/statements/matchowner/{id} [post]
func (h bankingController) matchStatementOwner(c *gin.Context) {

	id := c.Param("id")
	identifier, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		HandleError(c, err)
		return
	}

	adminId, err := h.accountingService.CheckCurrentAdminId(c.MustGet("adminId"))
	if err != nil {
		HandleError(c, err)
		return
	}
	username, err := h.accountingService.CheckCurrentUsername(c.MustGet("username"))
	if err != nil {
		HandleError(c, err)
		return
	}

	var req model.BankStatementMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(req); err != nil {
		HandleError(c, err)
		return
	}

	req.ConfirmedAt = time.Now()
	req.ConfirmedByUserId = *adminId
	req.ConfirmedByUsername = *username
	actionErr := h.bankingService.MatchStatementOwner(identifier, req)
	if actionErr != nil {
		HandleError(c, actionErr)
		return
	}
	c.JSON(201, model.Success{Message: "Confirm success"})
}

// @Summary IgnoreStatementOwner เพิกเฉยข้อมูลการเดินบัญชี
// @Description เพิกเฉยข้อมูลการเดินบัญชี
// @Tags Banking - Bank Account Statements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/statements/ignoreowner/{id} [post]
func (h bankingController) ignoreStatementOwner(c *gin.Context) {

	id := c.Param("id")
	identifier, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		HandleError(c, err)
		return
	}

	adminId, err := h.accountingService.CheckCurrentAdminId(c.MustGet("adminId"))
	if err != nil {
		HandleError(c, err)
		return
	}
	username, err := h.accountingService.CheckCurrentUsername(c.MustGet("username"))
	if err != nil {
		HandleError(c, err)
		return
	}

	var req model.BankStatementMatchRequest
	req.ConfirmedAt = time.Now()
	req.ConfirmedByUserId = *adminId
	req.ConfirmedByUsername = *username
	actionErr := h.bankingService.IgnoreStatementOwner(identifier, req)
	if actionErr != nil {
		HandleError(c, actionErr)
		return
	}
	c.JSON(201, model.Success{Message: "Ignore success"})
}

// @Summary DeleteBankStatement
// @Description ลบข้อมูลการโอนด้วย id ใช้ในหน้า จัดการธนาคาร - ธุรกรรม ส่งรหัสผ่านมาเพื่อยืนยันด้วย
// @Tags Banking - Bank Account Statements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/statements/{id} [delete]
func (h bankingController) deleteBankStatement(c *gin.Context) {

	id := c.Param("id")
	identifier, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		HandleError(c, err)
		return
	}

	actionErr := h.bankingService.DeleteBankStatement(identifier)
	if actionErr != nil {
		HandleError(c, actionErr)
		return
	}
	c.JSON(201, model.Success{Message: "Deleted success"})
}

// @Summary GetBankTransactions
// @Description ดึงข้อมูลลิสการฝากถอน
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.BankTransactionListRequest true "BankTransactionListRequest"
// @Success 200 {object} model.SuccessWithPagination
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/list [get]
func (h bankingController) getBankTransactions(c *gin.Context) {

	var query model.BankTransactionListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetBankTransactions(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary GetBankTransactions ดึงข้อมูลลิสการฝาก ที่จะมีทั้ง ฝาก และ โบนัส
// @Description ดึงข้อมูลลิสการฝาก ที่จะมีทั้ง ฝาก และ โบนัส
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.BankTransactionListRequest true "BankTransactionListRequest"
// @Success 200 {object} model.SuccessWithPagination
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/deposit_list [get]
func (h bankingController) getBankDepositTransactions(c *gin.Context) {

	var query model.BankTransactionListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetBankDepositTransactions(query)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary GetBankDepositTransStatusCounts ดึงข้อมูลจำนวนรายการฝากทั้งหมด ตามสถานะ
// @Description ดึงข้อมูลจำนวนรายการฝากทั้งหมด ตามสถานะ
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.BankTransactionListRequest true "BankTransactionListRequest"
// @Success 200 {object} model.SuccessWithPagination
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/count_deposit_statuses [get]
func (h bankingController) getBankDepositTransStatusCounts(c *gin.Context) {

	var query model.BankTransactionListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err2 := h.bankingService.GetBankDepositTransStatusCounts(query)
	if err2 != nil {
		HandleError(c, err2)
		return
	}
	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary GetBankTransactions ดึงข้อมูลลิสการถอน ที่จะมีทั้ง ถอน และ ดึงเครดิตกลับ
// @Description ดึงข้อมูลลิสการถอน ที่จะมีทั้ง ถอน และ ดึงเครดิตกลับ
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.BankTransactionListRequest true "BankTransactionListRequest"
// @Success 200 {object} model.SuccessWithPagination
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/withdraw_list [get]
func (h bankingController) getBankWithdrawTransactions(c *gin.Context) {

	var query model.BankTransactionListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetBankWithdrawTransactions(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary GetBankWithdrawTransStatusCounts ดึงข้อมูลจำนวนรายการถอนทั้งหมด ตามสถานะ
// @Description ดึงข้อมูลจำนวนรายการถอนทั้งหมด ตามสถานะ
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.BankTransactionListRequest true "BankTransactionListRequest"
// @Success 200 {object} model.SuccessWithPagination
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/count_withdraw_statuses [get]
func (h bankingController) getBankWithdrawTransStatusCounts(c *gin.Context) {

	var query model.BankTransactionListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetBankWithdrawTransStatusCounts(query)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary GetBankTransactionById
// @Description ดึงข้อมูลการฝากถอน ด้วย id
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/detail/{id} [get]
func (h bankingController) getBankTransactionById(c *gin.Context) {

	var req model.BankTransactionGetRequest

	if err := c.ShouldBindUri(&req); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetBankTransactionById(req)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary CreateBankTransaction สร้าง บันทึกรายการฝาก-ถอน ของลูกค้า
// @Description สร้างข้อมูล บันทึกรายการฝาก-ถอน
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.BankTransactionCreateBody true "*บังคับกรอก memberCode และ creditAmount และ transferType"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions [post]
func (h bankingController) createBankTransaction(c *gin.Context) {

	adminId, err := h.accountingService.CheckCurrentAdminId(c.MustGet("adminId"))
	if err != nil {
		HandleError(c, err)
		return
	}
	username, err := h.accountingService.CheckCurrentUsername(c.MustGet("username"))
	if err != nil {
		HandleError(c, err)
		return
	}

	var banking model.BankTransactionCreateBody
	if err := c.ShouldBindJSON(&banking); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(banking); err != nil {
		HandleError(c, err)
		return
	}
	banking.CreatedByUserId = *adminId
	banking.CreatedByUsername = *username
	if err := h.bankingService.CreateBankTransaction(banking); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(201, model.Success{Message: "Created success"})
}

// @Summary CreateBonusTransaction
// @Description สร้างข้อมูล บันทึกแจกโบนัส
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.BonusTransactionCreateBody true "body description"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/bonus [post]
func (h bankingController) createBonusTransaction(c *gin.Context) {

	username, err := h.accountingService.CheckCurrentUsername(c.MustGet("username"))
	if err != nil {
		HandleError(c, err)
		return
	}
	adminId, err := h.accountingService.CheckCurrentAdminId(c.MustGet("adminId"))
	if err != nil {
		HandleError(c, err)
		return
	}

	var banking model.BonusTransactionCreateBody
	if err := c.ShouldBindJSON(&banking); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(banking); err != nil {
		HandleError(c, err)
		return
	}
	banking.CreatedByUserId = *adminId
	banking.CreatedByUsername = *username

	if err := h.bankingService.CreateBonusTransaction(banking); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(201, model.Success{Message: "Created success"})
}

// @Summary UpdateBankTransaction อัพเดท บันทึกรายการฝาก-ถอน ของลูกค้า
// @Description อัพเดท บันทึกรายการฝาก-ถอน ส่งมาเฉพาะ ช่องที่อัพเดทเลย
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param body body model.BankTransactionUpdateRequest true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/{id} [patch]
func (h bankingController) updateBankTransaction(c *gin.Context) {

	adminId, err := h.accountingService.CheckCurrentAdminId(c.MustGet("adminId"))
	if err != nil {
		HandleError(c, err)
		return
	}
	username, err := h.accountingService.CheckCurrentUsername(c.MustGet("username"))
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

	var req model.BankTransactionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(req); err != nil {
		HandleError(c, err)
		return
	}
	req.UpdatedByUserId = *adminId
	req.UpdatedByUsername = *username

	if err := h.bankingService.UpdateBankTransaction(identifier, req); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(201, model.Success{Message: "Update success"})
}

// @Summary GetPendingDepositTransactions
// @Description ดึงข้อมูลลิสการฝาก ที่รออนุมัติ
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.PendingDepositTransactionListRequest true "query"
// @Success 200 {object} model.SuccessWithPagination
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/pendingdepositlist [get]
func (h bankingController) getPendingDepositTransactions(c *gin.Context) {

	var query model.PendingDepositTransactionListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetPendingDepositTransactions(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary GetPendingWithdrawTransactions
// @Description ดึงข้อมูลลิสการถอน ที่รออนุมัติ
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.PendingWithdrawTransactionListRequest true "query"
// @Success 200 {object} model.SuccessWithPagination
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/pendingwithdrawlist [get]
func (h bankingController) getPendingWithdrawTransactions(c *gin.Context) {

	var query model.PendingWithdrawTransactionListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetPendingWithdrawTransactions(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary CancelPendingTransaction ยกเลิก ข้อมูลการฝากและถอน ที่รอยืนยัน
// @Description ยกเลิก ไม่ยืนยัน ข้อมูลการฝากและถอน ที่รอยืนยัน
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param body body model.BankTransactionCancelBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/cancel/{id} [post]
func (h bankingController) cancelPendingTransaction(c *gin.Context) {

	adminId, err := h.accountingService.CheckCurrentAdminId(c.MustGet("adminId"))
	if err != nil {
		HandleError(c, err)
		return
	}
	username, err := h.accountingService.CheckCurrentUsername(c.MustGet("username"))
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

	var data model.BankTransactionCancelBody
	if err := c.ShouldBind(&data); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(data); err != nil {
		HandleError(c, err)
		return
	}

	data.Status = "canceled"
	// data.CancelRemark = data.CancelRemark
	data.CanceledAt = time.Now()
	data.CanceledByUserId = *adminId
	data.CanceledByUsername = *username

	actionErr := h.bankingService.CancelPendingTransaction(identifier, data)
	if actionErr != nil {
		HandleError(c, actionErr)
		return
	}
	c.JSON(201, model.Success{Message: "Cancel success"})
}

// @Summary ConfirmDepositTransaction ยืนยันข้อมูลการฝาก เฉยๆ
// @Description ยืนยันข้อมูลการฝาก เฉยๆ
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param body body model.BankConfirmDepositRequest true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/confirmdeposit/{id} [post]
func (h bankingController) confirmDepositTransaction(c *gin.Context) {

	username, err := h.accountingService.CheckCurrentUsername(c.MustGet("username"))
	if err != nil {
		HandleError(c, err)
		return
	}
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

	var req model.BankConfirmDepositRequest
	if err := c.ShouldBind(&req); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(req); err != nil {
		HandleError(c, err)
		return
	}
	req.ConfirmedAt = time.Now()
	req.ConfirmedByUserId = *adminId
	req.ConfirmedByUsername = *username
	actionErr := h.bankingService.ConfirmDepositTransaction(identifier, req)
	if actionErr != nil {
		HandleError(c, actionErr)
		return
	}
	c.JSON(201, model.Success{Message: "Confirm success"})
}

// @Summary ConfirmDepositCreditTransaction ยืนยันข้อมูลการฝาก เพื่ออนุมัติเครดิต
// @Description ยืนยันข้อมูลการฝาก เพื่ออนุมัติเครดิต
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param body body model.BankConfirmDepositRequest true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/confirmdepositcredit/{id} [post]
func (h bankingController) confirmDepositCreditTransaction(c *gin.Context) {

	username, err := h.accountingService.CheckCurrentUsername(c.MustGet("username"))
	if err != nil {
		HandleError(c, err)
		return
	}
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

	var req model.BankConfirmDepositRequest
	if err := c.ShouldBind(&req); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(req); err != nil {
		HandleError(c, err)
		return
	}
	req.ConfirmedAt = time.Now()
	req.ConfirmedByUserId = *adminId
	req.ConfirmedByUsername = *username
	actionCreditErr := h.bankingService.ConfirmDepositCredit(identifier, req)
	if actionCreditErr != nil {
		HandleError(c, actionCreditErr)
		return
	}
	c.JSON(201, model.Success{Message: "Confirm success"})
}

// @Summary confirmCreditWithdrawTransaction ยืนยัน/อนุมัติ รายการถอน ในสถานะ รอปรับเครดิต
// @Description ยืนยัน/อนุมัติ รายการถอน ในสถานะ รอปรับเครดิต
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param body body model.BankConfirmCreditWithdrawRequest true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/confirmcreditwithdraw/{id} [post]
func (h bankingController) confirmCreditWithdrawTransaction(c *gin.Context) {

	id := c.Param("id")
	identifier, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		HandleError(c, err)
		return
	}
	username, err := h.accountingService.CheckCurrentUsername(c.MustGet("username"))
	if err != nil {
		HandleError(c, err)
		return
	}
	adminId, err := h.accountingService.CheckCurrentAdminId(c.MustGet("adminId"))
	if err != nil {
		HandleError(c, err)
		return
	}

	var req model.BankConfirmCreditWithdrawRequest
	if err := c.ShouldBind(&req); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(req); err != nil {
		HandleError(c, err)
		return
	}
	req.ConfirmedAt = time.Now()
	req.ConfirmedByUserId = *adminId
	req.ConfirmedByUsername = *username

	actionErr := h.bankingService.ConfirmWithdrawTransaction(identifier, req)
	if actionErr != nil {
		HandleError(c, actionErr)
		return
	}
	c.JSON(201, model.Success{Message: "Confirm success"})
}

// @Summary confirmTransferWithdrawTransaction ยืนยัน/อนุมัติ รายการถอน ในสถานะ รอโอน
// @Description ยืนยัน/อนุมัติ รายการถอน ในสถานะ รอโอน
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param body body model.BankConfirmTransferWithdrawRequest true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/confirmtransferwithdraw/{id} [post]
func (h bankingController) confirmTransferWithdrawTransaction(c *gin.Context) {

	id := c.Param("id")
	identifier, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		HandleError(c, err)
		return
	}
	username, err := h.accountingService.CheckCurrentUsername(c.MustGet("username"))
	if err != nil {
		HandleError(c, err)
		return
	}
	adminId, err := h.accountingService.CheckCurrentAdminId(c.MustGet("adminId"))
	if err != nil {
		HandleError(c, err)
		return
	}

	var req model.BankConfirmTransferWithdrawRequest
	if err := c.ShouldBind(&req); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(req); err != nil {
		HandleError(c, err)
		return
	}
	req.ConfirmedAt = time.Now()
	req.ConfirmedByUserId = *adminId
	req.ConfirmedByUsername = *username

	actionErr := h.bankingService.ConfirmWithdrawTransfer(identifier, req)
	if actionErr != nil {
		HandleError(c, actionErr)
		return
	}
	c.JSON(201, model.Success{Message: "Confirm success"})
}

// @Summary continueAutoWithdrawTransaction ให้ระบบ
// @Description ยืนยัน/อนุมัติ รายการถอน ในสถานะ รอโอน
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/continueautowithdraw/{id} [post]
func (h bankingController) continueAutoWithdrawTransaction(c *gin.Context) {

	id := c.Param("id")
	identifier, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		HandleError(c, err)
		return
	}
	username, err := h.accountingService.CheckCurrentUsername(c.MustGet("username"))
	if err != nil {
		HandleError(c, err)
		return
	}
	adminId, err := h.accountingService.CheckCurrentAdminId(c.MustGet("adminId"))
	if err != nil {
		HandleError(c, err)
		return
	}

	var req model.BankConfirmTransferWithdrawRequest
	if err := c.ShouldBind(&req); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(req); err != nil {
		HandleError(c, err)
		return
	}
	req.ConfirmedAt = time.Now()
	req.ConfirmedByUserId = *adminId
	req.ConfirmedByUsername = *username

	actionErr := h.bankingService.ContinueAutoWithdrawTransaction(identifier)
	if actionErr != nil {
		HandleError(c, actionErr)
		return
	}
	c.JSON(201, model.Success{Message: "Confirm success"})
}

// @Summary GetFinishedTransactions
// @Description ดึงข้อมูลลิสการถอน ที่รออนุมัติ
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.FinishedTransactionListRequest true "query"
// @Success 200 {object} model.SuccessWithPagination
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/finishedlist [get]
func (h bankingController) getFinishedTransactions(c *gin.Context) {

	var query model.FinishedTransactionListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetFinishedTransactions(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary RemoveFinishedTransaction
// @Description ลบข้อมูลการฝากถอน ที่เสร็จสิ้นไปแล้ว
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/remove/{id} [post]
func (h bankingController) removeFinishedTransaction(c *gin.Context) {

	username, err := h.accountingService.CheckCurrentUsername(c.MustGet("username"))
	if err != nil {
		HandleError(c, err)
		return
	}
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

	var data model.BankTransactionRemoveBody
	data.Status = "removed"
	data.RemovedAt = time.Now()
	data.RemovedByUserId = *adminId
	data.RemovedByUsername = *username

	actionErr := h.bankingService.RemoveFinishedTransaction(identifier, data)
	if actionErr != nil {
		HandleError(c, actionErr)
		return
	}
	c.JSON(201, model.Success{Message: "Remove success"})
}

// @Summary GetRemovedTransactions
// @Description ดึงข้อมูลลิสการฝากถอนที่ถูกลบ
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.RemovedTransactionListRequest true "query"
// @Success 200 {object} model.SuccessWithPagination
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/removedlist [get]
func (h bankingController) getRemovedTransactions(c *gin.Context) {

	var query model.RemovedTransactionListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetRemovedTransactions(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary GetMemberByCode
// @Description ดึงข้อมูลสมาชิกด้วยโค้ด
// @Tags Banking - Member Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param code path string true "memberCode"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/member/info/{code} [get]
func (h bankingController) getMemberByCode(c *gin.Context) {

	memberCode := c.Param("code")

	data, err := h.bankingService.GetMemberByCode(memberCode)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary GetMembers
// @Description ดึงข้อมูลลิสมาชิก
// @Tags Banking - Bank Account Statements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.MemberListRequest true "MemberListRequest"
// @Success 200 {object} model.SuccessWithPagination
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/member/list [get]
func (h bankingController) getMembers(c *gin.Context) {

	var query model.MemberListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetMembers(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary GetPossibleStatementOwners
// @Description ดึงข้อมูลลิสสมาชิก ที่มีข้อมูลใกล้เคียงกับรายการสเตทเม้นที่รอดำเนินการ
// @Tags Banking - Bank Account Statements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.MemberPossibleListRequest true "MemberPossibleListRequest"
// @Success 200 {object} model.SuccessWithPagination
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/member/possibleowner [get]
func (h bankingController) getPossibleStatementOwners(c *gin.Context) {

	var query model.MemberPossibleListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetPossibleStatementOwners(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary GetMemberTransactionSummary
// @Description ดึงข้อมูลสรุปการฝากถอนของสมาชิก
// @Tags Banking - Member Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.MemberTransactionListRequest true "query"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/member/transactionsummary [get]
func (h bankingController) getMemberTransactionSummary(c *gin.Context) {

	var query model.MemberTransactionListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetMemberTransactionSummary(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary GetMemberTransactions
// @Description ดึงข้อมูลลิสการฝากถอนของสมาชิก
// @Tags Banking - Member Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.MemberTransactionListRequest true "query"
// @Success 200 {object} model.SuccessWithPagination
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/member/transactions [get]
func (h bankingController) getMemberTransactions(c *gin.Context) {

	var query model.MemberTransactionListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetMemberTransactions(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary GetMemberStatements ดึงข้อมูลรายการทางการเงินทั้งหมดของสมาชิก
// @Description ดึงข้อมูลรายการทางการเงินทั้งหมดของสมาชิก
// @Tags Banking - Member Statement
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.MemberTransactionListRequest true "query"
// @Success 200 {object} model.SuccessWithPagination
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/member/statements [get]
func (h bankingController) getMemberStatements(c *gin.Context) {

	var query model.MemberStatementListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetMemberStatements(query)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary GetMemberStatementById
// @Description ดึงข้อมูลรายการทางการเงินทั้งหมดของสมาชิกด้วย id
// @Tags Banking - Member Statement
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/member/statements/detail/{id} [get]
func (h bankingController) getMemberStatementById(c *gin.Context) {

	var req model.GetByIdRequest
	if err := c.ShouldBindUri(&req); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetMemberStatementById(req)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary ProcessMemberDepositCredit
// @Description สร้างข้อมูลรายการทางการเงิน
// @Tags Banking - Member Statement
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.MemberStatementCreateRequest true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/member/statements1 [post]
func (h bankingController) processMemberDepositCredit(c *gin.Context) {

	var req model.MemberStatementCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(req); err != nil {
		HandleError(c, err)
		return
	}

	if err := h.bankingService.ProcessMemberDepositCredit(req.UserId, req.Amount); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(201, model.Success{Message: "Created success"})
}

// @Summary ProcessMemberWithdrawCredit
// @Description สร้างข้อมูลรายการทางการเงิน
// @Tags Banking - Member Statement
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.MemberStatementCreateRequest true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/member/statements2 [post]
func (h bankingController) processMemberWithdrawCredit(c *gin.Context) {

	var req model.MemberStatementCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(req); err != nil {
		HandleError(c, err)
		return
	}

	if err := h.bankingService.ProcessMemberWithdrawCredit(req.UserId, req.Amount); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(201, model.Success{Message: "Created success"})
}

// @Summary ProcessMemberBonusCredit
// @Description สร้างข้อมูลรายการทางการเงิน
// @Tags Banking - Member Statement
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.MemberStatementCreateRequest true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/member/statements3 [post]
func (h bankingController) processMemberBonusCredit(c *gin.Context) {

	var req model.MemberStatementCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(req); err != nil {
		HandleError(c, err)
		return
	}

	if err := h.bankingService.ProcessMemberBonusCredit(req.UserId, req.Amount); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(201, model.Success{Message: "Created success"})
}

// @Summary ProcessMemberGetbackCredit
// @Description สร้างข้อมูลรายการทางการเงิน
// @Tags Banking - Member Statement
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.MemberStatementCreateRequest true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/member/statements4 [post]
func (h bankingController) processMemberGetbackCredit(c *gin.Context) {

	var req model.MemberStatementCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(req); err != nil {
		HandleError(c, err)
		return
	}

	if err := h.bankingService.ProcessMemberGetbackCredit(req.UserId, req.Amount); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(201, model.Success{Message: "Created success"})
}
