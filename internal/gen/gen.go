package gen

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"

	"github.com/mmcloughlin/avo/internal/inst"
)

type Interface interface {
	Generate([]inst.Instruction) ([]byte, error)
}

type Func func([]inst.Instruction) ([]byte, error)

func (f Func) Generate(is []inst.Instruction) ([]byte, error) {
	return f(is)
}

type Config struct {
	Name string
	Argv []string
}

func (c Config) GeneratedBy() string {
	if c.Argv == nil {
		return c.Name
	}
	return fmt.Sprintf("command: %s", strings.Join(c.Argv, " "))
}

func (c Config) GeneratedWarning() string {
	return fmt.Sprintf("Code generated by %s. DO NOT EDIT.", c.GeneratedBy())
}

type Builder func(Config) Interface

// GoFmt formats Go code produced from the given generator.
func GoFmt(i Interface) Interface {
	return Func(func(is []inst.Instruction) ([]byte, error) {
		b, err := i.Generate(is)
		if err != nil {
			return nil, err
		}
		return format.Source(b)
	})
}

type generator struct {
	buf bytes.Buffer
	err error
}

func (g *generator) Printf(format string, args ...interface{}) {
	if g.err != nil {
		return
	}
	if _, err := fmt.Fprintf(&g.buf, format, args...); err != nil {
		g.AddError(err)
	}
}

func (g *generator) AddError(err error) {
	if err != nil && g.err == nil {
		g.err = err
	}
}

func (g *generator) Result() ([]byte, error) {
	return g.buf.Bytes(), g.err
}
