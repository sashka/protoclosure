// Copyright (c) 2014 SameGoal LLC. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protoclosure

import (
	"fmt"
	"strconv"
)

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

