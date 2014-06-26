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

type pbObject map[string]interface{}

func toPBObjectKey(ft *reflect.StructField, tagNumbers bool) string {
	p := &proto.Properties{}
	p.Init(ft.Type, ft.Name, ft.Tag.Get("protobuf"), ft)

	k := strings.ToLower(p.OrigName)
	if tagNumbers {
		k = strconv.FormatInt(int64(p.Tag), 10)
	}
	return k
}

func toPBObjectValue(v interface{}, tagNumbers bool) interface{} {
	switch vt := v.(type) {
	case *int64:
		return strconv.FormatInt(*vt, 10)
	case *uint64:
		return strconv.FormatUint(*vt, 10)
	case []uint8:
		return string(vt)
	case proto.Message:
		return toPBObject(vt, tagNumbers)
	default:
		return v
	}
}

func toPBObject(pb proto.Message, tagNumbers bool) *pbObject {
	pbo := pbObject{}

	pbType := reflect.TypeOf(pb).Elem()
	pbValue := reflect.ValueOf(pb).Elem()
	for i := 0; i < pbType.NumField(); i++ {
		ft := pbType.Field(i)
		fv := pbValue.Field(i)

		// skip unimportant fields
		if strings.HasPrefix(ft.Name, "XXX_") {
			continue
		}
		if fv.IsNil() {
			continue
		}

		// populate pbo map with rewritten key, value pairs
		k := toPBObjectKey(&ft, tagNumbers)
		v := toPBObjectValue(fv.Interface(), tagNumbers)
		pbo[k] = v
	}

	return &pbo
}

func fromPBObject(pbo *pbObject, pb proto.Message, tagNumbers bool) error {
	return fmt.Errorf("TODO(hochhaus): Implement fromPBObject")
}
