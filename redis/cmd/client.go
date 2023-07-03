package main

import (
	"CometDB"
	"CometDB/redis"
	"CometDB/utils"
	"errors"
	"fmt"
	"github.com/tidwall/redcon"
	"strings"
)

func newWrongNumberOfArgsError(cmd string) error {
	return fmt.Errorf("ERR wrong number of arguments for '%s' command", cmd)
}

type cmdHandler func(cli *CometDBClient, args [][]byte) (interface{}, error)

var supportedCommands = map[string]cmdHandler{
	// string
	"set":    set,
	"get":    get,
	"strlen": strlen,

	// hash
	"hset": hset,
	"hget": hget,
	"hdel": hdel,
	"hlen": hlen,

	// set
	"sadd":      sadd,
	"sismember": sismember,
	"srem":      srem,
	"scard":     scard,

	// list
	"lpush": lpush,
	"rpush": rpush,
	"lpop":  lpop,
	"rpop":  rpop,
	"llen":  llen,

	// zset
	"zadd":   zadd,
	"zscore": zscore,

	// key
	"ping": nil,
	"type": typ,
	"del":  del,
}

type CometDBClient struct {
	server *CometDBServer
	db     *redis.RedisObject
}

func execClientCommand(conn redcon.Conn, cmd redcon.Command) {
	command := strings.ToLower(string(cmd.Args[0]))
	cmdFunc, ok := supportedCommands[command]
	if !ok {
		conn.WriteError("Err unsupported command: '" + command + "'")
		return
	}

	client, _ := conn.Context().(*CometDBClient)
	switch command {
	case "quit":
		_ = conn.Close()
	case "ping":
		conn.WriteString("PONG")
	default:
		res, err := cmdFunc(client, cmd.Args[1:])
		if err != nil {
			if err == CometDB.ErrKeyNotFound {
				conn.WriteNull()
			} else {
				conn.WriteError(err.Error())
			}
			return
		}
		conn.WriteAny(res)
	}
}

// =============== String ===============

func set(cli *CometDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 2 {
		return nil, newWrongNumberOfArgsError("set")
	}

	key, value := args[0], args[1]
	if err := cli.db.Set(key, 0, value); err != nil {
		return nil, err
	}
	return redcon.SimpleString("OK"), nil
}

func get(cli *CometDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, newWrongNumberOfArgsError("get")
	}

	value, err := cli.db.Get(args[0])
	if err != nil {
		return nil, err
	}
	return value, nil
}

func strlen(cli *CometDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, newWrongNumberOfArgsError("strlen")
	}

	value, err := cli.db.StrLen(args[0])
	if err != nil {
		return nil, err
	}
	return redcon.SimpleInt(value), nil
}

// =============== Hash ===============

func hset(cli *CometDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 3 {
		return nil, newWrongNumberOfArgsError("hset")
	}

	var ok = 0
	key, field, value := args[0], args[1], args[2]
	res, err := cli.db.HSet(key, field, value)
	if err != nil {
		return nil, err
	}
	if res {
		ok = 1
	}
	return redcon.SimpleInt(ok), nil
}
func hget(cli *CometDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 2 {
		return nil, newWrongNumberOfArgsError("hget")
	}
	key, field := args[0], args[1]
	return cli.db.HGet(key, field)
}

func hdel(cli *CometDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 2 {
		return nil, newWrongNumberOfArgsError("hdel")
	}
	var ok = 0
	key, field := args[0], args[1]
	del, err := cli.db.HDel(key, field)
	if err != nil {
		return nil, err
	}
	if del {
		ok = 1
	}
	return redcon.SimpleInt(ok), nil
}

func hlen(cli *CometDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, newWrongNumberOfArgsError("hlen")
	}

	key := args[0]
	res, err := cli.db.HLen(key)
	if err != nil {
		return nil, err
	}
	return redcon.SimpleInt(res), nil
}

// =============== Set ===============

