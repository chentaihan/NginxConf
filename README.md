#
NginxParse

Nginx源码解析，真的是将源码解析出来;
解析出源码中的所有宏、结构体、结构体实例;
从这些信息中我们就能得到，nginx所有自带配置、模块、变量等信息

# 项目的由来

一直想全面的学习一下Nginx，看书看源码，发现ginx代码规范性很好，每个模块一个.C文件实现，每个类型的模块采用同样的结构体，实现流程也很相识，  
于是决定通过代码将这些相似性的地方都摘出来，如Nginx所有配置、所有模块信息；
后来发现越解析越复杂，其实这个项目是在解析C语言，只不过用的是Nginx源码而已；
解析的内容有：
1. 宏
2. 联合体
3. 结构体
4. 预编译
5. 结构体实例
6. 别名

# 为什么用Go语言实现

当时正在学习Go和Nginx，那就用Go来解析Nginx源码吧

# 大致实现

1. 解析源码中所有的ngx\_command\_t，即解析出所有配置
2. 解析源码中所有的ngx\_module\_t，即解析出所有模块
3. 解析源码中所有的ngx\_http\_variable\_t，即解析出所有模块


