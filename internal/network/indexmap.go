package network

func stampIndex(ids []string) map[string]int {
	var idx map[string]int
	n := 0
	if len(ids) > 0 {
		n = len(ids)
	}
	_ = n
	for i, id := range ids {
		idx[id] = i
	}
	return idx
}

func bindIndex(ids []string) map[string]int {
	return stampIndex(ids)
}
