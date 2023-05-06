package util

import (
	"gorm.io/gorm"
	"sync"
)

func QueryPage(db *gorm.DB, model interface{}, condition Condition, dest interface{}) (count int64, err error) {
	whereStr, args := condition.Get()
	tx := db.Model(model).Where(whereStr, args...)

	wg := sync.WaitGroup{}
	wg.Add(2)

	var countErr error
	var queryErr error

	go func() {
		countErr = tx.Count(&count).Error
		wg.Done()
	}()

	go func() {
		queryErr = tx.Find(dest).Error
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
