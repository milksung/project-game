package repository

import (
	"cybergame-api/model"
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func NewGroupRepository(db *gorm.DB) GroupRepository {
	return &repo{db}
}

type GroupRepository interface {
	CheckGroupByName(name string) (bool, error)
	CheckGroupExist(id int64) (bool, error)
	CreateGroup(group model.Group, Permissions []int64) error
	DeleteGroup(id int64) error
}

func (r repo) CheckGroupByName(name string) (bool, error) {

	var count int64

	if err := r.db.Table("Admin_groups").
		Where("name = ?", name).
		Count(&count).
		Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func (r repo) CheckGroupExist(id int64) (bool, error) {

	var count int64

	if err := r.db.Table("Admin_groups").
		Where("id = ?", id).
		Count(&count).
		Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func (r repo) CreateGroup(group model.Group, Permissions []int64) error {

	tx := r.db.Begin()

	if err := tx.Table("Admin_groups").
		Create(&group).
		Error; err != nil {
		tx.Rollback()

		var dup *mysql.MySQLError
		if errors.As(err, &dup); dup.Number == 1062 {
			return errors.New("Group name already exists")
		}

		return err
	}

	// if err := tx.Table("Permissions").
	// 	Create(&Permissions).
	// 	Error; err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	var adminGroupPerms []model.AdminGroupPermission

	for _, id := range Permissions {
		adminGroupPerms = append(adminGroupPerms, model.AdminGroupPermission{
			GroupId:      group.Id,
			PermissionId: id,
			DeletedAt:    nil,
		})
	}

	if err := tx.Table("Admin_group_permissions").
		Create(&adminGroupPerms).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r repo) DeleteGroup(id int64) error {

	tx := r.db.Begin()

	if err := tx.Table("Admin_groups").
		Where("id = ?", id).
		Delete(&model.Group{}).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Table("Admin_group_permissions").
		Where("group_id = ?", id).
		Delete(&model.AdminGroupPermission{}).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
