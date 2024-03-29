package id

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	. "github.com/digisan/go-generics"
)

// level 6 | level 5   | level 4   | level 3   | level 2        | level 1        | level 0
// ****    | **** **** | **** **** | **** **** | **** **** **** | **** **** **** | **** **** ****

const N uint8 = 64
const F16 uint64 = 0xFFFFFFFFFFFFFFFF

var (
	_masks   = []uint64{}
	_segs    = []uint64{}
	_cap_lvl = []uint64{}
	_cap_std = uint64(0)
)

func clr64bits() {
	_masks = []uint64{}
	_segs = []uint64{}
	_cap_lvl = []uint64{}
	_cap_std = uint64(0)
}

func Init64bitsFromStr(s string) error {
	clr64bits()
	if ls := strings.Trim(s, "[]"); len(ls) == len(s)-2 && len(ls) > 0 {
		for _, seg_str := range strings.Split(ls, " ") {
			seg, e := strconv.ParseUint(seg_str, 16, 64) // seg in dump file is hex
			if e != nil {
				return fmt.Errorf("ingested error: _segs parsed error @%w", e)
			}
			_segs = append(_segs, seg)
		}
		_masks = genMasks(_segs)
		if err := initCaps(); err != nil {
			return err
		}
		mRecord.Store(ID(0), 0)
		mAlias.Store(ID(0), []any{ID_HRCHY_ROOT.String()})
		mRecord.Store(MaxID, 0)
		mAlias.Store(MaxID, []any{ID_STDAL_ROOT.String()})
		return nil
	}
	return fmt.Errorf("invalid s(%s) for init64bits", s)
}

// most left segment stores root class, then its descendants
func init64bits(segsFromLow ...uint8) error {
	if Sum(segsFromLow...) != 64 {
		return fmt.Errorf("error: sum of segments must be 64")
	}
	clr64bits()
	sum := uint8(0)
	for _, n := range segsFromLow {
		sum += n
		shift := N - sum
		_masks = append(_masks, F16>>uint64(shift))
	}
	// fmt.Printf("MASKS: %016x\n", _masks)

	_segs = genSegs(_masks)
	// fmt.Printf("SEGS : %016x\n", _segs)

	// fmt.Println(genMasks(_segs))
	if !reflect.DeepEqual(genMasks(_segs), _masks) {
		return fmt.Errorf("error: _masks & _segs are not consistent")
	}

	return initCaps()
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

func capOfDescendant(lvl int) int {
	if lvl < 0 {
		return -1
	}
	if lvl < len(_segs) {
		return int(trimLowB0(_segs[lvl]))
	}
	return 0
}

func initCaps() error {
	if len(_segs) == 0 || len(_masks) == 0 {
		return fmt.Errorf("_segs or _masks is not initialized")
	}
	_cap_std = F16 >> uint64(count1(_segs[0]))
	for _, seg := range _segs {
		_cap_lvl = append(_cap_lvl, trimLowB0(seg))
	}
	return nil
}

func PrintSegs(print bool) string {
	s := fmt.Sprintf("%016x", _segs)
	if print {
		fmt.Println(s)
	}
	return s
}

func PrintMasks(print bool) string {
	s := fmt.Sprintf("%016x", _masks)
	if print {
		fmt.Println(s)
	}
	return s
}

func PrintCapLevel(print bool) string {
	s := fmt.Sprintf("%d", _cap_lvl)
	if print {
		fmt.Println(s)
	}
	return s
}

func PrintCapStdal(print bool) string {
	s := fmt.Sprintf("%d", _cap_std)
	if print {
		fmt.Println(s)
	}
	return s
}
