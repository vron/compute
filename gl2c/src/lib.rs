// TODO: Handle shared types and variables so we share them correctly between threads...

// effectively it seems as if we can treet uniforms and buffers the same (uniform read only, but we do not care)
use std::cell::RefCell;
use std::clone::Clone;
//use std::fmt;
use std::fmt::Write as FWrite;
use std::io::Read;
use std::io::Write;
use std::marker::Copy;

extern crate serde;
extern crate serde_json;
#[macro_use]
extern crate serde_derive;

extern crate glsl;
use glsl::parser::Parse;
use glsl::syntax;

#[derive(Serialize, Deserialize, Debug)]
pub struct Struct {
    name: String,
    fields: Vec<Argument>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct Export {
    arguments: Vec<Argument>,
    shared: Vec<Argument>,
    structs: Vec<Struct>,
    wg_size: [i32; 3],
    body: String,
}

#[derive(Clone, Serialize, Deserialize, Debug)]
pub struct Argument {
    name: String,
    ty: String,
    arrno: Vec<i32>, // -1 for slice, 0 for is not array
}

/*
impl fmt::Display for Argument {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{}, {}, {}", self.name, self.ty, self.arrno)
    }
}
*/

pub struct State {
    body: Part,
    arguments: Vec<Argument>,
    shared: Vec<Argument>,
    structs: Vec<Struct>,
    wg_size: [i32; 3],

    output: Output,
    last_output: Vec<Output>,
}

pub struct Part {
    d: RefCell<String>,
}

impl Part {
    pub fn new() -> Part {
        Part {
            d: RefCell::new(String::new()),
        }
    }

    pub fn as_str(&mut self) -> String {
        self.d.borrow_mut().clone().clone()
    }
}

#[derive(Clone, Copy)]
enum Output {
    None,
    Body,
}

impl State {
    pub fn add_arg(&mut self, arg: Argument) {
        self.arguments.push(arg.clone());
    }
    pub fn add_shared(&mut self, arg: Argument) {
        self.shared.push(arg.clone());
    }

    pub fn write_json(mut self) -> String {
        let mut s = String::new();
        s.write_str(&self.body.as_str()[..]).expect("");
        let out = Export {
            arguments: self.arguments.iter().map(|b| b.clone()).collect(),
            shared: self.shared.iter().map(|b| b.clone()).collect(),
            structs: self.structs,
            wg_size: self.wg_size,
            body: s,
        };
        serde_json::to_string_pretty(&out).unwrap()
    }

    fn push_output(&mut self, o: Output) {
        self.last_output.push(self.output);
        self.output = o;
    }

    fn pop_output(&mut self) {
        self.output = self.last_output.pop().unwrap();
    }
}

impl std::io::Write for State {
    fn write(&mut self, buf: &[u8]) -> std::io::Result<usize> {
        match self.output {
            Output::None => (),
            Output::Body => {
                self.body.write(buf).expect("");
            }
        };
        std::io::Result::Ok(buf.len())
    }

    fn flush(&mut self) -> std::io::Result<()> {
        std::io::Result::Ok(())
    }
}

impl std::io::Write for Part {
    fn write(&mut self, buf: &[u8]) -> std::io::Result<usize> {
        self.d
            .borrow_mut()
            .write_str(std::str::from_utf8(buf).unwrap())
            .unwrap();
        std::io::Result::Ok(buf.len())
    }

