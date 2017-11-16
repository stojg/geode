package main

import (
	"math/rand"
	"os"

	"github.com/stojg/graphics/lib/core"
)

func main() {
	rand.Seed(19)
	l := NewLogger("gl.log")
	err := core.Main(l)
	if err != nil {
		l.Error(err)
		l.Close()
		os.Exit(1)
	}
	l.Close()
}
