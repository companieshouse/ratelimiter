package cache

import (
	"errors"
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRedisRateLimit(t *testing.T) {

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

	rl := &RedisLimiter{Pool: pool}

	Convey("Rate limit exceeded", t, func() {
		redigomock.Clear()
		redigomock.Command("EVALSHA").Expect(int64(0))
		redigomock.Command("EVALSHA").ExpectError(errors.New("Rate limit exceeded"))
		redigomock.Command("PTTL").Expect(int64(5))

		limited, remain, reset, err := rl.Limit("abc", 10, 60)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "Rate limit exceeded")
		So(remain, ShouldEqual, 0)
		So(reset, ShouldEqual, 5)
		So(limited, ShouldBeTrue)
	})

	Convey("Rate limit not exceeded", t, func() {
		redigomock.Clear()
		redigomock.Command("EVALSHA").Expect(int64(9))
		redigomock.Command("PTTL", "RateLimit:abc").Expect(int64(5))

		limited, remain, reset, err := rl.Limit("abc", 10, 60)
		So(err, ShouldBeNil)
		So(remain, ShouldEqual, 9)
		So(reset, ShouldEqual, 5)
		So(limited, ShouldBeFalse)
	})

	Convey("Error retrieving window remaining", t, func() {
		redigomock.Clear()
		redigomock.Command("PTTL", "RateLimit:abc").ExpectError(errors.New("An error"))
		redigomock.Command("EVALSHA").Expect(int64(0))

		_, _, _, err := rl.Limit("abc", 10, 60)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "An error")
	})

}
