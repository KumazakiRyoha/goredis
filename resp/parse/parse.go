package parse

import (
	"bufio"
	"errors"
	"github.com/hdt3213/godis/lib/logger"
	"goredis/goredis/interface/resp"
	"goredis/goredis/resp/reply"
	"io"
	"runtime/debug"
	"strconv"
	"strings"
)

type Payload struct {
	Data resp.Reply
	Err  error
}

// 解析器状态
type readState struct {
	// 该解析器是在解析单行还是多行数据
	readingMultiLine bool
	// 正在读取的指令应该有多少参数
	expectedArgsCount int
	// 用户消息的类型
	msgType byte
	// 消息
	args [][]byte
	// 字节组的长度
	bulkLen int64
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
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	bufReader := bufio.NewReader(reader)
	var state readState
	var err error
	var msg []byte
	for {
		// read line
		var ioErr bool
		msg, ioErr, err = readLine(bufReader, &state)
		if err != nil {
			if ioErr { // encounter io err, stop read
				ch <- &Payload{
					Err: err,
				}
				close(ch)
				return
			}
			// protocol err, reset read state
			ch <- &Payload{
				Err: err,
			}
			state = readState{}
			continue
		}

		// parse line
		// 判断是不是多行解析模式
		if !state.readingMultiLine {
			// receive new response
			if msg[0] == '*' {
				// multi bulk reply
				err = parseMultiBulkHeader(msg, &state)
				if err != nil {
					ch <- &Payload{
						Err: errors.New("protocol error: " + string(msg)),
					}
					state = readState{} // reset state
					continue
				}
				if state.expectedArgsCount == 0 {
					ch <- &Payload{
						Data: &reply.EmptyMultiBulReply{},
					}
					state = readState{} // reset state
					continue
				}
			} else if msg[0] == '$' { // bulk reply
				err = parseBulkHeader(msg, &state)
				if err != nil {
					ch <- &Payload{
						Err: errors.New("protocol error: " + string(msg)),
					}
					state = readState{} // reset state
					continue
				}
				if state.bulkLen == -1 { // null bulk reply
					ch <- &Payload{
						Data: &reply.NullBulReply{},
					}
					state = readState{} // reset state
					continue
				}
			} else {
				// single line reply
				result, err := parseSingleLineReply(msg)
				ch <- &Payload{
					Data: result,
					Err:  err,
				}
				state = readState{} // reset state
				continue
			}
		} else {
			// receive following bulk reply
			err = readBody(msg, &state)
			if err != nil {
				ch <- &Payload{
					Err: errors.New("protocol error: " + string(msg)),
				}
				state = readState{} // reset state
				continue
			}
			// if sending finished
			if state.finished() {
				var result resp.Reply
				if state.msgType == '*' {
					result = reply.MakeMultiBulReply(state.args)
				} else if state.msgType == '$' {
					result = reply.MakeBulReply(state.args[0])
				}
				ch <- &Payload{
					Data: result,
					Err:  err,
				}
				state = readState{}
			}
		}
	}
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

// *3\r\n$3\r\nSET\R\N$3\r\n
func parseMultiBulkHeader(msg []byte, state *readState) error {
	var err error
	var expectedLine uint64

	expectedLine, err = strconv.ParseUint(string(msg[1:len(msg)-2]), 10, 32)
	if err != nil {
		return errors.New("protocol error: " + string(msg))
	}
	if expectedLine == 0 {
		state.expectedArgsCount = 0
		return nil
	} else if expectedLine > 0 {
		// 解析器状态，第一个字符为*，表示在读取一个字符
		state.msgType = msg[0]
		// 表示在读取多行
		state.readingMultiLine = true
		state.expectedArgsCount = int(expectedLine)
		state.args = make([][]byte, 0, expectedLine)
		return nil
	} else {
		return errors.New("protocol error: " + string(msg))
	}
}

// $4\r\nPING\r\n
func parseBulkHeader(msg []byte, state *readState) error {
	var err error
	state.bulkLen, err = strconv.ParseInt(string(msg[1:len(msg)-2]), 10, 64)
	if err != nil {
		return errors.New("protocol error: " + string(msg))
	}
	if state.bulkLen == -1 {
		return nil
	} else if state.bulkLen > 0 {
		state.msgType = msg[0]
		state.readingMultiLine = true
		state.expectedArgsCount = 1
		state.args = make([][]byte, 0, 1)
		return nil
	} else {
		return errors.New("protocol error: " + string(msg))
	}
}

// +OK -err
func parseSingleLineReply(msg []byte) (resp.Reply, error) {
	str := strings.TrimSuffix(string(msg), "\r\n")
	var result resp.Reply
	switch msg[0] {
	case '+':
		result = reply.MakeStatusReply(str[1:])
	case '-':
		result = reply.MakeErrReply(str[1:])
	case ':':
		val, err := strconv.ParseInt(str[1:], 10, 64)
		if err != nil {
			return nil, errors.New("protocol error: " + string(msg))
		}
		result = reply.MakeIntReply(val)
	}
	return result, nil
}

// read the non-first lines of multi bulk reply or bulk reply
func readBody(msg []byte, state *readState) error {
	line := msg[0 : len(msg)-2]
	var err error
	if line[0] == '$' {
		state.bulkLen, err = strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return errors.New("protocol error: " + string(msg))
		}
		if state.bulkLen <= 0 { // null bulk in multi bulks
			state.args = append(state.args, []byte{})
			state.bulkLen = 0
		}
	} else {
		state.args = append(state.args, line)
	}
	return nil
}
