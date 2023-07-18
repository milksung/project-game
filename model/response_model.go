package model

type Success struct {
	Message string `json:"message"`
}

type SuccessWithData struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type SuccessWithList struct {
	Message string      `json:"message"`
	List    interface{} `json:"list"`
}

type SuccessWithPagination struct {
	Message string      `json:"message" validate:"required,min=1,max=255"`
	List    interface{} `json:"list"`
	Total   int64       `json:"total"`
}

type SuccessWithToken struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}
