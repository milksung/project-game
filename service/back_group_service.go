package service

import (
	"cybergame-api/model"
	"cybergame-api/repository"
	"fmt"
	"strings"
)

type GroupService interface {
	Create(user *model.CreateGroup) error
}

const GroupNotFound = "Permission not found"

type groupService struct {
	repo    repository.GroupRepository
	perRepo repository.PermissionRepository
}

func NewGroupService(
	repo repository.GroupRepository,
	perRepo repository.PermissionRepository,
) GroupService {
	return &groupService{repo, perRepo}
}

func (s *groupService) Create(data *model.CreateGroup) error {

	var group model.Group
	group.Name = data.Name

	perIdExists, err := s.perRepo.CheckPerListExist(data.PermissionIds)
	if err != nil {
		return internalServerError(err.Error())
	}

	var idNotFound []string
	for _, j := range data.PermissionIds {

		exist := false

		for _, k := range perIdExists {
			if j == k {
				exist = true
			}
		}

		if !exist {
			idNotFound = append(idNotFound, fmt.Sprintf("%v", j))
		}
	}

	if len(idNotFound) > 0 {
		return badRequest(fmt.Sprintf("Permission id %s not found", strings.Join(idNotFound, ",")))
	}

	if err := s.repo.CreateGroup(group, data.PermissionIds); err != nil {
		return err
	}

	return nil
}
