package database

import (
	"github.com/hdt3213/godis/lib/logger"
	"goredis/config"
	"goredis/interface/resp"
	"goredis/resp/reply"
	"strconv"
	"strings"
)

type DataBase struct {
	dbSet []*DB
}

func NewDataBase() *DataBase {
	database := &DataBase{}
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}
	database.dbSet = make([]*DB, config.Properties.Databases)
	for i := range database.dbSet {
		db := makeDB()
		db.index = i
		database.dbSet[i] = db
	}
	return database
}

// set k v
// get k
// select 2
func (db *DataBase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
	cmdName := strings.ToLower(string(args[0]))
	if cmdName == "select" {
		if len(args) == 2 {
			return reply.MakeArgNumErrReply("select")
		}
		return execSelect(client, db, args[1:])
	}

	dbIndex := client.GetDBIndex()
	database := db.dbSet[dbIndex]
	return database.Exec(client, args)
}

func (db *DataBase) Close() {

}

func (db *DataBase) AfterClientClose(c resp.Connection) {

}

// select 2
func execSelect(conn resp.Connection, database *DataBase, args [][]byte) resp.Reply {
	dbIndex, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.MakeErrReply("ERR invalid DB index")
	}
	if dbIndex > len(database.dbSet) {
		return reply.MakeErrReply("ERR DB index is out of range")
	}
	conn.SelectDB(dbIndex)
	return reply.MakeOkReply()
}
