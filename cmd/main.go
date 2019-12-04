package main


func main() {
	if err := apiserver.Start(); err != nil {
		log.Fatal(err)
	}
}
