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
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

func NewToken(f string) *Token {
	return &Token{
		fileString: f,
		messages:   map[string][]ProtoToken{},
	}
}

func (t *Token) ReadAll() {
	if !t.isProtobufFile() {
		log.Warn().Msg("not a protobuf file")
		return
	}

	ok := t.readProtobufName()
	if !ok {
		log.Error().Msg("cannot get protobuf file name")
		return
	}

	log.Info().Str("name", t.protobufName).Msg("success read protobuf file name")

	for t.readMessageName() {
		if !t.readFieldTags() {
			log.Info().Str("msg", t.currentMsgName).Msg("is empty")
			continue
		}
		if !t.readFieldNames() {
			continue
		}

		t.readFieldTypes()
		log.Info().Str("msg", t.currentMsgName).Msg("success read protobuf message")
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
		log.Warn().Msg("cannot find more protobuf message")
		return false
	}
	t.messages[t.currentMsgName] = []ProtoToken{}
	return true
}

func (t *Token) readFieldTags() bool {
	fmstr, ok := t.takeString("__fieldMap__ = MessageMicro.initFieldMap(new int[", "}") // FieldMap value string
	if !ok || fmstr[:1] == "0" {
		return false
	}

	for _, fme := range strings.Split(fmstr[2:], ", ") {
		i, err := strconv.ParseUint(fme, 10, 64)
		if err != nil {
			log.Error().Str("token", fme).Str("msg", t.currentMsgName).
				Msg("parse fieldmap tag string token to int failure")
			continue
		}

		t.messages[t.currentMsgName] = append(t.messages[t.currentMsgName], ProtoToken{tag: i >> 3})
	}

	return true
}

func (t *Token) readFieldNames() bool {
	fn, ok := t.takeString("}, new String[]{", "}, new Object[]{") // FieldMap key string
	if !ok {
		log.Error().Str("msg", t.currentMsgName).Msg("wtf the field names")
		return false
	}
	for i, fne := range strings.Split(fn, ", ") {
		t.messages[t.currentMsgName][i].name = strings.ReplaceAll(fne, `"`, "")
	}
	return true
}

func (t *Token) readFieldTypes() {
	for i, f := range t.messages[t.currentMsgName] {
		ft, ok := t.takeString(f.name+" = ", "(")
		if !ok {
			log.Error().Str("field", f.name).Msg("wtf the field type")
			t.messages[t.currentMsgName][i].typ = f.name
		}
		t.messages[t.currentMsgName][i].typ = ft
	}
}
