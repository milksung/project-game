package repository

import (
	"cybergame-api/model"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &repo{db}
}

type AdminRepository interface {
	GetAdmin(id int64) (*model.Admin, *[]model.PermissionList, *model.GroupDetail, error)
	GetAdminList(query model.AdminListQuery) (*[]model.AdminList, *int64, error)
	GetGroup(groupId int) (*model.AdminGroupPermissionResponse, error)
	GetGroupList(query model.AdminGroupQuery) (*[]model.GroupCountList, *int64, error)
	GetAdminByUsername(data model.LoginAdmin) (*model.Admin, error)
	GetAdminGroup(adminId int64) (*model.AdminGroupId, error)
	CheckAdmin(username string) (bool, error)
	CheckAdminById(id int64) (bool, error)
	CheckPhone(phone string) (bool, error)
	GetAdminCount(id int64) (int64, error)
	CreateAdmin(admin model.Admin) error
	CreateGroupAdmin(data []model.AdminPermissionList) error
	UpdateGroup(groupId int64, name string, data []model.AdminPermissionList) error
	UpdateAdmin(adminId int64, OldGroupId *int64, data model.UpdateAdmin) error
	UpdatePassword(adminId int64, data model.AdminUpdatePassword) error
	DeleteAdmin(adminId int64) error
}

func (r repo) GetAdminList(query model.AdminListQuery) (*[]model.AdminList, *int64, error) {

	var err error
	var list []model.AdminList
	var total int64

	exec := r.db.Model(model.Admin{}).Table("Admins a").
		Joins("LEFT JOIN Admin_groups ag ON ag.id = a.admin_group_id").
		Select("a.id, a.username, a.fullname, a.phone, a.email, ag.name AS admin_group_name, a.status, a.role")

	if query.Search != "" {
		exec = exec.Where("a.username LIKE ?", "%"+query.Search+"%")
	}

	if query.Status != "" {
		exec = exec.Where("a.status = ?", query.Status)
	}

	if query.AdminGroupId != "" {
		exec = exec.Where("a.admin_group_id = ?", query.AdminGroupId)
	}

	if err := exec.
		Limit(query.Limit).
		Offset(query.Limit * query.Page).
		Find(&list).
		Error; err != nil {
		return nil, nil, err
	}

	execTotal := r.db.Model(model.Admin{}).Table("Admins a").
		Joins("LEFT JOIN Admin_groups ag ON ag.id = a.admin_group_id").
		Select("a.id")

	if query.Search != "" {
		execTotal = execTotal.Where("a.username LIKE ?", "%"+query.Search+"%")
	}

	if query.Status != "" {
		execTotal = execTotal.Where("a.status = ?", query.Status)
	}

	if query.AdminGroupId != "" {
		execTotal = execTotal.Where("a.admin_group_id = ?", query.AdminGroupId)
	}

	if err = execTotal.
		Count(&total).
		Error; err != nil {
		return nil, nil, err
	}

	return &list, &total, nil
}

func (r repo) GetAdmin(id int64) (*model.Admin, *[]model.PermissionList, *model.GroupDetail, error) {

	var admin *model.Admin
	var permission *[]model.PermissionList
	var group *model.GroupDetail

	if err := r.db.Model(model.Admin{}).Table("Admins").
		Select("id, username, fullname, phone, email, role, status, admin_group_id").
		Where("id = ?", id).
		First(&admin).
		Error; err != nil {
		return nil, nil, nil, err
	}

	println("admin.AdminGroupId", admin.AdminGroupId)

	if err := r.db.Table("Admin_group_permissions agp").
		Joins("LEFT JOIN Permissions p ON p.id = agp.permission_id").
		Select("p.id, p.name").
		Where("agp.group_id = ?", admin.AdminGroupId).
		Find(&permission).
		Error; err != nil {
		return nil, nil, nil, err
	}

	if err := r.db.Table("Admin_groups").
		Select("id, name").
		Where("id = ?", admin.AdminGroupId).
		First(&group).
		Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return admin, permission, nil, nil
		}

		return nil, nil, nil, err
	}

	return admin, permission, group, nil
}

func (r repo) GetGroup(groupId int) (*model.AdminGroupPermissionResponse, error) {

	var group model.Group
	var permission []model.PermissionList

	if err := r.db.Table("Admin_groups").
		Select("id, name").
		Where("id = ?", groupId).
		First(&group).
		Error; err != nil {
		return nil, err
	}

	if err := r.db.Table("Permissions p").
		Select("p.id, p.name, gp.is_read, gp.is_write").
		Joins("LEFT JOIN Admin_group_permissions gp ON gp.permission_id = p.id").
		Where("gp.group_id = ?", groupId).
		Where("gp.deleted_at IS NULL").
		Find(&permission).
		Error; err != nil {
		return nil, err
	}

	var result model.AdminGroupPermissionResponse
	result.Id = group.Id
	result.Name = group.Name
	result.Permissions = permission

	return &result, nil
}

