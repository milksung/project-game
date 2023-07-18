package model

import (
	"time"

	"gorm.io/gorm"
)

type BankStatement struct {
	Id                int64          `json:"id" gorm:"primaryKey"`
	AccountId         int64          `json:"accountId"`
	ExternalId        int64          `json:"externalId"`
	Amount            float64        `json:"amount" sql:"type:decimal(14,2);"`
	Detail            string         `json:"detail"`
	BankId            int64          `json:"bankId"`
	StatementType     string         `json:"statementType"`
	FromBankId        int64          `json:"fromBankId"`
	FromBankCode      string         `json:"fromBankCode"`
	FromAccountNumber string         `json:"fromAccountNumber"`
	FromBankName      string         `json:"fromBankName"`
	FromBankIconUrl   string         `json:"fromBankIconUrl"`
	TransferAt        time.Time      `json:"transferAt"`
	Status            string         `json:"status"`
	CreatedAt         time.Time      `json:"createAt"`
	UpdatedAt         *time.Time     `json:"updateAt"`
	DeletedAt         gorm.DeletedAt `json:"deleteAt"`
}

type GetByIdRequest struct {
	Id int64 `uri:"id" binding:"required"`
}
type SimpleListRequest struct {
	Search  string `form:"search" extensions:"x-order:1"`
	Page    int    `form:"page" extensions:"x-order:2" default:"1" min:"1"`
	Limit   int    `form:"limit" extensions:"x-order:3" default:"10" min:"1" max:"100"`
	SortCol string `form:"sortCol" extensions:"x-order:4"`
	SortAsc string `form:"sortAsc" extensions:"x-order:5"`
}

type BankStatementSummary struct {
	TotalPendingStatementCount int64 `json:"totalPendingStatementCount"`
	TotalPendingDepositCount   int64 `json:"totalPendingDepositCount"`
	TotalPendingWithdrawCount  int64 `json:"totalPendingWithdrawCount"`
}

type BankStatementListRequest struct {
	AccountId        string `form:"accountId" extensions:"x-order:1"`
	StatementType    string `form:"statementType" extensions:"x-order:2"`
	FromTransferDate string `form:"fromTransferDate" extensions:"x-order:3"`
	ToTransferDate   string `form:"toTransferDate" extensions:"x-order:4"`
	Search           string `form:"search" extensions:"x-order:5"`
	SimilarId        int64  `form:"similarId" extensions:"x-order:6"`
	Status           string `form:"status" extensions:"x-order:7"`
	Page             int    `form:"page" extensions:"x-order:7" default:"1" min:"1"`
	Limit            int    `form:"limit" extensions:"x-order:8" default:"10" min:"1" max:"100"`
	SortCol          string `form:"sortCol" extensions:"x-order:9"`
	SortAsc          string `form:"sortAsc" extensions:"x-order:10"`
}

type BankStatementCreateBody struct {
	Id                int64     `json:"id"`
	AccountId         int64     `json:"accountId"`
	ExternalId        int64     `json:"externalId"`
	Amount            float64   `json:"amount" sql:"type:decimal(14,2);"`
	Detail            string    `json:"detail"`
	FromBankId        int64     `json:"fromBankId"`
	FromAccountNumber string    `json:"fromAccountNumber"`
	StatementType     string    `json:"statementType"`
	TransferAt        time.Time `json:"transferAt"`
	Status            string    `json:"-"`
}

type BankStatementMatchRequest struct {
	UserId              int64     `json:"userId" validate:"required"`
	ConfirmedAt         time.Time `json:"-"`
	ConfirmedByUserId   int64     `json:"-"`
	ConfirmedByUsername string    `json:"-"`
}

type BankStatementUpdateBody struct {
	Status string `json:"status" validate:"required"`
}

