package handler

import (
	"cybergame-api/middleware"
	"cybergame-api/model"
	"cybergame-api/service"

	"cybergame-api/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type menuController struct {
	menuService service.MenuService
}

func newMenuController(
	menuService service.MenuService,
) menuController {
	return menuController{menuService}
}

func MenuController(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewPermissionRepository(db)
	service := service.NewMenuService(repo)
	handler := newMenuController(service)

	// r = r.Group("/menu")
	r.GET("/menu", middleware.Authorize, handler.GetMenu)
}

// @Summary Get Menu
// @Description Get Menu
// @Tags Menu
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /menu [get]
func (h menuController) GetMenu(c *gin.Context) {

	adminId := c.MustGet("adminId").(float64)
	toInt := int64(adminId)

	list, err := h.menuService.GetMenu(toInt)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithList{Message: "Success", List: list})
}
