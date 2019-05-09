package ioutil

import (
	"strings"
	"testing"
)

func TestLimitedReader_Read(t *testing.T) {
	eofLimited := LimitReader(strings.NewReader("all"), 5).(*LimitedReader)

	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		l       *LimitedReader
		args    args
		wantN   int
		wantErr bool
	}{
		{
			name:    "not-limited",
			l:       LimitReader(strings.NewReader("all"), 3).(*LimitedReader),
			args:    args{p: make([]byte, 4)},
			wantN:   3,
			wantErr: false,
		},
		{
			name:    "not-limited-exact",
			l:       LimitReader(strings.NewReader("all"), 3).(*LimitedReader),
			args:    args{p: make([]byte, 3)},
			wantN:   3,
			wantErr: false,
		},
		{
			name:    "not-limited-EOF-OK",
			l:       eofLimited,
			args:    args{p: make([]byte, 4)},
			wantN:   3,
			wantErr: false,
		},
		{
			name:    "not-limited-EOF",
			l:       eofLimited,
			args:    args{p: make([]byte, 4)},
			wantN:   0,
			wantErr: true,
		},
		{
			name:    "limited",
			l:       LimitReader(strings.NewReader("limited"), 1).(*LimitedReader),
			args:    args{p: make([]byte, 3)},
			wantN:   2, // need to read one past to know we're past
			wantErr: true,
		},
		{
			name:    "limited-buff-under-N",
			l:       LimitReader(strings.NewReader("limited"), 0).(*LimitedReader),
			args:    args{p: make([]byte, 1)},
			wantN:   1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN, err := tt.l.Read(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("LimitedReader.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("LimitedReader.Read() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}
