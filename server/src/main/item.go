package main

import (
	"msg"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"code.google.com/p/goprotobuf/proto"
)

func LoadItem(db *sql.DB, user *User) {
	rows, err := db.Query("select item_id, level, data from item where role_id = ?", user.role.GetGuid())
	if err != nil {
		return
	}
	for rows.Next() {
		var id, level uint32
		var data string
		if err := rows.Scan(&id, &level, &data); err != nil {
			return
		}
		item := &msg.MItem{}
		item.Id = proto.Uint32(id)
		item.Level = proto.Uint32(level)
		item.Data = proto.String(data)
		user.item_list = append(user.item_list, item)
	}
}

func item_reply_list(user *User, req interface{}) {
	rep := &msg.MRItemList{}
	rep.ItemList = user.item_list
	user.Send(msg.MRITEMLIST, rep)
}

func init() {
    theMsgMgr.Register(msg.MQITEMLIST, item_reply_list)
}
