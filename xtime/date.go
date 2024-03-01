package xtime

import (
	"encoding/json"
	"fmt"
	"time"
)

type Date time.Time

func (dt *Date) MarshalJSON() ([]byte, error) {
	if dt == nil {
		return []byte(""), nil
	}
	var stamp = fmt.Sprintf("\"%s\"", time.Time(*dt).Format(time.DateOnly))
	return []byte(stamp), nil
}

// UnmarshalJSON 实现JSON反序列化方法
func (dt *Date) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err = json.Unmarshal(data, &s); err != nil {
		return err
	}
	for _, layout := range layouts {
		var t time.Time
		t, err = time.ParseInLocation(layout, s, time.Local)
		if err == nil {
			*dt = (Date)(t)
			return nil
		}
	}
	return err
}

func (dt *Date) ToTime() time.Time {
	return time.Time(*dt)
}

func (dt *Date) ToTimeP() *time.Time {
	if dt == nil {
		return nil
	}
	var temp = dt.ToTime()
	return &temp
}
