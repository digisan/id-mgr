package id

import (
	"fmt"
	"math"
	"sync"

	. "github.com/digisan/go-generics/v2"
)

type ID_TYPE int

const (
	ID_ROOT_HRCHY ID_TYPE = iota
	ID_ROOT_STDAL
	ID_HRCHY_ALLOC
	ID_HRCHY_UNALLOC
	ID_STDAL_ALLOC
	ID_STDAL_UNALLOC
	ID_UNKNOWN
)

func (t ID_TYPE) String() string {
	switch t {
	case ID_ROOT_HRCHY:
		return "ID_ROOT_HRCHY"
	case ID_ROOT_STDAL:
		return "ID_ROOT_STDAL"
	case ID_HRCHY_ALLOC:
		return "ID_HRCHY_ALLOC"
	case ID_HRCHY_UNALLOC:
		return "ID_HRCHY_UNALLOC"
	case ID_STDAL_ALLOC:
		return "ID_STDAL_ALLOC"
	case ID_STDAL_UNALLOC:
		return "ID_STDAL_UNALLOC"
	case ID_UNKNOWN:
		return "ID_UNKNOWN"
	default:
		return "ID_UNKNOWN"
	}
}

type ID uint64

const MaxID = ID(math.MaxUint64)

var (
	// key: super-id; value: descendants' amount. If key(super id) is 0, which means level0's amount
	mRecord = sync.Map{}
)

func Init(segsFromLow ...uint8) error {
	if err := init64bits(segsFromLow...); err != nil {
		return err
	} else {
		fmt.Printf("MASKS: %016x\n", _masks)
		fmt.Printf("SEGS:  %016x\n", _segs)
		fmt.Printf("CAP STDAL: %d\n", _cap_std)
		fmt.Printf("CAP HRCHY: %d\n", _cap_lvl)
	}
	mRecord.Store(ID(0), 0)
	mRecord.Store(MaxID, 0)
	return nil
}

func (id ID) Exists() bool {
	_, ok := mRecord.Load(id)
	return ok
}

func (id ID) Type() ID_TYPE {
	if id == 0 {
		return ID_ROOT_HRCHY
	}
	if id == MaxID {
		return ID_ROOT_STDAL
	}
	if len(_segs) <= 1 || _segs[0] == 0 || _segs[0] >= 64 {
		return ID_UNKNOWN
	}
	n := count1(_segs[0])
	switch lowBits(uint64(id), n) {
	case 0:
		return IF(id.Exists(), ID_STDAL_ALLOC, ID_STDAL_UNALLOC)
	default:
		return IF(id.Exists(), ID_HRCHY_ALLOC, ID_HRCHY_UNALLOC)
	}
}

func Level(id ID) int {
	for i, seg := range Reverse(_segs) {
		if uint64(id)&seg != 0 {
			return len(_segs) - i - 1
		}
	}
	return -1
}

func (id ID) Level() int {
	if id.Type() == ID_HRCHY_ALLOC {
		return Level(id)
	}
	return -1
}

func (id ID) ChildrenCount() int {
	if In(id.Type(), ID_HRCHY_ALLOC, ID_STDAL_ALLOC, ID_ROOT_HRCHY, ID_ROOT_STDAL) {
		if v, ok := mRecord.Load(id); ok {
			return v.(int)
		}
	}
	return -1
}

// func (id ID) Print() {
// 	parent, ok := id.Parent()
// 	fmt.Printf("%03x(%016b), lvl: %d, cap for self: %d, bit-idx for desc: %d, cap for desc: %d, ancestors: %x, parent: %v(%x)\n",
// 		uint(id),
// 		uint(id),
// 		id.Level(),
// 		id.Cap4SelfLvl(),
// 		id.BitIdx4Desc(),
// 		id.Cap4DescLvl(),
// 		id.Ancestors(),
// 		ok,
// 		parent,
// 	)
// }

