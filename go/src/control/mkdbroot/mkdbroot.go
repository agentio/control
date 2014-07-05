package main

import (
	"flag"
	"fmt"
	"labix.org/v2/mgo"
)

func getMongoSession() (mongoSession *mgo.Session) {
	mongoSession, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	return mongoSession
}

var username = flag.String("u", "root", "username for root database user")
var password = flag.String("p", "", "password for root database user")

func main() {
        flag.Parse()
	if len(*password) > 0 {
		mongoSession := getMongoSession()
		db := mongoSession.DB("admin")
		user := mgo.User{
			Username: *username,
			Password: *password,
			Roles:    []mgo.Role{"root"},
		}
		err := db.UpsertUser(&user)
		fmt.Printf("result: %+v\n", err)
	} else {
		fmt.Printf("Please specify a password for the root user with -p\n")
	}
}
