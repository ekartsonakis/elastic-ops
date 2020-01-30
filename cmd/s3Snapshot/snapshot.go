package s3snapshot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	// "io/ioutil"
	// "os/exec"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

// SnapData are Input data needed for snapshot
type SnapData struct {
	Indices            string `json:"indices"`
	IgnoreUnavailable  bool   `json:"ignore_unavailable"`
	IncludeGlobalState bool   `json:"include_global_state"`
}

// Snapshot is a typical Snapshot structure
type Snapshot struct {
	Snapshot           string        `json:"snapshot"`
	UUID               string        `json:"uuid"`
	Indices            []string      `json:"indices"`
	IncludeGlobalState bool          `json:"include_global_state"`
	State              string        `json:"state"`
	StartTime          time.Time     `json:"start_time"`
	EndTime            time.Time     `json:"end_time,omitempty"`
	DurationInMillis   int           `json:"duration_in_millis"`
	Failures           []interface{} `json:"failures"`
	Shards             struct {
		Total      int `json:"total"`
		Failed     int `json:"failed"`
		Successful int `json:"successful"`
	} `json:"shards"`
}

// GetSnapResp is Snapshots response
type GetSnapResp struct {
	Snapshots []Snapshot `json:"snapshots"`
}

// TakeSnapResp is taken Snapshots response
type TakeSnapResp struct {
	Accepted bool `json:"accepted"`
	Error interface{} `json:"error"`
}

// DeleteSnapResp is taken Snapshots response
type DeleteSnapResp struct {
	Acknowledged bool `json:"acknowledged"`
	Error interface{} `json:"error"`
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

func deleteReqWithAuth(username, password, server, endpoint string) (*http.Response, error) {

	url := fmt.Sprintf("https://%s:9200/%s", server, endpoint)
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func getSnapList(username, password, server, s3Repo string) (bool, int, error){
	inprogress := false
	log.Info("Getting snapshots available")
	getSnapEndpoint := fmt.Sprintf("_snapshot/%s/_all", s3Repo)
	resp, err := getReqWithAuth(username, password, server, getSnapEndpoint)
	if err != nil {
		return inprogress, 0, err
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	// log.Info(resp.StatusCode)

	var snapshots GetSnapResp
	if err := json.Unmarshal(bodyText, &snapshots); err != nil {
		return inprogress, 0, err
	}
	// log.Info("Unmarshaled snap")
	for _, s := range snapshots.Snapshots {
		if s.State == "SUCCESS" {
			log.Infof("Snapshot %s for indices %s is in %s state and took %d ms to complete.", s.Snapshot, s.Indices, s.State, s.DurationInMillis)
		} else {
			log.Infof("Snapshot %s for indices %s is in %s state.", s.Snapshot, s.Indices, s.State,)
			inprogress = true
		}
	}
	return inprogress, resp.StatusCode, nil
}

func takeNewSnap(username, password, server, s3Repo, indices, snapName string, globalState bool) error{
	log.Info("Taking a new snapshot")
	takeSnapEndpoint := fmt.Sprintf("_snapshot/%s/%s", s3Repo, snapName)
	resp, err := putReqWithAuth(username, password, server, takeSnapEndpoint, indices, globalState)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	log.Info(string(bodyText))
	// log.Info(resp.StatusCode)
	var t TakeSnapResp
	if err := json.Unmarshal(bodyText, &t); err != nil {
		return err
	}
	if !t.Accepted {
		return fmt.Errorf("%v", t.Error)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Response code %s is not accepted", resp.StatusCode)
	}
	return nil
}

func deleteSnap(username, password, server, s3Repo, snapName string) error{
	log.Infof("Deleting snapshot %s", snapName)
	deleteSnapEndpoint := fmt.Sprintf("_snapshot/%s/%s", s3Repo, snapName)
	resp, err := deleteReqWithAuth(username, password, server, deleteSnapEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// log.Info(string(bodyText))
	// log.Info(resp.StatusCode)
	var d DeleteSnapResp
	if err := json.Unmarshal(bodyText, &d); err != nil {
		return err
	}
	if !d.Acknowledged {
		return fmt.Errorf("%v", d.Error)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Response code %s is not accepted", resp.StatusCode)
	}
	return nil
}

// Snap represents the snapshot command
func Snap(username, password, server, s3Repo, indices, snapName string, onlyList, globalState, delete bool) {
	
	password = strings.Replace(password, "$", `\$`, -1)
	if onlyList {
		log.Infof("Only listing of snapshots under %s ordered", s3Repo)
		_, _, err := getSnapList(username, password, server, s3Repo)
		if err != nil {
			log.Fatalf("List Snapshots failed: %s", err)
		}
	} else if delete {
		log.Infof("Deleting of %s snapshot.", snapName)
		_, statusCode, err := getSnapList(username, password, server, s3Repo)
		if err != nil {
			log.Fatalf("List Snapshots failed: %s", err)
		}
		if statusCode == 200 {
			if err = deleteSnap(username, password, server, s3Repo, snapName); err != nil {
				log.Fatalf("Deleting of Snapshot %s was not accepted: %s", snapName, err)
			}

			// Checking if snapshot deleted
			log.Info("Checking deletion")
			_, _, err = getSnapList(username, password, server, s3Repo)
			if err != nil {
				log.Fatalf("List Snapshots failed: %s", err)
			}
		}

	} else {
		log.Infof("Taking new snapshot with %s name for indices %s for repository %s ordered", snapName, indices, s3Repo)
		inprogress, statusCode, err := getSnapList(username, password, server, s3Repo)
		if err != nil {
			log.Fatalf("List Snapshots failed: %s", err)
		}
		log.Info("New snapshot was accepted")
		if statusCode == 200 {
			if inprogress {
				log.Warn("Keep in mind another snapshot is in progress.")
			}
			//Taking the new snapshot
			if err = takeNewSnap(username, password, server, s3Repo, indices, snapName, globalState); err != nil {
				log.Fatalf("Taking of new Snapshot was not accepted: %s", err)
			}

			// Checking if snapshot started
			log.Info("Checking progress")
			_, _, err = getSnapList(username, password, server, s3Repo)
			if err != nil {
				log.Fatalf("List Snapshots failed: %s", err)
			}	
		}
	}
	
	log.Info("finished")
}
