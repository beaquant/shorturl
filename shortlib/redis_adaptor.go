/*************************************************************************
  > File Name: RedisAdaptor.go
  > Author: Wu Yinghao
  > Mail: wyh817@gmail.com
  > Created Time: 二  6/ 9 15:29:05 2015
 ************************************************************************/
package shortlib

import (
	//"errors"
	"fmt"
	"github.com/btstar/shorturl/models"
	"github.com/garyburd/redigo/redis"
)

type RedisAdaptor struct {
	conn   redis.Conn
	config *models.Redis
}

const SHORT_URL_COUNT_KEY = "short_url_count"

func NewRedisAdaptor(config *models.Redis) (*RedisAdaptor, error) {
	redis_cli := &RedisAdaptor{}
	redis_cli.config = config

	host := config.Redishost
	port := config.Redisport

	connStr := fmt.Sprintf("%v:%v", host, port)
	fmt.Printf(connStr)
	conn, err := redis.Dial("tcp", connStr)
	if err != nil {
		return nil, err
	}

	redis_cli.conn = conn

	return redis_cli, nil
}

func (r *RedisAdaptor) Release() {
	r.conn.Close()
}

func (r *RedisAdaptor) InitCountService() error {
	_, err := r.conn.Do("SET", SHORT_URL_COUNT_KEY, 0)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisAdaptor) NewShortUrlCount() (int64, error) {
	count, err := redis.Int64(r.conn.Do("INCR", SHORT_URL_COUNT_KEY))
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RedisAdaptor) SetUrl(shortUrl, originalUrl string) error {
	key := fmt.Sprintf("short:%v", shortUrl)
	_, err := r.conn.Do("SET", key, originalUrl)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisAdaptor) GetUrl(shortUrl string) (string, error) {
	key := fmt.Sprintf("short:%v", shortUrl)
	originalUrl, err := redis.String(r.conn.Do("GET", key))
	if err != nil {
		return "", err
	}
	return originalUrl, nil
}
