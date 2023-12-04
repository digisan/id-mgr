package idmgr

import "fmt"

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
