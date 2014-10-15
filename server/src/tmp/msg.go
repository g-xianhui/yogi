package msg

import (
	"net"
	"log"
)

var logger *log.Logger

type MsgHandler func (net.Conn, interface{}, []byte)
var msg_handlers = map[uint32]MsgHandler{}

func SetLogger(l *log.Logger) {
	logger = l
}

func Register(msg_type uint32, hanlder MsgHandler) {
	msg_handlers[msg_type] = hanlder
}

func Dispatch(conn net.Conn, target interface{}, msg_type uint32, pack []byte) {
	hanlder, ok := msg_handlers[msg_type]
	if ok != true {
		logger.Printf("unknow msg type[%d]", msg_type)
	} else {
		hanlder(conn, target, pack)
	}
}
