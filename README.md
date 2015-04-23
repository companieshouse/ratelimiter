# ratelimiter

_"No running in the API!"_

GO based ratelimiter

Usage
-----

Using redis caching

```go
pool := &redis.Pool{ ... }

limiter := NewRateLimiter(pool)

exceeded, remaining, reset, err := limiter.Limit(identity, limit, window)
```

With in-memory fake caching

```go
limiter := NewRateLimiter()

exceeded, remaining, reset, err := limiter.Limit(identity, limit, window)
```
