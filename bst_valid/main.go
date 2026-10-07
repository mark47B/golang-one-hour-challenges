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
	// "math"
)

type TreeNode struct {
    Val   int
    Left  *TreeNode
    Right *TreeNode
}

func IsValidBST(root *TreeNode) bool { // Через рекурсию
	if root == nil {
		return true
	}
	l := isValidBST(root.Left, nil, &root.Val)
	r := isValidBST(root.Right, &root.Val, nil)
	return l && r
}

func isValidBST(root *TreeNode, min *int, max *int) bool {
	if root == nil {
		return true
	}
	if min != nil {
		if root.Val < *min {
			return false
		}
	}
	if max != nil {
		if root.Val > *max {
			return false
		}
	}
	l := isValidBST(root.Left, min, &root.Val)
	r := isValidBST(root.Right, &root.Val, max)
	return l && r
}



func IsValidBST(root *TreeNode) bool { // через стек
	if root == nil {
		return true
	}
	stack := make([]*TreeNode, 0)
	current := root
	var prev *int

	for current != nil || len(stack) > 0 {
		for current != nil {
			stack = append(stack, current)
			current = current.Left
		}

		current = stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if prev != nil && current.Val <= *prev {
			return false
		}
		prevVal := current.Val
		prev = &prevVal

		current = current.Right
	}
	return true
}

func main(){
	a := &TreeNode{
		Val: 4,
		Left: nil,
		Right: nil,
	}
	b := &TreeNode{
		Val: 8,
		Left: nil,
		Right: nil,
	}

	c := &TreeNode{
		Val: 7,
		Left: a,
		Right: b,
	}
	d := &TreeNode{
		Val: 3,
		Left: nil,
		Right: nil,
	}
	root := &TreeNode{
		Val: 5,
		Left: d,
		Right: c,
	}

	fmt.Println(IsValidBST(root))

}

