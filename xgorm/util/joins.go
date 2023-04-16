package util

import (
	"fmt"
	"gorm.io/gorm/schema"
)

type Join struct {
	tableInfo tableInfo
	list      []*onInfo
}

func NewJoin(table schema.Tabler, as string) *Join {
	return &Join{
		tableInfo: tableInfo{
			join:      "join",
			tableName: table.TableName(),
			as:        as,
		},
	}
}

func NewLJoin(table schema.Tabler, as string) *Join {
	return &Join{
		tableInfo: tableInfo{
			join:      "left join",
			tableName: table.TableName(),
			as:        as,
		},
	}
}

func NewRJoin(table schema.Tabler, as string) *Join {
	return &Join{
		tableInfo: tableInfo{
			join:      "right join",
			tableName: table.TableName(),
			as:        as,
		},
	}
}

func (j *Join) On(lColumn, rt, rColumn string) *Join {
	j.list = append(j.list, &onInfo{
		lColumn: lColumn,
		rAs:     rt,
		rColumn: rColumn,
	})
	return j
}

func (j *Join) Column(column string) string {
	return fmt.Sprintf("%s.%s", j.tableInfo.as, column)
}

func (j *Join) Sql() (join string) {
	join += fmt.Sprintf("%s %s as %s", j.tableInfo.join, j.tableInfo.tableName, j.tableInfo.as)
	for i, on := range j.list {
		if i == 0 {
			join += fmt.Sprintf(" on `%s`.`%s` = `%s`.`%s`", j.tableInfo.as, on.lColumn, on.rAs, on.rColumn)
		} else {
			join += fmt.Sprintf(" and `%s`.`%s` = `%s`.`%s`", j.tableInfo.as, on.lColumn, on.rAs, on.rColumn)
		}
	}
	return
}

type tableInfo struct {
	join      string
	tableName string
	as        string
}

type onInfo struct {
	lColumn string
	rAs     string
	rColumn string
}
