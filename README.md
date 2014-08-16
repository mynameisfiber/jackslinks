Jacks Links
===========

> Non blocking doubly linked list with O(1) move and O(num_cursors) update

Based entirely on [Non-Blocking Doubly-Linked Lists with Good Amortized
Complexity](http://arxiv.org/abs/1408.1935) by Niloufar Shafiei.

How do I use this?
------------------

First you create an initial node:

```
node := jackslinks.NewNode("value")
```

And then you create your cursor over that node:

```
cursor := jackslinks.Newcursor(node)
```

And now you can get rockin!  Add and remove nodes with `InsertBefore` and
`InsertAfter`.  You can move back and forth with `MoveLeft` and `MoveRight`.
In addition, you can go to the top of the list with `MoveToHead`.  And, of
course, there is the obligitory `Delete` method to remove the current node.

```
cursor.InsertAfter("another thing")   # cursor -> "value"
cursor.InsertBefore("one last thing") # cursor -> "value"
cursor.MoveToHead()                   # cursor -> "one last thing"
cursor.Delete()                       # cursor -> "value"
cursor.MoveRight()                    # cursor -> "another thing"

fmt.Println("Current value under the cursor: ", cursor.Get().(string))
```

In addition, you can make as many cursors over the same list as you want and
just keep operating over them without worrying about locks,

```
for i := 0; i < 100; i++ {
    go func(id int) {
        local_cursor := cursor(node)
        for {
            local_cursor.InsertBefore(id)
        }
    }(i)
}
```
