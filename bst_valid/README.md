Реализуй функцию:

```
type TreeNode struct {
    Val   int
    Left  *TreeNode
    Right *TreeNode
}

func IsValidBST(root *TreeNode) bool 
```


Пример:

        5
       / \
      3   7
     / \   \
    2   4   8

→ true

А это:

        5
       / \
      3   7
         / \
        4   8

→ false

Ограничения

- Не используй сортировку.

- O(n) по времени.

- O(h) дополнительной памяти допустимо.

И здесь две реализации:
- рекурсивная;
- итеративная через stack.

