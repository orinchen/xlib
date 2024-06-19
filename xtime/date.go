package xtime

import (
	"encoding/json"
	"fmt"
	"github.com/orinchen/xlib/xstring"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"time"
)

type Date time.Time

func (d Date) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(time.Time(d))
}

func (d *Date) UnmarshalBSONValue(t bsontype.Type, value []byte) (err error) {
	if t != bson.TypeDateTime {
		return fmt.Errorf("invalid bson value type '%s'", t.String())
	}
	s, _, ok := bsoncore.ReadDateTime(value)
	if !ok {
		return fmt.Errorf("invalid bson string value")
	}

	*d = Date(time.UnixMilli(s))
	return
}

func (d Date) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(d).Format(time.DateOnly))
	return []byte(stamp), nil
}

// UnmarshalJSON 实现JSON反序列化方法
func (d *Date) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err = json.Unmarshal(data, &s); err != nil {
		return err
	}
	for _, layout := range layouts {
		var t time.Time
		t, err = time.ParseInLocation(layout, s, time.Local)
		if err == nil {
			*d = (Date)(t)
			return nil
		}
	}
	return err
}

func (d Date) MarshalText() (text []byte, err error) {
	var stamp = fmt.Sprintf("%s", time.Time(d).Format(time.DateOnly))
	return []byte(stamp), nil
}

func (d *Date) UnmarshalText(data []byte) (err error) {
	for _, layout := range layouts {
		var t time.Time
		t, err = time.ParseInLocation(layout, xstring.BytesToString(data), time.Local)
		if err == nil {
			*d = (Date)(t)
			return nil
		}
	}
	return err
}

func (d *Date) Time() time.Time {
	return time.Time(*d)
}

func NewDate(t time.Time) Date {
	return Date(t)
}
