package reply

type UnKnowErrReply struct {
}

var unKnowErrBytes = []byte("-Err unKnow\r\n")

func (u *UnKnowErrReply) Error() string {
	return "Err unKnow"
}

func (u *UnKnowErrReply) ToBytes() []byte {
	return unKnowErrBytes
}

type ArgNumErrReply struct {
	Cmd string
}

func (a *ArgNumErrReply) Error() string {
	return "-ERR wrong number of argument for '" + a.Cmd + "'command\r\n"
}

func (a *ArgNumErrReply) ToBytes() []byte {
	return []byte("-ERR wrong number of argument for '" + a.Cmd + "'command\r\n")
}

func MakeArgNumErrReply(cmd string) *ArgNumErrReply {
	return &ArgNumErrReply{
		Cmd: cmd,
	}
}

type SyntaxErrReply struct{}

var syntaxErrBytes = []byte("-Err syntax error\r\n")
var theSyntaxErrReply = &SyntaxErrReply{}

func MakeSyntaxErrReply() *SyntaxErrReply {
	return theSyntaxErrReply
}

func (s *SyntaxErrReply) Error() string {
	return "Err syntax error"
}

func (s *SyntaxErrReply) ToBytes() []byte {
	return syntaxErrBytes
}

type WrongTypeErrReply struct {
}

var wrongTypeErrReply = &WrongTypeErrReply{}

var wrongTypeErrBytes = []byte("-WRONGTYPE Operation against a key holding the wrong kind of value\r\n")

func (w *WrongTypeErrReply) Error() string {
	return "WRONGTYPE Operation against a key holding the wrong kind of value"
}

func (w *WrongTypeErrReply) ToBytes() []byte {
	return wrongTypeErrBytes
}

func MakeWrongTypeErrReply() *WrongTypeErrReply {
	return wrongTypeErrReply
}

type ProtocolErrReply struct {
	Msg string
}

func (p *ProtocolErrReply) Error() string {
	return "ERR Protocol error: '" + p.Msg
}

func (p *ProtocolErrReply) ToBytes() []byte {
	return []byte("-ERR Protocol error: '" + p.Msg + "'\r\n")
}

func MakeProtocolErrReply(msg string) *ProtocolErrReply {
	return &ProtocolErrReply{
		Msg: msg,
	}
}
