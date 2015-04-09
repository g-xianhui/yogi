package main

import (
	"database/sql"
)

type SLoadRole struct {
	guid uint32
}
func (req *SLoadRole) SqlMethod(db *sql.DB, reply chan<- DBResult) error {
	logger.Printf("SLoadRole[%d]", req.guid)
	result := DBResult{nil, nil}
	user, err := LoadSimple(db, req.guid)
	if err != nil && user != nil {
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

const (
    user_find int = iota
    user_load
    user_add
)

type UserMgrResult struct {
    user *User
    err error
}

type UserMgrReq struct {
    cmd int
    id uint32
    user *User
    result chan *UserMgrResult
}

type UserMgr struct {
    users map[uint32]*User
    channel chan *UserMgrReq
}
var theUserMgr = &UserMgr{}

func load_user(id uint32) (*User, error) {
	load_req := &SLoadRole{ guid : id }
	result, err := db_query(load_req)
	if err != nil {
		return nil, err
	} else {
		return result.(*User), nil
	}
}

func (mgr *UserMgr) Run() {
    mgr.users = map[uint32]*User{}
    mgr.channel = make(chan *UserMgrReq)
    for cmd := range mgr.channel {
        switch cmd.cmd {
        case user_find:
            user, _ := mgr.users[cmd.id]
            cmd.result <- &UserMgrResult{user, nil}
        case user_load:
            user, err := load_user(cmd.id)
            cmd.result <- &UserMgrResult{user, err}
        case user_add:
            mgr.users[cmd.id] = cmd.user
        }
    }
}

func (mgr *UserMgr) Find(id uint32) *User {
    reply := make(chan *UserMgrResult)
    mgr.channel <- &UserMgrReq{cmd : user_find, id : id, result : reply}
    result := <- reply
	return result.user
}

func (mgr *UserMgr) Add(id uint32, user *User) {
    mgr.channel <- &UserMgrReq{cmd : user_add, id : id, user : user}
}

func (mgr *UserMgr) Load(id uint32) (*User, error) {
	user := mgr.Find(id)
	if user != nil {
		return user, nil
	}

    reply := make(chan *UserMgrResult)
    mgr.channel <- &UserMgrReq{cmd : user_load, id : id, result : reply}
    result := <- reply

	user, err := result.user, result.err
	if err != nil {
		return nil, err
	}

	// maybe not create player yet
	if user == nil {
		user = &User{}
	}
    mgr.Add(id, user)

	return user, nil
}
