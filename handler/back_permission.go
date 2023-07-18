package handler

import (
	"cybergame-api/model"
	"cybergame-api/service"

	"cybergame-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type permissionController struct {
	permissionService service.PermissionService
}

func newPermissionController(
	permissionService service.PermissionService,
) permissionController {
	return permissionController{permissionService}
}

func PermissionController(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewPermissionRepository(db)
	service := service.NewPermissionService(repo)
	handler := newPermissionController(service)

	r = r.Group("/permissions")
	r.POST("/create", handler.create)
}

// @Summary Create Permission
// @Description Create Permission
// @Tags Permissions
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param register body model.CreatePermission true "Create Permission"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /permissions/create [post]
func (h permissionController) create(c *gin.Context) {

	data := &model.CreatePermission{}
	if err := c.ShouldBindJSON(data); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(data); err != nil {
		HandleError(c, err)
		return
	}

	err := h.permissionService.Create(data)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Created success"})
}
