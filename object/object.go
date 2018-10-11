package object

import "fmt"

type Kind int

const (
	Unknown Kind = iota
	Integer
	String
	Character
	FloatingPoint
	Boolean
)

type Object interface {
	Kind() Kind
	String() string
}

type IntegerLiteral struct {
	Value int
}

func (l IntegerLiteral) Kind() Kind {
	return Integer
}

func (l IntegerLiteral) String() string {
	return fmt.Sprintf("%d", l.Value)
}

type StringLiteral struct {
	Value string
}

func (l StringLiteral) Kind() Kind {
	return String
}

func (l StringLiteral) String() string {
	return l.Value
}

type CharacterLiteral struct {
	Value rune
}

func (l CharacterLiteral) Kind() Kind {
	return Character
}

func (l CharacterLiteral) String() string {
	return fmt.Sprintf("%d", l.Value)
}

type FloatingPointLiteral struct {
	Value float32
}

func (l FloatingPointLiteral) Kind() Kind {
	return FloatingPoint
}

func (l FloatingPointLiteral) String() string {
	return fmt.Sprintf("%e", l.Value)
}

var (
	True = &BooleanLiteral{
		value: true,
	}
	False = &BooleanLiteral{
		value: false,
	}
)

type BooleanLiteral struct {
	value bool
}

func (l BooleanLiteral) Kind() Kind {
	return Boolean
}

func (l BooleanLiteral) String() string {
	return fmt.Sprintf("%t", l.value)
}
