package ratelimiter

import (
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/companieshouse/ratelimiter.v1/cache"
)

func TestBaseRateLimiter(t *testing.T) {

	Convey("Instantiate with redis pool", t, func() {
		pool := &redis.Pool{
			MaxIdle:     1,
			IdleTimeout: 1,
			MaxActive:   1,
			Dial: func() (redis.Conn, error) {
				c := redigomock.NewConn()
				var e error
				return c, e
			},
		}

		rl := New(pool, nil)
		So(rl, ShouldNotBeNil)
		So(rl, ShouldHaveSameTypeAs, &Limiter{})
		So(rl.cache, ShouldHaveSameTypeAs, &cache.RedisLimiter{Pool: nil})
	})

	Convey("Instantiate with in memory", t, func() {
		rl := New(nil, nil)
		So(rl, ShouldNotBeNil)
		So(rl, ShouldHaveSameTypeAs, &Limiter{})
		So(rl.cache, ShouldHaveSameTypeAs, &cache.InMemoryLimiter{})
	})

	Convey("Unlimited user", t, func() {
		rl := New(nil, nil)
		redigomock.Clear()

		limited, remain, reset, err := rl.Limit("abc", -1, 60)
		So(err, ShouldBeNil)
		So(remain, ShouldEqual, -1)
		So(reset, ShouldEqual, 60)
		So(limited, ShouldBeFalse)
	})

}
