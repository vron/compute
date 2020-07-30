package input

type Input struct {
	Arguments []InputArgument
	Shared    []InputArgument
	Structs   []InputStruct
	Wg_size   [3]int
	Body      string
}

type InputArgument struct {
	Name  string
	Ty    string
	Arrno []int
}

type InputStruct struct {
	Name   string
	Fields []InputArgument
}
