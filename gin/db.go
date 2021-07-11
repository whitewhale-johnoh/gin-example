package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func GetUrMyUser(loginId string, password string) bool {
	var err error
	user := User{}

	err = db.Get(&user, "SELECT loginid, password FROM urmyusers WHERE loginid=$1 AND password=$2", loginId, password)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func AddUrMyUser(loginId string, password string, name string) (bool, error) {
	rows, err := db.Query("INSERT INTO urmysaju (loginId) VALUES ($1)", loginId)
	if err != nil {
		//panic(err)
		return false, err
	}

	defer rows.Close()

	return true, nil
}

func Close() {
	db.Close()
}

const (
	host     = "192.168.10.150"
	port     = 5432
	user     = "koreaogh"
	password = "ogh1898"
	dbname   = "urmydb"
)

func ConnectToDB() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err = sqlx.Open("postgres", psqlInfo)
	if db != nil {
		db.SetMaxOpenConns(100) //최대 커넥션
		db.SetMaxIdleConns(10)  //대기 커넥션
	}
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
