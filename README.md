# Web框架设计学习

### 核心

1.代表服务器的抽象：Server

2.代表上下文的抽象：Context

3.路由树

### 路由树设计

1.已经注册的路由，无法被覆盖

2.path必须以/开头，不能以/结尾

3.不能在同一个位置同时注册路径参数和通配符匹配

4.同名路径参数，在路径匹配时，值会被覆盖

5.路由树不是线程安全的，要求用户在启动HTTPServer前注册完路由

6.路由树算法：前缀树

7.路由匹配优先级：静态路由>路径参数>通配符匹配

8.路由查找不支持回溯

9.路由树组织方式：一个HTTP方法对应一棵路由树

### Context

1.context不是线程安全的，也不需要线程安全

2.原生API不可以重复读取HTTP协议的body内容，但是我们可以通过封装来允许重复读取，核心步骤是我们将body读取出来之后放到一个地方，后续都从这个地方读

3.原生API不可以修改HTTP协议的响应

4.Form和PostForm区别：

Form：URL里面的查询参数和PATCH，POST，PUT的表单数据。 所有表单数据都可以拿到。

PostForm：PATCH，POST，PUT body参数，在编码是x-www-form-urlencoded的时候才能拿到

5.路径参数的支持：web框架发现匹配上某个路径参数之后，将这段路径记录下来，作为路径参数的值，默认为string，用户可以根据需要转化为不同的类型

### AOP方案----Middleware：

(Aspect Oriented Programming)，面对切面编程。核心在于将横向关注点从业务中剥离出来。

横向关注点：就是那些跟业务没啥关系，但是每个业务都必须要处理。例如：

1.可观测性：logging ,metric 和tracing

2.安全相关：登录，鉴权和权限控制

3.错误处理：例如错误页面支持

4.可用性保证：熔断限流和降级等

在不同的框架，不用的语言中AOP 方案的叫法不同：Middleware，Handler，Chain,Filter,Filter-Chain,Interceptor,Wrapper...

洋葱模式：形如洋葱，拥有一个核心。核心一般为业务逻辑，使用middleware层层包裹核心，可以无侵入式地增强核心功能，或者解决AOP问题

责任链模式：不同的handler 组成一条链，链条的每一环都有自己的功能，一方面可以用责任链模式将复杂逻辑分成链条上的不同步骤，另一方面也可以灵活地在链条上添加新的handler

常见的可观测性框架：OpenTelemetry，SkyWalking，Prometheus...

集成可观测性框架：利用middleware机制

