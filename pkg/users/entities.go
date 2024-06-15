package users

import (
	"math"
	"time"
)

type User struct {
	// Username is unique
	Username string    `json:"name" db:"name"`
	DOB      time.Time `json:"dateOfBirth" db:"data_of_birth"`
}

func (u User) CalcDaysToBirthday() int {
	// Use server's local time.
	// NOTE: Might have some edge cases not handled as we are not storing date's timezone
	today := time.Now()

	if u.DOB.Month() == today.Month() && u.DOB.Day() == today.Day() {
		return 0
	}

	todayDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, nil)
	birthdayThisYear := time.Date(today.Year(), u.DOB.Month(), u.DOB.Day(), 0, 0, 0, 0, nil)

	// Birthday has not yet passed in the current year
	if birthdayThisYear.After(todayDate) {
		return int(math.Ceil(birthdayThisYear.Sub(today).Hours() / 24))
	}

	// Birthday has already passed in the current year,
	// we need to increase the year
	birthdayNextYear := time.Date(today.Year()+1, u.DOB.Month(), u.DOB.Day(), 0, 0, 0, 0, nil)
	return int(math.Ceil(birthdayNextYear.Sub(today).Hours() / 24))

}
