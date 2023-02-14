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

	i              int
	ei             int
	currentMsgName string

	protoName string
	messages  map[string][]ProtoToken
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

func (t *Token) takeString(h, e string) string {
	r := strings.Index(t.fileString[t.ei:], h)
	if r < 0 {
		fmt.Println("[warn] target first index is -1")
	}

	t.i = t.ei + len(h) + r

	r = strings.Index(t.fileString[t.i:], e)
	if r < 0 {
		fmt.Println("[warn] target last index is -1")
	}

	t.ei = t.i + r
	return t.fileString[t.i:t.ei]
}

func (t *Token) readMessageName() {
	t.currentMsgName = t.takeString(" extends MessageMicro<", ">")
	t.messages[t.currentMsgName] = []ProtoToken{}
}

func (t *Token) readFieldMapID() {
	fmstr := t.takeString("__fieldMap__ = MessageMicro.initFieldMap(new int[]{", "}") // FieldMap value string
	for _, fme := range strings.Split(fmstr, ", ") {
		i, err := strconv.ParseUint(fme, 10, 64)
		if err != nil {
			fmt.Println("[warn] parse fieldmap string to int token failure in", fme, "of", t.currentMsgName, ".")
			continue
		}

		t.messages[t.currentMsgName] = append(t.messages[t.currentMsgName], ProtoToken{tag: i >> 3})
	}
	return
}

func (t *Token) readFieldNames() {
	fmstr := t.takeString("}, new String[]{", "}, new Object[]{") // FieldMap key string
	for i, fme := range strings.Split(fmstr, ", ") {
		t.messages[t.currentMsgName][i].name = fme[1 : len(fme)-1]
	}
	return
}
