package reply

import (
	"bytes"
	"goredis/interface/resp"
	"strconv"
)

var (
	nullBulReplyBytes = []byte("-1")
	CRLF              = "\r\n"
)

type BulReply struct {
	Arg []byte // "moody" "$5\r\nmoody\r\n"
}

func (b BulReply) ToBytes() []byte {
	if len(b.Arg) == 0 {
		return nullBulReplyBytes
	}
	return []byte("$" + strconv.Itoa(len(b.Arg)) + CRLF + string(b.Arg) + CRLF)
}

func MakeBulReply(arg []byte) *BulReply {
	return &BulReply{
		Arg: arg,
	}
}

type MultiBulReply struct {
	Args [][]byte
}

func (r *MultiBulReply) ToBytes() []byte {
	argLen := len(r.Args)
	var buf bytes.Buffer
	buf.WriteString("*" + strconv.Itoa(argLen) + CRLF)
	for _, arg := range r.Args {
		if arg == nil {
			buf.WriteString(string(nullBulReplyBytes) + CRLF)
		} else {
			buf.WriteString("$" + strconv.Itoa(len(arg)) + CRLF + string(arg) + CRLF)
		}
	}
	return buf.Bytes()

}

func MakeMultiBulReply(arg [][]byte) *MultiBulReply {
	return &MultiBulReply{
		Args: arg,
	}
}

type StatusReply struct {
	Status string
}

func MakeStatusReply(status string) *StatusReply {
	return &StatusReply{
		Status: status,
	}
}

func (s *StatusReply) ToBytes() []byte {
	return []byte("+" + s.Status + CRLF)
}

type IntReply struct {
	Code int64
}

func MakeIntReply(code int64) *IntReply {
	return &IntReply{
		Code: code,
	}
}

func (i *IntReply) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(i.Code, 10) + CRLF)
}

type ErrorReply interface {
	Error() string
	ToBytes() []byte
}

type StandardErrReply struct {
	Status string
}

func (s *StandardErrReply) ToBytes() []byte {
	return []byte("-" + s.Status + CRLF)
}

func (s *StandardErrReply) Error() string {
	return s.Status
}

// MakeErrReply creates StandardErrReply
func MakeErrReply(status string) *StandardErrReply {
	return &StandardErrReply{
		Status: status,
	}
}

// IsErrorReply returns true if the given reply is error
func IsErrReply(reply resp.Reply) bool {
	return reply.ToBytes()[0] == '-'
}
