//--------------------- Copyright Block ----------------------
/*

PrayTimes.js: Prayer Times Calculator (ver 2.3)
Copyright (C) 2007-2011 PrayTimes.org

Developer: Hamid Zarrabi-Zadeh
License: GNU LGPL v3.0

TERMS OF USE:
Permission is granted to use this code, with or
without modification, in any website or application
provided that credit is given to the original work
with a link back to PrayTimes.org.

This program is distributed in the hope that it will
be useful, but WITHOUT ANY WARRANTY.

PLEASE DO NOT REMOVE THIS COPYRIGHT BLOCK.

*/

//--------------------- Help and Manual ----------------------
/*

User's Manual:
http://praytimes.org/manual

Calculation Formulas:
http://praytimes.org/calculation

* Go package created by Said Ali Samed on 17/4/17.
*
*/

package praytimes

import (
	"math"
	"time"
)

// Conventions
const (
	ConventionJafari = iota
	ConventionKarachi
	ConventionISNA
	ConventionMWL
	ConventionMakkah
	ConventionEgypt
	ConventionTehran
	ConventionCustom
)

// Fajr angle
const (
	AngleFajr = iota
	AngleDhuhr
	AngleMaghrib
	AngleNAN
	AngleIsha
)

// High latitude modes
const (
	HighlatMethodNone = iota
	HighlatMethodNightMiddle
	HighlatMethodOneSeventh
	HighlatMethodAngleBased
)

// Angle of the sun
var sunAngle = [8][5]float64{
	{16, 0, 4, 0, 14},     // Jafari
	{18, 1, 0, 0, 18},     // Karachi
	{15, 1, 0, 0, 15},     // ISNA
	{18, 1, 0, 0, 17},     // MWL
	{18.5, 1, 0, 1, 90},   // Makkah
	{19.5, 1, 0, 0, 17.5}, // Egypt
	{17.7, 0, 4.5, 0, 14}, // Tehran
	{18, 1, 0, 0, 17},     // Custom
}

// Times of the day
const (
	TimeImsak = iota
	TimeFajr
	TimeSunrise
	TimeDhuhr
	TimeAsr
	TimeSunset
	TimeMaghrib
	TimeIsha
	TimeMidnight
)

// Hours and minutes
const (
	Hours = iota
	Minutes
)

// AsrFactorStandard Asr factor for all madhabs other than Hanafi
const AsrFactorStandard = 1

// AsrFactorHanafi Hanafi Asr factor
const AsrFactorHanafi = 2

var defaultTimes = []float64{5, 5, 6, 12, 13, 18, 18, 18, 18} // default times

// HoursMinutes Hours and minutes data structure
type HoursMinutes struct {
	Hours   float64
	Minutes float64
}

// PrayTimes Prayer times data structure
type PrayTimes struct {
	Imsak    HoursMinutes
	Fajr     HoursMinutes
	Sunrise  HoursMinutes
	Dhuhr    HoursMinutes
	Asr      HoursMinutes
	Sunset   HoursMinutes
	Maghrib  HoursMinutes
	Isha     HoursMinutes
	Midnight HoursMinutes
}

// Config Configuration data structure
type Config struct {
	Convention    int
	ImsakMinutes  int
	DhuhrMinutes  int
	AsrFactor     int
	HighLatMethod int
	Latitude      float64
	Longitude     float64
	Date          time.Time
	TimeZone      string
	Offsets       [9]float64
}

// Default default settings
func Default() PrayTimes {
	cfg := Config{}
	cfg.Date = time.Now()
	cfg.TimeZone = "Local"
	cfg.AsrFactor = AsrFactorStandard
	cfg.Convention = ConventionJafari
	cfg.DhuhrMinutes = 0
	cfg.HighLatMethod = HighlatMethodAngleBased
	cfg.ImsakMinutes = 10
	cfg.Latitude = -33.7640187
	cfg.Longitude = 150.8202351
	cfg.Offsets = [9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0}

	return computeTimes(&cfg)
}

// Custom customised configuration
func Custom(convention int, dhuhrMinutes int, asrFactor int, highLatMethod int, latitude float64, longitude float64,
	date time.Time, timeZone string, offsets [9]float64) PrayTimes {
	cfg := Config{}
	cfg.Date = date
	cfg.TimeZone = timeZone
	cfg.AsrFactor = asrFactor
	cfg.Convention = convention
	cfg.DhuhrMinutes = dhuhrMinutes
	cfg.HighLatMethod = highLatMethod
	cfg.ImsakMinutes = 10
	cfg.Latitude = latitude
	cfg.Longitude = longitude
	cfg.Offsets = offsets

	return computeTimes(&cfg)
}

