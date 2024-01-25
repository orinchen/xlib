package util

import (
	"errors"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"math"
	"sync"
)

func QueryPageWithGen[T any](do gen.DO, page, size int, order ...field.Expr) (items T, recordCount, pageCount int64, err error) {
	offset := (page - 1) * size
	limit := size
	items, recordCount, err = QueryListWithGen[T](do, offset, limit, order...)
	pageCount = int64(math.Ceil(float64(recordCount) / float64(size)))
	return
}

func QueryListWithGen[T any](do gen.DO, offset, limit int, order ...field.Expr) (items T, count int64, err error) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	var countErr, queryErr error
	var tempItems any
	go func() {
		tempItems, queryErr = do.Order(order...).Offset(offset).Limit(limit).Find()
		if queryErr == nil {
			ok := false
			if items, ok = tempItems.(T); !ok {
				queryErr = errors.New("query result type is not match")
			}
		}
		wg.Done()
	}()
	go func() {
		count, countErr = do.Count()
		wg.Done()
	}()
	wg.Wait()
	if countErr != nil || queryErr != nil {
		err = errors.Join(countErr, queryErr)
	}
	return
}

func QueryPage(db *gorm.DB, order, dest interface{}, page, size int) (count int64, err error) {
	var countErr error
	var queryErr error

	var countTx = db.Begin()
	var queryTx = db.Begin()
	defer func() {
		countTx.Commit()
		queryTx.Commit()
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		countErr = countTx.Count(&count).Error
		wg.Done()
	}()

	go func() {
		queryErr = queryTx.Order(order).Limit(size).Offset(size * (page - 1)).Find(dest).Error
		wg.Done()
	}()

	wg.Wait()

	if countErr != nil {
		return 0, countErr
	}

	if queryErr != nil {
		return 0, queryErr
	}

	return
}

func QueryPageWithCondition(db *gorm.DB, condition Condition, order, dest interface{}, page, size int) (count int64, err error) {
	whereStr, args := condition.Get()
	tx := db.Where(whereStr, args...)

	return QueryPage(tx, order, dest, page, size)
}
