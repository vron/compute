package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type vinfo struct {
	underlyingType string
	tdName         string
	typeSpec       string
}

func vi(a ...string) vinfo {
	return vinfo{
		underlyingType: a[0],
		tdName:         a[1],
		typeSpec:       a[2],
	}
}

var vtypes = []vinfo{
	vi("float", "float", ""),
	vi("int32_t", "int", "i"),
	vi("uint32_t", "uint", "u"),
}

var vsizes = []int{2, 3, 4}

func generateTypes(inp Input) {
	f, err := os.Create(filepath.Join(fOut, "types.hpp"))
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	buf := bufio.NewWriter(f)
	defer buf.Flush()

	buf.WriteString("#pragma once\n")
	buf.WriteString("// Code generated DO NOT EDIT\n")
	buf.WriteString("#include <cmath>\n\n")

	// Write some generiv defined
	buf.WriteString("#define INL  __attribute__((always_inline))\n\n")
	buf.WriteString(`
	

`)
	buf.WriteString("typedef int32_t Bool;\n\n")

	// TODO: This entire thing could probably be made much smaller with templates... ?
	generateVectorTypes(buf)
	generateCreateVectors(buf)
	generateImageTypes(buf)
	generateBuiltinFunctions(buf)
	generateMatrixTypes(buf)
	generateCreateMatrices(buf)
	generateUserStructs(buf, inp)
}

func generateVectorTypes(buf *bufio.Writer) {
	// Forward decl.
	for _, typ := range vtypes {
		for _, size := range vsizes {
			fmt.Fprintf(buf, "typedef %v %vvec%v __attribute__((ext_vector_type(%v)));;\n", typ.underlyingType, typ.typeSpec, size, size)
		}
	}
	buf.WriteString("\n\n")
}

