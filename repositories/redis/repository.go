package redis

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/vrnvu/tiny-to-go/tiny"
)

type redisRepository struct {
	client *redis.Client
}

func newRedisClient(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewRedisRepository(redisURL string) (tiny.Repository, error) {
	repo := &redisRepository{}
	client, err := newRedisClient(redisURL)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewRedisRepository")
	}
	repo.client = client
	return repo, nil
}

func (r *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (r *redisRepository) Find(code string) (*tiny.Redirect, error) {
	redirect := &tiny.Redirect{}
	key := r.generateKey(code)
	data, err := r.client.HGetAll(key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	if len(data) == 0 {
		return nil, errors.Wrap(tiny.ErrRedirectNotFound, "repository.Redirect.Find")
	}
	timestamp, err := strconv.ParseInt(data["timestamp"], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	redirect.Code = data["code"]
	redirect.URL = data["url"]
	redirect.Timestamp = timestamp
	return redirect, nil
}

func (r *redisRepository) Save(redirect *tiny.Redirect) error {
	key := r.generateKey(redirect.Code)
	data := map[string]interface{}{
		"code":      redirect.Code,
		"url":       redirect.URL,
		"timestamp": redirect.Timestamp,
	}
	_, err := r.client.HMSet(key, data).Result()
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}
