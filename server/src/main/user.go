package main

import (
    "net"
    "msg"
	"code.google.com/p/goprotobuf/proto"
)

type User struct {
	seq_num uint32
	crypt_key []byte
	conn net.Conn
	msg_chan chan Msg
	role msg.MRoleSimple
	item_list []*msg.MItem
}

func (user *User) GetId() uint32 {
	return user.role.GetGuid()
}

func (user *User) Send(msg_type uint32, msg proto.Message) (int, error) {
	return NetSend(user.conn, user.seq_num, msg_type, msg, user.crypt_key)
}

