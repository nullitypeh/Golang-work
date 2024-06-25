# Golang-work
6.24：项目基本完成，四个接口经检测可以使用
6.25：已经完成了任务，在6.24的基础上补充了时区。但是接口实现分布式锁的时候用到了SET的扩展命令（SET EX PX NX），这个命令可能会导致锁过期释放了，业务还没执行完或者锁被别的线程误删。可能考虑继续优化。
![{DBB7438D-30F4-447d-91BE-68A55C7E5A4C}](https://github.com/nullitypeh/Golang-work/assets/159264792/3bbbd912-d7c0-4f00-bbcf-427024000862)
![{A78014C7-2B18-4a12-BDDE-A1504DAE6FD8}](https://github.com/nullitypeh/Golang-work/assets/159264792/76286858-9c16-4813-bb2b-d17d24813f6a)
![{F54F0C46-E240-4b2d-B730-8E8631F76032}](https://github.com/nullitypeh/Golang-work/assets/159264792/8526c06d-8f1d-4e59-b572-32fff8451539)
![image](https://github.com/nullitypeh/Golang-work/assets/159264792/d4000e93-c409-429f-83a9-d8b24f10ade8)
