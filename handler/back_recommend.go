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

type recommendController struct {
	recommendService service.RecommendService
}

func newRecommendController(
	recommendService service.RecommendService,
) recommendController {
	return recommendController{recommendService}
}

func RecommendController(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewRecommendRepository(db)
	service := service.NewRecommendService(repo)
	handler := newRecommendController(service)

	r = r.Group("/recommends")
	r.GET("/list", middleware.Authorize, handler.getRecommendList)
	r.POST("/create", middleware.Authorize, handler.createRecommend)
	r.PUT("/update/:id", middleware.Authorize, handler.updateRecommend)
	r.DELETE("/:id", middleware.Authorize, handler.deleteRecommend)

}

// @Summary Get Recommend List
// @Description Get Recommend List
// @Tags Recommends
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param _ query model.RecommendQuery true "Query Recommend"
// @Success 200 {object} model.SuccessWithList
// @Failure 400 {object} handler.ErrorResponse
// @Router /recommends/list [get]
func (h recommendController) getRecommendList(c *gin.Context) {

	var query model.RecommendQuery
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}

	list, total, err := h.recommendService.GetRecommendList(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{Message: "Success", List: list, Total: total})
}

// @Summary Create Recommend
// @Description Create Recommend
// @Tags Recommends
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param Body body model.CreateRecommend true "Create Recommend"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /recommends/create [post]
func (h recommendController) createRecommend(c *gin.Context) {

	var body model.CreateRecommend
	if err := c.ShouldBindJSON(&body); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(body); err != nil {
		HandleError(c, err)
		return
	}

	if err := h.recommendService.CreateRecommend(body); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Created success"})
}

// @Summary Update Recommend
// @Description Update Recommend
// @Tags Recommends
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "Recommend ID"
// @Param Body body model.CreateRecommend true "Update Recommend"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /recommends/update/{id} [put]
func (h recommendController) updateRecommend(c *gin.Context) {

	var body model.CreateRecommend
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

	if err := h.recommendService.UpdateRecommend(int64(toInt), body); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Updated success"})
}

// @Summary Delete Recommend
// @Description Delete Recommend
// @Tags Recommends
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "Recommend ID"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /recommends/{id} [delete]
func (h recommendController) deleteRecommend(c *gin.Context) {

	id := c.Param("id")
	toInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	if err := h.recommendService.DeleteRecommend(int64(toInt)); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Deleted success"})
}
