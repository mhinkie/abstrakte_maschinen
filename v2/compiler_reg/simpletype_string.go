// Code generated by "stringer -type=SimpleType"; DO NOT EDIT.

package main

import "strconv"

const _SimpleType_name = "NothingIntStr"

var _SimpleType_index = [...]uint8{0, 7, 10, 13}

func (i SimpleType) String() string {
	if i < 0 || i >= SimpleType(len(_SimpleType_index)-1) {
		return "SimpleType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SimpleType_name[_SimpleType_index[i]:_SimpleType_index[i+1]]
}
