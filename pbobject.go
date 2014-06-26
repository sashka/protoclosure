// Copyright (c) 2014 SameGoal LLC. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protoclosure

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"code.google.com/p/goprotobuf/proto"
)

var (
	typeOfInt64 = reflect.TypeOf(int64(0))
)

type pbObject map[string]interface{}

func toPBObjectKey(ft *reflect.StructField, tn bool) string {
	p := &proto.Properties{}
	p.Init(ft.Type, ft.Name, ft.Tag.Get("protobuf"), ft)

	k := strings.ToLower(p.OrigName)
	if tn {
		k = strconv.FormatInt(int64(p.Tag), 10)
	}
	return k
}

func toPBObjectValue(v interface{}, tn bool) interface{} {
	switch vt := v.(type) {
	case *int64:
		return strconv.FormatInt(*vt, 10)
	case *uint64:
		return strconv.FormatUint(*vt, 10)
	case []uint8:
		return string(vt)
	case proto.Message:
		return toPBObject(vt, tn)
	default:
		return v
	}
}

func toPBObject(pb proto.Message, tn bool) *pbObject {
	pbo := pbObject{}

	pbType := reflect.TypeOf(pb).Elem()
	pbValue := reflect.ValueOf(pb).Elem()
	for i := 0; i < pbType.NumField(); i++ {
		ft := pbType.Field(i)
		fv := pbValue.Field(i)

		// skip unimportant and unset fields
		if strings.HasPrefix(ft.Name, "XXX_") {
			continue
		}
		if fv.IsNil() {
			continue
		}

		// populate pbo map with rewritten key, value pairs
		k := toPBObjectKey(&ft, tn)
		v := toPBObjectValue(fv.Interface(), tn)
		pbo[k] = v
	}

	return &pbo
}

func setPBObjectFieldPtr(fv *reflect.Value, v interface{}, tn bool) error {
	newFV := reflect.New(fv.Type().Elem())

	if newFV.Type().Implements(typeOfMessage) {
		subMessage, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("Cannot convert %T to %v", v, fv.Type())
		}
		pblSM := pbObject(subMessage)
		err := fromPBObject(&pblSM, newFV.Interface().(proto.Message), tn)
		if err != nil {
			return err
		}
		fv.Set(newFV)
		return nil
	}

	var err error
	switch fv.Type().Elem() {
	case typeOfInt64:
		switch vt := v.(type) {
		case float64:
			// legal conversion
		case string:
			// legal conversion
			v, err = strconv.ParseInt(vt, 10, 64)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("Cannot convert %T to %v", v, fv.Type().Elem())
		}

	case typeOfUint64:
		switch vt := v.(type) {
		case float64:
			// legal conversion
		case string:
			// legal conversion
			v, err = strconv.ParseUint(vt, 10, 64)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("Cannot convert %T to %v", v, fv.Type().Elem())
		}

	case typeOfBool:
		switch vt := v.(type) {
		case bool:
			// legal conversion
		case float64:
			// legal conversion
			v = vt != 0
		default:
			return fmt.Errorf("Cannot convert %T to %v", v, fv.Type().Elem())
		}

	case typeOfString:
		switch v.(type) {
		case string:
			// legal conversion
		default:
			return fmt.Errorf("Cannot convert %T to %v", v, fv.Type().Elem())
		}

	default:
		switch v.(type) {
		case float64:
			// legal conversion
		default:
			return fmt.Errorf("Cannot convert %T to %v", v, fv.Type().Elem())
		}
	}

	fve := reflect.ValueOf(v).Convert(fv.Type().Elem())
	newFV.Elem().Set(fve)
	fv.Set(newFV)
	return nil
}

func setPBObjectFieldSlice(fv *reflect.Value, v interface{}) error {
	newFV := reflect.MakeSlice(fv.Type(), 0, 0)

	var err error
	switch fv.Type() {
	case typeOfSliceUint8:
		switch vt := v.(type) {
		case string:
			// legal conversion
			v = []uint8(vt)
		default:
			return fmt.Errorf("Unable to set Bytes value from %T", v)
		}

	case typeOfSliceInt32:
		switch vt := v.(type) {
		case []interface{}:
			// legal conversion
			v, err = toSliceInt32(vt)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("Cannot convert %T to []int32", vt)
		}

	case typeOfSliceInt64:
		switch vt := v.(type) {
		case []interface{}:
			// legal conversion
			v, err = toSliceInt64(vt)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("Cannot convert %T to []int64", vt)
		}

	case typeOfSliceUint32:
		switch vt := v.(type) {
		case []interface{}:
			// legal conversion
			v, err = toSliceUint32(vt)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("Cannot convert %T to []uint32", vt)
		}

	case typeOfSliceUint64:
		switch vt := v.(type) {
		case []interface{}:
			// legal conversion
			v, err = toSliceUint64(vt)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("Cannot convert %T to []uint64", vt)
		}

	case typeOfSliceFloat32:
		switch vt := v.(type) {
		case []interface{}:
			// legal conversion
			v, err = toSliceFloat32(vt)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("Cannot convert %T to []float32", vt)
		}

	case typeOfSliceFloat64:
		switch vt := v.(type) {
		case []interface{}:
			// legal conversion
			v, err = toSliceFloat64(vt)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("Cannot convert %T to []float64", vt)
		}

	case typeOfSliceBool:
		switch vt := v.(type) {
		case []interface{}:
			// legal conversion
			v, err = toSliceBool(vt)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("Cannot convert %T to []bool", vt)
		}

	case typeOfSliceString:
		switch vt := v.(type) {
		case []interface{}:
			// legal conversion
			v, err = toSliceString(vt)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("Cannot convert %T to []string", vt)
		}

	default:
		return fmt.Errorf("Unsupported slice type: %v", fv.Type())
	}

	vo := reflect.ValueOf(v)
	fv.Set(reflect.AppendSlice(newFV, vo))
	return nil
}

func setPBObjectField(fv *reflect.Value, v interface{}, tn bool) error {
	switch fv.Kind() {
	case reflect.Ptr:
		return setPBObjectFieldPtr(fv, v, tn)

	case reflect.Slice:
		return setPBObjectFieldSlice(fv, v)

	default:
		return fmt.Errorf("Unsupported PBObject Kind: %v", fv.Kind())
	}
}

func fromPBObject(pbo *pbObject, pb proto.Message, tn bool) error {
	pbType := reflect.TypeOf(pb).Elem()
	pbValue := reflect.ValueOf(pb).Elem()
	for i := 0; i < pbType.NumField(); i++ {
		ft := pbType.Field(i)
		fv := pbValue.Field(i)
		k := toPBObjectKey(&ft, tn)

		// skip unimportant and unset fields fields
		if strings.HasPrefix(ft.Name, "XXX_") {
			continue
		}
		v, ok := (*pbo)[k]
		if !ok {
			continue
		}

		// populate fv with rewritten value
		//fmt.Printf("i: %d k: %s fv: %v fv.Type(): %v fv.Kind(): %v v: %v\n",
		//	i, k, fv, fv.Type(), fv.Kind(), v)
		err := setPBObjectField(&fv, v, tn)
		if err != nil {
			return err
		}
	}

	return nil
}
