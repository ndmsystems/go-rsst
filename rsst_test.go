package rsst

import (
	"bytes"
	"testing"

	rsstApi "github.com/tdx/go/api/rsst"
)

var (
	data  = []byte{1, 2, 3, 4, 5, 6, 7}
	input = []byte{
		0x10, 0, 1, 2, 3, 4, 5, 6, 7, 0,
		0x20, 0,
		0x40, 0,
		0x40, 4,
	}
	output = []byte{
		// 0x10, 0, 1, 2, 3, 4, 5, 6, 7, 0,
		0x20, 0, 0, 5,
		0x40, 0, 48, 46, 51, 48, 54, 0, // 0.306
		0x40, 4, 78, 68, 78, 83, 0, // NDNS
	}
	infos = []rsstApi.Info{
		{ID: 0x1000, Data: data, Ok: false},
		{ID: 0x2000, Data: []byte{0, 5}, Ok: true},
		{ID: 0x4000, Data: []byte("0.306"), Ok: true},
		{ID: 0x4004, Data: []byte("NDNS"), Ok: true},
	}
)

func TestPack(t *testing.T) {
	buf := Pack(infos)
	if !bytes.Equal(buf, output) {
		t.Fatalf("exprected buf=%v, got %v", output, buf)
	}
}

func TestUnpack(t *testing.T) {
	infos := Unpack(input)
	if len(infos) != 4 {
		t.Fatalf("exprected 4, got %d infos", len(infos))
	}

	var tests = []struct {
		id   uint16
		data []byte
	}{
		{0x1000, data},
		{0x2000, nil},
		{0x4000, nil},
		{0x4004, nil},
	}

	for i, test := range tests {
		info := infos[i]
		if info.ID != test.id {
			t.Fatalf("exprected infos[%d].ID=%x, got %x", i, test.id, info.ID)
		}
		if !bytes.Equal(test.data, info.Data) {
			t.Fatalf("exprected infos[%d].Data=%v, got %v",
				i, test.data, info.Data)
		}
	}

}