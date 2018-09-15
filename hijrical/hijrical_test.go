package hijrical

import (
	"testing"
	"time"
)

func TestConversion(t *testing.T) {
	hijriDate := SimpleDate(time.Date(1981, 10, 14, 0, 0, 0, 0, time.UTC))
	if hijriDate != "Arbaa, 16 Dhulhijja 1401H" {
		t.Errorf("Incorrect Hijri date generated. Got: %s, want: %s",
			hijriDate,
			"Arbaa, 16 Dhulhijja 1401H")
	}
}
