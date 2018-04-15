package sqlhelpers

type FieldValueList struct {
	tname string
	als   []*fieldValue
	kls   []*fieldValue
}

func (f *FieldValueList) AddAttributeField(fname string, val interface{}) {
	f.als = append(f.als, &fieldValue{FieldName: fname, Value: val})
}

func (f *FieldValueList) AddKeyField(fname string, val interface{}) {
	f.kls = append(f.kls, &fieldValue{FieldName: fname, Value: val})
}

func (f *FieldValueList) AttributeCount() int {
	return len(f.als)
}

func (f *FieldValueList) KeyCount() int {
	return len(f.kls)
}

func (f *FieldValueList) Table() string {
	return f.tname
}

func (f *FieldValueList) FieldNameAndArgs() (fs []string, fargs []interface{}, ks []string, kargs []interface{}) {
	fs = make([]string, f.AttributeCount())
	fargs = make([]interface{}, f.AttributeCount())
	ks = make([]string, f.KeyCount())
	kargs = make([]interface{}, f.KeyCount())

	for i, v := range f.als {
		fs[i], fargs[i] = v.FieldName, v.Value
	}
	for i, v := range f.kls {
		ks[i], kargs[i] = v.FieldName, v.Value
	}

	return
}

type fieldValue struct {
	FieldName string
	Value     interface{}
}
