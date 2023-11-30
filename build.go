package idmgr

import "fmt"

func makeID(sid ID, idx int) ID {
	return ID(idx<<sid.BitIdx4Desc()) | sid
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
func CreateHierarchyID(sid ID) (ID, error) {
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
		lvl := sid.Level()
		if sid == 0 {
			lvl = 0
		}
		if int(nUsed) == capOfDescendant(lvl) {
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

func CreateStandaloneID() (ID, error) {
	n := count1(_segs[0])
	sid := MaxID
	for i := 1; i < int(F16>>n); i++ {
		id := ID(i << n)
		if _, ok := mRecord[id]; !ok {
			mRecord[sid]++
			mRecord[id] = 0
			return id, nil
		}
	}
	return MaxID, fmt.Errorf("cannot generate standalone id")
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
		id, err := CreateHierarchyID(sid)
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
		id, err := CreateStandaloneID()
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
