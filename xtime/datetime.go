package xtime

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"time"
)

type DateTime time.Time

func (dt DateTime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(time.Time(dt))
}

func (dt *DateTime) UnmarshalBSONValue(t bsontype.Type, value []byte) (err error) {
	if t != bson.TypeDateTime {
		return fmt.Errorf("invalid bson value type '%s'", t.String())
	}
	s, _, ok := bsoncore.ReadDateTime(value)
	if !ok {
		return fmt.Errorf("invalid bson string value")
	}

	*dt = DateTime(time.UnixMilli(s))
	return
}

func (dt DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(dt))
}

// UnmarshalJSON 实现JSON反序列化方法
func (dt *DateTime) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err = json.Unmarshal(data, &s); err != nil {
		return err
	}
	var t *time.Time
	if t, err = AutoParseInLocation(s, time.Local); err != nil {
		return
	}
	*dt = (DateTime)(*t)
	return
}

func (dt DateTime) MarshalText() (text []byte, err error) {
	var stamp = fmt.Sprintf("%s", time.Time(dt).Format(time.DateTime))
	return []byte(stamp), nil
}

func (dt *DateTime) UnmarshalText(data []byte) (err error) {
	var t *time.Time
	if t, err = AutoParseInLocation(string(data), time.Local); err != nil {
		return
	}
	*dt = (DateTime)(*t)
	return
}

func (dt *DateTime) Time() time.Time {
	return time.Time(*dt)
}
