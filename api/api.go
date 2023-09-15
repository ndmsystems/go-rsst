package api

type AuthorizeLevel int

const (
	NDSSTrust AuthorizeLevel = iota
	TokenTrust
)

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
	PackResponse([]Info) []byte
	UnpackRequest([]byte) []Info
	PackRequest([]Info) []byte
	UnpackResponse([]byte) []Info
}

// Rsst ...
type Rsst interface {
	Registrator
	Process(AuthorizeLevel, []Info)
}
