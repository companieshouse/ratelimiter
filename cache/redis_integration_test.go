// +build integration

package cache

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	. "github.com/smartystreets/goconvey/convey"
)

// func TestMain(m *testing.M) {
//
// 	setup()
// 	ret := m.Run()
// 	teardown()
// 	os.Exit(ret)
// }

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

	rl := &RedisLimiter{Pool: pool}

	conn := pool.Get()
	defer conn.Close()

	Convey("Ensure redis available and setup", t, func() {
		rand.Seed(time.Now().UTC().UnixNano())
		clientKey = "RATELIMITER_TEST:" + "#" + strconv.Itoa(rand.Intn(100))
		storedClientKey = "RateLimit:" + clientKey
		Printf("Using test user id [%s]\n", clientKey)

		_, err := conn.Do("SETEX", "TESTREDISUP", true, 1)
		So(err, ShouldBeNil)

		conn.Do("DEL", storedClientKey)
	})

	Convey("Instantiated OK", t, func() {
		So(rl, ShouldHaveSameTypeAs, &RedisLimiter{})
	})

	Convey("Uncached user", t, func() {
		exceeded, remain, reset, err := rl.Limit(clientKey, 1, 60)
		So(err, ShouldBeNil)
		So(exceeded, ShouldBeFalse)
		So(remain, ShouldEqual, 0)
		So(reset, ShouldBeLessThanOrEqualTo, 60000) // >= to allow for delay in execution
		conn.Do("DEL", storedClientKey)
	})

	Convey("Expired user", t, func() {
		conn.Do("SET", storedClientKey, 0)
		exceeded, _, _, err := rl.Limit(clientKey, 1, 60)
		So(err, ShouldEqual, "Rate limit exceeded")
		So(exceeded, ShouldBeTrue)
		conn.Do("DEL", storedClientKey)
	})

	conn.Do("DEL", storedClientKey) // cleanup

}
