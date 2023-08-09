package abis

func isBatch(raw []byte) bool {
	for _, c := range raw {
		if c == 0x20 || c == 0x09 || c == 0x0a || c == 0x0d {
			continue
		}
		return c == '['
	}
	return false
}
