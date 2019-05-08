package notifications

import "time"

type TimeConverter struct {
	currentTime time.Time
}

const minutesInADay = 1440
const minutesInAnHour = 60

func NewTimeConverter() TimeConverter {
	location, _ := time.LoadLocation("Europe/Sofia")
	currentTime := time.Now().In(location)
	tc := TimeConverter{
		currentTime: currentTime,
	}
	return tc
}

//GetTimeInMinutes gets the current minute of the week [0, 10079]
func (tc TimeConverter) GetTimeInMinutes() int {
	day := tc.getCurrentDay()
	minutes := tc.getCurrentMinutesInDay()
	return (day-1)*minutesInADay + minutes
}

func (tc TimeConverter) getCurrentDay() int {
	return int(tc.currentTime.Weekday())
}

func (tc TimeConverter) getCurrentMinutesInDay() int {
	hours := tc.currentTime.Hour()
	minutes := tc.currentTime.Minute()
	return hours*minutesInAnHour + minutes
}
