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

func LevelOrder(root *TreeNode) [][]int {
	levelOrder := make([][]int, 0, 0)
	if root == nil {
		return levelOrder
	}
	currentLevel := make([]*TreeNode, 0, 0)
	currentLevel = append(currentLevel, root)
	

	level := 0
	for len(currentLevel) != 0 {
		lenCurrent := len(currentLevel)
		levelOrder = append(levelOrder, make([]int, 0, lenCurrent))
		for i:=0; i < lenCurrent; i++ {
			levelOrder[level] = append(levelOrder[level], currentLevel[i].Val)
			if currentLevel[i].Right != nil {
				currentLevel = append(currentLevel, currentLevel[i].Right)
			}
			if currentLevel[i].Left != nil {
				currentLevel = append(currentLevel, currentLevel[i].Left)
			}
			
		}
		currentLevel = currentLevel[lenCurrent:]
		level++
	}
	return levelOrder
}


func main(){
	a := &TreeNode{
		Val: 15,
		Left: nil,
		Right: nil,
	}
	b := &TreeNode{
		Val: 7,
		Left: nil,
		Right: nil,
	}

	c := &TreeNode{
		Val: 20,
		Left: a,
		Right: b,
	}
	d := &TreeNode{
		Val: 9,
		Left: nil,
		Right: nil,
	}
	root := &TreeNode{
		Val: 3,
		Left: d,
		Right: c,
	}

	fmt.Println(LevelOrder(root))

}
