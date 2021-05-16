package handler

import (
	"bytes"
	"chat/data"
	"chat/db"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var connections = map[string]*websocket.Conn{}

func Chat(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorf("failed to upgrade, %+v", err)
		c.String(http.StatusInternalServerError, "failed to upgrade websocket")
		return
	}

	userName := c.Request.Header.Get("userName")
	if userName == "" {
		log.Errorf("failed to get username")
		c.String(http.StatusInternalServerError, "failed to get username")
		return
	}

	connections[userName] = conn

	conn.SetCloseHandler(func(code int, text string) error {
		delete(connections, userName)
		return nil
	})

	// 读取历史数据
	dataList, err := db.ReadChatData(userName)
	if err != nil {
		log.Errorf("failed to read history data, %+v", err)
		c.String(http.StatusInternalServerError, "failed to read history data")
		conn.Close()
		return
	}

	// 发送历史数据
	for _, chatData := range dataList {
		b, err := json.Marshal(chatData)
		if err != nil {
			log.Errorf("failed to marshal char data, %+v", err)
			c.String(http.StatusInternalServerError, "failed to marshal char data")
			conn.Close()
			return
		}
		if err := conn.WriteMessage(websocket.TextMessage, b); err != nil {
			log.Errorf("failed to write history data, %+v", err)
			c.String(http.StatusInternalServerError, "failed to write history data")
			conn.Close()
			return
		}
	}

	go func() {
		for {
			messageType, payload, err := conn.ReadMessage()
			if err != nil {
				log.Errorf("failed to read message, %+v", err)
				conn.Close()
				break
			}
			if messageType != websocket.TextMessage {
				continue
			}
			chatData := &data.ChatData{}
			if err := json.NewDecoder(bytes.NewReader(payload)).Decode(chatData); err != nil {
				log.Errorf("failed to decode json")
				conn.Close()
				break
			}
			log.Infof("succeed to receive chat data: %+v", chatData)
			toConn, ok := connections[chatData.ToUsername]
			if !ok {
				// 不在线
				log.Warnf("toUsername is not online")
				if err := db.SaveChatData(chatData); err != nil {
					log.Errorf("failed to save chat data")
				}
			} else {
				// 在线
				if err := toConn.WriteMessage(websocket.TextMessage, payload); err != nil {
					log.Errorf("failed to write message")
					toConn.Close()
				}
			}
		}
	}()
}