func sadd(cli *CometDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 2 {
		return nil, newWrongNumberOfArgsError("sadd")
	}

	var ok = 0
	key, member := args[0], args[1]
	res, err := cli.db.SAdd(key, member)
	if err != nil {
		return nil, err
	}
	if res {
		ok = 1
	}
	return redcon.SimpleInt(ok), nil
}

func sismember(cli *CometDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 2 {
		return nil, newWrongNumberOfArgsError("sismember")
	}

	var ok = 0
	key, member := args[0], args[1]
	res, err := cli.db.SIsMember(key, member)
	if err != nil {
		return nil, err
	}
	if res {
		ok = 1
	}
	return redcon.SimpleInt(ok), nil
}

func srem(cli *CometDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 2 {
		return nil, newWrongNumberOfArgsError("srem")
	}

	var ok = 0
	key, member := args[0], args[1]
	res, err := cli.db.SRem(key, member)
	if err != nil {
		return nil, err
	}
	if res {
		ok = 1
	}
	return redcon.SimpleInt(ok), nil
}

func scard(cli *CometDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, newWrongNumberOfArgsError("scard")
	}

	key := args[0]
	res, err := cli.db.SCard(key)
	if err != nil {
		return nil, err
	}
	return redcon.SimpleInt(res), nil
}

// =============== List ===============

func lpush(cli *CometDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 2 {
		return nil, newWrongNumberOfArgsError("lpush")
	}

	key, value := args[0], args[1]
	res, err := cli.db.LPush(key, value)
	if err != nil {
		return nil, err
	}
	return redcon.SimpleInt(res), nil
}

func rpush(cli *CometDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 2 {
		return nil, newWrongNumberOfArgsError("rpush")
	}

	key, value := args[0], args[1]
	res, err := cli.db.RPush(key, value)
	if err != nil {
		return nil, err
	}
	return redcon.SimpleInt(res), nil
}

func lpop(cli *CometDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, newWrongNumberOfArgsError("lpop")
	}

	key := args[0]
	res, err := cli.db.LPop(key)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func rpop(cli *CometDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, newWrongNumberOfArgsError("rpop")
	}

	key := args[0]
	res, err := cli.db.RPop(key)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func llen(cli *CometDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, newWrongNumberOfArgsError("llen")
	}

	key := args[0]
	res, err := cli.db.LLen(key)
	if err != nil {
		return nil, err
	}
	return redcon.SimpleInt(res), nil
}

// =============== ZSet ===============

func zadd(cli *CometDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 3 {
		return nil, newWrongNumberOfArgsError("zadd")
	}

	var ok = 0
	key, score, member := args[0], args[1], args[2]
	res, err := cli.db.ZAdd(key, utils.FloatFromBytes(score), member)
	if err != nil {
		return nil, err
	}
	if res {
		ok = 1
	}
	return redcon.SimpleInt(ok), nil
}

func zscore(cli *CometDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 2 {
		return nil, newWrongNumberOfArgsError("zscore")
	}
	key, field := args[0], args[1]
	return cli.db.ZScore(key, field)
}

// =============== Key ===============
func del(cli *CometDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, newWrongNumberOfArgsError("del")
	}
	key := args[0]
	err := cli.db.Del(key)
	if err != nil {
		return nil, errors.New("fail to delete")
	}
	return redcon.SimpleInt(1), nil
}

func typ(cli *CometDBClient, args [][]byte) (interface{}, error) {
	if len(args) != 1 {
		return nil, newWrongNumberOfArgsError("type")
	}
	key := args[0]
	typ, err := cli.db.Type(key)
	if err != nil {
		return nil, err
	}
	switch typ {
	case redis.String:
		return "string", nil
	case redis.Hash:
		return "hash", nil
	case redis.Set:
		return "set", nil
	case redis.List:
		return "list", nil
	case redis.ZSet:
		return "zset", nil
	default:
		return nil, nil
	}
}
