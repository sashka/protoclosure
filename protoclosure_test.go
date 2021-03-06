// Copyright (c) 2014 SameGoal LLC. All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protoclosure

import (
	"bytes"
	"testing"

	"github.com/golang/protobuf/proto"

	package_test_pb "protoclosure/package_test_pb"
	test_pb "protoclosure/test_pb"
)

// JS limits integers stored in Numbers (not strings) to the IEEE-754 floating
// point range: (-2^53+1), ..., (2^53-1)
//
// Number.MAX_SAFE_INTEGER = 9007199254740991
// Number.MIN_SAFE_INTEGER -9007199254740991

const (
	oobJSInt     = 9007199254740993
	oobJSStr     = "9007199254740993"
	maxSafeJSInt = 9007199254740991
	maxSafeJSStr = "9007199254740991"

	// golden values were extracted from closure-library unit tests.
	pbLiteGolden = "[null," +
		"101," +
		"\"102\"," +
		"103," +
		"\"104\"," +
		"105," +
		"\"106\"," +
		"107," +
		"\"108\"," +
		"109," +
		"\"110\"," +
		"111.5," +
		"112.5," +
		"1," +
		"\"test\"," +
		"\"abcd\"," +
		"[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,111]," +
		"null," +
		"[null,112]," +
		"null," +
		"null," +
		"0," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"[201,202]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[\"foo\",\"bar\"]," +
		"[]," +
		"[]," +
		"null," +
		"[]," +
		"[0,2]" +
		"]"

	pbLiteZeroIndexGolden = "[" +
		"101," +
		"\"102\"," +
		"103," +
		"\"104\"," +
		"105," +
		"\"106\"," +
		"107," +
		"\"108\"," +
		"109," +
		"\"110\"," +
		"111.5," +
		"112.5," +
		"1," +
		"\"test\"," +
		"\"abcd\"," +
		"[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,111]," +
		"null," +
		"[112]," +
		"null," +
		"null," +
		"0," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"[201,202]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[\"foo\",\"bar\"]," +
		"[]," +
		"[]," +
		"null," +
		"[]," +
		"[0,2]" +
		"]"

	largeIntPBLiteGolden = "[null," +
		"null," +
		"null," +
		"null," +
		"\"" + oobJSStr + "\"," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"null," +
		"[]," +
		"[]," +
		maxSafeJSStr + "," +
		"\"" + oobJSStr + "\"" +
		"]"

	largeIntPBLiteZeroIndexGolden = "[" +
		"null," +
		"null," +
		"null," +
		"\"" + oobJSStr + "\"," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"null," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"[]," +
		"null," +
		"[]," +
		"[]," +
		maxSafeJSStr + "," +
		"\"" + oobJSStr + "\"" +
		"]"

	pbLitePackageGolden = "[null," +
		"1," +
		pbLiteGolden + "," +
		"[" + pbLiteGolden + "," + pbLiteGolden + "]" +
		"]"

	pbLitePackageZeroIndexGolden = "[" +
		"1," +
		pbLiteZeroIndexGolden + "," +
		"[" + pbLiteZeroIndexGolden + "," + pbLiteZeroIndexGolden + "]" +
		"]"

	objectKeyNameGolden = "{" +
		"\"optional_bool\":true," +
		"\"optional_bytes\":\"abcd\"," +
		"\"optional_double\":112.5," +
		"\"optional_fixed32\":107," +
		"\"optional_fixed64\":\"108\"," +
		"\"optional_float\":111.5," +
		"\"optional_int32\":101," +
		"\"optional_int64\":\"102\"," +
		"\"optional_nested_enum\":0," +
		"\"optional_nested_message\":{\"b\":112}," +
		"\"optional_sfixed32\":109," +
		"\"optional_sfixed64\":\"110\"," +
		"\"optional_sint32\":105," +
		"\"optional_sint64\":\"106\"," +
		"\"optional_string\":\"test\"," +
		"\"optional_uint32\":103," +
		"\"optional_uint64\":\"104\"," +
		"\"optionalgroup\":{\"a\":111}," +
		"\"repeated_int32\":[201,202]," +
		"\"repeated_nested_enum\":[0,2]," +
		"\"repeated_string\":[\"foo\",\"bar\"]" +
		"}"

	largeIntObjectKeyNameGolden = "{" +
		"\"optional_int64_number\":" + maxSafeJSStr + "," +
		"\"optional_int64_string\":\"" + oobJSStr + "\"," +
		"\"optional_uint64\":\"" + oobJSStr + "\"" +
		"}"

	objectKeyNamePackageGolden = "{" +
		"\"optional_int32\":1," +
		"\"other_all\":" + objectKeyNameGolden + "," +
		"\"rep_other_all\":[" + objectKeyNameGolden + "," + objectKeyNameGolden + "]" +
		"}"

	objectKeyTagGolden = "{" +
		"\"1\":101," +
		"\"10\":\"110\"," +
		"\"11\":111.5," +
		"\"12\":112.5," +
		"\"13\":true," +
		"\"14\":\"test\"," +
		"\"15\":\"abcd\"," +
		"\"16\":{\"17\":111}," +
		"\"18\":{\"1\":112}," +
		"\"2\":\"102\"," +
		"\"21\":0," +
		"\"3\":103," +
		"\"31\":[201,202]," +
		"\"4\":\"104\"," +
		"\"44\":[\"foo\",\"bar\"]," +
		"\"49\":[0,2]," +
		"\"5\":105," +
		"\"6\":\"106\"," +
		"\"7\":107," +
		"\"8\":\"108\"," +
		"\"9\":109" +
		"}"

	largeIntObjectKeyTagGolden = "{" +
		"\"4\":\"" + oobJSStr + "\"," +
		"\"50\":" + maxSafeJSStr + "," +
		"\"51\":\"" + oobJSStr + "\"" +
		"}"

	objectKeyTagPackageGolden = "{" +
		"\"1\":1," +
		"\"2\":" + objectKeyTagGolden + "," +
		"\"3\":[" + objectKeyTagGolden + "," + objectKeyTagGolden + "]" +
		"}"

	specialCharString         = "\x04\"\\/\b\f\n\r\tÄúɠ"
	objectKeyTagEscapesGolden = "{" +
		"\"14\":\"\\u0004\\\"\\\\/\\u0008\\u000c\\n\\r\\u0009Äúɠ\"," +
		"\"15\":\"\\u0004\\\"\\\\/\\u0008\\u000c\\n\\r\\u0009Äúɠ\"" +
		"}"
)

