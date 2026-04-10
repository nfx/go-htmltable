package htmltable

const (
	baseUserAgent    = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) "
	engine           = "AppleWebKit/537.36 (KHTML, like Gecko) "
	browser          = "Chrome/121.0.0.0 Safari/537.36 "
	packageInfo      = "nfx/go-htmltable (+https://github.com/nfx/go-htmltable)"
	DefaultUserAgent = baseUserAgent + engine + browser + packageInfo
)

type options struct {
	userAgent string
	innerHTML bool
}

type Option func(*options)

// WithUserAgent sets the User-Agent header used when fetching URLs
func WithUserAgent(ua string) Option {
	return func(o *options) {
		o.userAgent = ua
	}
}

// WithInnerHTML instructs the parser to keep the inner HTML of each cell
func WithInnerHTML() Option {
	return func(o *options) {
		o.innerHTML = true
	}
}

func applyOptions(opts []Option) options {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}
