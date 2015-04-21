package ratelimiter

import (
	"github.com/companieshouse/gotools/log"
	"github.com/garyburd/redigo/redis"
)

var logger log.Glogger // Replace with generic via interface

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

// RateLimiter implements a redis backed rate limiter
func RateLimiter(identity string, limit int, window int, conn redis.Conn) (rateLimitExceeded bool, remaining int, reset int, e error) {

	logger.Debug("Rate limiting for identity: [%s] Limit: [%s] Window: [%s]", identity, limit, window)

	var r int64
	var err error

	if limit == -1 {
		return false, limit, window, nil // Unlimited user
	}

	r, err = redis.Int64(rlScript.Do(conn, "RateLimit:"+identity, limit, window))
	logger.Debug("Get and Decrement rate limit for identity: [%s] Remaining: [%s]", identity, r)

	if err != nil && err.Error() != errRateLimitExceeded {
		// Unexpected error
		logger.Error("Error: %s", err.Error())
		return true, 0, 0, err // Assume rate limit reached if there's an error to stop glitches causing someone to become unlimited
	}

	t, pttlErr := redis.Int64(conn.Do("PTTL", "RateLimit:"+identity))
	if pttlErr != nil {
		logger.Error("Error: %s", pttlErr.Error())
		return true, 0, 0, pttlErr // Assume rate limit reached if there's an error to stop glitches causing someone to become unlimited
	}

	if err != nil && err.Error() == errRateLimitExceeded {
		// Rate limit exceeded
		return true, int(r), int(t), err
	}

	return false, int(r), int(t), err
}
