package sqlhelpers

func NewFieldValueList(tn string) *fieldValueList {
	return &fieldValueList{tname: tn}
}

type fieldValueList struct {
	tname        string
	als          []*fieldValue
	kls          []*fieldValue
	isUpdateable bool
}

func (f *fieldValueList) AddAttributeField(fname string, val interface{}) {
	f.als = append(f.als, &fieldValue{FieldName: fname, Value: val})
}

func (f *fieldValueList) AddKeyField(fname string, val interface{}) {
	f.kls = append(f.kls, &fieldValue{FieldName: fname, Value: val})
}

func (f *fieldValueList) AttributeCount() int {
	return len(f.als)
}

func (f *fieldValueList) KeyCount() int {
	return len(f.kls)
}

func (f *fieldValueList) Table() string {
	return f.tname
}

func (f *fieldValueList) IsUpdateable() bool {
	return f.isUpdateable
}

func (f *fieldValueList) FieldNameAndArgs() (fs []string, fargs []interface{}, ks []string, kargs []interface{}) {
	fs = make([]string, f.AttributeCount())
	fargs = make([]interface{}, f.AttributeCount())
	ks = make([]string, f.KeyCount())
	kargs = make([]interface{}, f.KeyCount())

	for i, v := range f.als {
		fs[i], fargs[i] = v.FieldName, v.Value
	}
	if len(fs) > 0 {
		f.isUpdateable = true
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
