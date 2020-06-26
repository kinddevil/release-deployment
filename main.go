package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type ReleaseInfo struct {
	Assets []Asset `json:"assets"`
}

type Asset struct {
	ID int64 `json:"id"`
}

type Header struct {
	Name  string
	Value string
}

const (
	authHeaderName  = "Authorization"
	authHeader      = "token 1aef08e73e9c535db356ba5126ab10aae1880af6"
	acceptHaderName = "Accept"
	acceptHader     = "application/octet-stream"
	apiPrefix       = "https://api.github.com/repos/zhizhuoxingyi/intelligence_campus_back/releases/"
	releaseInfoPath = "latest"
	assetPath       = "assets/"
)

func doRequest(method, url string, headers ...*Header) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	for _, header := range headers {
		req.Header.Add(header.Name, header.Value)
	}

	return client.Do(req)
}

func doGet(url string, headers ...*Header) (*http.Response, error) {
	return doRequest("GET", url, headers...)
}

func main() {

	if _, err := os.Stat("lock"); err == nil {
		log.Println("running deployment, ignore...")
		os.Exit(0)
	}

	url := apiPrefix + releaseInfoPath
	log.Printf("Get latest release via %s\n", url)
	res, err := doGet(url, &Header{authHeaderName, authHeader})
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// responseData, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// responseString := string(responseData)

	info := &ReleaseInfo{}
	err = json.NewDecoder(res.Body).Decode(info)
	if err != nil {
		log.Fatal(err)
	}

	assetIDStr := strconv.FormatInt(info.Assets[0].ID, 10)
	log.Printf("Get latest asset id %s\n", assetIDStr)

	versionPath := "release_version"
	if _, err := os.Stat(versionPath); err == nil {
		data, err := ioutil.ReadFile(versionPath)
		if err != nil {
			log.Fatal(err)
		}
		if assetIDStr == string(data) {
			log.Println("Already have the latest release...")
			os.Exit(0)
		}
	} else if os.IsNotExist(err) {
		out, err := os.Create(versionPath)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
	} else {
		log.Fatal(err)
	}

	url = apiPrefix + assetPath + assetIDStr
	log.Printf("Download jar via %s\n", url)
	res, err = doGet(url, &Header{authHeaderName, authHeader}, &Header{acceptHaderName, acceptHader})
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	out, err := os.Create("app.jar")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(versionPath, []byte(assetIDStr), 0755)
}
