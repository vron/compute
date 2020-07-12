package kernel

/*
	A numeric-heavy bench where we sequentially do a couple
	of trigonometric functions on a grid.
*/
import (
	"math"
	"reflect"
	"testing"
	"unsafe"
)

var shader = `
#version 450

layout(local_size_x = 64, local_size_y = 1, local_size_z = 1) in;

layout(std430) buffer In {
	float din[];
};

layout(std430) buffer Out {
	float dout[];
};

void main() {
	uint base = gl_GlobalInvocationID.x*1024;

	// copy over the data
	for(uint i = 0; i < 1024; i++) {
		dout[base+i] = din[base+i];
	}

	// sin of each array element
	for(uint i = 0; i < 1024; i++) {
		dout[base+i] = sin(dout[base+i]);
	}

	// cos of each array element
	for(uint i = 0; i < 1024; i++) {
		dout[base+i] = cos(dout[base+i]);
	}

	// tan of each array element
	for(uint i = 0; i < 1024; i++) {
		dout[base+i] = tan(dout[base+i]);
	}
}
`

func BenchmarkTrig(b *testing.B) {
	// create the input data
	din := make([]float32, 64*1024*24)
	for i := range din {
		din[i] = math.Pi / 4
	}
	dout := make([]float32, 64*1024*24)

	d := Data{Din: floatToByte(din), Dout: floatToByte(dout)}
	k, err := New(-1)
	if err != nil {
		b.Error(err)
		b.FailNow()
	}
	defer k.Free()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := k.Dispatch(d, 24, 1, 1)
		if err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()

	for i := range dout {
		if math.Abs(float64(dout[i])-0.9509171) > 1e-5 {
			b.Error("bad value", i, dout[i])
		}
	}
}

func BenchmarkTrigRef(b *testing.B) {
	// create the input data
	din := make([]float32, 64*1024*24)
	for i := range din {
		din[i] = math.Pi / 4
	}
	dout := make([]float32, 64*1024*24)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(dout, din)
		for i := range dout {
			dout[i] = float32(math.Sin(float64(dout[i])))
			dout[i] = float32(math.Cos(float64(dout[i])))
			dout[i] = float32(math.Tan(float64(dout[i])))
		}
	}
	b.StopTimer()

	for i := range dout {
		if math.Abs(float64(dout[i])-0.9509171) > 1e-5 {
			b.Error("bad value", i, dout[i])
		}
	}
}

func floatToByte(raw []float32) []byte {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&raw))
	header.Len *= 4
	header.Cap *= 4
	data := *(*[]byte)(unsafe.Pointer(&header))
	return data
}
