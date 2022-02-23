package block

import "fmt"

// DoubleFlowerType represents a type of double flower.
type DoubleFlowerType struct {
	doubleFlower
}

type doubleFlower uint8

// Sunflower is a sunflower plant.
func Sunflower() DoubleFlowerType {
	return DoubleFlowerType{doubleFlower(0)}
}

// Lilac is a lilac plant.
func Lilac() DoubleFlowerType {
	return DoubleFlowerType{doubleFlower(1)}
}

// RoseBush is a rose bush plant.
func RoseBush() DoubleFlowerType {
	return DoubleFlowerType{doubleFlower(4)}
}

// Peony is a peony plant.
func Peony() DoubleFlowerType {
	return DoubleFlowerType{doubleFlower(5)}
}

// Uint8 returns the double plant as a uint8.
func (d doubleFlower) Uint8() uint8 {
	return uint8(d)
}

// Name ...
func (d doubleFlower) Name() string {
	switch d {
	case 0:
		return "Sunflower"
	case 1:
		return "Lilac"
	case 4:
		return "Rose Bush"
	case 5:
		return "Peony"
	}
	panic("unknown double plant type")
}

// FromString ...
func (d doubleFlower) FromString(s string) (interface{}, error) {
	switch s {
	case "sunflower":
		return DoubleFlowerType{doubleFlower(0)}, nil
	case "syringa":
		return DoubleFlowerType{doubleFlower(1)}, nil
	case "rose":
		return DoubleFlowerType{doubleFlower(4)}, nil
	case "paeonia":
		return DoubleFlowerType{doubleFlower(5)}, nil
	}
	return nil, fmt.Errorf("unexpected double plant type '%v', expecting one of 'sunflower', 'syringa', 'rose', or 'paeonia'", s)
}

// String ...
func (d doubleFlower) String() string {
	switch d {
	case 0:
		return "sunflower"
	case 1:
		return "syringa"
	case 4:
		return "rose"
	case 5:
		return "paeonia"
	}
	panic("unknown double plant type")
}

// DoubleFlowerTypes ...
func DoubleFlowerTypes() []DoubleFlowerType {
	return []DoubleFlowerType{Sunflower(), Lilac(), RoseBush(), Peony()}
}
