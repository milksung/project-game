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

type userController struct {
	userService service.UserService
}

func newUserController(
	userService service.UserService,
) userController {
	return userController{userService}
}

func UserController(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewUserRepository(db)
	perRepo := repository.NewPermissionRepository(db)
	groupRepo := repository.NewGroupRepository(db)
	agentConnectRepo := repository.NewAgentConnectRepository(db)
	service := service.NewUserService(repo, perRepo, groupRepo, agentConnectRepo)
	handler := newUserController(service)

	// role := middleware.Role(db)

	r = r.Group("/users")
	r.GET("/loginlogs/:id", middleware.Authorize, handler.getLoginLogs)
	r.GET("/detail/:id", middleware.Authorize, handler.GetUser)
	r.GET("/list", middleware.Authorize, handler.getUserList)
	r.GET("/updatelogs", middleware.Authorize, handler.getUpdateLogs)
	r.POST("/create", middleware.Authorize, handler.create)
	r.PUT("/update/:id", middleware.Authorize, handler.updateUser)
	r.PUT("/password/:id", middleware.Authorize, handler.resetPassword)
	r.DELETE("/:id", middleware.Authorize, handler.deleteUser)
}

// @Summary แสดงลิสประวัติการเข้าสู่ระบบของ User
// @Description Login Logs
// @Tags Users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} model.SuccessWithList
// @Failure 400 {object} handler.ErrorResponse
// @Router /users/loginlogs/{id} [get]
func (h userController) getLoginLogs(c *gin.Context) {

	id := c.Param("id")
	toInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	list, err := h.userService.GetUserLoginLogs(int64(toInt))
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithList{Message: "Success", List: list})
}

// @Summary Get User
// @Description Get User
// @Tags Users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /users/detail/{id} [get]
func (h userController) GetUser(c *gin.Context) {

	id := c.Param("id")
	toInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.userService.GetUser(int64(toInt))
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "Success", Data: data})
}

// @Summary Get User List
// @Description Get User List
// @Tags Users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param _ query model.UserListQuery false "Queries"
// @Success 200 {object} model.SuccessWithList
// @Failure 400 {object} handler.ErrorResponse
// @Router /users/list [get]
func (h userController) getUserList(c *gin.Context) {

	query := model.UserListQuery{}
	if err := c.ShouldBindQuery(&query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.userService.GetUserList(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, data)
}

// @Summary Get User Update Logs
// @Description Get User Update Logs
// @Tags Users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param _ query model.UserUpdateQuery false "Queries"
// @Success 200 {object} model.SuccessWithList
// @Failure 400 {object} handler.ErrorResponse
// @Router /users/updatelogs [get]
func (h userController) getUpdateLogs(c *gin.Context) {

	query := model.UserUpdateQuery{}
	if err := c.ShouldBindQuery(&query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.userService.GetUpdateLogs(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, data)
}

// @Summary Create User
// @Description Create User
// @Tags Users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param body body model.CreateUser false "Create User"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /users/create [post]
func (h userController) create(c *gin.Context) {

	data := &model.CreateUser{}
	if err := c.ShouldBindJSON(data); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(data); err != nil {
		HandleError(c, err)
		return
	}

	err := h.userService.Create(data)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Register success"})
}

// @Summary Update User
// @Description Update User
// @Tags Users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param body body model.UpdateUser true "Update User"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /users/update/{id} [put]
func (h userController) updateUser(c *gin.Context) {

	adminName := c.MustGet("username").(string)

	data := model.UpdateUser{}
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

	err = h.userService.UpdateUser(int64(toInt), data, adminName)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Updated success"})
}

// @Summary Update User Password
// @Description Update User Password
// @Tags Users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param body body model.UserUpdatePassword true "Update User Password"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /users/password/{id} [put]
func (h userController) resetPassword(c *gin.Context) {

	id := c.Param("id")
	toInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	data := model.UserUpdatePassword{}
	if err := c.ShouldBindJSON(&data); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(data); err != nil {
		HandleError(c, err)
		return
	}

	if err := h.userService.ResetPassword(int64(toInt), data); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Reset password success"})
}

// @Summary Delete User
// @Description Delete User
// @Tags Users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /users/{id} [delete]
func (h userController) deleteUser(c *gin.Context) {

	id := c.Param("id")
	toInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	err = h.userService.DeleteUser(int64(toInt))
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Deleted success"})
}
