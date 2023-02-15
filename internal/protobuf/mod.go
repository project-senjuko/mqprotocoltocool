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

type Token struct {
	fileString string

	protobufName string
	messages     ProtoFile

	i              int
	ei             int
	currentMsgName string
}

type ProtoFile map[string][]ProtoToken

type ProtoToken struct {
	name string
	tag  uint64
	typ  string
}
