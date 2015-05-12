# ratelimiter

_"No running in the API!"_

Go based ratelimiter with redis backend

Usage
-----

Using redis caching

```go
pool := &redis.Pool{ ... }

limiter := NewRateLimiter(pool)

// Decrement and get current limit
exceeded, remaining, reset, err := limiter.Limit(identity, limit, window)

// Just fetch current remaining limit
remaining, err := limiter.QueryLimit(identity)
```

With in-memory fake caching (**nb** Doesn't perform actual rate limiting. Provides a simulation for testing and novelty purposes)

```go
limiter := NewRateLimiter()

exceeded, remaining, reset, err := limiter.Limit(identity, limit, window)
remaining, err := limiter.QueryLimit(identity)
```
