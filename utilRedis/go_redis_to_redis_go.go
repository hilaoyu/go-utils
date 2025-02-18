package utilRedis

import (
	"context"
	"fmt"
	"github.com/hilaoyu/go-utils/utils"
	"github.com/redis/go-redis/v9"
	"io"
	"sync"
)

var (
	ctxDefault = context.Background()
)

type RedisGoConn struct {
	client   *redis.Client
	commands []*redis.StatusCmd
	mu       sync.Mutex
	err      error
}

func GoRedisToRedisGoConn(client *redis.Client) *RedisGoConn {
	return &RedisGoConn{
		client: client,
	}
}

func (rgc *RedisGoConn) Close() error {
	rgc.mu.Lock()
	defer rgc.mu.Unlock()

	if nil == rgc.client {
		return rgc.Err()
	}
	rgc.err = rgc.client.Close()
	rgc.client = nil
	if nil == rgc.err {
		rgc.err = fmt.Errorf("redis: closed")
	}
	return rgc.err
}
func (rgc *RedisGoConn) Err() error {
	return rgc.err
}
func (rgc *RedisGoConn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if nil == rgc.client {
		err = rgc.err
		return
	}

	cmdArgs := append([]interface{}{commandName}, args...)
	cmd := redis.NewStatusCmd(ctxDefault, cmdArgs...)
	err = rgc.client.Process(ctxDefault, cmd)
	if nil == err {
		reply, err = cmd.Bytes()
	}

	return
}
func (rgc *RedisGoConn) Send(commandName string, args ...interface{}) error {
	if nil == rgc.client {
		return rgc.err
	}
	cmdArgs := append([]interface{}{commandName}, args...)
	cmd := redis.NewStatusCmd(ctxDefault, cmdArgs...)
	rgc.err = rgc.client.Process(ctxDefault, cmd)

	if nil == rgc.err {
		rgc.mu.Lock()
		rgc.commands = append(rgc.commands, cmd)
		defer rgc.mu.Unlock()
	}

	return rgc.err
}

func (rgc *RedisGoConn) Flush() error {
	return rgc.Err()
}

func (rgc *RedisGoConn) Receive() (reply interface{}, err error) {
	rgc.mu.Lock()
	defer rgc.mu.Unlock()

	if len(rgc.commands) > 0 {
		cmd := utils.SliceShift(&rgc.commands)
		if nil != cmd.Err() {
			err = cmd.Err()
			return
		}
		reply, err = cmd.Bytes()
		return
	}
	return nil, io.EOF
}
