package redis

import (
	"context"
	"fmt"
	"reflect"

	"github.com/akatsukisun2020/wx_common/logger"
)

// KVRedisClient kv数据结构封装
type KVRedisClient struct {
	redisCli *redisClient
}

func NewKVRedisClient(conf *RedisConf) *KVRedisClient {
	return &KVRedisClient{
		redisCli: NewRedisClient(conf),
	}
}

func (cli *KVRedisClient) SetData(ctx context.Context, key string, data []byte) error {
	replySet, err := cli.redisCli.Do("SET", key, data)
	if err != nil {
		logger.Errorf("SET error: %v\n", err)
		return err
	}

	logger.Debugf("in SetData, key:%s, replySet%v", key, replySet)
	return nil
}

func (cli *KVRedisClient) GetData(ctx context.Context, key string) ([]byte, error) {
	reply, err := cli.redisCli.Do("GET", key)
	if err != nil {
		logger.Errorf("GET error, key:%s, err: %v\n", key, err)
		return []byte{}, err
	}

	if reply == nil {
		return []byte{}, nil
	}

	ret, ok := reply.([]byte)
	if !ok {
		logger.Errorf("GetData error, key:%s, reply.type:%v", key, reflect.TypeOf(reply))
		return []byte{}, fmt.Errorf("type error")
	}

	return ret, nil
}
