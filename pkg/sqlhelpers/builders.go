package sqlhelpers

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

func BuildUpdateQuery(tn string, fs []string, ks []string) string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "UPDATE %s SET ", tn)

	pnum := 1

	//field_name=$1
	for _, v := range fs {
		fmt.Fprintf(buf, "%s=$%d", v, pnum)
		pnum = pnum + 1
		if len(fs) >= pnum {
			fmt.Fprint(buf, ", ")
		}
	}

	fmt.Fprintf(buf, " WHERE (%s) IN (", strings.Join(ks, ", "))

	for i, _ := range ks {
		fmt.Fprintf(buf, "$%d", pnum)
		pnum = pnum + 1
		if len(ks) > i+1 {
			fmt.Fprint(buf, ", ")
		}
	}

	fmt.Fprint(buf, ") ")

	return buf.String()

}

func InQueryPlaceholder(total int, keyLen int) string {
	buf := &bytes.Buffer{}
	all := total * keyLen
	als := make([]string, all)
	tls := make([]string, total)

	for i := 0; i < all; i++ {
		als[i] = fmt.Sprintf("$%d", i+1)
	}

	for i, s, als := 0, als[:keyLen], als[keyLen:]; ; i, s, als = i+1, als[:keyLen], als[keyLen:] {
		tls[i] = fmt.Sprintf("(%s)", strings.Join(s, ","))
		if len(als) == 0 {
			break
		}
	}

	fmt.Fprintf(buf, "(%s)", strings.Join(tls, ", "))
	return buf.String()
}

func StringNotEqualAndNotEmpty(src string, dest string) bool {
	return (src != dest) && (src != "")
}

func TimeNotEqualAndNotEmpty(src time.Time, dest time.Time) bool {
	return !src.Equal(dest) && !src.IsZero()
}