func populateMessage(pb *test_pb.TestAllTypes) {
	pb.OptionalInt32 = proto.Int32(101)
	pb.OptionalInt64 = proto.Int64(102)
	pb.OptionalUint32 = proto.Uint32(103)
	pb.OptionalUint64 = proto.Uint64(104)
	pb.OptionalSint32 = proto.Int32(105)
	pb.OptionalSint64 = proto.Int64(106)
	pb.OptionalFixed32 = proto.Uint32(107)
	pb.OptionalFixed64 = proto.Uint64(108)
	pb.OptionalSfixed32 = proto.Int32(109)
	pb.OptionalSfixed64 = proto.Int64(110)
	pb.OptionalFloat = proto.Float32(111.5)
	pb.OptionalDouble = proto.Float64(112.5)
	pb.OptionalBool = proto.Bool(true)
	pb.OptionalString = proto.String("test")
	pb.OptionalBytes = []byte("abcd")

	group := &test_pb.TestAllTypes_OptionalGroup{}
	group.A = proto.Int32(111)
	pb.Optionalgroup = group

	nestedMessage := &test_pb.TestAllTypes_NestedMessage{}
	nestedMessage.B = proto.Int32(112)
	pb.OptionalNestedMessage = nestedMessage

	pb.OptionalNestedEnum = test_pb.TestAllTypes_FOO.Enum()

	pb.RepeatedInt32 = append(pb.RepeatedInt32, 201)
	pb.RepeatedInt32 = append(pb.RepeatedInt32, 202)

	pb.RepeatedString = append(pb.RepeatedString, "foo")
	pb.RepeatedString = append(pb.RepeatedString, "bar")

	pb.RepeatedNestedEnum = append(pb.RepeatedNestedEnum,
		test_pb.TestAllTypes_FOO)
	pb.RepeatedNestedEnum = append(pb.RepeatedNestedEnum,
		test_pb.TestAllTypes_BAR)
}