type BankStatementResponse struct {
	Id              int64          `json:"id" gorm:"primaryKey"`
	AccountId       int64          `json:"accountId"`
	ExternalId      int64          `json:"externalId"`
	AccountName     string         `json:"accountName"`
	AccountNumber   string         `json:"accountNumber"`
	BankName        string         `json:"bankName"`
	Amount          float64        `json:"amount" sql:"type:decimal(14,2);"`
	Detail          string         `json:"detail"`
	FromBankId      int64          `json:"fromBankId"`
	FromBankName    string         `json:"fromBankName"`
	FromBankIconUrl string         `json:"fromBankIconUrl"`
	StatementType   string         `json:"statementType"`
	TransferAt      time.Time      `json:"transferAt"`
	Status          string         `json:"status"`
	CreatedAt       time.Time      `json:"createAt"`
	UpdatedAt       *time.Time     `json:"updateAt"`
	DeletedAt       gorm.DeletedAt `json:"deleteAt"`
}

type BankTransaction struct {
	Id                  int64          `json:"id" gorm:"primaryKey"`
	MemberCode          string         `json:"memberCode"`
	UserId              int64          `json:"userId"`
	TransferType        string         `json:"transferType"`
	PromotionId         int64          `json:"promotionId"`
	FromAccountId       int64          `json:"fromAccountId"`
	FromBankId          int64          `json:"fromBankId"`
	FromBankCode        string         `json:"fromBankCode"`
	FromBankName        string         `json:"fromBankName"`
	FromAccountName     string         `json:"fromAccountName"`
	FromAccountNumber   string         `json:"fromAccountNumber"`
	ToAccountId         int64          `json:"toAccountId"`
	ToBankId            int64          `json:"toBankId"`
	ToBankName          string         `json:"toBankName"`
	ToBankCode          string         `json:"toBankCode"`
	ToAccountName       string         `json:"toAccountName"`
	ToAccountNumber     string         `json:"toAccountNumber"`
	CreditAmount        float64        `json:"creditAmount" sql:"type:decimal(14,2);"`
	PaidAmount          float64        `json:"paidAmount" sql:"type:decimal(14,2);"`
	DepositChannel      string         `json:"depositChannel"`
	OverAmount          float64        `json:"overAmount" sql:"type:decimal(14,2);"`
	BonusAmount         float64        `json:"bonusAmount" sql:"type:decimal(14,2);"`
	BonusReason         string         `json:"bonusReason"`
	BeforeAmount        float64        `json:"beforeAmount" sql:"type:decimal(14,2);"`
	AfterAmount         float64        `json:"afterAmount" sql:"type:decimal(14,2);"`
	BankChargeAmount    float64        `json:"bankChargeAmount" sql:"type:decimal(14,2);"`
	TransferAt          *time.Time     `json:"transferAt"`
	CreatedByUserId     int64          `json:"createdByUserId"`
	CreatedByUsername   string         `json:"createdByUsername"`
	CancelRemark        string         `json:"cancelRemark"`
	CanceledAt          *time.Time     `json:"canceledAt"`
	CanceledByUserId    int64          `json:"canceledByUserId"`
	CanceledByUsername  string         `json:"canceledByUsername"`
	ConfirmedAt         *time.Time     `json:"confirmedAt"`
	ConfirmedByUserId   int64          `json:"confirmedByUserId"`
	ConfirmedByUsername string         `json:"confirmedByUsername"`
	RemovedAt           *time.Time     `json:"removedAt"`
	RemovedByUserId     int64          `json:"removedByUserId"`
	RemovedByUsername   string         `json:"removedByUsername"`
	Status              string         `json:"status"`
	StatusDetail        string         `json:"statusDetail"`
	IsAutoCredit        bool           `json:"isAutoCredit"`
	CreatedAt           time.Time      `json:"createAt"`
	UpdatedAt           *time.Time     `json:"updateAt"`
	DeletedAt           gorm.DeletedAt `json:"deleteAt"`
}

type BankTransactionGetRequest struct {
	Id int64 `uri:"id" binding:"required"`
}

