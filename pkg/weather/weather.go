package weather

import (
	"encoding/json"
	"fmt"
	"github.com/rivo/tview"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
		Module: "weather",
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

func PopulateWeatherDisplay(weatherTextView *tview.TextView) {
	// Weather info
	weatherInfo, getWeatherError := GetWeather()
	if getWeatherError != nil {
		file, err := os.OpenFile("error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		log.SetOutput(file)
		log.Print(getWeatherError)
		fmt.Fprintf(weatherTextView, "Error, check error.log")
		return
	}

	go fmt.Fprintf(
		weatherTextView,
		"Weather in: %s\nCurrent temp: [red]%d °F[white]\nFeels like: [red]%d °F[white]\n",
		weatherInfo.Name,
		int(weatherInfo.Main.Temp),
		int(weatherInfo.Main.FeelsLike))
}

