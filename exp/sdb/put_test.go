// Copyright 2012 Aaron Jacobs. All Rights Reserved.
// Author: aaronjjacobs@gmail.com (Aaron Jacobs)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sdb

import (
	"github.com/jacobsa/aws/exp/sdb/conn"
	. "github.com/jacobsa/oglematchers"
	. "github.com/jacobsa/ogletest"
	"sort"
	"testing"
)

func TestPut(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

func getSortedKeys(r conn.Request) []string {
	result := sort.StringSlice{}
	for key, _ := range r {
		result = append(result, key)
	}

	sort.Sort(result)
	return result
}

////////////////////////////////////////////////////////////////////////
// PutAttributes
////////////////////////////////////////////////////////////////////////

type PutTest struct {
	domainTest

	item ItemName
	updates []PutUpdate
	preconditions []Precondition

	err error
}

func init() { RegisterTestSuite(&PutTest{}) }

func (t *PutTest) SetUp(i *TestInfo) {
	// Call common setup code.
	t.domainTest.SetUp(i)

	// Make the request legal by default.
	t.item = "foo"
	t.updates = []PutUpdate{PutUpdate{"bar", "baz", false}}
}

func (t *PutTest) callDomain() {
	t.err = t.domain.PutAttributes(t.item, t.updates, t.preconditions)
}

func (t *PutTest) EmptyItemName() {
	t.item = ""

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("item name")))
}

func (t *PutTest) InvalidItemName() {
	t.item = "taco\x80\x81\x82"

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("item name")))
	ExpectThat(t.err, Error(HasSubstr(string(t.item))))
}

func (t *PutTest) ZeroUpdates() {
	t.updates = []PutUpdate{}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("number")))
	ExpectThat(t.err, Error(HasSubstr("updates")))
	ExpectThat(t.err, Error(HasSubstr("0")))
}

func (t *PutTest) TooManyUpdates() {
	t.updates = make([]PutUpdate, 257)

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("number")))
	ExpectThat(t.err, Error(HasSubstr("updates")))
	ExpectThat(t.err, Error(HasSubstr("256")))
}

func (t *PutTest) OneAttributeNameEmpty() {
	t.updates = []PutUpdate{
		PutUpdate{Name: "foo"},
		PutUpdate{Name: "", Value: "taco"},
		PutUpdate{Name: "bar"},
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("attribute")))
	ExpectThat(t.err, Error(HasSubstr("name")))
	ExpectThat(t.err, Error(HasSubstr("taco")))
}

func (t *PutTest) OneAttributeNameInvalid() {
	t.updates = []PutUpdate{
		PutUpdate{Name: "foo"},
		PutUpdate{Name: "taco\x80\x81\x82"},
		PutUpdate{Name: "bar"},
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("attribute")))
	ExpectThat(t.err, Error(HasSubstr("name")))
	ExpectThat(t.err, Error(HasSubstr(t.updates[1].Name)))
}

func (t *PutTest) OneAttributeValueInvalid() {
	t.updates = []PutUpdate{
		PutUpdate{Name: "foo"},
		PutUpdate{Name: "bar", Value: "taco\x80\x81\x82"},
		PutUpdate{Name: "baz"},
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("attribute")))
	ExpectThat(t.err, Error(HasSubstr("value")))
	ExpectThat(t.err, Error(HasSubstr(t.updates[1].Value)))
}

func (t *PutTest) OnePreconditionNameEmpty() {
	t.preconditions = []Precondition{
		Precondition{Name: "foo", Exists: new(bool)},
		Precondition{Name: "", Exists: new(bool)},
		Precondition{Name: "baz", Exists: new(bool)},
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("attribute")))
	ExpectThat(t.err, Error(HasSubstr("name")))
}

func (t *PutTest) OnePreconditionNameInvalid() {
	t.preconditions = []Precondition{
		Precondition{Name: "foo", Exists: new(bool)},
		Precondition{Name: "taco\x80\x81\x82", Exists: new(bool)},
		Precondition{Name: "baz", Exists: new(bool)},
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("attribute")))
	ExpectThat(t.err, Error(HasSubstr("name")))
	ExpectThat(t.err, Error(HasSubstr(t.preconditions[1].Name)))
}