type BankTransactionListRequest struct {
	MemberCode       string `form:"memberCode" extensions:"x-order:2"`
	UserId           string `form:"userId" extensions:"x-order:3"`
	FromTransferDate string `form:"fromTransferDate" extensions:"x-order:4"`
	ToTransferDate   string `form:"toTransferDate" extensions:"x-order:5"`
	TransferType     string `form:"transferType" extensions:"x-order:6"`
	TransferStatus   string `form:"transferStatus" extensions:"x-order:7"`
	Search           string `form:"search" extensions:"x-order:8"`
	Page             int    `form:"page" extensions:"x-order:9" default:"1" min:"1"`
	Limit            int    `form:"limit" extensions:"x-order:10" default:"10" min:"1" max:"100"`
	SortCol          string `form:"sortCol" extensions:"x-order:11"`
	SortAsc          string `form:"sortAsc" extensions:"x-order:12"`
}

type BankTransactionCreateBody struct {
	Id                int64      `json:"-"`
	MemberCode        string     `json:"memberCode" validate:"required"`
	UserId            int64      `json:"-"`
	TransferType      string     `json:"transferType" validate:"required" example:"deposit"`
	PromotionId       *int64     `json:"promotionId"`
	FromAccountId     *int64     `json:"fromAccountId"`
	FromBankId        *int64     `json:"-"`
	FromAccountName   *string    `json:"-"`
	FromAccountNumber *string    `json:"-"`
	ToAccountId       *int64     `json:"toAccountId"`
	ToBankId          *int64     `json:"-"`
	ToAccountName     *string    `json:"-"`
	ToAccountNumber   *string    `json:"-"`
	CreditAmount      float64    `json:"creditAmount" validate:"required"`
	PaidAmount        float64    `json:"-"`
	DepositChannel    string     `json:"depositChannel"`
	OverAmount        float64    `json:"overAmount"`
	BonusAmount       float64    `json:"bonusAmount"`
	BeforeAmount      float64    `json:"-"`
	AfterAmount       float64    `json:"-"`
	TransferAt        *time.Time `json:"transferAt" example:"2023-05-31T22:33:44+07:00"`
	CreatedByUserId   int64      `json:"-"`
	CreatedByUsername string     `json:"-"`
	Status            string     `json:"-"`
	IsAutoCredit      bool       `json:"isAutoCredit"`
}

type BonusTransactionCreateBody struct {
	MemberCode        string    `json:"memberCode" validate:"required"`
	UserId            int64     `json:"-"`
	TransferType      string    `json:"-"`
	ToAccountId       int64     `json:"-"`
	ToBankId          int64     `json:"-"`
	ToAccountName     string    `json:"-"`
	ToAccountNumber   string    `json:"-"`
	BonusAmount       float64   `json:"bonusAmount" validate:"required"`
	BonusReason       string    `json:"bonusReason"`
	BeforeAmount      float64   `json:"-"`
	AfterAmount       float64   `json:"-"`
	TransferAt        time.Time `json:"transferAt" validate:"required" example:"2023-05-31T22:33:44+07:00"`
	CreatedByUserId   int64     `json:"-"`
	CreatedByUsername string    `json:"-"`
	Status            string    `json:"-"`
}

type BankTransactionUpdateRequest struct {
	IsAutoCredit      *bool  `json:"isAutoCredit"`
	UpdatedByUserId   int64  `json:"-"`
	UpdatedByUsername string `json:"-"`
}

type BankTransactionUpdateBody struct {
	IsAutoCredit      bool      `json:"isAutoCredit"`
	Status            string    `json:"status"`
	RemovedAt         time.Time `json:"removedAt" example:"2023-05-31T22:33:44+07:00"`
	RemovedByUserId   int64     `json:"removedByUserId"`
	RemovedByUsername string    `json:"removedByUsername"`
}

