package redis

import (
	"context"

	"github.com/akatsukisun2020/wx_common/logger"
	"github.com/gomodule/redigo/redis"
)

// ZsetRedisClient zset数据结构封装
type ZsetRedisClient struct {
	redisCli *redisClient
}

func NewZsetRedisClient(conf *RedisConf) *ZsetRedisClient {
	return &ZsetRedisClient{
		redisCli: NewRedisClient(conf),
	}
}

type ZsetItem struct {
	Member string
	Score  int64
}

// ZAdd 加入zset
func (cli *ZsetRedisClient) ZAdd(ctx context.Context, key string, itemList []*ZsetItem) error {
	var args []interface{}
	args = append(args, key)
	for _, item := range itemList {
		args = append(args, item.Score)
		args = append(args, item.Member)
	}

	replySet, err := cli.redisCli.Do("ZADD", args...)
	if err != nil {
		logger.Errorf("ZADD error, key:%s, itemList:%v, err: %v", key, itemList, err)
		return err
	}
	logger.Debugf("rsp:%v", replySet)
	return nil
}

// Zrangebyscore 返回有序集合中指定分数区间的成员列表
func (cli *ZsetRedisClient) Zrangebyscore(ctx context.Context, key string, min, max, count int64) ([]*ZsetItem, error) {
	args := []interface{}{key, min, max, "WITHSCORES", "LIMIT", 0, count} // TODO：测试一下，是不是从min开始的index
	//args := []interface{}{key, min, max, "WITHSCORES"}

	replySet, err := redis.Strings(cli.redisCli.Do("ZRANGEBYSCORE", args...))
	if err != nil {
		logger.Errorf("ZRANGEBYSCORE error, key:%s, min:%d, max:%d, count:%d, err: %v", key, min, max, count, err)
		return []*ZsetItem{}, err
	}

	if len(replySet) == 0 {
		return []*ZsetItem{}, nil
	}

	res := make([]*ZsetItem, 0, len(replySet)/2)
	for i := 0; i < len(replySet); i += 2 {
		member := replySet[i]
		scoreFloat, err := redis.Float64([]byte(replySet[i+1]), nil)
		if err != nil {
			logger.Errorf("Zrangebyscore error, replySet[i]:%v", replySet[i])
		}
		res = append(res, &ZsetItem{Member: member, Score: int64(scoreFloat)})
	}

	return res, nil
}

// Zremrangebyscore 加入zset
func (cli *ZsetRedisClient) Zremrangebyscore(ctx context.Context, key string, min, max int64) error {
	replySet, err := cli.redisCli.Do("ZREMRANGEBYSCORE", key, min, max)
	if err != nil {
		logger.Errorf("ZADD error, key:%s, min:%d, max:%d, err: %v", key, min, max, err)
		return err
	}
	logger.Debugf("rsp:%v", replySet)
	return nil
}
