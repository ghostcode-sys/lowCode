package utility

func HandlePanic() {
	if r := recover(); r != nil {
		println("Recovered from panic:", r)
	}
}
