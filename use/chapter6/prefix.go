package chapter6

import "github.com/go-redis/redis/v8"

type AutoComplete interface {
	// 根据前缀自动补全
	PrefixAutoComplete(string, string)
	// 随机补全
	RandomAutoComplete(string, string)
}

type PrefixStruct struct {
	Conn *redis.Client
}

func NewClient(conn *redis.Client) *PrefixStruct {
	return &PrefixStruct{
		Conn: conn,
	}
}

func (*PrefixStruct) PrefixAutoComplete(guild, prefix string) {
}

func (*PrefixStruct) findPrefixRange(prefix string) {
}
