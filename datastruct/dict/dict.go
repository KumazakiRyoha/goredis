package dict

// 定义Consumer方法
type Consumer func(key string, val any) bool

type Dict interface {
	// 获取数据
	Get(key string) (val interface{}, exists bool)
	// 数据长度
	Len() int
	// 插入数据
	Put(key string, val interface{}) (result int)
	// 如果不存在，则插入
	PutIfAbsent(key string, val interface{}) (result int)
	// 如果存在，则插入
	PutIfExists(key string, val interface{}) (result int)
	// 移除数据
	Remove(key string) (result int)
	// 遍历
	ForEach(consumer Consumer)
	// 返回key
	Keys() []string
	// 随机返回key
	Randomkeys(limit int) []string
	//
	RandomDistinctKeys(limit int) []string
	//
	Clear()
}
