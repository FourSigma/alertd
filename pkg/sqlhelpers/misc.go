package sqlhelpers

import (
	"fmt"
	"strings"
)

func ModifyStringList(ls []string, mFn func(string) string) (rs []string) {
	rs = make([]string, len(ls))
	for i, _ := range ls {
		rs[i] = mFn(ls[i])
	}
	return rs
}

func ModifyString(src string, mFn func(string) string) string {
	return mFn(src)
}

func PostgresPlaceholder(l int) string {
	pls := make([]string, l)
	for i, _ := range pls {
		pls[i] = fmt.Sprintf("$%d", i+1)
	}
	return strings.Join(pls, ", ")
}
