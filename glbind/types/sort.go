package types

func less(ts []*GlslType) func(i, j int) bool {
	return func(i, j int) bool {
		// first chec if one of the types directly or indirectly imports the other, if so
		// it should be sorted before...
		ti, tj := ts[i], ts[j]

		if ti.GlslOrder > 0 && tj.GlslOrder > 0 {
			return ti.GlslOrder < tj.GlslOrder
		}

		if dependsOn(ti.C, tj) {
			return false
		}
		if dependsOn(tj.C, ti) {
			return true
		}
		return ts[i].Name < ts[j].Name
	}
}

// true if ta depends on tb directly or indirectly
func dependsOn(ta *CType, tb *GlslType) bool {
	if ta.GlslType.Name == tb.Name {
		return true
	}
	for _, f := range ta.Fields {
		if dependsOn(f.Ty, tb) {
			return true
		}
	}
	return false
}
