package main

func FlagBitIsSet(flag byte, bitIndex int) bool {
	return (flag>>bitIndex)&1 == 1
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
