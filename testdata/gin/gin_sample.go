package gin

func Routes() {
	router.GET("/users", authMiddleware, getUsers)
	router.POST("/users", createUser)
}
