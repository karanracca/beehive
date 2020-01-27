package openweathermapbee

import (
	"errors"
	"fmt"
	"strings"
	"time"

	owm "github.com/briandowns/openweathermap"
)

type Condition struct {
	isPresent bool
	time      []time.Time
}

type Weather struct {
	low  float64
	high float64
}

type FormattedWeather5Data struct {
	city    string
	snow    *Condition
	rain    *Condition
	weather *Weather
}

func (wd *FormattedWeather5Data) GetForecaseWeatherData(data owm.ForecastWeatherJson) error {

	forecast5WeatherData, ok := data.(*owm.Forecast5WeatherData)
	if !ok {
		return errors.New("Invalid data type in Forecast service in")
	}

	wd.city = forecast5WeatherData.City.Name

	for _, d := range forecast5WeatherData.List {
		t := time.Unix(int64(d.Dt), 0)

		if (t.Day() == time.Now().Day()+1) && (t.Hour() > 6) {

			if wd.weather.low == 0 || wd.weather.low > d.Main.Temp {
				wd.weather.low = d.Main.Temp
			}

			if wd.weather.high == 0 || wd.weather.high < d.Main.Temp {
				wd.weather.high = d.Main.Temp
			}

			if d.Weather[0].Main == "Snow" {
				wd.snow.isPresent = true
				wd.snow.time = append(wd.snow.time, t)
			}

			if d.Weather[0].Main == "Rain" {
				wd.rain.isPresent = true
				wd.rain.time = append(wd.rain.time, t)
			}

		}
	}

	return nil
}

func (wd *FormattedWeather5Data) String() string {
	var s strings.Builder

	fmt.Fprintf(&s, "The max/min temprature for today is %f/%f.", wd.weather.high, wd.weather.low)

	if wd.snow.isPresent {
		s.WriteString("\nIts going to snow today around ")
		for _, v := range wd.snow.time {
			fmt.Fprintf(&s, "%v, ", v)
		}
		s.WriteString(".\n")
	}

	if wd.snow.isPresent {
		s.WriteString("\nIts going to snow today around ")
		for _, v := range wd.snow.time {
			fmt.Fprintf(&s, "%v, ", v)
		}
		s.WriteString(".\n")
	}

	if wd.rain.isPresent {
		s.WriteString("\nIts going to rain today around ")
		for _, v := range wd.rain.time {
			fmt.Fprintf(&s, "%v, ", v)
		}
		s.WriteString(".\n")
	}

	return s.String()
}