type BankTransactionResponse struct {
	Id                  int64          `json:"id" gorm:"primaryKey"`
	UserId              int64          `json:"userId"`
	MemberCode          string         `json:"memberCode"`
	UserUsername        string         `json:"userUsername"`
	UserFullname        string         `json:"userFullname"`
	TransferType        string         `json:"transferType"`
	PromotionId         int64          `json:"promotionId"`
	FromAccountId       int64          `json:"fromAccountId"`
	FromBankId          int64          `json:"fromBankId"`
	FromBankName        string         `json:"fromBankName"`
	FromAccountName     string         `json:"fromAccountName"`
	FromAccountNumber   string         `json:"fromAccountNumber"`
	ToAccountId         int64          `json:"toAccountId"`
	ToBankId            int64          `json:"toBankId"`
	ToBankName          string         `json:"toBankName"`
	ToAccountName       string         `json:"toAccountName"`
	ToAccountNumber     string         `json:"toAccountNumber"`
	CreditAmount        float64        `json:"creditAmount" sql:"type:decimal(14,2);"`
	PaidAmount          float64        `json:"paidAmount" sql:"type:decimal(14,2);"`
	DepositChannel      string         `json:"depositChannel"`
	OverAmount          float64        `json:"overAmount" sql:"type:decimal(14,2);"`
	BonusAmount         float64        `json:"bonusAmount" sql:"type:decimal(14,2);"`
	BonusReason         string         `json:"bonusReason"`
	BeforeAmount        float64        `json:"beforeAmount" sql:"type:decimal(14,2);"`
	AfterAmount         float64        `json:"afterAmount" sql:"type:decimal(14,2);"`
	BankChargeAmount    float64        `json:"bankChargeAmount" sql:"type:decimal(14,2);"`
	TransferAt          *time.Time     `json:"transferAt"`
	CreatedByUserId     int64          `json:"createdByUserId"`
	CreatedByUsername   string         `json:"createdByUsername"`
	CancelRemark        string         `json:"cancelRemark"`
	CanceledAt          *time.Time     `json:"canceledAt"`
	CanceledByUserId    int64          `json:"canceledByUserId"`
	CanceledByUsername  string         `json:"canceledByUsername"`
	ConfirmedAt         *time.Time     `json:"confirmedAt"`
	ConfirmedByUserId   int64          `json:"confirmedByUserId"`
	ConfirmedByUsername string         `json:"confirmedByUsername"`
	RemovedAt           *time.Time     `json:"removedAt"`
	RemovedByUserId     int64          `json:"removedByUserId"`
	RemovedByUsername   string         `json:"removedByUsername"`
	Status              string         `json:"status"`
	StatusDetail        string         `json:"statusDetail"`
	IsAutoCredit        bool           `json:"isAutoCredit"`
	CreatedAt           time.Time      `json:"createAt"`
	UpdatedAt           *time.Time     `json:"updateAt"`
	DeletedAt           gorm.DeletedAt `json:"deleteAt"`
}

type BankTransactionStatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type BankDepositTransStatusCounts struct {
	AllCount           int64 `json:"allCount"`
	PendingCount       int64 `json:"pendingCount"`
	PendingCreditCount int64 `json:"pendingCreditCount"`
	FinishedCount      int64 `json:"finishedCount"`
	FailedCount        int64 `json:"failedCount"`
}
type BankWithdrawTransStatusCounts struct {
	AllCount             int64 `json:"allCount"`
	PendingCreditCount   int64 `json:"pendingCreditCount"`
	PendingTransferCount int64 `json:"pendingTransferCount"`
	FinishedCount        int64 `json:"finishedCount"`
	FailedCount          int64 `json:"failedCount"`
}

type PendingDepositTransactionListRequest struct {
	FromTransferDate string `form:"fromTransferDate" extensions:"x-order:3"`
	ToTransferDate   string `form:"toTransferDate" extensions:"x-order:4"`
	Search           string `form:"search" extensions:"x-order:5"`
	Page             int    `form:"page" extensions:"x-order:6" default:"1" min:"1"`
	Limit            int    `form:"limit" extensions:"x-order:7" default:"10" min:"1" max:"100"`
	SortCol          string `form:"sortCol" extensions:"x-order:8"`
	SortAsc          string `form:"sortAsc" extensions:"x-order:9"`
}

