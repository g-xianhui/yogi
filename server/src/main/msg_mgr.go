package main

import (
    "msg"
)

type MsgHandler func (*User, interface{})
//var msg_handlers = map[uint32]MsgHandler{}

type Msg struct {
	msg_type uint32
	data interface{}
}

type MsgMgr struct {
	msg_handlers map[uint32]MsgHandler
}

var theMsgMgr = &MsgMgr{msg_handlers : map[uint32]MsgHandler{}}

func (mgr *MsgMgr) Register(msg_type uint32, hanlder MsgHandler) {
	mgr.msg_handlers[msg_type] = hanlder
}

func (mgr *MsgMgr) Dispatch(user *User, msg_type uint32, msg interface{}) {
	hanlder, ok := mgr.msg_handlers[msg_type]
	if ok != true {
		logger.Printf("unknow msg type[%d]", msg_type)
	} else {
		hanlder(user, msg)
	}
}

func (mgr *MsgMgr) UserMsgLoop(user *User) {
	for {
		new_msg := <- user.msg_chan
		if new_msg.msg_type == msg.FMQuit {
			RoleQuit(user)
			break
		}
		mgr.Dispatch(user, new_msg.msg_type, new_msg.data)
	}
}
