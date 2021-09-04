package webhook

type HandlerOption uint8

var (
	WithApply = HandlerOption(1)
	WithPlan  = HandlerOption(2)
)
