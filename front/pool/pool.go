package pool

import (
	"context"
	"errors"
	"sync"
	"time"

	"fmt"
	"os"

	"google.golang.org/grpc"
)

var nowFunc = time.Now

var (
	ErrPoolExhausted = errors.New("connection pool exhausted")
	ErrPoolHasClosed = errors.New("connection pool closed")
	ErrConnTimeout   = errors.New("connection status [timeout]")
	ErrConnInvalid   = errors.New("connection status [invalid]")
	ErrConnShutdown  = errors.New("connection status [shutdown]")
	ErrConnFailture  = errors.New("connection status [transfer-error]")

	ErrConnNil = errors.New("nil connection!")
)

type DialFunc func() (*grpc.ClientConn, error)
type TestOnBorrowFunc func(*grpc.ClientConn, time.Time) error
type PoolOption func(*poolOptions)

type poolOptions struct {
	size         int
	capacity     int
	timeout      time.Duration
	testOnBorrow TestOnBorrowFunc
}

func WithPoolSize(size int) PoolOption {
	return func(opt *poolOptions) {
		opt.size = size
		opt.capacity = size + size/2*3
	}
}

func WithPoolCapacity(cap int) PoolOption {
	return func(opt *poolOptions) {
		opt.capacity = cap
	}
}

func WithPoolTimeout(timeout time.Duration) PoolOption {
	return func(opt *poolOptions) {
		opt.timeout = timeout
	}
}

func WithTestOnBorrow(borrowFunc TestOnBorrowFunc) PoolOption {
	return func(opt *poolOptions) {
		opt.testOnBorrow = borrowFunc
	}
}

type PoolStats struct {
	Availabl int
	Capacite int
}

func (ps PoolStats) String() string {
	return fmt.Sprintf("[PoolStat] Capacite: %d, Availabl: %d .", ps.Capacite, ps.Availabl)
}

type Pool struct {
	Dial    DialFunc
	popts   poolOptions
	clients chan ClientConn
	mu      sync.RWMutex
}

type ClientConn struct {
	*grpc.ClientConn
	pool     *Pool
	timeUsed time.Time
}

func (c ClientConn) RpcConnection() (*grpc.ClientConn, error) {
	if c.ClientConn == nil {
		return nil, ErrConnNil
	}
	return c.ClientConn, nil
}

// New alloc grpc.ClientConn pool
// Not thread-safe!
func New(dial DialFunc, opts ...PoolOption) (*Pool, error) {

	p := &Pool{
		Dial: dial,
	}

	for _, opt := range opts {
		opt(&p.popts)
	}

	p.clients = make(chan ClientConn, p.popts.capacity)

	for i := 0; i < p.popts.size; i++ {
		c, err := dial()
		if err != nil {
			return nil, err
		}

		p.clients <- ClientConn{
			ClientConn: c,
			pool:       p,
			timeUsed:   nowFunc(),
		}
	}

	for i := 0; i < p.popts.capacity-p.popts.size; i++ {
		p.clients <- ClientConn{
			pool: p,
		}
	}

	return p, nil
}

func (p *Pool) getClients() chan ClientConn {
	p.mu.RLock()
	clients := p.clients
	p.mu.RUnlock()

	return clients
}

func (p *Pool) Close() {
	p.mu.Lock()
	clients := p.clients
	p.clients = nil
	p.mu.Unlock()

	if clients == nil {
		return
	}

	close(clients)

	for i := 0; i < p.Capacity(); i++ {
		client := <-clients
		if client.ClientConn == nil {
			continue
		}
		client.ClientConn.Close()
	}

}

func (p *Pool) IsClosed() bool {
	return p == nil || p.getClients() == nil
}

func (p *Pool) Get(ctx context.Context) (*ClientConn, error) {
	var err error

	clients := p.getClients()
	if clients == nil {
		return nil, ErrConnInvalid
	}

	cc := ClientConn{
		pool: p,
	}

	select {
	case cc = <-clients:
	case <-ctx.Done():
		return nil, ErrConnTimeout
	}

	cc.timeUsed = nowFunc()
	if cc.ClientConn != nil {
		/*
			idleTimeout := p.popts.timeout
			if idleTimeout > 0 && cc.timeUsed.Add(idleTimeout).Before(nowFunc()) {
				cc.ClientConn.Close()
				cc.ClientConn = nil
			}
		*/

		// fixed on 2018-05-09
		// 测试失败则尝试重建
		testOnBorrow := p.popts.testOnBorrow
		if testOnBorrow != nil && testOnBorrow(cc.ClientConn, cc.timeUsed) != nil {
			cc.ClientConn.Close()
			cc.ClientConn, err = p.Dial()
			if err != nil {
				fmt.Fprintln(os.Stderr, "[E] TestOnBorrow error:", err)
			}
		}

	} else {
		cc.ClientConn, err = p.Dial()
		if err != nil {
			clients <- ClientConn{
				pool: p,
			}
		}
	}

	return &cc, err
}

func (c *ClientConn) Close() error {
	if c == nil {
		return nil
	}

	if c.ClientConn == nil {
		return ErrConnInvalid
	}

	if c.pool.IsClosed() {
		return ErrPoolHasClosed
	}

	cc := ClientConn{
		pool:       c.pool,
		ClientConn: c.ClientConn,
		timeUsed:   nowFunc(),
	}

	select {
	case c.pool.clients <- cc:
	default:
		return ErrPoolExhausted
	}

	return nil
}

func (p *Pool) Capacity() int {
	if p.IsClosed() {
		return 0
	}
	return cap(p.clients)
}

func (p *Pool) Available() int {
	if p.IsClosed() {
		return 0
	}
	return len(p.clients)
}

func (p *Pool) Stats() PoolStats {
	return PoolStats{
		Availabl: p.Available(),
		Capacite: p.Capacity(),
	}
}
