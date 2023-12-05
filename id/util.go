package id

import (
	"fmt"
	"strings"

	. "github.com/digisan/go-generics/v2"
)

func DuplMark(v any) any {
	s := fmt.Sprint(v)
	io, ic := strings.LastIndex(s, "("), strings.LastIndex(s, ")")
	if io == -1 || ic == -1 {
		return fmt.Sprintf("%v(2)", s)
	}
	name, idxstr := s[:io], s[io+1:ic]
	if idx, ok := AnyTryToType[int](idxstr); ok && strings.HasSuffix(s, ")") {
		return fmt.Sprintf("%v(%d)", name, idx+1)
	}
	return s + "*"
}
