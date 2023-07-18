package handler

import (
	"cybergame-api/model"
	"cybergame-api/service"
	"strconv"

	"cybergame-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type frontAuthController struct {
	frontUserService service.FrontUserService
}

func newFrontAuthController(
	frontUserService service.FrontUserService,
) frontAuthController {
	return frontAuthController{frontUserService}
}

func FrontAuthController(r *gin.RouterGroup, db *gorm.DB) {

	frontUserRepo := repository.NewFrontUserRepository(db)
	agentConnectRepo := repository.NewAgentConnectRepository(db)
	service := service.NewFrontUserService(frontUserRepo, agentConnectRepo)
	handler := newFrontAuthController(service)

	r.POST("/sendotp/register", handler.frontSendOTPRegister)
	r.POST("/sendotp/forget", handler.frontSendOTPForget)
	r.POST("/verifyotp/register", handler.verifyOTPRegister)
	r.POST("/verifyotp/forget", handler.verifyOTPForget)
	r.POST("/resetpassword/:userId", handler.resetPassword)
	r.POST("/login", handler.frontLogin)
	r.POST("/updateinfo/:id", handler.updateInfo)
}

// @Summary ส่ง OTP ไปยังเบอร์โทรศัพท์ สำหรับลงทะเบียน
// @Description Send OTP
// @Tags Front - Auth
// @Accept  json
// @Produce  json
// @Param body body model.UserSendOTP true "Send OTP"
// @Success 201 {object} model.SuccessWithData
// @Failure 400 {object} ErrorResponse
// @Router /v1/frontend/sendotp/register [post]
func (h frontAuthController) frontSendOTPRegister(c *gin.Context) {

	var body model.UserSendOTP

	if err := c.ShouldBindJSON(&body); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(body); err != nil {
		HandleError(c, err)
		return
	}

	if err := h.frontUserService.FrontUserSendOTPRegister(body); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Sent success"})
}

// @Summary ส่ง OTP ไปยังเบอร์โทรศัพท์ สำหรับลืมรหัสผ่าน
// @Description Send OTP
// @Tags Front - Auth
// @Accept  json
// @Produce  json
// @Param body body model.UserSendOTP true "Send OTP"
// @Success 201 {object} model.Success
// @Failure 400 {object} ErrorResponse
// @Router /v1/frontend/sendotp/forget [post]
func (h frontAuthController) frontSendOTPForget(c *gin.Context) {

	var body model.UserSendOTP

	if err := c.ShouldBindJSON(&body); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(body); err != nil {
		HandleError(c, err)
		return
	}

	if err := h.frontUserService.FrontUserSendOTPFotget(body); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Sent success"})
}

// @Summary ยืนยัน OTP สำหรับลงทะเบียน
// @Description Verify OTP Register
// @Tags Front - Auth
// @Accept  json
// @Produce  json
// @Param body body model.UserVerifyOTP true "Verify OTP"
// @Success 201 {object} model.SuccessWithData
// @Failure 400 {object} ErrorResponse
// @Router /v1/frontend/verifyotp/register [post]
func (h frontAuthController) verifyOTPRegister(c *gin.Context) {

	var body model.UserVerifyOTP

	if err := c.ShouldBindJSON(&body); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(body); err != nil {
		HandleError(c, err)
		return
	}

	userId, err := h.frontUserService.FrontUserVerifyOTPRegister(body)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.SuccessWithData{Message: "Verify success", Data: userId})
}

// @Summary ยืนยัน OTP สำหรับลืมรหัสผ่าน
// @Description Verify OTP Forget
// @Tags Front - Auth
// @Accept  json
// @Produce  json
// @Param body body model.UserVerifyOTP true "Verify OTP"
// @Success 201 {object} model.SuccessWithData
// @Failure 400 {object} ErrorResponse
// @Router /v1/frontend/verifyotp/forget [post]
func (h frontAuthController) verifyOTPForget(c *gin.Context) {

	var body model.UserVerifyOTP

	if err := c.ShouldBindJSON(&body); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(body); err != nil {
		HandleError(c, err)
		return
	}

	userId, err := h.frontUserService.FrontUserVerifyOTPForget(body)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.SuccessWithData{Data: userId})
}

// @Summary reset password สำหรับลืมรหัสผ่าน
// @Description Reset Password
// @Tags Front - Auth
// @Accept  json
// @Produce  json
// @Param userId path string true "User ID"
// @Param body body model.UserUpdatePassword true "Reset Password"
// @Success 201 {object} model.Success
// @Failure 400 {object} ErrorResponse
// @Router /v1/frontend/resetpassword/{userId} [post]
func (h frontAuthController) resetPassword(c *gin.Context) {

	id := c.Param("userId")
	toInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	var body model.UserUpdatePassword

	if err := c.ShouldBindJSON(&body); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(body); err != nil {
		HandleError(c, err)
		return
	}

	if err := h.frontUserService.FrontUserResetPassword(int64(toInt), body); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Change password success"})
}

// @Summary ล็อกอิน
// @Description Login
// @Tags Front - Auth
// @Accept  json
// @Produce  json
// @Param body body model.UserLogin true "Login"
// @Success 201 {object} model.SuccessWithToken
// @Failure 400 {object} ErrorResponse
// @Router /v1/frontend/login [post]
func (h frontAuthController) frontLogin(c *gin.Context) {

	var body model.UserLogin

	if err := c.ShouldBindJSON(&body); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(body); err != nil {
		HandleError(c, err)
		return
	}

	token, err := h.frontUserService.FrontUserLogin(body)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.SuccessWithToken{Token: *token})
}

// @Summary อัพเดทข้อมูลผู้ใช้งาน
// @Description Update User Info
// @Tags Front - Auth
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param body body model.FrontUserUpdate true "Update User Info"
// @Success 201 {object} model.Success
// @Failure 400 {object} ErrorResponse
// @Router /v1/frontend/updateinfo/{id} [post]
func (h frontAuthController) updateInfo(c *gin.Context) {

	data := model.FrontUserUpdate{}
	if err := c.ShouldBindJSON(&data); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(data); err != nil {
		HandleError(c, err)
		return
	}

	id := c.Param("id")
	toInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	err = h.frontUserService.FrontUserUpdateInfo(int64(toInt), data)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Updated success"})
}
