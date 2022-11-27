package redis

import (
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

// 协程池子
type redisPool struct {
	pools map[string]*redis.Pool // ipport:池子
	lock  sync.RWMutex
}

var gRedisPools *redisPool

func init() {
	if gRedisPools == nil {
		gRedisPools = &redisPool{
			pools: make(map[string]*redis.Pool),
		}
	}
}

// GetRedisPool 尝试从连接池中获取连接，获取不到，则弄一个
func (p *redisPool) GetRedisPool(conf *RedisConf) *redis.Pool {
	// 尝试获取池子
	p.lock.RLock()
	if pool, ok := p.pools[conf.Addr]; ok {
		p.lock.RUnlock()
		return pool
	}
	p.lock.RUnlock()

	// 创建池子
	pool := &redis.Pool{
		MaxIdle:     conf.MaxIdle,
		MaxActive:   conf.MaxActive,
		IdleTimeout: conf.IdleTimeOut,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.Addr)
			if err != nil {
				return nil, err
			}

			if _, err := c.Do("AUTH", conf.PassWord); err != nil {
				// 秘钥认证： https://blog.csdn.net/corruptwww/article/details/126047077
				// https://help.aliyun.com/document_detail/425605.html
				c.Close()
				return nil, err
			}

			// 选择db
			//c.Do("SELECT", REDIS_DB)   // TODO: 看看要不要
			return c, nil
		},
	}
	p.lock.Lock()
	p.pools[conf.Addr] = pool
	p.lock.Unlock()
	return pool
}

type RedisConf struct {
	Addr           string
	MaxIdle        int
	MaxActive      int
	IdleTimeOut    time.Duration
	RequestTimeOut time.Duration // 请求超时配置
	UserName       string
	PassWord       string
}

type redisClient struct {
	conf *RedisConf
}

func NewRedisClient(conf *RedisConf) *redisClient {
	return &redisClient{
		conf: conf,
	}
}

// 使用连接池中的连接，来请求redis服务端
func (cli *redisClient) Do(commandName string, args ...interface{}) (interface{}, error) {
	pool := gRedisPools.GetRedisPool(cli.conf)
	conn := pool.Get()
	defer conn.Close() // 这个Close只是会清除一些连接状态，恢复成空闲连接，然后放入到pool中，给下一次使用。不是真实关闭。

	return redis.DoWithTimeout(conn, cli.conf.RequestTimeOut, commandName, args...)
}
