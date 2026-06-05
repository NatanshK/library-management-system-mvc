package utils

import (
	"math"
	"time"
)

func CalculateFine(dueDate time.Time, returnDate time.Time, dailyRate float64) (int, float64) {

	timeDifference := returnDate.Sub(dueDate)

	rawDays := timeDifference.Hours() / 24.0

	daysLate := int(math.Ceil(rawDays))

	if daysLate <= 0 {
		return 0, 0.0
	}

	totalFine := float64(daysLate) * dailyRate

	return daysLate, totalFine
}
