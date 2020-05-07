package rsst

import (
	rsstApi "github.com/tdx/go-rsst/api"
)

// Unpack request
func Unpack(buf []byte) []rsstApi.Info {
	bufLen := len(buf)
	if bufLen == 0 {
		return nil
	}

	infos := make([]rsstApi.Info, 0)
	var (
		i = 0
		n = len(buf)
	)
	for n > 0 {

		id := uint16(buf[i])<<8 + uint16(buf[i+1])
		i += 2
		n -= 2

		if id < 0x1000 || id > 0x7999 {
			return infos
		}

		info := rsstApi.Info{
			ID: id,
		}

		// contains zero terminated string parameter
		if info.ID < 0x2000 {
			if buf[i] == 0 {
				continue
			}
			for buf[i] != 0 {
				info.Data = append(info.Data, buf[i])
				i++
				n--
			}
			i++
			n--
			if i >= bufLen {
				return infos
			}
		}
		infos = append(infos, info)
	}

	return infos
}

// Pack response
func Pack(infos []rsstApi.Info) []byte {
	if len(infos) == 0 {
		return nil
	}

	n := 0
	for i := range infos {
		if !infos[i].Ok {
			continue
		}
		n += 2
		n += len(infos[i].Data)
		if infos[i].ID >= 0x4000 {
			n++
		}
	}

	buf := make([]byte, 0, n)
	for _, info := range infos {
		if !info.Ok {
			continue
		}
		buf = append(buf, byte(info.ID>>8))
		buf = append(buf, byte(info.ID&0xff))
		buf = append(buf, info.Data...)
		if info.ID >= 0x4000 {
			buf = append(buf, 0)
		}
	}

	return buf
}

// UnpackResponse ...
func UnpackResponse(buf []byte) []rsstApi.Info {
	bufLen := len(buf)
	if bufLen == 0 {
		return nil
	}

	infos := make([]rsstApi.Info, 0)
	var (
		i = 0
		n = len(buf)
	)
	for n > 0 {
		id := uint16(buf[i])<<8 + uint16(buf[i+1])
		i += 2
		n -= 2

		if id < 0x1000 || id > 0x7999 {
			return infos
		}

		info := rsstApi.Info{
			ID: id,
		}

		// contains zero terminated string parameter
		if info.ID >= 0x4000 {
			if buf[i] == 0 {
				continue
			}
			for buf[i] != 0 {
				info.Data = append(info.Data, buf[i])
				i++
				n--
				if i >= bufLen {
					return infos
				}
			}
			i++
			n--
		}
		infos = append(infos, info)
	}

	return infos
}

// PackRequest ...
func PackRequest(infos []rsstApi.Info) []byte {
	if len(infos) == 0 {
		return nil
	}

	n := 0
	for i := range infos {
		if !infos[i].Ok {
			continue
		}
		n += 2
		n += len(infos[i].Data)
		if infos[i].ID < 0x2000 {
			n++
		}
	}

	buf := make([]byte, 0, n)
	for _, info := range infos {
		if !info.Ok {
			continue
		}
		buf = append(buf, byte(info.ID>>8))
		buf = append(buf, byte(info.ID&0xff))
		if len(info.Data) > 0 {
			buf = append(buf, info.Data...)
			if info.ID < 0x2000 {
				buf = append(buf, 0)
			}
		}
	}

	return buf
}