type PendingWithdrawTransactionListRequest struct {
	FromTransferDate string `form:"fromTransferDate" extensions:"x-order:3"`
	ToTransferDate   string `form:"toTransferDate" extensions:"x-order:4"`
	Search           string `form:"search" extensions:"x-order:5"`
	Page             int    `form:"page" extensions:"x-order:6" default:"1" min:"1"`
	Limit            int    `form:"limit" extensions:"x-order:7" default:"10" min:"1" max:"100"`
	SortCol          string `form:"sortCol" extensions:"x-order:8"`
	SortAsc          string `form:"sortAsc" extensions:"x-order:9"`
}

type FinishedTransactionListRequest struct {
	FromTransferDate string `form:"fromTransferDate" extensions:"x-order:1"`
	ToTransferDate   string `form:"toTransferDate" extensions:"x-order:2"`
	AccountId        string `form:"accountId" extensions:"x-order:3"`
	TransferType     string `form:"transferType" extensions:"x-order:4"`
	Search           string `form:"search" extensions:"x-order:5"`
	Page             int    `form:"page" extensions:"x-order:6" default:"1" min:"1"`
	Limit            int    `form:"limit" extensions:"x-order:7" default:"10" min:"1" max:"100"`
	SortCol          string `form:"sortCol" extensions:"x-order:8"`
	SortAsc          string `form:"sortAsc" extensions:"x-order:9"`
}

type RemovedTransactionListRequest struct {
	FromTransferDate string `form:"fromTransferDate" extensions:"x-order:1"`
	ToTransferDate   string `form:"toTransferDate" extensions:"x-order:2"`
	AccountId        string `form:"accountId" extensions:"x-order:3"`
	TransferType     string `form:"transferType" extensions:"x-order:4"`
	Search           string `form:"search" extensions:"x-order:5"`
	Page             int    `form:"page" extensions:"x-order:6" default:"1" min:"1"`
	Limit            int    `form:"limit" extensions:"x-order:7" default:"10" min:"1" max:"100"`
	SortCol          string `form:"sortCol" extensions:"x-order:8"`
	SortAsc          string `form:"sortAsc" extensions:"x-order:9"`
}

type BankTransactionCancelBody struct {
	Status             string    `json:"-"`
	CancelRemark       string    `json:"cancelRemark" validate:"required"`
	CanceledAt         time.Time `json:"-"`
	CanceledByUserId   int64     `json:"-"`
	CanceledByUsername string    `json:"-"`
}

type BankConfirmDepositRequest struct {
	TransferAt          *time.Time `json:"transferAt"`
	SlipUrl             *string    `json:"slipUrl"`
	BonusAmount         *float64   `json:"bonusAmount"`
	ConfirmedAt         time.Time  `json:"-"`
	ConfirmedByUserId   int64      `json:"-"`
	ConfirmedByUsername string     `json:"-"`
}

type BankConfirmCreditWithdrawRequest struct {
	FromAccountId       *int64    `json:"fromAccountId"`
	CreditAmount        *float64  `json:"creditAmount"`
	BankChargeAmount    *float64  `json:"bankChargeAmount"`
	ConfirmedAt         time.Time `json:"-"`
	ConfirmedByUserId   int64     `json:"-"`
	ConfirmedByUsername string    `json:"-"`
}
type BankConfirmTransferWithdrawRequest struct {
	FromAccountId       *int64    `json:"fromAccountId"`
	BankChargeAmount    *float64  `json:"bankChargeAmount"`
	ConfirmedAt         time.Time `json:"-"`
	ConfirmedByUserId   int64     `json:"-"`
	ConfirmedByUsername string    `json:"-"`
}

type BankDepositTransactionConfirmBody struct {
	TransferAt          time.Time `json:"transferAt"`
	BonusAmount         float64   `json:"bonusAmount"`
	Status              string    `json:"status"`
	ConfirmedAt         time.Time `json:"confirmedAt"`
	ConfirmedByUserId   int64     `json:"confirmedByUserId"`
	ConfirmedByUsername string    `json:"confirmedByUsername"`
}

