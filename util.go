package schedules

import (
	"strings"
	"time"
)

// Parses a time range string from MicroStrategy to native time formats
//
// Example Time String:
//     11:00 am - 2:00 pm
func parseTimeRange(timeRange string) (start time.Time, end time.Time, err error) {
	timeRangeArray := strings.Split(timeRange, " - ")
	start, err = time.Parse("3:04 pm", timeRangeArray[0])
	end, err = time.Parse("3:04 pm", timeRangeArray[1])
	return
}
