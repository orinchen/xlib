package util

import "fmt"
import "strings"

// 自定义sql查询
type Condition struct {
	list []*conditionInfo
}

func NewCondition() *Condition {
	return &Condition{}
}

func (c *Condition) AndWithCondition(condition bool, column string, cases string, value interface{}) *Condition {
	if condition {
		c.list = append(c.list, &conditionInfo{
			andor:  "and",
			column: column,        // 列名
			case_:  string(cases), // 条件(and,or,in,>=,<=)
			value:  value,
		})
	}
	return c
}

// And a Condition by and .and 一个条件
func (c *Condition) And(column string, cases string, value interface{}) *Condition {
	return c.AndWithCondition(true, column, cases, value)
}

func (c *Condition) AndEq(column string, value interface{}) *Condition {
	return c.AndWithCondition(true, column, "=", value)
}

func (c *Condition) AndLt(column string, value interface{}) *Condition {
	return c.AndWithCondition(true, column, "<", value)
}

func (c *Condition) AndGt(column string, value interface{}) *Condition {
	return c.AndWithCondition(true, column, ">", value)
}

func (c *Condition) AndLte(column string, value interface{}) *Condition {
	return c.AndWithCondition(true, column, "<=", value)
}

func (c *Condition) AndGte(column string, value interface{}) *Condition {
	return c.AndWithCondition(true, column, ">=", value)
}

func (c *Condition) AndIsNull(column string) *Condition {
	return c.AndWithCondition(true, column, "is", nil)
}

func (c *Condition) AndIsNotNull(column string) *Condition {
	return c.AndWithCondition(true, column, "is Not", nil)
}

func (c *Condition) OrWithCondition(condition bool, column string, cases string, value interface{}) *Condition {
	if condition {
		c.list = append(c.list, &conditionInfo{
			andor:  "or",
			column: column,        // 列名
			case_:  string(cases), // 条件(and,or,in,>=,<=)
			value:  value,
		})
	}
	return c
}

// Or a Condition by or .or 一个条件
func (c *Condition) Or(column string, cases string, value interface{}) *Condition {
	return c.OrWithCondition(true, column, cases, value)
}

func (c *Condition) OrEq(column string, value interface{}) *Condition {
	return c.OrWithCondition(true, column, "=", value)
}

func (c *Condition) OrLt(column string, value interface{}) *Condition {
	return c.OrWithCondition(true, column, "<", value)
}

func (c *Condition) OrGt(column string, value interface{}) *Condition {
	return c.OrWithCondition(true, column, ">", value)
}

func (c *Condition) OrLte(column string, value interface{}) *Condition {
	return c.OrWithCondition(true, column, "<=", value)
}

func (c *Condition) OrGte(column string, value interface{}) *Condition {
	return c.OrWithCondition(true, column, ">=", value)
}

func (c *Condition) Get() (where string, out []interface{}) {
	firstAnd := -1
	for i := 0; i < len(c.list); i++ { // 查找第一个and
		if c.list[i].andor == "and" {
			_, column, _case := parseCondition(c.list[i])
			where = fmt.Sprintf("%v %v ?", column, _case)
			out = append(out, c.list[i].value)
			firstAnd = i
			break
		}
	}

	if firstAnd < 0 && len(c.list) > 0 { // 补刀
		_, column, _case := parseCondition(c.list[0])
		where = fmt.Sprintf("%v %v ?", column, _case)
		out = append(out, c.list[0].value)
		firstAnd = 0
	}

	for i := 0; i < len(c.list); i++ { // 添加剩余的
		if firstAnd != i {
			andor, column, _case := parseCondition(c.list[i])
			where += fmt.Sprintf(" %v %v %v ?", andor, column, _case)
			out = append(out, c.list[i].value)
		}
	}

	return
}

func parseCondition(condition *conditionInfo) (andor, column, _case string) {
	andor = condition.andor
	_case = condition.case_
	cis := strings.Split(strings.ReplaceAll(condition.column, "`", ""), ".")
	for i, ci := range cis {
		column += "`" + ci + "`"
		if i < len(cis)-1 {
			column += "."
		}
	}
	return
}

type conditionInfo struct {
	andor  string
	column string // 列名
	case_  string // 条件(in,>=,<=)
	value  interface{}
}
