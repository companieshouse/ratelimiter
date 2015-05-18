package cache

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInMemoryRateLimit(t *testing.T) {

	rl := &InMemoryLimiter{}

	Convey("Always return unlimited", t, func() {
		limited, remain, reset, err := rl.Limit("abc", 10, 60)
		So(err, ShouldBeNil)
		So(remain, ShouldEqual, 10)
		So(reset, ShouldEqual, 60)
		So(limited, ShouldBeFalse)
	})

	Convey("QueryLimit", t, func() {
		remain, window, err := rl.QueryLimit("abc", 10, 60)
		So(err, ShouldBeNil)
		So(remain, ShouldEqual, 10)
		So(window, ShouldEqual, 60)
	})

}