func (r repo) GetGroupList(query model.AdminGroupQuery) (*[]model.GroupCountList, *int64, error) {

	var exec = r.db.Model(model.Group{}).
		Table("Admin_groups").
		Select("id, name, admin_count").
		Limit(query.Limit).
		Offset(query.Limit * query.Page)

	if query.Search != "" {
		exec = exec.Where("name LIKE ?", "%"+query.Search+"%")
	}

	var list []model.GroupCountList
	if err := exec.
		Find(&list).
		Error; err != nil {
		return nil, nil, err
	}

	var execTotal = r.db.Model(model.Group{}).
		Table("Admin_groups").
		Select("id")

	if query.Search != "" {
		execTotal = execTotal.Where("name LIKE ?", "%"+query.Search+"%")
	}

	var total int64
	if err := execTotal.
		Count(&total).
		Error; err != nil {
		return nil, nil, err
	}

	return &list, &total, nil
}

func (r repo) GetAdminByUsername(data model.LoginAdmin) (*model.Admin, error) {
	var admin model.Admin

	if err := r.db.Table("Admins").
		Select("id, username, phone, password, email, role").
		Where("username = ?", data.Username).
		First(&admin).
		Error; err != nil {
		return nil, err
	}

	if admin.Id != 0 {
		if err := r.db.Table("Admins").
			Where("id = ?", admin.Id).
			Updates(model.AdminLoginUpdate{
				Ip:        data.Ip,
				LogedinAt: time.Now(),
			}).
			Error; err != nil {
			return nil, err
		}
	}

	return &admin, nil
}

func (r repo) GetAdminGroup(adminId int64) (*model.AdminGroupId, error) {

	var admin *model.AdminGroupId

	if err := r.db.Table("Admins").
		Select("admin_group_id").
		Where("id = ?", adminId).
		First(&admin).
		Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return admin, nil
}

func (r repo) CheckAdmin(username string) (bool, error) {
	var user model.Admin

	if err := r.db.Table("Admins").
		Where("username = ?", username).
		First(&user).
		Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, err
	}

	if user.Id != 0 {
		return false, nil
	}

	return true, nil
}

func (r repo) CheckAdminById(id int64) (bool, error) {
	var user model.Admin

	if err := r.db.Table("Admins").
		Where("id = ?", id).
		First(&user).
		Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r repo) CheckPhone(phone string) (bool, error) {
	var user model.Admin

	if err := r.db.Table("Admins").
		Where("phone = ?", phone).
		First(&user).
		Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r repo) GetAdminCount(id int64) (int64, error) {

	var adminCount int64
	if err := r.db.Table("Admin_groups").
		Select("admin_count").
		Where("id = ?", id).
		Find(&adminCount).
		Limit(1).
		Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}

		return 0, err
	}

	return adminCount, nil
}

func (r repo) CreateAdmin(admin model.Admin) error {

	tx := r.db.Begin()

	if err := tx.Table("Admins").
		Create(&admin).
		Error; err != nil {
		tx.Rollback()

		var dup *mysql.MySQLError
		if errors.As(err, &dup); dup.Number == 1062 {
			return errors.New("Username or Email already exists")
		}

		return err
	}

	if err := tx.Table("Admin_groups").
		Where("id = ?", admin.AdminGroupId).
		Update("admin_count", gorm.Expr("admin_count + ?", 1)).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r repo) CreateGroupAdmin(data []model.AdminPermissionList) error {

	if err := r.db.Table("Admin_group_permissions").
		Create(&data).
		Error; err != nil {
		return err
	}

	return nil
}

func (r repo) UpdateGroup(groupId int64, name string, data []model.AdminPermissionList) error {

	tx := r.db.Begin()

	if err := tx.Table("Admin_groups").
		Where("id = ?", groupId).
		Update("name", name).
		Error; err != nil {

		tx.Rollback()

		var dup *mysql.MySQLError
		if errors.As(err, &dup); dup.Number == 1062 {
			return errors.New("มีชื่อกลุ่มนี้อยู่แล้ว")
		}

		return err
	}

	if err := tx.Table("Admin_group_permissions").
		Where("group_id = ?", groupId).
		Delete(&model.AdminGroupPermission{}).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Table("Admin_group_permissions").
		Create(&data).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r repo) UpdateAdmin(adminId int64, OldGroupId *int64, data model.UpdateAdmin) error {

	tx := r.db.Begin()

	if err := tx.Model(model.Admin{}).Table("Admins").
		Where("id = ?", adminId).
		Updates(&data).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if *data.AdminGroupId != *OldGroupId {

		if err := tx.Table("Admin_groups").
			Where("id = ?", data.AdminGroupId).
			Update("admin_count", gorm.Expr("admin_count + ?", 1)).
			Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Table("Admin_groups").
			Where("id = ?", OldGroupId).
			Update("admin_count", gorm.Expr("admin_count - ?", 1)).
			Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r repo) UpdatePassword(adminId int64, data model.AdminUpdatePassword) error {

	if err := r.db.Model(model.Admin{}).Table("Admins").
		Where("id = ?", adminId).
		Update("password", data.Password).
		Error; err != nil {
		return err
	}

	return nil
}

func (r repo) DeleteAdmin(adminId int64) error {

	var groupId int64

	tx := r.db.Begin()

	if err := tx.Table("Admins").
		Select("admin_group_id").
		Where("id = ?", adminId).
		Find(&groupId).
		Limit(1).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Table("Admins").
		Where("id = ?", adminId).
		Delete(&model.Admin{}).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Table("Admin_groups").
		Where("id = ?", groupId).
		Update("admin_count", gorm.Expr("admin_count - ?", 1)).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
