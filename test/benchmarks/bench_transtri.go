package kernel

/*
	A numeric-heavy and struct heavy bench where we get a sequence of triangle
	strips and a transform and find the geometric center of them all. Intentionally
	written with some unneeded complexity to include that in the optimization testing.
*/
import (
	"math"
	"runtime"
	"sync"
	"testing"
)

var shader = `
#version 450

layout(local_size_x = 64, local_size_y = 1, local_size_z = 1) in;

struct triangle {
	vec2 vertices[3];
};

struct polygon {
	triangle triangles[64];
};

struct cog_res {
	vec2 cog;
	float area;
};

layout(std430) buffer In {
	mat3 transform;
	polygon polygons[];
};

layout(std430) buffer Out {
	vec2 cogs[];
};

shared cog_res shared_data[64];

float area_tri(triangle t) {
	float a = (t.vertices[1][0]-t.vertices[0][0])*(t.vertices[2][1]-t.vertices[0][1]) - (t.vertices[2][0]-t.vertices[0][0])*(t.vertices[1][1]-t.vertices[0][1]);
	if (a < 0.0f) {
		a *= -0.5f;
	} else {
		a *= 0.5f;
	}
	return a;
}

vec2 cog_tri(triangle t) {
	return vec2((t.vertices[1][0] + t.vertices[2][0] + t.vertices[0][0]) / 3.0f,
		(t.vertices[1][1] + t.vertices[2][1] + t.vertices[0][1]) / 3.0f);
}

cog_res tri(triangle t) {
	for(int i = 0; i < 3; i++) {
		t.vertices[i] = (transform*vec3(t.vertices[i], 1.0f)).xy;
	}
	cog_res r;	
	r.area = area_tri(t);
	r.cog = cog_tri(t);
	return r;
}

cog_res cog_poly(polygon p) {
	float area = 0.0f;
	vec2 cog = vec2(0, 0);
	for(int i = 0; i < 64; i++) {
		cog_res tr = tri(p.triangles[i]);
		area += tr.area;
		cog += tr.area*tr.cog;
	}
	cog /= area;
	cog_res r;
	r.area = area;
	r.cog = cog;
	return r;
}

void main() {
	// where we should do our job
	uint base_index = gl_WorkGroupID.x*gl_WorkGroupSize.x;
	uint local_index = gl_LocalInvocationID.x;
	uint index = base_index + local_index;

	// actual calc of this invocation
	cog_res my_res = cog_poly(polygons[index]);

	// sync with others in our WG by using shared memory
	// and a barrier
	shared_data[local_index] = my_res;
	barrier();
	
	// one in each WG should sum them up and return the
	// result.
	if (local_index == 0) {
		my_res.cog *= my_res.area;
		for( int i = 1; i < 64; i++) {
			cog_res fr = shared_data[i];
			my_res.area += fr.area;
			my_res.cog += fr.area*fr.cog;
		}
		my_res.cog /= my_res.area;
		cogs[gl_WorkGroupID.x] = my_res.cog;
	}
}
`

func BenchmarkTransTri(b *testing.B) {
	noi := 128
	data := d(noi)

	k, err := New(runtime.GOMAXPROCS(-1), 1024*1024)
	if err != nil {
		b.Error(err)
		b.FailNow()
	}
	defer k.Free()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := k.Dispatch(data, noi, 1, 1)
		if err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()

	// chec the cog's
	for i := 0; i < noi; i++ {
		cog := data.Cogs[i]
		if math.Abs(float64(cog[0]-1)) > 1e-4 || math.Abs(float64(cog[1]-1)) > 1e-4 {
			b.Error("bas cog data", i, cog)
		}
	}
}

func TestTransTri(t *testing.T) {
	noi := 2
	ensureRun(t, 1, noi, 1, 1, func() Data { return d(noi) }, func(res Data) {
		// chec the cog's
		for i := 0; i < noi; i++ {
			cog := res.Cogs[i]
			if math.Abs(float64(cog[0]-1)) > 1e-4 || math.Abs(float64(cog[1]-1)) > 1e-4 {
				t.Error("bas cog data", i, cog)
			}
		}
	})
}