func (id ID) BitIdx4TopLvl() int {
	lvl := id.Level()
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

func (id ID) BitIdx4NextDescLvl() int {
	if In(id.Type(), ID_HRCHY_ALLOC, ID_ROOT_HRCHY) {
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
	}
	return -1
}

func (id ID) Cap4TopLvl() uint64 {
	if id.Type() == ID_HRCHY_ALLOC {
		lvl := id.Level()
		if lvl == -1 {
			return 1
		}
		return _cap_lvl[lvl]
	}
	return 0
}

func (id ID) Cap4NextDescLvl() uint64 {
	if In(id.Type(), ID_HRCHY_ALLOC, ID_ROOT_HRCHY) {
		lvl := id.Level()
		if lvl == -1 {
			return _cap_lvl[0]
		}
		if lvl >= len(_cap_lvl)-1 {
			return 0
		}
		return _cap_lvl[lvl+1]
	}
	return 0
}

func (id ID) GenDescID() (ID, error) {
	if In(id.Type(), ID_HRCHY_ALLOC, ID_ROOT_HRCHY) {
		idx := id.BitIdx4NextDescLvl()
		cap := id.Cap4NextDescLvl()
		for i := uint64(1); i <= cap; i++ {
			desc := ID(i<<idx) | id
			if desc.Type() == ID_HRCHY_UNALLOC {
				if v, ok := mRecord.Load(id); ok {
					mRecord.Store(id, v.(int)+1)
				} else {
					return 0, fmt.Errorf("mRecord Load error on id(%x)", id)
				}
				if _, ok := mRecord.Load(desc); !ok {
					mRecord.Store(desc, 0)
				} else {
					return 0, fmt.Errorf("mRecord Load error on desc")
				}
				return desc, nil
			}
		}
		return 0, fmt.Errorf("no space for a new descendant id(%x)", id)
	}
	return 0, fmt.Errorf("only root or allocated hierarchy id can generate descendant, id(%x) cannot do", id)
}

func (id ID) Ancestors() (ids []ID) {
	if id.Type() == ID_HRCHY_ALLOC {
		for i := 0; i < id.Level(); i++ {
			sid := id & ID(_masks[i])
			if sid.Type() == ID_HRCHY_ALLOC {
				ids = append(ids, sid)
			}
		}
	}
	return
}

func (id ID) AncestorsWithSelf() (ids []ID) {
	if id.Type() == ID_HRCHY_ALLOC {
		for i := 0; i <= id.Level(); i++ {
			sid := id & ID(_masks[i])
			if sid.Type() == ID_HRCHY_ALLOC {
				ids = append(ids, sid)
			}
		}
	}
	return
}

// 0 is valid parent for level0's ID
func (id ID) Parent() (ID, bool) {
	if ancestors := id.Ancestors(); len(ancestors) > 0 {
		return ancestors[len(ancestors)-1], true
	}
	if id.Type() == ID_HRCHY_ALLOC && id.Level() == 0 {
		return 0, true
	}
	if id.Type() == ID_STDAL_ALLOC {
		return MaxID, true
	}
	return 0, false
}

func (id ID) iterDesc(out *[]ID) {
	if In(id.Type(), ID_HRCHY_ALLOC, ID_ROOT_HRCHY) {
		idx := id.BitIdx4NextDescLvl()
		cap := id.Cap4NextDescLvl()
		count := 0
		for i := uint64(1); i <= cap; i++ {
			desc_id := ID(i<<idx) | id
			if desc_id.Type() == ID_HRCHY_ALLOC {
				*out = append(*out, desc_id)
				count++
			}
			desc_id.iterDesc(out)
			if n := id.ChildrenCount(); uint64(count) == uint64(n) || n <= 0 {
				break
			}
		}
	}
}

func (id ID) Descendants(nextGenerations int) []ID {
	if In(id.Type(), ID_HRCHY_ALLOC, ID_ROOT_HRCHY) {
		rt := []ID{}
		id.iterDesc(&rt)
		nd := id.Level() + nextGenerations
		return Filter(rt, func(i int, e ID) bool {
			return e.Level() <= nd
		})
	}
	return nil
}

func (id ID) DescendantsWithSelf(nextGenerations int) []ID {
	if In(id.Type(), ID_HRCHY_ALLOC, ID_ROOT_HRCHY) {
		return append([]ID{id}, id.Descendants(nextGenerations)...)
	}
	return nil
}

func (id ID) PrintDescendants(nextGenerations int, inclSelf bool) {
	if In(id.Type(), ID_HRCHY_ALLOC, ID_ROOT_HRCHY) {
		descendants := []ID{}
		if inclSelf {
			descendants = id.DescendantsWithSelf(nextGenerations)
		} else {
			descendants = id.Descendants(nextGenerations)
		}
		for i, id := range descendants {
			switch i {
			case 0:
				fmt.Printf("[%06x", id)
			case len(descendants) - 1:
				fmt.Printf(" %06x]\n", id)
			default:
				fmt.Printf(" %06x", id)
			}
		}
	}
}

///////////////////////////////////////////////////////////////////////

func DeleteID(id ID, inclDesc bool) (rt []ID, err error) {
	if In(id.Type(), ID_ROOT_HRCHY, ID_ROOT_STDAL) {
		return nil, fmt.Errorf("root id cannot be deleted")
	}
	if !inclDesc {
		if id.ChildrenCount() > 0 {
			return nil, fmt.Errorf("id(%x) has children, cannot be deleted", id)
		}
		if sid, ok := id.Parent(); ok {
			n, ok := mRecord.Load(sid)
			if !ok {
				return nil, fmt.Errorf("id's parent(%x) record error", sid)
			}
			mRecord.Store(sid, n.(int)-1)
			mRecord.Delete(id)
			rt = []ID{id}
			// fmt.Printf("(0x%x) is deleted\n", id)
		}
	} else {
		// fmt.Println(id.DescendantsWithSelf(100))
		for _, desc := range Reverse(id.DescendantsWithSelf(100)) {
			if _, err := DeleteID(desc, false); err != nil {
				return nil, err
			}
			rt = append(rt, desc)
		}
	}
	return rt, nil
}

///////////////////////////////////////////////////////////////////////

func BitIdx4Stdal() int {
	if len(_segs) <= 1 || _segs[0] == 0 || _segs[0] >= 64 {
		return -1
	}
	return int(count1(_segs[0]))
}

func Cap4Stdal() uint64 {
	if len(_segs) <= 1 || _segs[0] == 0 || _segs[0] >= 64 {
		return 0
	}
	return F16 >> count1(_segs[0])
}

func GenStdalID() (ID, error) {
	idx := BitIdx4Stdal()
	cap := Cap4Stdal()
	for i := uint64(1); i <= cap; i++ {
		id := ID(i << idx)
		if id.Type() == ID_STDAL_UNALLOC {
			if v, ok := mRecord.Load(MaxID); ok {
				mRecord.Store(MaxID, v.(int)+1)
			} else {
				return 0, fmt.Errorf("mRecord Load error on MaxID")
			}
			if _, ok := mRecord.Load(id); !ok {
				mRecord.Store(id, 0)
			} else {
				return 0, fmt.Errorf("mRecord Load error on id")
			}
			return id, nil
		}
	}
	return 0, fmt.Errorf("no space for a new independent id")
}

///////////////////////////////////////////////////////////////////////

func HierarchyIDs() []ID {
	return ID(0).Descendants(len(_segs))
}

func StandaloneIDs() (rt []ID) {
	n := count1(_segs[0])
	nStdal := MaxID.ChildrenCount()
	capStdal := Cap4Stdal()
	for i := uint64(1); i <= capStdal; i++ {
		if id := ID(i << n); id.Type() == ID_STDAL_ALLOC {
			rt = append(rt, id)
		}
		if len(rt) == nStdal {
			break
		}
	}
	return
}

func WholeIDs() []ID {
	hIDs := HierarchyIDs()
	sIDs := StandaloneIDs()
	return append(hIDs, sIDs...)
}

///////////////////////////////////////////////////////////////////////

func SetID(id ID) (ID, error) {
	if In(id.Type(), ID_HRCHY_UNALLOC, ID_STDAL_UNALLOC) {
		var parent ID
		if lvl := Level(id); lvl == 0 {
			parent = 0
		} else if id > 0 && ID(_segs[0])&id == 0 {
			parent = MaxID
		} else {
			mask := _masks[lvl-1]
			parent = ID(mask & uint64(id))
		}

		// fmt.Printf("id(%x) - parent type: %s\n", id, parent.Type())

		switch parent.Type() {
		case ID_ROOT_HRCHY, ID_HRCHY_ALLOC:
			if NotIn(id, parent.Descendants(100)...) {
				if v, ok := mRecord.Load(parent); ok {
					mRecord.Store(id, 0)
					mRecord.Store(parent, v.(int)+1)
				} else {
					return 0, fmt.Errorf("id(0x%x)'s parent(0x%x) load error", id, parent)
				}
				return id, nil
			} else {
				return 0, fmt.Errorf("id(0x%x) already exists", id)
			}
		case ID_ROOT_STDAL:
			if NotIn(id, StandaloneIDs()...) {
				if v, ok := mRecord.Load(MaxID); ok {
					mRecord.Store(id, 0)
					mRecord.Store(MaxID, v.(int)+1)
				} else {
					return 0, fmt.Errorf("id(0x%x)'s parent(0x%x) load error", id, parent)
				}
				return id, nil
			} else {
				return 0, fmt.Errorf("id(0x%x) already exists", id)
			}
		default:
			return 0, fmt.Errorf("parent(0x%x) doesn't exist", parent)
		}
	}
	return 0, fmt.Errorf("id(0x%x) cannot be set", id)
}

func (id *ID) leftShift(nSeg int) error {
	if len(_masks) <= 1 || nSeg == 0 || nSeg >= len(_masks) {
		return fmt.Errorf("nSeg(%d) error or segs error", nSeg)
	}
	nShiftBit := count1(_masks[nSeg-1])
	*id <<= nShiftBit
	return nil
}

///////////////////////////////////////////////////////////////////////

func HookIDTree(id ID, tree ...ID) error {
	panic("not implemented")
}
