package main

import (
	"fmt"
	"os"
)

func debug(values ...any) {
	fmt.Fprintln(os.Stderr, values...)
}
