package common

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"

	"google.golang.org/grpc"
)

//TODO:https://zhuanlan.zhihu.com/p/612406500  https://cloud.tencent.com/developer/article/2230478
//TODO: 调研grpc连接池是否必要，长连接短连接，http2多路复用
//连接池性能更好（9w tps vs 40w tps）：https://www.jianshu.com/p/062b5462d3aa   队头阻塞（Head-of-line blocking）。
//多路复用：https://zhuanlan.zhihu.com/p/555785707

//连接池设计：基于GRPC的多路复用、超时重连特性，我们很容易实现GRPC连接池。
//https://zhuanlan.zhihu.com/p/100200985
//https://github.com/shimingyah/pool

//面试官:如何实现一个连接池，我当场懵了
//https://zhuanlan.zhihu.com/p/574549786

//长连接和短连接：https://zhuanlan.zhihu.com/p/433961209

type GRPCInterface interface {
	Get()
	Put(conn *grpc.ClientConn)
	Close()
	NewGRPCPool(o *Options, dialOptions ...grpc.DialOption)
	CreateGrpcConn(data map[string]interface{})
}

// GRPCPool pool info
type GRPCPool struct {
	Mu          sync.Mutex
	IdleTimeout time.Duration
	conns       chan *grpcIdleConn
	factory     func() (*grpc.ClientConn, error)
	close       func(*grpc.ClientConn) error
}

type grpcIdleConn struct {
	conn *grpc.ClientConn
	t    time.Time
}

func init() {
	rand.NewSource(time.Now().UnixNano())
}

// Get get from pool
func (c *GRPCPool) Get() (*grpc.ClientConn, error) {
	c.Mu.Lock()
	conns := c.conns
	c.Mu.Unlock()

	if conns == nil {
		return nil, errClosed
	}
	for {
		select {
		case wrapConn := <-conns:
			if wrapConn == nil {
				return nil, errClosed
			}
			//判断是否超时，超时则丢弃
			if timeout := c.IdleTimeout; timeout > 0 {
				if wrapConn.t.Add(timeout).Before(time.Now()) {
					//丢弃并关闭该链接
					c.close(wrapConn.conn)
					continue
				}
			}
			return wrapConn.conn, nil
		default:
			conn, err := c.factory()
			if err != nil {
				return nil, err
			}

			return conn, nil
		}
	}
}

// Put put back to pool
func (c *GRPCPool) Put(conn *grpc.ClientConn) error {
	if conn == nil {
		return errRejected
	}

	c.Mu.Lock()
	defer c.Mu.Unlock()

	if c.conns == nil {
		return c.close(conn)
	}

	select {
	case c.conns <- &grpcIdleConn{conn: conn, t: time.Now()}:
		return nil
	default:
		//连接池已满，直接关闭该链接
		return c.close(conn)
	}
}

// Close close pool
func (c *GRPCPool) Close() {
	c.Mu.Lock()
	conns := c.conns
	c.conns = nil
	c.factory = nil
	closeFun := c.close
	c.close = nil
	c.Mu.Unlock()

	if conns == nil {
		return
	}

	close(conns)
	for wrapConn := range conns {
		closeFun(wrapConn.conn)
	}
}

// IdleCount idle connection count
func (c *GRPCPool) IdleCount() int {
	c.Mu.Lock()
	conns := c.conns
	c.Mu.Unlock()
	return len(conns)
}

// NewGRPCPool init grpc pool
func NewGRPCPool(o *Options, dialOptions ...grpc.DialOption) (*GRPCPool, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	//init pool
	pool := &GRPCPool{
		conns: make(chan *grpcIdleConn, o.MaxCap),
		factory: func() (*grpc.ClientConn, error) {
			target := o.nextTarget()
			if target == "" {
				return nil, errTargets
			}

			ctx, cancel := context.WithTimeout(context.Background(), o.DialTimeout)
			defer cancel()

			return grpc.DialContext(ctx, target, dialOptions...)
		},
		close:       func(v *grpc.ClientConn) error { return v.Close() },
		IdleTimeout: o.IdleTimeout,
	}

	//danamic update targets
	o.update()

	//init make conns
	for i := 0; i < o.InitCap; i++ {
		conn, err := pool.factory()
		if err != nil {
			pool.Close()
			return nil, err
		}
		pool.conns <- &grpcIdleConn{conn: conn, t: time.Now()}
	}

	return pool, nil
}

