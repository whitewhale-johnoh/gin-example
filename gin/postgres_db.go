package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Userinfo struct {
	LoginUuid      string `db:"loginuuid"`
	AccessToken    string
	Email          string `db:"email"`
	PhoneNo        string `db:"phoneno"`
	PhoneCode      string `db:"phonecode"`
	Name           string `db:"name"`
	Birthdate      string `db:"birthday"`
	Hometown       string `db:"hometown"`
	Country        string `db:"country"`
	Gender         bool   `db:"gender"`
	ProfilePicPath string
}

var db *sqlx.DB

func GetUrMyUserEmail(signinuser *SignInUser) (string, error) {

	var urmyuuid string
	//user := SignInUser{}

	emailerr := db.Get(&urmyuuid, "SELECT loginuuid FROM urmyuserinfo WHERE email=$1 AND password=$2", signinuser.LoginID, signinuser.Password)
	if emailerr != nil {
		fmt.Println("emailerr")
		fmt.Println(emailerr.Error())
		return "emailerr", emailerr
	}

	rows, lastloginerr := db.Query("UPDATE urmyusers SET lastlogin=current_timestamp WHERE loginuuid=$1", urmyuuid)
	if lastloginerr != nil {
		fmt.Println("lastloginerr")
		fmt.Println(lastloginerr.Error())

		return "lastloginerr", lastloginerr
	}
	defer rows.Close()
	/*
		if user.Token != signinuser.Token {
			rows, err := db.Query("UPDATE urmyusers SET notificationtoken=$1 WHERE loginuuid=$2", signinuser.Token, urmyuuid)
			if err != nil {
				fmt.Println(err)
				return "", err
			}
			defer rows.Close()
		}
	*/

	return urmyuuid, nil
}

func GetUrMyUserPhone(signinuser *SignInUser) (string, error) {
	var err error
	var urmyuuid string
	user := SignInUser{}

	phoneerr := db.Get(&urmyuuid, "SELECT loginuuid FROM urmyuserinfo WHERE phoneno=$1 AND password=$2", signinuser.LoginID, signinuser.Password)
	if phoneerr != nil {
		return "phoneerr", phoneerr
	}

	rows, err := db.Query("UPDATE urmyusers SET lastlogin=current_timestamp WHERE loginuuid=$1", urmyuuid)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer rows.Close()

	if user.Token != signinuser.Token {
		rows, err := db.Query("UPDATE urmyusers SET notificationtoken=$1 WHERE loginuuid=$2", signinuser.Token, urmyuuid)
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		defer rows.Close()
	}

	return urmyuuid, nil
}

func GetUrMyUserData(uuid string) (Userinfo, error) {
	var userinfo Userinfo
	userdataerr := db.Get(&userinfo, "SELECT * FROM urmyuserinfo WHERE loginuuid=$1", uuid)
	if userdataerr != nil {
		return userinfo, userdataerr
	}
	return userinfo, nil
}

func AddUrMyUser(signupuser *SignUpUser) (string, error) {
	birthdateinsert := pq.QuoteLiteral(signupuser.Birthdate)
	var urmyuuid string
	var checkUUIDExist string
	for {
		uuid, generr := generateuuid()
		if generr != nil {
			fmt.Println("generr error")
			fmt.Println(generr.Error())
			return "", generr
		}
		sameuuiderr := db.Get(&checkUUIDExist, "SELECT loginuuid FROM urmyusers WHERE loginuuid=$1", uuid)
		if sameuuiderr.Error() == "sql: no rows in result set" {
			urmyuuid = uuid
			break
		} else {
			fmt.Println("sameuuiderr error")
			fmt.Println(sameuuiderr.Error())
		}
	}

	inserturmyuser, inserturmyuseriderr := db.Query("INSERT INTO urmyusers (loginuuid) VALUES ($1)", urmyuuid)
	if inserturmyuseriderr != nil {
		fmt.Println("inserturmyuseriderr error")
		fmt.Println(inserturmyuseriderr.Error())
		return "", inserturmyuseriderr
	}
	defer inserturmyuser.Close()

	inserturmybasicinfo, inserturmybasicinfoerr := db.Query("INSERT INTO urmyuserinfo (loginuuid, email, password, phoneno, phonecode, name, hometown, country, birthday, gender) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		urmyuuid, signupuser.LoginID, signupuser.Password, signupuser.PhoneNo, signupuser.Phonecode, signupuser.Name, signupuser.Hometown, signupuser.Country, birthdateinsert, signupuser.Gender)
	if inserturmybasicinfoerr != nil {
		fmt.Println("inserturmybasicinfoerr error")
		fmt.Println(inserturmybasicinfoerr.Error())
		return "", inserturmybasicinfoerr
	}
	defer inserturmybasicinfo.Close()

	insertagreement, insertagreementerr := db.Query("INSERT INTO urmyusersagreement (loginuuid, isoverage, urmyaccount, urmyoverallservice, urmynotiad, urmypersonaldataacc, urmylocation, urmyprofileadditional ,createdat) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, current_timestamp)",
		urmyuuid, signupuser.Isoverage, signupuser.Urmyaccount, signupuser.Urmyoverallservice, signupuser.Urmynotiad, signupuser.Urmypersonaldataacc, signupuser.Urmylocation, signupuser.Urmyprofileadditional)
	if insertagreementerr != nil {
		fmt.Println("insertagreementerr error")
		fmt.Println(insertagreementerr.Error())
		return "", insertagreementerr
	}
	defer insertagreement.Close()

	return urmyuuid, nil
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
