package output

// DiscardingOutput discards all output, which can be useful for testing, among other purposes.
type DiscardingOutput struct{ noopOutput }

func (o DiscardingOutput) WithValues(keysAndValues ...interface{}) Output { return o }
func (o DiscardingOutput) V(level int) Output                             { return o }

// Convention used to verify, at compile time, that DiscardingOutput implements the Output interface.
var _ Output = DiscardingOutput{}
