package service

import (
	"cybergame-api/helper"
	"cybergame-api/model"
	"cybergame-api/repository"
	"fmt"
	"strings"
)

type AdminService interface {
	GetAdmin(id int64) (*model.AdminDetail, error)
	GetAdminList(query model.AdminListQuery) (*model.SuccessWithPagination, error)
	GetGroup(id int) (*model.AdminGroupPermissionResponse, error)
	GetGroupList(query model.AdminGroupQuery) (*model.SuccessWithPagination, error)
	Login(data model.LoginAdmin) (*string, error)
	Create(user *model.CreateAdmin) (error, *[]string)
	CreateGroup(data *model.AdminCreateGroup) error
	UpdateAdmin(adminId int64, data model.AdminBody) (error, *[]string)
	UpdateGroup(groupId int64, data *model.AdminUpdateGroup) error
	ResetPassword(adminId int64, body model.AdminUpdatePassword) error
	DeleteGroup(id int64) error
	DeletePermission(perm model.DeletePermission) error
	DeleteAdmin(id int64) error
}

const AdminLoginFailed = "Phone Or Password is incorrect"
const AdminNotFound = "Admin not found"
const AdminExist = "Admin already exist"
const AdminPhoneExist = "Phone already exist"
const AdminGroupNotFound = "Group not found"
const AdminFirstNameLastName = "Fullname must be firstname lastname"

type adminService struct {
	repo      repository.AdminRepository
	perRepo   repository.PermissionRepository
	groupRepo repository.GroupRepository
}

func NewAdminService(
	repo repository.AdminRepository,
	perRepo repository.PermissionRepository,
	groupRepo repository.GroupRepository,
) AdminService {
	return &adminService{repo, perRepo, groupRepo}
}

func (s *adminService) GetAdmin(id int64) (*model.AdminDetail, error) {

	admin, perList, group, err := s.repo.GetAdmin(id)
	if err != nil {

		if err.Error() == "record not found" {
			return nil, notFound(AdminNotFound)
		}

		return nil, err
	}

	var result model.AdminDetail
	result.Id = admin.Id
	result.Username = admin.Username
	result.Fullname = admin.Fullname
	result.Phone = admin.Phone
	result.Email = admin.Email
	result.Status = admin.Status
	result.Role = admin.Role
	result.PermissionList = *perList

	if group != nil {
		result.Group = group
	}

	return &result, nil
}

func (s *adminService) GetAdminList(query model.AdminListQuery) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&query.Page, &query.Limit); err != nil {
		return nil, err
	}

	list, total, err := s.repo.GetAdminList(query)
	if err != nil {
		return nil, err
	}

	result := &model.SuccessWithPagination{
		Message: "Success",
		List:    list,
		Total:   *total,
	}

	return result, nil
}

func (s *adminService) GetGroup(id int) (*model.AdminGroupPermissionResponse, error) {

	group, err := s.repo.GetGroup(id)
	if err != nil {

		if err.Error() == "record not found" {
			return nil, notFound(AdminGroupNotFound)
		}

		return nil, err
	}

	return group, nil
}

func (s *adminService) GetGroupList(query model.AdminGroupQuery) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&query.Page, &query.Limit); err != nil {
		return nil, err
	}

	list, total, err := s.repo.GetGroupList(query)
	if err != nil {
		return nil, err
	}

	result := &model.SuccessWithPagination{
		Message: "Success",
		List:    list,
		Total:   *total,
	}

	return result, nil
}

func (s *adminService) Login(data model.LoginAdmin) (*string, error) {

	user, err := s.repo.GetAdminByUsername(data)
	if err != nil {

		if err.Error() == "record not found" {
			return nil, notFound(AdminLoginFailed)
		}

		return nil, internalServerError(err.Error())
	}

	if user == nil {
		return nil, badRequest(AdminLoginFailed)
	}

	if err := helper.CompareAdminPassword(data.Password, user.Password); err != nil {
		return nil, badRequest(AdminLoginFailed)
	}

	token, err := helper.CreateJWTAdmin(*user)
	if err != nil {
		return nil, internalServerError(err.Error())
	}

	return &token, nil
}

func (s *adminService) Create(data *model.CreateAdmin) (error, *[]string) {

	username, err := s.repo.CheckAdmin(data.Username)
	if err != nil {
		return err, nil
	}

	if username {
		return badRequest(AdminExist), nil
	}

	email, err := s.repo.CheckPhone(data.Phone)
	if err != nil {
		return err, nil
	}

	if email {
		return badRequest(AdminPhoneExist), nil
	}

	checkGroup, err := s.groupRepo.CheckGroupExist(data.AdminGroupId)
	if err != nil {
		return internalServerError(err.Error()), nil
	}

	if !checkGroup {
		return badRequest(AdminGroupNotFound), nil
	}

	hashedPassword, err := helper.GenAdminPassword(data.Password)
	if err != nil {
		return internalServerError(err.Error()), nil
	}

	splitFullname := strings.Split(data.Fullname, " ")
	if len(splitFullname) == 1 || strings.Trim(data.Fullname, " ") == "" {
		return badRequest(AdminFirstNameLastName), nil
	}

	var firstname, lastname *string
	if len(splitFullname) == 2 {
		firstname = &splitFullname[0]
		lastname = &splitFullname[1]
	}

	if len(splitFullname) == 3 {
		firstname = &splitFullname[1]
		lastname = &splitFullname[2]
	}

	newUser := model.Admin{}
	newUser.Email = data.Email
	newUser.Username = data.Username
	newUser.Fullname = data.Fullname
	newUser.Firstname = *firstname
	newUser.Lastname = *lastname
	newUser.Password = string(hashedPassword)
	newUser.Role = "ADMIN"
	newUser.Status = data.Status
	newUser.Phone = data.Phone
	newUser.AdminGroupId = data.AdminGroupId

	return s.repo.CreateAdmin(newUser), nil
}

