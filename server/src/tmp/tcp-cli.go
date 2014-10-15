package main

import "code.google.com/p/goprotobuf/proto"
import (
	"os"
	"fmt"
	"log"
	"net"
	"encoding/binary"
    "./msg"
    "./srp"
    "./crypt"
    "crypto/sha256"
)

var seq_num uint32
var send_buf []byte
var des_key []byte
func send(conn net.Conn, seq_num uint32, msg_type uint32, data []byte) {
	binary.BigEndian.PutUint16(send_buf, uint16(len(data)))
	binary.BigEndian.PutUint32(send_buf[2:], seq_num)
	binary.BigEndian.PutUint32(send_buf[6:], msg_type)
	copy(send_buf[10:], data)
	fmt.Printf("send %d bytes\n", 10 + len(data))
	conn.Write(send_buf[:10+len(data)])
}

func Send(conn net.Conn, seq_num uint32, msg_type uint32, data []byte) {
	binary.BigEndian.PutUint16(send_buf, uint16(len(data)))
	binary.BigEndian.PutUint32(send_buf[2:], seq_num)
	binary.BigEndian.PutUint32(send_buf[6:], msg_type)
	copy(send_buf[10:], data)

	ciphertext, err := crypt.AesEncrypt(send_buf[2:10+len(data)], des_key)
	if err != nil {
		log.Fatal(err)
	}
	copy(send_buf[2:], ciphertext)
	conn.Write(send_buf[:2+len(ciphertext)])
}

func receive(conn net.Conn) {
	cur_pos := 0
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf[cur_pos:])
		if err != nil {
			fmt.Println(err.Error())
			break
		}

		fmt.Printf("receive %d bytes\n", n)
		cur_pos += n
		// length(2) | seq_num(4) | type(4)
		if cur_pos < 10 {
			continue
		}

		msg_len := binary.BigEndian.Uint16(buf[:2])
		// package not complete
		if uint16(cur_pos) < msg_len + 10 {
			continue
		}

		cur_seqnum := binary.BigEndian.Uint32(buf[2:6])
		msg_type := binary.BigEndian.Uint32(buf[6:10])
		pack := buf[10:10+msg_len]
		cur_pos = cur_pos - int(10 + msg_len)
		buf = buf[10+msg_len:]

		if msg_type == msg.MRLOGIN {
			rep := &msg.MRLogin{}
			err := proto.Unmarshal(pack, rep)
			if err != nil {
				break
			}
			fmt.Printf("login result: %d, seq_num: %d\n", rep.GetResult(), cur_seqnum)
			seq_num = cur_seqnum
		} else {
			fmt.Printf("%s", pack)
		}
	}
}

func login(conn net.Conn) {
	fmt.Println("login")
	auth , err := srp.NewSRP("openssl.1024", sha256.New, nil)
	if err != nil {
		log.Fatal(err)
	}
	cs := auth.NewClientSession([]byte("agan"), []byte("9004"))

	bytes_A := cs.GetA()
	log.Printf("bytes_A: %v", bytes_A)
	// challenge
	log.Println("challenge")
	login_req := &msg.MQLogin{
		Name : proto.String("agan"),
		Bytes_A : bytes_A,
	}
	buffer, err := proto.Marshal(login_req)
	if err != nil {
		log.Fatal(err)
	}
	send(conn, 0, msg.MQLOGIN, buffer)

	// challenge reply
	log.Println("challenge reply")
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	} else if n < 10 {
		log.Println("login req package not complete")
		os.Exit(1)
	}
	
	challenge_rep := &msg.MRLoginChallenge{}
	err = proto.Unmarshal(buf[10:n], challenge_rep)
	if err != nil {
		log.Fatal(err)
	}

	salt := challenge_rep.GetSalt()
	bytes_B := challenge_rep.GetBytes_B()
	log.Printf("login salt: %v, bytes_B: %v", salt, bytes_B)

	// compute session key
	log.Println("compute session key")
	ckey, err := cs.ComputeKey(salt, bytes_B)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ckey[%d]: %v", len(ckey), ckey)

	// final verify
	log.Println("final verify")
	cauth := cs.ComputeAuthenticator()
	verify_req := &msg.MQLoginVerify{
		Bytes_M : cauth,
	}
	buffer, err = proto.Marshal(verify_req)
	if err != nil {
		log.Fatal(err)
	}
	send(conn, 0, msg.MQLOGINVERIFY, buffer)

	// verify reply
	log.Println("verify reply")
	n, err = conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	} else if n < 10 {
		log.Println("login req package not complete")
		os.Exit(1)
	}
	
	verify_rep := &msg.MRLoginVerify{}
	err = proto.Unmarshal(buf[10:n], verify_rep)
	if err != nil {
		log.Fatal(err)
	}

	bytes_HAMK := verify_rep.GetBytes_HAMK()
	if !cs.VerifyServerAuthenticator(bytes_HAMK) {
		log.Println("auth failed")
		os.Exit(1)
	}
	log.Println("verify success")

	des_key = ckey
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8088")
	if err != nil {
		fmt.Println("err:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	send_buf = make([]byte, 1024)

	login(conn)

	go receive(conn)
	buf := make([]byte, 1024)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			fmt.Println("os.Stdin.Read err:", err.Error())
			break
		}
		seq_num += 1
		send(conn, seq_num, 0, buf[:n])
	}
}