func validateMessage(t *testing.T, pb *test_pb.TestAllTypes) {
	if pb == nil {
		t.Errorf("nil pb")
		t.FailNow()
	}
	if pb.OptionalInt32 == nil {
		t.Errorf("Field expected, OptionalInt32")
		t.FailNow()
	}
	if pb.OptionalInt64 == nil {
		t.Errorf("Field expected, OptionalInt64")
		t.FailNow()
	}
	if pb.OptionalUint32 == nil {
		t.Errorf("Field expected, OptionalUint32")
		t.FailNow()
	}
	if pb.OptionalUint64 == nil {
		t.Errorf("Field expected, OptionalUint64")
		t.FailNow()
	}
	if pb.OptionalSint32 == nil {
		t.Errorf("Field expected, OptionalSint32")
		t.FailNow()
	}
	if pb.OptionalFixed32 == nil {
		t.Errorf("Field expected, OptionalFixed32")
		t.FailNow()
	}
	if pb.OptionalFixed64 == nil {
		t.Errorf("Field expected, OptionalFixed64")
		t.FailNow()
	}
	if pb.OptionalSfixed32 == nil {
		t.Errorf("Field expected, OptionalSfixed32")
		t.FailNow()
	}
	if pb.OptionalSfixed64 == nil {
		t.Errorf("Field expected, OptionalSfixed64")
		t.FailNow()
	}
	if pb.OptionalFloat == nil {
		t.Errorf("Field expected, OptionalFloat")
		t.FailNow()
	}
	if pb.OptionalDouble == nil {
		t.Errorf("Field expected, OptionalDouble")
		t.FailNow()
	}
	if pb.OptionalBool == nil {
		t.Errorf("Field expected, OptionalBool")
		t.FailNow()
	}
	if pb.OptionalString == nil {
		t.Errorf("Field expected, OptionalString")
		t.FailNow()
	}
	if pb.OptionalBytes == nil {
		t.Errorf("Field expected, OptionalBytes")
		t.FailNow()
	}
	if pb.Optionalgroup == nil {
		t.Errorf("Field expected, Optionalgroup")
		t.FailNow()
	}
	if pb.OptionalNestedMessage == nil {
		t.Errorf("Field expected, OptionalNestedMessage")
		t.FailNow()
	}
	if pb.OptionalNestedEnum == nil {
		t.Errorf("Field expected, OptionalNestedEnum")
		t.FailNow()
	}

	if len(pb.RepeatedInt32) != 2 {
		t.Errorf("Found len %d, want 0 (RepeatedInt32)", len(pb.RepeatedInt32))
		t.FailNow()
	}
	if len(pb.RepeatedInt64) != 0 {
		t.Errorf("Found len %d, want 0 (RepeatedInt64)", len(pb.RepeatedInt64))
		t.FailNow()
	}
	if len(pb.RepeatedUint32) != 0 {
		t.Errorf("Found len %d, want 0 (RepeateUint32)", len(pb.RepeatedUint32))
		t.FailNow()
	}
	if len(pb.RepeatedUint64) != 0 {
		t.Errorf("Found len %d, want 0 (RepeatedUint64)", len(pb.RepeatedUint64))
		t.FailNow()
	}
	if len(pb.RepeatedSint32) != 0 {
		t.Errorf("Found len %d, want 0 (RepeatedSint32)", len(pb.RepeatedSint32))
		t.FailNow()
	}
	if len(pb.RepeatedSint64) != 0 {
		t.Errorf("Found len %d, want 0 (RepeatedSint64)", len(pb.RepeatedSint64))
		t.FailNow()
	}
	if len(pb.RepeatedFixed32) != 0 {
		t.Errorf("Found len %d, want 0 (RepeatedFixed32)", len(pb.RepeatedFixed32))
		t.FailNow()
	}
	if len(pb.RepeatedFixed64) != 0 {
		t.Errorf("Found len %d, want 0 (RepeatedFixed64)", len(pb.RepeatedFixed64))
		t.FailNow()
	}
	if len(pb.RepeatedSfixed32) != 0 {
		t.Errorf("Found len %d, want 0 (RepeatedSfixed32)",
			len(pb.RepeatedSfixed32))
		t.FailNow()
	}
	if len(pb.RepeatedSfixed64) != 0 {
		t.Errorf("Found len %d, want 0 (RepeatedSfixed64)",
			len(pb.RepeatedSfixed64))
		t.FailNow()
	}
	if len(pb.RepeatedFloat) != 0 {
		t.Errorf("Found len %d, want 0 (RepeatedFloat)", len(pb.RepeatedFloat))
		t.FailNow()
	}
	if len(pb.RepeatedDouble) != 0 {
		t.Errorf("Found len %d, want 0 (RepeatedDouble)", len(pb.RepeatedDouble))
		t.FailNow()
	}
	if len(pb.RepeatedBool) != 0 {
		t.Errorf("Found len %d, want 0 (RepeatedBool)", len(pb.RepeatedBool))
		t.FailNow()
	}
	if len(pb.RepeatedString) != 2 {
		t.Errorf("Found len %d, want 0 (RepeatedString)", len(pb.RepeatedString))
		t.FailNow()
	}
	if len(pb.RepeatedBytes) != 0 {
		t.Errorf("Found len %d, want 0 (RepeatedBytes)", len(pb.RepeatedBytes))
		t.FailNow()
	}
	if len(pb.Repeatedgroup) != 0 {
		t.Errorf("Found len %d, want 0 (Repeatedgroup)", len(pb.Repeatedgroup))
		t.FailNow()
	}
	if len(pb.RepeatedNestedMessage) != 0 {
		t.Errorf("Found len %d, want 0 (RepeatedNestedMessage)",
			len(pb.RepeatedNestedMessage))
		t.FailNow()
	}
	if len(pb.RepeatedNestedEnum) != 2 {
		t.Errorf("Found len %d, want 2 (RepeatedNestedEnum)",
			len(pb.RepeatedNestedEnum))
		t.FailNow()
	}

	if *pb.OptionalInt32 != 101 {
		t.Errorf("Found %d, want 101 (OptionalInt32)", *pb.OptionalInt32)
	}
	if *pb.OptionalInt64 != 102 {
		t.Errorf("Found %d, want 102 (OptionalInt64)", *pb.OptionalInt64)
	}
	if *pb.OptionalUint32 != 103 {
		t.Errorf("Found %d, want 103 (OptionalUint32)", *pb.OptionalUint32)
	}
	if *pb.OptionalUint64 != 104 {
		t.Errorf("Found %d, want 104 (OptionalUint64)", *pb.OptionalUint64)
	}
	if *pb.OptionalSint32 != 105 {
		t.Errorf("Found %d, want 105 (OptionalSint32)", *pb.OptionalSint32)
	}
	if *pb.OptionalSint64 != 106 {
		t.Errorf("Found %d, want 106 (OptionalSint64)", *pb.OptionalSint64)
	}
	if *pb.OptionalFixed32 != 107 {
		t.Errorf("Found %d, want 107 (OptionalFixed32)", *pb.OptionalFixed32)
	}
	if *pb.OptionalFixed64 != 108 {
		t.Errorf("Found %d, want 108 (OptionalFixed64)", *pb.OptionalFixed64)
	}
	if *pb.OptionalSfixed32 != 109 {
		t.Errorf("Found %d, want 109 (OptionalSfixed32)", *pb.OptionalSfixed32)
	}
	if *pb.OptionalSfixed64 != 110 {
		t.Errorf("Found %d, want 110 (OptionalSfixed64)", *pb.OptionalSfixed64)
	}
	if *pb.OptionalFloat != 111.5 {
		t.Errorf("Found %f, want 111.5 (OptionalFloat)", *pb.OptionalFloat)
	}
	if *pb.OptionalDouble != 112.5 {
		t.Errorf("Found %f, want 112.5 (OptionalDouble)", *pb.OptionalDouble)
	}
	if !*pb.OptionalBool {
		t.Errorf("Found %t, want true (OptionalBool)", *pb.OptionalBool)
	}
	if *pb.OptionalString != "test" {
		t.Errorf("Found %s, want 'test' (OptionalString)", *pb.OptionalString)
	}
	if !bytes.Equal(pb.OptionalBytes, []byte("abcd")) {
		t.Errorf("Found %s, want 'abcd' (OptionalBytes)", pb.OptionalBytes)
	}
	if *pb.Optionalgroup.A != 111 {
		t.Errorf("Found %d, want 111 (Optionalgroup.A)", *pb.Optionalgroup.A)
	}
	if *pb.OptionalNestedMessage.B != 112 {
		t.Errorf("Found %d, want 112 (OptionalNestedMessage.B)",
			*pb.OptionalNestedMessage.B)
	}
	if *pb.OptionalNestedEnum != test_pb.TestAllTypes_FOO {
		t.Errorf("Found %d, want FOO (OptionalNestedEnum)", *pb.OptionalNestedEnum)
	}
	if pb.RepeatedInt32[0] != 201 {
		t.Errorf("Found %d, want 201 (RepeatedInt32[0])", pb.RepeatedInt32[0])
	}
	if pb.RepeatedInt32[1] != 202 {
		t.Errorf("Found %d, want 202 (RepeatedInt32[1])", pb.RepeatedInt32[1])
	}
	if pb.RepeatedNestedEnum[0] != test_pb.TestAllTypes_FOO {
		t.Errorf("Found %d, want 201 (RepeatedInt32[0])", pb.RepeatedNestedEnum[0])
	}
	if pb.RepeatedNestedEnum[1] != test_pb.TestAllTypes_BAR {
		t.Errorf("Found %d, want 202 (RepeatedInt32[1])", pb.RepeatedNestedEnum[1])
	}
}

