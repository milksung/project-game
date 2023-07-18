package model

type Menu struct {
	Id      int64      `json:"id"`
	Title   string     `json:"title"`
	Name    string     `json:"name"`
	Managed bool       `json:"managed"`
	List    *[]SubMenu `json:"list"`
}

type SubMenu struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Name    string `json:"name"`
	Managed bool   `json:"managed"`
}
