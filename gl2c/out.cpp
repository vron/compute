#define _cpt_WG_SIZE_X 8
#define _cpt_WG_SIZE_X 8
#define _cpt_WG_SIZE_X 1

struct sample_comp {
uvec3 gl_GlobalInvocationID;
uvec3 gl_LocalInvocationID;
image2DRGBA32F img;
// image2DXXX img;
void main() {
 vec4 pixel = make_vec4(0.f, 0.f, 0.f, 1.f);
 uvec2 pixel_coords;
 pixel_coords = (gl_GlobalInvocationID).sel(_ind_X, _ind_Y);
 pixel_coords *= 8;
 pixel_coords += (gl_LocalInvocationID).sel(_ind_X, _ind_Y);
 auto _c2_ = (((pixel_coords).sel(_ind_X))%(2))==(0);
 {
  pixel = if_then_else(_c2_,make_vec4(1.f, 1.f, 1.f, 1.f),pixel);
 }
 imageStore(img, make_ivec2(pixel_coords), pixel);
}
}
;

