#pragma once
// Code generated DO NOT EDIT
#include "../types/types.hpp"
#include "./shared.h"



struct cog_res;
struct triangle;
struct polygon;

struct cog_res {// size=16 alignment=8
	vec2 cog;  // size=8 alignment=8 offset=0
	float area;  // size=4 alignment=4 offset=8
	 cog_res () {};
	cog_res(vec2 cog, float area) : cog(cog), area(area) {};
};

struct triangle {// size=24 alignment=8
	vec2 vertices[3];  // size=24 alignment=8 offset=0
	 triangle () {};
	triangle(vec2 vertices[3]) : vertices{vertices[0], vertices[1], vertices[2]} {};
};

struct polygon {// size=1536 alignment=8
	triangle triangles[64];  // size=1536 alignment=8 offset=0
	 polygon () {};
	polygon( struct triangle triangles[64]) : triangles{triangles[0], triangles[1], triangles[2], triangles[3], triangles[4], triangles[5], triangles[6], triangles[7], triangles[8], triangles[9], triangles[10], triangles[11], triangles[12], triangles[13], triangles[14], triangles[15], triangles[16], triangles[17], triangles[18], triangles[19], triangles[20], triangles[21], triangles[22], triangles[23], triangles[24], triangles[25], triangles[26], triangles[27], triangles[28], triangles[29], triangles[30], triangles[31], triangles[32], triangles[33], triangles[34], triangles[35], triangles[36], triangles[37], triangles[38], triangles[39], triangles[40], triangles[41], triangles[42], triangles[43], triangles[44], triangles[45], triangles[46], triangles[47], triangles[48], triangles[49], triangles[50], triangles[51], triangles[52], triangles[53], triangles[54], triangles[55], triangles[56], triangles[57], triangles[58], triangles[59], triangles[60], triangles[61], triangles[62], triangles[63]} {};
};

void from_api(triangle *me, cpt_triangle d) {
	me->vertices[0][0] = d.vertices[0];
	me->vertices[0][1] = d.vertices[1];
	me->vertices[1][0] = d.vertices[2];
	me->vertices[1][1] = d.vertices[3];
	me->vertices[2][0] = d.vertices[4];
	me->vertices[2][1] = d.vertices[5];
};

void from_api(polygon *me, cpt_polygon d) {
	from_api(&(me->triangles[0]), d.triangles[0]);
	from_api(&(me->triangles[1]), d.triangles[1]);
	from_api(&(me->triangles[2]), d.triangles[2]);
	from_api(&(me->triangles[3]), d.triangles[3]);
	from_api(&(me->triangles[4]), d.triangles[4]);
	from_api(&(me->triangles[5]), d.triangles[5]);
	from_api(&(me->triangles[6]), d.triangles[6]);
	from_api(&(me->triangles[7]), d.triangles[7]);
	from_api(&(me->triangles[8]), d.triangles[8]);
	from_api(&(me->triangles[9]), d.triangles[9]);
	from_api(&(me->triangles[10]), d.triangles[10]);
	from_api(&(me->triangles[11]), d.triangles[11]);
	from_api(&(me->triangles[12]), d.triangles[12]);
	from_api(&(me->triangles[13]), d.triangles[13]);
	from_api(&(me->triangles[14]), d.triangles[14]);
	from_api(&(me->triangles[15]), d.triangles[15]);
	from_api(&(me->triangles[16]), d.triangles[16]);
	from_api(&(me->triangles[17]), d.triangles[17]);
	from_api(&(me->triangles[18]), d.triangles[18]);
	from_api(&(me->triangles[19]), d.triangles[19]);
	from_api(&(me->triangles[20]), d.triangles[20]);
	from_api(&(me->triangles[21]), d.triangles[21]);
	from_api(&(me->triangles[22]), d.triangles[22]);
	from_api(&(me->triangles[23]), d.triangles[23]);
	from_api(&(me->triangles[24]), d.triangles[24]);
	from_api(&(me->triangles[25]), d.triangles[25]);
	from_api(&(me->triangles[26]), d.triangles[26]);
	from_api(&(me->triangles[27]), d.triangles[27]);
	from_api(&(me->triangles[28]), d.triangles[28]);
	from_api(&(me->triangles[29]), d.triangles[29]);
	from_api(&(me->triangles[30]), d.triangles[30]);
	from_api(&(me->triangles[31]), d.triangles[31]);
	from_api(&(me->triangles[32]), d.triangles[32]);
	from_api(&(me->triangles[33]), d.triangles[33]);
	from_api(&(me->triangles[34]), d.triangles[34]);
	from_api(&(me->triangles[35]), d.triangles[35]);
	from_api(&(me->triangles[36]), d.triangles[36]);
	from_api(&(me->triangles[37]), d.triangles[37]);
	from_api(&(me->triangles[38]), d.triangles[38]);
	from_api(&(me->triangles[39]), d.triangles[39]);
	from_api(&(me->triangles[40]), d.triangles[40]);
	from_api(&(me->triangles[41]), d.triangles[41]);
	from_api(&(me->triangles[42]), d.triangles[42]);
	from_api(&(me->triangles[43]), d.triangles[43]);
	from_api(&(me->triangles[44]), d.triangles[44]);
	from_api(&(me->triangles[45]), d.triangles[45]);
	from_api(&(me->triangles[46]), d.triangles[46]);
	from_api(&(me->triangles[47]), d.triangles[47]);
	from_api(&(me->triangles[48]), d.triangles[48]);
	from_api(&(me->triangles[49]), d.triangles[49]);
	from_api(&(me->triangles[50]), d.triangles[50]);
	from_api(&(me->triangles[51]), d.triangles[51]);
	from_api(&(me->triangles[52]), d.triangles[52]);
	from_api(&(me->triangles[53]), d.triangles[53]);
	from_api(&(me->triangles[54]), d.triangles[54]);
	from_api(&(me->triangles[55]), d.triangles[55]);
	from_api(&(me->triangles[56]), d.triangles[56]);
	from_api(&(me->triangles[57]), d.triangles[57]);
	from_api(&(me->triangles[58]), d.triangles[58]);
	from_api(&(me->triangles[59]), d.triangles[59]);
	from_api(&(me->triangles[60]), d.triangles[60]);
	from_api(&(me->triangles[61]), d.triangles[61]);
	from_api(&(me->triangles[62]), d.triangles[62]);
	from_api(&(me->triangles[63]), d.triangles[63]);
};


