package util

import (
	"gorm.io/gorm"
	"sync"
)

func QueryPage(db *gorm.DB, dest interface{}, page, size int) (count int64, err error) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	var countErr error
	var queryErr error

	go func() {
		countErr = db.Count(&count).Error
		wg.Done()
	}()

	go func() {
		queryErr = db.Offset(size * (page - 1)).Limit(size).Find(dest).Error
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
