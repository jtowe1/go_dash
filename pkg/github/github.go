package github

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	issuesUrl string = "https://github.dev.shootproof.com/api/v3/search/issues"
	author string = "jtowe"
)

type Issues struct {
	Items []Issue `json:"items"`
}

type Issue struct {
	PullRequest PullRequest `json:"pull_request"`
}

type PullRequest struct {
	Url	string `json:"url"`
	Number int `json:"number"`
	Id int `json:"id"`
	Title string `json:"title"`
	State string `json:"state"`
	NumberOfComments int `json:"comments"`
	Labels []Label `json:"labels"`
	Additions int `json:"additions"`
	Deletions int `json:"deletions"`
}

type Label struct {
	Name string `json:"name"`
	Color string `json:"color"`
}

type client struct {
	client *http.Client
}

func newGithubClient() *client {
	return &client {
		client: &http.Client{
			Timeout: time.Second * 3,
		},
	}
}

func GetPullRequests() (*[]PullRequest, error) {
	issues, err := getIssues()
	if err != nil {
		return nil, err
	}

	var pullRequests []PullRequest
	var response *http.Response

	for _, issue := range issues.Items {
		var pullRequest PullRequest
		githubClient := newGithubClient()

		accessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
		request, err := http.NewRequest("GET", issue.PullRequest.Url, nil)
		if err != nil {
			return nil, err
		}
		request.SetBasicAuth("jtowe", accessToken)

		response, err = githubClient.client.Do(request)
		if err != nil {
			return nil, err
		}

		data, _ := ioutil.ReadAll(response.Body)

		err = json.Unmarshal(data, &pullRequest)
		if err != nil {
			return nil, err
		}

		pullRequests = append(pullRequests, pullRequest)
	}

	defer response.Body.Close()

	return &pullRequests, nil
}

func getIssues() (*Issues, error){
	githubClient := newGithubClient()

	accessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	request, err := http.NewRequest("GET", issuesUrl + "?q=state:open+type:pr+author:" + author, nil)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth("jtowe", accessToken)

	response, err := githubClient.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	data, _ := ioutil.ReadAll(response.Body)

	var issues Issues

	err = json.Unmarshal(data, &issues)
	if err != nil {
		return nil, err
	}

	return &issues, nil
}
