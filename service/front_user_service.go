package service

import (
	"cybergame-api/helper"
	"cybergame-api/model"
	"cybergame-api/repository"
	"os"
	"strconv"
	"strings"
	"time"
)

type FrontUserService interface {
	GetFrontUserLoginLogs(id int64) (*[]model.UserLoginLog, error)
	GetFrontUser(id int64) (*model.UserDetail, error)
	FrontUserLogin(body model.UserLogin) (*string, error)
	CreateFrontUser(user *model.CreateUser) error
	FrontUserSendOTPRegister(body model.UserSendOTP) error
	FrontUserVerifyOTPRegister(body model.UserVerifyOTP) (*int64, error)
	FrontUserSendOTPFotget(body model.UserSendOTP) error
	FrontUserVerifyOTPForget(body model.UserVerifyOTP) (*int64, error)
	FrontUserUpdateInfo(userId int64, body model.FrontUserUpdate) error
	FrontUserResetPassword(userId int64, body model.UserUpdatePassword) error
	FrontUserChangePassword(userId int64, body model.UserUpdatePassword) error
}

const FrontUserLoginFailed = "เบอร์โทรศัพท์หรือรหัสผ่านไม่ถูกต้อง"
const FrontUserNotFound = "ไม่พบผู้ใช้งาน"
const FrontUserExist = "ผู้ใช้งานนี้มีอยู่ในระบบแล้ว"
const FrontUserFullName = "ต้องมีชื่อและนามสกุล"
const FrontUserOtpExpired = "OTP หมดอายุแล้ว"
const FrontUserOtpNotMatch = "OTP ไม่ถูกต้อง"
const FrontUserOtpUsed = "OTP นี้ได้ถูกใช้งานไปแล้ว"
const FrontUserIsVerified = "ผู้ใช้งานนี้ได้ทำการยืนยันข้อมูลแล้ว"
const FrontUserFullnameMustBeSpace = "ชื่อและนามสกุลต้องมีช่องว่าง"
const FrontUserBankExist = "เลขที่บัญชีธนาคารมีอยู่ในระบบแล้ว"
const FrontUserTrueWalletExist = "บัญชีทรูวอลเล็ตมีอยู่ในระบบแล้ว"
const FrontUserPasswordIsReset = "รหัสผ่านถูกรีเซ็ตแล้ว"

type frontUserService struct {
	repo             repository.FrontUserRepository
	agentConnectRepo repository.AgentConnectRepository
}

func NewFrontUserService(
	repo repository.FrontUserRepository,
	agentConnectRepo repository.AgentConnectRepository,
) FrontUserService {
	return &frontUserService{repo, agentConnectRepo}
}

func (s *frontUserService) GetFrontUserLoginLogs(id int64) (*[]model.UserLoginLog, error) {

	logs, err := s.repo.GetFrontUserLoginLogs(id)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (s *frontUserService) GetFrontUser(id int64) (*model.UserDetail, error) {

	admin, err := s.repo.GetFrontUser(id)
	if err != nil {

		if err.Error() == "record not found" {
			return nil, notFound(UserNotFound)
		}

		return nil, err
	}

	return admin, nil
}

func (s *frontUserService) FrontUserLogin(body model.UserLogin) (*string, error) {

	user, err := s.repo.GetFrontUserByPhone(body.Phone)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, notFound(FrontUserNotFound)
	}

	if err := helper.CompareUserPassword(body.Password, user.Password); err != nil {
		return nil, badRequest(FrontUserLoginFailed)
	}

	token, err := helper.CreateJWTUser(*user)
	if err != nil {
		return nil, err
	}

	agentName := os.Getenv("AGENT_NAME")
	sign := agentName + *user.Username + body.Password
	timeNow := time.Now()
	agentData := model.AGCLogin{
		Username:  *user.Username,
		Partner:   agentName,
		Timestamp: timeNow.Unix(),
		Sign:      helper.CreateSign(sign, timeNow),
		Domain:    "http://test.com",
		Lang:      "th-th",
		IsMobile:  false,
		Ip:        body.IP,
	}

	if err = s.agentConnectRepo.Login(agentData); err != nil {
		return nil, err
	}

	return &token, nil
}

