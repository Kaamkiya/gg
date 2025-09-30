package shape

import (
	"reflect"
	"testing"

	"github.com/Kaamkiya/gg/internal/app/tetris/color"
)

func TestShapeMoveDown(t *testing.T) {
	shape := CreateNew(0, 0, NewRandomizer())
	movedDownShape := shape.MoveDown()

	if shape.color != movedDownShape.color {
		t.Fatal("Moving down should not modify color")
	}

	if shape.posY+1 != movedDownShape.posY {
		t.Fatal("Moving down should increase Y position by 1")
	}

	if shape.posX != movedDownShape.posX {
		t.Fatal("Moving down should not modify X position")
	}
}

func TestShapeRotateRight(t *testing.T) {
	shape := Shape{
		0,
		0,
		[][]bool{
			{true, true, true, true, true},
			{false, true, false, false, true},
		},
		color.None,
	}

	rotatedShape := shape.RotateRight()

	expectedGrid := [][]bool{
		{false, true},
		{true, true},
		{false, true},
		{false, true},
		{true, true},
	}

	if !reflect.DeepEqual(expectedGrid, rotatedShape.grid) {
		t.Fatal("Right rotation does not produce a correct grid")
	}
}

func TestShapeRotateLeft(t *testing.T) {
	shape := Shape{
		0,
		0,
		[][]bool{
			{true, true, true, true, true},
			{false, true, false, false, true},
			{true, true, false, true, true},
		},
		color.None,
	}

	rotatedShape := shape.RotateLeft()

	expectedGrid := [][]bool{
		{true, true, true},
		{true, false, true},
		{true, false, false},
		{true, true, true},
		{true, false, true},
	}

	if !reflect.DeepEqual(expectedGrid, rotatedShape.grid) {
		t.Fatal("Left rotation does not produce a correct grid")
	}
}

func testShapeOppositeRotationsCancelEachOther(t *testing.T) {
	shape := Shape{
		0,
		0,
		[][]bool{
			{true, true, true, true, true},
			{true, true, false, false, true},
			{false, true, false, true, false},
		},
		color.None,
	}

	rotatedShape := shape.RotateLeft().RotateRight()

	if !reflect.DeepEqual(shape.grid, rotatedShape.grid) {
		t.Fatal("Opposite rotations don't cancel each other")
	}
}
