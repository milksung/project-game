package repository

import (
	"cybergame-api/model"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func NewSettingWebRepository(db *gorm.DB) SettingWebRepository {
	return &repo{db}
}

type SettingWebRepository interface {
	GetSettingWeb(req model.SettingwebListRequest) (*model.SuccessWithPagination, error)
	GetSettingWebById(id int64) (*model.Settingweb, error)
	CreateSettingWeb(data model.SettingwebCreateBody) error
	UpdateSettingWeb(id int64, data model.SettingwebUpdateBody) error
}

func (r repo) GetSettingWebById(id int64) (*model.Settingweb, error) {
	var settingweb model.Settingweb

	if err := r.db.Table("setting_web").
		Select("id, logo, backgrond_color, user_auto, otp_register, tran_withdraw,register,deposit_first,deposit_next,withdraw,line,url,opt").
		Where("id = ?", id).
		First(&settingweb).
		Error; err != nil {
		return nil, err
	}
	return &settingweb, nil
}

func (r repo) GetSettingWeb(req model.SettingwebListRequest) (*model.SuccessWithPagination, error) {

	var list []model.SettingwebResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("setting_web")
	count = count.Select("id")
	if req.Search != "" {
		count = count.Where("id = ?", req.Search)
	}
	if err = count.
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	if total > 0 {
		// SELECT //
		query := r.db.Table("setting_web")
		query = query.Select("id, logo, backgrond_color, user_auto, otp_register, tran_withdraw,register,deposit_first,deposit_next,withdraw,line,url,op")
		if req.Search != "" {
			query = query.Where("id = ?", req.Search)
		}

		// Sort by ANY //
		req.SortCol = strings.TrimSpace(req.SortCol)
		if req.SortCol != "" {
			if strings.ToLower(strings.TrimSpace(req.SortAsc)) == "desc" {
				req.SortAsc = "DESC"
			} else {
				req.SortAsc = "ASC"
			}
			query = query.Order(req.SortCol + " " + req.SortAsc)
		}
		if err = query.
			Limit(req.Limit).
			Offset(req.Page * req.Limit).
			Scan(&list).
			Error; err != nil {
			return nil, err
		}
	}

	// End count total records for pagination purposes (without limit and offset) //
	var result model.SuccessWithPagination
	if list == nil {
		list = []model.SettingwebResponse{}
	}
	result.List = list
	result.Total = total
	return &result, nil
}
func (r repo) CreateSettingWeb(data model.SettingwebCreateBody) error {
	if err := r.db.Table("setting_web").Create(&data).Error; err != nil {
		fmt.Println(data)
		return err
	}

	return nil
}

func (r repo) UpdateSettingWeb(id int64, data model.SettingwebUpdateBody) error {
	if err := r.db.Table("setting_web").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}
