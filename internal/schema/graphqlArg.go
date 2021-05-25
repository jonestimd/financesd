package schema

// AsMaps converts a []interface{} to a []map[string]interface{}.
func asMaps(arg interface{}, nested ...string) []map[string]interface{} {
	items := arg.([]interface{})
	maps := make([]map[string]interface{}, len(items))
	for i, m := range items {
		maps[i] = m.(map[string]interface{})
		for _, name := range nested {
			if value, ok := maps[i][name]; ok {
				maps[i][name] = asMaps(value)
			}
		}
	}
	return maps
}
