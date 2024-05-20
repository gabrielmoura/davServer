package ternary

func Or(a, b interface{}) interface{} {
	if a != nil {
		return a
	}
	return b
}
func OrString(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
func Ternary(statement bool, a, b interface{}) interface{} {
	if statement {
		return a
	}
	return b
}
func And(a, b interface{}) interface{} {
	if a == nil {
		return nil
	}
	return b
}
