package unsafeext

import "testing"

func TestBytesToString(t *testing.T) {
	b := []byte{'g', 'o', '-', 'p', 'l', 'a', 'y', 'g', 'r', 'o', 'u', 'n', 'd'}
	s := BytesToString(b)
	expected := string(b)
	if s != expected {
		t.Fatalf("expected '%s' got '%s'", expected, s)
	}
}

func TestStringToBytes(t *testing.T) {
	s := "go-playground"
	b := StringToBytes(s)

	if string(b) != s {
		t.Fatalf("expected '%s' got '%s'", s, string(b))
	}
}
