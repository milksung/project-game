package service

import (
	"cybergame-api/helper"
	"cybergame-api/model"
	"cybergame-api/repository"
)

func NewLineNotifyService(
	repo repository.LineNotifyRepository,
) LineNotifyService {
	return &lineNotifyService{repo}
}

type LineNotifyService interface {
	//line notify group
	CreateLineNotify(data model.LinenotifyCreateBody) error
	GetLineNotify(data model.LinenotifyListRequest) (*model.SuccessWithPagination, error)
	GetLineNotifyById(data model.LinenotifyParam) (*model.Linenotify, error)
	UpdateLineNotify(id int64, data model.LinenotifyUpdateBody) error

	//line notify game
	CreateLineNoifyCyberGame(data model.LineNoifyCyberGameBody) error
	UpdateLinenotifyCyberGame(id int64, body model.UpdateStatusCyberGame) error
	GetLineNoifyCyberGameById(model.LineNoifyCyberGameParam) (*model.LineNoifyCyberGame, error)
	UpdateLinenotifyTypeCyberGame(id int64, body model.UpdateStatusTypeCyberGame) error
	GetLinenotifyTypeCyberGameList(query model.CyberGameQuery) ([]model.CyberGameList, int64, error)
}

const Admin_NotFound = "Admin not found"

type lineNotifyService struct {
	repo repository.LineNotifyRepository
}

func (s *lineNotifyService) GetLinenotifyTypeCyberGameList(query model.CyberGameQuery) ([]model.CyberGameList, int64, error) {

	if err := helper.Pagination(&query.Page, &query.Limit); err != nil {
		return nil, 0, err
	}

	return s.repo.GetLinenotifyCyberGameList(query)
}

func (s *lineNotifyService) CreateLineNotify(data model.LinenotifyCreateBody) error {

	var admin model.Admin

	checkAdmin, err := s.repo.CheckAdminId(admin.Id)
	if err != nil {
		return internalServerError(err.Error())
	}

	if !checkAdmin {
		return notFound(Admin_NotFound)
	}

	if err := s.repo.CreateLineNotify(data); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *lineNotifyService) GetLineNotify(params model.LinenotifyListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&params.Page, &params.Limit); err != nil {
		return nil, badRequest(err.Error())
	}

	records, err := s.repo.GetLineNotify(params)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *lineNotifyService) GetLineNotifyById(data model.LinenotifyParam) (*model.Linenotify, error) {

	line, err := s.repo.GetLineNotifyById(data.Id)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, notFound("record NotFound")
		}
		if err.Error() == "record not found" {
			return nil, notFound("record NotFound")
		}
		return nil, internalServerError(err.Error())
	}
	return line, nil
}

func (s *lineNotifyService) UpdateLineNotify(id int64, data model.LinenotifyUpdateBody) error {
	if err := s.repo.UpdateLineNotify(id, data); err != nil {
		return internalServerError(err.Error())
	}

	return nil
}

func (s *lineNotifyService) CreateLineNoifyCyberGame(data model.LineNoifyCyberGameBody) error {

	var list []model.LineNoifyCyberGameBody

	typeId := []int{}

	// เพิ่มค่าลงใน slice ด้วย loop
	for i := 1; i < 12; i++ {
		typeId = append(typeId, i)
	}

	for _, id := range typeId {

		list = append(list, model.LineNoifyCyberGameBody{
			Token:        data.Token,
			TypenotifyId: int64(id),
			AdminId:      data.AdminId,
			Type:         []int64{1, 2, 3, 4, 5, 6},
		})

		if err := s.repo.CreateLinenotifyCyberGame(list); err != nil {
			return internalServerError(err.Error())
		}
	}

	return nil

	// if err := s.repo.CreateLinenotifyCyberGame(data); err != nil {
	// 	return internalServerError(err.Error())
	// }

	// return nil
}

func (s *lineNotifyService) UpdateLinenotifyCyberGame(id int64, body model.UpdateStatusCyberGame) error {

	if err := s.repo.UpdateLineNotifyCyberGame(id, body); err != nil {
		return err
	}

	return nil
}

func (s *lineNotifyService) GetLineNoifyCyberGameById(data model.LineNoifyCyberGameParam) (*model.LineNoifyCyberGame, error) {

	line, err := s.repo.GetLinenotifyCyberGameById(data.Id)

	if err != nil {
		return nil, notFound("record NotFound")
	}

	if err != nil {
		if err.Error() == "record not found" {
			return nil, notFound("record NotFound")
		}
		if err.Error() == "record not found" {
			return nil, notFound("record NotFound")
		}
		return nil, internalServerError(err.Error())
	}
	return line, nil
}

func (s *lineNotifyService) UpdateLinenotifyTypeCyberGame(id int64, body model.UpdateStatusTypeCyberGame) error {

	if err := s.repo.UpdateLineNotifyTypeCyberGame(id, body); err != nil {
		return err
	}

	return nil
}
