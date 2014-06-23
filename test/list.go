package main

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
        "encoding/json"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func getMongoSession() (mongoSession *mgo.Session) {
	mongoSession, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	mongoSession.SetMode(mgo.Monotonic, true)
	return mongoSession
}

type Worker struct {
	Port      string
	Host      string
	Container string
}

type Deployment struct {
	Workers []Worker
	App     string
	Version string
}

type Version struct {
	Version    string
	Filename   string
	Created_at time.Time
}

type App struct {
	Id          bson.ObjectId `bson:"_id"`
	Name        string
	Description string
	Path        string
	Workers     uint32
	Versions    []Version
	Deployment  Deployment
}

func main() {
	fmt.Printf("Hello\n")

	mongoSession := getMongoSession()
	appsCollection := mongoSession.DB("control").C("apps")
	//var apps []map[string]interface{}
	var apps []App
	err := appsCollection.Find(nil).All(&apps)
	check(err)
<<<<<<< Local Changes
<<<<<<< Local Changes
	fmt.Printf("%+v", apps)
	
	
	
=======
	fmt.Printf("%+v\n\n", apps)

        b, err := json.Marshal(apps)
        check(err)
        fmt.Printf("%v\n\n", string(b))

        var moreapps []App
        err = json.Unmarshal(b, &moreapps)
        check(err)
	fmt.Printf("%+v\n\n", moreapps)
 

>>>>>>> External Changes
=======
	fmt.Printf("%+v\n\n", apps)

        b, err := json.Marshal(apps)
        check(err)
        fmt.Printf("%v\n\n", string(b))

        var moreapps []App
        err = json.Unmarshal(b, &moreapps)
        check(err)
	fmt.Printf("%+v\n\n", moreapps)
 

>>>>>>> External Changes
}
