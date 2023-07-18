package service

import (
	"cybergame-api/helper"
	"cybergame-api/model"
	"cybergame-api/repository"
)

type RecommendService interface {
	GetRecommendList(query model.RecommendQuery) ([]model.RecommendList, int64, error)
	CreateRecommend(user model.CreateRecommend) error
	UpdateRecommend(id int64, body model.CreateRecommend) error
	DeleteRecommend(id int64) error
}

type recommendService struct {
	repo repository.RecommendRepository
}

func NewRecommendService(
	repo repository.RecommendRepository,
) RecommendService {
	return &recommendService{repo}
}

func (s *recommendService) GetRecommendList(query model.RecommendQuery) ([]model.RecommendList, int64, error) {

	if err := helper.Pagination(&query.Page, &query.Limit); err != nil {
		return nil, 0, err
	}

	return s.repo.GetRecommendList(query)
}

func (s *recommendService) CreateRecommend(body model.CreateRecommend) error {

	if err := s.repo.CreateRecommend(body); err != nil {
		return err
	}

	return nil
}

func (s *recommendService) UpdateRecommend(id int64, body model.CreateRecommend) error {

	if err := s.repo.UpdateRecommend(id, body); err != nil {
		return err
	}

	return nil
}

func (s *recommendService) DeleteRecommend(id int64) error {

	if err := s.repo.DeleteRecommend(id); err != nil {
		return err
	}

	return nil
}
