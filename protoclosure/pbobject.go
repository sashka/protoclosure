// Copyright (c) 2014 SameGoal LLC. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protoclosure

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
)

type pbObject map[string]interface{}

func toPBObjectKey(ft *reflect.StructField, tagName bool) (string, bool) {
	p := &proto.Properties{}
	p.Init(ft.Type, ft.Name, ft.Tag.Get("protobuf"), ft)

	k := strings.ToLower(p.OrigName)
	if tagName {
		k = strconv.FormatInt(int64(p.Tag), 10)
	}
	numEnc := false
	if strings.HasSuffix(strings.ToLower(p.OrigName), "_number") {
		numEnc = true
	}
	return k, numEnc
}

func toPBObjectValue(v interface{}, tagName, numEnc bool) interface{} {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Slice &&
		val.Type().Elem().Implements(typeOfMessage) {
		messages := make([]interface{}, val.Len())
		for i := 0; i < val.Len(); i++ {
			pb := val.Index(i).Interface().(proto.Message)
			messages[i] = toPBObject(pb, tagName)
		}
		return messages
	}

	switch vt := v.(type) {
	case *int64:
		if numEnc {
			return v
		}
		return strconv.FormatInt(*vt, 10)
	case *uint64:
		if numEnc {
			return v
		}
		return strconv.FormatUint(*vt, 10)
	case []uint8:
		return string(vt)
	case proto.Message:
		return toPBObject(vt, tagName)
	default:
		return v
	}
}

func toPBObject(pb proto.Message, tagName bool) *pbObject {
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
		k, numEnc := toPBObjectKey(&ft, tagName)
		v := toPBObjectValue(fv.Interface(), tagName, numEnc)
		pbo[k] = v
	}

	return &pbo
}

func setPBObjectField(fv *reflect.Value, v interface{}, tagName bool) error {
	if v == nil {
		return nil
	}

	switch fv.Kind() {
	case reflect.Ptr:
		if fv.Type().Implements(typeOfMessage) {
			subMessage, ok := v.(map[string]interface{})
			if !ok {
				return fmt.Errorf("Cannot convert %T to %v", v, fv.Type())
			}
			pboSM := pbObject(subMessage)
			newFV := reflect.New(fv.Type().Elem())
			err := fromPBObject(&pboSM, newFV.Interface().(proto.Message), tagName)
			if err != nil {
				return err
			}
			fv.Set(newFV)
			return nil
		}
		return setPBFieldPtr(fv, v)

	case reflect.Slice:
		if fv.Type().Elem().Kind() == reflect.Ptr &&
			fv.Type().Elem().Implements(typeOfMessage) {
			subMessageSlice, ok := v.([]interface{})
			if !ok {
				return fmt.Errorf("Cannot convert %T to %v", v, fv.Type())
			}
			newFV := reflect.MakeSlice(fv.Type(), 0, 0)
			for _, sm := range subMessageSlice {
				subMessage, ok := sm.(map[string]interface{})
				if !ok {
					return fmt.Errorf("Cannot convert %T to %v", v, fv.Type())
				}
				pboSM := pbObject(subMessage)
				newPB := reflect.New(fv.Type().Elem().Elem())
				err := fromPBObject(&pboSM, newPB.Interface().(proto.Message), tagName)
				if err != nil {
					return err
				}
				newFV = reflect.Append(newFV, newPB)
			}
			fv.Set(newFV)
			return nil
		}
		return setPBFieldSlice(fv, v)

	default:
		return fmt.Errorf("Unsupported PBObject Kind: %v", fv.Kind())
	}
}

func fromPBObject(pbo *pbObject, pb proto.Message, tagName bool) error {
	pbType := reflect.TypeOf(pb).Elem()
	pbValue := reflect.ValueOf(pb).Elem()
	for i := 0; i < pbType.NumField(); i++ {
		ft := pbType.Field(i)
		fv := pbValue.Field(i)
		k, _ := toPBObjectKey(&ft, tagName)

		// skip unimportant and unset fields
		if strings.HasPrefix(ft.Name, "XXX_") {
			continue
		}
		v, ok := (*pbo)[k]
		if !ok {
			continue
		}

		// populate fv with rewritten value
		err := setPBObjectField(&fv, v, tagName)
		if err != nil {
			return err
		}
	}

	return nil
}
