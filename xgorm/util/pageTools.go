package util

import (
	"gorm.io/gorm"
	"sync"
)

func queryPage(db *gorm.DB, model interface{}, condition Condition, dest interface{}) (count int64, err error) {
	whereStr, args := condition.Get()
	tx := db.Model(model).Where(whereStr, args...)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		err = tx.Count(&count).Error
		wg.Done()
	}()

	go func() {
		err = tx.Find(dest).Error
		wg.Done()
	}()

	wg.Wait()

	return
}
