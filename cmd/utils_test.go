package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

/*
	it("assigns strings", () => {
	  expect(assignmentString("i", 7)).toBe("\ni = 7");
	});
*/
func Test_assignmentString(t *testing.T) {
	const expected = "\ni = 7"
	if actual := assignmentString("i", 7); actual != expected {
		t.Errorf(`assignment string expected "\ni = 7" but got %s`, actual)
	}
}

/*
	it("maps contents", () => {
	  const contents = { a_cool_key: "a cool value", another_key: "another value" };
	  expect(map(contents, assignmentString)).toMatchInlineSnapshot(`
	    "
	    a_cool_key = \\"a cool value\\"
	    another_key = \\"another value\\""
	  `);
	});
*/
func Test_mapContents(t *testing.T) {
	input := jmap{"a_cool_key": "a cool value", "another_key": "another value"}
	const expected = `
a_cool_key = "a cool value"
another_key = "another value"`
	if actual := mapContents(input, assignmentString); actual != expected {
		t.Errorf("mapContents output did not match: %v\n", cmp.Diff(expected, actual))
	}
}

/*
	it("converts a block", () => {
	  const contents = { a_cool_key: "a cool value", another_key: "another value" };
	  expect(block("my_block", contents, assignmentString)).toMatchInlineSnapshot(`
	    "
	    my_block {
	    a_cool_key = \\"a cool value\\"
	    another_key = \\"another value\\"
	    }"
	  `);
	});
*/
func Test_block(t *testing.T) {
	input := jmap{"a_cool_key": "a cool value", "another_key": "another value"}
	const expected = `
my_block {
a_cool_key = "a cool value"
another_key = "another value"
}`

	if actual := block("my_block", input, assignmentString); actual != expected {
		t.Errorf("block output did not match: %v\n", cmp.Diff(expected, actual))
	}
}

/*
	it("converts a block list", () => {
	  const contents = [
	    { a_cool_key: "a cool value", another_key: "another value" },
	    { a_cool_key: "another cool value", another_key: "another ANOTHER value" },
	  ];
	  expect(blockList(contents, "my_block", assignmentString)).toMatchInlineSnapshot(`
	    "

	    my_block {
	    a_cool_key = \\"a cool value\\"
	    another_key = \\"another value\\"
	    }
	    my_block {
	    a_cool_key = \\"another cool value\\"
	    another_key = \\"another ANOTHER value\\"
	    }"
	  `);
	});
*/
func Test_blockList(t *testing.T) {
	input := jmaps{
		jmap{"a_cool_key": "a cool value", "another_key": "another value"},
		jmap{"a_cool_key": "another cool value", "another_key": "another ANOTHER value"},
	}
	const expected = `

my_block {
a_cool_key = "a cool value"
another_key = "another value"
}
my_block {
a_cool_key = "another cool value"
another_key = "another ANOTHER value"
}`
	if actual := blockList(input, "my_block", assignmentString); actual != expected {
		t.Errorf("blockList output did not match: %v\n", cmp.Diff(expected, actual))
	}
}
