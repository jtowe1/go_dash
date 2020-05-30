package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Data struct {
	Weather []WeatherStationData `json:"weather"`
	Name string `json:"name"`
	Main Main `json:"main"`

}

type WeatherStationData struct {
	Id int `json:"id"`
	Main string `json:"main"`
	Description string `json:"description"`
	Icon string `json:"icon"`
}

type Main struct {
	Temp float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
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

