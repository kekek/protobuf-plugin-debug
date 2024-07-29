# protobuf 插件goland调试

在go开发中，goland的debug调试功能可以帮我们直观的看到代码运行的过程，对开发起到事半功倍的作用；

由于 protobuf 的go插件使用的是protoc运行，而不是使用go直接运行的，因此不能直接使用goland的调试功能

有没有其它的方式能在goland中使用 protobuf 调试呢 ？ 本文将给出一种方式

阅读本文，假设您已经对 protobuf 插件的开发有一定的了解。

###  说明

下面以开发一个 protoc-gen-test 为例，演示如何使用goland 调试。

本项目的目录结构
``` 
.
├── api         // proto 文件，及根据proto生成的文件
│ └── person
│     └── person.proto
├── cmd         // 命令
│ └── protoc-gen-debug 
│     └── main.go
│ └── protoc-gen-test
│     └── main.go
├── go.mod
├── internal      
│ └── messager
│     ├── messager.go
│     └── packager.go
└── readme.md

```

本文的目的是演示debug调试， 所以在这里我们设定：protoc-gen-test 的功能是生成 api/person/person.proto下所有 message 都实现 messager.Messager 接口的golang代码

此项目中有两个二进制文件
- protoc-gen-debug

这个命令是为了生成 protoc-gen-test 的输入参数

- protoc-gen-test 

这个是我们最终需要的调试的程序

### 思路

- 正常的流程

    正常的流程 完成 protoc-gen-test的编码后，使用如下命令来生成文件

```sh 
 
 protoc -I=./ --test_out=./ --plugin=protoc-gen-test=/path/to/protoc-gen-test   api/person/person.proto

```

此过程中 protoc 解析 api/person/person.proto 为AST树，然后传给 protoc-gen-test 生成相应的代码，此过程无法使用 goland来进行调试
 
- 调试流程

因此我们使用中间程序，proto-gen-debug，将 protoc 解析 api/person/person.proto 解析的ast树，序列化后，保存为参数文件，

然后将参数文件作为输入传给 protoc-gen-test，直接将protoc-gen-test 作为一个golang执行程序，即可使用 goland来调试

``` bash 

[//]: # ( proto-gen-debug 生成参数文件 testdata/request.pb.bin)

protoc -I ./ --debug_out=./ --plugin=protoc-gen-debug=/path/to/protoc-gen-debug --debug_opt=file_binary=testdata/request.pb.bin

[//]: # (testdata/request.pb.bin 参数文件传给 proto-gen-test，开启调试)

./bin/protoc-gen-test < testdata/request.pb.bin

```


### 操作过程

##### 1. 使用 protoc-gen-debug 生成一个二进制文件 request.pb.bin

- 安装 protoc-gen-debug

```bash 
 go build -o bin/protoc-gen-debug cmd/protoc-gen-debug/main.go
 ```

- 生成文件 testdata/request.pb.bin

request.pb.bin 的作用是替代 protoc 生成的AST(抽象语法)树

[protoc-gen-debug使用方法](https://github.com/pubg/protoc-gen-debug)

_protoc-gen-debug中也说明了一种在goland使用debug的方法，实测不可用_

下面是生成 request.pb.bin 的命令，注意替换 {/path/to/protoc-gen-debug} 为真实路径
```bash 
protoc -I ./ \
          --debug_out=./ \
          --plugin=protoc-gen-debug=./bin/protoc-gen-debug \
          --debug_opt=dump_binary=true \
          --debug_opt=dump_json=true \
          --debug_opt=file_binary=testdata/request.pb.bin \
          --debug_opt=file_json=testdata/request.pb.json \
          --debug_opt=parameter=expose_all=true:foo=bar \
          api/person/person.proto
```

![requtest-pb-bin.png](png%2Frequtest-pb-bin.png)

##### 2. 编写代码，然后编译出 一个未经过优化的插件执行文件  protoc-gen-test

此处省略代码过程

生成未经优化的插件

```bash 
 go build -o bin/protoc-gen-test -gcflags "all=-N -l" cmd/protoc-gen-test/main.go
 
 ```

##### 3. 使用go的调试工具，启动 dlv 远程调试

在本机 :2345 端口启动 dlv调试

此处会使用到 第1步生成的request.pb.bin文件

```bash 
 dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./bin/protoc-gen-test < testdata/request.pb.bin
```

##### 4.设置goland的 go-remote 调试

Run->Edit Configurations->➕(Add New Configuration) Go Remote -> 配置好 host和 port， 至此即可在 goland中打断点 （啊，真好）

![add.png](png%2Fadd.png)

![config.png](png%2Fconfig.png)

![debug.png](png%2Fdebug.png)


_注意：此调试过程并不会生成 相应的文件。文件的生成功能是由protoc负责完成的_

##### 5. 调试完成后，生成正式的文件

``` bash 

 go build -o bin/protoc-gen-test cmd/protoc-gen-test/main.go

 protoc -I=./ --test_out=paths=source_relative:. --plugin=protoc-gen-test=./bin/protoc-gen-test  api/person/person.proto
```