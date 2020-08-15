#pragma once

#define _cpt_WG_SIZE_X 64
#define _cpt_WG_SIZE_Y 1
#define _cpt_WG_SIZE_Z 1
#define _cpt_WG_SIZE (_cpt_WG_SIZE_Z * _cpt_WG_SIZE_Y * _cpt_WG_SIZE_X)

#include <cmath>
#include "../types/types.hpp"
#include "usertypes.hpp"
#include "../routines/routines.hpp"

class  shared_data_t {
public:
	cog_res* shared_data;

} ;

struct shader {
	uvec3 gl_NumWorkGroups;
	uvec3 gl_WorkGroupSize;
	uvec3 gl_WorkGroupID;
	uvec3 gl_LocalInvocationID;
	uvec3 gl_GlobalInvocationID;
	uint32_t gl_LocalInvocationIndex;
	Invocation<struct shader*>  *invocation;

	mat3 transform;
	polygon* polygons;

	vec2* cogs;

	cog_res* shared_data;


	shader() {};

;
;
;
float area_tri(triangle t) {
float a = ((((t).vertices[((int32_t)(1))][((int32_t)(0))])-((t).vertices[((int32_t)(0))][((int32_t)(0))]))*(((t).vertices[((int32_t)(2))][((int32_t)(1))])-((t).vertices[((int32_t)(0))][((int32_t)(1))])))-((((t).vertices[((int32_t)(2))][((int32_t)(0))])-((t).vertices[((int32_t)(0))][((int32_t)(0))]))*(((t).vertices[((int32_t)(1))][((int32_t)(1))])-((t).vertices[((int32_t)(0))][((int32_t)(1))])));
if ((a)<(((float)(0.)))) {
{
a *= -(((float)(0.5)));
}
} else {
a *= ((float)(0.5));
}
return a;
}
vec2 cog_tri(triangle t) {
return make_vec2(((((t).vertices[((int32_t)(1))][((int32_t)(0))])+((t).vertices[((int32_t)(2))][((int32_t)(0))]))+((t).vertices[((int32_t)(0))][((int32_t)(0))]))/(((float)(3.))), ((((t).vertices[((int32_t)(1))][((int32_t)(1))])+((t).vertices[((int32_t)(2))][((int32_t)(1))]))+((t).vertices[((int32_t)(0))][((int32_t)(1))]))/(((float)(3.))));
}
cog_res tri(triangle t) {
for (int32_t i = ((int32_t)(0));
(i)<(((int32_t)(3))); i++) {
(t).vertices[i] = ((transform)*(make_vec3((t).vertices[i], ((float)(1.))))).xy;
}
cog_res r;
(r).area = area_tri(t);
(r).cog = cog_tri(t);
return r;
}
cog_res cog_poly(polygon p) {
float area = ((float)(0.));
vec2 cog = make_vec2(((int32_t)(0)), ((int32_t)(0)));
for (int32_t i = ((int32_t)(0));
(i)<(((int32_t)(64))); i++) {
cog_res tr = tri((p).triangles[i]);
area += (tr).area;
cog += ((tr).area)*((tr).cog);
}
cog /= area;
cog_res r;
(r).area = area;
(r).cog = cog;
return r;
}
void main() {
uint32_t base_index = ((gl_WorkGroupID).x)*((gl_WorkGroupSize).x);
uint32_t local_index = (gl_LocalInvocationID).x;
uint32_t index = (base_index)+(local_index);
cog_res my_res = cog_poly(polygons[index]);
shared_data[local_index] = my_res;
this->barrier();
if ((local_index)==(((int32_t)(0)))) {
{
(my_res).cog *= (my_res).area;
for (int32_t i = ((int32_t)(1));
(i)<(((int32_t)(64))); i++) {
cog_res fr = shared_data[i];
(my_res).area += (fr).area;
(my_res).cog += ((fr).area)*((fr).cog);
}
(my_res).cog /= (my_res).area;
cogs[(gl_WorkGroupID).x] = (my_res).cog;
}
}
}



	void set_data(cpt_data d) {
	auto me = this;
	from_api(&(me->transform), d.transform);// mat bind
	me->polygons = (polygon*)d.polygons;
	me->cogs = (vec2*)d.cogs;
	return;
	}

	void barrier();

	static shared_data_t* create_shared_data() {
		shared_data_t *sd = new shared_data_t();
		sd->shared_data = (cog_res*)malloc(64*sizeof(cog_res));

		return sd;
	}

	static void free_shared_data(shared_data_t *sd) {
		free(sd->shared_data);

		delete sd;
	}

	void set_shared_data(shared_data_t *sd) {
		(void)sd;
		this->shared_data = sd->shared_data;

	}
};
