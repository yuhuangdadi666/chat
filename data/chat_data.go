package data

type ChatData struct {
	ToUsername string // 发给谁
	Content    string // 内容
	IsRead     bool   //是否已读
	SendTime   int64
}
