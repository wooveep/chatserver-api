/*
 * @Author: cloudyi.li
 * @Date: 2023-06-15 11:29:29
 * @LastEditTime: 2023-06-21 16:20:06
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/search/weather.go
 */
package search

import (
	"chatserver-api/utils/tools"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Forecast10Response struct {
	CalendarDayTemperatureMax []int     `json:"calendarDayTemperatureMax"`
	CalendarDayTemperatureMin []int     `json:"calendarDayTemperatureMin"`
	DayOfWeek                 []string  `json:"dayOfWeek"`
	ExpirationTimeUtc         []int     `json:"expirationTimeUtc"`
	MoonPhase                 []string  `json:"moonPhase"`
	MoonPhaseCode             []string  `json:"moonPhaseCode"`
	MoonPhaseDay              []int     `json:"moonPhaseDay"`
	MoonriseTimeLocal         []string  `json:"moonriseTimeLocal"`
	MoonriseTimeUtc           []int     `json:"moonriseTimeUtc"`
	MoonsetTimeLocal          []string  `json:"moonsetTimeLocal"`
	MoonsetTimeUtc            []int     `json:"moonsetTimeUtc"`
	Narrative                 []string  `json:"narrative"`
	Qpf                       []float64 `json:"qpf"`
	QpfSnow                   []float64 `json:"qpfSnow"`
	SunriseTimeLocal          []string  `json:"sunriseTimeLocal"`
	SunriseTimeUtc            []int     `json:"sunriseTimeUtc"`
	SunsetTimeLocal           []string  `json:"sunsetTimeLocal"`
	SunsetTimeUtc             []int     `json:"sunsetTimeUtc"`
	TemperatureMax            []int     `json:"temperatureMax"`
	TemperatureMin            []int     `json:"temperatureMin"`
	ValidTimeLocal            []string  `json:"validTimeLocal"`
	ValidTimeUtc              []int     `json:"validTimeUtc"`
	Daypart                   []struct {
		CloudCover            []int     `json:"cloudCover"`
		DayOrNight            []string  `json:"dayOrNight"`
		DaypartName           []string  `json:"daypartName"`
		IconCode              []int     `json:"iconCode"`
		IconCodeExtend        []int     `json:"iconCodeExtend"`
		Narrative             []string  `json:"narrative"`
		PrecipChance          []int     `json:"precipChance"`
		PrecipType            []string  `json:"precipType"`
		Qpf                   []float64 `json:"qpf"`
		QpfSnow               []float64 `json:"qpfSnow"`
		QualifierCode         []string  `json:"qualifierCode"`
		QualifierPhrase       []string  `json:"qualifierPhrase"`
		RelativeHumidity      []int     `json:"relativeHumidity"`
		SnowRange             []string  `json:"snowRange"`
		Temperature           []int     `json:"temperature"`
		TemperatureHeatIndex  []int     `json:"temperatureHeatIndex"`
		TemperatureWindChill  []int     `json:"temperatureWindChill"`
		ThunderCategory       []string  `json:"thunderCategory"`
		ThunderIndex          []int     `json:"thunderIndex"`
		UvDescription         []string  `json:"uvDescription"`
		UvIndex               []int     `json:"uvIndex"`
		WindDirection         []int     `json:"windDirection"`
		WindDirectionCardinal []string  `json:"windDirectionCardinal"`
		WindPhrase            []string  `json:"windPhrase"`
		WindSpeed             []int     `json:"windSpeed"`
		WxPhraseLong          []string  `json:"wxPhraseLong"`
		WxPhraseShort         []string  `json:"wxPhraseShort"`
	} `json:"daypart"`
}

type WeatherClient struct {
	api_key     string
	http_client http.Client
}

func NewWeatherClient(api_key string) WeatherClient {
	return WeatherClient{
		api_key,
		http.Client{},
	}
}

func (c *WeatherClient) Weathermakeapiurl(lat float64, lng float64, units string) string {
	if units == "" {
		units = "m"
	}
	url := fmt.Sprintf("https://api.weather.com/v3/wx/forecast/daily/5day?geocode=%f,%f&format=json&units=%s&language=zh-CN&apiKey=%s",
		lat, lng,
		units, c.api_key)
	//log.Debug(url)
	return url
}

func (c *WeatherClient) Weathermakeapirequest(url string, payload interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.New("Could not send request: " + err.Error())
	}

	res, err := c.http_client.Do(req)
	if err != nil {
		return errors.New("Could not read response: " + err.Error())
	}

	defer res.Body.Close()

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(payload)
	if err != nil {
		return errors.New("Could not decode: " + err.Error())
	}

	return nil
}

func (c *WeatherClient) WeatherdoGetForecast5(url string) (string, error) {
	var payload Forecast10Response
	var content []string
	content = append(content, "Day|Week|TemperatureMax(Celsius)|TemperatureMin(Celsius)|narrative|qpf(percent)|sunriseTime|sunsetTime")
	err := c.Weathermakeapirequest(url, &payload)
	if err != nil {
		return "", err
	}
	for i := 0; i < 6; i++ {
		day, SunriseTime, week, err := tools.TimeConvert(payload.SunriseTimeLocal[i])
		_, SunsetTime, _, err := tools.TimeConvert(payload.SunsetTimeLocal[i])
		if err != nil {
			return "", err
		}
		content = append(content, fmt.Sprintf("%s|%s|%d|%d|%s|%.2f%%|%s|%s", day, week, payload.CalendarDayTemperatureMax[i], payload.CalendarDayTemperatureMin[i], payload.Narrative[i], payload.Qpf[i], SunriseTime, SunsetTime))
	}
	return strings.Join(content, "\n"), nil
}
