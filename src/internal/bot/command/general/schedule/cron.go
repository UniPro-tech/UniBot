package schedule

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unibot/internal/scheduler"
)

type scheduleSpec interface {
	Next(after time.Time) time.Time
}

type intervalUnit int

const (
	unitMinutes intervalUnit = iota
	unitHours
	unitDays
	unitWeeks
	unitMonths
	unitYears
)

type intervalSchedule struct {
	interval int
	unit     intervalUnit
}

type dailyAtSchedule struct {
	hour int
	min  int
}

type weeklyAtSchedule struct {
	weekday time.Weekday
	hour    int
	min     int
}

func (s intervalSchedule) Next(after time.Time) time.Time {
	base := normalizeTime(after)
	var next time.Time

	switch s.unit {
	case unitMinutes:
		next = base.Add(time.Duration(s.interval) * time.Minute)
	case unitHours:
		next = base.Add(time.Duration(s.interval) * time.Hour)
	case unitDays:
		next = base.AddDate(0, 0, s.interval)
	case unitWeeks:
		next = base.AddDate(0, 0, 7*s.interval)
	case unitMonths:
		next = base.AddDate(0, s.interval, 0)
	case unitYears:
		next = base.AddDate(s.interval, 0, 0)
	}

	if !next.After(after) {
		next = next.Add(time.Minute)
	}

	return next
}

func (s dailyAtSchedule) Next(after time.Time) time.Time {
	base := after.In(scheduler.JST())
	candidate := time.Date(base.Year(), base.Month(), base.Day(), s.hour, s.min, 0, 0, base.Location())
	if !candidate.After(after) {
		candidate = candidate.AddDate(0, 0, 1)
	}
	return candidate
}

func (s weeklyAtSchedule) Next(after time.Time) time.Time {
	base := after.In(scheduler.JST())
	daysUntil := (int(s.weekday) - int(base.Weekday()) + 7) % 7
	candidate := time.Date(base.Year(), base.Month(), base.Day(), s.hour, s.min, 0, 0, base.Location()).AddDate(0, 0, daysUntil)
	if !candidate.After(after) {
		candidate = candidate.AddDate(0, 0, 7)
	}
	return candidate
}

func normalizeTime(t time.Time) time.Time {
	loc := scheduler.JST()
	return time.Date(t.In(loc).Year(), t.In(loc).Month(), t.In(loc).Day(), t.In(loc).Hour(), t.In(loc).Minute(), 0, 0, loc)
}

// convertToCron は自然言語の繰り返し指定をcronに変換する
func convertToCron(input string) (string, error) {
	text := preprocessScheduleText(input)
	spec, err := parseScheduleText(text)
	if err != nil {
		return "", err
	}

	now := time.Now().In(scheduler.JST())
	first := spec.Next(now)
	second := spec.Next(first)

	diffMinutes := int(second.Sub(first).Minutes())
	if diffMinutes <= 0 {
		return "", errors.New("invalid schedule")
	}

	min := first.Minute()
	hour := first.Hour()
	date := first.Day()
	month := int(first.Month())
	weekDay := int(first.Weekday())

	switch {
	case diffMinutes < 60:
		return fmt.Sprintf("*/%d * * * *", diffMinutes), nil
	case diffMinutes%60 == 0 && diffMinutes < 1440:
		hours := diffMinutes / 60
		return fmt.Sprintf("%d */%d * * *", min, hours), nil
	case diffMinutes >= 1440 && diffMinutes < 10080:
		return fmt.Sprintf("%d %d * * *", min, hour), nil
	case diffMinutes >= 10080 && diffMinutes < 40320:
		return fmt.Sprintf("%d %d * * %d", min, hour, weekDay), nil
	case diffMinutes >= 40320 && diffMinutes < 525600:
		return fmt.Sprintf("%d %d %d * *", min, hour, date), nil
	case diffMinutes >= 525600:
		return fmt.Sprintf("%d %d %d %d *", min, hour, date, month), nil
	default:
		return "", errors.New("invalid schedule")
	}
}

