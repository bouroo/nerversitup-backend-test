package models

type Result struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Errors     []Error     `json:"errors,omitempty"`
	ResultInfo *ResultInfo `json:"result_info,omitempty"`
}

type ResultInfo struct {
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	Count     int `json:"count"`
	TotalCont int `json:"total_count"`
}

type Error struct {
	Code    int    `json:"code"`
	Source  string `json:"source,omitempty"`
	Title   string `json:"title,omitempty"`
	Message string `json:"message,omitempty"`
}
