package dpet

type ParseOptionSet struct {
	notParseData bool
	onlyHeader   bool
}

type ParseOption func(*ParseOptionSet)

func genParseOption(opts ...ParseOption) *ParseOptionSet {
	option := &ParseOptionSet{}
	for _, opt := range opts {
		opt(option)
	}
	return option
}

// NotParseData 不解析数据区，将数据区解压并直接写入dataset中的buffer
func NotParseData() ParseOption {
	return func(set *ParseOptionSet) {
		set.notParseData = true
	}
}

// OnlyParseHeader 只解析文件头，解析产生的dataset不包含数据信息
func OnlyParseHeader() ParseOption {
	return func(set *ParseOptionSet) {
		set.onlyHeader = true
	}
}
