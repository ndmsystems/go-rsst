package rsst

import (
	rsstApi "github.com/tdx/go-rsst/api"
)

// UnpackRequest ...
func UnpackRequest(buf []byte) []rsstApi.Info {
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
			if buf[i] == 0 { // zero string ?
				info.Ok = true
				infos = append(infos, info)
				i++
				n--
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
			info.Ok = true
			infos = append(infos, info)
			continue
		}

		// > 0x19999
		// no data
		info.Ok = true

		infos = append(infos, info)
	}

	return infos
}

// PackResponse ...
func PackResponse(infos []rsstApi.Info) []byte {
	if len(infos) == 0 {
		return nil
	}

	n := 0
	for i := range infos {
		if !infos[i].Ok {
			continue
		}

		// check data len
		if infos[i].ID < 0x2000 && len(infos[i].Data) != 0 {
			continue
		}
		if infos[i].ID >= 0x2000 && infos[i].ID < 0x4000 &&
			len(infos[i].Data) != 2 {

			continue
		}

		// id
		n += 2

		// data len
		n += len(infos[i].Data)

		if infos[i].ID >= 0x4000 {
			// trailing zero byte
			n++
		}
	}

	buf := make([]byte, 0, n)

	for i := range infos {
		info := infos[i]

		if !info.Ok {
			continue
		}

		// check data len
		if info.ID < 0x2000 && len(info.Data) != 0 {
			continue
		}
		if info.ID >= 0x2000 && info.ID < 0x4000 && len(info.Data) != 2 {
			continue
		}

		// store ID
		buf = append(buf, byte(info.ID>>8))
		buf = append(buf, byte(info.ID&0xff))

		// data
		buf = append(buf, info.Data...)

		if info.ID >= 0x4000 {
			// trailing zero byte
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
	for n > 1 {
		id := uint16(buf[i])<<8 + uint16(buf[i+1])
		i += 2
		n -= 2

		// fmt.Printf("%x: i=%d n=%d\n", id, i, n)

		if id < 0x1000 || id > 0x7999 {
			return infos
		}

		info := rsstApi.Info{
			ID: id,
		}

		if id < 0x2000 {
			// only id
			info.Ok = true
			infos = append(infos, info)
			continue
		}

		// contains zero terminated string parameter
		if info.ID >= 0x4000 {
			if buf[i] == 0 { // zero string ?
				info.Ok = true
				infos = append(infos, info)
				i++
				n--
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
			info.Ok = true
			infos = append(infos, info)
			continue
		}

		// < 0x4000
		if len(buf[i:]) < 2 {
			return infos
		}
		info.Data = append(info.Data, buf[i], buf[i+1])
		info.Ok = true
		i += 2
		n -= 2

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
