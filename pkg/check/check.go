package check

import (
	"fmt"
	"math/rand"
	"time"
)

func TimeNow() time.Time {

	time := time.Now()

	return time
}

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

func GenerateBarCode() int {

	rand.Seed(TimeNow().UnixNano())

	min := 1000000000

	max := 9999999999

	random := rand.Intn(max-min) + min

	return random

}
