package hashid

import (
	"fmt"
	"github.com/orinchen/xlib/xstring"
	"github.com/pkg/errors"
	"github.com/sqids/sqids-go"
	"strings"
	"sync"
)

var _hashid *sqids.Sqids

var once sync.Once

func Init(conf Config) {
	once.Do(func() {
		_hashid = conf.Build()
	})
}

func Encode(numbers ...uint64) (string, error) {
	if _hashid == nil {
		return "", errors.New("HashID 没有初始化")
	}

	return _hashid.Encode(numbers)
}

func DecodeSlice(hash string) ([]uint64, error) {
	if _hashid == nil {
		return nil, errors.New("HashID 没有初始化")
	}

	return _hashid.Decode(hash), nil
}

func Decode(hash string) (uint64, error) {
	if _hashid == nil {
		return 0, errors.New("HashID 没有初始化")
	}
	result, err := DecodeSlice(hash)
	if err != nil {
		return 0, err
	} else if len(result) < 1 {
		return 0, nil
	} else {
		return result[0], nil
	}
}

type HashId uint64

func (id HashId) MarshalJSON() ([]byte, error) {
	str, err := Encode(uint64(id))
	if err != nil {
		return nil, err
	}
	return xstring.StringToBytes(fmt.Sprintf("\"%s\"", str)), nil
}

func (id *HashId) UnmarshalJSON(data []byte) (err error) {
	ints, err := Decode(strings.ReplaceAll(xstring.BytesToString(data), "\"", ""))
	if err != nil {
		return err
	}
	*id = HashId(ints)
	return
}

func (id *HashId) MarshalText() (text []byte, err error) {
	str, err := Encode(uint64(*id))
	if err != nil {
		return nil, err
	}
	return xstring.StringToBytes(fmt.Sprintf("%s", str)), nil
}

func (id *HashId) UnmarshalText(data []byte) (err error) {
	ints, err := Decode(xstring.BytesToString(data))
	if err != nil {
		return err
	}
	*id = HashId(ints)
	return
}

type HashIds []uint64

func (ids *HashIds) MarshalJSON() ([]byte, error) {
	str, err := Encode(*ids...)
	if err != nil {
		return nil, err
	}
	return xstring.StringToBytes(fmt.Sprintf("\"%s\"", str)), nil
}

func (ids *HashIds) UnmarshalJSON(data []byte) (err error) {
	ints, err := DecodeSlice(strings.ReplaceAll(xstring.BytesToString(data), "\"", ""))
	if err != nil {
		return err
	}
	*ids = ints
	return
}

func (ids *HashIds) MarshalText() (text []byte, err error) {
	str, err := Encode(*ids...)
	if err != nil {
		return nil, err
	}
	return xstring.StringToBytes(fmt.Sprintf("%s", str)), nil
}

func (ids *HashIds) UnmarshalText(data []byte) (err error) {
	hashIds := strings.Split(string(data), ",")
	for _, s := range hashIds {
		var id uint64
		id, err = Decode(s)
		if err != nil {
			return err
		}
		*ids = append(*ids, id)
	}
	return
}
