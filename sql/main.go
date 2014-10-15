package main

import (
	"log"
	"srp"
	"crypto/sha256"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db_user = "root"
	db_pwd = ""
	db_name = "super_malio"
)

func main() {
	db, err := sql.Open("mysql", db_user + ":" + db_pwd + "@/" + db_name)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	username := []byte("xxxx")
	password := []byte("xxxx")

	srp, err := srp.NewSRP("openssl.1024", sha256.New, nil)
	if err != nil {
		log.Fatal(err)
	}

	srp.NewClientSession(username, password)
	salt, v, err := srp.ComputeVerifier(password)
	if err != nil {
		log.Fatal(err)
	}	

	log.Printf("salt(%d):%v", len(salt), salt)
	log.Printf("v:(%d)%v", len(v), v)

	_, err = db.Exec("insert into account(name, salt, vkey, guid) values(?, ?, ?, ?)", username, salt, v, 1)
	if err != nil {
		log.Fatal(err)
	}
}
