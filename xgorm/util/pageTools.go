package util

import (
	"errors"
	"github.com/orinchen/xlib/xwg"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"math"
)

func QueryPageWithGen[T any](do gen.DO, page, size int, order ...field.Expr) (items T, recordCount, pageCount int64, err error) {
	offset := (page - 1) * size
	limit := size
	items, recordCount, err = QueryListWithGen[T](do, offset, limit, order...)
	pageCount = int64(math.Ceil(float64(recordCount) / float64(size)))
	return
}

func QueryListWithGen[T any](do gen.DO, offset, limit int, order ...field.Expr) (items T, count int64, err error) {
	wg := xwg.Group{}
	var tempItems any

	wg.Go(func() error {
		var _err error
		tempItems, _err = do.Order(order...).Offset(offset).Limit(limit).Find()
		if _err == nil {
			ok := false
			if items, ok = tempItems.(T); !ok {
				_err = errors.New("query result type is not match")
			}
		}

		return _err
	})

	wg.Go(func() error {
		var _err error
		count, _err = do.Count()
		return _err
	})

	err = wg.Wait()
	return
}

func QueryPage(db *gorm.DB, order, dest interface{}, page, size int) (count int64, err error) {
	var countTx = db.Begin()
	var queryTx = db.Begin()
	defer func() {
		countTx.Commit()
		queryTx.Commit()
	}()

	wg := xwg.Group{}

	wg.Go(func() error {
		return countTx.Count(&count).Error
	})

	wg.Go(func() error {
		return queryTx.Order(order).Limit(size).Offset(size * (page - 1)).Find(dest).Error
	})

	err = wg.Wait()

	return
}

func QueryPageWithCondition(db *gorm.DB, condition Condition, order, dest interface{}, page, size int) (count int64, err error) {
	whereStr, args := condition.Get()
	tx := db.Where(whereStr, args...)

	return QueryPage(tx, order, dest, page, size)
}
