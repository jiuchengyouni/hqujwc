# hqu-jwc

## gin+grpc+etcd 实现
华侨大学新版教务处查成绩学分绩点排名等不同功能的_WEUcookie之间的交换和获得，基于https://github.com/Zakiaatot/hqu_ual_interface
实现，正在逐渐补充中

## 项目结构
```
hqujwc/
├── app                   // 各个微服务
│   ├── gateway           // 网关
│   ├── WX                // 微信模块微服务
│   └── jwc               // 教务处模块微服务
├── config                // 配置文件
├── idl                   // protoc文件
│   └── pb                // 放置生成的pb文件
├── pkg                   // 各种包
│   └── discovery         // etcd服务注册、keep-alive、获取服务信息等等
└── types                 // 定义各种通用结构体
```

## jwc 教务处模块
```
jwc/
├── cmd                   // 启动入口
├── service               // 业务服务
├── utils                 // ******教务处功能获取
└── repository            // 持久层(未完成)
    └── db                // 视图层
        ├── dao           // 对数据库进行操作
        └── model         // 定义数据库的模型
```

## wx 微信基本操作模块
```
wx/
├── cmd                   // 启动入口
├── service               // 业务服务
├── utils                 // wx操作工具
└── repository            // 持久层(未完成)
    └── db                // 视图层
        ├── dao           // 对数据库进行操作
        └── model         // 定义数据库的模型
```