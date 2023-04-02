package database

import "strings"

// 记录指令河command结构体之间的关系
var cmdTable = make(map[string]*command)

type command struct {
	exector ExecFunc
	arity   int
}

func RegisterCommand(name string, exector ExecFunc, arity int) {
	name = strings.ToLower(name)
	cmdTable[name] = &command{
		exector: exector,
		arity:   arity,
	}
}
