package handler

import (
	"cybergame-api/middleware"
	"cybergame-api/model"
	"cybergame-api/service"

	"cybergame-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type scammerController struct {
	scammerService service.ScammerService
}

func newScammerController(
	scammerService service.ScammerService,
) scammerController {
	return scammerController{scammerService}
}

func ScammerController(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewScammerRepository(db)
	service := service.NewScammerService(repo)
	handler := newScammerController(service)

	r = r.Group("/scammers")
	r.GET("/list", middleware.Authorize, handler.getScammerList)
	r.POST("/create", middleware.Authorize, handler.CreateScammer)

}

// @Summary Get Scammer List
// @Description Get Scammer List
// @Tags Scammers
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param _ query model.ScammerQuery true "Query Scammer"
// @Success 200 {object} model.SuccessWithList
// @Failure 400 {object} handler.ErrorResponse
// @Router /scammers/list [get]
func (h scammerController) getScammerList(c *gin.Context) {

	var query model.ScammerQuery
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}

	list, total, err := h.scammerService.GetScammerList(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{Message: "Success", List: list, Total: total})
}

// @Summary Create Scammer
// @Description Create Scammer
// @Tags Scammers
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param Body body model.CreateScammer true "Create Scammer"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /scammers/create [post]
func (h scammerController) CreateScammer(c *gin.Context) {

	var body model.CreateScammer
	if err := c.ShouldBindJSON(&body); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(body); err != nil {
		HandleError(c, err)
		return
	}

	err := h.scammerService.Create(body)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Created success"})
}
