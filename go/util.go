package gomfk

func maybePanic(x interface{}) {
	if x != nil {
		panic(x)
	}
}