package repository

import (
	"cybergame-api/model"
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func NewUserRepository(db *gorm.DB) UserRepository {
	return &repo{db}
}

type UserRepository interface {
	GetUserLoginLogs(id int64) (*[]model.UserLoginLog, error)
	GetUser(id int64) (*model.UserDetail, error)
	GetUserList(query model.UserListQuery) (*[]model.UserList, *int64, error)
	GetUserByPhone(phone string) (*model.UserByPhone, error)
	GetUpdateLogs(query model.UserUpdateQuery) (*[]model.UserUpdateLogResponse, *int64, error)
	CheckUser(username string) (bool, error)
	CheckUserPhone(phone string) (bool, error)
	CheckUserById(id int64) (bool, error)
	CountUser() (*int64, error)
	CreateUser(user model.User) (int64, error)
	UpdateUser(userId int64, data model.UpdateUser, changes []model.UserUpdateLogs) error
	UpdateUserPassword(userId int64, data model.UserUpdatePassword) error
	DeleteUser(id int64) error
}

func (r repo) GetUserLoginLogs(id int64) (*[]model.UserLoginLog, error) {

	var logs []model.UserLoginLog

	if err := r.db.Table("User_login_logs").
		Where("user_id = ?", id).
		Find(&logs).
		Order("created_at DESC").
		Error; err != nil {
		return nil, err
	}

	return &logs, nil
}

func (r repo) GetUserList(query model.UserListQuery) (*[]model.UserList, *int64, error) {

	var err error
	var list []model.UserList
	var total int64

	exec := r.db.Model(model.User{}).Table("Users").
		Select("id, member_code, promotion, fullname, bankname, bank_account, channel, credit, ip, ip_registered, created_at, updated_at, logedin_at")

	if query.Search != "" {
		exec = exec.Where("username LIKE ?", "%"+query.Search+"%")
	}

	if query.From != nil && query.To != nil {
		exec = exec.Where("created_at BETWEEN ? AND ?", query.From, query.To)
	}

	if query.NonMember {
		exec = exec.Where("member_code = ?", "")
	}

	if !query.NonMember {
		exec = exec.Where("member_code != ?", "")
	}

	if err := exec.
		Limit(query.Limit).
		Offset(query.Limit * query.Page).
		Find(&list).
		Error; err != nil {
		return nil, nil, err
	}

	execTotal := r.db.Model(model.User{}).Table("Users").
		Select("id")

	if query.Search != "" {
		execTotal = execTotal.Where("username LIKE ?", "%"+query.Search+"%")
	}

	if query.From != nil && query.To != nil {
		execTotal = execTotal.Where("created_at BETWEEN ? AND ?", query.From, query.To)
	}

	if query.NonMember {
		execTotal = execTotal.Where("member_code = ?", "")
	}

	if !query.NonMember {
		execTotal = execTotal.Where("member_code != ?", "")
	}

	if err = execTotal.
		Count(&total).
		Error; err != nil {
		return nil, nil, err
	}

	return &list, &total, nil
}

func (r repo) GetUserByPhone(phone string) (*model.UserByPhone, error) {

	var user *model.UserByPhone

	if err := r.db.Table("Users").
		Select("id, phone").
		Where("phone = ?", phone).
		First(&user).
		Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

func (r repo) GetUser(id int64) (*model.UserDetail, error) {

	var admin *model.UserDetail

	if err := r.db.Model(model.User{}).Table("Users").
		Select("id, partner, member_code, phone, promotion, bankname, bank_account, fullname, channel, true_wallet, contact, note, course").
		Where("id = ?", id).
		First(&admin).
		Error; err != nil {
		return nil, err
	}

	return admin, nil
}

func (r repo) GetUpdateLogs(query model.UserUpdateQuery) (*[]model.UserUpdateLogResponse, *int64, error) {

	var logs []model.UserUpdateLogResponse
	var total int64

	exec := r.db.Table("User_update_logs")

	if query.Search != "" {
		exec = exec.Where("description LIKE ?", "%"+query.Search+"%")
	}

	if query.From != nil && query.To != nil {
		exec = exec.Where("created_at BETWEEN ? AND ?", query.From, query.To)
	}

	if err := exec.
		Limit(query.Limit).
		Offset(query.Limit * query.Page).
		Find(&logs).
		Order("created_at DESC").
		Error; err != nil {
		return nil, nil, err
	}

	execTotal := r.db.Table("User_update_logs")

	if query.Search != "" {
		execTotal = execTotal.Where("description LIKE ?", "%"+query.Search+"%")
	}

	if query.From != nil && query.To != nil {
		execTotal = execTotal.Where("created_at BETWEEN ? AND ?", query.From, query.To)
	}

	if err := execTotal.
		Count(&total).
		Error; err != nil {
		return nil, nil, err
	}

	return &logs, &total, nil
}

func (r repo) CheckUser(username string) (bool, error) {

	var user model.User

	if err := r.db.Table("Users").
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

func (r repo) CheckUserPhone(phone string) (bool, error) {

	var user model.User

	if err := r.db.Table("Users").
		Where("phone = ?", phone).
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

func (r repo) CheckUserById(id int64) (bool, error) {
	var user model.User

	if err := r.db.Table("Users").
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

func (r repo) CountUser() (*int64, error) {

	var total int64

	if err := r.db.Model(model.User{}).Table("Users").
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	return &total, nil
}

func (r repo) CreateUser(user model.User) (int64, error) {

	if err := r.db.Table("Users").
		Create(&user).
		Error; err != nil {

		var dup *mysql.MySQLError
		if errors.As(err, &dup); dup.Number == 1062 {
			return 0, errors.New("Phone already exists")
		}

		return 0, err
	}

	return user.Id, nil
}

func (r repo) UpdateUser(userId int64, data model.UpdateUser, changes []model.UserUpdateLogs) error {

	tx := r.db.Begin()

	if err := tx.Model(model.User{}).Table("Users").
		Where("id = ?", userId).
		Updates(&data).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if len(changes) > 0 {
		if err := tx.Table("User_update_logs").
			Create(&changes).
			Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r repo) UpdateUserPassword(userId int64, data model.UserUpdatePassword) error {

	if err := r.db.Table("Users").
		Where("id = ?", userId).
		Update("password", data.Password).
		Error; err != nil {
		return err
	}

	return nil
}

func (r repo) DeleteUser(id int64) error {

	if err := r.db.Table("Users").
		Where("id = ?", id).
		Delete(&model.User{}).
		Error; err != nil {
		return err
	}

	return nil
}
