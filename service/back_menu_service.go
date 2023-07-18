package service

import (
	"cybergame-api/model"
	"cybergame-api/repository"
	"fmt"
	"strings"
)

type MenuService interface {
	GetMenu(adminId int64) ([]model.Menu, error)
}

type menuService struct {
	PermRepo repository.PermissionRepository
}

func NewMenuService(
	PermRepo repository.PermissionRepository,
) MenuService {
	return &menuService{PermRepo}
}

func (s *menuService) GetMenu(adminId int64) ([]model.Menu, error) {

	perms, adminPers, err := s.PermRepo.GetPermissions(adminId)
	if err != nil {
		return nil, err
	}

	fmt.Println("adminPers", adminPers)

	var menu []model.Menu

	for _, per := range perms {

		subMenu := []model.SubMenu{}

		for _, subPer := range perms {
			if !subPer.Main {

				count := len(strings.Split(subPer.PermissionKey, "_"))

				if per.PermissionKey == strings.Join(strings.Split(subPer.PermissionKey, "_")[:count-1], "_") {
					subMenu = append(subMenu, model.SubMenu{
						Id:    subPer.Id,
						Title: subPer.Name,
						Name:  subPer.PermissionKey,
					})
				}
			}
		}

		if per.Main {

			menu = append(menu, model.Menu{
				Id:      per.Id,
				Title:   per.Name,
				Name:    per.PermissionKey,
				List:    &subMenu,
				Managed: checkRead(*per, adminPers),
			})
		}
	}

	return menu, nil
}

func checkRead(per model.Permission, adminPers *[]model.AdminPermission) bool {

	read := false

	for _, adminPer := range *adminPers {
		if adminPer.PermissionId == per.Id {
			read = true
		}
	}

	return read
}
