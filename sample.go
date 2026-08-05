package main

func main() {
	id := "1"

	query := "SELECT * FROM users WHERE id=" + id

	println(query)
}
