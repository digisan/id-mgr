package idmgr

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"

	. "github.com/digisan/go-generics/v2"
	fd "github.com/digisan/gotk/file-dir"
	"github.com/digisan/gotk/strs"
)

const N uint8 = 64
const F16 uint64 = 0xFFFFFFFFFFFFFFFF

var (
	_masks = []uint64{}
	_segs  = []uint64{}
)

// most left segment stores root class, then its descendants
func Init64bits(segments ...uint8) error {
	if Sum(segments...) != 64 {
		return fmt.Errorf("error: sum of segments must be 64")
	}

	sum := uint8(0)
	for _, n := range segments {
		sum += n
		shift := N - sum
		_masks = append(_masks, F16>>uint64(shift))
	}
	fmt.Printf("MASKS: %016x\n", _masks)

	_segs = genSegs(_masks)
	fmt.Printf("SEGS : %016x\n", _segs)

	// fmt.Println(genMasks(_segs))
	if !reflect.DeepEqual(genMasks(_segs), _masks) {
		return fmt.Errorf("error: _masks & _segs are not consistent")
	}
	return nil
}

func genSegs(masks []uint64) (segs []uint64) {
	for i, mask := range masks {
		var seg uint64
		if i == 0 {
			seg = mask
		} else {
			seg = mask ^ masks[i-1]
		}
		// fmt.Printf("seg  %d: %016x\n\n", i, seg)
		segs = append(segs, seg)
	}
	return
}

func genMasks(segs []uint64) (masks []uint64) {
	for i, seg := range segs {
		var mask uint64
		if i == 0 {
			mask = seg
		} else {
			mask = seg | masks[i-1]
		}
		// fmt.Printf("mask  %d: %016x\n\n", i, mask)
		masks = append(masks, mask)
	}
	return
}

func lowBits(n uint64, nLow uint) uint64 {
	shift := 64 - nLow
	checker := F16 >> uint64(shift)
	return n & checker
}

func count1(n uint64) uint {
	count := uint(0)
	for i := 0; i < 64; i++ {
		flag1 := uint64(0b01 << i)
		if flag1&n == flag1 {
			count++
		}
	}
	return count
}

func countF(n uint64) uint {
	count := uint(0)
	for i := 0; i < 64; i += 4 {
		flagF := uint64(0x0F << i)
		if flagF&n == flagF {
			count++
		}
	}
	return count
}

///////////////////////////////////////////////////////////////////////////////

type ID uint64

// level 6 | level 5   | level 4   | level 3   | level 2        | level 1        | level 0
// ****    | **** **** | **** **** | **** **** | **** **** **** | **** **** **** | **** **** ****

var (

	// key: super-id; value: descendants' amount. If key(super id) is 0, which means level0's amount
	mRecord = make(map[ID]int)

	// alias: key id, value: aliases
	mAlias = make(map[ID][]any)
)

func MaxID() ID {
	return math.MaxUint64
}

func (id ID) Exists() bool {
	_, ok := mRecord[id]
	return ok
}

func (id ID) level() int {
	for i, s := range Reverse(_segs) {
		// fmt.Printf("---> %d %016x\n", i, s)
		if uint64(id)&s != 0 {
			return len(_segs) - i - 1
		}
	}
	return -1
}

func (id ID) selfStartBitIdx() int {
	lvl := id.level()
	if lvl < 0 {
		return -1
	}
	if lvl == 0 {
		return 0
	}

	mask := _masks[lvl-1]
	for i := 0; i < 64; i++ {
		if mask == 0 {
			return i
		}
		mask = mask >> 1
	}
	return 0
}

func (id ID) descAvailableBitIdx() int {
	lvl := id.level()
	if lvl == -1 {
		return 0
	}

	mask := _masks[lvl]
	for i := 0; i < 64; i++ {
		if mask == 0 {
			return i
		}
		mask = mask >> 1
	}
	return 0
}

func (id ID) Ancestors() (ids []ID) {
	for i := 0; i < id.level(); i++ {
		ids = append(ids, id&ID(_masks[i]))
	}
	return
}

// 0 is valid parent for level0's ID
func (id ID) Parent() (ID, bool) {
	if ancestors := id.Ancestors(); len(ancestors) > 0 {
		return ancestors[len(ancestors)-1], true
	}
	if id != 0 {
		return 0, true
	}
	return 0, false
}

func (id ID) Descendants(nextGenerations int) []ID {
	if _, ok := mRecord[id]; !ok {
		return nil
	}
	rt := []ID{}
	id.iterDesc(&rt)
	nd := id.level() + nextGenerations
	return Filter(rt, func(i int, e ID) bool {
		return e.level() <= nd
	})
}

func (id ID) iterDesc(out *[]ID) {
	n := mRecord[id]
	for i := 1; i <= n; i++ {
		desc_id := makeID(id, i)
		if _, ok := mRecord[desc_id]; ok {
			*out = append(*out, desc_id)
		}
		desc_id.iterDesc(out)
	}
}

