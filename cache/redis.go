package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// Here is the definition of the cache structure.
type Redis struct {
	chunk    *redis.Client
	state    *redis.Client
	addr     string
	password string
	size     int
}

func NewRedis(addr string, password string, size int) *Redis {
	return &Redis{
		chunk:    nil,
		state:    nil,
		addr:     addr,
		password: password,
		size:     size,
	}
}

// InitConnection initialize redis connection settings.
func (r *Redis) InitConnection() error {
	r.chunk = redis.NewClient(&redis.Options{
		Addr:     r.addr,
		Password: r.password,
		DB:       0,
		PoolSize: r.size,
	})
	r.state = redis.NewClient(&redis.Options{
		Addr:     r.addr,
		Password: r.password,
		DB:       1,
		PoolSize: r.size,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.chunk.Ping(ctx).Result()
	if err != nil {
		return err
	}
	_, err = r.state.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}

// CloseConnection close the connection with the cache.
func (r *Redis) CloseConnection() error {
	err := r.chunk.Close()
	return err
}

// InsertOneChunk is a function insert one chunk in the cache.
func (r *Redis) InsertOneChunk(hash string, chunk Chunk) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := r.chunk.SAdd(ctx, hash, chunk.Index).Result()
	if err != nil {
		return err
	}
	return nil
}

// GetChunkList will return this file's all chunks.
func (r *Redis) GetChunkList(hash string, name string) (*[]Chunk, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := r.chunk.SMembers(ctx, hash).Result()
	if err != nil {
		return nil, err
	}
	chunks := make([]Chunk, 0)
	for _, val := range res {
		index, _ := strconv.Atoi(val)
		chunks = append(chunks, Chunk{
			Name:  name,
			Hash:  hash,
			Index: index,
		})
	}
	return &chunks, nil
}

//	RemoveAllRecords is a function which run after merge or error.
func (r *Redis) RemoveAllRecords(hash string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := r.chunk.Del(ctx, hash).Result()
	if err != nil {
		return err
	}
	_, err = r.state.Del(ctx, hash).Result()
	if err != nil {
		return err
	}
	return nil
}

// ChangeFileState init the file's state for the server check.
func (r *Redis) ChangeFileState(hash string, saved bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := r.state.Set(ctx, hash, saved, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

// CheckFileUpload check the state of the file.
func (r *Redis) CheckFileUpload(hash string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := r.state.Get(ctx, hash).Result()
	if err != nil {
		return false, err
	}
	exist, _ := strconv.ParseBool(res)
	return exist, nil
}
