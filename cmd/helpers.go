package cmd

import "net/url"

// buildQuery converts a map of string key/value pairs into url.Values,
// skipping empty values.
func buildQuery(params map[string]string) url.Values {
	q := url.Values{}
	for k, v := range params {
		if v != "" {
			q.Set(k, v)
		}
	}
	return q
}
