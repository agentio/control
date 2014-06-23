package main

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"	
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func md5HashWithSalt(input, salt string) string {
	hasher := hmac.New(md5.New, []byte(salt))
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

func getMongoSession() (mongoSession *mgo.Session) {
	mongoSession, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	mongoSession.SetMode(mgo.Monotonic, true)
	return mongoSession
}

func main() {
	username := "tim"
	password := "123"
	saltedPassword := md5HashWithSalt(password, "agent.io")
	mongoSession := getMongoSession()
	usersCollection := mongoSession.DB("accounts").C("users")
	usersCollection.Insert(bson.M{"username": username, "password": saltedPassword})
}
