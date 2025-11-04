package config

type Options struct {
	placeholderReplace bool // Enable placeholder replacement when rewriting the request path.
	skipNamespace      bool // Skip adding namespace prefix to metric names.
}

func (opts *Options) SetPlaceholderReplace(flag bool) {
	opts.placeholderReplace = flag
}

func (opts *Options) EnablePlaceholderReplace() bool {
	return opts.placeholderReplace
}

func (opts *Options) SetSkipNamespace(flag bool) {
	opts.skipNamespace = flag
}

func (opts *Options) SkipNamespace() bool {
	return opts.skipNamespace
}
