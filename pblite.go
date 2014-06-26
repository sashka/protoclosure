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
	typeOfSliceBool    = reflect.TypeOf([]bool{})
	typeOfSliceUint8   = reflect.TypeOf([]uint8{})
	typeOfSliceFloat32 = reflect.TypeOf([]float32{})
	typeOfSliceFloat64 = reflect.TypeOf([]float64{})
	typeOfSliceUint32  = reflect.TypeOf([]uint32{})
	typeOfSliceUint64  = reflect.TypeOf([]uint64{})
	typeOfSliceInt32   = reflect.TypeOf([]int32{})
	typeOfSliceInt64   = reflect.TypeOf([]int64{})
	typeOfSliceString  = reflect.TypeOf([]string{})

	typeOfMessage = reflect.TypeOf((*proto.Message)(nil)).Elem()
	typeOfString  = reflect.TypeOf("")
	typeOfUint64  = reflect.TypeOf(uint64(0))
	typeOfUint32  = reflect.TypeOf(uint32(0))
	typeOfFloat64 = reflect.TypeOf(float64(0))
	typeOfFloat32 = reflect.TypeOf(float32(0))
	typeOfBool    = reflect.TypeOf(true)
)

type pbLite []interface{}

func genTagMap(pb proto.Message) (int, map[int]int) {
	maxTagNumber := -1
	tagMap := make(map[int]int)
	pbType := reflect.TypeOf(pb).Elem()
	for i := 0; i < pbType.NumField(); i++ {
		ft := pbType.Field(i)
		if strings.HasPrefix(ft.Name, "XXX_") {
			continue
		}
		p := &proto.Properties{}
		p.Init(ft.Type, ft.Name, ft.Tag.Get("protobuf"), &ft)
		if p.Tag > maxTagNumber {
			maxTagNumber = p.Tag
		}
		tagMap[p.Tag] = i
	}
	return maxTagNumber, tagMap
}

func toPBLiteValue(v interface{}, zeroIndex bool) interface{} {
	switch vt := v.(type) {
	case *int64:
		return strconv.FormatInt(*vt, 10)
	case *uint64:
		return strconv.FormatUint(*vt, 10)
	case []uint8:
		return string(vt)
	case proto.Message:
		return toPBLite(vt, zeroIndex)
	case *bool:
		if *vt {
			return int(1)
		}
		return int(0)
	default:
		return v
	}
}

func toPBLite(pb proto.Message, zeroIndex bool) *pbLite {
	pbl := pbLite{}

	maxTagNumber, tagMap := genTagMap(pb)
	pbValue := reflect.ValueOf(pb).Elem()

	startIndex := 0
	if zeroIndex {
		startIndex = 1
	}
	lastNonNil := -1
	for ti := startIndex; ti <= maxTagNumber; ti++ {
		i, ok := tagMap[ti]
		if !ok {
			pbl = append(pbl, nil)
			continue
		}
		fv := pbValue.Field(i)

		// write stub markers for empty fields
		if fv.IsNil() {
			if fv.Kind() == reflect.Slice {
				pbl = append(pbl, []string{})
			} else {
				pbl = append(pbl, nil)
			}
			continue
		}

		v := toPBLiteValue(fv.Interface(), zeroIndex)
		pbl = append(pbl, v)
		lastNonNil = len(pbl)
	}

	// Truncate trailing nils
	pbl = pbl[:lastNonNil]

	return &pbl
}

func setPBLiteField(fv *reflect.Value, v interface{}, zeroIndex bool) error {
	switch fv.Kind() {
	case reflect.Ptr:
		if fv.Type().Implements(typeOfMessage) {
			subMessage, ok := v.([]interface{})
			if !ok {
				return fmt.Errorf("Illegal JSON sub message format")
			}
			pblSM := pbLite(subMessage)
			newFV := reflect.New(fv.Type().Elem())
			err := fromPBLite(&pblSM, newFV.Interface().(proto.Message), zeroIndex)
			if err != nil {
				return err
			}
			fv.Set(newFV)
			return nil
		}
		return setPBFieldPtr(fv, v)

	case reflect.Slice:
		return setPBFieldSlice(fv, v)

	default:
		return fmt.Errorf("Unsupported PBObject Kind: %v", fv.Kind())
	}
}

func fromPBLite(pbl *pbLite, pb proto.Message, zeroIndex bool) error {
	maxTagNumber, tagMap := genTagMap(pb)
	pbValue := reflect.ValueOf(pb).Elem()

	startIndex := 1
	if zeroIndex {
		startIndex = 0
	}
	for ti := startIndex; ti <= maxTagNumber && ti < len(*pbl); ti++ {
		var i int
		var ok bool
		if zeroIndex {
			i, ok = tagMap[ti+1]
		} else {
			i, ok = tagMap[ti]
		}
		if !ok {
			continue
		}

		fv := pbValue.Field(i)
		v := (*pbl)[ti]

		err := setPBLiteField(&fv, v, zeroIndex)
		if err != nil {
			return err
		}
	}

	return nil
}
