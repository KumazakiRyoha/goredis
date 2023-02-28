package parse

import (
	"bufio"
	"errors"
	"goredis/goredis/interface/resp"
	"io"
)

type Payload struct {
	Data resp.Reply
	Err  error
}

type readState struct {
	readingMultiLine  bool
	expectedArgsCount int
	msgType           byte
	args              [][]byte
	bulkLen           int64
}

func (r *readState) finished() bool {
	return r.expectedArgsCount > 0 && len(r.args) == r.expectedArgsCount
}

func ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	go parse0(reader, ch)
	return ch
}

func parse0(reader io.Reader, ch chan<- *Payload) {

}

func readLine(bufReader *bufio.Reader, state *readState) ([]byte, bool, error) {
	//1.\r\n切分
	//2.之前读到了数字，严格取子政府个数

	var msg []byte
	var err error
	if state.bulkLen == 0 { //1.\r\n区分
		msg, err = bufReader.ReadBytes('\n')
		if err != nil {
			return nil, true, err
		}
		// 如果msg 长度为0或者消息的倒数第二个字符不为'\r'则说明协议错误
		if len(msg) == 0 || msg[len(msg)-2] != '\r' {
			return nil, false, errors.New("protocol error: " + string(msg))
		}

	} else { //2.表示之前读到了$数字，严格读取字符个数
		msg = make([]byte, state.bulkLen+2)
		_, err := io.ReadFull(bufReader, msg)
		if err != nil {
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' || msg[len(msg)-1] != '\n' {
			return nil, false, errors.New("protocol error: " + string(msg))
		}

	}
	return msg, false, nil
}