func (s *frontUserService) CreateFrontUser(data *model.CreateUser) error {

	phone, err := s.repo.CheckFrontUserPhone(data.Phone)
	if err != nil {
		return err
	}

	if phone {
		return badRequest(FrontUserExist)
	}

	hashedPassword, err := helper.GenUserPassword(data.Password)
	if err != nil {
		return internalServerError(err.Error())
	}

	var newUser model.User
	newUser.Partner = &data.Partner
	newUser.MemberCode = &data.MemberCode
	newUser.Phone = data.Phone
	newUser.Promotion = &data.Promotion
	newUser.Password = string(hashedPassword)
	newUser.Status = "ACTIVE"
	newUser.Fullname = data.Fullname
	newUser.Bankname = data.Bankname
	newUser.BankCode = data.BankCode
	newUser.BankAccount = data.BankAccount
	newUser.Channel = data.Channel
	newUser.TrueWallet = data.TrueWallet
	newUser.Contact = data.Contact
	newUser.Note = data.Note
	newUser.Course = data.Course
	newUser.IpRegistered = data.IpRegistered

	splitFullname := strings.Split(data.Fullname, " ")
	if len(splitFullname) == 1 {
		return badRequest(FrontUserFullName)
	}

	var firstname, lastname *string
	if len(splitFullname) == 2 {
		firstname = &splitFullname[0]
		lastname = &splitFullname[1]
		newUser.Firstname = *firstname
		newUser.Lastname = *lastname
	}

	if len(splitFullname) == 3 {
		firstname = &splitFullname[1]
		lastname = &splitFullname[2]
		newUser.Firstname = *firstname
		newUser.Lastname = *lastname
	}

	if err := s.repo.CreateFrontUser(newUser); err != nil {
		return err
	}

	return nil
}

func (s *frontUserService) FrontUserSendOTPRegister(body model.UserSendOTP) error {

	var otp model.UserOTP
	otp.Code = helper.GenNumber(6)
	otp.Ref = helper.GenStringUpper(6)
	otp.Type = "REGISTER"
	otp.ExpiredAt = time.Now().Add(time.Minute * 5)

	var user model.User
	user.Phone = body.Phone
	user.Status = "ACTIVE"

	userId, err := s.repo.FrontCreateUserOTP(user)
	if err != nil {
		return err
	}

	otp.UserId = *userId

	expired, err := s.repo.FrontCheckOTPExpired(*userId)
	if err != nil {
		return err
	}

	if expired {
		return badRequest("OTP ยังไม่หมดอายุ")
	}

	if err := s.repo.FrontAddUserOTP(&otp); err != nil {
		return err
	}

	return nil
}

func (s *frontUserService) FrontUserVerifyOTPRegister(body model.UserVerifyOTP) (*int64, error) {

	user, err := s.repo.GetFrontUserByPhone(body.Phone)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, notFound(FrontUserNotFound)
	}

	if user.VerifiedAt != nil {
		return nil, badRequest(FrontUserIsVerified)
	}

	body.UserId = user.Id

	otp, err := s.repo.FrontGetUserOTP(body)
	if err != nil {
		return nil, err
	}

	if otp == nil {
		return nil, notFound(FrontUserOtpNotMatch)
	}

	if otp.VerifiedAt != nil {
		return nil, badRequest(FrontUserOtpUsed)
	}

	if otp.ExpiredAt.Before(time.Now()) {
		return nil, badRequest(FrontUserOtpExpired)
	}

	if err := s.repo.FrontUpdateUserOTP(otp); err != nil {
		return nil, err
	}

	return &user.Id, nil
}

func (s *frontUserService) FrontUserSendOTPFotget(body model.UserSendOTP) error {

	user, err := s.repo.GetFrontUserByPhone(body.Phone)
	if err != nil {
		return err
	}

	if user == nil {
		return notFound(FrontUserExist)
	}

	checkOtp, err := s.repo.FrontUserOtpCheckExpired(user.Id)
	if err != nil {
		return err
	}

	if checkOtp {
		return badRequest("OTP ยังไม่หมดอายุ")
	}

	var otp model.UserOTP
	otp.Code = helper.GenNumber(6)
	otp.Ref = helper.GenStringUpper(6)
	otp.Type = "FORGET"
	otp.UserId = user.Id
	otp.ExpiredAt = time.Now().Add(time.Minute * 5)

	if err := s.repo.FrontAddUserOTP(&otp); err != nil {
		return err
	}

	return nil
}

func (s *frontUserService) FrontUserVerifyOTPForget(body model.UserVerifyOTP) (*int64, error) {

	user, err := s.repo.GetFrontUserByPhone(body.Phone)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, notFound(FrontUserNotFound)
	}

	verifyOtp := model.UserVerifyOTP{
		UserId: user.Id,
		Code:   body.Code,
	}

	otp, err := s.repo.FrontGetUserOTP(verifyOtp)
	if err != nil {
		return nil, err
	}

	if otp == nil || otp.VerifiedAt != nil {
		return nil, notFound(FrontUserOtpNotMatch)
	}

	if otp.ExpiredAt.Before(time.Now()) {
		return nil, badRequest(FrontUserOtpExpired)
	}

	if err := s.repo.FrontUpdateUserOTPForget(otp); err != nil {
		return nil, err
	}

	return &user.Id, nil
}

