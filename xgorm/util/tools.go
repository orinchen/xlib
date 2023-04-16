package util

import (
	"fmt"
	"gorm.io/gorm/schema"
)

func Tabler(t schema.Tabler) string {
	return fmt.Sprintf("%s", t)
}

func TablerAs(t schema.Tabler, as string) string {
	return fmt.Sprintf("%s as %s", t.TableName(), as)
}
