# ratelimiter

_"No running in the API!"_

GO based ratelimiter

Usage
-----

Using redis caching

```go
pool := &redis.Pool{ ... }

limiter := NewRateLimiter(pool)

hasExceededLimit, allowanceRemaining, secondsToReset, err := limiter.Limit(identify, limit, window)
```

With in-memory fake caching

```go
limiter := NewRateLimiter()

hasExceededLimit, allowanceRemaining, secondsToReset, err := limiter.Limit(identify, limit, window)
```
