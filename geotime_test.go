package geotime

import (
	"testing"
	"time"
)

var tests = []struct {
	latitude  float64
	longitude float64
	date      string
	sunrise   string
	sunset    string
	partofday string
}{
	{49.8, 24.03, "27.10.22 05:00:00 +0300", "27.10.22 08:04:45 +0300", "27.10.22 18:10:45 +0300", "night"},
	{49.8, 24.03, "30.10.22 13:00:00 +0300", "30.10.22 08:10:30 +0300", "30.10.22 18:04:30 +0300", "noon"},
	{49.8, 24.03, "31.10.22 07:00:00 +0200", "31.10.22 07:11:27 +0200", "31.10.22 17:03:27 +0200", "sunrise"},
	{39.73, -105.0, "01.02.06 15:04:05 -0700", "01.02.06 07:07:41 -0700", "01.02.06 17:19:41 -0700", "day"},
}

func TestCalculate(t *testing.T) {
	for _, test := range tests {
		date, err := time.Parse("02.01.06 15:04:05 -0700", test.date)
		if err != nil {
			t.Errorf("Error parsing string %s", test.date)
			continue
		}
		gt := Calculate(test.latitude, test.longitude, date)
		if gt.Sunrise.Format("02.01.06 15:04:05 -0700") != test.sunrise {
			t.Errorf("Calculate sunrise (%v, %v, %s) is %s, but must be %s", test.latitude, test.longitude, test.date,
				gt.Sunrise, test.sunrise)
		}
		if gt.Sunset.Format("02.01.06 15:04:05 -0700") != test.sunset {
			t.Errorf("Calculate sunset (%v, %v, %s) is %s, but must be %s", test.latitude, test.longitude, test.date,
				gt.Sunset, test.sunset)
		}
		if gt.PartOfDay != test.partofday {
			t.Errorf("Calculate parto of day (%v, %v, %s) is %s, but must be %s", test.latitude, test.longitude, test.date,
				gt.PartOfDay, test.partofday)
		}

	}
}

func TestSunrise(t *testing.T) {
	for _, test := range tests {
		date, err := time.Parse("02.01.06 15:04:05 -0700", test.date)
		if err != nil {
			t.Errorf("Error parsing string %s", test.date)
			continue
		}
		want, err := time.Parse("02.01.06 15:04:05 -0700", test.sunrise)
		if err != nil {
			t.Errorf("Error parsing string %s", test.date)
			continue
		}
		got := Sunrise(test.latitude, test.longitude, date)
		if !got.Equal(want) {
			t.Errorf("Sunrise(%v, %v, %s) in %v but must be in %s", test.latitude, test.longitude, test.date,
				got, want)
		}

	}
}

func TestSunset(t *testing.T) {
	for _, test := range tests {
		date, err := time.Parse("02.01.06 15:04:05 -0700", test.date)
		if err != nil {
			t.Errorf("Error parsing string %s", test.date)
			continue
		}
		want, err := time.Parse("02.01.06 15:04:05 -0700", test.sunset)
		if err != nil {
			t.Errorf("Error parsing string %s", test.date)
			continue
		}
		got := Sunset(test.latitude, test.longitude, date)
		if !got.Equal(want) {
			t.Errorf("Sunrise(%v, %v, %s) in %v but must be in %s", test.latitude, test.longitude, test.date,
				got, want)
		}

	}
}

func TestPartOfDay(t *testing.T) {
	for _, test := range tests {
		date, err := time.Parse("02.01.06 15:04:05 -0700", test.date)
		if err != nil {
			t.Errorf("Error parsing string %s", test.date)
			continue
		}
		got := PartOfDay(test.latitude, test.longitude, date, nil)
		if got != "" && got != test.partofday {
			t.Errorf("PartOfDay(%v, %v, %s, nil) is %v but must %s", test.latitude, test.longitude, test.date,
				got, test.partofday)
		}

	}
}
