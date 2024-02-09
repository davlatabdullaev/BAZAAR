package check

import (
	"fmt"
	"time"
)

func CalculateAge(birthDate string) int {
	layout := "2006-01-02"
	birthday, err := time.Parse(layout, birthDate)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return 0
	}

	now := time.Now()
	age := now.Year() - birthday.Year()

	if now.YearDay() < birthday.YearDay() {
		age--
	}

	return age
}
