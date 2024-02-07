package id

import (
	"errors"
	"fmt"
	"math"
	"sync"

	. "github.com/digisan/go-generics/v2"
	lk "github.com/digisan/logkit"
)

type ID_TYPE int

const (
	ID_HRCHY_ROOT ID_TYPE = iota
	ID_HRCHY_ALLOC
	ID_HRCHY_UNALLOC
	ID_STDAL_ROOT
	ID_STDAL_ALLOC
	ID_STDAL_UNALLOC
	ID_UNKNOWN
)

func (t ID_TYPE) String() string {
	switch t {
	case ID_HRCHY_ROOT:
		return "ID_HRCHY_ROOT"
	case ID_STDAL_ROOT:
		return "ID_STDAL_ROOT"
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
	if err := ClrAllID(); err != nil {
		return err
	}

	if err := init64bits(segsFromLow...); err != nil {
		return err
	} else {
		fmt.Printf("MASKS: %016x\n", _masks)
		fmt.Printf("SEGS:  %016x\n", _segs)
		fmt.Printf("CAP STDAL: %d\n", _cap_std)
		fmt.Printf("CAP HRCHY: %d\n", _cap_lvl)
	}
	mRecord.Store(ID(0), 0)
	mAlias.Store(ID(0), []any{ID_HRCHY_ROOT.String()})
	mRecord.Store(MaxID, 0)
	mAlias.Store(MaxID, []any{ID_STDAL_ROOT.String()})
	return nil
}

func (id ID) Exists() bool {
	_, ok := mRecord.Load(id)
	return ok
}

func (id ID) Type() ID_TYPE {
	if id == 0 {
		return ID_HRCHY_ROOT
	}
	if id == MaxID {
		return ID_STDAL_ROOT
	}
	if len(_segs) <= 1 || _segs[0] == 0 || _segs[0] >= F16>>1 {
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
	if In(id.Type(), ID_HRCHY_ALLOC, ID_STDAL_ALLOC, ID_HRCHY_ROOT, ID_STDAL_ROOT) {
		if v, ok := mRecord.Load(id); ok {
			return v.(int)
		}
	}
	return -1
}

func (id ID) Part(level int) ID {
	s := id & ID(_segs[level])
	if level == 0 {
		return s
	}
	return s >> count1(_masks[level-1])
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
	if In(id.Type(), ID_HRCHY_ALLOC, ID_HRCHY_ROOT) {
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
	lk.FailOnErrWhen(len(_cap_lvl) == 0, "%v", errors.New("[_cap_lvl] is not initialized"))
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
	lk.FailOnErrWhen(len(_cap_lvl) == 0, "%v", errors.New("[_cap_lvl] is not initialized"))
	if In(id.Type(), ID_HRCHY_ALLOC, ID_HRCHY_ROOT) {
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
	if In(id.Type(), ID_HRCHY_ALLOC, ID_HRCHY_ROOT) {
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

func (id ID) AvailableDescID() (ID, error) {
	desc, err := id.GenDescID()
	if err != nil {
		return 0, err
	}
	removed, err := DeleteID(desc, true)
	if err != nil {
		return 0, err
	}
	if len(removed) != 1 || removed[0] != desc {
		return 0, fmt.Errorf("desc(%x) from id(%x) deleted error", desc, id)
	}
	return desc, nil
}

func (id ID) Ancestors(inclSelf bool) (ids []ID) {
	if id.Type() == ID_HRCHY_ALLOC {
		n := id.Level()
		if inclSelf {
			n += 1
		}
		for i := 0; i < n; i++ {
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
	if ancestors := id.Ancestors(false); len(ancestors) > 0 {
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
	if In(id.Type(), ID_HRCHY_ALLOC, ID_HRCHY_ROOT) {
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

func (id ID) Descendants(nextGenerations int, inclSelf bool) []ID {
	if In(id.Type(), ID_HRCHY_ALLOC, ID_HRCHY_ROOT) {
		rt := []ID{}
		id.iterDesc(&rt)
		nd := id.Level() + nextGenerations
		desc := Filter(rt, func(i int, e ID) bool {
			return e.Level() <= nd
		})
		if inclSelf {
			return append([]ID{id}, desc...)
		} else {
			return desc
		}
	}
	return nil
}

func (id ID) PrintDescendants(nextGenerations int, inclSelf bool) {
	if In(id.Type(), ID_HRCHY_ALLOC, ID_HRCHY_ROOT) {
		descendants := id.Descendants(nextGenerations, inclSelf)
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

// by DFS
func (id ID) DescendantsInfo(inclSelf bool) (ids []ID, counts []int, aliases [][]any) {
	m := TreeNodeCount()
	for _, desc := range id.Descendants(100, inclSelf) {
		ids = append(ids, desc)
		counts = append(counts, m[desc])
		aliases = append(aliases, desc.Alias())
	}
	return
}

///////////////////////////////////////////////////////////////////////

func BitIdx4Stdal() int {
	if len(_segs) <= 1 || _segs[0] == 0 || _segs[0] >= F16>>1 {
		return -1
	}
	return int(count1(_segs[0]))
}

func Cap4Stdal() uint64 {
	if len(_segs) <= 1 || _segs[0] == 0 || _segs[0] >= F16>>1 {
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
	return ID(0).Descendants(len(_segs), false)
}

func StandaloneIDs() (rt []ID) {
	nStdal := MaxID.ChildrenCount()
	if nStdal <= 0 {
		return
	}
	n := count1(_segs[0])
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

func TreeDescription() map[ID][]ID {
	rt := make(map[ID][]ID)
	descendants := ID(0).Descendants(100, true)
	for _, d1 := range descendants {
		if d1.Level() == 0 {
			rt[0] = append(rt[0], d1)
		}
		for _, d2 := range descendants {
			if d2.Level()-d1.Level() == 1 {
				if d2.Level() == 0 {
					continue
				}
				if ID(_masks[d2.Level()-1])&d2 == d1 {
					rt[d1] = append(rt[d1], d2)
				}
			}
		}
	}
	return rt
}

func TreeNodeCount() map[ID]int {
	rt := make(map[ID]int)
	for _, desc := range ID(0).Descendants(100, true) {
		if n, ok := mRecord.Load(desc); ok {
			rt[desc] = n.(int)
		}
	}
	return rt
}

// TreeNodeCount is needed...
// the 1st node is existing id, then the others are new generated id
// the 1st nCountOfDesc is 1, then others are from TreeNodeCount[existing]

func CopyBranch(oriNode, underNode ID) error {
	if NotIn(oriNode.Type(), ID_HRCHY_ALLOC) {
		return fmt.Errorf("oriNode(0x%x) is invalid ID, cannot do CopyBranch", oriNode)
	}
	if NotIn(underNode.Type(), ID_HRCHY_ALLOC, ID_HRCHY_ROOT) {
		return fmt.Errorf("dstNode(0x%x) is invalid ID, cannot do CopyBranch", underNode)
	}

	_, listDescCount, listDescAliases := oriNode.DescendantsInfo(true)
	idx4Desc := 0
	// fmt.Println(listDescCount)
	// fmt.Println(listDescAliases)

	dstRootNode, err := underNode.GenDescID() // for root oriNode
	if err != nil {
		return err
	}
	fmt.Printf("root node ==> %x\n", dstRootNode)

	// ** copy root alias ** //
	aliases := oriNode.Alias()
	if err := dstRootNode.AddAlias(aliases...); err != nil {
		return err
	}
	// ** copy root alias ** //

	return copyBranch(oriNode, dstRootNode, listDescCount, listDescAliases, &idx4Desc)
}

func copyBranch(oriNode, dstNode ID, listDescCount []int, listDescAliases [][]any, pIdx4Desc *int) error {

	idxCount := listDescCount[*pIdx4Desc]
	if idxCount == 0 {
		// fmt.Printf("leaf node ==> %x\n", dstNode) // leaf node
		return nil
	}

	for i := 0; i < idxCount; i++ {
		nid, err := dstNode.GenDescID()
		if err != nil {
			return err
		}
		// fmt.Printf("new node ==> %x\n", nid)

		// ** copy alias ** //
		if *pIdx4Desc+1 < len(listDescAliases) {
			aliases := listDescAliases[*pIdx4Desc+1]
			if err := nid.AddAlias(aliases...); err != nil {
				return err
			}
		}
		// ** copy alias ** //

		(*pIdx4Desc)++
		if err := copyBranch(nid, nid, listDescCount, listDescAliases, pIdx4Desc); err != nil {
			return err
		}
	}
	return nil
}

func Transplant(oriNode, underNode ID) error {
	if err := CopyBranch(oriNode, underNode); err != nil {
		return err
	}
	if _, err := DeleteID(oriNode, true); err != nil {
		return err
	}
	return cleanupAlias()
}

// if id exists, do nothing and no error
func SetID(id ID, aliases ...any) (ID, error) {
	// if In(id.Type(), ID_HRCHY_UNALLOC, ID_STDAL_UNALLOC) {
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
	case ID_HRCHY_ROOT, ID_HRCHY_ALLOC:
		if NotIn(id, parent.Descendants(100, false)...) {
			if v, ok := mRecord.Load(parent); ok {
				mRecord.Store(id, 0)
				mRecord.Store(parent, v.(int)+1)
			} else {
				return 0, fmt.Errorf("id(0x%x)'s parent(0x%x) load error", id, parent)
			}

			// ** add the id's aliases ** //
			if err := id.AddAlias(aliases...); err != nil {
				return 0, err
			}

			return id, nil
		}
	case ID_STDAL_ROOT:
		if NotIn(id, StandaloneIDs()...) {
			if v, ok := mRecord.Load(MaxID); ok {
				mRecord.Store(id, 0)
				mRecord.Store(MaxID, v.(int)+1)
			} else {
				return 0, fmt.Errorf("id(0x%x)'s parent(0x%x) load error", id, parent)
			}

			// ** add the id's aliases ** //
			if err := id.AddAlias(aliases...); err != nil {
				return 0, err
			}

			return id, nil
		}
	default:
		return 0, fmt.Errorf("parent(0x%x) doesn't exist", parent)
	}
	// }
	return id, nil // fmt.Errorf("id(0x%x) already exists, SetID abort", id)
}

func DeleteID(id ID, inclDesc bool) (rt []ID, err error) {
	// if In(id.Type(), ID_HRCHY_ROOT, ID_STDAL_ROOT) {
	// 	return nil, fmt.Errorf("root id cannot be deleted")
	// }
	if !inclDesc {
		if id.ChildrenCount() > 0 {
			return nil, fmt.Errorf("id(%x) has children, cannot be deleted", id)
		}
		if sid, ok := id.Parent(); ok {
			n, ok := mRecord.Load(sid)
			if !ok {
				return nil, fmt.Errorf("id's parent(%x) record error", sid)
			}

			// ** before deleting id, its alias are removed in advance ** //
			if err := id.RmAlias(); err != nil {
				return nil, err
			}
			// ** //

			mRecord.Store(sid, n.(int)-1)
			mRecord.Delete(id)
			rt = []ID{id}
			// fmt.Printf("(0x%x) is deleted\n", id)
		}
	} else {
		// fmt.Println(id.Descendants(100, true))
		for _, desc := range Reverse(id.Descendants(100, true)) {
			if _, err := DeleteID(desc, false); err != nil {
				return nil, err
			}
			rt = append(rt, desc)
		}
	}
	return rt, nil
}

// only delete leaves id
func DeleteIDs(ids ...ID) error {
	for _, id := range ids {
		if _, err := DeleteID(id, false); err != nil {
			return err
		}
	}
	return nil
}

// all aliases are also deleted here
func ClrAllID() error {
	if len(_masks) == 0 || len(_segs) == 0 || len(_cap_lvl) == 0 || _cap_std == 0 {
		return nil
	}
	_, err := DeleteID(0, true)
	if err != nil {
		return err
	}
	return DeleteIDs(StandaloneIDs()...)
}

func IsValidID(id ID) bool {
	if id.Type() == ID_STDAL_ALLOC {
		return true
	}
	for _, id := range id.Ancestors(true) {
		if NotIn(id.Type(), ID_HRCHY_ROOT, ID_HRCHY_ALLOC) {
			return false
		}
	}
	return len(id.Ancestors(true)) > 0
}

func PrintRecord() {
	mRecord.Range(func(key, value any) bool {
		fmt.Println("mRecord:", key, "-->", value)
		return true
	})
}

// here the id could be temp id, i.e not existing. but still need to be shifted
// func (id *ID) leftShift(nSeg int, check bool) (ID, error) {
// 	if check && id.Type() != ID_HRCHY_ALLOC {
// 		return 0, fmt.Errorf("id(%x) is not existing", *id)
// 	}
// 	if len(_masks) <= 1 || nSeg >= len(_masks) {
// 		return 0, fmt.Errorf("nSeg(%d) error or segs error", nSeg)
// 	}
// 	if nSeg == 0 {
// 		return *id, nil
// 	}
// 	nShiftBit := count1(_masks[nSeg-1])
// 	*id <<= nShiftBit
// 	return *id, nil
// }

// func newLowSegID(id, low ID) ID {
// 	low = low & ID(_segs[0])
// 	r_seg0 := ^_segs[0]
// 	id_low0 := r_seg0 & uint64(id)
// 	return ID(id_low0 | uint64(low))
// }
