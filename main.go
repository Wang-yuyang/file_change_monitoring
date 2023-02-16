package file_change_monitoring

func main() {
	hello(On)
}

func hello(on func(msg string)) {
	on("hello")
}

func On(msg string) {
	println(msg)
}
