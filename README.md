# Golang-work
6.24：项目基本完成，四个接口经检测可以使用
6.25：已经完成了任务，在6.24的基础上补充了时区。但是接口实现分布式锁的时候用到了SET的扩展命令（SET EX PX NX），这个命令可能会导致锁过期释放了，业务还没执行完或者锁被别的线程误删。可能考虑继续优化。
![{DBB7438D-30F4-447d-91BE-68A55C7E5A4C}](https://github.com/nullitypeh/Golang-work/assets/159264792/3bbbd912-d7c0-4f00-bbcf-427024000862)
