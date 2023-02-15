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

type Token struct {
	fileString string

	protoName string
	messages  map[string][]ProtoToken

	i              int
	ei             int
	currentMsgName string
}

type ProtoToken struct {
	name string
	tag  uint64
	typ  string
}

func NewToken(f string) *Token {
	return &Token{
		fileString: f,
		messages:   map[string][]ProtoToken{},
	}
}

// 检测 import 是否有 MessageMicro，没有则跳过

func (t *Token) ReadFromPath() {}

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

func (t *Token) ReadAll() {
	for t.readMessageName() {
		if t.readFieldMapID() {
			t.readFieldNames()
			t.readFieldType()
			fmt.Println("[info] success read protobuf message:", t.currentMsgName)
			continue
		}

		fmt.Println("[info]", t.currentMsgName, "is empty message")
	}
}

func (t *Token) readMessageName() bool {
	var ok bool
	t.currentMsgName, ok = t.takeString(" extends MessageMicro<", ">")
	if !ok {
		fmt.Println("[warn] cannot find protobuf message")
		return false
	}
	t.messages[t.currentMsgName] = []ProtoToken{}
	return true
}

func (t *Token) readFieldMapID() bool {
	fmstr, _ := t.takeString("__fieldMap__ = MessageMicro.initFieldMap(new int[", "}") // FieldMap value string
	if fmstr[:1] == "0" {
		return false
	}

	for _, fme := range strings.Split(fmstr[2:], ", ") {
		i, err := strconv.ParseUint(fme, 10, 64)
		if err != nil {
			fmt.Println("[warn] parse fieldmap string to int token failure in", fme, "of", t.currentMsgName, ".")
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

func (t *Token) readFieldType() {
	for i, f := range t.messages[t.currentMsgName] {
		ft, _ := t.takeString(f.name+" = ", "(")
		t.messages[t.currentMsgName][i].typ = ft
	}
}
