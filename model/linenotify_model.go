package model

import (
	"time"
)

type Linenotify struct {
	Id          int64      `json:"id"`
	StartCredit float32    `json:"startcredit" sql:"type:decimal(14,2);"`
	Token       string     `json:"token" validate:"required"`
	NotifyId    int64      `json:"notifyId" validate:"required"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}
type LinenotifyResponse struct {
	Id          int64   `json:"id"`
	StartCredit float32 `json:"startcredit" sql:"type:decimal(14,2);"`
	Token       string  `json:"token" validate:"required"`
	NotifyId    int64   `json:"notifyId" validate:"required"`
	Status      string  `json:"status"`
}
type LinenotifyListResponse struct {
	Id    int `json:"id"`
	Total int `json:"total"`
}

type LinenotifyParam struct {
	Id int64 `uri:"id" binding:"required"`
}

type LinenotifyListRequest struct {
	Page    int    `form:"page" default:"1" min:"1"`
	Limit   int    `form:"limit" default:"10" min:"1" max:"100"`
	Search  string `form:"search"`
	SortCol string `form:"sortCol"`
	SortAsc string `form:"sortAsc"`
}

type LinenotifyCreateBody struct {
	StartCredit float32 `json:"startcredit" sql:"type:decimal(14,2);"`
	Token       string  `json:"token" validate:"required"`
	NotifyId    int64   `json:"notifyId" validate:"required"`
	Status      string  `json:"status"`
}
type LinenotifyUpdateBody struct {
	StartCredit float32 `json:"startcredit" sql:"type:decimal(14,2);"`
	Token       string  `json:"token" validate:"required"`
	NotifyId    int64   `json:"notifyId" validate:"required"`
	Status      string  `json:"status"`
}

type LinenotifyUpdateRequest struct {
	StartCredit float32 `json:"startcredit" sql:"type:decimal(14,2);"`
	Token       string  `json:"token" validate:"required"`
	NotifyId    int64   `json:"notifyId" validate:"required"`
	Status      string  `json:"status"`
}

type LinenotifyGame struct {
	Id        int64      `json:"id"`
	Name      string     `json:"name" validate:"required"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

type LinenotifyGameResponse struct {
	Id     int64  `json:"id"`
	Name   string `json:"name" validate:"required"`
	Status string `json:"status"`
}

type LinenotifyGameParam struct {
	Id int64 `uri:"id" binding:"required"`
}

type LineNoifyCyberGame struct {
	Token     string     `json:"token" validate:"required"`
	Name      string     `json:"name"`
	AdminId   string     `json:"adminId"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

type LineNoifyCyberGameBody struct {
	Token        string  `json:"token" validate:"required"`
	TypenotifyId int64   `json:"typenotifyId"`
	AdminId      int64   `json:"adminId" validate:"required"`
	Type         []int64 `json:"type" validate:"required"`
}

type LineNoifyCyberGameList struct {
	TypenotifyId int64 `json:"typenotifyId"`
}

type LineNoifyCyberGameParam struct {
	Id int64 `uri:"id" binding:"required"`
}

type LineNoifyCyberGameUpdateBody struct {
	Id        string    `json:"-"`
	DeletedAt time.Time `json:"-"`
}

type LineNoifyCyberGameReponse struct {
	Token        int64  `json:"token" validate:"required"`
	Username     string `json:"username"`
	TypenotifyId string `json:"typeNotifyId"`
	Status       string `json:"status"`
}

type UpdateStatusCyberGame struct {
	Status    string    `json:"Status" validate:"required"`
	UpdatedAt time.Time `json:"-"`
}

type UpdateStatusTypeCyberGame struct {
	Status    string    `json:"Status" validate:"required"`
	UpdatedAt time.Time `json:"-"`
}

type CyberGameList struct {
	Id           int        `json:"id"`
	Name         string     `json:"name"`
	AdminId      int        `json:"adminId"`
	TypenotifyId int        `json:"typenotifyId"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
}

type CyberGameQuery struct {
	Page   int    `form:"page" validate:"min=1" default:"1"`
	Limit  int    `form:"limit" validate:"min=1,max=100" default:"10"`
	Filter string `form:"filter" example:""`
}