func preprocessScheduleText(input string) string {
	text := strings.TrimSpace(input)
	if text == "" {
		return ""
	}

	reEveryDay := regexp.MustCompile(`(?i)^\s*every\s+day\s*$`)
	if reEveryDay.MatchString(text) {
		text = "every day at 9:00 am"
	}

	reEveryDayAny := regexp.MustCompile(`(?i)\bevery\s+day\b`)
	if reEveryDayAny.MatchString(text) && !reEveryDay.MatchString(text) {
		text = strings.TrimSpace(reEveryDayAny.ReplaceAllString(text, ""))
	}

	reZeroHour := regexp.MustCompile(`\b0:([0-5][0-9])\b`)
	text = reZeroHour.ReplaceAllString(text, "12:$1")

	reTwelveWithoutAmPm := regexp.MustCompile(`\b12:([0-5][0-9])\b(?!\s?(am|pm))`)
	text = reTwelveWithoutAmPm.ReplaceAllString(text, "12:$1 am")

	re24Hour := regexp.MustCompile(`\b([1][3-9]|2[0-3]):([0-5][0-9])\b`)
	text = re24Hour.ReplaceAllStringFunc(text, func(match string) string {
		parts := strings.Split(match, ":")
		if len(parts) != 2 {
			return match
		}
		hour, err := strconv.Atoi(parts[0])
		if err != nil {
			return match
		}
		minute := parts[1]
		ampmHour := hour - 12
		period := "pm"
		if hour < 12 {
			ampmHour = hour
			period = "am"
		}
		return fmt.Sprintf("%d:%s %s", ampmHour, minute, period)
	})

	return strings.TrimSpace(text)
}

func parseScheduleText(text string) (scheduleSpec, error) {
	if text == "" {
		return nil, errors.New("empty schedule")
	}

	reAt := regexp.MustCompile(`(?i)^at\s+([0-9]{1,2}:[0-9]{2}(?:\s*(?:am|pm))?)$`)
	if match := reAt.FindStringSubmatch(text); match != nil {
		hour, min, err := parseTime(match[1])
		if err != nil {
			return nil, err
		}
		return dailyAtSchedule{hour: hour, min: min}, nil
	}

	reDaily := regexp.MustCompile(`(?i)^every\s+day\s+at\s+([0-9]{1,2}:[0-9]{2}(?:\s*(?:am|pm))?)$`)
	if match := reDaily.FindStringSubmatch(text); match != nil {
		hour, min, err := parseTime(match[1])
		if err != nil {
			return nil, err
		}
		return dailyAtSchedule{hour: hour, min: min}, nil
	}

	reWeekly := regexp.MustCompile(`(?i)^every\s+(monday|tuesday|wednesday|thursday|friday|saturday|sunday)\s+at\s+([0-9]{1,2}:[0-9]{2}(?:\s*(?:am|pm))?)$`)
	if match := reWeekly.FindStringSubmatch(text); match != nil {
		weekday := parseWeekday(match[1])
		hour, min, err := parseTime(match[2])
		if err != nil {
			return nil, err
		}
		return weeklyAtSchedule{weekday: weekday, hour: hour, min: min}, nil
	}

	reEvery := regexp.MustCompile(`(?i)^every\s+(\d+)\s+(minute|minutes|mins|hour|hours|day|days|week|weeks|month|months|year|years)$`)
	if match := reEvery.FindStringSubmatch(text); match != nil {
		interval, err := strconv.Atoi(match[1])
		if err != nil || interval <= 0 {
			return nil, errors.New("invalid interval")
		}

		switch strings.ToLower(match[2]) {
		case "minute", "minutes", "mins":
			return intervalSchedule{interval: interval, unit: unitMinutes}, nil
		case "hour", "hours":
			return intervalSchedule{interval: interval, unit: unitHours}, nil
		case "day", "days":
			return intervalSchedule{interval: interval, unit: unitDays}, nil
		case "week", "weeks":
			return intervalSchedule{interval: interval, unit: unitWeeks}, nil
		case "month", "months":
			return intervalSchedule{interval: interval, unit: unitMonths}, nil
		case "year", "years":
			return intervalSchedule{interval: interval, unit: unitYears}, nil
		}
	}

	return nil, errors.New("invalid schedule format")
}

