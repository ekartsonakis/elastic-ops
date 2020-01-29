package s3snapshot

import (
	"bytes"
	"encoding/json"
	"fmt"

	// "io/ioutil"
	// "os/exec"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type SnapData struct {
	Indices            string `json:"indices"`
	IgnoreUnavailable  bool   `json:"ignore_unavailable"`
	IncludeGlobalState bool   `json:"include_global_state"`
}

func getReqWithAuth(username, password, server, endpoint string) (*http.Response, error) {

	url := fmt.Sprintf("https://%s:9200/%s", server, endpoint)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func putReqWithAuth(username, password, server, endpoint, indices string, globalState bool) (*http.Response, error) {

	url := fmt.Sprintf("https://%s:9200/%s", server, endpoint)
	client := &http.Client{}
	data := SnapData{
		Indices:            indices,
		IgnoreUnavailable:  true,
		IncludeGlobalState: globalState,
	}
	json, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(json))
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(json))
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// Snap represents the snapshot command
func Snap(username, password, server, indices, s3Repo string, globalState bool) {
	log.Info("starting snaphsoting...")
	password = strings.Replace(password, "$", `\$`, -1)
	log.Info(password)
	log.Info(username)
	log.Info(server)
	log.Info(indices)
	log.Info(s3Repo)
	log.Info(globalState)

	//Get all available snapshots first
	log.Info("Getting snapshots available")
	endpoint := fmt.Sprintf("_snapshot/%s/_all", s3Repo)
	resp, err := getReqWithAuth(username, password, server, endpoint)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	log.Info(string(bodyText))
	log.Info(resp.StatusCode)
	if resp.StatusCode == 200 {
		//Taking the new snapshot
		log.Info("Taking a new snapshot")
		snapName := "%3Csnapshot-%7Bnow%2Fd%7D%3E"
		endpoint = fmt.Sprintf("_snapshot/%s/%s", s3Repo, snapName)
		resp, err = putReqWithAuth(username, password, server, endpoint, indices, globalState)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		bodyText, err = ioutil.ReadAll(resp.Body)
		log.Info(string(bodyText))
		log.Info(resp.StatusCode)

		if resp.StatusCode == 200 {
			// Checking if snapshot started
			log.Info("Snapshot started")
			log.Info("Checking progress")
			endpoint = fmt.Sprintf("_snapshot/%s/_all", s3Repo)
			resp, err = getReqWithAuth(username, password, server, endpoint)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			bodyText, err = ioutil.ReadAll(resp.Body)
			log.Info(string(bodyText))
			log.Info(resp.StatusCode)
		}
	}

	log.Info("finished")
}