func (s *frontUserService) FrontUserUpdateInfo(userId int64, body model.FrontUserUpdate) error {

	user, err := s.repo.GetFrontUser(userId)
	if err != nil {

		if err.Error() == "record not found" {
			return notFound(UserNotFound)
		}

		return internalServerError(err.Error())
	}

	if user == nil {
		return notFound(UserNotFound)
	}

	if user.Fullname != "" {
		return badRequest(FrontUserIsVerified)
	}

	if body.Fullname != "" {

		splitFullname := strings.Split(body.Fullname, " ")
		if len(splitFullname) == 1 {
			return badRequest(FrontUserFullnameMustBeSpace)
		}

		var firstname, lastname *string
		if len(splitFullname) == 2 {
			firstname = &splitFullname[0]
			lastname = &splitFullname[1]
			body.Firstname = firstname
			body.Lastname = lastname
		}

		if len(splitFullname) == 3 {
			firstname = &splitFullname[1]
			lastname = &splitFullname[2]
			body.Firstname = firstname
			body.Lastname = lastname
		}
	}

	passwordOrigin := body.Password

	if body.Password != "" {

		newPasword, err := helper.GenUserPassword(body.Password)
		if err != nil {
			return internalServerError(err.Error())
		}

		body.Password = newPasword
	}

	checkBank, err := s.repo.CheckFrontUserBank(body.BankAccount, body.BankCode)
	if err != nil {
		return err
	}

	if checkBank {
		return badRequest(FrontUserBankExist)
	}

	checkTrue, err := s.repo.CheckFrontUserTrueWallet(body.TrueWallet)
	if err != nil {
		return err
	}

	if checkTrue {
		return badRequest(FrontUserTrueWalletExist)
	}

	countUser, err := s.repo.CountFrontUser()
	if err != nil {
		return err
	}

	agentStart, _ := strconv.Atoi(os.Getenv("AGENT_START_NUMBER"))

	agentTotal := *countUser + int64(agentStart)
	conv := strconv.FormatInt(agentTotal, 10)
	agentName := os.Getenv("AGENT_NAME")
	username := agentName + helper.IncrementNumber(conv)
	body.Username = username
	sign := agentName + username
	timeNow := time.Now()
	agentData := model.AGCRegister{
		Username:  username,
		Agentname: agentName,
		Fullname:  body.Fullname,
		Password:  passwordOrigin,
		Currency:  "THB",
		Dob:       "1990-01-01",
		Mobile:    user.Phone,
		Ip:        body.Ip,
		Timestamp: timeNow.Unix(),
		Sign:      helper.CreateSign(sign, timeNow),
	}

	body.VerifiedAt = time.Now()

	if err := s.repo.FrontUpdateUser(userId, body); err != nil {
		return err
	}

	if err := s.agentConnectRepo.Register(agentData); err != nil {

		agentData.Fullname = ""
		if err := s.repo.FrontUpdateUser(userId, body); err != nil {
			return internalServerError(ServerError)
		}

		return internalServerError(ServerError)
	}

	return nil

}

func (s *frontUserService) FrontUserResetPassword(userId int64, body model.UserUpdatePassword) error {

	user, err := s.repo.CheckFrontUserById(userId)
	if err != nil {
		return internalServerError(err.Error())
	}

	if user == nil {
		return notFound(UserNotFound)
	}

	if !user.IsResetPassword {
		return badRequest(FrontUserPasswordIsReset)
	}

	newPasword, err := helper.GenUserPassword(body.Password)
	if err != nil {
		return internalServerError(err.Error())
	}

	agentName := os.Getenv("AGENT_NAME")
	sign := agentName + *user.Username + body.Password
	timeNow := time.Now()
	agentData := model.AGCChangePassword{
		PlayerName:  *user.Username,
		Partner:     agentName,
		NewPassword: body.Password,
		Timestamp:   timeNow.Unix(),
		Sign:        helper.CreateSign(sign, timeNow),
	}

	if err = s.agentConnectRepo.ChangePassword(agentData); err != nil {
		return err
	}

	body.Password = newPasword

	if err := s.repo.FrontUserUpdatePassword(userId, body); err != nil {
		return err
	}

	return nil
}

func (s *frontUserService) FrontUserChangePassword(userId int64, body model.UserUpdatePassword) error {

	user, err := s.repo.CheckFrontUserById(userId)
	if err != nil {
		return internalServerError(err.Error())
	}

	if user == nil {
		return notFound(UserNotFound)
	}

	newPasword, err := helper.GenUserPassword(body.Password)
	if err != nil {
		return internalServerError(err.Error())
	}

	body.Password = newPasword

	if err := s.repo.FrontUserUpdatePassword(userId, body); err != nil {
		return err
	}

	return nil
}
