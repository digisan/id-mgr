package idmgr

import "fmt"

func DelID(id ID) error {

	if _, ok := mRecord[id]; !ok {
		return nil
	}

	// mRecord for standalone, adjust it
	if id.IsStandalone() {
		delete(mRecord, id)
		delete(mAlias, id)
		mRecord[MaxID]--
	} else {
		if descIDs := id.Descendants(1); len(descIDs) > 0 {
			return fmt.Errorf("%x(%v) has descendants [%x], cannot delete, abort", id, id.Alias(), descIDs)
		}
		delete(mRecord, id)
		delete(mAlias, id)
		// BUT, DO NOT modify parent mRecord for hierarchy!!!
		// if parent, ok := id.Parent(); ok {
		// 	mRecord[parent]--
		// }
	}

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

// DelIDViaAlias incurs updated WholeIDs
func DelIDViaAlias(alias any) error {
	fmt.Println(WholeIDs())
	id := SearchIDByAlias(alias, WholeIDs()...)
	if len(fmt.Sprint(alias)) > 0 && id == 0 {
		return fmt.Errorf("alias [%s] cannot be found, nothing to delete", alias)
	}
	return DelID(id)
}

func DelIDsOnAlias(aliases ...any) error {
	for _, alias := range aliases {
		if err := DelIDViaAlias(alias); err != nil {
			return err
		}
	}
	return nil
}
