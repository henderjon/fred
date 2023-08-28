package main

import "testing"

func Test_setPager(t *testing.T) {
	expected := &cache{
		pager: 5,
	}

	result := &cache{}

	doSetPager("5", result)

	if result.getPager() != expected.getPager() {
		t.Error("")
	}
}

func Test_setColumn(t *testing.T) {
	expected := &cache{
		column: 5,
	}

	result := &cache{}

	doSetColumn("5", result)

	if result.getColumn() != expected.getColumn() {
		t.Error("")
	}
}