func TestMarshalPBLite(t *testing.T) {
	pb := &test_pb.TestAllTypes{}
	populateMessage(pb)

	s, err := MarshalPBLite(pb)
	if err != nil {
		t.Fatalf("unable to MarshalPBLite: %v", err)
	}
	if !bytes.Equal(s, []byte(pbLiteGolden)) {
		t.Errorf("Found %s, want %s", string(s), pbLiteGolden)
	}
}

func TestMarshalPBLiteLargeInt(t *testing.T) {
	pb := &test_pb.TestAllTypes{}
	pb.OptionalUint64 = proto.Uint64(oobJSInt)
	pb.OptionalInt64Number = proto.Int64(maxSafeJSInt)
	pb.OptionalInt64String = proto.Int64(oobJSInt)

	s, err := MarshalPBLite(pb)
	if err != nil {
		t.Fatalf("unable to MarshalPBLite: %v", err)
	}
	if !bytes.Equal(s, []byte(largeIntPBLiteGolden)) {
		t.Errorf("Found %s, want %s", string(s), largeIntPBLiteGolden)
	}
}

func TestMarshalPBLitePackage(t *testing.T) {
	pb := &package_test_pb.TestPackageTypes{}
	pb.OptionalInt32 = proto.Int32(1)
	testMessage := &test_pb.TestAllTypes{}
	populateMessage(testMessage)
	pb.OtherAll = testMessage
	pb.RepOtherAll = append(pb.RepOtherAll, testMessage)
	pb.RepOtherAll = append(pb.RepOtherAll, testMessage)

	s, err := MarshalPBLite(pb)
	if err != nil {
		t.Fatalf("unable to MarshalPBLite: %v", err)
	}
	if !bytes.Equal(s, []byte(pbLitePackageGolden)) {
		t.Errorf("Found %s, want %s", string(s), pbLitePackageGolden)
	}
}

