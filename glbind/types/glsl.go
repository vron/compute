package types

type GlslType struct {
	Name   string
	Export bool
	C      *CType
}

func (t GlslType) CName(prefix string) string {
	return prefix + t.Name
}
