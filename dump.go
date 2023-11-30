package idmgr

import (
	"fmt"
	"os"
	"strings"

	. "github.com/digisan/go-generics/v2"
)

func GenIDTree(print bool) string {

	// fmt.Println(mRecord)
	// fmt.Println(WholeIDs())

	lines := []string{}
	for i, id := range HierarchyIDs() {
		lvl := id.Level()
		indent := strings.Repeat("\t", lvl)
		aliasesStr := ""
		if aliases, ok := AnysTryToTypes[string](id.Alias()); ok {
			aliasesStr = strings.Join(aliases, "^")
		}
		lines = append(lines, fmt.Sprintf("%s%x|%d|%v", indent, id, mRecord[id], aliasesStr)) // generated string use hexadecimal id

		if print {
			fmt.Printf("%03d: %s%x|%d|%v\n", i+1, indent, id, mRecord[id], aliasesStr) // print with line number, use hexadecimal 0xid
		}
	}

	offset := len(lines)
	for i, id := range StandaloneIDs() {
		aliasesStr := ""
		if aliases, ok := AnysTryToTypes[string](id.Alias()); ok {
			aliasesStr = strings.Join(aliases, "^")
		}
		lines = append(lines, fmt.Sprintf("%s%x|%d|%v", "\t", id, mRecord[id], aliasesStr))
		if print {
			fmt.Printf("%03d: %s%x|%d|%v\n", i+1+offset, "\t", id, mRecord[id], aliasesStr) // print with line number, use hexadecimal 0xid
		}
	}

	rt := strings.Join(lines, "\n")
	return rt
}

func DumpIDTree(fpath string) error {
	out := fmt.Sprintf("%016x\n", _segs) + GenIDTree(false)
	return os.WriteFile(fpath, []byte(out), os.ModePerm)
}
