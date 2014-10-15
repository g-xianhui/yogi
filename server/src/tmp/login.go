package main

import (
	"net"
	"msg"
	"code.google.com/p/goprotobuf/proto"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type SLoadRole struct {
	name string
}
func (msg *SLoadRole) SqlMethod(db *sql.DB, reply chan<- DBResult) error {
	result := DBResult{nil, nil}
	user, err := LoadSimple(db, msg.name)	
	if user != nil {
		LoadItem(db, user)
		for _, item := range(user.item_list) {
			logger.Printf("item id[%d], level[%d], data[%s]", item.GetId(), item.GetLevel(), item.GetData())
		}
	}
	result.err = err
	result.data = user
	reply <- result
	return err
}

func LoginProcess(conn net.Conn) *User {
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		logger.Println(err.Error())
		return nil
	} else if n < 10 {
		logger.Println("login req package not complete")
		return nil
	}
	
	req := &msg.MQLogin{}
	err = proto.Unmarshal(buf[10:n], req)
	if err != nil {
		logger.Println(err.Error())
		return nil
	}

	login_name := req.GetName()
	login_pwd := req.GetPwd()
	logger.Printf("login req: name[%s], pwd[%s]", login_name, login_pwd)

	load_req := &SLoadRole{ name : login_name }
	result := make(chan DBResult)
	db_req := &DBMsg{result, load_req}
	db_chan <- db_req
	v := <- result
	if v.err != nil || v.data == nil {
		return nil
	} else {
		return v.data.(*User)
	}
}