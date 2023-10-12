package util

import (
	"gorm.io/gorm"
	"sync"
)

func QueryPage(db *gorm.DB, dest interface{}, page, size int) (count int64, err error) {
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
		queryErr = queryTx.Limit(size).Offset(size * (page - 1)).Find(dest).Error
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

func QueryPageWithCondition(db *gorm.DB, condition Condition, dest interface{}, page, size int) (count int64, err error) {
	whereStr, args := condition.Get()
	tx := db.Where(whereStr, args...)

	return QueryPage(tx, dest, page, size)
}
