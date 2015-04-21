package ratelimiter

import (
	"errors"
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	. "github.com/smartystreets/goconvey/convey"
)

func setDial(c *redis.Pool, f func() (redis.Conn, error)) {
	c.Dial = f
}

func TestIntegrationRateLimiter(t *testing.T) {

	conn := redigomock.NewConn()

	Convey("Unlimited user", t, func() {
		redigomock.Clear()

		limited, remain, reset, err := RateLimiter("abc", -1, 60, conn)
		So(err, ShouldBeNil)
		So(remain, ShouldEqual, -1)
		So(reset, ShouldEqual, 60)
		So(limited, ShouldBeFalse)
	})

	Convey("Rate limit exceeded", t, func() {
		redigomock.Clear()
		redigomock.Command("EVALSHA").Expect(int64(0))
		redigomock.Command("EVALSHA").ExpectError(errors.New("Rate limit exceeded"))
		redigomock.Command("PTTL").Expect(int64(5))

		limited, remain, reset, err := RateLimiter("abc", 10, 60, conn)
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

		limited, remain, reset, err := RateLimiter("abc", 10, 60, conn)
		So(err, ShouldBeNil)
		So(remain, ShouldEqual, 9)
		So(reset, ShouldEqual, 5)
		So(limited, ShouldBeFalse)
	})

	Convey("Error retrieving window remaining", t, func() {
		redigomock.Clear()
		redigomock.Command("PTTL", "RateLimit:abc").ExpectError(errors.New("An error"))
		redigomock.Command("EVALSHA").Expect(int64(0))

		_, _, _, err := RateLimiter("abc", 10, 60, conn)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "An error")
	})

}
