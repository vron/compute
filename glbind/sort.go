package main

// TODO: This is a highly inefficient way to do the sorting, fix it when it becomes a problem...
func less(ts []*Type) func(i, j int) bool {
	return func(i, j int) bool {
		// first chec if one of the types directly or indirectly imports the other, if so
		// it should be sorted before...
		ti, tj := ts[i], ts[j]

		if ti.userStructId > 0 && tj.userStructId > 0 {
			return ti.userStructId < tj.userStructId
		}

		if dependsOn(ti.cType, tj) {
			return false
		}
		if dependsOn(tj.cType, ti) {
			return true
		}
		return ts[i].Name < ts[j].Name
	}
}

// true if ta depends on tb directly or indirectly
func dependsOn(ta *CType, tb *Type) bool {
	if ta.ty.Name == tb.Name {
		return true
	}
	for _, f := range ta.Fields {
		if dependsOn(f.Ty, tb) {
			return true
		}
	}
	return false
}
