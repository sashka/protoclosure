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

func fromPBLite(pbl *pbLite, pb proto.Message, zeroIndex bool) error {
	maxTagNumber, tagMap := genTagMap(pb)
	pbType := reflect.TypeOf(pb).Elem()
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

		ft := pbType.Field(i)
		p := &proto.Properties{}
		p.Init(ft.Type, ft.Name, ft.Tag.Get("protobuf"), &ft)
		fv := pbValue.Field(i)
		pblv := (*pbl)[ti]
		switch fv.Kind() {
		case reflect.Ptr:
			shouldSetNewFV := false
			var newFV reflect.Value
			if fv.IsNil() {
				newFV = reflect.New(fv.Type().Elem())
			} else {
				newFV = fv
			}

			if fv.Type().Implements(typeOfMessage) {
				subMessage, ok := pblv.([]interface{})
				if !ok {
					return fmt.Errorf("Illegal JSON sub message format")
				}
				pblSM := pbLite(subMessage)
				err := fromPBLite(&pblSM, newFV.Interface().(proto.Message), zeroIndex)
				if err != nil {
					return err
				}
				fv.Set(newFV)
				continue
			}

			fve := newFV.Elem()
			switch fve.Type() {
			case typeOfUint64, typeOfUint32:
				switch pblt := pblv.(type) {
				case float64:
					shouldSetNewFV = true
					fve.SetUint(uint64(pblt))
				case string:
					shouldSetNewFV = true
					ui64, err := strconv.ParseUint(pblt, 10, 64)
					if err != nil {
						return fmt.Errorf("Illegal Uint64: %s", pblt)
					}
					fve.SetUint(ui64)
				default:
					return fmt.Errorf("Unable to set Int value from %T", pblt)
				}

			case typeOfBool:
				switch pblt := pblv.(type) {
				case bool:
					shouldSetNewFV = true
					fve.SetBool(pblt)
				case float64:
					shouldSetNewFV = true
					fve.SetBool(pblt != 0)
				default:
					return fmt.Errorf("Unable to set Bool value from %T", pblt)
				}

			case typeOfString:
				switch pblt := pblv.(type) {
				case string:
					shouldSetNewFV = true
					fve.SetString(pblt)
				default:
					return fmt.Errorf("Unable to set String value from %T", pblt)
				}

			case typeOfFloat64, typeOfFloat32:
				switch pblt := pblv.(type) {
				case float64:
					shouldSetNewFV = true
					fve.SetFloat(float64(pblt))
				default:
					return fmt.Errorf("Unable to set Uint value from %T", pblt)
				}

			default:
				switch pblt := pblv.(type) {
				case float64:
					shouldSetNewFV = true
					fve.SetInt(int64(pblt))
				case string:
					shouldSetNewFV = true
					i64, err := strconv.ParseInt(pblt, 10, 64)
					if err != nil {
						return fmt.Errorf("Illegal Int64: %s", pblt)
					}
					fve.SetInt(i64)
				default:
					return fmt.Errorf("Unable to set Int value from %T", pblt)
				}
			}

			if shouldSetNewFV {
				fv.Set(newFV)
			}

		case reflect.Slice:
			shouldSetNewFV := false
			var newFV reflect.Value
			if fv.IsNil() {
				newFV = reflect.MakeSlice(fv.Type(), 0, 0)
			} else {
				newFV = fv
			}

			switch newFV.Type() {
			case typeOfSliceUint8:
				switch pblt := pblv.(type) {
				case string:
					shouldSetNewFV = true
					newFV = reflect.AppendSlice(newFV, reflect.ValueOf([]uint8(pblt)))
				default:
					return fmt.Errorf("Unable to set Bytes value from %T", pblt)
				}

			case typeOfSliceInt32:
				switch pblt := pblv.(type) {
				case []interface{}:
					shouldSetNewFV = true
					si32, err := toSliceInt32(pblt)
					if err != nil {
						return err
					}
					newFV = reflect.AppendSlice(newFV, reflect.ValueOf(si32))

				default:
					return fmt.Errorf("Unable to set []Int32 value from %T", pblt)
				}

			case typeOfSliceInt64:
				switch pblt := pblv.(type) {
				case []interface{}:
					shouldSetNewFV = true
					si64, err := toSliceInt64(pblt)
					if err != nil {
						return err
					}
					newFV = reflect.AppendSlice(newFV, reflect.ValueOf(si64))

				default:
					return fmt.Errorf("Unable to set []Int64 value from %T", pblt)
				}

			case typeOfSliceUint32:
				switch pblt := pblv.(type) {
				case []interface{}:
					shouldSetNewFV = true
					sui32, err := toSliceUint32(pblt)
					if err != nil {
						return err
					}
					newFV = reflect.AppendSlice(newFV, reflect.ValueOf(sui32))

				default:
					return fmt.Errorf("Unable to set []Uint32 value from %T", pblt)
				}

			case typeOfSliceUint64:
				switch pblt := pblv.(type) {
				case []interface{}:
					shouldSetNewFV = true
					sui64, err := toSliceUint64(pblt)
					if err != nil {
						return err
					}
					newFV = reflect.AppendSlice(newFV, reflect.ValueOf(sui64))

				default:
					return fmt.Errorf("Unable to set []Uint64 value from %T", pblt)
				}

			case typeOfSliceFloat32:
				switch pblt := pblv.(type) {
				case []interface{}:
					shouldSetNewFV = true
					sui32, err := toSliceFloat32(pblt)
					if err != nil {
						return err
					}
					newFV = reflect.AppendSlice(newFV, reflect.ValueOf(sui32))

				default:
					return fmt.Errorf("Unable to set []Float32 value from %T", pblt)
				}

			case typeOfSliceFloat64:
				switch pblt := pblv.(type) {
				case []interface{}:
					shouldSetNewFV = true
					sui64, err := toSliceFloat64(pblt)
					if err != nil {
						return err
					}
					newFV = reflect.AppendSlice(newFV, reflect.ValueOf(sui64))

				default:
					return fmt.Errorf("Unable to set []Float64 value from %T", pblt)
				}

			case typeOfSliceBool:
				switch pblt := pblv.(type) {
				case []interface{}:
					shouldSetNewFV = true
					sui64, err := toSliceBool(pblt)
					if err != nil {
						return err
					}
					newFV = reflect.AppendSlice(newFV, reflect.ValueOf(sui64))

				default:
					return fmt.Errorf("Unable to set []Bool value from %T", pblt)
				}

			case typeOfSliceString:
				switch pblt := pblv.(type) {
				case []interface{}:
					shouldSetNewFV = true
					sui64, err := toSliceString(pblt)
					if err != nil {
						return err
					}
					newFV = reflect.AppendSlice(newFV, reflect.ValueOf(sui64))

				default:
					return fmt.Errorf("Unable to set []Bool value from %T", pblt)
				}

			default:
				return fmt.Errorf("Unsupported type: %v", newFV.Type())
			}

			if shouldSetNewFV {
				fv.Set(newFV)
			}

		default:
			return fmt.Errorf("Unsupported field type %v", fv.Kind())
		}
	}

	return nil
}
