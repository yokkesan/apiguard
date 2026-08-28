package auth

func authMiddleware(c *Context) {
	token := c.GetHeader("Authorization")
	_ = token
}

func RoutesWithAuth() {
	router.GET("/private", authMiddleware, getPrivate)
	router.GET("/public", getPublic)
}