func (id ID) PrintDescendants(nextGenerations int) {
	descendants := id.Descendants(nextGenerations)
	for i, id := range descendants {
		switch i {
		case 0:
			fmt.Printf("[%x", id)
		case len(descendants) - 1:
			fmt.Printf("%x]\n", id)
		default:
			fmt.Printf(" %x ", id)
		}
	}
}

func (id ID) AncestorsWithSelf() (ids []ID) {
	return append(id.Ancestors(), id)
}

func (id ID) Alias() []any {
	return mAlias[id]
}

func AllHierarchyIDs() []ID {
	return ID(0).Descendants(len(_segs))
}

func AllStandaloneIDs() []ID {
	n := count1(_segs[0])
	nStandalone := mRecord[MaxID()]
	rt := []ID{}
	for i := 1; i <= nStandalone; i++ {
		id := ID(i << n)
		if _, ok := mRecord[id]; ok {
			rt = append(rt, id)
		}
	}
	return rt
}

func WholeIDs() []ID {
	return append(AllHierarchyIDs(), AllStandaloneIDs()...)
}

func aliasOccupied(alias any, byIDs ...ID) (bool, ID) {
	if len(byIDs) == 0 {
		byIDs = WholeIDs()
	}
	for _, desc := range byIDs {
		if In(alias, desc.Alias()...) {
			return true, desc
		}
	}
	return false, 0
}

func SearchIDByAlias(alias any, fromIDs ...ID) ID {
	if len(fromIDs) == 0 {
		fromIDs = WholeIDs()
	}
	for _, id := range fromIDs {
		if In(alias, id.Alias()...) {
			return id
		}
	}
	return 0
}

var (
	exclChars = []string{"^", "|", ":", "[", "]"}
)

func validateAlias(alias any) bool {
	return !strs.ContainsAny(fmt.Sprint(alias), exclChars...)
}

// check alias conflict
func CheckAlias(aliases []any, fromIDs ...ID) error {
	if len(fromIDs) == 0 {
		fromIDs = WholeIDs()
	}
	for _, alias := range aliases {
		if !validateAlias(alias) {
			return fmt.Errorf("'%v' contains invalid characters like %+v", alias, exclChars)
		}
		if used, byId := aliasOccupied(alias, fromIDs...); used {
			return fmt.Errorf("'%v' is already used by [%x]", alias, byId)
		}
	}
	return nil
}

func (id ID) AddAliases(aliases []any, validRangeIDs ...ID) ([]any, error) {
	if !id.Exists() {
		return nil, fmt.Errorf("error: %v doesn't exist, cannot do AddAlias", id)
	}

	// check alias conflict
	if err := CheckAlias(aliases, validRangeIDs...); err != nil {
		return id.Alias(), err
	}

	mAlias[id] = append(mAlias[id], aliases...)
	mAlias[id] = Settify(mAlias[id]...)
	return id.Alias(), nil
}

func (id ID) RmAliases(aliases ...any) ([]any, error) {
	if !id.Exists() {
		return nil, fmt.Errorf("error: %v doesn't exist, cannot do RmAlias", id)
	}
	mAlias[id] = Filter(id.Alias(), func(i int, e any) bool {
		return NotIn(e, aliases...)
	})
	return id.Alias(), nil
}

func trimLowZero(s uint64, bitStep uint8) uint64 {
	var bitChecker = F16 >> (64 - bitStep)
	for {
		if bitChecker&s != 0 {
			return s
		}
		s = s >> bitStep
	}
}

func trimLowZeroBin(s uint64) uint64 {
	return trimLowZero(s, 1)
}

func trimLowZeroOct(s uint64) uint64 {
	return trimLowZero(s, 3)
}

func trimLowZeroHex(s uint64) uint64 {
	return trimLowZero(s, 4)
}

func nextLvlDescCap(lvl int) int {
	if lvl < 0 {
		return -1
	}
	if lvl < len(_segs) {
		return int(trimLowZeroBin(_segs[lvl]))
	}
	return 0
}

func makeID(sid ID, idx int) ID {
	return ID(idx<<ID(sid.descAvailableBitIdx())) | sid
}

func IsValidID(id ID) bool {
	if len(mRecord) == 0 && id == 0 {
		return true
	}
	for _, id := range id.AncestorsWithSelf() {
		if !id.Exists() {
			return false
		}
	}
	return true
}

// if sid is 0, generate level 0 class
func GenID(sid ID) (ID, error) {
	if !IsValidID(sid) {
		return 0, fmt.Errorf("error: %x(HEX) is invalid ID, cannot be another's super ID", sid)
	}
	if nUsed, ok := mRecord[sid]; !ok || nUsed == 0 { // the first descendant class comes
		id := makeID(sid, 1)
		defer func() {
			mRecord[sid] = 1
			mRecord[id] = 0
		}()
		return id, nil
	} else {
		lvl := sid.level()
		if sid == 0 {
			lvl = 0
		}
		if int(nUsed) == nextLvlDescCap(lvl) {
			return 0, fmt.Errorf("level [%d] has no space to store [%d]", lvl, nUsed+1)
		}
		id := makeID(sid, nUsed+1)
		defer func() {
			mRecord[sid]++
			mRecord[id] = 0
		}()
		return id, nil
	}
}

