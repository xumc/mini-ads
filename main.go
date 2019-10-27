package main

func intersectInt32Slice(a []int32, b []int32) []int32 {
	retSlice := make([]int32, 0)
	for i := range a {
		for j := range b {
			if a[i] == b[j] {
				retSlice = append(retSlice, a[i])
			}
		}
	}

	return retSlice
}

func main() {
}
