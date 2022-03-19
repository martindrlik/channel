package channel_test

import (
	"context"
	"testing"
	"time"

	"github.com/martindrlik/channel"
)

func TestChannel(t *testing.T) {
	foo := must(channel.Open("foo"))
	bar := must(channel.Open("bar"))
	defer foo.Close()
	defer bar.Close()
	go mustSend(foo, "bar", "hello bar")
	go mustSend(bar, "foo", "hello foo")
	receive := func(c *channel.Channel, to, expect string) {
		actual := string(must(c.Receive()))
		if actual != expect {
			t.Errorf("%s expected to receive %q, got %q", to, expect, actual)
		}
	}
	receive(foo, "foo", "hello foo")
	receive(bar, "bar", "hello bar")
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

func mustSend(c *channel.Channel, to, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := c.Send(ctx, to, []byte(message))
	if err != nil {
		panic(err)
	}
}
