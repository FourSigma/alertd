package sqlhelpers

import (
	"fmt"
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
)

type MyId uuid.UUID

func TestUUID(tst *testing.T) {
	my := MyId(uuid.NewV4())

	t := reflect.TypeOf(my)

	v := reflect.ValueOf(my)

	u := reflect.TypeOf(uuid.UUID{})

	fmt.Println("Convertable: ", t.ConvertibleTo(u))

	rUUID := reflect.TypeOf(uuid.UUID{})
	fmt.Println(v.Convert(rUUID).Interface().(uuid.UUID))
}
