package model

import (
	"time"
)

type Settingweb struct {
	Id             int64     `json:"id"`
	Logo           string    `json:"logo"`
	BackgrondColor string    `json:"backgrondcolor"`
	UserAuto       string    `json:"userAuto"`
	OtpRegister    string    `json:"otpRegister"`
	AutoWithdraw   string    `json:"autowithdraw"`
	TranWithdraw   string    `json:"tranWithdraw"`
	Register       string    `json:"register"`
	DepositFirst   string    `json:"depositFirst"`
	DepositNext    string    `json:"depositNext"`
	Withdraw       string    `json:"withdraw"`
	Line           string    `json:"line"`
	Url            string    `json:"url"`
	Opt            string    `json:"opt"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
type SettingwebResponse struct {
	Id             int64  `json:"id"`
	Logo           string `json:"logo"`
	BackgrondColor string `json:"backgrondcolor"`
	UserAuto       string `json:"userAuto"`
	OtpRegister    string `json:"otpRegister"`
	AutoWithdraw   string `json:"autowithdraw"`
	TranWithdraw   string `json:"tranWithdraw"`
	Register       string `json:"register"`
	DepositFirst   string `json:"depositFirst"`
	DepositNext    string `json:"depositNext"`
	Withdraw       string `json:"withdraw"`
	Line           string `json:"line"`
	Url            string `json:"url"`
	Opt            string `json:"opt"`
}
type SettingwebListResponse struct {
	Id    int `json:"id"`
	Total int `json:"total"`
}

type SettingwebParam struct {
	Id int64 `uri:"id" binding:"required"`
}

type SettingwebListRequest struct {
	Page    int    `form:"page" default:"1" min:"1"`
	Limit   int    `form:"limit" default:"10" min:"1" max:"100"`
	Search  string `form:"search"`
	SortCol string `form:"sortCol"`
	SortAsc string `form:"sortAsc"`
}

type SettingwebCreateBody struct {
	Logo           string `json:"logo"`
	BackgrondColor string `json:"backgrondcolor" validate:"required"`
	UserAuto       string `json:"userAuto" validate:"required"`
	OtpRegister    string `json:"otpRegister" validate:"required"`
	AutoWithdraw   string `json:"autowithdraw" validate:"required"`
	TranWithdraw   string `json:"tranWithdraw" validate:"required"`
	Register       string `json:"register" validate:"required"`
	DepositFirst   string `json:"depositFirst" validate:"required"`
	DepositNext    string `json:"depositNext" validate:"required"`
	Withdraw       string `json:"withdraw" validate:"required"`
	Line           string `json:"line" validate:"required"`
	Url            string `json:"url" validate:"required"`
	Opt            string `json:"opt" validate:"required"`
}
type SettingwebUpdateBody struct {
	Id             int64  `json:"id" validate:"required"`
	Logo           string `json:"logo" validate:"required"`
	BackgrondColor string `json:"backgrondcolor" validate:"required"`
	UserAuto       string `json:"userAuto" validate:"required"`
	OtpRegister    string `json:"otpRegister" validate:"required"`
	AutoWithdraw   string `json:"autowithdraw" validate:"required"`
	TranWithdraw   string `json:"tranWithdraw" validate:"required"`
	Register       string `json:"register" validate:"required"`
	DepositFirst   string `json:"depositFirst" validate:"required"`
	DepositNext    string `json:"depositNext" validate:"required"`
	Withdraw       string `json:"withdraw" validate:"required"`
	Line           string `json:"line" validate:"required"`
	Url            string `json:"url" validate:"required"`
	Opt            string `json:"opt" validate:"required"`
}
