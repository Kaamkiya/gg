package main

import (
	"fmt"
	"strings"
)

func Test(s string) {
	fmt.Println(strings.TrimRight(s, " ") + "+")
}
