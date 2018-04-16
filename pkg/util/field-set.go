package util

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Zeroer interface {
	IsZero() bool
}

func NewFieldSet(fs ...Field) FieldSet {
	for i, v := range fs {
		v.pos = uint8(i)
	}
	return FieldSet{fls: fs}
}

func NewField(name string, value interface{}, ptr interface{}, canUpdate bool) Field {
	return Field{
		name:      name,
		value:     value,
		ptr:       ptr,
		canUpdate: canUpdate,
	}
}

type Field struct {
	name      string
	value     interface{}
	ptr       interface{}
	pos       uint8
	canUpdate bool
}

func (f *Field) IsZero() bool {
	switch val := f.value.(type) {
	case time.Time:
		return val.IsZero()
	case *time.Time:
		if val != nil {
			return val.IsZero()
		}
		return val == nil
	case Zeroer:
		return val.IsZero()
	case string:
		return val == ""
	case *string:
		return val == nil
	case uint8, uint32, uint16, uint64:
		return val == 0
	case int8, int32, int16, int64:
		return val == 0
	case float32, float64:
		return val == 0.0
	case *float32:
		if val != nil {
			return *val == 0.0
		}
		return val == nil
	case *float64:
		if val != nil {
			return *val == 0.0
		}
		return val == nil
	case []uint8, []uint32, []uint16, []uint64:
		return val == nil
	case *uint8, *uint32, *uint16, *uint64:
		return val == nil
	case *int8, *int32, *int16, *int64:
		return val == nil
	}

	return false
}

type FieldSet struct {
	fls []Field
}

func (f FieldSet) Diff(cmp FieldSet) (diff FieldSet) {
	diff = FieldSet{}
	fs := cmp.Map()
	for _, v := range f.fls {
		val, ok := fs[v.name]
		if !ok {
			continue
		}
		switch v.value.(type) {
		case uuid.UUID:
			if uuid.Equal(v.value.(uuid.UUID), val.(uuid.UUID)) {
				continue
			}
		case *uuid.UUID:
			if uuid.Equal(*v.value.(*uuid.UUID), *val.(*uuid.UUID)) {
				continue
			}
		case time.Time:
			if v.value.(time.Time).Equal(val.(time.Time)) {
				continue
			}
		// case Data:  //use for polymorphic types
		// 	if v.value.(Hasher).Hash() == val.(Hasher).Hash() {
		// 		continue
		// 	}
		default:
			if v.value == val {
				continue
			}
		}
		diff.fls = append(diff.fls, v)
		fmt.Println("NOT EQUAL", v.value, val, v.value == val)
	}
	return
}
func (f FieldSet) Map() (m map[string]interface{}) {
	m = map[string]interface{}{}
	for _, v := range f.fls {
		m[v.name] = v.value
	}
	return
}

func (f *FieldSet) Add(nf Field) FieldSet {
	f.fls = append(f.fls, nf)
	return *f
}

func (f FieldSet) Args() (fl []string, al []interface{}, ptrs []interface{}) {
	fl = make([]string, len(f.fls))
	al = make([]interface{}, len(f.fls))
	ptrs = make([]interface{}, len(f.fls))

	for i, v := range f.fls {
		fl[i] = v.name
		al[i] = v.value
		ptrs[i] = v.ptr
	}
	return
}

func (f FieldSet) Ptrs() (ptrs []interface{}) {
	ptrs = make([]interface{}, len(f.fls))
	for i, v := range f.fls {
		ptrs[i] = v.ptr
	}
	return
}

func (f FieldSet) FieldNameList(fn func(string) string) (fl []string) {
	switch fn {
	case nil:
		for _, v := range f.fls {
			fl = append(fl, v.name)
		}
	default:
		for _, v := range f.fls {
			fl = append(fl, fn(v.name))
		}
	}

	return
}

func (f FieldSet) IsEmpty() bool {
	return len(f.fls) == 0
}

func (f FieldSet) Filter(filter func(*Field) bool) (n FieldSet) {
	if filter == nil {
		return f
	}
	n = FieldSet{}
	for _, v := range f.fls {
		if !filter(&v) {
			continue
		}
		n.fls = append(n.fls, v)
	}
	return
}
func UpdateableFields(f *Field) bool {
	return f.canUpdate && !f.IsZero()
}

func RemoveUpdatedAt(f *Field) bool {
	return !(f.name == "updatedAt")
}
