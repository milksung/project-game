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

type linenotifyController struct {
	linenotifyService service.LineNotifyService
}

func newLineNotifyController(
	linenotifyService service.LineNotifyService,
) linenotifyController {
	return linenotifyController{linenotifyService}
}

// @Summary CreateLineNotify
// @Description ตั้งค่าแจ้งเตือนไลน์
// @Tags LineNotify
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.LinenotifyCreateBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /linenotify/create [post]
func LineNotifyController(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewLineNotifyRepository(db)
	service := service.NewLineNotifyService(repo)
	handler := newLineNotifyController(service)

	linenotifRoute := r.Group("/linenotify")
	linenotifRoute.POST("/create", middleware.Authorize, handler.createLineNotify)
	linenotifRoute.GET("/detail/:id", middleware.Authorize, handler.getLineNotifyById)
	linenotifRoute.PUT("/update/:id", middleware.Authorize, handler.updateLineNotify)

	//GameCyberNoitfy
	linenotifRoute.GET("/game/detail/:id", handler.getLineNotifyCyberGameById)
	linenotifRoute.POST("game/create", middleware.Authorize, handler.createLineNotifyCyberGame)
	linenotifRoute.PUT("/game/update/:id", middleware.Authorize, handler.updateLinenotifyCyberGame)
	linenotifRoute.PUT("/game/type/update/:id", middleware.Authorize, handler.updateLinenotifyTypeCyberGame)
	linenotifRoute.GET("/game/list", middleware.Authorize, handler.getLinenotifyTypeCyberGameList)
}

// @Summary Get LinenotifyTypeCyberGame List
// @Description Get LinenotifyTypeCyberGame List
// @Tags LineNotify
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param _ query model.CyberGameQuery true "Query LinenotifyTypeCyberGame"
// @Success 200 {object} model.SuccessWithList
// @Failure 400 {object} handler.ErrorResponse
// @Router /linenotify/game/list [get]
func (h linenotifyController) getLinenotifyTypeCyberGameList(c *gin.Context) {

	var query model.CyberGameQuery
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}

	list, total, err := h.linenotifyService.GetLinenotifyTypeCyberGameList(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{Message: "Success", List: list, Total: total})
}

func (h linenotifyController) createLineNotify(c *gin.Context) {

	var line model.LinenotifyCreateBody
	if err := c.ShouldBindJSON(&line); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(line); err != nil {
		HandleError(c, err)
		return
	}

	errline := h.linenotifyService.CreateLineNotify(line)
	if errline != nil {
		HandleError(c, errline)
		return
	}
	c.JSON(201, model.Success{Message: "Created success"})

}

// @Summary GetLineNotifyById
// @Description ดึงข้อมูลการcแจ้งเตือนไลน์ ด้วย id
// @Tags LineNotify
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /linenotify/detail/{id} [get]
func (h linenotifyController) getLineNotifyById(c *gin.Context) {

	var line model.LinenotifyParam

	if err := c.ShouldBindUri(&line); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.linenotifyService.GetLineNotifyById(line)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary UpdateNotify
// @Description แก้ไข แจ้งเตือนไลน์
// @Tags LineNotify
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param body body model.LinenotifyUpdateBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /linenotify/update/{id} [put]
func (h linenotifyController) updateLineNotify(c *gin.Context) {

	var body model.LinenotifyUpdateBody
	if err := c.ShouldBindJSON(&body); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(body); err != nil {
		HandleError(c, err)
		return
	}

	id := c.Param("id")
	toInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	if err := h.linenotifyService.UpdateLineNotify(int64(toInt), body); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Updated success"})
}

// @Summary CreateLineNotifyCyberGame
// @Description ตั้งค่าแจ้งเตือนไลน์
// @Tags LineNotify
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.LineNoifyCyberGameBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /linenotify/game/create [post]
func (h linenotifyController) createLineNotifyCyberGame(c *gin.Context) {
	var bot model.LineNoifyCyberGameBody
	if err := c.ShouldBindJSON(&bot); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(bot); err != nil {
		HandleError(c, err)
		return
	}

	errline := h.linenotifyService.CreateLineNoifyCyberGame(bot)
	if errline != nil {
		HandleError(c, errline)
		return
	}
	c.JSON(201, model.Success{Message: "Created success"})

}

// @Summary UpdateNotifyCyberGame
// @Description แก้ไขการเชื่อมต่อไลน์
// @Tags LineNotify
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param body body model.UpdateStatusCyberGame true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /linenotify/game/update/{id} [put]
func (h linenotifyController) updateLinenotifyCyberGame(c *gin.Context) {

	var body model.UpdateStatusCyberGame
	if err := c.ShouldBindJSON(&body); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(body); err != nil {
		HandleError(c, err)
		return
	}

	id := c.Param("id")
	toInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	if err := h.linenotifyService.UpdateLinenotifyCyberGame(int64(toInt), body); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Updated success"})
}

// @Summary GetLineNotifyCyberGameById
// @Description ดึงข้อมูลการcแจ้งเตือนไลน์ ด้วย id
// @Tags LineNotify
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /linenotify/game/detail/{id} [get]
func (h linenotifyController) getLineNotifyCyberGameById(c *gin.Context) {

	var line model.LineNoifyCyberGameParam

	if err := c.ShouldBindUri(&line); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.linenotifyService.GetLineNoifyCyberGameById(line)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary UpdateNotifyCyberTypeGame
// @Description แก้ไขสถานะประเภทเกม
// @Tags LineNotify
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param body body model.UpdateStatusTypeCyberGame true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /linenotify/game/type/update/{id} [put]
func (h linenotifyController) updateLinenotifyTypeCyberGame(c *gin.Context) {

	var body model.UpdateStatusTypeCyberGame
	if err := c.ShouldBindJSON(&body); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(body); err != nil {
		HandleError(c, err)
		return
	}

	id := c.Param("id")
	toInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	if err := h.linenotifyService.UpdateLinenotifyTypeCyberGame(int64(toInt), body); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Updated success"})
}
