package rsst

import (
	"bytes"
	"testing"

	rsstApi "github.com/tdx/go-rsst/api"
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
	bad  = []byte{200, 153, 128, 45, 191, 148, 6, 91}
	inp2 = []byte{
		32, 0, 0, 0,
		32, 1, 0, 3,
		32, 2, 0, 5,
		64, 4, 75, 101, 101, 110, 101, 116, 105, 99, 32, 65, 105, 114, 0,
		64, 5, 75, 101, 101, 110, 101, 116, 105, 99, 32, 65, 105, 114, 0,
	}
	b1     = []byte{75, 101, 101, 110, 101, 116, 105, 99, 32, 65, 105, 114}
	b2     = []byte{75, 101, 101, 110, 101, 116, 105, 99, 32, 65, 105, 114}
	infos2 = []rsstApi.Info{
		{ID: 0x2000, Data: []byte{0, 0}, Ok: true},
		{ID: 0x2001, Data: []byte{0, 3}, Ok: true},
		{ID: 0x2002, Data: []byte{0, 5}, Ok: true},
		{ID: 0x4004, Data: b1, Ok: true},
		{ID: 0x4005, Data: b2, Ok: true},
	}
	zero = []byte{
		64, 4, 0,
		32, 0, 0, 0,
		64, 5, 0,
		32, 2, 0, 5,
	}
	zeroInfos = []rsstApi.Info{
		{ID: 16388, Data: nil, Ok: true},
		{ID: 8192, Data: []byte{0, 0}, Ok: true},
		{ID: 16389, Data: nil, Ok: true},
		{ID: 8194, Data: []byte{0, 5}, Ok: true},
	}
	// request
	reqInfos = []rsstApi.Info{
		{ID: 0x2000, Data: nil, Ok: true},
		{ID: 0x2001, Data: nil, Ok: true},
		{ID: 0x2002, Data: nil, Ok: true},
		{ID: 0x4004, Data: nil, Ok: true},
		{ID: 0x4005, Data: nil, Ok: true},
	}
	reqData = []byte{
		0x20, 0,
		0x20, 1,
		0x20, 2,
		0x40, 4,
		0x40, 5,
	}
	// 0x1xxx test
	d1101 = []rsstApi.Info{
		{ID: 4353,
			Data: []byte{110, 100, 110, 115, 47, 114, 101, 109, 111, 116, 101,
				73, 110, 102, 111, 32, 109, 101, 116, 104, 111, 100}, Ok: true},
	}
)

func TestPack(t *testing.T) {
	buf := PackResponse(infos)
	if !bytes.Equal(buf, output) {
		t.Fatalf("exprected buf=%v, got %v", output, buf)
	}
}

func TestBad(t *testing.T) {
	infos := UnpackRequest(bad)
	if len(infos) != 0 {
		t.Fatalf("expected empty Infos from bad input, got: %#v", infos)
	}
}

func TestUnpackResponseInp2(t *testing.T) {
	infos := UnpackResponse(inp2)
	if len(infos) != 5 {
		t.Fatalf("unpack failed: expected 5 Infos, got %d: %v",
			len(infos), infos)
	}
	for i := range infos {
		info := infos[i]
		good := infos2[i]

		if info.ID != good.ID {
			t.Fatalf("exprected infos[%d].ID=%x, got %x", i, good.ID, info.ID)
		}
		if !bytes.Equal(info.Data, good.Data) {
			t.Fatalf("exprected infos[%d].Data=%v, got %v",
				i, good.Data, info.Data)
		}
	}
}

func TestUnpackResponseWithZeroStrings(t *testing.T) {
	infos := UnpackResponse(zero)
	if len(infos) != 4 {
		t.Fatalf("unpack failed: expected 4 Infos, got %d: %v",
			len(infos), infos)
	}
	for i := range infos {
		info := infos[i]
		good := zeroInfos[i]

		if info.ID != good.ID {
			t.Fatalf("exprected infos[%d].ID=%x, got %x", i, good.ID, info.ID)
		}
		if !bytes.Equal(info.Data, good.Data) {
			t.Fatalf("exprected infos[%d].Data=%v, got %v",
				i, good.Data, info.Data)
		}
	}
}

func TestPackRequestInp2(t *testing.T) {
	data := PackRequest(reqInfos)
	if !bytes.Equal(reqData, data) {
		t.Fatalf("exprected %v, got %v", inp2, data)
	}

	infos := UnpackRequest(data)
	if len(infos) != len(reqInfos) {
		t.Fatalf("expected %d infos, got %d: %v",
			len(infos), len(reqInfos), infos)
	}
	for i := range infos {
		info := infos[i]
		good := reqInfos[i]

		if info.ID != good.ID {
			t.Fatalf("exprected infos[%d].ID=%x, got %x", i, good.ID, info.ID)
		}
	}
}

func TestUnpack(t *testing.T) {
	infos := UnpackRequest(input)
	if len(infos) != 4 {
		t.Fatalf("exprected 4, got %d infos: %v", len(infos), infos)
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

func Test1xxxPack(t *testing.T) {
	data := PackResponse(d1101)
	if len(data) != 0 {
		t.Fatalf("expected bad 0x1xxx response is skiped, got: %v", data)
	}
}
