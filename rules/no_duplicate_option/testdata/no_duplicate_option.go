package fixtures

type Widget struct {
	width  int
	height int
	color  string
}

type Option func(*Widget)

func withWidth(w int) Option {
	return func(widget *Widget) {
		widget.width = w
	}
}

func withHeight(h int) Option {
	return func(widget *Widget) {
		widget.height = h
	}
}

func withColor(c string) Option {
	return func(widget *Widget) {
		widget.color = c
	}
}

func NewWidget(opts ...Option) *Widget {
	w := &Widget{}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func NewWidgetWithName(name string, opts ...Option) *Widget {
	w := &Widget{}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func badDuplicateOption() {
	// withWidth(10) is passed twice (exact same call)
	_ = NewWidget(withWidth(10), withHeight(20), withWidth(10)) // MATCH /duplicated option function argument: withWidth(10)/
}

func badDuplicateOptionWithName() {
	// withWidth is passed twice with a leading non-variadic arg
	_ = NewWidgetWithName("test", withWidth(10), withHeight(20), withWidth(10)) // MATCH /duplicated option function argument: withWidth(10)/
}

func badMultipleDuplicates() {
	// Both withWidth and withColor are duplicated
	_ = NewWidget(withWidth(10), withColor("red"), withWidth(10), withColor("red")) /* MATCH /duplicated option function argument: withWidth(10)/ */ // MATCH /duplicated option function argument: withColor("red")/
}

func goodNoDuplicates() {
	_ = NewWidget(withWidth(10), withHeight(20))
}

func goodSingleOption() {
	_ = NewWidget(withWidth(10))
}

func goodNoOptions() {
	_ = NewWidget()
}

func goodDifferentValues() {
	// Same function name but different args — not flagged because string repr differs
	_ = NewWidget(withWidth(10), withWidth(20))
}

func nonVariadicFunc(a, b int) int {
	return a + b
}

func goodNonVariadic() {
	_ = nonVariadicFunc(1, 1)
}
