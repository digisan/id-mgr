package idmgr

import (
	"fmt"

	. "github.com/digisan/go-generics/v2"
)

const N uint8 = 64
const F16 uint64 = 0xFFFFFFFFFFFFFFFF

var (
	MASKS = []uint64{}
	SEGS  = []uint64{}
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
		MASKS = append(MASKS, F16>>uint64(shift))
	}
	// fmt.Println(MASKS)

	for i, mask := range MASKS {
		var seg uint64
		if i == 0 {
			seg = mask
		} else {
			seg = mask ^ MASKS[i-1]
		}
		// fmt.Printf("seg  %d: %016x\n\n", i, seg)
		SEGS = append(SEGS, seg)
	}
	// fmt.Println(SEGS)

	return nil
}

///////////////////////////////////////////////////////////////////////////////

type ID uint64

// level 6 | level 5   | level 4   | level 3   | level 2        | level 1        | level 0
// ****    | **** **** | **** **** | **** **** | **** **** **** | **** **** **** | **** **** ****

var (

	// key: super-id; value: descendants' amount. If key(super id) is 0, which means level0's amount
	mRecord = make(map[ID]uint)

	// key: id deleted
	mDeleted = make(map[ID]struct{})

	// alias: key id, value: aliases
	mAlias = make(map[ID][]string)
)

func (id ID) level() int {
	for i, s := range Reverse(SEGS) {
		// fmt.Printf("---> %d %016x\n", i, s)
		if uint64(id)&s != 0 {
			return len(SEGS) - i - 1
		}
	}
	return -1
}

func (id ID) availableSegBitIdx() int {
	lvl := id.level()
	if lvl == -1 {
		return 0
	}
	const bitChecker = 0x0F
	s := MASKS[lvl]
	for i := 0; i < 64; i++ {
		if bitChecker&s == 0 {
			return i
		}
		s = s >> 1
	}
	return 0
}

func (id ID) Ancestors() (ids []ID) {
	for i := 0; i < id.level(); i++ {
		ids = append(ids, id&ID(MASKS[i]))
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
		*out = append(*out, descID)
		descID.iterDesc(out)
	}
}

func (id ID) AncestorsWithSelf() (ids []ID) {
	return append(id.Ancestors(), id)
}

func (id ID) Alias() []string {
	return mAlias[id]
}

func aliasOccupied(alias string, byIDs ...ID) (bool, ID) {
	if len(byIDs) == 0 {
		byIDs = ID(0).Descendants(10)
	}
	for _, desc := range byIDs {
		if In(alias, desc.Alias()...) {
			return true, desc
		}
	}
	return false, 0
}

func SearchIDByAlias(alias string, fromIDs ...ID) ID {
	if len(fromIDs) == 0 {
		fromIDs = ID(0).Descendants(10)
	}
	for _, id := range fromIDs {
		if In(alias, id.Alias()...) {
			return id
		}
	}
	return 0
}

func (id ID) AddAlias(aliases ...string) ([]string, error) {

	// check alias conflict
	byIDs := ID(0).Descendants(10)
	for _, alias := range aliases {
		if used, byId := aliasOccupied(alias, byIDs...); used {
			return id.Alias(), fmt.Errorf("'%v' is already used by [%d], [%d] cannot use it", alias, byId, id)
		}
	}

	mAlias[id] = append(mAlias[id], aliases...)
	mAlias[id] = Settify(mAlias[id]...)
	return id.Alias(), nil
}

func (id ID) RmAlias(aliases ...string) []string {
	mAlias[id] = Filter(id.Alias(), func(i int, e string) bool {
		return NotIn(e, aliases...)
	})
	return id.Alias()
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

func maxDescCap(lvl int) uint {
	if lvl < len(SEGS) {
		return uint(trimLowZeroBin(SEGS[lvl]))
	}
	return 0
}

func makeID(sid ID, idx uint) ID {
	abi := sid.availableSegBitIdx()
	idx = idx << ID(abi)
	return sid | ID(idx)
}

func checkSuperID(sid ID) error {
	if sid == 0 {
		return nil
	}
	for _, anc := range sid.AncestorsWithSelf() {
		// fmt.Println("ancestor:", anc)
		if _, ok := mRecord[anc]; !ok {
			return fmt.Errorf("class value@%x(HEX) doesn't exist (level %d)", anc, anc.level())
		}
	}
	return nil
}

// if sid is 0, generate level 0 class
func GenID(sid ID) (ID, error) {
	if err := checkSuperID(sid); err != nil {
		return 0, err
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
		if nUsed == maxDescCap(lvl) {
			return 0, fmt.Errorf("level %d has no space to store", lvl)
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
	return nil
}

func MakeHierarchy(super, self string) {

}