type BankWithdrawTransactionConfirmBody struct {
	FromAccountId       *int64    `json:"fromAccountId"`
	TransferAt          time.Time `json:"transferAt"`
	CreditAmount        float64   `json:"creditAmount"`
	BankChargeAmount    float64   `json:"bankChargeAmount"`
	Status              string    `json:"status"`
	ConfirmedAt         time.Time `json:"confirmedAt"`
	ConfirmedByUserId   int64     `json:"confirmedByUserId"`
	ConfirmedByUsername string    `json:"confirmedByUsername"`
}

type CreateBankTransactionActionBody struct {
	Id                  int64      `json:"id"`
	ActionKey           string     `json:"actionKey"`
	TransactionId       int64      `json:"transactionId"`
	UserId              int64      `json:"userId"`
	TransferType        string     `json:"transferType"`
	FromAccountId       int64      `json:"fromAccountId"`
	ToAccountId         int64      `json:"toAccountId"`
	JsonBefore          string     `json:"jsonBefore"`
	TransferAt          *time.Time `json:"transferAt"`
	SlipUrl             string     `json:"slipUrl"`
	BonusAmount         float64    `json:"bonusAmount"`
	CreditAmount        float64    `json:"creditAmount"`
	BankChargeAmount    float64    `json:"bankChargeAmount"`
	ConfirmedAt         time.Time  `json:"confirmedAt"`
	ConfirmedByUserId   int64      `json:"confirmedByUserId"`
	ConfirmedByUsername string     `json:"confirmedByUsername"`
}

type CreateBankStatementActionBody struct {
	StatementId         int64     `json:"statementId"`
	UserId              int64     `json:"userId"`
	ActionType          string    `json:"actionType"`
	AccountId           int64     `json:"accountId"`
	JsonBefore          string    `json:"jsonBefore"`
	ConfirmedAt         time.Time `json:"confirmedAt"`
	ConfirmedByUserId   int64     `json:"confirmedByUserId"`
	ConfirmedByUsername string    `json:"confirmedByUsername"`
}
type BankTransactionRemoveBody struct {
	Status            string    `json:"-" validate:"required"`
	RemovedAt         time.Time `json:"removedAt"`
	RemovedByUserId   int64     `json:"removedByUserId"`
	RemovedByUsername string    `json:"removedByUsername"`
}
type BankAutoWithdrawCondition struct {
	TransId                 int64   `json:"TransId"`
	TransStatus             string  `json:"TransStatus"`
	UserId                  int64   `json:"UserId"`
	FromAccountId           int64   `json:"toAccountId"`
	CreditAmount            float64 `json:"CreditAmount"`
	BankChargeAmount        float64 `json:"BankChargeAmount"`
	MinCreditAmount         float64 `json:"minCreditAmount"`
	MaxCreditAmount         float64 `json:"maxCreditAmount"`
	AutoWithdrawCreditFlag  string  `json:"autoWithdrawCreditFlag"`
	AutoWithdrawConfirmFlag string  `json:"autoWithdrawConfirmFlag"`
}

