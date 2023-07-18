package model

type AGCRegister struct {
	Username         string `json:"Username" validate:"required"`
	Agentname        string `json:"Agentname" validate:"required"`
	Fullname         string `json:"Fullname" validate:"required"`
	Password         string `json:"Password" validate:"required"`
	Currency         string `json:"Currency" validate:"required"`
	Dob              string `json:"Dob"`
	Gender           int    `json:"Gender" validate:"required"`
	Email            string `json:"Email"`
	Mobile           string `json:"Mobile" validate:"required"`
	Ip               string `json:"Ip" validate:"required"`
	Timestamp        int64  `json:"Timestamp" validate:"required"`
	Sign             string `json:"Sign" validate:"required"`
	CommFollowUpline int    `json:"CommFollowUpline"`
	PTFollowUpline   int    `json:"PTFollowUpline"`
}

type AGCLogin struct {
	Username  string `json:"Username" validate:"required"`
	Partner   string `json:"Partner" validate:"required"`
	Timestamp int64  `json:"Timestamp" validate:"required"`
	Sign      string `json:"Sign" validate:"required"`
	Domain    string `json:"Domain" validate:"required"`
	Lang      string `json:"Lang" validate:"required"`
	IsMobile  bool   `json:"IsMobile" validate:"required"`
	Ip        string `json:"Ip" validate:"required"`
}

type AGCChangePassword struct {
	PlayerName  string `json:"PlayerName" validate:"required"`
	Partner     string `json:"Partner" validate:"required"`
	NewPassword string `json:"NewPassword" validate:"required"`
	Timestamp   int64  `json:"Timestamp" validate:"required"`
	Sign        string `json:"Sign" validate:"required"`
}

type AGCDeposit struct {
	Agentname     string  `json:"Agentname" validate:"required"`
	PlayerName    string  `json:"PlayerName" validate:"required"`
	Amount        float64 `json:"Amount" validate:"required"`
	Timestamp     int64   `json:"TimeStamp" validate:"required"`
	Sign          string  `json:"Sign" validate:"required"`
	TransactionId string  `json:"TransactionId" validate:"required"`
}

type AGCWithdraw struct {
	Agentname     string  `json:"Agentname" validate:"required"`
	PlayerName    string  `json:"PlayerName" validate:"required"`
	Amount        float64 `json:"Amount" validate:"required"`
	TimeStamp     int64   `json:"TimeStamp" validate:"required"`
	Sign          string  `json:"Sign" validate:"required"`
	TransactionId string  `json:"TransactionId" validate:"required"`
}