    fn flush(&mut self) -> std::io::Result<()> {
        std::io::Result::Ok(())
    }
}

pub fn visit_identifier(s: &mut State, i: &syntax::Identifier) {
    let _ = write!(s, "{}", &i.0);
}

pub fn visit_type_name(_s: &mut State, t: &syntax::TypeName) -> String {
    String::from(&t.0)
    //let _ = write!(s, "{}", &t.0);
}

pub fn type_specifier_non_array(s: &mut State, t: &syntax::TypeSpecifierNonArray) -> String {
    let temp: String;
    String::from(match *t {
        syntax::TypeSpecifierNonArray::Void => "void",
        syntax::TypeSpecifierNonArray::Bool => "Bool",
        syntax::TypeSpecifierNonArray::Int => "int32_t",
        syntax::TypeSpecifierNonArray::UInt => "uint32_t",
        syntax::TypeSpecifierNonArray::Float => "float",
        syntax::TypeSpecifierNonArray::Double => "double",
        syntax::TypeSpecifierNonArray::Vec2 => "vec2",
        syntax::TypeSpecifierNonArray::Vec3 => "vec3",
        syntax::TypeSpecifierNonArray::Vec4 => "vec4",
        syntax::TypeSpecifierNonArray::DVec2 => "dvec2",
        syntax::TypeSpecifierNonArray::DVec3 => "dvec3",
        syntax::TypeSpecifierNonArray::DVec4 => "dvec4",
        syntax::TypeSpecifierNonArray::BVec2 => "bvec2",
        syntax::TypeSpecifierNonArray::BVec3 => "bvec3",
        syntax::TypeSpecifierNonArray::BVec4 => "bvec4",
        syntax::TypeSpecifierNonArray::IVec2 => "ivec2",
        syntax::TypeSpecifierNonArray::IVec3 => "ivec3",
        syntax::TypeSpecifierNonArray::IVec4 => "ivec4",
        syntax::TypeSpecifierNonArray::UVec2 => "uvec2",
        syntax::TypeSpecifierNonArray::UVec3 => "uvec3",
        syntax::TypeSpecifierNonArray::UVec4 => "uvec4",
        syntax::TypeSpecifierNonArray::Mat2 => "mat2",
        syntax::TypeSpecifierNonArray::Mat3 => "mat3",
        syntax::TypeSpecifierNonArray::Mat4 => "mat4",
        syntax::TypeSpecifierNonArray::Mat23 => "mat23",
        syntax::TypeSpecifierNonArray::Mat24 => "mat24",
        syntax::TypeSpecifierNonArray::Mat32 => "mat32",
        syntax::TypeSpecifierNonArray::Mat34 => "mat34",
        syntax::TypeSpecifierNonArray::Mat42 => "mat42",
        syntax::TypeSpecifierNonArray::Mat43 => "mat43",
        syntax::TypeSpecifierNonArray::DMat2 => "dmat2",
        syntax::TypeSpecifierNonArray::DMat3 => "dmat3",
        syntax::TypeSpecifierNonArray::DMat4 => "dmat4",
        syntax::TypeSpecifierNonArray::DMat23 => "dmat23",
        syntax::TypeSpecifierNonArray::DMat24 => "dmat24",
        syntax::TypeSpecifierNonArray::DMat32 => "dmat32",
        syntax::TypeSpecifierNonArray::DMat34 => "dmat34",
        syntax::TypeSpecifierNonArray::DMat42 => "dmat42",
        syntax::TypeSpecifierNonArray::DMat43 => "dmat43",
        syntax::TypeSpecifierNonArray::Sampler1D => "sampler1D",
        syntax::TypeSpecifierNonArray::Image1D => "image1D",
        syntax::TypeSpecifierNonArray::Sampler2D => "sampler2D",
        syntax::TypeSpecifierNonArray::Image2D => "image2D",
        syntax::TypeSpecifierNonArray::Sampler3D => "sampler3D",
        syntax::TypeSpecifierNonArray::Image3D => "image3D",
        syntax::TypeSpecifierNonArray::SamplerCube => "samplerCube",
        syntax::TypeSpecifierNonArray::ImageCube => "imageCube",
        syntax::TypeSpecifierNonArray::Sampler2DRect => "sampler2DRect",
        syntax::TypeSpecifierNonArray::Image2DRect => "image2DRect",
        syntax::TypeSpecifierNonArray::Sampler1DArray => "sampler1DArray",
        syntax::TypeSpecifierNonArray::Image1DArray => "image1DArray",
        syntax::TypeSpecifierNonArray::Sampler2DArray => "sampler2DArray",
        syntax::TypeSpecifierNonArray::Image2DArray => "image2DArray",
        syntax::TypeSpecifierNonArray::SamplerBuffer => "samplerBuffer",
        syntax::TypeSpecifierNonArray::ImageBuffer => "imageBuffer",
        syntax::TypeSpecifierNonArray::Sampler2DMS => "sampler2DMS",
        syntax::TypeSpecifierNonArray::Image2DMS => "image2DMS",
        syntax::TypeSpecifierNonArray::Sampler2DMSArray => "sampler2DMSArray",
        syntax::TypeSpecifierNonArray::Image2DMSArray => "image2DMSArray",
        syntax::TypeSpecifierNonArray::SamplerCubeArray => "samplerCubeArray",
        syntax::TypeSpecifierNonArray::ImageCubeArray => "imageCubeArray",
        syntax::TypeSpecifierNonArray::Sampler1DShadow => "sampler1DShadow",
        syntax::TypeSpecifierNonArray::Sampler2DShadow => "sampler2DShadow",
        syntax::TypeSpecifierNonArray::Sampler2DRectShadow => "sampler2DRectShadow",
        syntax::TypeSpecifierNonArray::Sampler1DArrayShadow => "sampler1DArrayShadow",
        syntax::TypeSpecifierNonArray::Sampler2DArrayShadow => "sampler2DArrayShadow",
        syntax::TypeSpecifierNonArray::SamplerCubeShadow => "samplerCubeShadow",
        syntax::TypeSpecifierNonArray::SamplerCubeArrayShadow => "samplerCubeArrayShadow",
        syntax::TypeSpecifierNonArray::ISampler1D => "isampler1D",
        syntax::TypeSpecifierNonArray::IImage1D => "iimage1D",
        syntax::TypeSpecifierNonArray::ISampler2D => "isampler2D",
        syntax::TypeSpecifierNonArray::IImage2D => "iimage2D",
        syntax::TypeSpecifierNonArray::ISampler3D => "isampler3D",
        syntax::TypeSpecifierNonArray::IImage3D => "iimage3D",
        syntax::TypeSpecifierNonArray::ISamplerCube => "isamplerCube",
        syntax::TypeSpecifierNonArray::IImageCube => "iimageCube",
        syntax::TypeSpecifierNonArray::ISampler2DRect => "isampler2DRect",
        syntax::TypeSpecifierNonArray::IImage2DRect => "iimage2DRect",
        syntax::TypeSpecifierNonArray::ISampler1DArray => "isampler1DArray",
        syntax::TypeSpecifierNonArray::IImage1DArray => "iimage1DArray",
        syntax::TypeSpecifierNonArray::ISampler2DArray => "isampler2DArray",
        syntax::TypeSpecifierNonArray::IImage2DArray => "iimage2DArray",
        syntax::TypeSpecifierNonArray::ISamplerBuffer => "isamplerBuffer",
        syntax::TypeSpecifierNonArray::IImageBuffer => "iimageBuffer",
        syntax::TypeSpecifierNonArray::ISampler2DMS => "isampler2MS",
        syntax::TypeSpecifierNonArray::IImage2DMS => "iimage2DMS",
        syntax::TypeSpecifierNonArray::ISampler2DMSArray => "isampler2DMSArray",
        syntax::TypeSpecifierNonArray::IImage2DMSArray => "iimage2DMSArray",
        syntax::TypeSpecifierNonArray::ISamplerCubeArray => "isamplerCubeArray",
        syntax::TypeSpecifierNonArray::IImageCubeArray => "iimageCubeArray",
        syntax::TypeSpecifierNonArray::AtomicUInt => "atomic_uint",
        syntax::TypeSpecifierNonArray::USampler1D => "usampler1D",
        syntax::TypeSpecifierNonArray::UImage1D => "uimage1D",
        syntax::TypeSpecifierNonArray::USampler2D => "usampler2D",
        syntax::TypeSpecifierNonArray::UImage2D => "uimage2D",
        syntax::TypeSpecifierNonArray::USampler3D => "usampler3D",
        syntax::TypeSpecifierNonArray::UImage3D => "uimage3D",
        syntax::TypeSpecifierNonArray::USamplerCube => "usamplerCube",
        syntax::TypeSpecifierNonArray::UImageCube => "uimageCube",
        syntax::TypeSpecifierNonArray::USampler2DRect => "usampler2DRect",
        syntax::TypeSpecifierNonArray::UImage2DRect => "uimage2DRect",
        syntax::TypeSpecifierNonArray::USampler1DArray => "usampler1DArray",
        syntax::TypeSpecifierNonArray::UImage1DArray => "uimage1DArray",
        syntax::TypeSpecifierNonArray::USampler2DArray => "usampler2DArray",
        syntax::TypeSpecifierNonArray::UImage2DArray => "uimage2DArray",
        syntax::TypeSpecifierNonArray::USamplerBuffer => "usamplerBuffer",
        syntax::TypeSpecifierNonArray::UImageBuffer => "uimageBuffer",
        syntax::TypeSpecifierNonArray::USampler2DMS => "usampler2DMS",
        syntax::TypeSpecifierNonArray::UImage2DMS => "uimage2DMS",
        syntax::TypeSpecifierNonArray::USampler2DMSArray => "usamplerDMSArray",
        syntax::TypeSpecifierNonArray::UImage2DMSArray => "uimage2DMSArray",
        syntax::TypeSpecifierNonArray::USamplerCubeArray => "usamplerCubeArray",
        syntax::TypeSpecifierNonArray::UImageCubeArray => "uimageCubeArray",
        syntax::TypeSpecifierNonArray::Struct(ref ss) => {
            temp = visit_struct_non_declaration(s, ss);
            &temp[..]
        }
        syntax::TypeSpecifierNonArray::TypeName(ref tn) => {
            temp = visit_type_name(s, tn);
            &temp[..]
        }
    })
}

pub fn visit_type_specifier(s: &mut State, t: &syntax::TypeSpecifier) -> (String, i32) {
    // Habndle the array types here!
    let ty = type_specifier_non_array(s, &t.ty);
    let _ = write!(s, "{}", ty);

    if let Some(ref arr_spec) = t.array_specifier {
        let l = match visit_array_spec(s, arr_spec) {
            Some(v) => v,
            None => panic!("arrays with non - constant lengths not supported yet"),
        };
        return (ty, l);
    }
    (ty, 0)
}

pub fn visit_fully_specified_type(s: &mut State, t: &syntax::FullySpecifiedType) {
    if let Some(ref qual) = t.qualifier {
        visit_type_qualifier(s, &qual);
        let _ = write!(s, "{}", " ");
    }

    // So we are here in a global struct defn. let _ = write!(s, "{}", " FFF ");
    visit_type_specifier(s, &t.ty);
}

pub fn visit_struct_non_declaration(s: &mut State, ss: &syntax::StructSpecifier) -> String {
    s.push_output(Output::None);
    let _ = write!(s, "{}", "struct ");

    let mut nn = String::from("NONAMESTRUCT");
    if let Some(ref name) = ss.name {
        let _ = write!(s, "{} ", name);
        nn = String::from(name.as_str());
    }

    let _ = write!(s, "{}", "{\n");

    let mut fields: Vec<Argument> = Vec::new();
    for field in &ss.fields.0 {
        let fld = visit_struct_field(s, field);
        fields.push(fld);
    }

    let _ = write!(s, "{}", "}");
    s.pop_output();
    s.structs.push(Struct {
        name: String::from(&nn[..]),
        fields: fields,
    });
    nn
}

pub fn visit_struct_field(s: &mut State, field: &syntax::StructFieldSpecifier) -> Argument {
    if let Some(ref qual) = field.qualifier {
        visit_type_qualifier(s, &qual);
        let _ = write!(s, "{}", " ");
    }

    let (tys, _) = visit_type_specifier(s, &field.ty);
    let _ = write!(s, "{}", " ");

    // there’s at least one identifier
    let mut identifiers = field.identifiers.0.iter();
    let identifier = identifiers.next().unwrap();

    let arrno = match visit_arrayed_identifier(s, identifier) {
        None => vec!(),
        Some(v) => vec!(v),
    };

    // write the rest of the identifiers
    for _identifier in identifiers {
        panic!("not supported yet - does this mean one type mutiple variables??")
        //let _ = write!(s, "{}", ", ");
        //visit_arrayed_identifier(s, identifier);
    }

    let _ = write!(s, "{}", ";\n");
    Argument {
        name: String::from(identifier.ident.as_str()),
        ty: tys,
        arrno: arrno,
    }
}

pub fn visit_array_spec(s: &mut State, a: &syntax::ArraySpecifier) -> Option<i32> {
    match *a {
        syntax::ArraySpecifier::Unsized => {
            let _ = write!(s, "{}", "[]");
            return Some(-1);
        }
        syntax::ArraySpecifier::ExplicitlySized(ref e) => {
            let _ = write!(s, "{}", "[");
            let a = visit_expr(s, &e);
            let _ = write!(s, "{}", "]");
            a
        }
    }
}

pub fn visit_arrayed_identifier(s: &mut State, a: &syntax::ArrayedIdentifier) -> Option<i32> {
    let _ = write!(s, "{}", a.ident);

    if let Some(ref arr_spec) = a.array_spec {
        return visit_array_spec(s, arr_spec);
    }
    None
}

pub fn non_accepted_layout_spec(s: &mut State, q: &syntax::TypeQualifier) -> bool {
    let mut was: bool = false;
    let qualifiers = q.qualifiers.0.iter();
    for qual_spec in qualifiers {
        match *qual_spec {
            syntax::TypeQualifierSpec::Layout(ref l) => {
                let qualifiers = l.ids.0.iter();
                for qual_spec in qualifiers {
                    let _ = write!(s, "{}", ", ");
                    match *qual_spec {
                        syntax::LayoutQualifierSpec::Identifier(ref i, None) => {
                            if i.0 == "std430" {
                                was = true;
                            }
                        }
                        _ => (),
                    }
                }
            }
            _ => {
                return false;
            }
        }
    }
    was
}

pub fn visit_type_qualifier(s: &mut State, q: &syntax::TypeQualifier) {
    let mut qualifiers = q.qualifiers.0.iter();
    let first = qualifiers.next().unwrap();

    visit_type_qualifier_spec(s, first);

    for qual_spec in qualifiers {
        let _ = write!(s, "{}", " ");
        visit_type_qualifier_spec(s, qual_spec)
    }
}

pub fn visit_type_qualifier_spec(s: &mut State, q: &syntax::TypeQualifierSpec) {
    match *q {
        syntax::TypeQualifierSpec::Storage(ref ss) => visit_storage_qualifier(s, &ss),
        syntax::TypeQualifierSpec::Layout(ref l) => visit_layout_qualifier(s, &l),
        syntax::TypeQualifierSpec::Precision(ref p) => visit_precision_qualifier(s, &p),
        syntax::TypeQualifierSpec::Interpolation(ref i) => visit_interpolation_qualifier(s, &i),
        syntax::TypeQualifierSpec::Invariant => {
            let _ = write!(s, "{}", "invariant");
        }
        syntax::TypeQualifierSpec::Precise => {
            let _ = write!(s, "{}", "precise");
        }
    }
}

pub fn visit_storage_qualifier(s: &mut State, q: &syntax::StorageQualifier) {
    match *q {
        syntax::StorageQualifier::Const => {
            let _ = write!(s, "{}", "const");
        }
        syntax::StorageQualifier::InOut => {
            let _ = write!(s, "{}", "inout");
        }
        syntax::StorageQualifier::In => {
            let _ = write!(s, "{}", "in");
        }
        syntax::StorageQualifier::Out => {
            let _ = write!(s, "{}", "out");
        }
        syntax::StorageQualifier::Centroid => {
            let _ = write!(s, "{}", "centroid");
        }
        syntax::StorageQualifier::Patch => {
            let _ = write!(s, "{}", "patch");
        }
        syntax::StorageQualifier::Sample => {
            let _ = write!(s, "{}", "sample");
        }
        syntax::StorageQualifier::Uniform => {
            let _ = write!(s, "{}", "uniform");
        }
        syntax::StorageQualifier::Attribute => {
            let _ = write!(s, "{}", "attribute");
        }
        syntax::StorageQualifier::Varying => {
            let _ = write!(s, "{}", "varying");
        }
        syntax::StorageQualifier::Buffer => {
            let _ = write!(s, "{}", "buffer");
        }
        syntax::StorageQualifier::Shared => {
            let _ = write!(s, "{}", "shared");
        }
        syntax::StorageQualifier::Coherent => {
            let _ = write!(s, "{}", "coherent");
        }
        syntax::StorageQualifier::Volatile => {
            let _ = write!(s, "{}", "volatile");
        }
        syntax::StorageQualifier::Restrict => {
            let _ = write!(s, "{}", "restrict");
        }
        syntax::StorageQualifier::ReadOnly => {
            let _ = write!(s, "{}", "readonly");
        }
        syntax::StorageQualifier::WriteOnly => {
            let _ = write!(s, "{}", "writeonly");
        }
        syntax::StorageQualifier::Subroutine(ref n) => visit_subroutine(s, &n),
    }
}

pub fn visit_subroutine(s: &mut State, types: &Vec<syntax::TypeName>) {
    let _ = write!(s, "{}", "subroutine");

    if !types.is_empty() {
        let _ = write!(s, "{}", "(");

        let mut types_iter = types.iter();
        let first = types_iter.next().unwrap();

        visit_type_name(s, first);

        for type_name in types_iter {
            let _ = write!(s, "{}", ", ");
            visit_type_name(s, type_name);
        }

        let _ = write!(s, "{}", ")");
    }
}

pub fn visit_layout_qualifier(s: &mut State, l: &syntax::LayoutQualifier) {
    let mut qualifiers = l.ids.0.iter();
    let first = qualifiers.next().unwrap();

    let _ = write!(s, "{}", "layout (");
    let _ = visit_layout_qualifier_spec(s, first);

    for qual_spec in qualifiers {
        let _ = write!(s, "{}", ", ");
        visit_layout_qualifier_spec(s, qual_spec);
    }

    let _ = write!(s, "{}", ")");
}

pub fn maybe_handle_wg(state: &mut State, q: &syntax::TypeQualifier) -> Option<i32> {
    // try to chec if this is a wg size spec, if so eat it, else let it fall through
    let mut qualifiers = q.qualifiers.0.iter();
    let first = qualifiers.next().unwrap();

    let l = match first {
        syntax::TypeQualifierSpec::Layout(ref l) => l,
        _ => return None,
    };
    let qualifiers = l.ids.0.iter();
    for qual_spec in qualifiers {
        match qual_spec {
            syntax::LayoutQualifierSpec::Identifier(ref i, Some(ref e)) => {
                let s = i.as_str();
                if s == "local_size_x" {
                    state.wg_size[0] = as_number(e);
                } else if s == "local_size_y" {
                    state.wg_size[1] = as_number(e);
                } else if s == "local_size_z" {
                    state.wg_size[2] = as_number(e);
                } else {
                    return None; // was not a layout
                }
            }
            _ => return None, // they must specify size
        };
    }

    Some(0)
}

pub fn maybe_handle_global_buffer(state: &mut State, l: &syntax::InitDeclaratorList) -> Option<()> {
    let d = &l.head;

    let mut subtype = String::from("");

    // chec if this is a layout statement, and if so find the first qualifier in that
    if let Some(q) = &d.ty.qualifier {
        let mut qualifiers = q.qualifiers.0.iter();
        let first = qualifiers.next().unwrap();
        let l = match first {
            syntax::TypeQualifierSpec::Layout(ref l) => l,
            _ => return None,
        };
        {
            let qualifiers = l.ids.0.iter();
            for qual_spec in qualifiers {
                match qual_spec {
                    syntax::LayoutQualifierSpec::Identifier(ref i,None) => {
                        subtype = String::from(i.as_str());
                        break;
                    }
                    _ => (), // we require a specifier for images ( TODO: can this be other data?)
                };
            }
        }

    // find the subsequent ones, "buffer" and "uniform" we consider the same, then find the actual type?
        for qs in qualifiers {
            match qs {
                syntax::TypeQualifierSpec::Storage(ref sq) => {
                    match sq {
                        syntax::StorageQualifier::Uniform => (),
                        syntax::StorageQualifier::Buffer => (),
                        _ => return None // not a global input / output
                    }
                },
                _ => return None,
            }
        }
    } else {
        return None;
    }

    let mut typname = type_specifier_non_array(state, &d.ty.ty.ty);

    // if this has a tail we should not handle it
    for _decl in &l.tail {
        return None;
    }
    if let Some(ref name) = d.name {    // Add this argument
        typname.push_str(&subtype[..]);
        state.add_arg(Argument{
            name: String::from(name.as_str()),
            ty: typname,
            arrno: vec!(),
        });
        return Some(())
    }

    None

}

pub fn maybe_handle_shared(state: &mut State, l: &syntax::InitDeclaratorList) -> Option<()> {
    let d = &l.head;


    let subtype = String::from("");

    // chec if this is a layout statement, and if so find the first qualifier in that
    if let Some(q) = &d.ty.qualifier {
        let mut qualifiers = q.qualifiers.0.iter();
        let first = qualifiers.next().unwrap();
        let s= match first {
            syntax::TypeQualifierSpec::Storage(ref s) => s,
            _ => return None,
        };
    
        match s {
            syntax::StorageQualifier::Shared => (),
            _ => return None,
        }

    } else {
        return None;
    }

    let mut typname = type_specifier_non_array(state, &d.ty.ty.ty);

    // if this has a tail we should not handle it
    for _decl in &l.tail {
        return None;
    }

    // if there is any array part add that
    let mut arrno = vec!();
    if let Some(a) = &l.head.array_specifier {
        let len = visit_array_spec(state, &a);
        match len {
            Some(l) => {
                arrno.push(l);
            },
            _ => {},
        };
    }
    if let Some(ref name) = d.name {    // Add this argument
        typname.push_str(&subtype[..]);
        state.add_shared(Argument{
            name: String::from(name.as_str()),
            ty: typname,
            arrno: arrno,
        });
        return Some(())
    }

    None

}

pub fn as_number(expr: &syntax::Expr) -> i32 {
    // convert an expression to a number, including using constants, or
    // not
    match expr {
        syntax::Expr::IntConst(v) => *v,
        syntax::Expr::UIntConst(v) => *v as i32,
        _ => panic!("only hard coded numbers are supported for WG size (layout specifier)"),
    }
}

pub fn visit_layout_qualifier_spec(s: &mut State, l: &syntax::LayoutQualifierSpec) -> String {
    match *l {
        syntax::LayoutQualifierSpec::Identifier(ref i, Some(ref e)) => {
            let _ = write!(s, "{} = ", i);
            visit_expr(s, &e);
            String::from(String::from(i.0.as_str()))
        }
        syntax::LayoutQualifierSpec::Identifier(ref i, None) => {
            visit_identifier(s, &i);
            String::from(String::from(i.0.as_str()))
        }
        syntax::LayoutQualifierSpec::Shared => {
            let _ = write!(s, "{}", "shared");
            String::from("")
        }
    }
}

pub fn visit_precision_qualifier(s: &mut State, p: &syntax::PrecisionQualifier) {
    match *p {
        syntax::PrecisionQualifier::High => {
            let _ = write!(s, "{}", "highp");
        }
        syntax::PrecisionQualifier::Medium => {
            let _ = write!(s, "{}", "mediump");
        }
        syntax::PrecisionQualifier::Low => {
            let _ = write!(s, "{}", "low");
        }
    }
}

pub fn visit_interpolation_qualifier(s: &mut State, i: &syntax::InterpolationQualifier) {
    match *i {
        syntax::InterpolationQualifier::Smooth => {
            let _ = write!(s, "{}", "smooth");
        }
        syntax::InterpolationQualifier::Flat => {
            let _ = write!(s, "{}", "flat");
        }
        syntax::InterpolationQualifier::NoPerspective => {
            let _ = write!(s, "{}", "noperspective");
        }
    }
}

pub fn visit_float(s: &mut State, x: f32) {
    if x.fract() == 0. {
        let _ = write!(s, "((float)({}.))", x);
    } else {
        let _ = write!(s, "((float)({}))", x);
    }
}

pub fn visit_double(s: &mut State, x: f64) {
    if x.fract() == 0. {
        let _ = write!(s, "((double)({}.))", x);
    } else {
        let _ = write!(s, "((double)({}))", x);
    }
}

pub fn visit_expr(s: &mut State, expr: &syntax::Expr) -> Option<i32> {
    match *expr {
        syntax::Expr::Variable(ref i) => visit_identifier(s, &i),
        syntax::Expr::IntConst(ref x) => {
            let _ = write!(s, "((int32_t)({}))", x);
            return Some(*x);
        }
        syntax::Expr::UIntConst(ref x) => {
            let _ = write!(s, "((uint32_t)({}))", x);
            return Some((*x) as i32);
        }
        syntax::Expr::BoolConst(ref x) => {
            let _ = write!(s, "{}", x);
        }
        syntax::Expr::FloatConst(ref x) => visit_float(s, *x),
        syntax::Expr::DoubleConst(ref x) => visit_double(s, *x),
        syntax::Expr::Unary(ref op, ref e) => {
            visit_unary_op(s, &op);
            let _ = write!(s, "{}", "(");
            visit_expr(s, &e);
            let _ = write!(s, "{}", ")");
        }
        syntax::Expr::Binary(ref op, ref l, ref r) => {
            let _ = write!(s, "{}", "(");
            visit_expr(s, &l);
            let _ = write!(s, "{}", ")");
            visit_binary_op(s, &op);
            let _ = write!(s, "{}", "(");
            visit_expr(s, &r);
            let _ = write!(s, "{}", ")");
        }
        syntax::Expr::Ternary(ref c, ref ss, ref e) => {
            visit_expr(s, &c);
            let _ = write!(s, "{}", " ? ");
            visit_expr(s, &ss);
            let _ = write!(s, "{}", " : ");
            visit_expr(s, &e);
        }
        syntax::Expr::Assignment(ref v, ref op, ref e) => {
            visit_expr(s, &v);
            let _ = write!(s, "{}", " ");
            visit_assignment_op(s, &op);
            let _ = write!(s, "{}", " ");
            visit_expr(s, &e);
        }
        syntax::Expr::Bracket(ref e, ref a) => {
            visit_expr(s, &e);
            visit_array_spec(s, &a);
        }
        syntax::Expr::FunCall(ref fun, ref args) => {
            // we treat function calls to vector constructors specially since
            // they need to be remade inte mae calls
            match fun {
                syntax::FunIdentifier::Identifier(ref n) => {
                    let _ = match &n.0[..] {
                        "vec2" => write!(s, "{}", "make_vec2"),
                        "vec3" => write!(s, "{}", "make_vec3"),
                        "vec4" => write!(s, "{}", "make_vec4"),
                        "ivec2" => write!(s, "{}", "make_ivec2"),
                        "ivec3" => write!(s, "{}", "make_ivec3"),
                        "ivec4" => write!(s, "{}", "make_ivec4"),
                        "uvec2" => write!(s, "{}", "make_uvec2"),
                        "uvec3" => write!(s, "{}", "make_uvec3"),
                        "uvec4" => write!(s, "{}", "make_uvec4"),
                        "bvec2" => write!(s, "{}", "make_bvec2"),
                        "bvec3" => write!(s, "{}", "make_bvec3"),
                        "bvec4" => write!(s, "{}", "make_bvec4"),
                        "mat2" => write!(s, "{}", "make_mat2"),
                        "mat3" => write!(s, "{}", "make_mat3"),
                        "mat4" => write!(s, "{}", "make_mat4"),
                        "barrier" => write!(s, "{}", "this->barrier"),
                        _ =>  write!(s, "{}", n.0),
                    };
                },
                syntax::FunIdentifier::Expr(ref e) => {
                    visit_expr(s, &*e);
                }
            };
            //visit_function_identifier(s, &fun);

            let _ = write!(s, "{}", "(");

            let is_atomic = match fun {
                syntax::FunIdentifier::Identifier(ref n) => {
                    match &n.0[..] {
                        "atomicAdd" => true,
                        "atomicMin" => true,
                        "atomicMax" => true,
                        "atomicAnd" => true,
                        "atomicOr" => true,
                        "atomicXor" => true,
                        "atomicExchange" => true,
                        "atomicCompSwap" => true,
                        _ => false,
                    }
                },
                _ =>  false,
            };

            if !args.is_empty() {
                let mut args_iter = args.iter();
                let first = args_iter.next().unwrap();
                if is_atomic {
                    let _ = write!(s, "{}", "&(");
                    visit_expr(s, first);
                    let _ = write!(s, "{}", ")");
                } else {
                    visit_expr(s, first);
                }

                for e in args_iter {
                    let _ = write!(s, "{}", ", ");
                    visit_expr(s, e);
                }
            }

            let _ = write!(s, "{}", ")");
        }
        syntax::Expr::Dot(ref e, ref i) => {
            let _ = write!(s, "{}", "(");
            visit_expr(s, &e);
            let _ = write!(s, "{}", ")");
            let _ = write!(s, "{}", ".");
            visit_identifier(s, &i);
        }
        syntax::Expr::PostInc(ref e) => {
            visit_expr(s, &e);
            let _ = write!(s, "{}", "++");
        }
        syntax::Expr::PostDec(ref e) => {
            visit_expr(s, &e);
            let _ = write!(s, "{}", "--");
        }
        syntax::Expr::Comma(ref a, ref b) => {
            visit_expr(s, &a);
            let _ = write!(s, "{}", ", ");
            visit_expr(s, &b);
        }
    }
    return None;
}

pub fn visit_path(s: &mut State, path: &syntax::Path) {
    match path {
        syntax::Path::Absolute(ss) => {
            let _ = write!(s, "<{}>", ss);
        }
        syntax::Path::Relative(ss) => {
            let _ = write!(s, "\"{}\"", ss);
        }
    }
}

pub fn visit_unary_op(s: &mut State, op: &syntax::UnaryOp) {
    match *op {
        syntax::UnaryOp::Inc => {
            let _ = write!(s, "{}", "++");
        }
        syntax::UnaryOp::Dec => {
            let _ = write!(s, "{}", "--");
        }
        syntax::UnaryOp::Add => {
            let _ = write!(s, "{}", "+");
        }
        syntax::UnaryOp::Minus => {
            let _ = write!(s, "{}", "-");
        }
        syntax::UnaryOp::Not => {
            let _ = write!(s, "{}", "!");
        }
        syntax::UnaryOp::Complement => {
            let _ = write!(s, "{}", "~");
        }
    }
}

pub fn visit_binary_op(s: &mut State, op: &syntax::BinaryOp) {
    match *op {
        syntax::BinaryOp::Or => {
            let _ = write!(s, "{}", "||");
        }
        syntax::BinaryOp::Xor => {
            let _ = write!(s, "{}", "^^");
        }
        syntax::BinaryOp::And => {
            let _ = write!(s, "{}", "&&");
        }
        syntax::BinaryOp::BitOr => {
            let _ = write!(s, "{}", "|");
        }
        syntax::BinaryOp::BitXor => {
            let _ = write!(s, "{}", "^");
        }
        syntax::BinaryOp::BitAnd => {
            let _ = write!(s, "{}", "&");
        }
        syntax::BinaryOp::Equal => {
            let _ = write!(s, "{}", "==");
        }
        syntax::BinaryOp::NonEqual => {
            let _ = write!(s, "{}", "!=");
        }
        syntax::BinaryOp::LT => {
            let _ = write!(s, "{}", "<");
        }
        syntax::BinaryOp::GT => {
            let _ = write!(s, "{}", ">");
        }
        syntax::BinaryOp::LTE => {
            let _ = write!(s, "{}", "<=");
        }
        syntax::BinaryOp::GTE => {
            let _ = write!(s, "{}", ">=");
        }
        syntax::BinaryOp::LShift => {
            let _ = write!(s, "{}", "<<");
        }
        syntax::BinaryOp::RShift => {
            let _ = write!(s, "{}", ">>");
        }
        syntax::BinaryOp::Add => {
            let _ = write!(s, "{}", "+");
        }
        syntax::BinaryOp::Sub => {
            let _ = write!(s, "{}", "-");
        }
        syntax::BinaryOp::Mult => {
            let _ = write!(s, "{}", "*");
        }
        syntax::BinaryOp::Div => {
            let _ = write!(s, "{}", "/");
        }
        syntax::BinaryOp::Mod => {
            let _ = write!(s, "{}", "%");
        }
    }
}

pub fn visit_assignment_op(s: &mut State, op: &syntax::AssignmentOp) {
    match *op {
        syntax::AssignmentOp::Equal => {
            let _ = write!(s, "{}", "=");
        }
        syntax::AssignmentOp::Mult => {
            let _ = write!(s, "{}", "*=");
        }
        syntax::AssignmentOp::Div => {
            let _ = write!(s, "{}", "/=");
        }
        syntax::AssignmentOp::Mod => {
            let _ = write!(s, "{}", "%=");
        }
        syntax::AssignmentOp::Add => {
            let _ = write!(s, "{}", "+=");
        }
        syntax::AssignmentOp::Sub => {
            let _ = write!(s, "{}", "-=");
        }
        syntax::AssignmentOp::LShift => {
            let _ = write!(s, "{}", "<<=");
        }
        syntax::AssignmentOp::RShift => {
            let _ = write!(s, "{}", ">>=");
        }
        syntax::AssignmentOp::And => {
            let _ = write!(s, "{}", "&=");
        }
        syntax::AssignmentOp::Xor => {
            let _ = write!(s, "{}", "^=");
        }
        syntax::AssignmentOp::Or => {
            let _ = write!(s, "{}", "|=");
        }
    }
}

pub fn visit_function_identifier(s: &mut State, i: &syntax::FunIdentifier) {
    match *i {
        syntax::FunIdentifier::Identifier(ref n) => visit_identifier(s, &n),
        syntax::FunIdentifier::Expr(ref e) => {
            visit_expr(s, &*e);
        }
    }
}

pub fn visit_declaration(s: &mut State, d: &syntax::Declaration, global: bool) {
    match *d {
        syntax::Declaration::FunctionPrototype(ref proto) => {
            visit_function_prototype(s, &proto);
            let _ = write!(s, "{}", ";\n");
        }
        syntax::Declaration::InitDeclaratorList(ref list) => {
            // Global struct goes here
            // So does buffer inputs that are not part of blos

            s.push_output(Output::None);
            match maybe_handle_global_buffer(s, &list) {
                Some(..) => {
                    s.pop_output();
                    return
                },
                _ => (),
            }
            s.pop_output();

            // Handle declarations of shared varaibles that we need to parse out
            s.push_output(Output::None);
            match maybe_handle_shared(s, &list) {
                Some(..) => {
                    s.pop_output();
                    return
                },
                _ => (),
            }
            s.pop_output();

            //s.push_output(Output::None);

            if global {
                s.push_output(Output::None);
            }
            visit_init_declarator_list(s, &list);
            if global{
               s.pop_output();
            }
            let _ = write!(s, "{}", ";\n");
            //s.pop_output();
        }
        syntax::Declaration::Precision(ref qual, ref ty) => {
            visit_precision_qualifier(s, &qual);
            visit_type_specifier(s, &ty);
            let _ = write!(s, "{}", ";\n");
        }
        syntax::Declaration::Block(ref block) => {
            visit_block(s, &block);
            //let _ = write!(s, "{}", ";\n");
        }
        syntax::Declaration::Global(ref qual, ref identifiers) => {
            // The IN specifier goes straing here
            match maybe_handle_wg(s, &qual) {
                Some(..) => return,
                _ => (),
            }

            visit_type_qualifier(s, &qual);

            if !identifiers.is_empty() {
                let mut iter = identifiers.iter();
                let first = iter.next().unwrap();
                visit_identifier(s, first);

                for identifier in iter {
                    let _ = write!(s, ", {}", identifier);
                }
            }

            let _ = write!(s, "{}", ";\n");
        }
    }
}

pub fn visit_function_prototype(s: &mut State, fp: &syntax::FunctionPrototype) {
    visit_fully_specified_type(s, &fp.ty);
    let _ = write!(s, "{}", " ");
    visit_identifier(s, &fp.name);

    let _ = write!(s, "{}", "(");

    if !fp.parameters.is_empty() {
        let mut iter = fp.parameters.iter();
        let first = iter.next().unwrap();
        visit_function_parameter_declaration(s, first);

        for param in iter {
            let _ = write!(s, "{}", ", ");
            visit_function_parameter_declaration(s, param);
        }
    }

    let _ = write!(s, "{}", ")");
}
pub fn visit_function_parameter_declaration(
    s: &mut State,
    p: &syntax::FunctionParameterDeclaration,
) {
    match *p {
        syntax::FunctionParameterDeclaration::Named(ref qual, ref fpd) => {
            if let Some(ref q) = *qual {
                visit_type_qualifier(s, q);
                let _ = write!(s, "{}", " ");
            }

            visit_function_parameter_declarator(s, fpd);
        }
        syntax::FunctionParameterDeclaration::Unnamed(ref qual, ref ty) => {
            if let Some(ref q) = *qual {
                visit_type_qualifier(s, q);
                let _ = write!(s, "{}", " ");
            }

            visit_type_specifier(s, ty);
        }
    }
}

pub fn visit_function_parameter_declarator(s: &mut State, p: &syntax::FunctionParameterDeclarator) {
    visit_type_specifier(s, &p.ty);
    let _ = write!(s, "{}", " ");
    visit_arrayed_identifier(s, &p.ident);
}

pub fn visit_init_declarator_list(s: &mut State, i: &syntax::InitDeclaratorList) {
    visit_single_declaration(s, &i.head);

    for decl in &i.tail {
        let _ = write!(s, "{}", ", ");
        visit_single_declaration_no_type(s, decl);
    }
}

pub fn visit_single_declaration(s: &mut State, d: &syntax::SingleDeclaration) {
    visit_fully_specified_type(s, &d.ty);

    if let Some(ref name) = d.name {
        let _ = write!(s, "{}", " ");
        visit_identifier(s, name);
    }

    if let Some(ref arr_spec) = d.array_specifier {
        visit_array_spec(s, arr_spec);
    }

    if let Some(ref initializer) = d.initializer {
        let _ = write!(s, "{}", " = ");
        visit_initializer(s, initializer);
    }
}

pub fn visit_single_declaration_no_type(s: &mut State, d: &syntax::SingleDeclarationNoType) {
    visit_arrayed_identifier(s, &d.ident);

    if let Some(ref initializer) = d.initializer {
        let _ = write!(s, "{}", " = ");
        visit_initializer(s, initializer);
    }
}

pub fn visit_initializer(s: &mut State, i: &syntax::Initializer) {
    match *i {
        syntax::Initializer::Simple(ref e) => {
            visit_expr(s, e);
        }
        syntax::Initializer::List(ref list) => {
            let mut iter = list.0.iter();
            let first = iter.next().unwrap();

            let _ = write!(s, "{}", "{ ");
            visit_initializer(s, first);

            for ini in iter {
                let _ = write!(s, "{}", ", ");
                visit_initializer(s, ini);
            }

            let _ = write!(s, "{}", " }");
        }
    }
}

pub fn visit_block(s: &mut State, b: &syntax::Block) {
    s.push_output(Output::None);
    // Thus far we assume that this one is visided only when specifying
    // buffers / unifomrs that we should be able to access
    visit_type_qualifier(s, &b.qualifier);
    if non_accepted_layout_spec(s, &b.qualifier) {
        panic!("we only support std430 layouts!");
    }
    //let _ = write!(s, "{}", " ");
    visit_identifier(s, &b.name);
    //let _ = write!(s, "{}", " {");

    for field in &b.fields {
        let arg = visit_struct_field(s, field);
        s.add_arg(arg);
        //let _ = write!(s, "{}", "\n");
    }
    //let _ = write!(s, "{}", "}");

    if let Some(ref ident) = b.identifier {
        visit_arrayed_identifier(s, ident);
    }
    s.pop_output();
}

pub fn visit_function_definition(s: &mut State, fd: &syntax::FunctionDefinition) {
    visit_function_prototype(s, &fd.prototype);
    let _ = write!(s, "{}", " ");
    visit_compound_statement(s, &fd.statement);
}

pub fn visit_compound_statement(s: &mut State, cst: &syntax::CompoundStatement) {
    let _ = write!(s, "{}", "{\n");

    for st in &cst.statement_list {
        visit_statement(s, st);
    }

    let _ = write!(s, "{}", "}\n");
}

pub fn visit_statement(s: &mut State, st: &syntax::Statement) {
    match *st {
        syntax::Statement::Compound(ref cst) => visit_compound_statement(s, cst),
        syntax::Statement::Simple(ref sst) => visit_simple_statement(s, sst),
    }
}

pub fn visit_simple_statement(s: &mut State, sst: &syntax::SimpleStatement) {
    match *sst {
        syntax::SimpleStatement::Declaration(ref d) => visit_declaration(s, d, false),
        syntax::SimpleStatement::Expression(ref e) => visit_expression_statement(s, e),
        syntax::SimpleStatement::Selection(ref ss) => visit_selection_statement(s, ss),
        syntax::SimpleStatement::Switch(ref ss) => visit_switch_statement(s, ss),
        syntax::SimpleStatement::CaseLabel(ref cl) => visit_case_label(s, cl),
        syntax::SimpleStatement::Iteration(ref i) => visit_iteration_statement(s, i),
        syntax::SimpleStatement::Jump(ref j) => visit_jump_statement(s, j),
    }
}

pub fn visit_expression_statement(s: &mut State, est: &syntax::ExprStatement) {
    if let Some(ref e) = *est {
        visit_expr(s, e);
    }

    let _ = write!(s, "{}", ";\n");
}

pub fn visit_selection_statement(s: &mut State, sst: &syntax::SelectionStatement) {
    let _ = write!(s, "{}", "if (");
    visit_expr(s, &sst.cond);
    let _ = write!(s, "{}", ") {\n");
    visit_selection_rest_statement(s, &sst.rest);
}

pub fn visit_selection_rest_statement(s: &mut State, sst: &syntax::SelectionRestStatement) {
    match *sst {
        syntax::SelectionRestStatement::Statement(ref if_st) => {
            visit_statement(s, if_st);
            let _ = write!(s, "{}", "}\n");
        }
        syntax::SelectionRestStatement::Else(ref if_st, ref else_st) => {
            visit_statement(s, if_st);
            let _ = write!(s, "{}", "} else ");
            visit_statement(s, else_st);
        }
    }
}

pub fn visit_switch_statement(s: &mut State, sst: &syntax::SwitchStatement) {
    let _ = write!(s, "{}", "switch (");
    visit_expr(s, &sst.head);
    let _ = write!(s, "{}", ") {\n");

    for st in &sst.body {
        visit_statement(s, st);
    }

    let _ = write!(s, "{}", "}\n");
}

pub fn visit_case_label(s: &mut State, cl: &syntax::CaseLabel) {
    match *cl {
        syntax::CaseLabel::Case(ref e) => {
            let _ = write!(s, "{}", "case ");
            visit_expr(s, e);
            let _ = write!(s, "{}", ":\n");
        }
        syntax::CaseLabel::Def => {
            let _ = write!(s, "{}", "default:\n");
        }
    }
}

pub fn visit_iteration_statement(s: &mut State, ist: &syntax::IterationStatement) {
    match *ist {
        syntax::IterationStatement::While(ref cond, ref body) => {
            let _ = write!(s, "{}", "while (");
            visit_condition(s, cond);
            let _ = write!(s, "{}", ") ");
            visit_statement(s, body);
        }
        syntax::IterationStatement::DoWhile(ref body, ref cond) => {
            let _ = write!(s, "{}", "do ");
            visit_statement(s, body);
            let _ = write!(s, "{}", " while (");
            visit_expr(s, cond);
            let _ = write!(s, "{}", ")\n");
        }
        syntax::IterationStatement::For(ref init, ref rest, ref body) => {
            let _ = write!(s, "{}", "for (");
            visit_for_init_statement(s, init);
            visit_for_rest_statement(s, rest);
            let _ = write!(s, "{}", ") ");
            visit_statement(s, body);
        }
    }
}

pub fn visit_condition(s: &mut State, c: &syntax::Condition) {
    match *c {
        syntax::Condition::Expr(ref e) => {
            visit_expr(s, e);
        }
        syntax::Condition::Assignment(ref ty, ref name, ref initializer) => {
            visit_fully_specified_type(s, ty);
            let _ = write!(s, "{}", " ");
            visit_identifier(s, name);
            let _ = write!(s, "{}", " = ");
            visit_initializer(s, initializer);
        }
    }
}

pub fn visit_for_init_statement(s: &mut State, i: &syntax::ForInitStatement) {
    match *i {
        syntax::ForInitStatement::Expression(ref expr) => {
            if let Some(ref e) = *expr {
                visit_expr(s, e);
            }
        }
        syntax::ForInitStatement::Declaration(ref d) => visit_declaration(s, d, false),
    }
}

pub fn visit_for_rest_statement(s: &mut State, r: &syntax::ForRestStatement) {
    if let Some(ref cond) = r.condition {
        visit_condition(s, cond);
    }

    let _ = write!(s, "{}", "; ");

    if let Some(ref e) = r.post_expr {
        visit_expr(s, e);
    }
}

pub fn visit_jump_statement(s: &mut State, j: &syntax::JumpStatement) {
    match *j {
        syntax::JumpStatement::Continue => {
            let _ = write!(s, "{}", "continue;\n");
        }
        syntax::JumpStatement::Break => {
            let _ = write!(s, "{}", "break;\n");
        }
        syntax::JumpStatement::Discard => {
            let _ = write!(s, "{}", "discard;\n");
        }
        syntax::JumpStatement::Return(ref e) => {
            let _ = write!(s, "{}", "return ");
            if let Some(e) = e {
                visit_expr(s, e);
            }
            let _ = write!(s, "{}", ";\n");
        }
    }
}

pub fn visit_preprocessor(s: &mut State, pp: &syntax::Preprocessor) {
    match *pp {
        syntax::Preprocessor::Define(ref pd) => visit_preprocessor_define(s, pd),
        syntax::Preprocessor::Else => visit_preprocessor_else(s),
        syntax::Preprocessor::ElseIf(ref pei) => visit_preprocessor_elseif(s, pei),
        syntax::Preprocessor::EndIf => visit_preprocessor_endif(s),
        syntax::Preprocessor::Error(ref pe) => visit_preprocessor_error(s, pe),
        syntax::Preprocessor::If(ref pi) => visit_preprocessor_if(s, pi),
        syntax::Preprocessor::IfDef(ref pid) => visit_preprocessor_ifdef(s, pid),
        syntax::Preprocessor::IfNDef(ref pind) => visit_preprocessor_ifndef(s, pind),
        syntax::Preprocessor::Include(..) => {
            panic!("visited and #include that had not been handled")
        }
        syntax::Preprocessor::Line(ref pl) => visit_preprocessor_line(s, pl),
        syntax::Preprocessor::Pragma(ref pp) => visit_preprocessor_pragma(s, pp),
        syntax::Preprocessor::Undef(ref pu) => visit_preprocessor_undef(s, pu),
        syntax::Preprocessor::Version(ref pv) => visit_preprocessor_version(s, pv),
        syntax::Preprocessor::Extension(..) => (),
    }
}

pub fn visit_preprocessor_define(s: &mut State, pd: &syntax::PreprocessorDefine) {
    match *pd {
        syntax::PreprocessorDefine::ObjectLike {
            ref ident,
            ref value,
        } => {
            let _ = write!(s, "#define {} {}\n", ident, value);
        }

        syntax::PreprocessorDefine::FunctionLike {
            ref ident,
            ref args,
            ref value,
        } => {
            let _ = write!(s, "#define {}(", ident);

            if !args.is_empty() {
                let _ = write!(s, "{}", &args[0]);

                for arg in &args[1..args.len()] {
                    let _ = write!(s, ", {}", arg);
                }
            }

            let _ = write!(s, ") {}\n", value);
        }
    }
}

pub fn visit_preprocessor_else(s: &mut State) {
    let _ = write!(s, "{}", "#else\n");
}

pub fn visit_preprocessor_elseif(s: &mut State, pei: &syntax::PreprocessorElseIf) {
    let _ = write!(s, "#elseif {}\n", pei.condition);
}

pub fn visit_preprocessor_error(s: &mut State, pe: &syntax::PreprocessorError) {
    let _ = writeln!(s, "#error {}", pe.message);
}

pub fn visit_preprocessor_endif(s: &mut State) {
    let _ = write!(s, "{}", "#endif\n");
}

pub fn visit_preprocessor_if(s: &mut State, pi: &syntax::PreprocessorIf) {
    let _ = write!(s, "#if {}\n", pi.condition);
}

pub fn visit_preprocessor_ifdef(s: &mut State, pid: &syntax::PreprocessorIfDef) {
    let _ = write!(s, "{}", "#ifdef ");
    visit_identifier(s, &pid.ident);
    let _ = write!(s, "{}", "\n");
}

pub fn visit_preprocessor_ifndef(s: &mut State, pind: &syntax::PreprocessorIfNDef) {
    let _ = write!(s, "{}", "#ifndef ");
    visit_identifier(s, &pind.ident);
    let _ = write!(s, "{}", "\n");
}

pub fn visit_preprocessor_line(s: &mut State, pl: &syntax::PreprocessorLine) {
    let _ = write!(s, "#line {}", pl.line);
    if let Some(source_string_number) = pl.source_string_number {
        let _ = write!(s, " {}", source_string_number);
    }
    let _ = write!(s, "{}", "\n");
}

pub fn visit_preprocessor_pragma(s: &mut State, pp: &syntax::PreprocessorPragma) {
    let _ = writeln!(s, "#pragma {}", pp.command);
}

pub fn visit_preprocessor_undef(s: &mut State, pud: &syntax::PreprocessorUndef) {
    let _ = write!(s, "{}", "#undef ");
    visit_identifier(s, &pud.name);
    let _ = write!(s, "{}", "\n");
}

pub fn visit_preprocessor_version(_s: &mut State, pv: &syntax::PreprocessorVersion) {
    if pv.version != 450 {
        panic!("only support version 450")
    }
}

pub fn visit_external_declaration(s: &mut State, ed: &syntax::ExternalDeclaration) {
    match *ed {
        syntax::ExternalDeclaration::Preprocessor(ref pp) => visit_preprocessor(s, pp),
        syntax::ExternalDeclaration::FunctionDefinition(ref fd) => visit_function_definition(s, fd),
        syntax::ExternalDeclaration::Declaration(ref d) => visit_declaration(s, d, true),
    }
}

pub fn visit_translation_unit(s: &mut State, tu: &syntax::TranslationUnit) {
    for ed in &(tu.0).0 {
        visit_external_declaration(s, ed);
    }
}

pub fn translate(file: String) -> String {
    let mut contents = String::new();
    std::fs::File::open(&file)
        .unwrap()
        .read_to_string(&mut contents)
        .unwrap();

    let contents = contents.replace("\r\n", "\n");

    let tu = syntax::TranslationUnit::parse(contents).unwrap();
    let mut state = State {
        body: Part::new(),
        arguments: Vec::new(),
        shared: Vec::new(),
        structs: Vec::new(),
        output: Output::Body,
        last_output: Vec::new(),
        wg_size: [0, 1, 1],
    };

    visit_translation_unit(&mut state, &tu);

    state.write_json()
}