// ---------- Trigonometric functions ------------

// Range reduce angle in degree
func fixAngle(a float64) float64 {
	a = a - (360 * (math.Floor(a / 360.0)))
	if a < 0 {
		a += 360
	}
	return a
}

// Range reduce hours to 0..23
func fixHour(a float64) float64 {
	a = a - 24.0*math.Floor(a/24.0)
	if a < 0 {
		a += 24
	}
	return a
}

// Radian to degree
func radiansToDegrees(alpha float64) float64 {
	return (alpha * 180.0) / math.Pi
}

// Degree to radian
func degreesToRadians(alpha float64) float64 {
	return (alpha * math.Pi) / 180.0
}

// Degree to sin
func degreeToSin(d float64) float64 {
	return math.Sin(degreesToRadians(d))
}

// Degree to cos
func degreeToCos(d float64) float64 {
	return math.Cos(degreesToRadians(d))
}

// Degree to tan
func degreeToTan(d float64) float64 {
	return math.Tan(degreesToRadians(d))
}

// Degree arcsin
func degreeArcsin(x float64) float64 {
	return radiansToDegrees(math.Asin(x))
}

// Degree arccos
func degreeArccos(x float64) float64 {
	return radiansToDegrees(math.Acos(x))
}

// Degree arctan
func degreeArctan(x float64) float64 {
	return radiansToDegrees(math.Atan(x))
}

// Degree arctan2
func degreeArctan2(y float64, x float64) float64 {
	return radiansToDegrees(math.Atan2(y, x))
}

// Degree arccot
func degreeArccot(x float64) float64 {
	return radiansToDegrees(math.Atan2(1.0, x))
}

// ---------- Julian date functions --------------

// Calculate julian date from a calendar date
func julianDate(year float64, month float64, day float64) float64 {
	if month <= 2 {
		year--
		month += 12
	}

	A := math.Floor(year / 100.0)
	B := 2 - A + math.Floor(A/4.0)
	y := math.Floor(365.25 * (year + 4716))
	m := math.Floor(30.6001 * (month + 1))
	return (y + m + day + B) - 1524
}

// ---------- Calculation functions --------------

/*
References:
http://www.ummah.net/astronomy/saltime
http://aa.usno.navy.mil/faq/docs/SunApprox.html
Compute declination angle of sun and equation of time
*/
func sunPosition(jd float64) [2]float64 {
	D := jd - 2451545
	g := fixAngle(357.529 + 0.98560028*D)
	q := fixAngle(280.459 + 0.98564736*D)
	L := fixAngle(q + (1.915 * degreeToSin(g)) + (0.020 * degreeToSin(2*g)))
	//R := 1.00014 - 0.01671 * degreeToCos(g) - 0.00014 * degreeToCos(2 * g)
	e := 23.439 - (0.00000036 * D)
	d := degreeArcsin(degreeToSin(e) * degreeToSin(L))
	RA := (degreeArctan2(degreeToCos(e)*degreeToSin(L), degreeToCos(L))) / 15.0
	RA = fixHour(RA)
	EqT := q/15.0 - RA
	sPosition := [2]float64{d, EqT}
	return sPosition
}

// Compute equation of time
func equationOfTime(jd float64) float64 {
	return sunPosition(jd)[1]
}

// Compute declination angle of sun
func sunDeclination(jd float64) float64 {
	return sunPosition(jd)[0]
}

// Compute midday (Dhuhr, Zawal) time
func computeMidday(t float64, jDay float64) float64 {
	T := equationOfTime(jDay + t)
	return fixHour(12 - T)
}

// Compute time for the  given angle G
func computeTime(G float64, t float64, jDay float64, latitude float64) float64 {
	D := sunDeclination(jDay + t)
	Z := computeMidday(t, jDay)
	beg := -degreeToSin(G) - degreeToSin(D)*degreeToSin(latitude)
	mid := degreeToCos(D) * degreeToCos(latitude)
	V := degreeArccos(beg/mid) / 15.0
	if G > 90 {
		return Z + -V
	}

	return Z + V
}

