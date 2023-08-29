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

	if result.undo1 != nil || result.undo2.getRev() != 1 {
		t.Errorf("stage: %d / %d", result.undo1.getRev(), result.undo2.getRev())
		// t.Error("stage: nil / one")
	}

	result.stageUndo(two.clone())

	if result.undo1.getRev() != 1 || result.undo2.getRev() != 2 {
		t.Errorf("stage: %d / %d", result.undo1.getRev(), result.undo2.getRev())
		// t.Error("stage: one / two")
	}

	result.stageUndo(three.clone())

	if result.undo1.getRev() != 2 || result.undo2.getRev() != 3 {
		t.Errorf("stage: %d / %d", result.undo1.getRev(), result.undo2.getRev())
		// t.Error("stage: two / three")
	}

	result.unstageUndo()

	if result.undo1.getRev() != 3 || result.undo2.getRev() != 2 {
		t.Errorf("unstage: %d / %d", result.undo1.getRev(), result.undo2.getRev())
		// t.Error("unstage three / two")
	}

}
