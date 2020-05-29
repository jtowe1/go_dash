package github

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Data struct {
	Items []PullRequest `json:"items"`
}

type PullRequest struct {
	Number int `json:"number"`
	Id int `json:"id"`
	Title string `json:"title"`
	State string `json:"state"`
	NumberOfComments int `json:"comments"`
}

func GetPullRequests() (*Data, error){
	client := http.Client{
		Timeout: time.Second * 3,
	}

	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	request, err := http.NewRequest("GET", "https://github.dev.shootproof.com/api/v3/search/issues?q=state:open+type:pr+author:jtowe", nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	request.SetBasicAuth("jtowe", accessToken)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	data, _ := ioutil.ReadAll(response.Body)

	var pullRequestData Data

	err = json.Unmarshal(data, &pullRequestData)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &pullRequestData, nil
}
