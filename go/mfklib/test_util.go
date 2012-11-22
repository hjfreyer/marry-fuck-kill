package mfklib

import (
	"code.google.com/p/goprotobuf/proto"
	. "launchpad.net/gocheck"
)

// -----------------------------------------------------------------------
// ProtoEquals checker.

type protoEqualsChecker struct {
	*CheckerInfo
}

// The ProtoEquals checker verifies that the obtained and expected values are
// both Google protocol buffers, and that they are proto-equal.
var ProtoEquals Checker = &protoEqualsChecker{
	&CheckerInfo{Name: "ProtoEquals", Params: []string{"obtained", "expected"}},
}

func (checker *protoEqualsChecker) Check(params []interface{}, names []string) (result bool, error string) {
	obtained, ok := params[0].(proto.Message)
	if !ok {
		return false, "obtained is not a proto.Message"
	}

	expected, ok := params[1].(proto.Message)
	if !ok {
		return false, "expected is not a proto.Message"
	}

	params[0] = proto.MarshalTextString(obtained)
	params[1] = proto.MarshalTextString(expected)

	return proto.Equal(obtained, expected), ""
}
