// Copyright (c) 2014 SameGoal LLC. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package protoclosure implements a JSON-based interface for sharing protocol
// buffer messges between goprotobuf and closure-library's goog.proto2.
package protoclosure

import (
	"encoding/json"

	"github.com/golang/protobuf/proto"
)

// MarshalPBLite takes the protocol buffer and encodes it into the PBLite JSON
// format, returning the data.
func MarshalPBLite(pb proto.Message) ([]byte, error) {
	return json.Marshal(toPBLite(pb, false))
}

// MarshalPBLiteZeroIndex takes the protocol buffer and encodes it into the
// zero-indexed PBLite JSON format, returning the data.
func MarshalPBLiteZeroIndex(pb proto.Message) ([]byte, error) {
	return json.Marshal(toPBLite(pb, true))
}

// MarshalObjectKeyName takes the protocol buffer and encodes it into the
// Object JSON format using field names as the JSON keys, returning the data.
func MarshalObjectKeyName(pb proto.Message) ([]byte, error) {
	return json.Marshal(toPBObject(pb, false))
}

// MarshalObjectKeyTag takes the protocol buffer and encodes it into the Object
// JSON format using tag numbers as the JSON keys, returning the data.
func MarshalObjectKeyTag(pb proto.Message) ([]byte, error) {
	return json.Marshal(toPBObject(pb, true))
}

// UnmarshalPBLite parses the PBLite JSON format protocol buffer representation
// in data and places the decoded result in pb.
func UnmarshalPBLite(data []byte, pb proto.Message) error {
	pbl := &pbLite{}
	err := json.Unmarshal(data, pbl)
	if err != nil {
		return err
	}
	return fromPBLite(pbl, pb, false)
}

// UnmarshalPBLiteZeroIndex parses the zero-indexed PBLite JSON format protocol
// buffer representation in data and places the decoded result in pb.
func UnmarshalPBLiteZeroIndex(data []byte, pb proto.Message) error {
	pbl := &pbLite{}
	err := json.Unmarshal(data, pbl)
	if err != nil {
		return err
	}
	return fromPBLite(pbl, pb, true)
}

// UnmarshalObjectKeyName parses the field name based Object JSON format
// protocol buffer representation in data and places the decoded result in pb.
func UnmarshalObjectKeyName(data []byte, pb proto.Message) error {
	pbo := &pbObject{}
	err := json.Unmarshal(data, pbo)
	if err != nil {
		return err
	}
	return fromPBObject(pbo, pb, false)
}

// UnmarshalObjectKeyTag parses the tag number based Object JSON format
// protocol buffer representation in data and places the decoded result in pb.
func UnmarshalObjectKeyTag(data []byte, pb proto.Message) error {
	pbo := &pbObject{}
	err := json.Unmarshal(data, pbo)
	if err != nil {
		return err
	}
	return fromPBObject(pbo, pb, true)
}
