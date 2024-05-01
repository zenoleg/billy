package rating

import "time"

const (
	TopDay PeriodPreset = iota
	TopWeek
	TopMonth
	TopEver
)

type (
	PeriodPreset uint8
)

func (p PeriodPreset) MakeFromAndTo(now time.Time) (time.Time, time.Time) {
	switch p {
	case TopDay:
		return now.Add(-time.Hour * 24), now
	case TopWeek:
		return now.Add(-time.Hour * 24 * 7), now
	case TopMonth:
		return now.Add(-time.Hour * 24 * 30), now
	case TopEver:
		return time.Unix(0, 0), now
	}

	return time.Time{}, time.Time{}
}
