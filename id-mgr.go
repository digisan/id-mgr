package idmgr

import (
	"fmt"
	"os"
	"strings"

	. "github.com/digisan/go-generics/v2"
	"github.com/digisan/gotk/strs"
)

const N uint8 = 64
const F16 uint64 = 0xFFFFFFFFFFFFFFFF

var (
	masks = []uint64{}
	segs  = []uint64{}
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
		masks = append(masks, F16>>uint64(shift))
	}
	// fmt.Println(masks)

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
	// fmt.Println(segs)

	return nil
}

///////////////////////////////////////////////////////////////////////////////

type ID uint64

// level 6 | level 5   | level 4   | level 3   | level 2        | level 1        | level 0
// ****    | **** **** | **** **** | **** **** | **** **** **** | **** **** **** | **** **** ****

var (

	// key: super-id; value: descendants' amount. If key(super id) is 0, which means level0's amount
	mRecord = make(map[ID]uint)

	// alias: key id, value: aliases
	mAlias = make(map[ID][]any)
)

func (id ID) Exists() bool {
	_, ok := mRecord[id]
	return ok
}

func (id ID) level() int {
	for i, s := range Reverse(segs) {
		// fmt.Printf("---> %d %016x\n", i, s)
		if uint64(id)&s != 0 {
			return len(segs) - i - 1
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

	mask := masks[lvl-1]
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

	mask := masks[lvl]
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
		ids = append(ids, id&ID(masks[i]))
	}
	return
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
	for i := uint(1); i <= n; i++ {
		descID := makeID(id, i)
		// fmt.Println(descID)
		if _, ok := mRecord[descID]; ok {
			*out = append(*out, descID)
		}
		descID.iterDesc(out)
	}
}

func (id ID) AncestorsWithSelf() (ids []ID) {
	return append(id.Ancestors(), id)
}

func (id ID) Alias() []any {
	return mAlias[id]
}

func WholeIDs() []ID {
	return ID(0).Descendants(len(segs))
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
	exclChars = []string{"(", ")", "^"}
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
			return fmt.Errorf("'%v' is already used by [%d]", alias, byId)
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
	if lvl < len(segs) {
		return int(trimLowZeroBin(segs[lvl]))
	}
	return 0
}

func makeID(sid ID, idx uint) ID {
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
			mRecord[sid] = nUsed + 1
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
		return fmt.Errorf("%x has descendants [%x], cannot be deleted, nothing to delete", id, descIDs)
	}
	delete(mRecord, id)
	return nil
}

// DelIDViaAlias incurs updated WholeIDs
func DelIDViaAlias(alias any) error {
	id := SearchIDByAlias(alias, WholeIDs()...)
	if len(fmt.Sprint(alias)) > 0 && id == 0 {
		return fmt.Errorf("alias [%s] cannot be found, nothing to delete", alias)
	}
	return DelID(id)
}

// BuildHierarchy incurs updated WholeIDs. building one super with multiple descendants (each descendant with single alias)
func BuildHierarchy(super any, selves ...any) error {

	fromIDs := WholeIDs()

	sid := SearchIDByAlias(super, fromIDs...)
	if sid == 0 && len(fmt.Sprint(super)) > 0 {
		return fmt.Errorf("super must be empty string as root, but [%v] is given", super)
	}

	for _, self := range selves {
		if err := CheckAlias([]any{self}, fromIDs...); err != nil {
			return fmt.Errorf("%w, build nothing for [%s]-[%s]", err, super, selves)
		}
		id, err := GenID(sid)
		if err != nil {
			return err
		}
		fromIDs = WholeIDs()
		if _, err := id.AddAliases([]any{self}, fromIDs...); err != nil {
			return err
		}
	}
	return nil
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

func PrintHierarchy() string {
	// fmt.Println(mRecord)
	// fmt.Println(WholeIDs())
	lines := []string{}
	for _, id := range WholeIDs() {
		lvl := id.level()
		indent := strings.Repeat("\t", lvl)
		aliasesStr := ""
		if aliases, ok := AnysTryToTypes[string](id.Alias()); ok {
			aliasesStr = strings.Join(aliases, "^")
		}
		lines = append(lines, fmt.Sprintf("%s%d(%v)", indent, id, aliasesStr))
		// fmt.Printf("%s%d(%v)\n", indent, id, aliasesStr)
	}
	rt := strings.Join(lines, "\n")
	fmt.Println(rt)
	return rt
}

func DumpHierarchy(fpath string) error {
	return os.WriteFile(fpath, []byte(PrintHierarchy()), os.ModePerm)
}

func IngestHierarchy(fpath string) error {
	return nil
}
