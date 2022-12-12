package time

import "time"

var (
	TimezoneVietnam *time.Location
)

const (
	TimezoneVietnamString = "Asia/Ho_Chi_Minh"
)

func init() {
	TimezoneVietnam, _ = time.LoadLocation(TimezoneVietnamString)
}

// FromWooDateString returns time from format dd/mm/yyyy based on Vietnam timezone
func FromWooDateString(date string) (time.Time, error) {
	return time.ParseInLocation("02/01/2006", date, TimezoneVietnam)
}
