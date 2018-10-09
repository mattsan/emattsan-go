package owm

import (
    "encoding/json"
    "net/http"
    "net/url"
)

const (
    UrlCurrentWeather = "https://api.openweathermap.org/data/2.5/weather"
    UrlForecast5Days = "https://api.openweathermap.org/data/2.5/forecast"
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
    Coord   Coordinate `json:"coord"`
    Weather []Weather  `json:"weather"`
    Base    string     `json:"base"`
    Main    Main       `json:"main"`
    Wind    Wind       `json:"wind"`
    Clouds  Clouds     `json:"clouds"`
    Rain    Value3h    `json:"rain"`
    Snow    Value3h    `json:"snow"`
    Dt      int        `json:"dt"`
    Sys     Sys        `json:"sys"`
    Id      int        `json:"id"`
    Name    string     `json:"name"`
}

type City struct {
    Id      int        `json:"id"`
    Name    string     `json:"name"`
    Coord   Coordinate `json:"coord"`
    Country string     `json:"country"`
}

type ListItem struct {
    Dt      int       `json:"dt"`
    Main    Main      `json:"main"`
    Weather []Weather `json:"weather"`
    Clouds  Clouds    `json:"clouds"`
    Wind    Wind      `json:"wind"`
    Rain    Value3h   `json:"rain"`
    Snow    Value3h   `json:"snow"`
    DtTxt   string    `json:"dt_txt"`
}

type Forecast5DaysData struct {
    City City `json:"city"`
    Cnt  int  `json:"cnt"`
    List []ListItem `json:"list"`
}

type Query struct {
    AppId string
    City string
    Units string
    Language string
}

func (query *Query) Encode() string {
    values := url.Values{}
    values.Set("q", query.City)
    values.Set("units", query.Units)
    values.Set("lang", query.Language)
    values.Set("appid", query.AppId)
    return values.Encode()
}

func Tokyo(appid string) Query {
    return Query{
        AppId: appid,
        City: "Tokyo,jp",
        Units: "metric",
        Language: "ja",
    }
}

func (query *Query) GetchCurrentWeather() (*CurrentWeatherData, error) {
    url, _ := url.Parse(UrlCurrentWeather)
    url.RawQuery = query.Encode()

    resp, err := http.Get(url.String())
    if err != nil { return nil, err }
    defer resp.Body.Close()

    data := new(CurrentWeatherData)
    decoder := json.NewDecoder(resp.Body)
    err = decoder.Decode(data)
    return data, err
}

func (query *Query) GetchForecast5Days() (*Forecast5DaysData, error) {
    url, _ := url.Parse(UrlForecast5Days)
    url.RawQuery = query.Encode()

    println(url.String())

    resp, err := http.Get(url.String())
    if err != nil { return nil, err }
    defer resp.Body.Close()

    data := new(Forecast5DaysData)
    decoder := json.NewDecoder(resp.Body)
    err = decoder.Decode(data)
    return data, err
}
