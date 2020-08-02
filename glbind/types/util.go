package types

func align(s, a int) Alignment {
	if s < 0 {
		return Alignment{ByteSize: -s, ByteAlignment: a}
	}
	return Alignment{ByteSize: s, ByteAlignment: a}
}

// CreateArray returns a CType pointing to the provided CType in as
// many layers as needed to create arrays of dims.
func CreateArray(ct *CType, dims []int) *CType {
	if len(dims) == 0 {
		return ct
	}

	for i := len(dims) - 1; i >= 0; i-- {
		if dims[i] == 0 || dims[i] < -1 {
			panic("array dim cannot be 0 or < -1")
		}
		if dims[i] == -1 && i != 0 {
			panic("only last array index can be -1")
		}
		ct = &CType{
			Array: CArray{
				Len:   dims[i],
				CType: ct,
			},
			Size: align(ct.Size.ByteSize*dims[i], ct.Size.ByteAlignment),
		}
	}
	return ct
}
