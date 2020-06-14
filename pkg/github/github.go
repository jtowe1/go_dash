package github

import (
	"encoding/json"
	"github.com/rivo/tview"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	issuesUrl string = "https://github.dev.shootproof.com/api/v3/search/issues"
	author string = "jtowe"
)

type Widget struct {
	//goDash.TableWidget
	Row int
	Col int
	RowSpan int
	ColSpan int
	MinGridHeight int
	MinGridWidth int
	View *tview.Table
	Module string
}

func (w *Widget) GetView() interface{} {
	return w.View
}

func (w *Widget) GetRow() int {
	return w.Row
}

func (w *Widget) GetCol() int {
	return w.Col
}

func (w *Widget) GetRowSpan() int {
	return w.RowSpan
}

func (w *Widget) GetColSpan() int {
	return w.ColSpan
}

func (w *Widget) GetMinGridHeight() int {
	return w.MinGridHeight
}

func (w *Widget) GetMinGridWidth() int {
	return w.MinGridWidth
}

func (w *Widget) GetModule() string {
	return w.Module
}

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

func GetWidget() *Widget {
	githubTable := tview.NewTable()
	githubTable.SetBorders(true)

	widget := Widget{
		View: githubTable,
		Row: 1,
		Col: 1,
		RowSpan: 1,
		ColSpan: 2,
		MinGridHeight: 0,
		MinGridWidth: 100,
		Module: "github",
	}

	return &widget
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

func PopulateGithubDisplay(githubTable *tview.Table, app *tview.Application) {
	pullRequests, gitHubError := GetPullRequests()
	if gitHubError != nil {
		file, err := os.OpenFile("error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		log.SetOutput(file)
		log.Print(gitHubError)
		githubTable.SetCell(0, 0, tview.NewTableCell("Error, check error.log"))
		app.Draw()
		return
	}

	githubTable.SetCell(0, 0, tview.NewTableCell("️[aquamarine]Open Pull Requests authored by Jeremiah[white]"))
	githubTable.SetCell(0, 1, tview.NewTableCell("[aquamarine]Comments[white]"))
	githubTable.SetCell(0, 2, tview.NewTableCell("[aquamarine]Labels[white]"))
	githubTable.SetCell(0, 3, tview.NewTableCell("[aquamarine]Additions[white]"))
	githubTable.SetCell(0, 4, tview.NewTableCell("[aquamarine]Deletions[white]"))

	rowCounter := 1
	for _, pullRequest := range *pullRequests {
		githubTable.SetCell(rowCounter, 0, tview.NewTableCell(pullRequest.Title))
		githubTable.SetCell(rowCounter, 1, tview.NewTableCell(strconv.Itoa(pullRequest.NumberOfComments)).SetAlign(tview.AlignCenter))

		labels := ""
		for _, label := range pullRequest.Labels {
			labels += "[#" + label.Color +"]" + label.Name + " "
		}
		githubTable.SetCell(rowCounter, 2, tview.NewTableCell(labels))
		githubTable.SetCell(rowCounter, 3, tview.NewTableCell("[green]" + strconv.Itoa(pullRequest.Additions) + "[white]").SetAlign(tview.AlignCenter))
		githubTable.SetCell(rowCounter, 4, tview.NewTableCell("[red]" + strconv.Itoa(pullRequest.Deletions) + "[white]").SetAlign(tview.AlignCenter))
		rowCounter++
	}

	app.Draw()
}
