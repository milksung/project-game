package repository

import (
	"cybergame-api/model"

	"gorm.io/gorm"
)

func NewRecommendRepository(db *gorm.DB) RecommendRepository {
	return &repo{db}
}

type RecommendRepository interface {
	GetRecommendList(query model.RecommendQuery) ([]model.RecommendList, int64, error)
	CreateRecommend(recommend model.CreateRecommend) error
	UpdateRecommend(id int64, body model.CreateRecommend) error
	DeleteRecommend(id int64) error
}

func (r repo) GetRecommendList(query model.RecommendQuery) ([]model.RecommendList, int64, error) {

	var recommends []model.RecommendList

	db := r.db.Table("Recommend_channels").Select("id, title, status, url, created_at")

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	if query.Filter != "" {
		db = db.Where("title LIKE ?", "%"+query.Filter+"%")
	}

	if err := db.
		Limit(query.Limit).
		Offset(query.Limit * query.Page).
		Find(&recommends).
		Order("id ASC").
		Error; err != nil {
		return nil, 0, err
	}

	var total int64

	queryTotal := r.db.Table("Recommend_channels")

	if query.Status != "" {
		queryTotal = queryTotal.Where("status = ?", query.Status)
	}

	if query.Filter != "" {
		queryTotal = queryTotal.Where("title LIKE ?", "%"+query.Filter+"%")
	}

	if err := queryTotal.Table("Recommend_channels").
		Select("id").
		Count(&total).
		Error; err != nil {
		return nil, 0, err
	}

	return recommends, total, nil
}

func (r repo) CreateRecommend(recommend model.CreateRecommend) error {

	if err := r.db.Table("Recommend_channels").
		Create(&recommend).
		Error; err != nil {
		return err
	}

	return nil
}

func (r repo) UpdateRecommend(id int64, body model.CreateRecommend) error {

	if err := r.db.Table("Recommend_channels").
		Where("id = ?", id).
		Updates(&body).
		Error; err != nil {
		return err
	}

	return nil
}

func (r repo) DeleteRecommend(id int64) error {

	if err := r.db.Table("Recommend_channels").
		Where("id = ?", id).
		Delete(&model.Recommend{}).
		Error; err != nil {
		return err
	}

	return nil
}
