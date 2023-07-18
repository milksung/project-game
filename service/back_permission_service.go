package service

import (
	"cybergame-api/model"
	"cybergame-api/repository"
)

type PermissionService interface {
	Create(user *model.CreatePermission) error
}

const PermNotFound = "Permission not found"

type permissionService struct {
	repo repository.PermissionRepository
}

func NewPermissionService(
	repo repository.PermissionRepository,
) PermissionService {
	return &permissionService{repo}
}

func (s *permissionService) Create(data *model.CreatePermission) error {

	if err := s.repo.CreatePermission(data); err != nil {
		return err
	}

	return nil
}
