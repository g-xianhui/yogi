package main
import (
	"strconv"
	"time"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

const (
	find int = iota
	add
	update
	remove
)

type sysMsg struct {
	cmd int
	key string
	value string
	result chan string
}
type SysChan chan *sysMsg

type SSysQInfoLoad struct {
}
func (msg *SSysQInfoLoad)SqlMethod(db *sql.DB, reply chan<- DBResult) error {
	result := DBResult{nil, map[string]string{}}
	RETURN:
	for {
		rows, err := db.Query("select `key`, `value` from server_info")
		if err != nil {
			result.err = err
			break RETURN
		}

		m := map[string]string{}
		for rows.Next() {
			var key, value string
			if result.err = rows.Scan(&key, &value); result.err != nil {
				break RETURN
			}
			m[key] = value
		}
		result.data = m
		result.err = rows.Err()
		break
	}
	reply <- result
	return result.err
}

type SSysQInfoSave struct {
	info map[string]string
}
func (msg *SSysQInfoSave)SqlMethod(db *sql.DB, reply chan<- DBResult) error {
	if _, err := db.Exec("delete from server_info"); err != nil {
		return err
	}
	for k, v := range msg.info {
		if _, err := db.Exec("insert into server_info(`key`, `value`) values (?, ?)", k, v); err != nil {
			return err
		}
	}
	logger.Println("sysinfo save end")
	return nil
}

func (sys SysChan) Find(key string) string {
	logger.Printf("sys.Find: %s", key)
	reply := make(chan string)
	sys <- &sysMsg{cmd : find, key : key, result : reply}
	result := <- reply
	return result
}

func (sys SysChan) Add(key string, value string) {
	sys <- &sysMsg{cmd : add, key : key, value : value}
}

func (sys SysChan) Update(key string, value string) {
	sys <- &sysMsg{cmd : update, key : key, value : value}
}

func (sys SysChan) Run() {
	logger.Println("sys.Run")
	data := map[string]string{}
	sys_init(data)
	for k, v := range data {
		logger.Printf("data[k]:%s, [v]:%s", k, v)
	}
	for cmd := range sys {
		switch cmd.cmd {
		case find:
			value, _ := data[cmd.key]
			cmd.result <- value
		case add:
			data[cmd.key] = cmd.value
		case remove:
			delete(data, cmd.key)
		case update:
			data[cmd.key] = cmd.value
		}
	}
}

func sys_init(data map[string]string) {
	logger.Printf("load_sys_info")
	load_req := &SSysQInfoLoad{}

    ret, err := db_query(load_req)
    if err != nil {
        logger.Fatal("load sys info failed!")
    }

	rep := ret.(map[string]string)
	for k, v := range rep {
		logger.Printf("%s:%s", k, v)
		data[k] = v
	}

	now := time.Now().Unix()
	logger.Printf("time now:%d", now)
	data["open_time"] = strconv.FormatInt(now, 10)

	set_req := &SSysQInfoSave{data}
    db_exec(set_req)
}
