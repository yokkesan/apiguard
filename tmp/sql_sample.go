package tmp

func ExampleSQL(id string) {
	query := "SELECT * FROM users WHERE id=" + id
	_ = query
}
