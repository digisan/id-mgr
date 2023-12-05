package idmgr

import (
	"fmt"
	"os"
	"strings"

	. "github.com/digisan/go-generics/v2"
	. "github.com/digisan/id-mgr/id"
)

func GenIDTree(print bool) string {

	lines := []string{}
	for i, id := range HierarchyIDs() {
		lvl := id.Level()
		indent := strings.Repeat("\t", lvl)
		aliasesStr := ""
		if aliases, ok := AnysTryToTypes[string](id.Alias()); ok {
			aliasesStr = strings.Join(aliases, "^")
		}
		lines = append(lines, fmt.Sprintf("%s%x|%v", indent, id, aliasesStr)) // generated string use hexadecimal id

		if print {
			fmt.Printf("%03d: %s%x|%v\n", i+1, indent, id, aliasesStr) // print with line number, use hexadecimal 0xid
		}
	}

	offset := len(lines)
	for i, id := range StandaloneIDs() {
		aliasesStr := ""
		if aliases, ok := AnysTryToTypes[string](id.Alias()); ok {
			aliasesStr = strings.Join(aliases, "^")
		}
		lines = append(lines, fmt.Sprintf("%s%x|%v", "\t", id, aliasesStr))
		if print {
			fmt.Printf("%03d: %s%x|%v\n", i+1+offset, "\t", id, aliasesStr) // print with line number, use hexadecimal 0xid
		}
	}

	rt := strings.Join(lines, "\n")
	return rt
}

func DumpIDTree(fpath string) error {
	out := PrintSegs(false) + "\n" + GenIDTree(false)
	return os.WriteFile(fpath, []byte(out), os.ModePerm)
}
