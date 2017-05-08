package driver

import (
	"github.com/garyburd/redigo/redis"
	"sync"
	"time"
)

var (
	rdsMu sync.Mutex
	Rds   = make(map[int]*RedisDriver)
)

type RedisDriver struct {
	Connected bool
	cfg       *RedisConfig
	Pool      *redis.Pool
}

type RedisConfig struct {
	DbHost            string `yaml:"dbHost"`
	DbPwd             string `yaml:"dbPwd"`
	DbName            int    `yaml:"dbName"`
	DbMaxIdleConns    int    `yaml:"dbMaxIdleConns"`
	DbConnMaxLifeTime int    `yaml:"dbConnMaxLifeTime"`
}

func RegisterRedis(name int, cfg *RedisConfig) (err error) {
	rdsMu.Lock()
	defer rdsMu.Unlock()
	md := new(RedisDriver)
	md.cfg = cfg
	md.cfg.DbName = name
	err = md.Register(cfg)
	if err != nil {
		return
	}
	Rds[name] = md
	return
}

func GetRedis(name int) (conn redis.Conn) {
	md := Rds[name]
	if md == nil {
		return
	}
	if !md.Connected {
		md.Connect()
	}
	conn = md.Pool.Get()
	return
}

func (p *RedisDriver) Register(cfg *RedisConfig) (err error) {
	//todo cfg的基础判断
	p.cfg = cfg
	return
}

func (p *RedisDriver) Connect() {
	p.Pool = &redis.Pool{
		MaxIdle:     100,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", p.cfg.DbHost)
			if err != nil {
				return nil, err
			}
			if p.cfg.DbPwd != "" {
				if _, err := c.Do("AUTH", p.cfg.DbPwd); err != nil {
					c.Close()
					panic("Redis AUTH Error" + err.Error())
					return nil, err
				}
			}
			if p.cfg.DbName > 0 {
				if _, err := c.Do("SELECT", p.cfg.DbName); err != nil {
					c.Close()
					panic("Redis SEELCT Error" + err.Error())
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