func parseTime(text string) (int, int, error) {
	re := regexp.MustCompile(`(?i)^(\d{1,2}):(\d{2})\s*(am|pm)?$`)
	match := re.FindStringSubmatch(strings.TrimSpace(text))
	if match == nil {
		return 0, 0, errors.New("invalid time")
	}

	hour, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, 0, errors.New("invalid time")
	}

	min, err := strconv.Atoi(match[2])
	if err != nil {
		return 0, 0, errors.New("invalid time")
	}

	if min < 0 || min > 59 {
		return 0, 0, errors.New("invalid time")
	}

	ampm := strings.ToLower(match[3])
	if ampm != "" {
		if hour < 1 || hour > 12 {
			return 0, 0, errors.New("invalid time")
		}
		if ampm == "am" {
			if hour == 12 {
				hour = 0
			}
		}
		if ampm == "pm" {
			if hour != 12 {
				hour += 12
			}
		}
	} else {
		if hour < 0 || hour > 23 {
			return 0, 0, errors.New("invalid time")
		}
	}

	return hour, min, nil
}

func parseWeekday(text string) time.Weekday {
	switch strings.ToLower(text) {
	case "monday":
		return time.Monday
	case "tuesday":
		return time.Tuesday
	case "wednesday":
		return time.Wednesday
	case "thursday":
		return time.Thursday
	case "friday":
		return time.Friday
	case "saturday":
		return time.Saturday
	case "sunday":
		return time.Sunday
	default:
		return time.Sunday
	}
}

func describeCron(cronText string) string {
	fields := strings.Fields(cronText)
	if len(fields) != 5 {
		return cronText
	}

	min := fields[0]
	hour := fields[1]
	day := fields[2]
	month := fields[3]
	week := fields[4]

	if strings.HasPrefix(min, "*/") && hour == "*" && day == "*" && month == "*" && week == "*" {
		return fmt.Sprintf("Every %s minutes", strings.TrimPrefix(min, "*/"))
	}

	if strings.HasPrefix(hour, "*/") && day == "*" && month == "*" && week == "*" {
		return fmt.Sprintf("Every %s hours at minute %s", strings.TrimPrefix(hour, "*/"), min)
	}

	if day == "*" && month == "*" && week == "*" {
		return fmt.Sprintf("At %s:%s every day", padHour(hour), padMinute(min))
	}

	if day == "*" && month == "*" && week != "*" {
		return fmt.Sprintf("At %s:%s, only on %s", padHour(hour), padMinute(min), weekDayName(week))
	}

	if day != "*" && month == "*" && week == "*" {
		return fmt.Sprintf("At %s:%s, on day %s of the month", padHour(hour), padMinute(min), day)
	}

	if day != "*" && month != "*" && week == "*" {
		return fmt.Sprintf("At %s:%s, on day %s of month %s", padHour(hour), padMinute(min), day, month)
	}

	return cronText
}

func padHour(hour string) string {
	if len(hour) == 1 {
		return "0" + hour
	}
	return hour
}

func padMinute(min string) string {
	if len(min) == 1 {
		return "0" + min
	}
	return min
}

func weekDayName(week string) string {
	week = strings.TrimSpace(week)
	switch week {
	case "0", "7":
		return "Sunday"
	case "1":
		return "Monday"
	case "2":
		return "Tuesday"
	case "3":
		return "Wednesday"
	case "4":
		return "Thursday"
	case "5":
		return "Friday"
	case "6":
		return "Saturday"
	default:
		return week
	}
}
