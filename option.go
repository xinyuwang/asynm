package asynm

const (
	// never expired
	NoExpiration      int64 = -1
	DefaultExpiration int64 = -2
)

type options struct {
	// 默认过期时间
	expiration int64
}

func newOption() *options {
	return &options{
		expiration: DefaultExpiration,
	}
}

type Option interface {
	apply(*options)
}

// expirattion
type expirationOption int64

func (v expirationOption) apply(opt *options) {
	opt.expiration = int64(v)
}

func WithExpiration(v int64) Option {
	return expirationOption(v)
}
