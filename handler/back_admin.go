package handler

import (
	"cybergame-api/middleware"
	"cybergame-api/model"
	"cybergame-api/service"
	"strconv"

	"cybergame-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type adminController struct {
	adminService service.AdminService
}

func newAdminController(
	adminService service.AdminService,
) adminController {
	return adminController{adminService}
}

func AdminController(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewAdminRepository(db)
	perRepo := repository.NewPermissionRepository(db)
	groupRepo := repository.NewGroupRepository(db)
	service := service.NewAdminService(repo, perRepo, groupRepo)
	handler := newAdminController(service)

	// role := middleware.Role(db)

	r = r.Group("/admins")
	// r.GET("/detail/:id", middleware.Authorize, role.CheckAdmin("admin"), handler.GetAdmin)
	r.GET("/detail/:id", middleware.Authorize, handler.GetAdmin)
	r.GET("/list", middleware.Authorize, handler.getAdminList)
	r.POST("/create", middleware.Authorize, handler.create)
	r.PUT("/update/:id", middleware.Authorize, handler.updateAdmin)
	r.PUT("/password/:id", middleware.Authorize, handler.resetPassword)

	r.GET("/group", middleware.Authorize, handler.groupList)
	r.GET("/group/:id", middleware.Authorize, handler.getGroup)
	// r.POST("/group", middleware.Authorize, handler.createGroup)
	r.PUT("/group/:id", middleware.Authorize, handler.updateGroup)
	r.DELETE("/group/:id", middleware.Authorize, handler.deleteGroup)
	r.DELETE("/permission", middleware.Authorize, handler.deletePermission)
	r.DELETE("/:id", middleware.Authorize, handler.deleteAdmin)
}

// @Summary Group List
// @Description Group List
// @Tags Admins
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param _ query model.AdminGroupQuery false "Queries"
// @Success 200 {object} model.SuccessWithList
// @Failure 400 {object} handler.ErrorResponse
// @Router /admins/group [get]
func (h adminController) groupList(c *gin.Context) {

	query := model.AdminGroupQuery{}
	if err := c.ShouldBindQuery(&query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.adminService.GetGroupList(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, data)
}

// @Summary Get Group
// @Description Get Group
// @Tags Admins
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "Group ID"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /admins/group/{id} [get]
func (h adminController) getGroup(c *gin.Context) {

	id := c.Param("id")
	toInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.adminService.GetGroup(toInt)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "Success", Data: data})
}

// @Summary Get Admin
// @Description Get Admin
// @Tags Admins
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "Admin ID"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /admins/detail/{id} [get]
func (h adminController) GetAdmin(c *gin.Context) {

	id := c.Param("id")
	toInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.adminService.GetAdmin(int64(toInt))
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "Success", Data: data})
}

// @Summary Get Admin List
// @Description Get Admin List
// @Tags Admins
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param _ query model.AdminListQuery false "Queries"
// @Success 200 {object} model.SuccessWithList
// @Failure 400 {object} handler.ErrorResponse
// @Router /admins/list [get]
func (h adminController) getAdminList(c *gin.Context) {

	query := model.AdminListQuery{}
	if err := c.ShouldBindQuery(&query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.adminService.GetAdminList(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, data)
}

// @Summary Create Admin
// @Description Create Admin
// @Tags Admins
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param Body body model.CreateAdmin true "Create Admin"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /admins/create [post]
func (h adminController) create(c *gin.Context) {

	data := &model.CreateAdmin{}
	if err := c.ShouldBindJSON(data); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(data); err != nil {
		HandleError(c, err)
		return
	}

	err, perErrs := h.adminService.Create(data)
	if err != nil {
		HandleError(c, err)
		return
	}

	if perErrs != nil {
		HandleError(c, perErrs)
		return
	}

	c.JSON(201, model.Success{Message: "Register success"})
}

// // @Summary Create Group
// // @Description Create Group
// // @Tags Admins
// // @Security BearerAuth
// // @Accept  json
// // @Produce  json
// // @Param Body body model.AdminCreateGroup true "Create Group"
// // @Success 201 {object} model.Success
// // @Failure 400 {object} handler.ErrorResponse
// // @Router /admins/group [post]
// func (h adminController) createGroup(c *gin.Context) {

// 	data := &model.AdminCreateGroup{}
// 	if err := c.ShouldBindJSON(data); err != nil {
// 		HandleError(c, err)
// 		return
// 	}

// 	if err := validator.New().Struct(data); err != nil {
// 		HandleError(c, err)
// 		return
// 	}

// 	err := h.adminService.CreateGroup(data)
// 	if err != nil {
// 		HandleError(c, err)
// 		return
// 	}

// 	c.JSON(201, model.Success{Message: "Created success"})
// }

// @Summary Update Group
// @Description Update Group
// @Tags Admins
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "Group ID"
// @Param Body body model.AdminUpdateGroup true "Update Group"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /admins/group/{id} [put]
func (h adminController) updateGroup(c *gin.Context) {

	data := &model.AdminUpdateGroup{}
	if err := c.ShouldBindJSON(data); err != nil {
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

	err = h.adminService.UpdateGroup(int64(toInt), data)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Updated success"})
}

// @Summary Update Admin
// @Description Update Admin
// @Tags Admins
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "Admin ID"
// @Param Body body model.AdminBody true "Update Admin"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /admins/update/{id} [put]
func (h adminController) updateAdmin(c *gin.Context) {

	data := model.AdminBody{}
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

	err, perErrs := h.adminService.UpdateAdmin(int64(toInt), data)
	if err != nil {
		HandleError(c, err)
		return
	}

	if perErrs != nil {
		HandleError(c, perErrs)
		return
	}

	c.JSON(201, model.Success{Message: "Updated success"})
}

// @Summary Update Admin Password
// @Description Update Admin Password
// @Tags Admins
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "Admin ID"
// @Param Body body model.AdminUpdatePassword true "Update Admin Password"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /admins/password/{id} [put]
func (h adminController) resetPassword(c *gin.Context) {

	id := c.Param("id")
	toInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	data := model.AdminUpdatePassword{}
	if err := c.ShouldBindJSON(&data); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(data); err != nil {
		HandleError(c, err)
		return
	}

	if err := h.adminService.ResetPassword(int64(toInt), data); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Reset password success"})
}

// @Summary Delete Group
// @Description Delete Group
// @Tags Admins
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "Group ID"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /admins/group/{id} [delete]
func (h adminController) deleteGroup(c *gin.Context) {

	id := c.Param("id")
	toInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	var param model.DeleteGroup
	param.Id = int64(toInt)

	if err := h.adminService.DeleteGroup(param.Id); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Deleted success"})
}

// @Summary Delete Permission
// @Description Delete Permission
// @Tags Admins
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param delete body model.DeletePermission true "Delete Permission"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /admins/permission [delete]
func (h adminController) deletePermission(c *gin.Context) {

	var body model.DeletePermission
	if err := c.ShouldBindJSON(&body); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(body); err != nil {
		HandleError(c, err)
		return
	}

	if err := h.adminService.DeletePermission(body); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Deleted success"})
}

// @Summary Delete Admin
// @Description Delete Admin
// @Tags Admins
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "Admin ID"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /admins/{id} [delete]
func (h adminController) deleteAdmin(c *gin.Context) {

	id := c.Param("id")
	toInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	if err := h.adminService.DeleteAdmin(int64(toInt)); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Deleted success"})
}
