package handler

import (
	"cybergame-api/model"
	"cybergame-api/service"

	"cybergame-api/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type authController struct {
	adminService service.AdminService
}

func newAuthController(
	adminService service.AdminService,
) authController {
	return authController{adminService}
}

func AuthController(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewAdminRepository(db)
	perRepo := repository.NewPermissionRepository(db)
	groupRepo := repository.NewGroupRepository(db)
	service := service.NewAdminService(repo, perRepo, groupRepo)
	handler := newAuthController(service)

	r.POST("/login", handler.login)
}

// @Summary Login
// @Description Login
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param login body model.LoginAdmin true "Login"
// @Success 201 {object} model.LoginResponse
// @Failure 400 {object} handler.ErrorResponse
// @Failure 401 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Router /login [post]
func (h authController) login(c *gin.Context) {

	var body model.LoginAdmin

	if err := c.ShouldBindJSON(&body); err != nil {
		HandleError(c, err)
		return
	}

	// if err := validator.New().Struct(body); err != nil {
	// 	HandleError(c, err)
	// 	return
	// }

	token, err := h.adminService.Login(body)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.LoginResponse{Token: *token})
}
