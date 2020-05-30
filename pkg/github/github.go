package github

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Issues struct {
	Items []Issue `json:"items"`
}

type Issue struct {
	PullRequest PullRequest `json:"pull_request"`
}

type PullRequest struct {
	Url 				string 	`json:"url"`
	Number 				int 	`json:"number"`
	Id 					int 	`json:"id"`
	Title 				string 	`json:"title"`
	State 				string 	`json:"state"`
	NumberOfComments 	int 	`json:"comments"`
}

func GetPullRequests() (*[]PullRequest, error) {
	issues, err := getIssues()
	if err != nil {
		log.Fatal(err)
	}

	var pullRequests []PullRequest
	var response *http.Response

	for _, issue := range issues.Items {
		var pullRequest PullRequest
		client := http.Client{
			Timeout: time.Second * 3,
		}

		accessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
		request, err := http.NewRequest("GET", issue.PullRequest.Url, nil)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		request.SetBasicAuth("jtowe", accessToken)

		response, err = client.Do(request)
		if err != nil {
			return nil, err
		}

		data, _ := ioutil.ReadAll(response.Body)

		err = json.Unmarshal(data, &pullRequest)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		pullRequests = append(pullRequests, pullRequest)
	}

	defer response.Body.Close()

	return &pullRequests, nil
}

func getIssues() (*Issues, error){
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

	var issues Issues

	err = json.Unmarshal(data, &issues)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &issues, nil
}
