# kchain
基于tendermint的区块链

k链是基于tendermint底层的联盟链
k链是结合了tendermint,abci,web三者的抽象,共同打包成一个完整的binary
k链现在只有最基础的数据存储功能,后期会添加账户体系等

待完成:
1. 验证节点的管理
2. 智能合约层


## autoenv安装

```
https://github.com/direnv/direnv
https://github.com/direnv/direnv/releases
```

## gopm安装

```
go get -u -v github.com/go-task/task/cmd/task
https://github.com/gpmgo/gopm/
go get -u github.com/gpmgo/gopm
```

## 下载依赖

```
gopm get -s 依赖
gopm get -g 依赖
```
