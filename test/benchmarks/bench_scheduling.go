package kernel

/*
	A wide and small shader in order to benchmar scheduling stuff, such as
	launching wg's, syning etc.
*/
import (
	"reflect"
	"testing"
	"unsafe"
)

var shader = `
#version 450

layout(local_size_x = 4, local_size_y = 4, local_size_z = 4) in;

layout(std430) buffer Out {
	int data[];
};

void main() {
	uint index = gl_LocalInvocationIndex;
	int value = int(index);

	if (index % 2 == 0) {
		value *= 2;
	}

	barrier();

	if (index < 1024) {
		value = 1;
	}
	for (int i = 0; i < 4; i++) {
		if (index % 2 == 0) {
			value += 1;
		}

		barrier();
	}

	barrier();

	data[index] = value;
}
`

func BenchmarkScheduling(b *testing.B) {
	data := make([]int32, 64)
	d := Data{Data: intToByte(data)}
	k, err := New(-1)
	if err != nil {
		b.Error(err)
		b.FailNow()
	}
	defer k.Free()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := k.Dispatch(d, 8, 8, 8)
		if err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()

	for i := range data {
		ex := int32(1)
		if i%2 == 0 {
			ex = 5
		}
		if data[i] != ex {
			b.Error(i, "expected value: ", ex, "got", data[i])
		}
	}
}

func intToByte(raw []int32) []byte {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&raw))
	header.Len *= 4
	header.Cap *= 4
	data := *(*[]byte)(unsafe.Pointer(&header))
	return data
}
