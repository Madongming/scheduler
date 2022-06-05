# 调度器demo

## [需求描述](./Describe.md)

## [设计文档](./designs/design.md)
详细的设计，包括数据结构，api的设计
[详细](./designs/implement_v1.md)

## 操作命令
- docker
```shell
make run-docker
```shell

- test
```shell
make test
```

- build
```shell
make build
```
也可以增加-race检查冲突
```shell
make build-check-race
```

- 运行
```shell
bin/schedule
```
## 现状及限制
- 只是实现了函数级别的api，后续可以自由对接cmd活着http api
- 数据存储层是简单实现的文件(json)记录的数据
- 一些任务插件还没有做数据同步到状态，需要根据业务找一个合适的方式，在此demo就没有先实现
- 自动执行的任务没有被记录在history中，因为觉得放进去不久就会占满。应该有一个更好的管理方式。

## 扩展及优化
- 分离store和schedule(对应代码中的package)层，扩展的时候可以在中间加上cache，定时同步到store中。
- 这个cache可以映射为schedule中的对象，在cache中做全局锁，之后可以横向扩张无状态的schedule，提高任务的并发量
- 对于schedule中的关键函数需要做benchmark测试
- 使用trace分析gc占用，内存分配情况
