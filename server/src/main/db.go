package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db_user = "root"
	db_pwd = ""
	db_name = "super_malio"
)

type DBResult struct {
	err error
	data interface{}
}

type SqlMsg interface {
	SqlMethod(db *sql.DB, reply chan<- DBResult) error
}

type DBMsg struct {
	reply chan DBResult
	msg SqlMsg
}

var db_chan chan *DBMsg

func db_thread() {
    db_chan = make(chan *DBMsg)
	db, err := sql.Open("mysql", db_user + ":" + db_pwd + "@/" + db_name)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer db.Close()

	for v := range db_chan {
		req := v.msg
		if err := req.SqlMethod(db, v.reply); err != nil {
			logger.Println(err.Error())
		}
	}
}

func db_query(req SqlMsg) (interface{}, error) {
	reply := make(chan DBResult)
	db_chan <- &DBMsg{reply, req}
	v := <- reply
	return v.data, v.err
}

func db_exec(req SqlMsg) {
    db_chan <- &DBMsg{nil, req}
}
