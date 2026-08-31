package validation

type Context struct{}

func (c *Context) Query(key string) string {
	return ""
}

func (c *Context) ShouldBindJSON(value any) error {
	return nil
}

func getUser(c *Context) {
	id := c.Query("id")
	_ = id
}

func createUser(c *Context) {
	var request struct {
		Name string
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		return
	}
}
