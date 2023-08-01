package leetcode_test

import (
	"fmt"
	"testing"
)


func Test1Solve(t *testing.T) {
	solve1("abcd", "bc")
}

func solve1(a, b string) int {
	alen, blen := len(a), len(b)
	dp := make([][]int, alen)

	for i:=0; i < alen; i ++ {
		dp[i] = make([]int, blen)
	}

	fmt.Println(dp)
	return -1
}
