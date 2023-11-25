package gqlt

func assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}
