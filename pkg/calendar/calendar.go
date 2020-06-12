package calendar

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rivo/tview"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Widget struct {
	Row int
	Col int
	RowSpan int
	ColSpan int
	MinGridHeight int
	MinGridWidth int
	View *tview.TextView
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

type client struct {
	client *http.Client
}

func GetWidget(app *tview.Application) *Widget {
	calendarTextView := tview.NewTextView()
	calendarTextView.SetBorder(true).SetTitle("📅  Calendar")
	calendarTextView.SetChangedFunc(func() {
		app.Draw()
	})

	widget := Widget{
		View: calendarTextView,
		Row: 0,
		Col: 0,
		RowSpan: 2,
		ColSpan: 1,
		MinGridHeight: 0,
		MinGridWidth: 100,
		Module: "calendar",
	}

	return &widget
}

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func GetCalendar() (*calendar.Events, error) {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)
	client.Timeout = time.Second * 3

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}

	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		return events, nil
	}
	return nil, err
}