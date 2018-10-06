package owm

import (
    "fmt"
    "encoding/json"
    "net/http"
)

type Coordinate struct {
    longitude float64 `json:"lon"`
    latitude  float64 `json:"lat"`
}

type Weather struct {
    Id          int    `json:"id"`
    Main        string `json:main"`
    Description string `json:"description"`
    Icon        string `json:"icon"`
}

type Main struct {
    Temp      float64 `json:"temp"`
    Pressure  int     `json:"pressure"`
    Humidity  int     `json:"humidity"`
    TempMin   float64 `json:"temp_min"`
    TempMax   float64 `json:"temp_max"`
    SeaLevel  int     `json:"sea_level"`
    GrndLevel int     `json:"grnd_level"`
}

type Wind struct {
    Speed float64 `json:"speed"`
    Deg   int     `json:"deg"`
}

type Clouds struct {
    All int `json:"all"`
}

type Value3h struct {
    Value int `json:"3h"`
}

type Sys struct {
    Country string `json:"country"`
    Sunrise int    `json:"sunrise"`
    Sunset  int    `json:"sunset"`
}

type CurrentWeatherData struct {
    Coordinate Coordinate `json:"coord"`
    Weather    []Weather  `json:"weather"`
    Base       string     `json:"base"`
    Main       Main       `json:"main"`
    Wind       Wind       `json:"wind"`
    Clouds     Clouds     `json:"clouds"`
    Rain       Value3h    `json:"rain"`
    Snow       Value3h    `json:"snow"`
    Dt         int        `json:"dt"`
    Sys        Sys        `json:"sys"`
    Id         int        `json:"id"`
    Name       string     `json:"name"`
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