func GenIDStandalone() (ID, error) {
	n := count1(_segs[0])
	sid := MaxID()
	if nUsed, ok := mRecord[sid]; !ok || nUsed == 0 {
		id := ID(0b01 << n)
		defer func() {
			mRecord[sid] = 1
			mRecord[id] = 0
		}()
		return id, nil
	} else {
		id := ID((nUsed + 1) << n)
		defer func() {
			mRecord[sid]++
			mRecord[id] = 0
		}()
		return id, nil
	}
}

func DelID(id ID) error {
	if _, ok := mRecord[id]; !ok {
		return nil
	}
	if descIDs := id.Descendants(1); len(descIDs) > 0 {
		return fmt.Errorf("%x(%v) has descendants [%x], cannot delete, abort", id, id.Alias(), descIDs)
	}
	delete(mRecord, id)

	// DO NOT modify parent mRecord!!!
	// if parent, ok := id.Parent(); ok {
	// 	mRecord[parent]--
	// }

	return nil
}

func DelIDs(ids ...ID) error {
	for _, id := range ids {
		if err := DelID(id); err != nil {
			return err
		}
	}
	return nil
}

// DelIDOnAlias incurs updated WholeIDs
func DelIDOnAlias(alias any) error {
	id := SearchIDByAlias(alias, WholeIDs()...)
	if len(fmt.Sprint(alias)) > 0 && id == 0 {
		return fmt.Errorf("alias [%s] cannot be found, nothing to delete", alias)
	}
	return DelID(id)
}

func DelIDsOnAlias(aliases ...any) error {
	for _, alias := range aliases {
		if err := DelIDOnAlias(alias); err != nil {
			return err
		}
	}
	return nil
}

// BuildHierarchy incurs updated WholeIDs. building one super with multiple descendants (each descendant with single alias!)
func BuildHierarchy(super any, descAliases ...any) ([]ID, error) {

	fromIDs := WholeIDs()

	sid := SearchIDByAlias(super, fromIDs...)
	if sid == 0 && len(fmt.Sprint(super)) > 0 {
		return nil, fmt.Errorf("super must be empty string as root, but [%v] is given", super)
	}

	rt := []ID{}
	for _, self := range descAliases {
		if err := CheckAlias([]any{self}, fromIDs...); err != nil {
			return nil, fmt.Errorf("%w, build nothing for [%s]-[%s]", err, super, descAliases)
		}
		id, err := GenID(sid)
		if err != nil {
			return nil, err
		}
		fromIDs = WholeIDs()
		if _, err := id.AddAliases([]any{self}, fromIDs...); err != nil {
			return nil, err
		}
		rt = append(rt, id)
	}
	return rt, nil
}

func BuildStandalone(aliases ...any) ([]ID, error) {
	fromIDs := WholeIDs()
	if err := CheckAlias(aliases, fromIDs...); err != nil {
		return nil, err
	}
	rt := []ID{}
	for _, alias := range aliases {
		id, err := GenIDStandalone()
		if err != nil {
			return nil, err
		}
		fromIDs = WholeIDs()
		if _, err := id.AddAliases([]any{alias}, fromIDs...); err != nil {
			return nil, err
		}
		rt = append(rt, id)
	}
	return rt, nil
}

func (id ID) IsStandalone() bool {
	n := count1(_segs[0])
	_, ok := mRecord[id]
	return lowBits(uint64(id), n) == 0 && ok
}

func AddAliases(self any, aliases ...any) error {
	id := SearchIDByAlias(self)
	_, err := id.AddAliases(aliases)
	return err
}

func GetAliases(self any) []any {
	id := SearchIDByAlias(self)
	return id.Alias()
}

func RmAliases(self any, aliases ...any) error {
	id := SearchIDByAlias(self)
	_, err := id.RmAliases(aliases...)
	return err
}

func ChangeAlias(old, new any) error {
	if err := AddAliases(old, new); err != nil {
		return err
	}
	if err := RmAliases(new, old); err != nil {
		return err
	}
	return nil
}

func GenHierarchy(print bool) string {
	// fmt.Println(mRecord)
	// fmt.Println(WholeIDs())
	lines := []string{}
	for i, id := range WholeIDs() {
		lvl := id.level()
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
	rt := strings.Join(lines, "\n")
	return rt
}

func DumpHierarchy(fpath string) error {
	out := fmt.Sprintf("%016x\n", _segs) + GenHierarchy(false)
	return os.WriteFile(fpath, []byte(out), os.ModePerm)
}

// FILL 1._segs, 2._masks，3.mAlias 4.mRecord & REDO BuildHierarchy
func IngestHierarchy(fpath string) error {

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

	// *** mRecord[0] *** //
	for _, id := range idGroup {
		if pid, ok := id.Parent(); ok && pid == 0 {
			mRecord[0]++
		}
	}

	return err
}
