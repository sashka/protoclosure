// Copyright (c) 2014 SameGoal LLC. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protoclosure

import (
	"fmt"
	"reflect"
	"strconv"

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
	typeOfBool    = reflect.TypeOf(true)
	typeOfInt64   = reflect.TypeOf(int64(0))
)

func setPBFieldPtr(fv *reflect.Value, v interface{}) error {
	newFV := reflect.New(fv.Type().Elem())

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
			fmt.Printf("fv: %v v: %v\n", fv, v)
			return fmt.Errorf("Cannot convert %T to %v", v, fv.Type().Elem())
		}
	}

	fve := reflect.ValueOf(v).Convert(fv.Type().Elem())
	newFV.Elem().Set(fve)
	fv.Set(newFV)
	return nil
}

func setPBFieldSlice(fv *reflect.Value, v interface{}) error {
	if vt, ok := v.([]interface{}); ok && len(vt) == 0 {
		return nil
	}

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

func toSliceFloat32(s []interface{}) ([]float32, error) {
	d := []float32{}
	for _, v := range s {
		switch vt := v.(type) {
		case float64:
			d = append(d, float32(vt))
		default:
			return nil, fmt.Errorf("Illegal type in slice: %T", v)
		}
	}
	return d, nil
}

func toSliceFloat64(s []interface{}) ([]float64, error) {
	d := []float64{}
	for _, v := range s {
		switch vt := v.(type) {
		case float64:
			d = append(d, vt)
		default:
			return nil, fmt.Errorf("Illegal type in slice: %T", v)
		}
	}
	return d, nil
}

func toSliceInt32(s []interface{}) ([]int32, error) {
	d := []int32{}
	for _, v := range s {
		switch vt := v.(type) {
		case float64:
			d = append(d, int32(vt))
		default:
			return nil, fmt.Errorf("Illegal type in slice: %T", v)
		}
	}
	return d, nil
}

func toSliceInt64(s []interface{}) ([]int64, error) {
	d := []int64{}
	for _, v := range s {
		switch vt := v.(type) {
		case string:
			i64, err := strconv.ParseInt(vt, 10, 64)
			if err != nil {
				return nil, err
			}
			d = append(d, i64)
		default:
			return nil, fmt.Errorf("Illegal type in slice: %T", v)
		}
	}
	return d, nil
}

func toSliceUint32(s []interface{}) ([]uint32, error) {
	d := []uint32{}
	for _, v := range s {
		switch vt := v.(type) {
		case float64:
			d = append(d, uint32(vt))
		default:
			return nil, fmt.Errorf("Illegal type in slice: %T", v)
		}
	}
	return d, nil
}

func toSliceUint64(s []interface{}) ([]uint64, error) {
	d := []uint64{}
	for _, v := range s {
		switch vt := v.(type) {
		case string:
			ui64, err := strconv.ParseUint(vt, 10, 64)
			if err != nil {
				return nil, err
			}
			d = append(d, ui64)
		default:
			return nil, fmt.Errorf("Illegal type in slice: %T", v)
		}
	}
	return d, nil
}

func toSliceBool(s []interface{}) ([]bool, error) {
	d := []bool{}
	for _, v := range s {
		switch vt := v.(type) {
		case bool:
			d = append(d, vt)
		case float64:
			d = append(d, vt != 0)
		default:
			return nil, fmt.Errorf("Illegal type in slice: %T", v)
		}
	}
	return d, nil
}

func toSliceString(s []interface{}) ([]string, error) {
	d := []string{}
	for _, v := range s {
		switch vt := v.(type) {
		case string:
			d = append(d, vt)
		default:
			return nil, fmt.Errorf("Illegal type in slice: %T", v)
		}
	}
	return d, nil
}