func TestUnmarshalPBLite(t *testing.T) {
	pb := &test_pb.TestAllTypes{}
	err := UnmarshalPBLite([]byte(pbLiteGolden), pb)
	if err != nil {
		t.Fatalf("unable to UnmarshalPBLite: %v", err)
	}
	validateMessage(t, pb)
}

func TestUnmarshalPBLiteLargeInt(t *testing.T) {
	pb := &test_pb.TestAllTypes{}
	err := UnmarshalPBLite([]byte(largeIntPBLiteGolden), pb)
	if err != nil {
		t.Fatalf("unable to UnmarshalPBLite: %v", err)
	}
	if pb.OptionalUint64 == nil {
		t.Errorf("Field expected, OptionalUint64")
		t.FailNow()
	}
	if pb.OptionalInt64Number == nil {
		t.Errorf("Field expected, OptionalInt64Number")
		t.FailNow()
	}
	if pb.OptionalInt64String == nil {
		t.Errorf("Field expected, OptionalInt64String")
		t.FailNow()
	}
	if *pb.OptionalUint64 != oobJSInt {
		t.Errorf("Found %d, want %d", *pb.OptionalUint64, oobJSInt)
	}
	if *pb.OptionalInt64Number != maxSafeJSInt {
		t.Errorf("Found %d, want %d", *pb.OptionalInt64Number, maxSafeJSInt)
	}
	if *pb.OptionalInt64String != oobJSInt {
		t.Errorf("Found %d, want %d", *pb.OptionalInt64String, oobJSInt)
	}
}

func TestUnmarshalPBLitePackage(t *testing.T) {
	pb := &package_test_pb.TestPackageTypes{}
	err := UnmarshalPBLite([]byte(pbLitePackageGolden), pb)
	if err != nil {
		t.Fatalf("unable to UnmarshalPBLite: %v", err)
	}
	if *pb.OptionalInt32 != 1 {
		t.Errorf("Found %d, want 1", *pb.OptionalInt32)
	}
	validateMessage(t, pb.OtherAll)
}

func TestMarshalPBLiteZeroIndex(t *testing.T) {
	pb := &test_pb.TestAllTypes{}
	populateMessage(pb)

	s, err := MarshalPBLiteZeroIndex(pb)
	if err != nil {
		t.Fatalf("unable to MarshalPBLiteZeroIndex: %v", err)
	}
	if !bytes.Equal(s, []byte(pbLiteZeroIndexGolden)) {
		t.Errorf("Found %s, want %s", string(s), pbLiteZeroIndexGolden)
	}
}