// Compute the time of Asr (Standard: factor=1, Hanafi: factor=2)
func computeAsr(factor int, t float64, jDay float64, latitude float64) float64 {
	D := sunDeclination(jDay + t)
	G := -degreeArccot(float64(factor) + degreeToTan(math.Abs(latitude-D)))
	return computeTime(G, t, jDay, latitude)
}

func computeTimes(cfg *Config) PrayTimes {
	t := cfg.Date
	loc, _ := time.LoadLocation(cfg.TimeZone)
	_, offset := t.In(loc).Zone()
	tzOffset := float64(offset / 60 / 60)

	jDay := julianDate(float64(t.Year()), float64(t.Month()), float64(t.Day()))
	lonDiff := cfg.Longitude / (15.0 * 24.0)
	jDay = jDay - lonDiff

	fajrTime := computeTime(180-sunAngle[cfg.Convention][AngleFajr],
		dayPortion(defaultTimes[TimeFajr]), jDay, cfg.Latitude)
	sunriseTime := computeTime(180-0.833, dayPortion(defaultTimes[TimeSunrise]), jDay, cfg.Latitude)
	dhuhrTime := computeMidday(dayPortion(defaultTimes[TimeDhuhr]), jDay)
	asrTime := computeAsr(cfg.AsrFactor, dayPortion(defaultTimes[TimeAsr]), jDay, cfg.Latitude)
	sunsetTime := computeTime(0.833, dayPortion(defaultTimes[TimeSunset]), jDay, cfg.Latitude)
	maghribTime := computeTime(sunAngle[cfg.Convention][AngleMaghrib],
		dayPortion(defaultTimes[TimeMaghrib]), jDay, cfg.Latitude)
	ishaTime := computeTime(sunAngle[cfg.Convention][AngleIsha],
		dayPortion(defaultTimes[TimeIsha]), jDay, cfg.Latitude)

	computedTimes := []float64{fajrTime, fajrTime, sunriseTime, dhuhrTime,
		asrTime, sunsetTime, maghribTime, ishaTime, ishaTime}

	computedTimes = adjustTimes(cfg, computedTimes, tzOffset, cfg.Longitude)
	computedTimes = tuneTimes(computedTimes, cfg.Offsets)

	// Add midnight time
	if cfg.Convention == ConventionJafari {
		computedTimes[TimeMidnight] = computedTimes[TimeSunset] +
			timeDiff(computedTimes[TimeSunset], computedTimes[TimeFajr])/2
	} else {
		computedTimes[TimeMidnight] = computedTimes[TimeSunset] +
			timeDiff(computedTimes[TimeSunset], computedTimes[TimeSunrise])/2
	}

	var hoursMinutes [2]float64
	times := PrayTimes{}

	hoursMinutes = convertToHoursMinutes(computedTimes[TimeImsak])
	times.Imsak.Hours = hoursMinutes[Hours]
	times.Imsak.Minutes = hoursMinutes[Minutes]

	hoursMinutes = convertToHoursMinutes(computedTimes[TimeFajr])
	times.Fajr.Hours = hoursMinutes[Hours]
	times.Fajr.Minutes = hoursMinutes[Minutes]

	hoursMinutes = convertToHoursMinutes(computedTimes[TimeSunrise])
	times.Sunrise.Hours = hoursMinutes[Hours]
	times.Sunrise.Minutes = hoursMinutes[Minutes]

	hoursMinutes = convertToHoursMinutes(computedTimes[TimeDhuhr])
	times.Dhuhr.Hours = hoursMinutes[Hours]
	times.Dhuhr.Minutes = hoursMinutes[Minutes]

	hoursMinutes = convertToHoursMinutes(computedTimes[TimeAsr])
	times.Asr.Hours = hoursMinutes[Hours]
	times.Asr.Minutes = hoursMinutes[Minutes]

	hoursMinutes = convertToHoursMinutes(computedTimes[TimeSunset])
	times.Sunset.Hours = hoursMinutes[Hours]
	times.Sunset.Minutes = hoursMinutes[Minutes]

	hoursMinutes = convertToHoursMinutes(computedTimes[TimeMaghrib])
	times.Maghrib.Hours = hoursMinutes[Hours]
	times.Maghrib.Minutes = hoursMinutes[Minutes]

	hoursMinutes = convertToHoursMinutes(computedTimes[TimeIsha])
	times.Isha.Hours = hoursMinutes[Hours]
	times.Isha.Minutes = hoursMinutes[Minutes]

	hoursMinutes = convertToHoursMinutes(computedTimes[TimeMidnight])
	times.Midnight.Hours = hoursMinutes[Hours]
	times.Midnight.Minutes = hoursMinutes[Minutes]

	return times
}

