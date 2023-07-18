package service

import (
	"cybergame-api/helper"
	"cybergame-api/model"
	"cybergame-api/repository"
	"strings"
)

type ScammerService interface {
	GetScammerList(query model.ScammerQuery) ([]model.ScammertList, int64, error)
	Create(user model.CreateScammer) error
}

type scammerService struct {
	repo repository.ScammerRepository
}

func NewScammerService(
	repo repository.ScammerRepository,
) ScammerService {
	return &scammerService{repo}
}

func (s *scammerService) GetScammerList(query model.ScammerQuery) ([]model.ScammertList, int64, error) {

	if err := helper.Pagination(&query.Page, &query.Limit); err != nil {
		return nil, 0, err
	}

	return s.repo.GetScammerList(query)
}

func (s *scammerService) Create(body model.CreateScammer) error {

	var data model.Scammer

	splitFullname := strings.Split(*body.Fullname, " ")
	if len(splitFullname) > 1 {
		data.Firstname = &splitFullname[0]
		data.Lastname = &splitFullname[1]
	}

	data.Fullname = body.Fullname
	data.Bankname = body.Bankname
	data.BankAccount = body.BankAccount
	data.Phone = body.Phone
	data.Reason = body.Reason

	if err := s.repo.CreateScammer(data); err != nil {
		return err
	}

	return nil
}