func TestMarshalPBLiteZeroIndexLargeInt(t *testing.T) {
	pb := &test_pb.TestAllTypes{}
	pb.OptionalUint64 = proto.Uint64(oobJSInt)
	pb.OptionalInt64Number = proto.Int64(maxSafeJSInt)
	pb.OptionalInt64String = proto.Int64(oobJSInt)

	s, err := MarshalPBLiteZeroIndex(pb)
	if err != nil {
		t.Fatalf("unable to MarshalPBLiteZeroIndex: %v", err)
	}
	if !bytes.Equal(s, []byte(largeIntPBLiteZeroIndexGolden)) {
		t.Errorf("Found %s, want %s", string(s), largeIntPBLiteZeroIndexGolden)
	}
}

func TestMarshalPBLiteZeroIndexPackage(t *testing.T) {
	pb := &package_test_pb.TestPackageTypes{}
	pb.OptionalInt32 = proto.Int32(1)
	testMessage := &test_pb.TestAllTypes{}
	populateMessage(testMessage)
	pb.OtherAll = testMessage
	pb.RepOtherAll = append(pb.RepOtherAll, testMessage)
	pb.RepOtherAll = append(pb.RepOtherAll, testMessage)
	s, err := MarshalPBLiteZeroIndex(pb)
	if err != nil {
		t.Fatalf("unable to MarshalPBLiteZeroIndex: %v", err)
	}
	if !bytes.Equal(s, []byte(pbLitePackageZeroIndexGolden)) {
		t.Errorf("Found %s, want %s", string(s), pbLitePackageZeroIndexGolden)
	}
}

func TestUnmarshalPBLiteZeroIndex(t *testing.T) {
	pb := &test_pb.TestAllTypes{}
	err := UnmarshalPBLiteZeroIndex([]byte(pbLiteZeroIndexGolden), pb)
	if err != nil {
		t.Fatalf("unable to UnmarshalPBLiteZeroIndex: %v", err)
	}
	validateMessage(t, pb)
}

func TestUnmarshalPBLiteZeroIndexLargeInt(t *testing.T) {
	pb := &test_pb.TestAllTypes{}
	err := UnmarshalPBLiteZeroIndex([]byte(largeIntPBLiteZeroIndexGolden), pb)
	if err != nil {
		t.Fatalf("unable to UnmarshalPBLiteZeroIndex: %v", err)
	}
	if pb.OptionalUint64 == nil {
		t.Errorf("Field expected, OptionalUint64")
		t.FailNow()
	}
	if pb.OptionalInt64Number == nil {
		t.Errorf("Field expected, OptionalInt64Number")
		t.FailNow()
	}
	if pb.OptionalInt64String == nil {
		t.Errorf("Field expected, OptionalInt64String")
		t.FailNow()
	}
	if *pb.OptionalUint64 != oobJSInt {
		t.Errorf("Found %d, want %d", *pb.OptionalUint64, oobJSInt)
	}
	if *pb.OptionalInt64Number != maxSafeJSInt {
		t.Errorf("Found %d, want %d", *pb.OptionalInt64Number, maxSafeJSInt)
	}
	if *pb.OptionalInt64String != oobJSInt {
		t.Errorf("Found %d, want %d", *pb.OptionalInt64String, oobJSInt)
	}
}

func TestUnmarshalPBLitePackageZeroIndex(t *testing.T) {
	pb := &package_test_pb.TestPackageTypes{}
	err := UnmarshalPBLiteZeroIndex([]byte(pbLitePackageZeroIndexGolden), pb)
	if err != nil {
		t.Fatalf("unable to UnmarshalPBLiteZeroIndex: %v", err)
	}
	if *pb.OptionalInt32 != 1 {
		t.Errorf("Found %d, want 1", *pb.OptionalInt32)
	}
	validateMessage(t, pb.OtherAll)
}

func TestMarshalObjectKeyName(t *testing.T) {
	pb := &test_pb.TestAllTypes{}
	populateMessage(pb)

	s, err := MarshalObjectKeyName(pb)
	if err != nil {
		t.Fatalf("unable to MarshalObjectKeyName: %v", err)
	}
	if !bytes.Equal(s, []byte(objectKeyNameGolden)) {
		t.Errorf("Found %s, want %s", string(s), objectKeyNameGolden)
	}
}

