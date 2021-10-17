package webhook

type HandlerOption uint8

const (
	WithApply HandlerOption = iota
	WithPlan
)