func CreateGrpcConn(data map[string]interface{}) (*GRPCPool, error) {

	if _, ok := data["addrs"]; !ok {
		return nil, errors.New("addrs no found in grpc config")
	}

	if _, ok := data["pool_size"]; !ok {
		data["pool_size"] = float64(25) //20230104, 如果不传，默认赋值。推荐引擎
	}

	if _, ok := data["timeout"]; !ok {
		data["timeout"] = float64(1000) //20230104, 如果不传，默认赋值。推荐引擎
	}

	addrs_raw := data["addrs"].([]interface{})
	pool_size := int(data["pool_size"].(float64))
	readTimeout := time.Duration(int64(data["readTimeoutMs"].(float64))) * time.Millisecond
	writeTimeout := time.Duration(int64(data["writeTimeoutMs"].(float64))) * time.Millisecond
	dialTimeout := time.Duration(int64(data["dialTimeoutMs"].(float64))) * time.Millisecond
	idleTimeout := time.Duration(int64(data["idleTimeoutS"].(float64))) * time.Second
	initCap := int(data["initCap"].(float64))

	addrs := make([]string, 0)
	for _, addr := range addrs_raw {
		addrs = append(addrs, addr.(string))
	}

	//这里貌似有问题，传入的类型应该是time.duration
	options := &Options{
		InitTargets:  addrs,
		InitCap:      initCap,
		MaxCap:       pool_size,
		DialTimeout:  dialTimeout,
		IdleTimeout:  idleTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	grpc_pool, err := NewGRPCPool(options, grpc.WithInsecure())
	return grpc_pool, err
}

var (
	errClosed   = errors.New("pool is closed")
	errInvalid  = errors.New("invalid config")
	errRejected = errors.New("connection is nil. rejecting")
	errTargets  = errors.New("targets server is empty")
)

// Options pool options
type Options struct {
	lock sync.RWMutex
	//targets node
	targets *[]string
	//targets channel
	input chan *[]string

	//InitTargets init targets
	InitTargets []string
	// init connection
	InitCap int
	// max connections
	MaxCap       int
	DialTimeout  time.Duration
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Input is the input channel
func (o *Options) Input() chan<- *[]string {
	return o.input
}

// update targets
func (o *Options) update() {
	//init targets
	o.targets = &o.InitTargets

	go func() {
		for targets := range o.input {
			if targets == nil {
				continue
			}

			o.lock.Lock()
			o.targets = targets
			o.lock.Unlock()
		}
	}()

}

// NewOptions returns a new newOptions instance with sane defaults.
func NewOptions() *Options {
	o := &Options{}
	o.InitCap = 5
	o.MaxCap = 100
	o.DialTimeout = 5 * time.Second
	o.ReadTimeout = 5 * time.Second
	o.WriteTimeout = 5 * time.Second
	o.IdleTimeout = 60 * time.Second
	return o
}

// validate checks a Config instance.
func (o *Options) validate() error {
	if o.InitTargets == nil ||
		o.InitCap <= 0 ||
		o.MaxCap <= 0 ||
		o.InitCap > o.MaxCap ||
		o.DialTimeout == 0 ||
		o.ReadTimeout == 0 ||
		o.WriteTimeout == 0 {
		return errInvalid
	}
	return nil
}

// nextTarget next target implement load balance
func (o *Options) nextTarget() string {
	o.lock.RLock()
	defer o.lock.RUnlock()

	tlen := len(*o.targets)
	if tlen <= 0 {
		return ""
	}

	//rand server
	return (*o.targets)[rand.Int()%tlen]
}
