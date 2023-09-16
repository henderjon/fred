//go:build memory

package main

import (
	"testing"
)

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

func Test_undo(t *testing.T) {
	one := &memoryBuf{
		rev: 1,
	}

	two := &memoryBuf{
		rev: 2,
	}

	three := &memoryBuf{
		rev: 3,
	}

	result := &cache{}

	result.stageUndo(one.clone())

	if result.prevBuffer != nil || result.currBuffer.getRev() != 1 {
		t.Errorf("stage: %d / %d", result.prevBuffer.getRev(), result.currBuffer.getRev())
		// t.Error("stage: nil / one")
	}

	result.stageUndo(two.clone())

	if result.prevBuffer.getRev() != 1 || result.currBuffer.getRev() != 2 {
		t.Errorf("stage: %d / %d", result.prevBuffer.getRev(), result.currBuffer.getRev())
		// t.Error("stage: one / two")
	}

	result.stageUndo(three.clone())

	if result.prevBuffer.getRev() != 2 || result.currBuffer.getRev() != 3 {
		t.Errorf("stage: %d / %d", result.prevBuffer.getRev(), result.currBuffer.getRev())
		// t.Error("stage: two / three")
	}

	result.unstageUndo()

	if result.prevBuffer.getRev() != 3 || result.currBuffer.getRev() != 2 {
		t.Errorf("unstage: %d / %d", result.prevBuffer.getRev(), result.currBuffer.getRev())
		// t.Error("unstage three / two")
	}

}
