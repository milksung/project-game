package repository

import (
	"cybergame-api/model"

	"gorm.io/gorm"
)

func NewScammerRepository(db *gorm.DB) ScammerRepository {
	return &repo{db}
}

type ScammerRepository interface {
	GetScammerList(query model.ScammerQuery) ([]model.ScammertList, int64, error)
	CreateScammer(scammer model.Scammer) error
}

func (r repo) GetScammerList(query model.ScammerQuery) ([]model.ScammertList, int64, error) {

	var scammers []model.ScammertList

	db := r.db.Table("Scammers")

	if query.DateStart != nil && query.DateEnd != nil {
		db = db.Where("created_at BETWEEN ? AND ?", query.DateStart, query.DateEnd)
	}

	if query.BankName != nil && *query.BankName != "" {
		db = db.Where("bankname = ?", query.BankName)
	}

	if query.Filter != nil && *query.Filter != "" {
		db = db.Where("fullname LIKE ? OR bankname LIKE ? OR bank_account LIKE ?", "%"+*query.Filter+"%", "%"+*query.Filter+"%", "%"+*query.Filter+"%")
	}

	if err := db.
		Limit(query.Limit).
		Offset(query.Limit * query.Page).
		Find(&scammers).
		Order("id desc").
		Error; err != nil {
		return nil, 0, err
	}

	var total int64
	queryTotal := r.db.Table("Scammers")

	if err := queryTotal.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return scammers, total, nil
}

func (r repo) CreateScammer(scammer model.Scammer) error {

	if err := r.db.Table("Scammers").
		Create(&scammer).
		Error; err != nil {
		return err
	}

	return nil
}
