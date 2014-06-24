// Copyright (c) 2014 SameGoal LLC. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package protoclosure implements a JSON-based interface for sharing protocol
// buffer messges between goprotobuf and closure-library's goog.proto2.
package protoclosure

import (
	"encoding/json"
	"fmt"

	"code.google.com/p/goprotobuf/proto"
)

// MarshalPBLite takes the protocol buffer and encodes it into the PBLite JSON
// format, returning the data.
func MarshalPBLite(pb proto.Message) ([]byte, error) {
	return nil, fmt.Errorf("Umimplemented")
}

// MarshalPBLiteZeroIndex takes the protocol buffer and encodes it into the
// zero-indexed PBLite JSON format, returning the data.
func MarshalPBLiteZeroIndex(pb proto.Message) ([]byte, error) {
	return nil, fmt.Errorf("Umimplemented")
}

// MarshalObjectKeyName takes the protocol buffer and encodes it into the
// Object JSON format using field names as the JSON keys, returning the data.
func MarshalObjectKeyName(pb proto.Message) ([]byte, error) {
	j, err := json.Marshal(pb)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// MarshalObjectKeyTag takes the protocol buffer and encodes it into the Object
// JSON format using tag numbers as the JSON keys, returning the data.
func MarshalObjectKeyTag(pb proto.Message) ([]byte, error) {
	j, err := json.Marshal(pb)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// UnmarshalPBLite parses the PBLite JSON format protocol buffer representation
// in data and places the decoded result in pb.
func UnmarshalPBLite(data []byte, pb proto.Message) error {
	return fmt.Errorf("Umimplemented")
}

// UnmarshalPBLiteZeroIndex parses the zero-indexed PBLite JSON format protocol
// buffer representation in data and places the decoded result in pb.
func UnmarshalPBLiteZeroIndex(data []byte, pb proto.Message) error {
	return fmt.Errorf("Umimplemented")
}

// UnmarshalObjectKeyName parses the field name based Object JSON format
// protocol buffer representation in data and places the decoded result in pb.
func UnmarshalObjectKeyName(data []byte, pb proto.Message) error {
	err := json.Unmarshal(data, pb)
	if err != nil {
		return err
	}
	return nil
}

// UnmarshalObjectKeyTag parses the tag number based Object JSON format
// protocol buffer representation in data and places the decoded result in pb.
func UnmarshalObjectKeyTag(data []byte, pb proto.Message) error {
	err := json.Unmarshal(data, pb)
	if err != nil {
		return err
	}
	return nil
}
