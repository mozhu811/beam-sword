package lark

type Fields struct {
	Event  string  `json:"事件"`
	Type   string  `json:"类型"`
	Amount float64 `json:"金额"`
	Tag    string  `json:"分类"`
	Date   int64   `json:"日期"`
}
type Record struct {
	Fields Fields `json:"fields"`
}

type BatchCreateReq struct {
	Records []Record `json:"records"`
}
