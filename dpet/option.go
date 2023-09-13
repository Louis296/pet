package dpet

type ParseOptionSet struct {
	notParseData bool
}

type ParseOption func(*ParseOptionSet)

func genParseOption(opts ...ParseOption) *ParseOptionSet {
	option := &ParseOptionSet{}
	for _, opt := range opts {
		opt(option)
	}
	return option
}

func NotParseData() ParseOption {
	return func(set *ParseOptionSet) {
		set.notParseData = true
	}
}