func (t *PutTest) OnePreconditionValueInvalid() {
	t.preconditions = []Precondition{
		Precondition{Name: "foo", Value: new(string)},
		Precondition{Name: "bar", Value: new(string)},
		Precondition{Name: "baz", Value: new(string)},
	}

	*t.preconditions[0].Value = ""
	*t.preconditions[1].Value = "taco\x80\x81\x82"
	*t.preconditions[2].Value = "qux"

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("attribute")))
	ExpectThat(t.err, Error(HasSubstr("value")))
	ExpectThat(t.err, Error(HasSubstr(*t.preconditions[1].Value)))
}

func (t *PutTest) OnePreconditionMissingOperand() {
	t.preconditions = []Precondition{
		Precondition{Name: "foo", Exists: new(bool)},
		Precondition{Name: "bar"},
		Precondition{Name: "baz", Exists: new(bool)},
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("precondition")))
	ExpectThat(t.err, Error(HasSubstr("bar")))
}

func (t *PutTest) OnePreconditionHasTwoOperands() {
	t.preconditions = []Precondition{
		Precondition{Name: "foo", Exists: new(bool)},
		Precondition{Name: "bar", Exists: new(bool), Value: new(string)},
		Precondition{Name: "baz", Exists: new(bool)},
	}

	// Call
	t.callDomain()

	ExpectThat(t.err, Error(HasSubstr("Invalid")))
	ExpectThat(t.err, Error(HasSubstr("precondition")))
	ExpectThat(t.err, Error(HasSubstr("bar")))
}

func (t *PutTest) BasicParameters() {
	t.item = "some_item"
	t.updates = []PutUpdate{
		PutUpdate{Name: "foo"},
		PutUpdate{Name: "bar", Value: "taco", Replace: true},
		PutUpdate{Name: "baz", Value: "burrito"},
	}

	// Call
	t.callDomain()
	AssertNe(nil, t.c.req)

	AssertThat(
		getSortedKeys(t.c.req),
		ElementsAre(
			"Attribute.1.Name",
			"Attribute.1.Value",
			"Attribute.2.Name",
			"Attribute.2.Replace",
			"Attribute.2.Value",
			"Attribute.3.Name",
			"Attribute.3.Value",
			"DomainName",
			"ItemName",
		),
	)

	ExpectEq("foo", t.c.req["Attribute.1.Name"])
	ExpectEq("bar", t.c.req["Attribute.2.Name"])
	ExpectEq("baz", t.c.req["Attribute.3.Name"])

	ExpectEq("", t.c.req["Attribute.1.Value"])
	ExpectEq("taco", t.c.req["Attribute.2.Value"])
	ExpectEq("burrito", t.c.req["Attribute.3.Value"])

	ExpectEq("true", t.c.req["Attribute.2.Replace"])

	ExpectEq("some_item", t.c.req["ItemName"])
	ExpectEq(t.name, t.c.req["DomainName"])
}

func (t *PutTest) NoPreconditions() {
	// Call
	t.callDomain()
	AssertNe(nil, t.c.req)

	ExpectThat(getSortedKeys(t.c.req), Not(Contains(HasSubstr("Expected"))))
}

func (t *PutTest) SomePreconditions() {
	ExpectEq("TODO", "")
}

func (t *PutTest) ConnReturnsError() {
	ExpectEq("TODO", "")
}

func (t *PutTest) ConnSaysOkay() {
	ExpectEq("TODO", "")
}

////////////////////////////////////////////////////////////////////////
// BatchPutAttributes
////////////////////////////////////////////////////////////////////////

type BatchPutTest struct {
	domainTest
}

func init() { RegisterTestSuite(&BatchPutTest{}) }

func (t *BatchPutTest) DoesFoo() {
	ExpectEq("TODO", "")
}