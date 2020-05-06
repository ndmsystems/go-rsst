package api

// Info ...
type Info struct {
	ID   uint16
	Data []byte
	Ok   bool // handled in Process
}

// Handler ..
type Handler func(*Info)

// Registrator ...
type Registrator interface {
	AddHandler(ID uint16, f Handler)
}

// Packer ...
type Packer interface {
	Pack([]Info) []byte
	Unpack([]byte) []Info
	PackRequest([]Info) []byte
	UnpackResponse([]byte) []Info
}

// Rsst ...
type Rsst interface {
	Registrator
	Process([]Info)
}
