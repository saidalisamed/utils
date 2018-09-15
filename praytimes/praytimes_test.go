package praytimes

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
)

func TestJafariFajr(t *testing.T) {
	times := generatePrayerTimes()
	if !strings.Contains(times, "Fajr: 03:59AM") {
		t.Errorf(`Generated prayer times for Sydney Australia on 14th October 1981 was incorrect, got:
%s 
want:
Fajr: 03:59AM
Sunrise: 05:15AM
Dhuhr: 11:41AM
Sunset: 06:08PM
Maghrib: 06:23PM`, times)
	}
}

func generatePrayerTimes() string {
	loc, _ := time.LoadLocation("Australia/Sydney")
	p := Custom(ConventionJafari,
		0,
		AsrFactorStandard,
		HighlatMethodAngleBased,
		-33.8688,
		151.2093,
		time.Date(1981, 10, 14, 0, 0, 0, 0, loc),
		"Australia/Sydney",
		[9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0})

	timeNames := []string{"Fajr", "Sunrise", "Dhuhr", "Sunset", "Maghrib"}
	var data string
	for _, timeName := range timeNames {
		pTime := p.Midnight
		if timeName == "Fajr" {
			pTime = p.Fajr
		} else if timeName == "Sunrise" {
			pTime = p.Sunrise
		} else if timeName == "Dhuhr" {
			pTime = p.Dhuhr
		} else if timeName == "Sunset" {
			pTime = p.Sunset
		} else if timeName == "Maghrib" {
			pTime = p.Maghrib
		}

		data += fmt.Sprintf("%s: %s\n", timeName, time12Hours(pTime.Hours, pTime.Minutes))
	}

	return data
}

func time12Hours(hours float64, minutes float64) string {
	t, err := time.Parse("15:04", fmt.Sprintf("%02d:%02d", int64(hours), int64(minutes)))
	if err != nil {
		log.Println(err)
	}
	return t.Format("03:04PM")
}
