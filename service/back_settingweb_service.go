package service

import (
	"cybergame-api/helper"
	"cybergame-api/model"
	"cybergame-api/repository"
)

// var (
// 	storageClient *storage.Client
// )

func NewSettingWebService(
	repo repository.SettingWebRepository,
) SettingWebService {
	return &settingwebService{repo}
}

type SettingWebService interface {
	CreateSettingWeb(data model.SettingwebCreateBody) error
	GetSettingWeb(data model.SettingwebListRequest) (*model.SuccessWithPagination, error)
	GetSettingWebById(data model.SettingwebParam) (*model.Settingweb, error)
}

type settingwebService struct {
	repo repository.SettingWebRepository
}

// CreateSettingWeb implements SettingWebService
func (s *settingwebService) CreateSettingWeb(data model.SettingwebCreateBody) error {

	var web model.Settingweb
	web.Logo = data.Logo
	web.BackgrondColor = data.UserAuto
	web.UserAuto = data.UserAuto
	web.OtpRegister = data.OtpRegister
	web.TranWithdraw = data.TranWithdraw
	web.Register = data.Register
	web.DepositFirst = data.DepositFirst
	web.DepositNext = data.DepositNext
	web.Withdraw = data.Withdraw
	web.Line = data.Line
	web.Url = data.Url
	web.Opt = data.Opt

	if err := s.repo.CreateSettingWeb(data); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *settingwebService) GetSettingWeb(params model.SettingwebListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&params.Page, &params.Limit); err != nil {
		return nil, badRequest(err.Error())
	}

	records, err := s.repo.GetSettingWeb(params)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *settingwebService) GetSettingWebById(data model.SettingwebParam) (*model.Settingweb, error) {

	setting, err := s.repo.GetSettingWebById(data.Id)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, notFound("Setting NotFound")
		}
		if err.Error() == "Account not found" {
			return nil, notFound("Setting NotFound")
		}
		return nil, internalServerError(err.Error())
	}
	return setting, nil
}
