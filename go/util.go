package gomfk

func CheckOk(err error) {
	if err != nil {
		panic(err)
	}
}
