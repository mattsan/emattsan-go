package owm

import (
    "fmt"
    "encoding/json"
    "net/http"
)

type Weather struct {
    Id          int    `json:"id"`
    Main        string `json:main"`
    Description string `json:"description"`
}

type CurrentWeatherData struct {
    Weather []Weather `json:"weather"`
}

type Params struct {
    AppId string
    BaseUrl string
    City string
    Units string
}

func (params *Params) getCurrentWeatherUrl() string {
  return fmt.Sprintf("%s?q=%s&units=%s&appid=%s", params.BaseUrl, params.City, params.Units, params.AppId)
}

func (data *CurrentWeatherData) FetchCurrentWeather(params *Params) error {
    resp, err := http.Get(params.getCurrentWeatherUrl())

    if err != nil { return err }

    defer resp.Body.Close()

    decoder := json.NewDecoder(resp.Body)
    return decoder.Decode(data)
}
