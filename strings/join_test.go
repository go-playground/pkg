package stringsext

import (
	"strings"
	"testing"
)

func TestJoin(t *testing.T) {
	s1, s2, s3 := "a", "b", "c"
	arr := []string{s1, s2, s3}
	if strings.Join(arr, ",") != Join(",", s1, s2, s3) {
		t.Errorf("Join failed")
	}
}