func TestMarshalObjectKeyNameLargeInt(t *testing.T) {
	pb := &test_pb.TestAllTypes{}
	pb.OptionalUint64 = proto.Uint64(oobJSInt)
	pb.OptionalInt64Number = proto.Int64(maxSafeJSInt)
	pb.OptionalInt64String = proto.Int64(oobJSInt)

	s, err := MarshalObjectKeyName(pb)
	if err != nil {
		t.Fatalf("unable to MarshalObjectKeyName: %v", err)
	}
	if !bytes.Equal(s, []byte(largeIntObjectKeyNameGolden)) {
		t.Errorf("Found %s, want %s", string(s), largeIntObjectKeyNameGolden)
	}
}

func TestMarshalObjectKeyNamePackage(t *testing.T) {
	pb := &package_test_pb.TestPackageTypes{}
	pb.OptionalInt32 = proto.Int32(1)
	testMessage := &test_pb.TestAllTypes{}
	populateMessage(testMessage)
	pb.OtherAll = testMessage
	pb.RepOtherAll = append(pb.RepOtherAll, testMessage)
	pb.RepOtherAll = append(pb.RepOtherAll, testMessage)
	s, err := MarshalObjectKeyName(pb)
	if err != nil {
		t.Fatalf("unable to MarshalObjectKeyName: %v", err)
	}
	if !bytes.Equal(s, []byte(objectKeyNamePackageGolden)) {
		t.Errorf("Found %s, want %s", string(s), objectKeyNamePackageGolden)
	}
}

func TestUnmarshalObjectKeyName(t *testing.T) {
	pb := &test_pb.TestAllTypes{}
	err := UnmarshalObjectKeyName([]byte(objectKeyNameGolden), pb)
	if err != nil {
		t.Fatalf("unable to UnmarshalObjectKeyName: %v", err)
	}
	validateMessage(t, pb)
}

func TestUnmarshalObjectKeyNameLargeInt(t *testing.T) {
	pb := &test_pb.TestAllTypes{}
	err := UnmarshalObjectKeyName([]byte(largeIntObjectKeyNameGolden), pb)
	if err != nil {
		t.Fatalf("unable to UnmarshalObjectKeyName: %v", err)
	}
	if pb.OptionalUint64 == nil {
		t.Errorf("Field expected, OptionalUint64")
		t.FailNow()
	}
	if pb.OptionalInt64Number == nil {
		t.Errorf("Field expected, OptionalInt64Number")
		t.FailNow()
	}
	if pb.OptionalInt64String == nil {
		t.Errorf("Field expected, OptionalInt64String")
		t.FailNow()
	}
	if *pb.OptionalUint64 != oobJSInt {
		t.Errorf("Found %d, want %d", *pb.OptionalUint64, oobJSInt)
	}
	if *pb.OptionalInt64Number != maxSafeJSInt {
		t.Errorf("Found %d, want %d", *pb.OptionalInt64Number, maxSafeJSInt)
	}
	if *pb.OptionalInt64String != oobJSInt {
		t.Errorf("Found %d, want %d", *pb.OptionalInt64String, oobJSInt)
	}
}

func TestUnmarshalObjectKeyNamePackage(t *testing.T) {
	pb := &package_test_pb.TestPackageTypes{}
	err := UnmarshalObjectKeyName([]byte(objectKeyNamePackageGolden), pb)
	if err != nil {
		t.Fatalf("unable to UnmarshalObjectKeyName: %v", err)
	}
	if pb.OptionalInt32 == nil {
		t.Errorf("Field expected, OptionalInt32")
		t.FailNow()
	}
	if *pb.OptionalInt32 != 1 {
		t.Errorf("Found %d, want 1", *pb.OptionalInt32)
	}
	validateMessage(t, pb.OtherAll)
}

func TestMarshalObjectKeyTag(t *testing.T) {
	pb := &test_pb.TestAllTypes{}
	populateMessage(pb)

	s, err := MarshalObjectKeyTag(pb)
	if err != nil {
		t.Fatalf("unable to MarshalObjectKeyTag: %v", err)
	}
	if !bytes.Equal(s, []byte(objectKeyTagGolden)) {
		t.Errorf("Found %s, want %s", string(s), objectKeyTagGolden)
	}
}

func TestMarshalObjectKeyTagLargeInt(t *testing.T) {
	pb := &test_pb.TestAllTypes{}
	pb.OptionalUint64 = proto.Uint64(oobJSInt)
	pb.OptionalInt64Number = proto.Int64(maxSafeJSInt)
	pb.OptionalInt64String = proto.Int64(oobJSInt)

	s, err := MarshalObjectKeyTag(pb)
	if err != nil {
		t.Fatalf("unable to MarshalObjectKeyTag: %v", err)
	}
	if !bytes.Equal(s, []byte(largeIntObjectKeyTagGolden)) {
		t.Errorf("Found %s, want %s", string(s), largeIntObjectKeyTagGolden)
	}
}