type Member struct {
	Id            int64     `json:"id"`
	MemberCode    string    `json:"memberCode"`
	Username      string    `json:"username"`
	Phone         string    `json:"phone"`
	Firstname     string    `json:"firstname"`
	Lastname      string    `json:"lastname"`
	Fullname      string    `json:"fullname"`
	Credit        float64   `json:"credit"`
	Bankname      string    `json:"bankname"`
	BankAccount   string    `json:"bankAccount"`
	Promotion     string    `json:"promotion"`
	Status        string    `json:"status"`
	Channel       string    `json:"channel"`
	TrueWallet    string    `json:"trueWallet"`
	Note          string    `json:"note"`
	TurnoverLimit int       `json:"turnoverLimit"`
	CreatedAt     time.Time `json:"createdAt"`
}
type MemberListRequest struct {
	Search  string `form:"search" extensions:"x-order:1"`
	Page    int    `form:"page" extensions:"x-order:7" default:"1" min:"1"`
	Limit   int    `form:"limit" extensions:"x-order:8" default:"10" min:"1" max:"100"`
	SortCol string `form:"sortCol" extensions:"x-order:9"`
	SortAsc string `form:"sortAsc" extensions:"x-order:10"`
}
type MemberPossibleListRequest struct {
	UnknownStatementId int64   `form:"unknownStatementId" extensions:"x-order:1"`
	UserAccountNumber  *string `form:"userAccountNumber" extensions:"x-order:2"`
	UserBankCode       *string `form:"userBankCode" extensions:"x-order:3"`
	Page               int     `form:"page" extensions:"x-order:7" default:"1" min:"1"`
	Limit              int     `form:"limit" extensions:"x-order:8" default:"10" min:"1" max:"100"`
	SortCol            string  `form:"sortCol" extensions:"x-order:9"`
	SortAsc            string  `form:"sortAsc" extensions:"x-order:10"`
}
type MemberTransaction struct {
	Id                  int64          `json:"id" gorm:"primaryKey"`
	UserId              int64          `json:"userId"`
	MemberCode          string         `json:"memberCode"`
	UserUsername        string         `json:"userUsername"`
	UserFullname        string         `json:"userFullname"`
	TransferType        string         `json:"transferType"`
	PromotionId         int64          `json:"promotionId"`
	FromAccountId       int64          `json:"fromAccountId"`
	FromBankId          int64          `json:"fromBankId"`
	FromBankName        string         `json:"fromBankName"`
	FromAccountName     string         `json:"fromAccountName"`
	FromAccountNumber   string         `json:"fromAccountNumber"`
	ToAccountId         int64          `json:"toAccountId"`
	ToBankId            int64          `json:"toBankId"`
	ToBankName          string         `json:"toBankName"`
	ToAccountName       string         `json:"toAccountName"`
	ToAccountNumber     string         `json:"toAccountNumber"`
	CreditAmount        float64        `json:"creditAmount" sql:"type:decimal(14,2);"`
	PaidAmount          float64        `json:"paidAmount" sql:"type:decimal(14,2);"`
	DepositChannel      string         `json:"depositChannel"`
	OverAmount          float64        `json:"overAmount" sql:"type:decimal(14,2);"`
	BonusAmount         float64        `json:"bonusAmount" sql:"type:decimal(14,2);"`
	BonusReason         string         `json:"bonusReason"`
	BeforeAmount        float64        `json:"beforeAmount" sql:"type:decimal(14,2);"`
	AfterAmount         float64        `json:"afterAmount" sql:"type:decimal(14,2);"`
	BankChargeAmount    float64        `json:"bankChargeAmount" sql:"type:decimal(14,2);"`
	TransferAt          time.Time      `json:"transferAt"`
	CreatedByUserId     int64          `json:"createdByUserId"`
	CreatedByUsername   string         `json:"createdByUsername"`
	CancelRemark        string         `json:"cancelRemark"`
	CanceledAt          time.Time      `json:"canceledAt"`
	CanceledByUserId    int64          `json:"canceledByUserId"`
	CanceledByUsername  string         `json:"canceledByUsername"`
	ConfirmedAt         *time.Time     `json:"confirmedAt"`
	ConfirmedByUserId   int64          `json:"confirmedByUserId"`
	ConfirmedByUsername string         `json:"confirmedByUsername"`
	RemovedAt           time.Time      `json:"removedAt"`
	RemovedByUserId     int64          `json:"removedByUserId"`
	RemovedByUsername   string         `json:"removedByUsername"`
	Status              string         `json:"status"`
	StatusDetail        string         `json:"statusDetail"`
	IsAutoCredit        bool           `json:"isAutoCredit"`
	CreatedAt           time.Time      `json:"createAt"`
	UpdatedAt           *time.Time     `json:"updateAt"`
	DeletedAt           gorm.DeletedAt `json:"deleteAt"`
}
type MemberTransactionListRequest struct {
	UserId           string `form:"userId" extensions:"x-order:1"`
	FromTransferDate string `form:"fromTransferDate" extensions:"x-order:2"`
	ToTransferDate   string `form:"toTransferDate" extensions:"x-order:3"`
	TransferType     string `form:"transferType" extensions:"x-order:4"`
	Search           string `form:"search" extensions:"x-order:5"`
	Page             int    `form:"page" extensions:"x-order:6" default:"1" min:"1"`
	Limit            int    `form:"limit" extensions:"x-order:7" default:"10" min:"1" max:"100"`
	SortCol          string `form:"sortCol" extensions:"x-order:8"`
	SortAsc          string `form:"sortAsc" extensions:"x-order:9"`
}
type MemberTransactionSummary struct {
	TotalDepositAmount  float64 `json:"totalDepositAmount"`
	TotalWithdrawAmount float64 `json:"totalWithdrawAmount"`
	TotalBonusAmount    float64 `json:"totalBonusAmount"`
}

