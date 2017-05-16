package hey

func makeKey(keys ...interface{}) interface{} {
	return keys
}

func newUpdateOp(operator string, fieldNo int, value interface{}) interface{} {
	return []interface{}{operator, fieldNo, value}
}

func makeUpdate(ops ...interface{}) interface{} {
	return ops
}
