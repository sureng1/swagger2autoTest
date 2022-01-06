# get start
到根目录
go run main.go


# swagger2autoTest

现状：
query，path这种param.Props 直接是string，name放在Param那里
改进：
改进find，set prop的方法，将param作为单独的一层放进去
2，把param封装成一个独立的object一层。

要求:
1,不管什么类型，都能够marshal出来，{"name":""}
2,不管什么类型，都能set改动