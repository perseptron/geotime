// Mandrakesoft 2022

// Package geotime calculate astronomical events as sunrise, solar noon,
// sunset and some other geographical date-time important variables
package geotime

import (
	"math"
	"time"
)

type Geotime struct {
	Jd        float64
	Jc        float64
	Date      time.Time
	Lat       float64
	Long      float64
	SolarNoon time.Time
	HourAngle time.Duration
	Sunrise   time.Time
	Sunset    time.Time
	PartOfDay string
	obliq     float64
	ganom     float64
	glong     float64
}

const (
	night   = "night"
	sunrise = "sunrise"
	morning = "morning"
	day     = "day"
	noon    = "noon"
	sunset  = "sunset"
	evening = "evening"
)

func Calculate(lat, long float64, date time.Time) Geotime {
	gt := new(Geotime)
	preCalc(date, gt)
	gt.Date = date
	gt.Lat = lat
	gt.Long = long
	gt.Jd = JD(date)
	gt.Jc = JC(date)
	gt.SolarNoon = SolarNoon(long, date, gt)
	gt.HourAngle = HourAngle(lat, date, gt)
	gt.Sunrise = gt.SolarNoon.Add(-gt.HourAngle)
	gt.Sunset = gt.SolarNoon.Add(gt.HourAngle)
	gt.PartOfDay = PartOfDay(lat, long, date, gt)
	return *gt
}

func JD(date time.Time) float64 {
	epoch := time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC)
	days := date.Sub(epoch).Hours() / 24
	_, tz := date.Zone()
	return days + 2 + 2415018.5 + float64(date.Hour()/24) - float64(tz/60/60)/24
}

func JC(date time.Time) float64 {
	return (JD(date) - 2451545) / 36525
}

func obliquity(jc float64) float64 {
	return 23 + (26+(21.448-jc*(46.815+jc*(0.00059-jc*0.001813)))/60)/60 +
		0.00256*math.Cos(toRad(125.04-1934.136*jc))
}

func geomAnomaly(jc float64) float64 {
	return 357.52911 + jc*(35999.05029-0.0001537*jc)
}

func geomLong(jc float64) float64 {
	return math.Mod(280.46646+jc*(36000.76983+jc*0.0003032), 360)
}

func preCalc(date time.Time, gt *Geotime) {
	gt.Jd = JD(date)
	gt.Jc = JC(date)
	gt.obliq = obliquity(gt.Jc)
	gt.ganom = geomAnomaly(gt.Jc)
	gt.glong = geomLong(gt.Jc)
}

func SolarNoon(long float64, date time.Time, gt *Geotime) time.Time {
	if gt == nil {
		gt = new(Geotime)
		preCalc(date, gt)
	}

	_, tz := date.Zone()
	y := math.Tan(toRad(gt.obliq/2)) * math.Tan(toRad(gt.obliq/2))
	eo := 0.016708634 - gt.Jc*(0.000042037+0.0000001267*gt.Jc)
	eqt := 4 * toDeg(y*math.Sin(2*toRad(gt.glong))-2*eo*math.Sin(toRad(gt.ganom))+
		4*eo*y*math.Sin(toRad(gt.ganom))*math.Cos(2*toRad(gt.glong))-
		0.5*y*y*math.Sin(4*toRad(gt.glong))-1.25*eo*eo*math.Sin(2*toRad(gt.ganom)))
	sn := 720 - 4*long - eqt + float64(tz/60)
	gt.SolarNoon = truncTime(date).Add(time.Duration(sn*60) * time.Second)
	return gt.SolarNoon
}

func HourAngle(lat float64, date time.Time, gt *Geotime) time.Duration {
	if gt == nil {
		gt = new(Geotime)
		preCalc(date, gt)
	}
	eqCtr := math.Sin(toRad(gt.ganom))*(1.914602-gt.Jc*(0.004817+0.000014*gt.Jc)) +
		math.Sin(toRad(2*gt.ganom))*(0.019993-0.000101*gt.Jc) + math.Sin(toRad(3*gt.ganom))*0.000289
	truelong := gt.glong + eqCtr
	app := truelong - 0.00569 - 0.00478*math.Sin(toRad(125.04-1934.136*gt.Jc))
	decl := toDeg(math.Asin(math.Sin(toRad(gt.obliq)) * math.Sin(toRad(app))))
	HA := toDeg(math.Acos(math.Cos(toRad(90.833))/(math.Cos(toRad(lat))*math.Cos(toRad(decl))) -
		math.Tan(toRad(lat))*math.Tan(toRad(decl))))
	return time.Duration(HA*4) * time.Minute
}

func Sunrise(lat, long float64, date time.Time) time.Time {
	var gt Geotime
	preCalc(date, &gt)
	return SolarNoon(long, date, &gt).Add(-HourAngle(lat, date, &gt))

}

func Sunset(lat, long float64, date time.Time) time.Time {
	var gt Geotime
	preCalc(date, &gt)
	return SolarNoon(long, date, &gt).Add(HourAngle(lat, date, &gt))
}

func PartOfDay(lat, long float64, date time.Time, gt *Geotime) string {
	if gt == nil {
		gt = new(Geotime)
		preCalc(date, gt)
		SolarNoon(long, date, gt)
		HourAngle(lat, date, gt)
		gt.Sunrise = gt.SolarNoon.Add(-gt.HourAngle)
		gt.Sunset = gt.SolarNoon.Add(gt.HourAngle)
	}

	if date.After(gt.Sunset.Add(time.Minute * 60)) {
		return night
	}
	if date.After(gt.Sunset) {
		return evening
	}
	if date.After(gt.Sunset.Add(-time.Minute * 20)) {
		return sunset
	}
	if date.After(gt.SolarNoon.Add(-time.Minute*10)) && date.Before(gt.SolarNoon.Add(time.Minute*10)) {
		return noon
	}
	if date.After(gt.Sunrise.Add(time.Minute * 60)) {
		return day
	}
	if date.After(gt.Sunrise) {
		return morning
	}
	if date.After(gt.Sunrise.Add(-time.Minute * 20)) {
		return sunrise
	}
	return night
}

func truncTime(date time.Time) time.Time {
	//TODO: correct error at the day of time shift
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}

func toDeg(rad float64) float64 {
	return rad * 180 / math.Pi
}

func toRad(deg float64) float64 {
	return deg * math.Pi / 180
}
