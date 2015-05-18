// +build integration

package ratelimiter

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	. "github.com/smartystreets/goconvey/convey"
)

var clientKey string
var storedClientKey string

func TestIntegrationRedisRateLimit(t *testing.T) {

	pool := &redis.Pool{
		MaxIdle:     1,
		IdleTimeout: 1,
		MaxActive:   10,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "localhost:6379")
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}

	rl := New(pool, nil)

	conn := pool.Get()
	defer conn.Close()

	d, err := time.ParseDuration("1m")
	if err != nil {
		panic(err)
	}

	Convey("Integration tests", t, func() {

		Convey("Ensure redis available and setup", func() {
			rand.Seed(time.Now().UTC().UnixNano())
			clientKey = "RATELIMITER_TEST:" + "#" + strconv.Itoa(rand.Intn(100))
			storedClientKey = "RateLimit:" + clientKey
			Printf("Using test user id [%s]\n", clientKey)

			_, err := conn.Do("SETEX", "TESTREDISUP", true, 1)
			So(err, ShouldBeNil)
		})

		Convey("Instantiated OK", func() {
			So(rl, ShouldHaveSameTypeAs, &Limiter{})
		})

		Convey("Uncached user", func() {
			exceeded, remain, reset, err := rl.Limit(clientKey, 1, d)
			So(err, ShouldBeNil)
			So(exceeded, ShouldBeFalse)
			So(remain, ShouldEqual, 0)
			So(reset.Seconds(), ShouldBeLessThanOrEqualTo, 60000) // >= to allow for delay in execution
		})

		Convey("Expired user", func() {
			conn.Do("SET", storedClientKey, 0)
			exceeded, _, _, err := rl.Limit(clientKey, 1, d)
			So(err, ShouldBeNil)
			So(exceeded, ShouldBeTrue)
		})

		Convey("Unlimited user", func() {
			exceeded, _, _, err := rl.Limit(clientKey, -1, d)
			So(err, ShouldBeNil)
			So(exceeded, ShouldBeFalse)
		})

		Reset(func() {
			conn.Do("DEL", storedClientKey)
		})

	})
}