func d(noi int) Data {
	d := Data{
		Transform: &Mat3{Vec3{1, 0, 0}, Vec3{0, 1, 0}, Vec3{2.0 / 3, 2.0 / 3, 1}},
		Polygons:  make([]Polygon, noi*64),
		Cogs:      make([]Vec2, noi),
	}
	// fill with polygons
	for i := 0; i < noi*64; i++ {
		d.Polygons[i] = p()
	}
	return d
}

func dgo(noi int) (tr [9]float32, ps []Polygon, cogs []float32) {
	tr = [9]float32{1, 0, 0, 0, 1, 0, 2.0 / 3, 2.0 / 3, 1}
	ps = make([]Polygon, noi*64)
	cogs = make([]float32, noi*2)
	for i := 0; i < noi*64; i++ {
		ps[i] = p()
	}
	return
}

func p() Polygon {
	p := Polygon{}
	for i := 0; i < 64; i++ {
		p.Triangles[i] = t()
	}
	return p
}

func t() Triangle {
	// area 0.5 and COG at (1/3, 1/3)
	return Triangle{
		Vertices: [3]Vec2{{0, 0}, {0, 1}, {1, 0}},
	}
}

func BenchmarkTransTriRef(b *testing.B) {
	noi := 128
	tr, ps, cogs := dgo(noi)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		refimpl(tr, ps, cogs)
	}
	b.StopTimer()

	// chec the cog's
	for i := 0; i < noi; i++ {
		cog := [2]float32{cogs[i*2], cogs[i*2+1]}
		if math.Abs(float64(cog[0]-1)) > 1e-4 || math.Abs(float64(cog[1]-1)) > 1e-4 {
			b.Error("bas cog data", i, cog)
		}
	}
}

func refimpl(tr [9]float32, ps []Polygon, cogs []float32) {
	// simply do it in parts of 64
	wg := sync.WaitGroup{}
	si := 0
	noe := 64
	for i := 0; i < len(ps)/64; i++ {
		ei := si + noe
		wg.Add(1)
		go func(tr [9]float32, ps []Polygon, cogs []float32) {
			impl(tr, ps, cogs)
			wg.Done()
		}(tr, ps[si:ei], cogs[i*2:i*2+2])
		si = ei
	}
	wg.Wait()
}

func impl(tr [9]float32, ps []Polygon, cogs []float32) {
	if len(cogs) != 2 {
		panic("should be")
	}

	area, x, y := float32(0), float32(0), float32(0)
	for i := 0; i < 64; i++ {
		a, b, c := cog_poly(ps[i], tr)
		area += a
		x += a * b
		y += a * c
	}
	x /= area
	y /= area

	cogs[0] = x
	cogs[1] = y
}

func area_tri(t Triangle) float32 {
	a := (t.Vertices[1][0]-t.Vertices[0][0])*(t.Vertices[2][1]-t.Vertices[0][1]) - (t.Vertices[2][0]-t.Vertices[0][0])*(t.Vertices[1][1]-t.Vertices[0][1])
	if a < 0.0 {
		a *= -0.5
	} else {
		a *= 0.5
	}
	return a
}

func cog_tri(t Triangle) (float32, float32) {
	return (t.Vertices[1][0] + t.Vertices[2][0] + t.Vertices[0][0]) / 3.0,
		(t.Vertices[1][1] + t.Vertices[2][1] + t.Vertices[0][1]) / 3.0
}

func tri(t *Triangle, transform [9]float32) (area, x, y float32) {
	for i := 0; i < 3; i++ {
		t.Vertices[i][0], t.Vertices[i][1], _ = matmul(transform, t.Vertices[i][0], t.Vertices[i][1], 1.0)
	}
	area = area_tri(*t)
	x, y = cog_tri(*t)
	return
}

func matmul(m [9]float32, x, y, z float32) (float32, float32, float32) {
	return m[0]*x + m[3]*y + m[6]*z, m[0+1]*x + m[3+1]*y + m[6+1]*z, m[0+1+1]*x + m[3+1+1]*y + m[6+1+1]*z
}

func cog_poly(p Polygon, t [9]float32) (area, x, y float32) {
	for i := 0; i < 64; i++ {
		a, b, c := tri(&p.Triangles[i], t)
		area += a
		x += b * a
		y += c * a
	}
	x /= area
	y /= area
	return
}