type MemberStatementType struct {
	Id        int64          `json:"id" gorm:"primaryKey"`
	Code      string         `json:"code"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"createAt"`
	UpdatedAt *time.Time     `json:"updateAt"`
	DeletedAt gorm.DeletedAt `json:"deleteAt"`
}
type MemberStatement struct {
	Id              int64          `json:"id" gorm:"primaryKey"`
	UserId          int64          `json:"userId"`
	StatementTypeId int64          `json:"statementTypeId"`
	TransferAt      time.Time      `json:"transferAt"`
	Info            string         `json:"info"`
	BeforeBalance   float64        `json:"beforeBalance" sql:"type:decimal(14,2);"`
	Amount          float64        `json:"amount" sql:"type:decimal(14,2);"`
	AfterBalance    float64        `json:"afterBalance" sql:"type:decimal(14,2);"`
	CreatedAt       time.Time      `json:"createAt"`
	UpdatedAt       *time.Time     `json:"updateAt"`
	DeletedAt       gorm.DeletedAt `json:"deleteAt"`
}
type MemberStatementCreateRequest struct {
	UserId int64   `json:"userId"`
	Amount float64 `json:"amount"`
}
type MemberStatementCreateBody struct {
	Id              int64 `json:"id"`
	UserId          int64 `json:"userId"`
	StatementTypeId int64 `json:"statementTypeId"`
	// TransferAt      time.Time `json:"transferAt"`
	Info string `json:"info"`
	// BeforeBalance   float64   `json:"beforeBalance" sql:"type:decimal(14,2);"`
	Amount float64 `json:"amount" sql:"type:decimal(14,2);"`
	// AfterBalance    float64   `json:"afterBalance" sql:"type:decimal(14,2);"`
}
type MemberStatementListRequest struct {
	UserId           string `form:"userId" extensions:"x-order:1"`
	StatementTypeId  string `form:"statementTypeId" extensions:"x-order:2"`
	FromTransferDate string `form:"fromTransferDate" extensions:"x-order:3"`
	ToTransferDate   string `form:"toTransferDate" extensions:"x-order:4"`
	Search           string `form:"search" extensions:"x-order:5"`
	Page             int    `form:"page" extensions:"x-order:6" default:"1" min:"1"`
	Limit            int    `form:"limit" extensions:"x-order:7" default:"10" min:"1" max:"100"`
	SortCol          string `form:"sortCol" extensions:"x-order:8"`
	SortAsc          string `form:"sortAsc" extensions:"x-order:9"`
}
type MemberStatementResponse struct {
	Id                int64          `json:"id" gorm:"primaryKey"`
	UserId            int64          `json:"userId"`
	MemberCode        string         `json:"memberCode"`
	UserUsername      string         `json:"userUsername"`
	UserFullname      string         `json:"userFullname"`
	StatementTypeId   int64          `json:"statementTypeId"`
	StatementTypeName string         `json:"statementTypeName"`
	TransferAt        time.Time      `json:"transferAt"`
	Info              string         `json:"info"`
	BeforeBalance     float64        `json:"beforeBalance" sql:"type:decimal(14,2);"`
	Amount            float64        `json:"amount" sql:"type:decimal(14,2);"`
	AfterBalance      float64        `json:"afterBalance" sql:"type:decimal(14,2);"`
	CreatedAt         time.Time      `json:"createAt"`
	UpdatedAt         *time.Time     `json:"updateAt"`
	DeletedAt         gorm.DeletedAt `json:"deleteAt"`
}
