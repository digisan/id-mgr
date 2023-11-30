package idmgr

import (
	"fmt"
	"strconv"
	"strings"

	. "github.com/digisan/go-generics/v2"
	fd "github.com/digisan/gotk/file-dir"
)

// FILL 1._segs, 2._masks，3.mAlias 4.mRecord & REDO BuildHierarchy
func IngestTree(fpath string) error {

	// clear global variables
	_masks = []uint64{}
	_segs = []uint64{}

	var (
		idGroup []ID
		iLn     int
		err     error
	)
	fd.FileLineScan(fpath, func(line string) (bool, string) {
		iLn++

		// *** _masks, _segs *** //
		if ln := strings.Trim(line, "[]"); len(ln) == len(line)-2 {
			fmt.Println(ln)

			for _, seg_str := range strings.Split(ln, " ") {
				seg, e := strconv.ParseUint(seg_str, 16, 64) // seg in dump file is hex
				if e != nil {
					err = fmt.Errorf("ingested error: _segs parsed error @%w", e)
					return false, ""
				}
				_segs = append(_segs, seg)
			}
			_masks = genMasks(_segs)
			return true, ""
		}

		// id | children-count | aliases
		if strings.Count(line, "|") != 2 {
			err = fmt.Errorf("ingested error: ID line incorrect format @%v", line)
			return false, ""
		}
		ln := strings.TrimSpace(line)
		id_cnt_alias := strings.Split(ln, "|")

		// fmt.Printf("--> %02d %v\n", iLn, id_cnt_alias)

		// *** ID *** //
		id, e := strconv.ParseUint(id_cnt_alias[0], 16, 64) // id in dump file is hex
		if e != nil {
			err = fmt.Errorf("ingested error: id parsed error @%w", e)
			return false, ""
		}
		idGroup = append(idGroup, ID(id))

		// *** mRecord (missing mRecord[0]) *** //
		mRecord[ID(id)], _ = AnyTryToType[int](id_cnt_alias[1]) // count is dec

		// *** mAlias *** //
		mAlias[ID(id)] = TypesAsAnyToAnys(strings.Split(id_cnt_alias[2], "^"))

		return true, ""

	}, "")

	if len(_masks) == 0 || len(_segs) == 0 {
		return fmt.Errorf("_masks or _segs ingested error")
	}

	// *** mRecord[0], mRecord[MaxID] *** //
	for _, id := range idGroup {
		if pid, ok := id.Parent(); ok && pid == 0 && !id.IsStandalone() {
			mRecord[0]++
			continue
		}
		if id.IsStandalone() {
			mRecord[MaxID]++
			continue
		}
	}

	return err
}
