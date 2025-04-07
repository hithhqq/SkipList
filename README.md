## 跳表

跳表是基于二维单向链表实现的，使用二分查找的思想，查找的时间复杂度是log(n)级别的。数据结构如下图所示：

![img](https://i-blog.csdnimg.cn/blog_migrate/6276df0c6fadbac855ba40d058a84b6b.png)

跳表的节点在插入时，通过随机数的形式得到出现在哪一层。