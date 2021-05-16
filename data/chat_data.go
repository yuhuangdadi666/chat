package data

type ChatData struct {
	ToUsername string `json:"to_username"` // 发给谁
	Content    string `json:"content"` // 内容
	IsRead     bool `json:"is_read"`   //是否已读
	SendTime   int64 `json:"send_time"`
}
