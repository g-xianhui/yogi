package main

import (
	"code.google.com/p/goprotobuf/proto"
	"msg"
)

func EchoHandler(user *User, data interface{}) {
	logger.Println("EchoHandler")
	m := data.([]byte)
	req := &msg.MQEcho{}
	if err := proto.Unmarshal(m, req); err != nil {
		logger.Println(err)
		return
	}

	logger.Printf("user: %d, Echo: %s", user.role.GetGuid(), req.GetData())
	rep := &msg.MREcho{}
	rep.Data = proto.String(req.GetData())
	user.Send(msg.MRECHO, rep)
}

func init() {
	theMsgMgr.Register(msg.MQECHO, EchoHandler)
}