func TestMarshalObjectKeyTagPackage(t *testing.T) {
	pb := &package_test_pb.TestPackageTypes{}
	pb.OptionalInt32 = proto.Int32(1)
	testMessage := &test_pb.TestAllTypes{}
	populateMessage(testMessage)
	pb.OtherAll = testMessage
	pb.RepOtherAll = append(pb.RepOtherAll, testMessage)
	pb.RepOtherAll = append(pb.RepOtherAll, testMessage)
	s, err := MarshalObjectKeyTag(pb)
	if err != nil {
		t.Fatalf("unable to MarshalObjectKeyTag: %v", err)
	}
	if !bytes.Equal(s, []byte(objectKeyTagPackageGolden)) {
		t.Errorf("Found %s, want %s", string(s), objectKeyTagPackageGolden)
	}
}

func TestUnmarshalObjectKeyTag(t *testing.T) {
	pb := &test_pb.TestAllTypes{}
	err := UnmarshalObjectKeyTag([]byte(objectKeyTagGolden), pb)
	if err != nil {
		t.Fatalf("unable to UnmarshalObjectKeyTag: %v", err)
	}
	validateMessage(t, pb)
}

func TestUnmarshalObjectKeyTagLargeInt(t *testing.T) {
	pb := &test_pb.TestAllTypes{}
	err := UnmarshalObjectKeyTag([]byte(largeIntObjectKeyTagGolden), pb)
	if err != nil {
		t.Fatalf("unable to UnmarshalObjectKeyTag: %v", err)
	}
	if pb.OptionalUint64 == nil {
		t.Errorf("Field expected, OptionalUint64")
		t.FailNow()
	}
	if pb.OptionalInt64Number == nil {
		t.Errorf("Field expected, OptionalInt64Number")
		t.FailNow()
	}
	if pb.OptionalInt64String == nil {
		t.Errorf("Field expected, OptionalInt64String")
		t.FailNow()
	}
	if *pb.OptionalUint64 != oobJSInt {
		t.Errorf("Found %d, want %d", *pb.OptionalUint64, oobJSInt)
	}
	if *pb.OptionalInt64Number != maxSafeJSInt {
		t.Errorf("Found %d, want %d", *pb.OptionalInt64Number, maxSafeJSInt)
	}
	if *pb.OptionalInt64String != oobJSInt {
		t.Errorf("Found %d, want %d", *pb.OptionalInt64String, oobJSInt)
	}
}

func TestUnmarshalObjectKeyTagPackage(t *testing.T) {
	pb := &package_test_pb.TestPackageTypes{}
	err := UnmarshalObjectKeyTag([]byte(objectKeyTagPackageGolden), pb)
	if err != nil {
		t.Fatalf("unable to UnmarshalObjectKeyTag: %v", err)
	}
	if pb.OptionalInt32 == nil {
		t.Errorf("Field expected, OptionalInt32")
		t.FailNow()
	}
	if *pb.OptionalInt32 != 1 {
		t.Errorf("Found %d, want 1", *pb.OptionalInt32)
	}
	validateMessage(t, pb.OtherAll)
}

// func TestMarshalObjectKeyTagEscape(t *testing.T) {
// 	pb := &test_pb.TestAllTypes{}
// 	pb.OptionalString = proto.String(specialCharString)
// 	pb.OptionalBytes = []byte(specialCharString)
//
// 	s, err := MarshalObjectKeyTag(pb)
// 	if err != nil {
// 		t.Fatalf("unable to MarshalObjectKeyTag: %v", err)
// 	}
// 	if !bytes.Equal(s, []byte(objectKeyTagEscapesGolden)) {
// 		t.Errorf("Found %s, want %s", string(s), objectKeyTagEscapesGolden)
// 	}
// }

func TestUnmarshalObjectKeyTagEscape(t *testing.T) {
	pb := &test_pb.TestAllTypes{}
	err := UnmarshalObjectKeyTag([]byte(objectKeyTagEscapesGolden), pb)
	if err != nil {
		t.Fatalf("unable to MarshalObjectKeyTag: %v", err)
	}
	if pb.OptionalString == nil {
		t.Errorf("Field expected, OptionalString")
		t.FailNow()
	}
	if *pb.OptionalString != specialCharString {
		t.Errorf("Found %s, want %s", *pb.OptionalString, specialCharString)
	}
	if !bytes.Equal(pb.OptionalBytes, []byte(specialCharString)) {
		t.Errorf("Found %s, want %s", string(pb.OptionalBytes), specialCharString)
	}
}
