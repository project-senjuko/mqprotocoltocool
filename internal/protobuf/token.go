////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Copyright (c) 2023-present qianjunakasumi <i@qianjunakasumi.ren>                                                    /
//                            project-senjuko/mqprotocoltocool Contributors                                            /
//                                                                                                                     /
//           https://github.com/qianjunakasumi                                                                         /
//           https://github.com/project-senjuko/mqprotocoltocool/graphs/contributors                                   /
//                                                                                                                     /
//   This Source Code Form is subject to the terms of the Mozilla Public                                               /
//   License, v. 2.0. If a copy of the MPL was not distributed with this                                               /
//   file, You can obtain one at http://mozilla.org/MPL/2.0/.                                                          /
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

package protobuf

import (
	"fmt"
	"strconv"
	"strings"
)

func NewToken(f string) *Token {
	return &Token{
		fileString: f,
		messages:   map[string][]ProtoToken{},
	}
}

func (t *Token) ReadAll() {
	if !t.isProtobufFile() {
		fmt.Println("[warn] not a protobuf file")
		return
	}

	ok := t.readProtobufName()
	if !ok {
		fmt.Println("[error] cannot get protobuf file name")
		return
	}

	fmt.Println("[info] success read protobuf file name:", t.protobufName)

	for t.readMessageName() {
		if t.readFieldTags() {
			t.readFieldNames()
			t.readFieldTypes()
			fmt.Println("[info] success read protobuf message:", t.currentMsgName)
			continue
		}

		fmt.Println("[info]", t.currentMsgName, "is empty message")
	}
}

func (t *Token) takeString(h, e string) (string, bool) {
	r := strings.Index(t.fileString[t.ei:], h)
	if r < 0 {
		return "", false
	}

	t.i = t.ei + len(h) + r
	r = strings.Index(t.fileString[t.i:], e)
	if r < 0 {
		return "", false
	}

	t.ei = t.i + r
	return t.fileString[t.i:t.ei], true
}

func (t *Token) isProtobufFile() bool {
	if strings.Index(t.fileString, "import com.tencent.mobileqq.pb.MessageMicro;") < 0 {
		return false
	}
	return true
}

func (t *Token) readProtobufName() bool {
	var ok bool
	t.protobufName, ok = t.takeString("public final class ", " {")
	return ok
}

func (t *Token) readMessageName() bool {
	var ok bool
	t.currentMsgName, ok = t.takeString(" extends MessageMicro<", ">")
	if !ok {
		fmt.Println("[warn] cannot find more protobuf message")
		return false
	}
	t.messages[t.currentMsgName] = []ProtoToken{}
	return true
}

func (t *Token) readFieldTags() bool {
	fmstr, _ := t.takeString("__fieldMap__ = MessageMicro.initFieldMap(new int[", "}") // FieldMap value string
	if fmstr[:1] == "0" {
		return false
	}

	for _, fme := range strings.Split(fmstr[2:], ", ") {
		i, err := strconv.ParseUint(fme, 10, 64)
		if err != nil {
			fmt.Println("[erro] parse fieldmap string to int token failure in", fme, "of", t.currentMsgName, ".")
			continue
		}

		t.messages[t.currentMsgName] = append(t.messages[t.currentMsgName], ProtoToken{tag: i >> 3})
	}

	return true
}

func (t *Token) readFieldNames() {
	fn, _ := t.takeString("}, new String[]{", "}, new Object[]{") // FieldMap key string
	for i, fne := range strings.Split(fn, ", ") {
		t.messages[t.currentMsgName][i].name = fne[1 : len(fne)-1]
	}
	return
}

func (t *Token) readFieldTypes() {
	for i, f := range t.messages[t.currentMsgName] {
		ft, _ := t.takeString(f.name+" = ", "(")
		t.messages[t.currentMsgName][i].typ = ft
	}
}
