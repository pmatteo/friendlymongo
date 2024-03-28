package friendlymongo

type unwindOpts struct {
	includeIndex         *string
	preserveNullAndEmpty *bool
}

type unwindOptsFunc func(*unwindOpts)

func IncludeIndex(index string) unwindOptsFunc {

	return func(opts *unwindOpts) {
		opts.includeIndex = &index
	}
}

func PreserveNullEmpty(preserve bool) unwindOptsFunc {

	return func(opts *unwindOpts) {
		opts.preserveNullAndEmpty = &preserve
	}
}
