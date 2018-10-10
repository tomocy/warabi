package repl

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/tomocy/warabi/evaluator"
)

const packageStatement = "package main\n"

type REPLer interface {
	REPL()
}

func NewStandard(r io.Reader, w io.Writer) REPLer {
	return newStandard(r, w)
}

func NewWarabi(r io.Reader, w io.Writer) REPLer {
	return newWarabi(r, w)
}

type standard struct {
	*repler
}

func newStandard(r io.Reader, w io.Writer) *standard {
	return &standard{
		repler: new(r, w),
	}
}

func (repler standard) REPL() {
	go repler.waitSignal()
	repler.repl()
}

func (repler standard) repl() {
	scanner := bufio.NewScanner(repler.r)
	fileSet := token.NewFileSet()
	for scanner.Scan() {
		src := packageStatement + scanner.Text()
		file, _ := parser.ParseFile(fileSet, "example.go", src, parser.Mode(0))
		ast.Print(fileSet, file)
		repler.println()
	}
}

type warabi struct {
	*repler
}

func newWarabi(r io.Reader, w io.Writer) *warabi {
	return &warabi{
		repler: new(r, w),
	}
}

func (repler warabi) REPL() {
	go repler.waitSignal()
	repler.repl()
}

func (repler warabi) repl() {
	scanner := bufio.NewScanner(repler.r)
	for scanner.Scan() {
		objs := evaluator.Evaluate(scanner.Text())
		strs := make([]string, len(objs))
		for i, obj := range objs {
			strs[i] = obj.String()
		}
		repler.println(strings.Join(strs, ", "))
	}
}

type repler struct {
	r     io.Reader
	w     io.Writer
	sigCh chan os.Signal
}

func new(r io.Reader, w io.Writer) *repler {
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT)
	return &repler{
		r:     r,
		w:     w,
		sigCh: sigCh,
	}
}

func (repler repler) waitSignal() {
	for {
		select {
		case <-repler.sigCh:
			repler.println()
			repler.println("See you later")
			os.Exit(1)
		}
	}
}

func (repler repler) print(a ...interface{}) {
	fmt.Fprint(repler.w, a...)
}

func (repler repler) println(a ...interface{}) {
	fmt.Fprintln(repler.w, a...)
}

func (repler repler) printf(format string, a ...interface{}) {
	fmt.Fprintf(repler.w, format, a...)
}
