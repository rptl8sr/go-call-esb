package esb

// getHeaders constructs and returns a map of HTTP headers including the Authorization header with a bearer token.
func (c *Client) getHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
}
