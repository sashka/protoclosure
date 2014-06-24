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
	typeOfSliceBytes = reflect.TypeOf([]byte(nil))
	typeOfSliceUint8 = reflect.TypeOf([]uint8{})
	typeOfMessage = reflect.TypeOf((*proto.Message)(nil)).Elem()
	typeOfString = reflect.TypeOf("")
)

type pbLite []interface{}

func toPBLite(pb proto.Message, zeroIndex bool) (*pbLite, error) {
	pbl := pbLite{}

	maxTagNumber := -1
	tagMap := make(map[int]int)
	pbType := reflect.TypeOf(pb).Elem()
	pbValue := reflect.ValueOf(pb).Elem()
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

		ft := pbType.Field(i)
		p := &proto.Properties{}
		p.Init(ft.Type, ft.Name, ft.Tag.Get("protobuf"), &ft)
		fv := pbValue.Field(i)
		switch fv.Kind() {
		case reflect.Ptr:
			if fv.IsNil() {
				pbl = append(pbl, nil)
				continue
			}

			if fv.Type().Implements(typeOfMessage) {
				subMessage, err := toPBLite(fv.Interface().(proto.Message), zeroIndex)
				if err != nil {
					return nil, err
				}
				pbl = append(pbl, subMessage)
			} else {
				fve := fv.Elem()
				switch fve.Kind() {
				case reflect.Int64:
					val := fve.Int()
					pbl = append(pbl, strconv.FormatInt(val, 10))
				case reflect.Uint64:
					val := fve.Uint()
					pbl = append(pbl, strconv.FormatUint(val, 10))
				case reflect.Bool:
					if fve.Bool() {
						pbl = append(pbl, 1)
					} else {
						pbl = append(pbl, 0)
					}
				default:
					pbl = append(pbl, fv.Interface())
				}
			}
		case reflect.Slice:
			switch fv.Type() {
			case typeOfSliceBytes, typeOfSliceUint8:
				pbl = append(pbl, fv.Convert(typeOfString).String())
			default:
				if fv.IsNil() {
					pbl = append(pbl, []string{})
					continue
				} else {
					pbl = append(pbl, fv.Interface())
				}
			}
		default:
			pbl = append(pbl, fv.Interface())
		}
		lastNonNil = len(pbl)
	}

	// Truncate trailing nils
	pbl = pbl[:lastNonNil]

	return &pbl, nil
}

func toPB(pbl *pbLite, pb proto.Message, zeroIndex bool) error {
	return fmt.Errorf("TODO(hochhaus): Implement pbLiteToPB")
}
