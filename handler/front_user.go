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

type frontUserController struct {
	frontUserService service.FrontUserService
}

func newFrontUserController(
	frontUserService service.FrontUserService,
) frontUserController {
	return frontUserController{frontUserService}
}

func FrontUserController(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewFrontUserRepository(db)
	agentConnectRepo := repository.NewAgentConnectRepository(db)
	service := service.NewFrontUserService(repo, agentConnectRepo)
	handler := newFrontUserController(service)

	// role := middleware.Role(db)

	r = r.Group("/users")
	r.GET("/detail/:id", middleware.UserAuthorize, handler.getFrontUser)
	r.PUT("/password/:id", middleware.UserAuthorize, handler.resetFrontPassword)
}

// @Summary Get User
// @Description Get User
// @Tags Front - Users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /v1/front/users/detail/{id} [get]
func (h frontUserController) getFrontUser(c *gin.Context) {

	id := c.Param("id")
	toInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.frontUserService.GetFrontUser(int64(toInt))
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "Success", Data: data})
}

// @Summary Update User Password
// @Description Update User Password
// @Tags Front - Users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param body body model.UserUpdatePassword true "Update User Password"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /v1/front/users/password/{id} [put]
func (h frontUserController) resetFrontPassword(c *gin.Context) {

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

	if err := h.frontUserService.FrontUserChangePassword(int64(toInt), data); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Reset password success"})
}
