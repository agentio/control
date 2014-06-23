package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

//const SERVICE = "http://localhost:2010"
//const CREDENTIALS = "tim:123"

const SERVICE = "http://beta.agent.io"
const CREDENTIALS = "tim:123"

type User struct {
	Id       string
	Username string
	Password string
}

type Version struct {
	Version    string
	Filename   string
	Created_at time.Time
}

type Worker struct {
	Port      uint32
	Host      string
	Container string
	Version   string
}

type App struct {
	Id          string
	Name        string
	Description string
	Path        string
	Domains     string
	Capacity    uint32
	Versions    []Version
	Workers     []Worker
}

func check(err error) {
	if err != nil {
		fmt.Printf("Error! %v\n", err)
		panic(err)
	}
}

func performRequest(result *map[string]interface{}, req *http.Request) {
	client := &http.Client{}
	encoding := base64.StdEncoding.EncodeToString([]byte(CREDENTIALS))
	authorization := fmt.Sprintf("Basic %v", encoding)
	req.Header.Add("Authorization", authorization)
	fmt.Printf("REQUESTING...\n")
	response, err := client.Do(req)
	check(err)
	buffer, err := ioutil.ReadAll(response.Body)
	check(err)
	defer response.Body.Close()
	if response.StatusCode != 200 {
		fmt.Printf("status code %v\n", response.StatusCode)
		fmt.Println(string(buffer))
		fmt.Println("Done")
		panic("goodbye")
	}
	fmt.Printf("RESPONSE %v\n", string(buffer))
	err = json.Unmarshal(buffer, result)
	fmt.Printf("UNMARSHALED\n")
	check(err)
}

func performRequestIntoBuffer(req *http.Request) []byte {
	client := &http.Client{}
	encoding := base64.StdEncoding.EncodeToString([]byte(CREDENTIALS))
	authorization := fmt.Sprintf("Basic %v", encoding)
	req.Header.Add("Authorization", authorization)
	fmt.Printf("REQUESTING...\n")
	response, err := client.Do(req)
	check(err)
	buffer, err := ioutil.ReadAll(response.Body)
	check(err)
	defer response.Body.Close()
	if response.StatusCode != 200 {
		fmt.Printf("status code %v\n", response.StatusCode)
		fmt.Println(string(buffer))
		fmt.Println("Done")
		panic("goodbye")
	}
	fmt.Printf("RESPONSE %v\n", string(buffer))
	return buffer
}

func getApps(result *[]App) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/control/apps", SERVICE), nil)
	check(err)
	buffer := performRequestIntoBuffer(req)
	err = json.Unmarshal(buffer, result)
	check(err)
}

func createApp(result *map[string]interface{}, app map[string]interface{}) {
	jsonData, err := json.Marshal(app)
	post_data := bytes.NewReader(jsonData)
	req, err := http.NewRequest("POST", fmt.Sprintf("%v/control/apps", SERVICE), post_data)
	check(err)
	req.Header.Add("Content-Type", "application/json")
	performRequest(result, req)
}

func deleteApps(result *map[string]interface{}) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%v/control/apps", SERVICE), nil)
	check(err)
	performRequest(result, req)
}

func getApp(result *App, appid string) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/control/apps/%v", SERVICE, appid), nil)
	check(err)
	buffer := performRequestIntoBuffer(req)
	err = json.Unmarshal(buffer, result)
	check(err)
}

func createAppVersion(result *map[string]interface{}, appid string, version []byte) {
	post_data := bytes.NewReader(version)
	req, err := http.NewRequest("POST", fmt.Sprintf("%v/control/apps/%v/versions", SERVICE, appid), post_data)
	check(err)
	performRequest(result, req)
}

func deployAppVersion(result *map[string]interface{}, appid string, versionid string) {
	data := url.Values{}
	data.Set("command", "start")
	req, err := http.NewRequest("POST", fmt.Sprintf("%v/control/apps/%v/versions/%v", SERVICE, appid, versionid),
		bytes.NewBufferString(data.Encode()))
	check(err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	performRequest(result, req)
}

func stopAppVersion(result *map[string]interface{}, appid string, versionid string) {
	data := url.Values{}
	data.Set("command", "stop")
	req, err := http.NewRequest("POST", fmt.Sprintf("%v/control/apps/%v/versions/%v", SERVICE, appid, versionid),
		bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	check(err)
	performRequest(result, req)
}

func stopApp(result *map[string]interface{}, appid string) {
	data := url.Values{}
	data.Set("command", "stop")
	req, err := http.NewRequest("POST", fmt.Sprintf("%v/control/apps/%v", SERVICE, appid),
		bytes.NewBufferString(data.Encode()))
	check(err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	performRequest(result, req)
}

func main() {

	if false {
		{
			// delete current apps
			var result map[string]interface{}
			deleteApps(&result)
			fmt.Printf("DELETING %+v\n", result)
		}

		{
			// verify that there are no apps
			var apps []App
			getApps(&apps)
			fmt.Printf("GETTING %+v\n", apps)
		}

		newapps := []map[string]interface{}{
			{
				"name":        "sample",
				"path":        "sample",
				"capacity":    2,
				"description": "A sample app.",
			},
		}

		for _, newapp := range newapps {
			var result map[string]interface{}
			createApp(&result, newapp)
			fmt.Printf("CREATING %+v\n", result)
		}

		var apps []App
		getApps(&apps)
		fmt.Printf("GETTING %+v\n", apps)

		for _, app := range apps {
			fmt.Printf("=========\n")
			fmt.Printf("APP IN RANGE %+v\n", app)

			appid := app.Id

			fmt.Printf("APP ID %+v\n", appid)

			var app2 App
			getApp(&app2, appid)
			fmt.Printf("RECEIVED APP %+v\n", app2)

			data, _ := ioutil.ReadFile("sample.zip")
			fmt.Printf("%v\n", len(data))

			{
				fmt.Printf("CREATE VERSION\n")

				var result map[string]interface{}
				createAppVersion(&result, appid, data)
			}

			{
				fmt.Printf("DEPLOY VERSION\n")

				var app App
				getApp(&app, appid)
				fmt.Printf("%+v\n", app)

				versions := app.Versions
				fmt.Printf("versioninfo %+v\n", versions)

				version := versions[0]
				versionid := version.Version

				fmt.Printf("VERSION ID: %+v\n", versionid)

				var result map[string]interface{}
				deployAppVersion(&result, appid, versionid)

				fmt.Printf("DEPLOY COMPLETE\n")
			}
		}
	}

	if true {
		// verify that there are no apps
		var apps []App
		getApps(&apps)
		fmt.Printf("GETTING %+v\n", apps)

		if false {
			for _, info := range apps {
				app := info
				appid := app.Id

				fmt.Printf("STOPPING %+v\n", appid)
				{
					var result map[string]interface{}
					stopApp(&result, appid)
				}
			}
		}
	}
}
