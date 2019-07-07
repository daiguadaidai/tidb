// Copyright 2019 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"github.com/daiguadaidai/parser/mysql"
	. "github.com/pingcap/check"
)

type ColumnTestSuite struct {
}

var _ = Suite(new(ColumnTestSuite))

func (s ColumnTestSuite) TestDumpColumn(c *C) {
	info := ColumnInfo{
		Schema:             "testSchema",
		Table:              "testTable",
		OrgTable:           "testOrgTable",
		Name:               "testName",
		OrgName:            "testOrgName",
		ColumnLength:       1,
		Charset:            106,
		Flag:               0,
		Decimal:            1,
		Type:               14,
		DefaultValueLength: 2,
		DefaultValue:       []byte{5, 2},
	}
	r := info.Dump(nil)
	exp := []byte{0x3, 0x64, 0x65, 0x66, 0xa, 0x74, 0x65, 0x73, 0x74, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x9, 0x74, 0x65, 0x73, 0x74, 0x54, 0x61, 0x62, 0x6c, 0x65, 0xc, 0x74, 0x65, 0x73, 0x74, 0x4f, 0x72, 0x67, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x8, 0x74, 0x65, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0xb, 0x74, 0x65, 0x73, 0x74, 0x4f, 0x72, 0x67, 0x4e, 0x61, 0x6d, 0x65, 0xc, 0x6a, 0x0, 0x1, 0x0, 0x0, 0x0, 0xe, 0x0, 0x0, 0x1, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5, 0x2}
	c.Assert(r, DeepEquals, exp)

	c.Assert(dumpFlag(mysql.TypeSet, 0), Equals, uint16(mysql.SetFlag))
	c.Assert(dumpFlag(mysql.TypeEnum, 0), Equals, uint16(mysql.EnumFlag))
	c.Assert(dumpFlag(mysql.TypeString, 0), Equals, uint16(0))

	c.Assert(dumpType(mysql.TypeSet), Equals, mysql.TypeString)
	c.Assert(dumpType(mysql.TypeEnum), Equals, mysql.TypeString)
	c.Assert(dumpType(mysql.TypeBit), Equals, mysql.TypeBit)
}
