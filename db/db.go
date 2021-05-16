package db

import (
	"bytes"
	"chat/data"
	"encoding/json"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	LvDB *leveldb.DB
)

func SaveChatData(chatData *data.ChatData) error {
	v, err := json.Marshal(chatData)
	if err != nil {
		return err
	}
	k := fmt.Sprintf("%s_%d", chatData.ToUsername, chatData.SendTime)

	if err := LvDB.Put([]byte(k), v, nil); err != nil {
		return err
	}
	return nil
}

func ReadChatData(username string) ([]*data.ChatData, error) {
	kPrefix := fmt.Sprintf("%s_", username)
	iter := LvDB.NewIterator(util.BytesPrefix([]byte(kPrefix)), nil)

	dataList := make([]*data.ChatData, 0)
	for iter.Next() {
		b := iter.Value()
		chatData := &data.ChatData{}
		if err := json.NewDecoder(bytes.NewReader(b)).Decode(chatData); err != nil {
			return nil, err
		}
		if chatData.ToUsername != username {
			break
		}
		if chatData.IsRead {
			continue
		}
		LvDB.Delete(iter.Key(), nil)
		dataList = append(dataList, chatData)
	}
	return dataList, nil
}
