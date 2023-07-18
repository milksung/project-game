package repository

import (
	"cybergame-api/model"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func NewLineNotifyRepository(db *gorm.DB) LineNotifyRepository {
	return &repo{db}
}

type LineNotifyRepository interface {
	//linenotify
	GetLineNotify(req model.LinenotifyListRequest) (*model.SuccessWithPagination, error)
	GetLineNotifyById(id int64) (*model.Linenotify, error)
	CreateLineNotify(data model.LinenotifyCreateBody) error
	UpdateLineNotify(id int64, data model.LinenotifyUpdateBody) error

	//linenotifygame
	GetLinenotifyGameById(id int64) (*model.LinenotifyGame, error)
	CreateLinenotifyCyberGame(data []model.LineNoifyCyberGameBody) error
	UpdateLineNotifyCyberGame(id int64, body model.UpdateStatusCyberGame) error
	GetLinenotifyCyberGameById(id int64) (*model.LineNoifyCyberGame, error)
	UpdateLineNotifyTypeCyberGame(id int64, body model.UpdateStatusTypeCyberGame) error
	GetLinenotifyCyberGameList(query model.CyberGameQuery) ([]model.CyberGameList, int64, error)
	CheckAdminId(adminid int64) (bool, error)
}

func (r repo) CheckAdminId(adminid int64) (bool, error) {

	var lineGame model.LineNoifyCyberGame

	if err := r.db.Table("Admins").
		Where("id = ?", adminid).
		First(&lineGame).
		Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r repo) GetLinenotifyCyberGameList(query model.CyberGameQuery) ([]model.CyberGameList, int64, error) {

	var cybergame []model.CyberGameList

	db := r.db.Table("Line_notifygame as line").
		Joins("LEFT JOIN Type_notify AS typenotify ON typenotify.id = line.typenotify_id").
		Joins("LEFT JOIN Admins AS ad ON ad.id = line.admin_id").
		Select("typenotify.name ,line.admin_id,line.id, line.token, line.typenotify_id, line.tag, line.status , line.created_at, line.updated_at")

	if query.Filter != "" {
		db = db.Where("line.admin_id = ?", query.Filter)
	}

	if err := db.
		Limit(query.Limit).
		Offset(query.Limit * query.Page).
		Find(&cybergame).
		Order("line.id ASC").
		Error; err != nil {
		return nil, 0, err
	}

	var total int64

	queryTotal := r.db.Table("Line_notifygame")

	if query.Filter != "" {
		queryTotal = queryTotal.Where("id = ?", query.Filter)
	}
	if err := queryTotal.Table("Line_notifygame").
		Select("id").
		Count(&total).
		Error; err != nil {
		return nil, 0, err
	}

	return cybergame, total, nil
}

func (r repo) GetLineNotifyById(id int64) (*model.Linenotify, error) {
	var linenotify model.Linenotify

	if err := r.db.Table("line_notify").
		Select("id, start_credit, token, notify_id, status , created_at, updated_at").
		Where("id = ?", id).
		First(&linenotify).
		Error; err != nil {
		return nil, err
	}
	return &linenotify, nil
}

func (r repo) GetLineNotify(req model.LinenotifyListRequest) (*model.SuccessWithPagination, error) {

	var list []model.LinenotifyListResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("line_notify")
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
		query := r.db.Table("line_notify")
		query = query.Select("id,start_credit, token, notify_id, status")
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
		list = []model.LinenotifyListResponse{}
	}
	result.List = list
	result.Total = total
	return &result, nil
}
func (r repo) CreateLineNotify(data model.LinenotifyCreateBody) error {
	if err := r.db.Table("line_notify").Create(&data).Error; err != nil {
		fmt.Println(data)
		return err
	}

	return nil
}

func (r repo) UpdateLineNotify(id int64, data model.LinenotifyUpdateBody) error {
	if err := r.db.Table("line_notify").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) GetLinenotifyGameById(id int64) (*model.LinenotifyGame, error) {
	var linenotifygame model.LinenotifyGame

	if err := r.db.Table("type_notify").
		Select("id, name, client_id, client_secret, response_type, redirect_uri, scope, state, status, created_at , updated_at").
		Where("status = 'ACTIVE' AND  id = ?", id).
		First(&linenotifygame).
		Error; err != nil {
		return nil, err
	}
	return &linenotifygame, nil
}

func (r repo) CreateLinenotifyCyberGame(data []model.LineNoifyCyberGameBody) error {
	if err := r.db.Table("Line_notifygame").Create(&data).Error; err != nil {
		fmt.Println(data)
		return err
	}

	return nil
}

func (r repo) UpdateLineNotifyCyberGame(id int64, data model.UpdateStatusCyberGame) error {

	if err := r.db.Table("Line_notifygame").
		Where("id = ?", id).
		Update("status", data.Status).
		Error; err != nil {
		return err
	}

	return nil
}

func (r repo) GetLinenotifyCyberGameById(id int64) (*model.LineNoifyCyberGame, error) {
	var linenotify model.LineNoifyCyberGame

	if err := r.db.Table("Line_notifygame as line").
		Joins("LEFT JOIN Type_notify AS typenotify ON typenotify.id = line.typenotify_id").
		Joins("LEFT JOIN Admins AS ad ON ad.id = line.admin_id").
		Select("typenotify.name ,line.admin_id,line.id, line.token, line.typenotify_id, line.tag, line.status , line.created_at, line.updated_at").
		Where("line.id = ?", id).
		First(&linenotify).
		Error; err != nil {
		return nil, err
	}
	return &linenotify, nil
}

func (r repo) UpdateLineNotifyTypeCyberGame(id int64, data model.UpdateStatusTypeCyberGame) error {

	if err := r.db.Table("Line_notifygame as line").
		Joins("LEFT JOIN Type_notify AS typenotify ON typenotify.id = line.typenotify_id").
		Joins("LEFT JOIN Admins AS ad ON ad.id = line.admin_id").
		Where("line.id = ?", id).
		Update("line.status", data.Status).
		Error; err != nil {
		return err
	}

	return nil
}
