package main

import (
	"chat/db"
	"chat/handler"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	d, err := leveldb.OpenFile("chat_data", nil)
	if err != nil {
		log.Errorf("failed to open db, %+v", err)
		panic(err)
	}
	db.LvDB = d
	defer db.LvDB.Close()

	router := gin.Default()
	router.GET("/chat", handler.Chat)
	if err := router.Run(":9090"); err != nil {
		log.Errorf("failed to run server")
		panic(err)
	}
}
