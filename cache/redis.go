package cache

import "github.com/garyburd/redigo/redis"

// RedisLimiter defines a redis backed rate limiter implementation
type RedisLimiter struct {
	Pool *redis.Pool
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
func (rl *RedisLimiter) Limit(identity string, limit int, window int) (rateLimitExceeded bool, remaining int, reset int, lastError error) {

	logger.Debug("Rate limiting for identity: [%s] Limit: [%d] Window: [%d]", identity, limit, window)

	var r int64
	var err error

	conn := rl.Pool.Get()
	defer conn.Close()

	rateLimitExceeded = false

	r, err = redis.Int64(rlScript.Do(conn, "RateLimit:"+identity, limit, window, nil))
	logger.Debug("Get and Decrement rate limit for identity: [%s] Remaining: [%d] Window: [%d]", identity, r, window)

	if err != nil && err.Error() != errRateLimitExceeded {
		rateLimitExceeded, lastError = handleUnexpected(err)
		return
	}
	remaining = int(r)

	t, pttlErr := redis.Int64(conn.Do("PTTL", "RateLimit:"+identity))
	if pttlErr != nil {
		rateLimitExceeded, lastError = handleUnexpected(pttlErr)
		return
	}
	reset = int(t)

	if err != nil && err.Error() == errRateLimitExceeded {
		logger.Debug("Rate limit exceeded for identity: [%s] Time to reset: [%s]", identity, t)
		rateLimitExceeded = true
		lastError = err
	}

	return
}

// QueryLimit allows querying of the current remaining limit for an identity
func (rl *RedisLimiter) QueryLimit(identity string) (remain int, err error) {

	conn := rl.Pool.Get()
	defer conn.Close()

	remain64, err := redis.Int64(conn.Do("GET", "RateLimit:"+identity))
	if err != nil {
		logger.Error("ID [%s] Failed to fetch limit remaining", identity)
		return
	}

	return int(remain64), nil
}
