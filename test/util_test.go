package kernel

import "testing"

func ensureRun(t *testing.T, nt int, d Data, numx, numy, numz int) {
	err := run(t, nt, d, numx, numy, numz)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func run(t *testing.T, nt int, d Data, numx, numy, numz int) error {
	k, err := New(nt, -1)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer k.Free()
	return k.Dispatch(d, numx, numy, numz)
}
