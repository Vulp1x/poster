package instagrapi

type Configuration struct {
	Hostname string
}

// Default sets default values in config variables.
func (c *Configuration) Default() {
	c.Hostname = "http://0.0.0.0:8000"
}