func generateMatrixTypes(buf *bufio.Writer) {

	// TOOD: Is this modulu thing a reasonable thing to do?
	buf.WriteString(`
struct mat2;
struct mat3;
struct mat4;

struct mat2 {
	vec2 c[2];
	
	void from_api(float *arr) { this->c[0][0] = arr[0];this->c[0][1] = arr[1];this->c[1][0] = arr[2];this->c[1][1] = arr[3];};
	//void from_api(float arr[4]) { this->c[0][0] = arr[0];this->c[0][1] = arr[1];this->c[1][0] = arr[2];this->c[1][1] = arr[3];};
	vec2 &operator[](int index);

};

struct mat3 {
	vec3 c[3];

	void from_api(float *arr) { this->c[0][0] = arr[0];this->c[0][1] = arr[1];this->c[0][2] = arr[2];this->c[1][0] = arr[3];this->c[1][1] = arr[4];this->c[1][2] = arr[5];this->c[2][0] = arr[6];this->c[2][1] = arr[7];this->c[2][2] = arr[8];};
	//void from_api(float arr[9]) { this->c[0][0] = arr[0];this->c[0][1] = arr[1];this->c[0][2] = arr[2];this->c[1][0] = arr[3];this->c[1][1] = arr[4];this->c[1][2] = arr[5];this->c[2][0] = arr[6];this->c[2][1] = arr[7];this->c[2][2] = arr[8];};
	
	vec3 &operator[](int index);
};

struct mat4 {
	vec4 c[4];

	void from_api(float *arr) { this->c[0][0] = arr[0];this->c[0][1] = arr[1];this->c[0][2] = arr[2];this->c[0][3] = arr[3];this->c[1][0] = arr[4];this->c[1][1] = arr[5];this->c[1][2] = arr[6];this->c[1][3] = arr[7];this->c[2][0] = arr[8];this->c[2][1] = arr[9];this->c[2][2] = arr[10];this->c[2][3] = arr[11];this->c[3][0] = arr[12];this->c[3][1] = arr[13];this->c[3][2] = arr[14];this->c[3][3] = arr[15];};
	//void from_api(float arr[16]) { this->c[0][0] = arr[0];this->c[0][1] = arr[1];this->c[0][2] = arr[2];this->c[0][3] = arr[3];this->c[1][0] = arr[4];this->c[1][1] = arr[5];this->c[1][2] = arr[6];this->c[1][3] = arr[7];this->c[2][0] = arr[8];this->c[2][1] = arr[9];this->c[2][2] = arr[10];this->c[2][3] = arr[11];this->c[3][0] = arr[12];this->c[3][1] = arr[13];this->c[3][2] = arr[14];this->c[3][3] = arr[15];};
	
	vec4 &operator[](int index);


};
mat2 operator*(mat2 lhs, const mat2& rhs) {
	for(int i = 0; i < 2; i++) {
		vec2 row = make_vec2(lhs.c[0][i], lhs.c[1][i]);
		lhs.c[0][i] = dot(row, rhs.c[0]);
		lhs.c[1][i] = dot(row, rhs.c[1]);
	}
	return lhs;
};
mat3 operator*(mat3 lhs, const mat3& rhs) {
	for(int i = 0; i < 3; i++) {
		vec3 row = make_vec3(lhs.c[0][i], lhs.c[1][i], lhs.c[2][i]);
		lhs.c[0][i] = dot(row, rhs.c[0]);
		lhs.c[1][i] = dot(row, rhs.c[1]);
		lhs.c[2][i] = dot(row, rhs.c[2]);
	}
	return lhs;
};
mat4 operator*(mat4 lhs, const mat4& rhs) {
	for(int i = 0; i < 4; i++) {
		vec4 row = make_vec4(lhs.c[0][i], lhs.c[1][i], lhs.c[2][i], lhs.c[3][i]);
		lhs.c[0][i] = dot(row, rhs.c[0]);
		lhs.c[1][i] = dot(row, rhs.c[1]);
		lhs.c[2][i] = dot(row, rhs.c[2]);
		lhs.c[3][i] = dot(row, rhs.c[3]);
	}
	return lhs;
};
vec2 operator*(mat2 lhs, const vec2& rhs) {
	return make_vec2(lhs[0][0]*rhs[0] + lhs[1][0]*rhs[1],
				lhs[0][1]*rhs[1] + lhs[1][1]*rhs[1]);
};
vec3 operator*(mat3 lhs, const vec3& rhs) {
	return make_vec3(lhs[0][0]*rhs[0] + lhs[1][0]*rhs[1] + lhs[2][0]*rhs[2],
				lhs[0][1]*rhs[0] + lhs[1][1]*rhs[1] + lhs[2][1]*rhs[2],
				lhs[0][2]*rhs[0] + lhs[1][2]*rhs[1] + lhs[2][2]*rhs[2]);
};
vec4 operator*(mat4 lhs, const vec4& rhs) {
	return make_vec4(lhs[0][0]*rhs[0] + lhs[1][0]*rhs[1] + lhs[2][0]*rhs[2] + lhs[3][0]*rhs[3],
				lhs[0][1]*rhs[0] + lhs[1][1]*rhs[1] + lhs[2][1]*rhs[2] + lhs[3][1]*rhs[3],
				lhs[0][2]*rhs[0] + lhs[1][2]*rhs[1] + lhs[2][2]*rhs[2] + lhs[3][2]*rhs[3],
				lhs[0][3]*rhs[0] + lhs[1][3]*rhs[1] + lhs[2][3]*rhs[2] + lhs[3][3]*rhs[3]);
};
vec2 &mat2::operator[](int index) {
	index = index % 2;
	switch (index) {
	case 0:
	return (this->c[0]);
	case 1:
	return (this->c[1]);
	}
	__builtin_unreachable();
}
vec3 &mat3::operator[](int index) {
	index = index % 3;
	switch (index) {
	case 0:
	return (this->c[0]);
	case 1:
	return (this->c[1]);
	case 2:
	return (this->c[2]);
	}
	__builtin_unreachable();
}
vec4 &mat4::operator[](int index) {
	index = index % 4;
	switch (index) {
	case 0:
	return (this->c[0]);
	case 1:
	return (this->c[1]);
	case 2:
	return (this->c[2]);
	case 3:
	return (this->c[3]);
	}
	__builtin_unreachable();
}

`)
}

func compName(i int) string {
	switch i {
	case 0:
		return string('x')
	case 1:
		return string('y')
	case 2:
		return string('z')
	case 3:
		return string('w')
	}
	panic("bad length")
}

func generateImageTypes(buf *bufio.Writer) {
	// TODO: loop multiple image types, for now only support one...
	// TODO: can we add attributes such that the compiler will now the alignemnt of data (and chec in set_data), so it
	// can use aligned sse operations to store in e.g. a vec4 into it?

	buf.WriteString(`struct image2Drgba32f {
	float* data;
	int width;
  
	`)
	if types.Get("image2Drgba32f").apiType {
		buf.WriteString(`	void from_api(cpt_image2Drgba32f d) {
			this->data = (float*)d.data;
			this->width = d.width;
		};

		`)
	}
	buf.WriteString("};\n")

	buf.WriteString("\n\n")
}