// ---------- Compute Prayer Times functions -----

// Night portion used for adjusting times in higher latitudes
func nightPortion(cfg *Config, angle float64) float64 {
	if cfg.HighLatMethod == HighlatMethodAngleBased {
		return angle / 60.0
	} else if cfg.HighLatMethod == HighlatMethodNightMiddle {
		return 0.5
	} else if cfg.HighLatMethod == HighlatMethodOneSeventh {
		return 0.14286
	} else {
		return 0
	}
}

func dayPortion(time float64) float64 {
	time /= 24
	return time
}

// Compute difference between two defaultTimes
func timeDiff(time1 float64, time2 float64) float64 {
	return fixHour(time2 - time1)
}

// Adjust fajr, maghrib and isha for locations in higher latitudes
func adjustHighLatTime(cfg *Config, timeName int, time float64, times []float64) float64 {
	nightTime := timeDiff(times[TimeSunset], times[TimeSunrise]) // sunset to sunrise

	if timeName == TimeFajr {
		// Adjust fajr
		fajrDiff := nightPortion(cfg, sunAngle[cfg.Convention][AngleFajr]) * nightTime
		if times[TimeFajr] < 1 || timeDiff(times[TimeFajr], times[TimeSunrise]) > fajrDiff {
			return times[TimeSunrise] - fajrDiff
		}
	} else if timeName == TimeMaghrib {
		// Adjust maghrib
		maghribAngle := 4.0
		if sunAngle[cfg.Convention][AngleDhuhr] == 0 {
			maghribAngle = sunAngle[cfg.Convention][AngleMaghrib]
		}
		maghribDiff := nightPortion(cfg, maghribAngle) * nightTime
		if times[TimeMaghrib] < 1 || timeDiff(times[TimeSunset], times[TimeMaghrib]) > maghribDiff {
			return times[TimeSunset] + maghribDiff
		}
	} else if timeName == TimeIsha {
		ishaAngle := 18.0
		if sunAngle[cfg.Convention][AngleNAN] == 0 {
			ishaAngle = sunAngle[cfg.Convention][AngleIsha]
		}
		ishaDiff := nightPortion(cfg, ishaAngle) * nightTime
		if times[TimeIsha] < 1 || timeDiff(times[TimeSunset], times[TimeIsha]) > ishaDiff {
			return times[TimeSunset] + ishaDiff
		}
	}

	return time
}

func adjustTimes(cfg *Config, computedTimes []float64, tzOffset float64, longitude float64) []float64 {
	for i := 0; i < len(computedTimes); i++ {
		computedTimes[i] += tzOffset - longitude/15

		if cfg.HighLatMethod != HighlatMethodNone {
			computedTimes[i] = adjustHighLatTime(cfg, i, computedTimes[i], computedTimes)
		}
	}

	// Imsak 10 minutes before fajr
	computedTimes[TimeImsak] = computedTimes[TimeFajr] - float64(cfg.ImsakMinutes)/60.0

	// Add dhuhr minutes if any
	computedTimes[TimeDhuhr] += float64(cfg.DhuhrMinutes) / 60.0

	if sunAngle[cfg.Convention][AngleDhuhr] == 1 {
		computedTimes[TimeMaghrib] = computedTimes[TimeSunset] + sunAngle[cfg.Convention][AngleMaghrib]/60
	}

	if sunAngle[cfg.Convention][AngleNAN] == 1 {
		computedTimes[TimeIsha] = computedTimes[TimeMaghrib] + sunAngle[cfg.Convention][AngleIsha]/60
	}

	return computedTimes
}

func tuneTimes(computedTimes []float64, offsets [9]float64) []float64 {
	for t := 0; t < len(computedTimes); t++ {
		computedTimes[t] = computedTimes[t] + offsets[t]/60.0
	}

	return computedTimes
}

func convertToHoursMinutes(time float64) [2]float64 {
	hoursMinute := [2]float64{0, 0}
	time = fixHour(time + 0.5/60.0) // add 0.5 minutes to round
	hours := math.Floor(time)
	minutes := math.Floor((time - hours) * 60.0)

	hoursMinute[Hours] = hours
	hoursMinute[Minutes] = minutes
	return hoursMinute
}
