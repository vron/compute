package types

type Alignment struct {
	ByteSize      int
	ByteAlignment int
}

func (ts *Types) calculateAlignments() {
	for _, t := range ts.l {
		alignVisit(t.C)
	}
}

func alignVisit(ct *CType) Alignment {
	// from https://www.oreilly.com/library/view/opengl-programming-guide/9780132748445/app09lev1sec3.html
	// but we treat vec3 as vec4 to allow for vector operations

	if ct.Size != (Alignment{}) {
		return ct.Size // allready handled, nothing to do
	}
	if ct.IsStruct() {
		return alignVisitStruct(ct)
	} else if ct.IsArray() {
		return alignVisitArray(ct)
	}
	panic("all basics and vectors should have pre-set alignments")
}

func alignVisitStruct(ct *CType) Alignment {
	offset := 0
	maxFieldAlignment := 0
	for i, f := range ct.Struct.Fields {
		fa := alignVisit(f.CType)

		for ; offset%fa.ByteAlignment != 0; offset++ {
		}
		ct.Struct.Fields[i].ByteOffset = offset
		if fa.ByteAlignment > maxFieldAlignment {
			maxFieldAlignment = fa.ByteAlignment
		}
		offset += fa.ByteSize
	}

	for ; offset%maxFieldAlignment != 0; offset++ {
	}

	ct.Size.ByteSize = offset
	ct.Size.ByteAlignment = maxFieldAlignment
	return ct.Size
}

func alignVisitArray(ct *CType) Alignment {
	fa := alignVisit(ct.Array.CType)
	ct.Size = align(fa.ByteSize*ct.Array.Len, fa.ByteAlignment)
	return ct.Size
}
