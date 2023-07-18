package repository

import (
	"cybergame-api/model"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func NewFrontUserRepository(db *gorm.DB) FrontUserRepository {
	return &repo{db}
}

type FrontUserRepository interface {
	GetFrontUserLoginLogs(id int64) (*[]model.UserLoginLog, error)
	GetFrontUser(id int64) (*model.UserDetail, error)
	GetFrontUserList(query model.UserListQuery) (*[]model.UserList, *int64, error)
	GetFrontUserByPhone(phone string) (*model.User, error)
	GetFrontUpdateLogs(query model.UserUpdateQuery) (*[]model.UserUpdateLogResponse, *int64, error)
	CountFrontUser() (*int64, error)
	FrontGetUserOTP(data model.UserVerifyOTP) (*model.UserOTP, error)
	CheckFrontUser(username string) (bool, error)
	CheckFrontUserPhone(phone string) (bool, error)
	CheckFrontUserById(id int64) (*model.User, error)
	CheckFrontUserBank(bankNumber, bankCode string) (bool, error)
	CheckFrontUserTrueWallet(trueNumber string) (bool, error)
	FrontUserOtpCheckExpired(userId int64) (bool, error)
	CreateFrontUser(admin model.User) error
	FrontCreateUserOTP(user model.User) (*int64, error)
	FrontCheckOTPExpired(userId int64) (bool, error)
	FrontAddUserOTP(data *model.UserOTP) error
	FrontUpdateUser(userId int64, data model.FrontUserUpdate) error
	FrontUserUpdatePassword(userId int64, data model.UserUpdatePassword) error
	FrontUpdateUserOTP(data *model.UserOTP) error
	FrontUpdateUserOTPForget(data *model.UserOTP) error
	DeleteFrontUser(id int64) error
}

func (r repo) GetFrontUserLoginLogs(id int64) (*[]model.UserLoginLog, error) {

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

func (r repo) GetFrontUserList(query model.UserListQuery) (*[]model.UserList, *int64, error) {

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

func (r repo) GetFrontUserByPhone(phone string) (*model.User, error) {

	var user *model.User

	if err := r.db.Table("Users").
		Select("id, username, phone, password, verified_at").
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

func (r repo) GetFrontUser(id int64) (*model.UserDetail, error) {

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

func (r repo) GetFrontUpdateLogs(query model.UserUpdateQuery) (*[]model.UserUpdateLogResponse, *int64, error) {

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

func (r repo) CountFrontUser() (*int64, error) {

	var total int64

	if err := r.db.Model(model.User{}).Table("Users").
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	return &total, nil
}

func (r repo) FrontGetUserOTP(data model.UserVerifyOTP) (*model.UserOTP, error) {

	var otp *model.UserOTP

	if err := r.db.Model(&model.UserOTP{}).Table("User_otps").
		Select("id, code, verified_at, user_id, ref, expired_at").
		Where("user_id = ?", data.UserId).
		Where("code = ?", data.Code).
		First(&otp).
		Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return otp, nil
}

func (r repo) GetFrontUserOTPVerify(otpType string, userId int64) (*model.UserOTP, error) {

	var otp *model.UserOTP

	if err := r.db.Model(&model.UserOTP{}).Table("User_otps").
		Select("id, code, verified_at, user_id, ref, expired_at").
		Where("type = ?", otpType).
		Where("user_id = ?", userId).
		Order("created_at DESC").
		First(&otp).
		Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return otp, nil
}

func (r repo) CheckFrontUser(username string) (bool, error) {

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

func (r repo) CheckFrontUserPhone(phone string) (bool, error) {

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

func (r repo) CheckFrontUserById(id int64) (*model.User, error) {

	var user model.User

	if err := r.db.Table("Users").
		Select("is_reset_password, username").
		Where("id = ?", id).
		First(&user).
		Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r repo) CheckFrontUserBank(bankNumber, bankCode string) (bool, error) {

	var user model.User

	if err := r.db.Table("Users").
		Select("id").
		Where("bank_account = ?", bankNumber).
		Where("bank_code = ?", bankCode).
		First(&user).
		Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, err
	}

	if user.Id == 0 {
		return false, nil
	}

	return true, nil
}

func (r repo) CheckFrontUserTrueWallet(trueNumber string) (bool, error) {

	var user model.User

	if err := r.db.Table("Users").
		Select("id").
		Where("true_wallet = ?", trueNumber).
		First(&user).
		Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, err
	}

	if user.Id == 0 {
		return false, nil
	}

	return true, nil
}

func (r repo) CreateFrontUser(admin model.User) error {

	if err := r.db.Table("Users").
		Create(&admin).
		Error; err != nil {

		var dup *mysql.MySQLError
		if errors.As(err, &dup); dup.Number == 1062 {
			return errors.New("Phone already exists")
		}

		return err
	}

	return nil
}

func (r repo) FrontCreateUserOTP(user model.User) (*int64, error) {

	if err := r.db.Table("Users").
		Create(&user).
		Error; err != nil {

		var dup *mysql.MySQLError
		if errors.As(err, &dup); dup.Number == 1062 {

			if err := r.db.Table("Users").
				Select("id").
				Where("phone = ?", user.Phone).
				First(&user).
				Error; err != nil {
				return nil, err
			}

			return &user.Id, nil
		}

		return nil, err
	}

	return &user.Id, nil
}

func (r repo) FrontCheckOTPExpired(userId int64) (bool, error) {

	var total int64

	if err := r.db.Table("User_otps").
		Where("user_id = ?", userId).
		Where("verified_at IS NULL").
		Where("expired_at > ? ", time.Now()).
		Count(&total).
		Error; err != nil {
		return false, err
	}

	fmt.Println("count", total)

	return total > 0, nil
}

func (r repo) FrontAddUserOTP(data *model.UserOTP) error {

	if err := r.db.Table("User_otps").
		Create(&data).
		Error; err != nil {
		return err
	}

	return nil
}

func (r repo) FrontUpdateUser(userId int64, data model.FrontUserUpdate) error {

	if err := r.db.Model(model.User{}).Table("Users").
		Where("id = ?", userId).
		Updates(&data).
		Error; err != nil {
		return err
	}

	return nil
}

func (r repo) FrontUserUpdatePassword(userId int64, data model.UserUpdatePassword) error {

	obj := map[string]interface{}{
		"password":          data.Password,
		"is_reset_password": false,
	}

	if err := r.db.Table("Users").
		Where("id = ?", userId).
		Updates(obj).
		Error; err != nil {
		return err
	}

	return nil
}

func (r repo) FrontUpdateUserOTP(data *model.UserOTP) error {

	tx := r.db.Begin()

	if err := tx.Table("Users").
		Where("id = ?", data.UserId).
		Where("verified_at IS NOT NULL").
		Update("verified_at", time.Now()).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := r.db.Table("User_otps").
		Where("id = ?", data.Id).
		Where("user_id = ?", data.UserId).
		Where("code = ?", data.Code).
		Update("verified_at", time.Now()).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r repo) FrontUserOtpCheckExpired(userId int64) (bool, error) {

	var result *model.UserOTP

	if err := r.db.Table("User_otps").
		Where("user_id = ?", userId).
		Where("verified_at IS NULL").
		Where("expired_at > ?", time.Now()).
		First(&result).
		Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, err
	}

	return result != nil, nil
}

func (r repo) FrontUpdateUserOTPForget(data *model.UserOTP) error {

	tx := r.db.Begin()

	if err := tx.Table("Users").
		Where("id = ?", data.UserId).
		Where("is_reset_password = ?", false).
		Update("is_reset_password", true).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := r.db.Table("User_otps").
		Where("id = ?", data.Id).
		Where("user_id = ?", data.UserId).
		Where("code = ?", data.Code).
		Update("verified_at", time.Now()).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r repo) DeleteFrontUser(id int64) error {

	if err := r.db.Table("Users").
		Where("id = ?", id).
		Delete(&model.User{}).
		Error; err != nil {
		return err
	}

	return nil
}
