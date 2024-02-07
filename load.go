package idmgr

import (
	"fmt"
	"strconv"
	"strings"

	. "github.com/digisan/go-generics/v2"
	fd "github.com/digisan/gotk/file-dir"
	. "github.com/digisan/id-mgr/id"
)

// FILL 1._segs, 2._masks，3.mAlias 4.mRecord & REDO BuildHierarchy
func IngestTree(fPath string) error {
	var (
		i   int
		err error
	)

	// *** clearing all ID which includes clearing their aliases
	if err := ClrAllID(); err != nil {
		return err
	}

	_, e := fd.FileLineScan(fPath, func(line string) (bool, string) {
		i++

		if i == 1 {
			err = Init64bitsFromStr(line)
			return false, ""
		}

		// id | aliases
		if strings.Count(line, "|") != 1 {
			err = fmt.Errorf("ingested fail: incorrect ID format @%v", line)
			return false, ""
		}

		ln := strings.TrimSpace(line)
		id_alias := strings.Split(ln, "|")

		// *** ID *** //
		id, e := strconv.ParseUint(id_alias[0], 16, 64) // id in dump file is hex
		if e != nil {
			err = fmt.Errorf("ingested fail: id parsed error @%w", e)
			return false, ""
		}

		// *** Alias *** //
		_, err = SetID(ID(id), TypesAsAnyToAnys(strings.Split(id_alias[1], "^"))...)

		return true, ""
	}, "")

	if err != nil {
		return err
	}
	return e
}
