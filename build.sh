#!/bin/bash

# 创建exec文件夹
mkdir -p exec

# 编译http/tikstart.go到exec/http文件
go build -o exec/http http/tikstart.go

# 编译rpc/user、video、contact文件夹中的三个go
for dir in user video contact
do
  go build -o "exec/$dir" rpc/"$dir"/"$dir".go
done

# 编译task/server/server.go到exec/task_server
go build -o exec/task_server task/server/server.go

# 编译task/client/client.go到exec/task_client
go build -o exec/task_client task/client/client.go
