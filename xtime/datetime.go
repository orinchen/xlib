package xtime

import (
	"encoding/json"
	"fmt"
	"time"
)

type DateTime time.Time

func (dt *DateTime) MarshalJSON() ([]byte, error) {
	if dt == nil {
		return []byte(""), nil
	}
	var stamp = fmt.Sprintf("\"%s\"", time.Time(*dt).Format(time.DateTime))
	return []byte(stamp), nil
}

// UnmarshalJSON 实现JSON反序列化方法
func (dt *DateTime) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err = json.Unmarshal(data, &s); err != nil {
		return err
	}
	for _, layout := range layouts {
		t, err := time.ParseInLocation(layout, s, time.Local)
		if err == nil {
			*dt = (DateTime)(t)
			return nil
		}
	}
	return err
}

func (dt *DateTime) ToTime() time.Time {
	return time.Time(*dt)
}

func (dt *DateTime) ToTimeP() *time.Time {
	if dt == nil {
		return nil
	}
	var temp = dt.ToTime()
	return &temp
}
