package buffer

import "testing"

func TestSlice(t *testing.T) {
	s := make([]byte, 100, 101)
	s[0] = 1

	d := s[:103]
	t.Logf("s:%v, d:%v", s, d)
}
