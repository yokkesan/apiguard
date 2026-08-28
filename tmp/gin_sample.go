package tmp

func Routes() {
	router.GET("/users", authMiddleware, getUsers)
	router.POST("/users", createUser)
}
