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

**Input**

  - **identity** - (string) a unique string to identify the user for which you're ratelimiting
  - **limit** - (int) number of requests allowable within window
  - **window** - (time.Duration) length of the window

**Output**

  - **exceeded** - (boolean) true if rate limit has been exceeded
  - **remaining** - (int) number of requests still allowable in current window
  - **reset** - (time.Duration) length of time until window resets
  
With in-memory fake caching (**nb** Doesn't perform actual rate limiting. Provides a simulation for testing and novelty purposes)

```go
limiter := NewRateLimiter()

exceeded, remaining, reset, err := limiter.Limit(identity, limit, window)
remaining, err := limiter.QueryLimit(identity)
```
