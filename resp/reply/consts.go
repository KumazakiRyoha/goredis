package reply

type PongReply struct {
}

var pongbytes = []byte("+PONG\r\n")

func (p *PongReply) ToBytes() []byte {
	return pongbytes
}

func MakePongReply() *PongReply {
	return &PongReply{}
}

type OkReply struct{}

var okBytes = []byte("+OK\r\n")

var theOkReply = new(OkReply)

func (o *OkReply) ToBytes() []byte {
	return okBytes
}

func MakeOkReply() *OkReply {
	return theOkReply
}

type NullBulReply struct{}

var nullBulBytes = []byte("$-1\r\n")

func (n *NullBulReply) ToBytes() []byte {
	return nullBulBytes
}

func MakeNullBulReply() *NullBulReply {
	return &NullBulReply{}
}

var emptyMultiBulBytes = []byte("*0\r\n")

type EmptyMultiBulReply struct{}

func (e *EmptyMultiBulReply) ToBytes() []byte {
	return emptyMultiBulBytes
}

func MakeEmptyMultiBulReply() *EmptyMultiBulReply {
	return &EmptyMultiBulReply{}
}

type NoReply struct{}

var noBytes = []byte("")

func (n NoReply) ToBytes() []byte {
	return noBytes
}

func MakeNoReply() *NoReply {
	return &NoReply{}
}
