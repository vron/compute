package kernel

import (
	"reflect"
	"testing"
)

// TODO: deepequal to ensure the two methods return the same (in effect verifying the enc dec methods)

func ensureRun(t *testing.T, nt int, numx, numy, numz int, b func() Data, c func(Data)) {

	d2, err := runRaw(t, nt, numx, numy, numz, b, c)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	d1, err := run(t, nt, numx, numy, numz, b, c)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(d1, d2) {
		t.Log(d1)
		t.Log(d2)
		t.Error("data not equal from raw and not raw")
		t.FailNow()
	}
}

func run(t *testing.T, nt int, numx, numy, numz int, b func() Data, c func(Data)) (*Data, error) {
	k, err := New(nt, -1)
	if err != nil {
		return nil, err
	}
	defer k.Free()
	d := b()
	err = k.Dispatch(d, numx, numy, numz)
	if err != nil {
		return nil, err
	}
	c(d)
	return &d, nil
}

func runRaw(t *testing.T, nt int, numx, numy, numz int, b func() Data, c func(Data)) (*Data, error) {
	// encode it to a DataRaw struct and also compare to ensure it is all good
	k, err := New(nt, -1)
	if err != nil {
		return nil, err
	}
	defer k.Free()
	d := b()
	rd := encodeData(d)
	err = k.DispatchRaw(rd, numx, numy, numz)
	if err != nil {
		return nil, err
	}
	d = decodeData(rd)
	c(d)
	return &d, nil
}
