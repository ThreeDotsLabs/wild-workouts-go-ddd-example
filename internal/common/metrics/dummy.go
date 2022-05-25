package metrics

type NoOp struct{}

func (d NoOp) Inc(_ string, _ int) {
	// todo - add some implementation!
}
