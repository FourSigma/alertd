package sqlhelpers

import (
	"bytes"
	"fmt"
	"strings"
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