func generateCreateVectors(buf *bufio.Writer) {
	for _, typ := range vtypes {
		for _, size := range vsizes {
			printSingleCreate(buf, typ, size)
			printMultipeCreate(buf, typ, size)
		}
	}
	// we also need make_int and friends as effectively no-ops but generated by the translation
	scalarTypes := [][2]string{ /*{"bool", "bool"},*/ {"int", "int32_t"}, {"uint", "uint32_t"}, {"float", "float"}}
	for _, s := range scalarTypes {
		for _, t := range scalarTypes {
			fmt.Fprintf(buf, "%v INL make_%v(%v n) { return (%v)n; }; \n", s[1], s[0], t[1], s[1])
		}
	}

	buf.WriteString("\n\n")
}

func printSingleCreate(buf *bufio.Writer, typ vinfo, size int) {
	fmt.Fprintf(buf, "%vvec%v INL make_%vvec%v(%v v) { return ", typ.typeSpec, size, typ.typeSpec, size, typ.underlyingType)

	fmt.Fprintf(buf, "%vvec%v{", typ.typeSpec, size)
	for i := 0; i < size; i++ {
		fmt.Fprintf(buf, "v")
		if i != size-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString("}; };\n")

}

func printSingleCreateMat(buf *bufio.Writer, size int) {
	// a single value only will be placed on diagonal
	for _, t := range []string{"uint32_t", "int32_t", "float", "double"} {
		fmt.Fprintf(buf, "mat%v INL make_mat%v(%v v) { mat%v m; ", size, size, t, size)
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				if i != j {
					fmt.Fprintf(buf, "m[%v][%v] = 0.0;", i, j)
				} else {
					fmt.Fprintf(buf, "m[%v][%v] = (float)v;", i, j)
				}
				buf.WriteString("; ")
			}
		}
		buf.WriteString("return m; };\n")
	}
}

func printMultipleCreateMat(buf *bufio.Writer, size int) {
	// for now only support floats (to limit the number that becomes large..)
	for no := 2; no <= size*size; no++ {
		args := ""
		body := ""
		for i := 0; i < no; i++ {
			args += fmt.Sprintf("float v%v,", i)

			cn := i / size
			rn := i % size

			body += fmt.Sprintf("m.c[%v][%v] = v%v; ", cn, rn, i)
		}
		args = args[:len(args)-1]

		fmt.Fprintf(buf, "mat%v INL make_mat%v("+args+") { mat%v m; ", size, size, size)

		buf.WriteString(body + "return m; };\n")
	}
}

func printVectorCreateMat(buf *bufio.Writer, size int) {
	// for now only support vectors of the correct type and the correct size:

	args := ""
	body := ""
	for i := 0; i < size; i++ {
		args += fmt.Sprintf("vec%v v%v,", size, i)
		body += fmt.Sprintf("m.c[%v] = v%v; ", i, i)
	}
	args = args[:len(args)-1]

	fmt.Fprintf(buf, "mat%v INL make_mat%v("+args+") { mat%v m; ", size, size, size)

	buf.WriteString(body + "return m; };\n")
}

func printMultipeCreate(buf *bufio.Writer, typ vinfo, size int) {
	options := []aopt{{1, "float"}, {1, "int32_t"}, {1, "uint32_t"}}
	for _, typ := range vtypes {
		for _, size := range vsizes {
			options = append(options, aopt{size, fmt.Sprintf("%vvec%v", typ.typeSpec, size)})
		}
	}
	sizes := make([][]aopt, 0)
	genPossibleSizes(&sizes, &[]aopt{}, size, options)
	for _, args := range sizes {
		// but each argument can also be of different sizes
		printCreate(buf, typ, size, args)
		_ = args
	}
}