func (s *adminService) CreateGroup(data *model.AdminCreateGroup) error {

	checkGroup, err := s.groupRepo.CheckGroupExist(data.GroupId)
	if err != nil {
		return internalServerError(err.Error())
	}

	if !checkGroup {
		return badRequest(AdminGroupNotFound)
	}

	checkPermission, err := s.perRepo.CheckPerListExist(data.PermissionIds)
	if err != nil {
		return internalServerError(err.Error())
	}

	var idNotFound []string
	for _, j := range data.PermissionIds {

		exist := false

		for _, k := range checkPermission {
			if j == k {
				exist = true
			}
		}

		if !exist {
			idNotFound = append(idNotFound, fmt.Sprintf("%d", j))
		}
	}

	if len(idNotFound) > 0 {
		return badRequest(fmt.Sprintf("Permission id %s not found", strings.Join(idNotFound, ",")))
	}

	var list []model.AdminPermissionList

	for _, id := range data.PermissionIds {
		list = append(list, model.AdminPermissionList{
			GroupId:      data.GroupId,
			PermissionId: id,
		})
	}

	if err := s.repo.CreateGroupAdmin(list); err != nil {
		return err
	}

	return nil
}

func (s *adminService) UpdateAdmin(adminId int64, body model.AdminBody) (error, *[]string) {

	var data model.UpdateAdmin

	if body.AdminGroupId != nil {
		checkGroup, err := s.groupRepo.CheckGroupExist(*body.AdminGroupId)
		if err != nil {
			return internalServerError(err.Error()), nil
		}

		if !checkGroup {
			return notFound(AdminGroupNotFound), nil
		}

		data.AdminGroupId = body.AdminGroupId
	}

	var oldGroupId *int64

	if body.AdminGroupId != nil {

		getGroupId, err := s.repo.GetAdminGroup(adminId)
		if err != nil {
			return internalServerError(err.Error()), nil
		}

		oldGroupId = &getGroupId.AdminGroupId
	}

	data.Email = body.Email
	data.Status = body.Status

	splitFullname := strings.Split(body.Fullname, " ")
	if len(splitFullname) == 1 || strings.Trim(body.Fullname, " ") == "" {
		return badRequest(AdminFirstNameLastName), nil
	}

	if len(splitFullname) == 2 {
		data.Firstname = splitFullname[0]
		data.Lastname = splitFullname[1]
	}

	if len(splitFullname) == 3 {
		data.Firstname = splitFullname[1]
		data.Lastname = splitFullname[2]
	}

	data.Fullname = body.Fullname

	return s.repo.UpdateAdmin(adminId, oldGroupId, data), nil
}

func (s *adminService) UpdateGroup(groupId int64, data *model.AdminUpdateGroup) error {

	checkGroup, err := s.groupRepo.CheckGroupExist(groupId)
	if err != nil {
		return internalServerError(err.Error())
	}

	if !checkGroup {
		return badRequest(AdminGroupNotFound)
	}

	checkPermission, err := s.perRepo.CheckPerListExist(data.PermissionIds)
	if err != nil {
		return internalServerError(err.Error())
	}

	var idNotFound []string
	for _, j := range data.PermissionIds {

		exist := false

		for _, k := range checkPermission {
			if j == k {
				exist = true
			}
		}

		if !exist {
			idNotFound = append(idNotFound, fmt.Sprintf("%d", j))
		}
	}

	if len(idNotFound) > 0 {
		return badRequest(fmt.Sprintf("ไม่พบ Permission id %s", strings.Join(idNotFound, ",")))
	}

	var list []model.AdminPermissionList

	for _, id := range data.PermissionIds {
		list = append(list, model.AdminPermissionList{
			GroupId:      groupId,
			PermissionId: id,
		})
	}

	if err := s.repo.UpdateGroup(groupId, data.Name, list); err != nil {
		return err
	}

	return nil
}

func (s *adminService) ResetPassword(adminId int64, body model.AdminUpdatePassword) error {

	checkAdmin, err := s.repo.CheckAdminById(adminId)
	if err != nil {
		return internalServerError(err.Error())
	}

	if !checkAdmin {
		return notFound(AdminNotFound)
	}

	newPasword, err := helper.GenAdminPassword(body.Password)
	if err != nil {
		return internalServerError(err.Error())
	}

	body.Password = newPasword

	if err := s.repo.UpdatePassword(adminId, body); err != nil {
		return err
	}

	return nil
}

func (s *adminService) DeleteGroup(id int64) error {

	exist, err := s.repo.GetAdminCount(id)
	if err != nil {
		return err
	}

	if exist > 0 {
		return badRequest("ยังมีผู้ใช้งานอยู่ในกลุ่มนี้ กรุณาย้ายกลุ่มผู้ใช้งานออกก่อน")
	}

	if err := s.groupRepo.DeleteGroup(id); err != nil {
		return err
	}

	return nil
}

func (s *adminService) DeletePermission(perm model.DeletePermission) error {

	if err := s.perRepo.DeletePermission(perm); err != nil {
		return err
	}

	return nil
}

func (s *adminService) DeleteAdmin(id int64) error {

	if err := s.repo.DeleteAdmin(id); err != nil {
		return err
	}

	return nil
}
