package cache

import (
	"testing"

	"github.com/companieshouse/ratelimiter/generic"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInMemoryRateLimit(t *testing.T) {

	log := &generic.DefaultLogger{}
	rl := &InMemoryLimiter{Logger: log}

	Convey("Always return unlimited", t, func() {
		limited, remain, reset, err := rl.Limit("abc", 10, 60)
		So(err, ShouldBeNil)
		So(remain, ShouldEqual, 10)
		So(reset, ShouldEqual, 60)
		So(limited, ShouldBeFalse)
	})

	Convey("QueryLimit", t, func() {
		remain, err := rl.QueryLimit("abc")
		So(err, ShouldBeNil)
		So(remain, ShouldEqual, 0)
	})

}
