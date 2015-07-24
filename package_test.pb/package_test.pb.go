// Code generated by protoc-gen-go.
// source: gopkg.in/samegoal/protoclosure.v0/package_test.proto
// DO NOT EDIT!

/*
Package someprotopackage is a generated protocol buffer package.

It is generated from these files:
	gopkg.in/samegoal/protoclosure.v0/package_test.proto

It has these top-level messages:
	TestPackageTypes
*/
package someprotopackage

import proto "github.com/golang/protobuf/proto"
import math "math"
import test "gopkg.in/samegoal/protoclosure.v0/test.pb"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type TestPackageTypes struct {
	OptionalInt32    *int32               `protobuf:"varint,1,opt,name=optional_int32" json:"optional_int32,omitempty"`
	OtherAll         *test.TestAllTypes   `protobuf:"bytes,2,opt,name=other_all" json:"other_all,omitempty"`
	RepOtherAll      []*test.TestAllTypes `protobuf:"bytes,3,rep,name=rep_other_all" json:"rep_other_all,omitempty"`
	XXX_unrecognized []byte               `json:"-"`
}

func (m *TestPackageTypes) Reset()         { *m = TestPackageTypes{} }
func (m *TestPackageTypes) String() string { return proto.CompactTextString(m) }
func (*TestPackageTypes) ProtoMessage()    {}

func (m *TestPackageTypes) GetOptionalInt32() int32 {
	if m != nil && m.OptionalInt32 != nil {
		return *m.OptionalInt32
	}
	return 0
}

func (m *TestPackageTypes) GetOtherAll() *test.TestAllTypes {
	if m != nil {
		return m.OtherAll
	}
	return nil
}

func (m *TestPackageTypes) GetRepOtherAll() []*test.TestAllTypes {
	if m != nil {
		return m.RepOtherAll
	}
	return nil
}

func init() {
}
