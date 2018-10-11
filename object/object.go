package object

import "fmt"

type Kind int

const (
	Unknown Kind = iota
	Integer
	String
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
