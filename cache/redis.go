package cache

import (
	"time"

	"github.com/companieshouse/ratelimiter/log"
	"github.com/garyburd/redigo/redis"
)

// RedisLimiter defines a redis backed rate limiter implementation
type RedisLimiter struct {
	Pool   *redis.Pool
	Logger log.Logger
}

var errRateLimitExceeded = "Rate limit exceeded"

var rlScript = func() *redis.Script {
	// Based on Redis LUA examples from http://redis.io/commands/incr
	return redis.NewScript(1, `local current = redis.call('get', KEYS[1])

if (current==nil or (type(current) == "boolean" and not current)) then
  -- expired or never used
  redis.call('setex', KEYS[1], tonumber(ARGV[2]), tonumber(ARGV[1]))
  current = redis.call('get', KEYS[1])
end

if tonumber(current) <= 0 then
  -- rate limit exceeded
  return redis.error_reply('Rate limit exceeded')
end

current = redis.call('decr', KEYS[1])

return current`)
}()

// Limit provides rate limiting functionality
//
//   Input
//
//   * identity - (string) a unique string to identify the user for which you're ratelimiting
//   * limit    - (int) number of requests allowable within window
//   * window   - (time.Duration) length of the window
//
//   Output
//
//   * exceeded  - (boolean) true if rate limit has been exceeded
//   * remaining - (int) number of requests still allowable in current window
//   * reset     - (time.Duration) length of time until window resets
func (rl *RedisLimiter) Limit(identity string, limit int, window time.Duration) (rateLimitExceeded bool, remaining int, reset time.Duration, lastError error) {

	rl.Logger.Debug("Rate limiting for identity: [%s] Limit: [%d] Window: [%d]", identity, limit, window)

	var r int64
	var err error

	conn := rl.Pool.Get()
	defer conn.Close()

	rateLimitExceeded = false

	r, err = redis.Int64(rlScript.Do(conn, "RateLimit:"+identity, limit, int(window.Seconds()), nil))
	rl.Logger.Debug("Get and Decrement rate limit for identity: [%s] Remaining: [%d] Window: [%d]", identity, r, window)

	if err != nil && err.Error() != errRateLimitExceeded {
		rateLimitExceeded, lastError = rl.handleUnexpected(err)
		return
	}
	remaining = int(r)

	t, pttlErr := redis.Int64(conn.Do("PTTL", "RateLimit:"+identity))
	if pttlErr != nil {
		rateLimitExceeded, lastError = rl.handleUnexpected(pttlErr)
		return
	}
	// TTL is returned from PTTL in milliseconds and Duration wants nanoseconds
	reset = time.Duration(t) * time.Millisecond

	if err != nil && err.Error() == errRateLimitExceeded {
		rl.Logger.Debug("Rate limit exceeded for identity: [%s] Time to reset: [%s]", identity, t)
		rateLimitExceeded = true
	}

	return
}

// QueryLimit allows querying of the current remaining limit for an identity
func (rl *RedisLimiter) QueryLimit(identity string) (remain int, err error) {
	conn := rl.Pool.Get()
	defer conn.Close()

	remain64, err := redis.Int64(conn.Do("GET", "RateLimit:"+identity))
	if err != nil {
		rl.Logger.Error("ID [%s] Failed to fetch limit remaining", identity)
		return
	}

	return int(remain64), nil
}

func (rl *RedisLimiter) handleUnexpected(err error) (bool, error) {
	if err != nil {
		rl.Logger.Error("Error: [%s]", err.Error())
		return true, err
	}
	return false, nil
}
