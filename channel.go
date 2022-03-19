package channel

import (
	"context"
	"sync"
)

var channels = struct {
	m map[string]*Channel
	sync.Mutex
}{
	m: map[string]*Channel{},
}

type Channel struct {
	Ch chan []byte
}

func Open(name string) (*Channel, error) {
	defer channels.Unlock()
	channels.Lock()
	c, ok := channels.m[name]
	if !ok {
		c = new(Channel)
		c.Ch = make(chan []byte)
		channels.m[name] = c
	}
	return c, nil
}

func (ch *Channel) Send(ctx context.Context, to string, b []byte) error {
	c, ok := tryGetChannel(to)
	if !ok {
		return nil
	}
	select {
	case c <- b:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

func tryGetChannel(name string) (chan []byte, bool) {
	defer channels.Unlock()
	channels.Lock()
	c, ok := channels.m[name]
	if ok {
		return c.Ch, true
	}
	return nil, false
}

func (ch *Channel) Receive() ([]byte, error) {
	return <-ch.Ch, nil
}

func (ch *Channel) Close() error {
	return nil
}
