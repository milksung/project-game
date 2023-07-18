package model

import (
	"time"
)

type Scammer struct {
	Id          int64      `json:"id"`
	Fullname    *string    `json:"fullname"`
	Firstname   *string    `json:"firstname"`
	Lastname    *string    `json:"lastname"`
	Bankname    *string    `json:"bankname"`
	BankAccount *string    `json:"bankAccount"`
	Phone       *string    `json:"phone"`
	Reason      *string    `json:"reason"`
	CreatedAt   *time.Time `json:"createdAt"`
}

type CreateScammer struct {
	Fullname    *string `json:"fullname"`
	Bankname    *string `json:"bankname" validate:"max=50"`
	BankAccount *string `json:"bankAccount" validate:"max=15"`
	Phone       *string `json:"phone" validate:"required,min=10,max=12"`
	Reason      *string `json:"reason" validate:"max=255"`
}

type ScammertList struct {
	Id          int64      `json:"id"`
	Fullname    *string    `json:"fullname"`
	Bankname    *string    `json:"bankname"`
	BankAccount *string    `json:"bankAccount"`
	Phone       *string    `json:"phone"`
	Reason      *string    `json:"reason"`
	CreatedAt   *time.Time `json:"createdAt"`
}

type ScammerDetail struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type ScammerQuery struct {
	Page      int        `form:"page" default:"1"`
	Limit     int        `form:"limit" default:"10"`
	DateStart *time.Time `form:"dateStart" time_format:"2006-01-02T15:04:05+07:00" default:"2021-01-01T00:00:00+07:00"`
	DateEnd   *time.Time `form:"dateEnd" time_format:"2006-01-02T15:04:05+07:00" default:"2021-01-01T00:00:00+07:00"`
	BankName  *string    `form:"bankName"`
	Filter    *string    `form:"filter"`
}
