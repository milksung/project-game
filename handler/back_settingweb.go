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

type settingwebController struct {
	settingebService service.SettingWebService
}

func newSettingwebController(
	settingebService service.SettingWebService,
) settingwebController {
	return settingwebController{settingebService}
}

// @Summary CreateSettingWeb
// @Description ตั้งค่าหน้าเว็บไซต์
// @Tags Settingweb
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.SettingwebCreateBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /settingweb/create [post]
func SettingwebController(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewSettingWebRepository(db)
	service := service.NewSettingWebService(repo)
	handler := newSettingwebController(service)

	settingWebRoute := r.Group("/settingweb")
	settingWebRoute.POST("/create", middleware.Authorize, handler.createsettingweb)
	settingWebRoute.GET("/detail/:id", middleware.Authorize, handler.getSettingWebById)

}
func (h settingwebController) createsettingweb(c *gin.Context) {

	var settingweb model.SettingwebCreateBody
	if err := c.ShouldBindJSON(&settingweb); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(settingweb); err != nil {
		HandleError(c, err)
		return
	}

	errsettingweb := h.settingebService.CreateSettingWeb(settingweb)
	if errsettingweb != nil {
		HandleError(c, errsettingweb)
		return
	}
	c.JSON(201, model.Success{Message: "Created success"})

}

// @Summary GetSettingWebById
// @Description ดึงข้อมูลการตั้งค่า ด้วย id
// @Tags Settingweb
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /settingweb/detail/{id} [get]
func (h settingwebController) getSettingWebById(c *gin.Context) {

	var settingweb model.SettingwebParam

	if err := c.ShouldBindUri(&settingweb); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.settingebService.GetSettingWebById(settingweb)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}
