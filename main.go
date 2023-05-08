package main

func main() {
	server := NewAPIServer(":3001")
	server.Run()
}
