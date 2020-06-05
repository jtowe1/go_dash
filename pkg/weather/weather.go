package weather

import (
	"encoding/json"
	"fmt"
	"github.com/rivo/tview"
	"io/ioutil"
	"jeremiahtowe.com/go_dash/goDash"
	"log"
	"net/http"
	"os"
)

type Widget struct {
	goDash.TextViewWidget
	Row int
	Col int
	RowSpan int
	ColSpan int
	MinGridHeight int
	MinGridWidth int
	View *tview.TextView
}

type Data struct {
	Weather []weatherStationData `json:"weather"`
	Name string `json:"name"`
	Main main `json:"main"`

}

type weatherStationData struct {
	Id int `json:"id"`
	Main string `json:"main"`
	Description string `json:"description"`
	Icon string `json:"icon"`
}

type main struct {
	Temp float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
}

func GetWidget(app *tview.Application) *Widget {
	weatherTextView := tview.NewTextView().SetDynamicColors(true)
	weatherTextView.SetBorder(true).SetTitle("☁️  Weather")
	weatherTextView.SetChangedFunc(func() {
		app.Draw()
	})

	widget := Widget{
		View: weatherTextView,
		Row: 0,
		Col: 2,
		RowSpan: 1,
		ColSpan: 1,
		MinGridHeight: 0,
		MinGridWidth: 100,
	}

	return &widget
}

func GetWeather() (*Data, error){
	apiKey := os.Getenv("OPEN_WEATHER_MAP_API_KEY")
	response, err := http.Get("https://api.openweathermap.org/data/2.5/weather?q=kennesaw&appid=" + apiKey)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("status code %d returned", response.StatusCode)
	}

	data, _ := ioutil.ReadAll(response.Body)

	var weatherData Data

	err = json.Unmarshal(data, &weatherData)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	kelvinToFahrenheit(&weatherData.Main.Temp)
	kelvinToFahrenheit(&weatherData.Main.FeelsLike)

	return &weatherData, nil
}

func kelvinToFahrenheit(tempInKelvin *float64) {
	*tempInKelvin = (((*tempInKelvin - 273.15) * 9) / 5) + 32
}

