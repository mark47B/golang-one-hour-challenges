package main

import ( 
//"runtime"
// "time"
 "fmt"
// "context"
// _ "log"
// "errors"
// _ "net/http/pprof"
// "net/http"
	// "sync"
	// _ "sync/atomic"
)

func lengthOfLongestSubstring(s string) int {
	if len(s) == 0 {
		return 0
	}

	var l, r = 0, 0
	maxLength := 0
	
	unique := make(map[rune]int)
	sRunes := []rune(s)

	for r < len(sRunes) {
		if pos, ok := unique[sRunes[r]]; ok {
			if pos >= l {
				l = pos + 1
			}
		}
		if r - l + 1 > maxLength {
			maxLength = r - l + 1
		}
		unique[sRunes[r]] = r
		r++
	}

	return maxLength
}

func main(){
	fmt.Println(lengthOfLongestSubstring("abcdabcbbcdefgimnbvcxzlkjhgfdsab"))
}
