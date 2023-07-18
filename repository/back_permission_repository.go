package repository

import (
	"cybergame-api/model"
	"errors"

	"gorm.io/gorm"
)

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &repo{db}
}

type PermissionRepository interface {
	GetPermissions(adminId int64) ([]*model.Permission, *[]model.AdminPermission, error)
	CheckPerListExist(ids []int64) ([]int64, error)
	CheckPerListAndGroupId(groupId int64, perIds []int64) ([]model.PermissionList, error)
	CreatePermission(data *model.CreatePermission) error
	DeletePermission(perm model.DeletePermission) error
}

func (r repo) GetPermissions(adminId int64) ([]*model.Permission, *[]model.AdminPermission, error) {

	var list []*model.Permission
	var adminPers *[]model.AdminPermission
	var admin model.AdminGroupId

	if err := r.db.Table("Admins").
		Select("admin_group_id").
		Where("id = ?", adminId).
		First(&admin).Error; err != nil {

		if err.Error() == "record not found" {
			return nil, nil, errors.New("record not found")
		}

		return nil, nil, err
	}

	if err := r.db.Table("Permissions").
		Select("id, name, permission_key, main").
		Where("Permission_key IS NOT NULL").
		Order("position ASC").
		Find(&list).Error; err != nil {
		return nil, nil, err
	}

	if err := r.db.Table("Admin_group_permissions agp").
		Joins("LEFT JOIN Permissions p ON p.id = agp.permission_id").
		Select("agp.permission_id").
		Where("agp.group_id = ?", admin.AdminGroupId).
		Find(&adminPers).Error; err != nil {
		return nil, nil, err
	}

	return list, adminPers, nil
}

func (r repo) CheckPerListExist(ids []int64) ([]int64, error) {

	var list []int64

	if err := r.db.Table("Permissions").Select("id").Where("id IN ?", ids).Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (r repo) CheckPerListAndGroupId(groupId int64, perIds []int64) ([]model.PermissionList, error) {

	var list []model.PermissionList

	if err := r.db.Table("Admin_group_permissions agp").
		Joins("LEFT JOIN Permissions p ON p.id = agp.permission_id").
		Select("agp.permission_id AS id, p.name, agp.is_read, agp.is_write").
		Where("agp.group_id = ? AND agp.permission_id IN (?)", groupId, perIds).
		Find(&list).
		Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (r repo) CreatePermission(data *model.CreatePermission) error {

	if err := r.db.Table("Permissions").
		Create(&data.Permissions).
		Error; err != nil {
		return err
	}

	return nil
}

func (r repo) DeletePermission(perm model.DeletePermission) error {

	tx := r.db.Begin()

	if err := tx.Table("Permissions").
		Where("id IN (?)", perm.PermissionIds).
		Delete(&model.Permission{}).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Table("Admin_group_permissions").
		Where("permission_id IN (?)", perm.PermissionIds).
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
