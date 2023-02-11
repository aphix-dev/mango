package mango

import (
	"reflect"
)

func MakePub[Type any](in Type) Type {
	var a Type
	pub := reflect.New(reflect.TypeOf(a))
	v := reflect.ValueOf(in)
	for i := 0; i < v.NumField(); i++ {
		curVal := v.Field(i)
		accessTag := v.Type().Field(i).Tag.Get("access")
		if accessTag == "private" {
		} else if accessTag == "public" || accessTag == "" {
			reflect.Indirect(pub).Field(i).Set(curVal)
		}
	}

	g := pub.Interface()

	return *(g.(*Type))
}
