package httpext

import "testing"

func TestQualityValue(t *testing.T) {
	type args struct {
		v  string
		qv float32
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "in-range",
			args: args{v: "test", qv: 0.5},
			want: "test;q=0.5",
		},
		{
			name: "in-range-trailing-zeros",
			args: args{v: "test", qv: 0.500},
			want: "test;q=0.5",
		},
		{
			name: "greater-than-range",
			args: args{v: "test", qv: 1.500},
			want: "test;q=1",
		},
		{
			name: "less-than-range",
			args: args{v: "test", qv: 0.0000001},
			want: "test;q=0.001",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := QualityValue(tt.args.v, tt.args.qv); got != tt.want {
				t.Errorf("QualityValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