func printCreate(buf *bufio.Writer, typ vinfo, size int, o []aopt) {
	fmt.Fprintf(buf, "%vvec%v INL make_%vvec%v(", typ.typeSpec, size, typ.typeSpec, size)
	for i, arg := range o {
		fmt.Fprintf(buf, "%v a%v", arg.typName, i)
		if i != len(o)-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString(") { return ")
	fmt.Fprintf(buf, "%vvec%v{", typ.typeSpec, size)
	nod := -1
	for i, arg := range o {
		if arg.size == 1 {
			nod++
			if arg.typName != typ.underlyingType {
				fmt.Fprintf(buf, "(%v)", typ.underlyingType)
			}
			fmt.Fprintf(buf, "(a%v)", i)
			if nod != size-1 {
				buf.WriteString(", ")
			}
		} else {
			for j := 0; j < arg.size; j++ {
				nod++
				if arg.typName != typ.underlyingType { // this should be underlying type?
					fmt.Fprintf(buf, "(%v)", typ.underlyingType)
				}
				fmt.Fprintf(buf, "(a%v.%v)", i, compName(j))
				if nod != size-1 {
					buf.WriteString(", ")
				}
			}
		}
	}
	buf.WriteString("}; };\n")
}

type aopt struct {
	size    int
	typName string
}

// generate all ways in which we can choose sizes from options such
// that the total adds up to target
func genPossibleSizes(all *[][]aopt, current *[]aopt, target int, options []aopt) {
	if target == 0 {
		*all = append(*all, *current)
		return // done on this one
	}
	for _, o := range options {
		if o.size > target {
			continue
		}
		cd := make([]aopt, len(*current))
		copy(cd, *current)
		cd = append(cd, o)
		genPossibleSizes(all, &cd, target-o.size, options)
	}
}

func generateCreateMatrices(buf *bufio.Writer) {
	// first similarly to vectors out of scalar values
	for _, size := range vsizes {
		printSingleCreateMat(buf, size)
		printMultipleCreateMat(buf, size)
	}

	// can also create vrom vector columns
	for _, size := range vsizes {
		printVectorCreateMat(buf, size)
	}
	// or from other matrices

	buf.WriteString("\n\n")
}

func generateBuiltinFunctions(buf *bufio.Writer) {
	//buf.WriteString("void barrier();\n")

	fimageStoreLoad(buf)
	fAtomic(buf)
	fGeometricFunctions(buf)
}

func fimageStoreLoad(buf *bufio.Writer) {
	// TODO: generate for types of images, also for non-2D images
	buf.WriteString(`void imageStore(image2Drgba32f image, ivec2 P, vec4 data) {
	int32_t index = 4*P.x+4*P.y*image.width;
	image.data[index + 0] = data.x;
	image.data[index + 1] = data.y;
	image.data[index + 2] = data.z;
	image.data[index + 3] = data.w;
}

`)
	buf.WriteString(`vec4 imageLoad(image2Drgba32f image, ivec2 P) {
	int32_t index = 4*P.x+4*P.y*image.width;
	return make_vec4(image.data[index + 0], image.data[index + 1], image.data[index + 2], image.data[index + 3]);
}

`)
}
func fAtomic(buf *bufio.Writer) {
	for _, t := range []string{"uint32_t", "int32_t"} {
		fmt.Fprintf(buf, "%v INL atomicAdd(%v *mem, %v data) {\n", t, t, t)
		fmt.Fprintf(buf, "\treturn __atomic_add_fetch(mem, data, __ATOMIC_SEQ_CST);\n")
		fmt.Fprintf(buf, "}\n\n")

		fmt.Fprintf(buf, "%v INL atomicAnd(%v *mem, %v data) {\n", t, t, t)
		fmt.Fprintf(buf, "\treturn __atomic_and_fetch(mem, data, __ATOMIC_SEQ_CST);\n")
		fmt.Fprintf(buf, "}\n\n")

		fmt.Fprintf(buf, "%v INL atomicOr(%v *mem, %v data) {\n", t, t, t)
		fmt.Fprintf(buf, "\treturn __atomic_or_fetch(mem, data, __ATOMIC_SEQ_CST);\n")
		fmt.Fprintf(buf, "}\n\n")

		fmt.Fprintf(buf, "%v INL atomicXor(%v *mem, %v data) {\n", t, t, t)
		fmt.Fprintf(buf, "\treturn __atomic_xor_fetch(mem, data, __ATOMIC_SEQ_CST);\n")
		fmt.Fprintf(buf, "}\n\n")

		fmt.Fprintf(buf, "%v INL atomicMin(%v *mem, %v data) {\n", t, t, t)
		fmt.Fprintf(buf, "\treturn __atomic_fetch_min(mem, data, __ATOMIC_SEQ_CST);\n")
		fmt.Fprintf(buf, "}\n\n")

		fmt.Fprintf(buf, "%v INL atomicMax(%v *mem, %v data) {\n", t, t, t)
		fmt.Fprintf(buf, "\treturn __atomic_fetch_max(mem, data, __ATOMIC_SEQ_CST);\n")
		fmt.Fprintf(buf, "}\n\n")

		fmt.Fprintf(buf, "%v INL atomicExchange(%v *mem, %v data) {\n", t, t, t)
		fmt.Fprintf(buf, "\treturn __atomic_exchange_n(mem, data, __ATOMIC_SEQ_CST);\n")
		fmt.Fprintf(buf, "}\n\n")

		// TODO: Is this one really correct?!
		fmt.Fprintf(buf, "%v INL atomicCompSwap(%v *mem, %v compare, %v data) {\n", t, t, t, t)
		fmt.Fprintf(buf, "\t__atomic_compare_exchange_n(mem, &compare, data, true, __ATOMIC_SEQ_CST, __ATOMIC_SEQ_CST);\n")
		fmt.Fprintf(buf, "\treturn compare;\n")
		fmt.Fprintf(buf, "}\n\n")

	}

}

func fGeometricFunctions(buf *bufio.Writer) {
	// There seem to only be defined for float vectors, is that correct?
	// TODO: Liely these implementations suffer from numeric problems and should be implemented in smarter ways...

	buf.WriteString(`
float cross(vec2 a, vec2 b) { return a[0]*b[1]-a[1]*b[0]; };
float dot(vec2 a, vec2 b) { return a[0] * b[0] + a[1] * b[1]; };
float length(vec2 a) { return sqrt(dot(a, a)); };
float distance(vec2 a, vec2 b) { return length(a - b); };
vec2 normalize(vec2 a) { return a / length(a); };
vec3 cross(vec3 a, vec3 b) {
	return make_vec3(a[1] * b[2] - a[2] * b[1], a[2] * b[0] - a[0] * b[2], a[0] * b[1] - a[1] * b[0]);
}
float dot(vec3 a, vec3 b) { return a[0] * b[0] + a[1] * b[1] + a[2] * b[2]; };
float length(vec3 a) { return sqrt(dot(a, a)); };
float distance(vec3 a, vec3 b) { return length(a - b); };
vec3 normalize(vec3 a) { return a / length(a); };
float dot(vec4 a, vec4 b) { return a[0] * b[0] + a[1] * b[1] + a[2] * b[2] + a[3] * b[3]; };
float length(vec4 a) { return sqrt(dot(a, a)); };
float distance(vec4 a, vec4 b) { return length(a - b); };
vec4 normalize(vec4 a) { return a / length(a); };
	`)

}

func generateUserStructs(buf *bufio.Writer, inp Input) {
	buf.WriteString("\n\n\n")

	for _, s := range types.UserStructs() {
		fmt.Fprintf(buf, "struct %v;\n", s.Name)
	}
	buf.WriteString("\n")
	for _, s := range types.UserStructs() {
		fmt.Fprintf(buf, "struct %v {// size=%v alignment=%v\n", s.Name, s.CType().Size.ByteSize, s.CType().Size.ByteAlignment)

		for _, f := range s.CType().Fields {
			//fmt.Fprintf(buf, "\t%v %v;  // size=%v alignment=%v offset=%v\n", f.Ty.ty.Name, f.Name)
			fmt.Fprintf(buf, "%v  // size=%v alignment=%v offset=%v\n", f.CxxFieldString(), f.Ty.Size.ByteSize, f.Ty.Size.ByteAlignment, f.ByteOffset)
		}

		// empty constructor
		fmt.Fprintf(buf, "\t %v () {};\n", s.Name)

		// write a constructor to allow function-style initialization
		//notused(vec3 ac[3], Bool bb, float cc) : aa{ac[0], ac[1], ac[2]}, bb(bb), cc(cc) {};

		fmt.Fprintf(buf, "\t%v(", s.Name)
		for i, f := range s.CType().Fields {
			// TODO: Must we allow type cases here to?
			arl := f.CxxArrayLen()
			ex := ""
			if arl > 0 {
				ex += fmt.Sprintf("[%v]", arl)
			}
			strt := ""
			if len(f.Ty.Fields) > 0 {
				strt = " struct "
			}
			fmt.Fprintf(buf, "%v%v %v"+ex, strt, f.Ty.ty.Name, f.Name)
			if i != len(s.CType().Fields)-1 {
				buf.WriteString(", ")
			}
		}
		buf.WriteString(") : ")
		for i, f := range s.CType().Fields {
			arl := f.CxxArrayLen()
			if arl > 0 {
				fmt.Fprintf(buf, "%v{", f.Name)
				fmt.Fprintf(buf, "%v[0]", f.Name)
				for j := 1; j < arl; j++ {
					fmt.Fprintf(buf, ", %v[%v]", f.Name, j)
				}
				fmt.Fprintf(buf, "}")
			} else {
				fmt.Fprintf(buf, "%v(%v)", f.Name, f.Name)
			}
			if i != len(s.CType().Fields)-1 {
				buf.WriteString(", ")
			}
		}
		fmt.Fprintf(buf, " {};\n")

		// we also need to construct a from_api function to copy over data from
		// the api structs - only arrays we enforce user to handle the alignment
		// WE simply declare it here and define it below, if we have to refer to
		// ourselves somewhere
		if s.apiType {
			fmt.Fprintf(buf, "\tvoid from_api(%v);\n", s.CType().Name)
		}

		buf.WriteString("};\n\n")
	}

	// Write implementation for the from api methods, copying field by field (rec if needed)
	for _, s := range types.UserStructs() {
		if s.apiType {
			fmt.Fprintf(buf, "void %v::from_api(%v d) {\n", s.Name, s.CType().Name)

			for _, cf := range s.CType().Fields {
				cf.CxxBinding(buf)
			}

			fmt.Fprintf(buf, "};\n\n")
		}

	}

	buf.WriteString("\n")
}

/*

func bDirect(arg InputArgument) string {
	return fmt.Sprintf("\tthis->%v = d.%v;\n", arg.CxxName(), arg.CName())
}

func bVec(size int) func(InputArgument) string {
	return func(arg InputArgument) string {
		// bind a value by copying vector data
		s := ""
		for i := 0; i < size; i++ {
			s += fmt.Sprintf("\tthis->%v.%v = d.%v[%v];\n", arg.CxxName(), compName(i), arg.CName(), i)
		}
		return s
	}
}


*/

func (cf CField) CxxBinding(buf io.Writer) {
	if cf.Ty.IsSlice {
		// slice type, the incomeing is *void and we assume everyhting is layed out, assign!
		fmt.Fprintf(buf, "\tthis->%v = (%v*)d.%v;\n", cf.Name, cf.Ty.ty.Name, cf.Name)
		return
	}
	if cf.Ty.ArrayLen == 0 && (len(cf.Ty.Fields) > 0 || cf.Ty.ty.Name == "mat2" || cf.Ty.ty.Name == "mat3" || cf.Ty.ty.Name == "mat4") {
		// this is a struct, assign each one of them, this must be done recursively!
		fmt.Fprintf(buf, "\t(&(this->%v))->from_api(d.%v);\n", cf.Name, cf.Name)
		return
	}
	if cf.Ty.ArrayLen == 0 {
		fmt.Fprintf(buf, "\tthis->%v = d.%v;\n", cf.Name, cf.Name)
		return
	}

	// this is a matrix, for now handle it specifically until we figure out a way

	// so it is an array of stuff, do the same for each once of the elements, but as a temp hac
	// chec for the underlying type to set from that if needed
	arrlen := cf.CxxArrayLen()
	vecSize := cf.Ty.ArrayLen
	if arrlen > 0 {
		vecSize = cf.Ty.ArrayLen / arrlen
	} else {
		arrlen = 1
	}
	for i := 0; i < arrlen; i++ {
		if len(cf.Ty.Fields) > 0 {
			// this is a struct, assign each one of them, this must be done recursively!
			fmt.Fprintf(buf, "\t(&(this->%v[%v]))->from_api(d.%v[%v]);\n", cf.Name, i, cf.Name, i)
		} else if cf.Ty.ty.Name == "mat2" || cf.Ty.ty.Name == "mat3" || cf.Ty.ty.Name == "mat4" {
			// binf it...

			if arrlen > 1 {
				fmt.Fprintf(buf, "\t(&(this->%v[%v]))->from_api(&d.%v[%v]);// mat bind\n", cf.Name, i, cf.Name, i*vecSize)
			} else {
				fmt.Fprintf(buf, "\t(&(this->%v))->from_api(d.%v);// mat bind\n", cf.Name, cf.Name)
			}
		} else {
			if arrlen > 1 {
				for j := 0; j < vecSize; j++ {
					fmt.Fprintf(buf, "\tthis->%v[%v][%v] = d.%v[%v];\n", cf.Name, i, j, cf.Name, i*vecSize+j)
				}
			} else {
				for j := 0; j < vecSize; j++ {
					fmt.Fprintf(buf, "\tthis->%v[%v] = d.%v[%v];\n", cf.Name, j, cf.Name, j)
				}
			}
		}
	}
}
