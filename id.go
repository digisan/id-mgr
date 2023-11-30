package idmgr

import (
	"fmt"
	"math"

	. "github.com/digisan/go-generics/v2"
)

type ID uint64

const MaxID = ID(math.MaxUint64)

var (
	// key: super-id; value: descendants' amount. If key(super id) is 0, which means level0's amount
	mRecord = make(map[ID]int)
	// alias: key id, value: aliases
	mAlias = make(map[ID][]any)
)

func (id ID) Exists() bool {
	_, ok := mRecord[id]
	return ok
}

func (id ID) Level() int {
	for i, seg := range Reverse(_segs) {
		// fmt.Printf("---> %d %016x\n", i, seg)
		if uint64(id)&seg != 0 {
			return len(_segs) - i - 1
		}
	}
	return -1
}

func (id ID) Print() {
	parent, ok := id.Parent()
	fmt.Printf("%03x(%016b), lvl: %d, cap for self: %d, bit-idx for desc: %d, cap for desc: %d, ancestors: %x, parent: %v(%x)\n",
		uint(id),
		uint(id),
		id.Level(),
		id.Cap4SelfLvl(),
		id.BitIdx4Desc(),
		id.Cap4DescLvl(),
		id.Ancestors(),
		ok,
		parent,
	)
}

// func (id ID) selfStartBitIdx() int {
// 	lvl := id.Level()
// 	if lvl < 0 {
// 		return -1
// 	}
// 	if lvl == 0 {
// 		return 0
// 	}
// 	mask := _masks[lvl-1]
// 	for i := 0; i < 64; i++ {
// 		if mask == 0 {
// 			return i
// 		}
// 		mask = mask >> 1
// 	}
// 	return 0
// }

func (id ID) BitIdx4Desc() int {
	lvl := id.Level()
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
	return -1
}

func (id ID) Cap4DescLvl() int {
	lvl := id.Level()
	if lvl == -1 {
		return int(_cap_lvl[0])
	}
	if lvl >= len(_cap_lvl)-1 {
		return 0
	}
	return int(_cap_lvl[lvl+1])
}

func (id ID) Cap4SelfLvl() int {
	lvl := id.Level()
	if lvl == -1 {
		return 1
	}
	return int(_cap_lvl[lvl])
}

// func (id ID) GenerateDescID() (ID, error) {
// 	idx := id.BitIdx4Desc()
// 	cap := id.Cap4DescLvl()
// 	for i := 1; i <= cap; i++ {
// 		desc := ID(i<<idx) | id
// 		if _, ok := mRecord[desc]; !ok {
//          // parent on mRecord
// 			mRecord[desc] = 0
// 			return desc, nil
// 		}
// 	}
// 	return 0, fmt.Errorf("no space for a new descendant id")
// }

func (id ID) Ancestors() (ids []ID) {
	if id.IsStandaloneFmt() {
		return nil
	}
	for i := 0; i < id.Level(); i++ {
		ids = append(ids, id&ID(_masks[i]))
	}
	return
}

// 0 is valid parent for level0's ID
func (id ID) Parent() (ID, bool) {
	if id.IsStandaloneFmt() {
		return MaxID, true
	}
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
	nd := id.Level() + nextGenerations
	return Filter(rt, func(i int, e ID) bool {
		return e.Level() <= nd
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
	if id == MaxID {
		return []any{"standalone"}
	}
	return mAlias[id]
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

func (id ID) IsStandalone() bool {
	_, ok := mRecord[id]
	return id.IsStandaloneFmt() && ok
}

func (id ID) IsStandaloneFmt() bool {
	n := count1(_segs[0])
	return lowBits(uint64(id), n) == 0
}

func HierarchyIDs() []ID {
	return ID(0).Descendants(len(_segs))
}

func StandaloneIDs() []ID {
	n := count1(_segs[0])
	nStandalone := mRecord[MaxID]
	rt := []ID{}
	for i := 1; true; i++ {
		id := ID(i << n)
		if _, ok := mRecord[id]; ok {
			rt = append(rt, id)
		}
		if len(rt) == nStandalone {
			break
		}
	}
	return rt
}

func WholeIDs() []ID {
	hIDs := HierarchyIDs()
	sIDs := StandaloneIDs()
	return append(hIDs, sIDs...)
}
