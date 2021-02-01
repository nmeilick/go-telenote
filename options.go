package telenote

type Options struct {
	DisableWebPreview bool
	ParseMode         string
}

func NewOptions() *Options {
	return &Options{
		DisableWebPreview: false,
		ParseMode:         "Markdown",
	}
}

type Option func(*Options) error

func Preview() Option {
	return func(o *Options) error {
		o.DisableWebPreview = false
		return nil
	}
}

func NoPreview() Option {
	return func(o *Options) error {
		o.DisableWebPreview = true
		return nil
	}
}

func ParseMode(mode string) Option {
	return func(o *Options) error {
		o.ParseMode = mode
		return nil
	}
}
