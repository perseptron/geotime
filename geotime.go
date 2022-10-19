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
	HourAngle  time.Duration
	Sunrise   time.Time
	Sunset    time.Time
	PartOfDay string
}

const (
	night = "night"
	sunrise = "sunrise"
	morning = "morning"
	day = "day"
	noon = "noon"
	sunset = "sunset"
	evening = "evening"
)

var gt Geotime
var jd, jc, ganom, glong, obliq float64

func Calculate(lat, long float64, date time.Time) (gt Geotime) {
	preCalc(date)
	gt.Date = date
	gt.Lat = lat
	gt.Long = long
	gt.Jd = jd
	gt.Jc = jc
	gt.SolarNoon = SolarNoon(long, date)
	gt.HourAngle = HourAngle(lat,date)
	gt.Sunrise = gt.SolarNoon.Add( -gt.HourAngle)
	gt.Sunset = gt.SolarNoon.Add(gt.HourAngle)
	gt.PartOfDay = PartOfDay(lat,long,date)
	return gt
}

func JD(date time.Time) float64 {
	epoch := time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC)
	days := date.Sub(epoch).Hours() / 24
	_, tz := date.Zone()
	return days+2 + 2415018.5 + float64(date.Hour()/24) - float64(tz/60/60)/24
}

func JC() float64 {
	return (jd - 2451545) / 36525
}

func obliquity() float64 {
	return 23 + (26+(21.448-jc*(46.815+jc*(0.00059-jc*0.001813)))/60)/60 +
		0.00256*math.Cos(toRad(125.04-1934.136*jc))
}

func geomAnomaly() float64 {
	return 357.52911 + jc*(35999.05029-0.0001537*jc)
}

func geomLong() float64 {
	return math.Mod(280.46646+jc*(36000.76983+jc*0.0003032), 360)
}

func preCalc(date time.Time) {
	if jd ==0 {
		jd = JD(date)
	}
	if jc == 0 {
		jc = JC()
	}
	if ganom == 0 {
		ganom = geomAnomaly()
	}
	if glong == 0 {
		glong = geomLong()
	}
	if obliq == 0 {
		obliq = obliquity()
	}
}

func SolarNoon(long float64, date time.Time) time.Time {
	preCalc(date)
	_, tz := date.Zone()
	y := math.Tan(toRad(obliq/2)) * math.Tan(toRad(obliq/2))
	eo := 0.016708634 - jc*(0.000042037+0.0000001267*jc)
	eqt := 4 * toDeg(y*math.Sin(2*toRad(glong))-2*eo*math.Sin(toRad(ganom))+
		4*eo*y*math.Sin(toRad(ganom))*math.Cos(2*toRad(glong))-
		0.5*y*y*math.Sin(4*toRad(glong))-1.25*eo*eo*math.Sin(2*toRad(ganom)))
	sn := 720 - 4*long - eqt + float64(tz/60)
	return truncTime(date).Add(time.Duration(sn*60) * time.Second)
}

func HourAngle(lat float64, date time.Time) time.Duration {
	preCalc(date)
	eqCtr := math.Sin(toRad(ganom))*(1.914602-jc*(0.004817+0.000014*jc)) +
		math.Sin(toRad(2*ganom))*(0.019993-0.000101*jc) + math.Sin(toRad(3*ganom))*0.000289
	truelong := glong + eqCtr
	app := truelong - 0.00569 - 0.00478*math.Sin(toRad(125.04-1934.136*jc))
	decl := toDeg(math.Asin(math.Sin(toRad(obliq)) * math.Sin(toRad(app))))
	HA := toDeg(math.Acos(math.Cos(toRad(90.833))/(math.Cos(toRad(lat))*math.Cos(toRad(decl))) -
		math.Tan(toRad(lat))*math.Tan(toRad(decl))))
	return time.Duration(HA*4) * time.Minute
}


func Sunrise(lat, long float64, date time.Time) time.Time {
		return SolarNoon(long, date).Add( -HourAngle(lat, date))

}

func Sunset(lat, long float64, date time.Time) time.Time {
		return SolarNoon(long, date).Add(HourAngle(lat, date))
}

func PartOfDay(lat, long float64, date time.Time) string {
	var snoon, srise, sset time.Time
	var ha time.Duration
	if snoon = gt.SolarNoon; snoon.IsZero() {
		snoon = SolarNoon(long, date)
	}
	if ha = gt.HourAngle; ha == 0 {
		ha =HourAngle(lat, date)
	}
	if srise = gt.Sunrise; srise.IsZero() {
		srise = snoon.Add(-ha)
	}
	if sset = gt.Sunset; sset.IsZero() {
		sset = snoon.Add(ha)
	}

	if date.After(sset.Add(time.Minute * 60)) {
		return night
	}
	if date.After(sset) {
		return evening
	}
	if date.After(sset.Add(-time.Minute*20)){
		return sunset
	}
	if date.After(snoon.Add(-time.Minute * 10)) && date.Before(snoon.Add(time.Minute * 10)) {
		return noon
	}
	if date.After(srise.Add(time.Minute*60)) {
		return day
	}
	if date.After(srise) {
		return morning
	}
	if date.After(srise.Add(-time.Minute*20)) {
		return sunrise
	}
	return night
}

func truncTime(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0,0,0,0, date.Location())
}

func toDeg(rad float64) float64 {
	return rad * 180 / math.Pi
}

func toRad(deg float64) float64 {
	return deg * math.Pi / 180
}
