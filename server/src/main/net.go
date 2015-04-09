package main

import (
	"net"
	"msg"
	"encoding/binary"
	"code.google.com/p/goprotobuf/proto"
	"crypt"
	"errors"
)

const MAX_CLI_BUF_SIZE = 2048

func NewConn(conn net.Conn) {
	defer conn.Close()
	user, err := RoleLogin(conn)
	if err != nil {
		return
	}

	// TODO reply rolesimple in this point, guid zero means need create player
	// ReplyRoleSimple(user)

	// one user got two goroutine, because it will get msg
	// from outside(network) and inside(framework)
	user.msg_chan = make(chan Msg)
	go theMsgMgr.UserMsgLoop(user)
	receive_loop(user)
}

func receive_loop(user *User) {
	conn := user.conn
	msg_chan := user.msg_chan
	seq_num := user.seq_num
	crypt_key := user.crypt_key

	buf := make([]byte, MAX_CLI_BUF_SIZE)
	cur_pos := 0

	RETURN:
	for {
		// conn could be closed by client or framefork
		n, err := conn.Read(buf[cur_pos:])
		if err != nil {
			logger.Println(err)
			break
		}

		logger.Printf("user[%d] receive %d bytes", user.GetId(), n)
		cur_pos += n
		// length(2)
		if cur_pos < 2 {
			continue
		}

		at := 0
		for {
			left_buf := buf[at:cur_pos]
			msg_len := int(binary.BigEndian.Uint16(left_buf[:2]))

			// package not complete
			if cur_pos - at < msg_len + 2 {
				break
			}

			// decrypt
			plaintext, err := crypt.AesDecrypt(left_buf[2:2+msg_len], crypt_key)
			if err != nil {
				logger.Printf("user[%d] message decrypt error[%s]", user.GetId(), err.Error())
				break RETURN
			}

			// seq_num(4) | type(4)
			msg_seqnum := binary.BigEndian.Uint32(plaintext[:4])
			msg_type := binary.BigEndian.Uint32(plaintext[4:8])
			logger.Printf("user: %d, msg_len: %d, msg_type: %d, seq_num: %d", user.GetId(), msg_len, msg_type, msg_seqnum)

			// package not in sequence
			seq_num += 1
			if msg_seqnum != seq_num {
				logger.Printf("user[%d] seq_num not equal server[%d], client[%d]", user.GetId(), seq_num, msg_seqnum)
				break RETURN
			}

			// TODO how to eliminate this alloc and copy?
			pack_len := len(plaintext) - 8
			pack := make([]byte, pack_len)
			copy(pack, plaintext[8:])
			new_msg := Msg{msg_type, pack}
			msg_chan <- new_msg

			at += msg_len + 2
			if at == cur_pos {
				break
			}
		}
		copy(buf, buf[at:cur_pos])
		cur_pos -= at

		if cur_pos >= MAX_CLI_BUF_SIZE {
			logger.Printf("user[%d] recv buf full", user.GetId())
			break RETURN
		}
	}

	msg_chan <- Msg{msg.FMQuit, nil}
}

func NetSend(conn net.Conn, seq_num uint32, msg_type uint32, msg proto.Message, crypt_key []byte) (int, error) {
	pack, err := proto.Marshal(msg)
	if err != nil {
		return 0, err
	}

	n := len(pack)
	buf := make([]byte, 8 + n)
	binary.BigEndian.PutUint32(buf[:4], seq_num)
	binary.BigEndian.PutUint32(buf[4:], msg_type)
	copy(buf[8:], pack)

	ciphertext := buf[:]
	if crypt_key != nil {
		ciphertext, err = crypt.AesEncrypt(buf, crypt_key)
		if err != nil {
			return 0, err
		}
	}

	send_buf := make([]byte, 2 + len(ciphertext))
	binary.BigEndian.PutUint16(send_buf, uint16(len(ciphertext)))
	copy(send_buf[2:], ciphertext)
	return conn.Write(send_buf)
}

func ReadPackage(conn net.Conn, buf []byte, msg proto.Message) error {
	n, err := conn.Read(buf)
	if err != nil {
		logger.Printf("%s: %s", conn.RemoteAddr(), err.Error())
		return err
	} else if n < 10 {
		logger.Printf("%s login req package not complete", conn.RemoteAddr())
		return errors.New("package not complete")
	}

	if err := proto.Unmarshal(buf[10:n], msg); err != nil {
		return err
	}
	return nil
}
