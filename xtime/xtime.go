/*
 *  Author: Orin Chen
 *   Email: orinchen@gmail.com
 *    Time: 2020/9/24 3:39 下午
 * Project: mars
 *    File: time.go
 *     IDE: GoLand
 */

package xtime

import (
	"fmt"
	"time"
)

const (
	TimeOnlyCN       = "15时04分05秒"
	DateTimeCN       = "2006年01月02日 15时04分05秒"
	DateOnlyCN       = "2006年01月02日"
	RFC3339WithoutTZ = "2006-01-02T15:04:05"
	RFC3399Msec      = "2006-01-02T15:04:05.000Z07:00"
)

var layouts = []string{
	time.DateTime,
	time.DateOnly,
	DateTimeCN,
	DateOnlyCN,
	RFC3339WithoutTZ,
	time.RFC3339,
	RFC3399Msec,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339Nano,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
	TimeOnlyCN,
	time.TimeOnly,
	time.Layout,
	"2006-01-02 15:04:05Z07:00",
	"02 Jan 06 15:04 MST",
	"02 Jan 2006",
	"2006-01-02 15:04:05 -07:00",
	"2006-01-02 15:04:05 -0700",
	"20060102150405",
}

func AutoParse(value string) (t *time.Time, err error) {
	return AutoParseInLocation(value, time.Now().Location())
}

func AutoParseInLocation(value string, loc *time.Location) (t *time.Time, err error) {
	var tt time.Time
	for _, layout := range layouts {
		tt, err = time.ParseInLocation(layout, value, loc)
		if err == nil {
			t = &tt
			return
		}
	}
	return t, fmt.Errorf("%s, 该时间格式无法解析", value)
}

type HalfYear int

const (
	FirstHalfYear HalfYear = iota
	SecondHalfYear
)

type Quarter int

const (
	FirstQuarter Quarter = iota + 1
	SecondQuarter
	ThirdQuarter
	FourthQuarter
)

func GetDateTimeSpanByStartEnd(start, end *DateTime) (s *time.Time, e *time.Time) {
	if start != nil {
		t := time.Time(*start)
		temp := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		s = &temp
	}
	if end != nil {
		t := time.Time(*end)
		temp := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999, t.Location())
		e = &temp
	}
	return
}

func GetTimeSpanByStartEnd(start, end *time.Time) (s *time.Time, e *time.Time) {
	if start != nil {
		temp := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
		s = &temp
	}
	if end != nil {
		temp := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999, end.Location())
		e = &temp
	}
	return
}

func GetDaySpanByTime(t time.Time) (start, end time.Time) {
	return GetDaySpan(t.Year(), t.Month(), t.Day(), t.Location())
}

func GetNearDaysSpanByTime(t time.Time, nearDays int) (start, end time.Time) {
	return GetNearDaysSpan(t.Year(), t.Month(), t.Day(), t.Location(), nearDays)
}

func GetNearMonthsSpanByTime(t time.Time, nearMonths int) (start, end time.Time) {
	return GetNearMonthsSpan(t.Year(), t.Month(), t.Day(), t.Location(), nearMonths)
}

func GetDaySpan(year int, month time.Month, day int, loc *time.Location) (start, end time.Time) {
	start = time.Date(year, month, day, 0, 0, 0, 0, loc)
	end = start.Add(time.Hour*24 - time.Nanosecond)
	return
}

func GetNearDaysSpan(year int, month time.Month, day int, loc *time.Location, nearDays int) (start, end time.Time) {
	start = time.Date(year, month, day, 0, 0, 0, 0, loc).AddDate(0, 0, nearDays)
	end = time.Date(year, month, day, 23, 59, 59, 0, loc)
	return
}

func GetNearMonthsSpan(year int, month time.Month, day int, loc *time.Location, nearMonths int) (start, end time.Time) {
	start = time.Date(year, month, day, 0, 0, 0, 0, loc).AddDate(0, nearMonths, 0)
	end = time.Date(year, month, day, 23, 59, 59, 0, loc)
	return
}

func GetWeekSpan(year int, week int, loc *time.Location) (start, end time.Time) {
	//yearDay := start.YearDay()
	yearFirstDay := time.Date(year, 1, 1, 0, 0, 0, 0, loc)
	firstDayInWeek := int(yearFirstDay.Weekday())

	//今年第一周有几天
	firstWeekDays := 1
	if firstDayInWeek != 0 {
		firstWeekDays = 7 - firstDayInWeek + 1
	}

	if week == 1 {
		start = yearFirstDay
		end = yearFirstDay.AddDate(0, 0, firstWeekDays).Add(-time.Nanosecond)
		return
	}

	startDays := 7*(week-2) + 1 + firstWeekDays // 本周第一天距离本年第一天的天数

	start = yearFirstDay.AddDate(0, 0, startDays)
	end = yearFirstDay.AddDate(0, 0, startDays+7).Add(-time.Microsecond)
	return
}

func GetHalfYearSpan(year int, halfYear HalfYear, loc *time.Location) (start, end time.Time) {
	if halfYear == FirstHalfYear {
		start = time.Date(year, 1, 1, 0, 0, 0, 0, loc)
		end = time.Date(year, 6, 30, 23, 59, 59, 99999, loc)
	} else if halfYear == SecondHalfYear {
		start = time.Date(year, 7, 1, 0, 0, 0, 0, loc)
		end = time.Date(year, 12, 31, 23, 59, 59, 99999, loc)
	}
	return
}

func GetMonthSpanByTime(t time.Time) (start, end time.Time) {
	return GetMonthSpan(t.Year(), t.Month(), t.Location())
}

func GetMonthSpan(year int, month time.Month, loc *time.Location) (start, end time.Time) {
	start = time.Date(year, month, 1, 0, 0, 0, 0, loc)
	end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return
}

func GetYearSpan(year int, location *time.Location) (start, end time.Time) {
	start = time.Date(year, 1, 1, 0, 0, 0, 0, location)
	end = start.AddDate(1, 0, 0).Add(-time.Nanosecond)
	return
}

func GetQuarterSpan(year int, quarter Quarter, loc *time.Location) (start, end time.Time) {
	switch quarter {
	case FirstQuarter:
		start = time.Date(year, 1, 1, 0, 0, 0, 0, loc)
		end = start.AddDate(0, 3, 0).Add(-time.Nanosecond)
	case SecondQuarter:
		start = time.Date(year, 4, 1, 0, 0, 0, 0, loc)
		end = start.AddDate(0, 3, 0).Add(-time.Nanosecond)
	case ThirdQuarter:
		start = time.Date(year, 7, 1, 0, 0, 0, 0, loc)
		end = start.AddDate(0, 3, 0).Add(-time.Nanosecond)
	case FourthQuarter:
		start = time.Date(year, 10, 1, 0, 0, 0, 0, loc)
		end = start.AddDate(0, 3, 0).Add(-time.Nanosecond)
	}
	return
}
