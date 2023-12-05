package idmgr

import (
	"fmt"
	"strconv"
	"strings"

	. "github.com/digisan/go-generics/v2"
	fd "github.com/digisan/gotk/file-dir"
	. "github.com/digisan/id-mgr/id"
	lk "github.com/digisan/logkit"
)

// FILL 1._segs, 2._masks，3.mAlias 4.mRecord & REDO BuildHierarchy
func IngestTree(fpath string) error {
	var (
		i   int
		err error
	)

	lk.FailOnErr("%v", ClrAllID()) // clearing all ID which includes clearing their aliases

	fd.FileLineScan(fpath, func(line string) (bool, string) {
		i++

		if i == 1 {
			lk.FailOnErr("%v", Init64bitsFromStr(line))
			return true, ""
		}

		// id | aliases
		lk.FailOnErrWhen(strings.Count(line, "|") != 1, "%v", fmt.Errorf("ingested fail: incorrect ID format @%v", line))
		ln := strings.TrimSpace(line)
		id_alias := strings.Split(ln, "|")

		// *** ID *** //
		id, err := strconv.ParseUint(id_alias[0], 16, 64) // id in dump file is hex
		if err != nil {
			lk.FailOnErr("%v", fmt.Errorf("ingested fail: id parsed error @%w", err))
		}

		_, err = SetID(ID(id), TypesAsAnyToAnys(strings.Split(id_alias[1], "^"))...)
		lk.FailOnErr("%v", err)

		return true, ""
	}, "")

	return err
}
