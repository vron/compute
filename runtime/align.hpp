#include <cerrno>

bool kernel::ensure_alignments(cpt_data d) {
	(void)d;
	if(sizeof(Bool) != 4) { return this->set_error(EINVAL, "static check failed: sizeof(Bool) != 4"); };
	if(alignof(Bool) != 4) { return this->set_error(EINVAL, "static check failed: alignof(Bool) != 4"); };
	if(sizeof(float) != 4) { return this->set_error(EINVAL, "static check failed: sizeof(float) != 4"); };
	if(alignof(float) != 4) { return this->set_error(EINVAL, "static check failed: alignof(float) != 4"); };
	if(sizeof(int32_t) != 4) { return this->set_error(EINVAL, "static check failed: sizeof(int32_t) != 4"); };
	if(alignof(int32_t) != 4) { return this->set_error(EINVAL, "static check failed: alignof(int32_t) != 4"); };
	if(sizeof(image2Drgba32f) != 16) { return this->set_error(EINVAL, "static check failed: sizeof(image2Drgba32f) != 16"); };
	if(alignof(image2Drgba32f) != 8) { return this->set_error(EINVAL, "static check failed: alignof(image2Drgba32f) != 8"); };
	if(sizeof(ivec2) != 8) { return this->set_error(EINVAL, "static check failed: sizeof(ivec2) != 8"); };
	if(alignof(ivec2) != 8) { return this->set_error(EINVAL, "static check failed: alignof(ivec2) != 8"); };
	if(sizeof(ivec3) != 16) { return this->set_error(EINVAL, "static check failed: sizeof(ivec3) != 16"); };
	if(alignof(ivec3) != 16) { return this->set_error(EINVAL, "static check failed: alignof(ivec3) != 16"); };
	if(sizeof(ivec4) != 16) { return this->set_error(EINVAL, "static check failed: sizeof(ivec4) != 16"); };
	if(alignof(ivec4) != 16) { return this->set_error(EINVAL, "static check failed: alignof(ivec4) != 16"); };
	if(sizeof(mat2) != 16) { return this->set_error(EINVAL, "static check failed: sizeof(mat2) != 16"); };
	if(alignof(mat2) != 8) { return this->set_error(EINVAL, "static check failed: alignof(mat2) != 8"); };
	if(sizeof(mat3) != 48) { return this->set_error(EINVAL, "static check failed: sizeof(mat3) != 48"); };
	if(alignof(mat3) != 16) { return this->set_error(EINVAL, "static check failed: alignof(mat3) != 16"); };
	if(sizeof(mat4) != 64) { return this->set_error(EINVAL, "static check failed: sizeof(mat4) != 64"); };
	if(alignof(mat4) != 16) { return this->set_error(EINVAL, "static check failed: alignof(mat4) != 16"); };
	if(sizeof(uint32_t) != 4) { return this->set_error(EINVAL, "static check failed: sizeof(uint32_t) != 4"); };
	if(alignof(uint32_t) != 4) { return this->set_error(EINVAL, "static check failed: alignof(uint32_t) != 4"); };
	if(sizeof(uvec2) != 8) { return this->set_error(EINVAL, "static check failed: sizeof(uvec2) != 8"); };
	if(alignof(uvec2) != 8) { return this->set_error(EINVAL, "static check failed: alignof(uvec2) != 8"); };
	if(sizeof(uvec3) != 16) { return this->set_error(EINVAL, "static check failed: sizeof(uvec3) != 16"); };
	if(alignof(uvec3) != 16) { return this->set_error(EINVAL, "static check failed: alignof(uvec3) != 16"); };
	if(sizeof(uvec4) != 16) { return this->set_error(EINVAL, "static check failed: sizeof(uvec4) != 16"); };
	if(alignof(uvec4) != 16) { return this->set_error(EINVAL, "static check failed: alignof(uvec4) != 16"); };
	if(sizeof(vec2) != 8) { return this->set_error(EINVAL, "static check failed: sizeof(vec2) != 8"); };
	if(alignof(vec2) != 8) { return this->set_error(EINVAL, "static check failed: alignof(vec2) != 8"); };
	if(sizeof(vec3) != 16) { return this->set_error(EINVAL, "static check failed: sizeof(vec3) != 16"); };
	if(alignof(vec3) != 16) { return this->set_error(EINVAL, "static check failed: alignof(vec3) != 16"); };
	if(sizeof(vec4) != 16) { return this->set_error(EINVAL, "static check failed: sizeof(vec4) != 16"); };
	if(alignof(vec4) != 16) { return this->set_error(EINVAL, "static check failed: alignof(vec4) != 16"); };
	if((((uintptr_t)(const void *)(d.img.data)) % (8)) != 0) { return this->set_error(EINVAL, "the argument d.img.data provided was not aligned to a 8 byte boundary as required"); };

	return true;
}
