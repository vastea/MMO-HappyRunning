# MMO-HappyRunning

## 一、简介

MMO-HappyRunning是基于[myzinx](https://github.com/vastea/myzinx)框架开发的MMO游戏服务器，该项目主要学习使用，因此只实现了简单的功能。

## 二、基础概念扫盲

1. aoi 兴趣点算法，aoi，即 Area Of Interest。

## 三、服务器应用基础协议

> Client 和 Server，谁下面有协议名，说明是谁发送的。比如，MsgID 为 1，SyncPid 在 Server 下，说明这个消息应该是 Server 发送给
> Client 的。

| MsgID |   Client    |   Server   |                     描述                      |
|:-----:|:-----------:|:----------:|:-------------------------------------------:|
|   1   |      -      |  SyncPid   |             同步玩家本次登录的ID(用来标识玩家)             |
|   2   |    Talk     |     -      |                    世界聊天                     |
|   3   | MovePackage |     -      |                     移动                      |
|  200  |      -      | BroadCast  | 广播消息(TP 1-世界聊天 2-出生点坐标同步 3-动作 4-移动之后坐标信息更新) |
|  201  |      -      |  SyncPid   |              广播消息，掉线/aoi消失在视野               |
|  202  |      -      | SyncPlayer |              同步周围的人的位置信息(包括自己)              |