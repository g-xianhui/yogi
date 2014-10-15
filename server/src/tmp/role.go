package main

import (
	"net"
	"msg"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"code.google.com/p/goprotobuf/proto"
	"srp"
	"crypto/sha256"
	"math/rand"
	"errors"
)

type SLoadAccout struct {
	name string
}
type Account struct {
	guid uint32
	salt []byte
	vkey []byte
}

func (req *SLoadAccout) SqlMethod(db *sql.DB, reply chan<- DBResult) error {
	result := DBResult{nil, nil}	
	var account Account
	err := db.QueryRow("select guid, salt, vkey from account where name = ?", req.name).Scan(&account.guid, &account.salt, &account.vkey)

	result.err = err
	result.data = &account
	reply <- result
	return err
}

func LoadSimple(db *sql.DB, guid uint32) (*User, error) {
	logger.Printf("LoadSimple: %d", guid)
	var user *User
	var err error
	RETURN:
	for {
		rows, _err := db.Query("select name, level from role_simple where guid = ?", guid)
		if _err != nil {
			err = _err
			break RETURN
		}

		for rows.Next() {
			var name string
			var level uint32
			if err = rows.Scan(&name, &level); err != nil {
				break RETURN
			}

			user = &User{}
			user.role.Guid = proto.Uint32(guid)
			user.role.Name = proto.String(name)
			user.role.Level = proto.Uint32(level)
		}
		break
	}
	return user, err
}

func Login(conn net.Conn) (*User, error) {
	logger.Println("Login")
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		logger.Println(err.Error())
		return nil, err
	} else if n < 10 {
		logger.Println("login req package not complete")
		return nil, errors.New("package not complete")
	}
	
	req := &msg.MQLoginChallenge{}
	err = proto.Unmarshal(buf[10:n], req)
	if err != nil {
		logger.Println(err.Error())
		return nil, err
	}

	login_name := req.GetName()
	bytes_A := req.GetBytes_A()
	logger.Printf("login req: name[%s], bytes_A: %v", login_name, bytes_A)

	load_req := &SLoadAccout{ name : login_name }
	err, ret := db_query(load_req)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	account := ret.(*Account)

	// challenge reply
	auth, err := srp.NewSRP("openssl.1024", sha256.New, nil)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	logger.Printf("salt: %v", account.salt)
	ss := auth.NewServerSession([]byte(login_name), account.salt, account.vkey)
	challenge_rep := &msg.MRLoginChallenge{Salt : account.salt, Bytes_B : ss.GetB()}
	send(conn, 0, msg.MRLOGINCHALLENGE, challenge_rep, nil)

	// compute session key
	skey, err := ss.ComputeKey(bytes_A)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	logger.Printf("server computes key[%d]: %v", len(skey), skey)

	// final verify
	n, err = conn.Read(buf)
	if err != nil {
		logger.Println(err.Error())
		return nil, err
	} else if n < 10 {
		logger.Println("login req package not complete")
		return nil, errors.New("package not complete")
	}
	verify_req := &msg.MQLoginVerify{}
	err = proto.Unmarshal(buf[10:n], verify_req)
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	bytes_M := verify_req.GetBytes_M()
	if !ss.VerifyClientAuthenticator(bytes_M) {
		logger.Printf("verify failed")
		return nil, errors.New("verify failed")
	}

	// verify reply
	verify_rep := &msg.MRLoginVerify{Bytes_HAMK : ss.ComputeAuthenticator(bytes_M)}
	send(conn, 0, msg.MRLOGINVERIFY, verify_rep, nil)
	logger.Printf("verify success: %d", account.guid)

	// login result
	user, err := Load(account.guid)
	result_rep := &msg.MRLoginResult{}
	seq_num := rand.Uint32()
	if err != nil {
		result_rep.Result = proto.Int32(0)
	} else {
		result_rep.Result = proto.Int32(1)
	}
	send(conn, seq_num, msg.MRLOGINRESULT, result_rep, skey)

	if err != nil {
		return nil, err
	} else {
		user.crypt_key = skey
		user.conn = conn
		user.seq_num = seq_num
		return user, nil
	}
}

func Quit(user *User, msg interface{}) {
	logger.Printf("user[%d] quit", user.GetId())	
}

func Kick(user *User, msg interface{}) {
	logger.Printf("user[%d] kick", user.GetId())	
	// force ReceiveLoop to break
	user.conn.Close()
}

func init() {
	theMsgMgr.Register(msg.FMQuit, Quit)	
	theMsgMgr.Register(msg.FMKick, Kick)	
}