```
type TreeNode struct {
    Val   int
    Left  *TreeNode
    Right *TreeNode
}
```

Реализуй:
```
func LevelOrder(root *TreeNode) [][]int
```

Для дерева:
```
        3
       / \
      9   20
         /  \
        15   7
```

нужно получить:
```
[][]int{
    {3},
    {9, 20},
    {15, 7},
}
```
Требования:

- O(n) времени;
- O(n) памяти в худшем случае;
- без рекурсии;
- каждый уровень должен оказаться отдельным []int
