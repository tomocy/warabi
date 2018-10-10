package object

import "fmt"

type Kind int

const (
	Unknown Kind = iota
	Integer
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
