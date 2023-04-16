package util

import (
	"fmt"
	"strings"
)

type Sl struct {
	list []*slInfo
}

func NewSc(tab, col, as string) *Sl {
	sl := &Sl{}
	sl.list = append(sl.list, &slInfo{
		tab:      tab,
		col:      col,
		as:       as,
		isColumn: true,
	})
	return sl
}

func NewSv(val, as string) *Sl {
	sl := &Sl{}
	sl.list = append(sl.list, &slInfo{
		tab:      "",
		col:      val,
		as:       as,
		isColumn: false,
	})
	return sl
}

func (sl *Sl) Sc(tab, col, as string) *Sl {
	sl.list = append(sl.list, &slInfo{
		tab:      tab,
		col:      col,
		as:       as,
		isColumn: true,
	})
	return sl
}

func (sl *Sl) Sv(val, as string) *Sl {
	sl.list = append(sl.list, &slInfo{
		tab:      "",
		col:      val,
		as:       as,
		isColumn: false,
	})
	return sl
}

func (sl *Sl) Sql() (sql string) {
	for i, si := range sl.list {
		tab, col, as := parseSi(si)
		if i == 0 {
			sql += fmt.Sprintf("%s%s%s", tab, col, as)
		} else {
			sql += fmt.Sprintf(",%s%s%s", tab, col, as)
		}
	}
	return
}

func parseSi(si *slInfo) (tab, col, as string) {
	tab = strings.ReplaceAll(si.tab, "`", "")
	if tab != "" {
		tab = fmt.Sprintf("`%s`.", tab)
	}

	col = strings.ReplaceAll(si.col, "`", "")
	if col != "*" && si.isColumn {
		col = fmt.Sprintf("`%s`", col)
	}

	as = strings.ReplaceAll(si.as, "`", "")

	if as != "" {
		as = fmt.Sprintf(" as `%s`", as)
	}
	return
}

type slInfo struct {
	tab      string
	col      string
	as       string
	isColumn bool
}
