

## **内存相关**
### [内存分配原理](https://www.topgoer.cn/docs/gozhuanjia/gozhuanjiachapter044.1-memory_alloc)
### [垃圾回收原理](https://www.topgoer.cn/docs/gozhuanjia/chapter044.2-garbage_collection)
### [逃逸分析](https://www.topgoer.cn/docs/gozhuanjia/chapter044.3-escape_analysis)
### [Go语言的内存模型及堆的分配管理](https://zhuanlan.zhihu.com/p/76802887)
### 谈谈内存泄露，什么情况下内存会泄露？怎么定位排查内存泄漏问题？
答：go 中的内存泄漏一般都是 goroutine 泄漏，就是 goroutine 没有被关闭，或者没有添加超时控制，让 goroutine 一只处于阻塞状态，不能被 GC。  
**内存泄露有下面一些情况**  
1）如果 goroutine 在执行时被阻塞而无法退出，就会导致 goroutine 的内存泄漏，一个 goroutine 的最低栈大小为 2KB，在高并发的场景下，对内存的消耗也是非常恐怖的。  
2）互斥锁未释放或者造成死锁会造成内存泄漏  
3）time.Ticker 是每隔指定的时间就会向通道内写数据。作为循环触发器，必须调用 stop 方法才会停止，从而被 GC 掉，否则会一直占用内存空间。  
4）字符串的截取引发临时性的内存泄漏

```go
func main() {
    var str0 = "12345678901234567890"
    str1 := str0[:10]
}
```

5）切片截取引起子切片内存泄漏

```go
func main() {
    var s0 = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
    s1 := s0[:3]
}
```

6）函数数组传参引发内存泄漏【如果我们在函数传参的时候用到了数组传参，且这个数组够大（我们假设数组大小为 100 万，64 位机上消耗的内存约为 800w 字节，即 8MB 内存），或者该函数短时间内被调用 N 次，那么可想而知，会消耗大量内存，对性能产生极大的影响，如果短时间内分配大量内存，而又来不及 GC，那么就会产生临时性的内存泄漏，对于高并发场景相当可怕。】  
**排查方式：**  
一般通过 pprof 是 Go 的性能分析工具，在程序运行过程中，可以记录程序的运行信息，可以是 CPU 使用情况、内存使用情况、goroutine 运行情况等，当需要性能调优或者定位 Bug 时候，这些记录的信息是相当重要。  
**当然你能说说具体的分析指标更加分咯，有的面试官就喜欢他问什么，你简洁的回答什么，不喜欢巴拉巴拉详细解释一通，比如虾P面试官，不过他考察的内容特别多，可能是为了节约时间。**

### golang 的内存逃逸吗？什么情况下会发生内存逃逸？（必问）
答：1)本该分配到栈上的变量，跑到了堆上，这就导致了内存逃逸。2)栈是高地址到低地址，栈上的变量，函数结束后变量会跟着回收掉，不会有额外性能的开销。3)变量从栈逃逸到堆上，如果要回收掉，需要进行 gc，那么 gc 一定会带来额外的性能开销。编程语言不断优化 gc 算法，主要目的都是为了减少 gc 带来的额外性能开销，变量一旦逃逸会导致性能开销变大。  
**内存逃逸的情况如下：**  
1）方法内返回局部变量指针。  
2）向 channel 发送指针数据。  
3）在闭包中引用包外的值。  
4）在 slice 或 map 中存储指针。  
5）切片（扩容后）长度太大。  
6）在 interface 类型上调用方法。

### 请简述 Go 是如何分配内存的？
<font style="color:rgb(36, 41, 46);">Golang内存分配是个相当复杂的过程，其中还掺杂了GC的处理，这里仅仅对其关键数据结构进行了说明，了解其原理而又不至于深陷实现细节。</font>

1. <font style="color:rgb(36, 41, 46);">Golang程序启动时申请一大块内存，并划分成spans、bitmap、arena区域</font>
2. <font style="color:rgb(36, 41, 46);">arena区域按页划分成一个个小块</font>
3. <font style="color:rgb(36, 41, 46);">span管理一个或多个页</font>
4. <font style="color:rgb(36, 41, 46);">mcentral管理多个span供线程申请使用</font>
5. <font style="color:rgb(36, 41, 46);">mcache作为线程私有资源，资源来源于mcentral</font>

****

### [go内存分配器](https://zhuanlan.zhihu.com/p/410317967)
### Channel 分配在栈上还是堆上？哪些对象分配在堆上，哪些对象分配在栈上？
Channel 被设计用来实现协程间通信的组件，其作用域和生命周期不可能仅限于某个函数内部，所以 golang 直接将其分配在堆上  
准确地说，你并不需要知道。Golang 中的变量只要被引用就一直会存活，存储在堆上还是栈上由内部实现决定而和具体的语法没有关系。  
知道变量的存储位置确实和效率编程有关系。如果可能，Golang 编译器会将函数的局部变量分配到函数栈帧（stack frame）上。然而，如果编译器不能确保变量在函数 return 之后不再被引用，编译器就会将变量分配到堆上。而且，如果一个局部变量非常大，那么它也应该被分配到堆上而不是栈上。  
当前情况下，如果一个变量被取地址，那么它就有可能被分配到堆上,然而，还要对这些变量做逃逸分析，如果函数 return 之后，变量不再被引用，则将其分配到栈上。

### 介绍一下大对象小对象，为什么小对象多了会造成 gc 压力？
小于等于 32k 的对象就是小对象，其它都是大对象。一般小对象通过 mspan 分配内存；大对象则直接由 mheap 分配内存。通常小对象过多会导致 GC 三色法消耗过多的 CPU。优化思路是，减少对象分配。  
小对象：如果申请小对象时，发现当前内存空间不存在空闲跨度时，将会需要调用 nextFree 方法获取新的可用的对象，可能会触发 GC 行为。  
大对象：如果申请大于 32k 以上的大对象时，可能会触发 GC 行为。

## 编译
### [逃逸分析是怎么进行的](http://golang.design/go-questions/compile/escape/)
在编译原理中，分析指针动态范围的方法称之为逃逸分析。通俗来讲，当一个对象的指针被多个方法或线程引用时，我们称这个指针发生了逃逸。

Go语言的逃逸分析是编译器执行静态代码分析后，对内存管理进行的优化和简化，它可以决定一个变量是分配到堆还栈上。

写过C/C++的同学都知道，调用著名的malloc和new函数可以在堆上分配一块内存，这块内存的使用和销毁的责任都在程序员。一不小心，就会发生内存泄露。

Go语言里，基本不用担心内存泄露了。虽然也有new函数，但是使用new函数得到的内存不一定就在堆上。堆和栈的区别对程序员“模糊化”了，当然这一切都是Go编译器在背后帮我们完成的。

Go语言逃逸分析最基本的原则是：如果一个函数返回对一个变量的引用，那么它就会发生逃逸。

简单来说，编译器会分析代码的特征和代码生命周期，Go中的变量只有在编译器可以证明在函数返回后不会再被引用的，才分配到栈上，其他情况下都是分配到堆上。

Go语言里没有一个关键字或者函数可以直接让变量被编译器分配到堆上，相反，编译器通过分析代码来决定将变量分配到何处。

对一个变量取地址，可能会被分配到堆上。但是编译器进行逃逸分析后，如果考察到在函数返回后，此变量不会被引用，那么还是会被分配到栈上。

编译器会根据变量是否被外部引用来决定是否逃逸：

> 1. 如果函数外部没有引用，则优先放到栈中；
> 2. 如果函数外部存在引用，则必定放到堆中；
>

Go的垃圾回收，让堆和栈对程序员保持透明。真正解放了程序员的双手，让他们可以专注于业务，“高效”地完成代码编写。把那些内存管理的复杂机制交给编译器，而程序员可以去享受生活。

逃逸分析这种“骚操作”把变量合理地分配到它该去的地方。即使你是用new申请到的内存，如果我发现你竟然在退出函数后没有用了，那么就把你丢到栈上，毕竟栈上的内存分配比堆上快很多；反之，即使你表面上只是一个普通的变量，但是经过逃逸分析后发现在退出函数之后还有其他地方在引用，那我就把你分配到堆上。

如果变量都分配到堆上，堆不像栈可以自动清理。它会引起Go频繁地进行垃圾回收，而垃圾回收会占用比较大的系统开销（占用CPU容量的25%）。

堆和栈相比，堆适合不可预知大小的内存分配。但是为此付出的代价是分配速度较慢，而且会形成内存碎片。栈内存分配则会非常快。栈分配内存只需要两个CPU指令：“PUSH”和“RELEASE”，分配和释放；而堆分配内存首先需要去找到一块大小合适的内存块，之后要通过垃圾回收才能释放。

通过逃逸分析，可以尽量把那些不需要分配到堆上的变量直接分配到栈上，堆上的变量少了，会减轻分配堆内存的开销，同时也会减少gc的压力，提高程序的运行速度。

### [GoRoot 和 GoPath 有什么用](http://golang.design/go-questions/compile/gopath/)
GoRoot 是 Go 的安装路径。mac 或 unix 是在 `/usr/local/go` 路径上，来看下这里都装了些什么：

![](https://golang.design/go-questions/compile/assets/1.png#id=a7Ono&originHeight=144&originWidth=2074&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=)

bin 目录下面：

![](https://golang.design/go-questions/compile/assets/2.png#id=pbBTl&originHeight=142&originWidth=1576&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=)

pkg 目录下面：

![](https://golang.design/go-questions/compile/assets/3.png#id=iycbe&originHeight=464&originWidth=2302&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=)

Go 工具目录如下，其中比较重要的有编译器 `compile`，链接器 `link`：

![](https://golang.design/go-questions/compile/assets/4.png#id=kVJYb&originHeight=184&originWidth=1624&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=)

GoPath 的作用在于提供一个可以寻找 `.go` 源码的路径，它是一个工作空间的概念，可以设置多个目录。Go 官方要求，GoPath 下面需要包含三个文件夹：

```go
src
pkg
bin
```

src 存放源文件，pkg 存放源文件编译后的库文件，后缀为 `.a`；bin 则存放可执行文件。

### [Go 编译链接过程概述](http://golang.design/go-questions/compile/link-process/)
Go 程序并不能直接运行，每条 Go 语句必须转化为一系列的低级机器语言指令，将这些指令打包到一起，并以二进制磁盘文件的形式存储起来，也就是可执行目标文件。

从源文件到可执行目标文件的转化过程：

![](https://golang.design/go-questions/compile/assets/7.png#id=sXo2v&originHeight=396&originWidth=1920&originalType=binary&ratio=1&rotation=0&showTitle=false&status=done&style=none&title=)

完成以上各个阶段的就是 Go 编译系统。你肯定知道大名鼎鼎的 GCC（GNU Compile Collection），中文名为 GNU 编译器套装，它支持像 C，C++，Java，Python，Objective-C，Ada，Fortran，Pascal，能够为很多不同的机器生成机器码。

可执行目标文件可以直接在机器上执行。一般而言，先执行一些初始化的工作；找到 main 函数的入口，执行用户写的代码；执行完成后，main 函数退出；再执行一些收尾的工作，整个过程完毕。

在接下来的文章里，我们将探索`编译`和`运行`的过程。

Go 源码里的编译器源码位于 `src/cmd/compile` 路径下，链接器源码位于 `src/cmd/link` 路径下。

### [Go 编译相关的命令详解](http://golang.design/go-questions/compile/cmd/)
和编译相关的命令主要是：

```go
go build
go install
go run
```

#### go build
`go build` 用来编译指定 packages 里的源码文件以及它们的依赖包，编译的时候会到 `$GoPath/src/package` 路径下寻找源码文件。`go build` 还可以直接编译指定的源码文件，并且可以同时指定多个。

通过执行 `go help build` 命令得到 `go build` 的使用方法：

```go
usage: go build [-o output] [-i] [build flags] [packages]
```

`-o` 只能在编译单个包的时候出现，它指定输出的可执行文件的名字。

`-i` 会安装编译目标所依赖的包，安装是指生成与代码包相对应的 `.a` 文件，即静态库文件（后面要参与链接），并且放置到当前工作区的 pkg 目录下，且库文件的目录层级和源码层级一致。

至于 build flags 参数，`build, clean, get, install, list, run, test` 这些命令会共用一套：

| 参数 | 作用 |
| --- | --- |
| -a | 强制重新编译所有涉及到的包，包括标准库中的代码包，这会重写 /usr/local/go 目录下的 `.a` |
| 文件 |  |
| -n | 打印命令执行过程，不真正执行 |
| -p n | 指定编译过程中命令执行的并行数，n 默认为 CPU 核数 |
| -race | 检测并报告程序中的数据竞争问题 |
| -v | 打印命令执行过程中所涉及到的代码包名称 |
| -x | 打印命令执行过程中所涉及到的命令，并执行 |
| -work | 打印编译过程中的临时文件夹。通常情况下，编译完成后会被删除 |


我们知道，Go 语言的源码文件分为三类：命令源码、库源码、测试源码。

> 命令源码文件：是 Go 程序的入口，包含 `func main()` 函数，且第一行用 `package main` 声明属于 main 包。
>

> 库源码文件：主要是各种函数、接口等，例如工具类的函数。
>

> 测试源码文件：以 `_test.go` 为后缀的文件，用于测试程序的功能和性能。
>

注意，`go build` 会忽略 `*_test.go` 文件。

#### go install
`go install` 用于编译并安装指定的代码包及它们的依赖包。相比 `go build`，它只是多了一个“安装编译后的结果文件到指定目录”的步骤。

还是使用之前 hello-world 项目的例子，我们先将 pkg 目录删掉，在项目根目录执行：

```go
go install src/main.go

或者

go install util
```

两者都会在根目录下新建一个 `pkg` 目录，并且生成一个 `util.a` 文件。

并且，在执行前者的时候，会在 GOBIN 目录下生成名为 main 的可执行文件。

所以，运行 `go install` 命令，库源码包对应的 `.a` 文件会被放置到 `pkg` 目录下，命令源码包生成的可执行文件会被放到 GOBIN 目录。

`go install` 在 GoPath 有多个目录的时候，会产生一些问题，具体可以去看郝林老师的 `Go 命令教程`，这里不展开了。

#### go run
`go run` 用于编译并运行命令源码文件。

### [Go 程序启动过程是怎样的](http://golang.design/go-questions/compile/booting/)
## 框架
### Gin
文档：[https://gin-gonic.com/zh-cn/docs/introduction/](https://gin-gonic.com/zh-cn/docs/introduction/)

<font style="color:rgb(85, 85, 85);">Gin框架介绍及使用-李文周：</font>[https://www.liwenzhou.com/posts/Go/Gin_framework/#autoid-0-0-0](https://www.liwenzhou.com/posts/Go/Gin_framework/#autoid-0-0-0)

Gin源码阅读与分析：[https://www.yuque.com/iveryimportantpig/huchao/zd24cb3z2bco5304](https://www.yuque.com/iveryimportantpig/huchao/zd24cb3z2bco5304)

#### 理解
1. Gin 是一个用 Go 语言编写的<font style="color:#DF2A3F;">轻量级 Web 框架，</font>专注于高效的 HTTP 路由和中间件管理。它以简洁易用的 API 和极高的性能著称，<font style="color:#DF2A3F;">适合开发 RESTful API 和 Web 服务</font>。
2. Gin 的**<font style="color:#DF2A3F;">核心是路由机制</font>**，通过将<font style="color:#DF2A3F;"> HTTP 请求路由到相应的处理函数来实现。它支持路由分组，便于组织和管理复杂的路由结构</font>。
3. 同时，<font style="color:#DF2A3F;">Gin 提供了一套强大的中间件机制，允许在请求到达处理函数之前进行预处理，如日志记录、认证、错误处理等</font>。
4. Gin 的**<font style="color:#DF2A3F;">另一个亮点是它的 JSON 解析和响应处理能力</font>**，通过内置的 `c.JSON` 方法，可以轻松地将数据结构序列化为 JSON 格式返回给客户端。

总的来说，Gin 适合用于开发性能要求高的 Web 应用，尤其是对于需要处理大量并发请求的场景。

#### 特性
1. **快速**
    1. 基于 Radix 树的路由，小内存占用。没有反射。可预测的 API 性能。
2. **支持中间件**
    1. 传入的 HTTP 请求可以由一系列中间件和最终操作来处理。 例如：Logger，Authorization，GZIP，最终操作 DB。
3. **Crash 处理**
    1. Gin 可以 catch 一个发生在 HTTP 请求中的 panic 并 recover 它。这样，你的服务器将始终可用。例如，你可以向 Sentry 报告这个 panic！
4. **JSON 验证**
    1. Gin 可以解析并验证请求的 JSON，例如检查所需值的存在。
5. **路由组**
    1. <font style="color:#DF2A3F;">更好地组织路由。是否需要授权，不同的 API 版本</font>…… 此外，这些<font style="color:#DF2A3F;">组可以无限制地嵌套而不会降低性能</font>。
    2.  Gin 使用<font style="color:#DF2A3F;">基于树状结构的路由匹配算法，能够快速地匹配 URL 路径 </font> 
6. **错误管理**
    1. Gin 提供了一种方便的方法来收集 HTTP 请求期间发生的所有错误。最终，中间件可以将它们写入日志文件，数据库并通过网络发送。
7. **内置渲染**
    1. Gin 为 JSON，XML 和 HTML 渲染提供了易于使用的 API。
8. **可扩展性**
    1. 新建一个中间件非常简单，去查看[示例代码](https://gin-gonic.com/zh-cn/docs/examples/using-middleware/)吧。

#### Gin路由机制  
**Gin 的路由机制非常灵活和高效，主要有以下几个方面：**

1. **路由定义**：  
Gin 使用 `*gin.Engine` 对象来定义路由。可以通过 `GET`、`POST` 等方法为不同的 HTTP 请求定义处理函数。例如：

```go
r := gin.Default()
r.GET("/ping", func(c *gin.Context) {
    c.String(http.StatusOK, "pong")
})
r.Run()
```

2. **路由组**：  
Gin 支持通过 `Group` 方法创建路由组，方便管理相关的路由。例如：

```go
v1 := r.Group("/v1")
{
    v1.GET("/users", getUsers)
    v1.GET("/posts", getPosts)
}
```

3. **路由参数**：  
<font style="color:#262626;">可以在路由中使用动态参数</font>，Gin 会自动提取这些参数。例如：

```go
r.GET("/user/:id", func(c *gin.Context) {
    id := c.Param("id")
    c.String(http.StatusOK, "User ID: %s", id)
})
```

4. **路由方法**：  
<font style="color:#DF2A3F;">Gin 支持 HTTP 的各种请求方法，包括 </font>`<font style="color:#DF2A3F;">GET</font>`<font style="color:#DF2A3F;">、</font>`<font style="color:#DF2A3F;">POST</font>`<font style="color:#DF2A3F;">、</font>`<font style="color:#DF2A3F;">PUT</font>`<font style="color:#DF2A3F;">、</font>`<font style="color:#DF2A3F;">DELETE</font>`<font style="color:#DF2A3F;"> 等</font>，通过对应的方法定义不同的路由处理函数。
5. **路由优先级**：  
更具体的路由定义优先匹配，例如带有路径参数的路由会比通用的路由更优先匹配。
6. **中间件**：  
<font style="color:#DF2A3F;">Gin 允许为路由定义中间件，用于处理请求的预处理、认证、日志记录等任务</font>。例如：

```go
r.Use(gin.Logger())
r.Use(gin.Recovery())
```

**总结**：Gin 的路由机制通过提供清晰的路由定义、灵活的路由分组、动态参数支持、方法匹配和中间件支持，使得构建高效、结构化的 Web 应用变得简单和高效。

####  请求打入到响应的一个过程  
**Gin 框架的请求处理过程大致分为以下几个步骤：**

1. **<font style="color:#DF2A3F;">请求接收</font>**<font style="color:#DF2A3F;">：</font>  
当 HTTP 请求到达 Gin 应用时，<font style="color:#DF2A3F;">Gin 框架会首先接收到请求。这些请求会被 </font>`<font style="color:#DF2A3F;">*gin.Engine</font>`<font style="color:#DF2A3F;"> 对象处理</font>，`Engine` 是 Gin 的核心组件。
2. **路由匹配**：  
<font style="color:#DF2A3F;">Gin 根据请求的 URL 和 HTTP 方法（如 GET、POST）来匹配路由</font>。框架会查找定义的路由规则，并<font style="color:#DF2A3F;">找到与请求最匹配的处理函数（Handler）</font>。
3. **中间件处理**：  
<font style="color:#DF2A3F;">在执行路由处理函数之前，Gin 会依次执行与该路由关联的中间件</font>。中间件可以用于请求的预处理，如认证、日志记录等。
4. **执行处理函数**：  
<font style="color:#DF2A3F;">中间件执行完毕后，Gin 会调用匹配的路由处理函数</font>。处理函数可以访问请求数据、处理业务逻辑，并准备响应数据。
5. **生成响应**：  
处理函数会通过 `<font style="color:#DF2A3F;">*gin.Context</font>`<font style="color:#DF2A3F;"> 对象生成响应。可以设置响应状态码、响应头以及响应体</font>。Gin 提供了多种方法来构造响应，比如 `c.String()`、`<font style="color:#DF2A3F;">c.JSON()</font>`、`c.XML()` 等。
6. **响应返回**：  
最终，<font style="color:#DF2A3F;">Gin 将响应数据发送回客户端，完成请求-响应周期</font>。

**总结**：Gin 框架处理请求的流程从接收请求开始，经过<font style="color:#DF2A3F;">路由匹配和中间件处理，执行处理函数，生成并返回响应</font>。整个过程高效且结构清晰，帮助开发者快速构建 Web 应用。

#### [gin目录结构](https://blog.csdn.net/qq_34877350/article/details/107959381)
文档：[https://blog.csdn.net/qq_34877350/article/details/107959381](https://blog.csdn.net/qq_34877350/article/details/107959381)

```plain
├── gin
│   ├──  Router
│          └── router.go
│   ├──  Middlewares
│          └── corsMiddleware.go
│   ├──  Controllers
│          └── testController.go
│   ├──  Services
│          └── testService.go
│   ├──  Models
│          └── testModel.go
│   ├──  Databases
│          └── mysql.go
│   ├──  Sessions
│          └── session.go
└── main.go

```

+ 使用gorm访问数据库
+ gin 为项目根目录
+ main.go 为入口文件
+ Router 为路由目录
+ Middlewares 为中间件目录
+ Controllers 为控制器目录（MVC）
+ Services 为服务层目录，这里把DAO逻辑也写入其中，如果分开也可以
+ Models 为模型目录
+ Databases 为数据库初始化目录
+ Sessions 为session初始化目录
+ 文件 引用顺序 大致如下：
+ main.go(在main中关闭数据库) - router(Middlewares) - Controllers - Services(sessions) - Models - Databases

### go-zero
文档：[https://go-zero.dev/cn/docs/introduction](https://go-zero.dev/cn/docs/introduction)

> go-zero 是一个集成了各种工程实践的 web 和 rpc 框架。通过弹性设计保障了大并发服务端的稳定性，经受了充分的实战检验。  
go-zero 包含极简的 API 定义和生成工具 goctl，可以根据定义的 api 文件一键生成 Go, iOS, Android, Kotlin, Dart, TypeScript, JavaScript 代码，并可直接运行。
>

使用 go-zero 的好处：

+ 轻松获得支撑千万日活服务的稳定性
+ 内建级联超时控制、限流、自适应熔断、自适应降载等微服务治理能力，无需配置和额外代码
+ 微服务治理中间件可无缝集成到其它现有框架使用
+ 极简的 API 描述，一键生成各端代码
+ 自动校验客户端请求参数合法性
+ 大量微服务治理和并发工具包

## 场景
### 有没有遇到过cpu不高但是内存高的场景？怎么排查的
 在实际开发中，遇到 CPU 使用率不高但内存占用很高的情况并不少见。这种现象通常<font style="color:#DF2A3F;">表明程序中存在内存泄漏、内存占用过大、或者内存管理不当</font>的问题。下面是一个排查的步骤：  

在实际开发中，遇到 CPU 使用率不高但内存占用很高的情况并不少见。这种现象通常表明程序中存在内存泄漏、内存占用过大、或者内存管理不当的问题。下面是一个排查的步骤：

#### 检查内存占用情况
+ **工具：**`top`**, **`htop`**, **`ps`  
使用这些系统工具查看内存占用较高的进程，确认是否是你的 Go 程序导致的内存消耗。
+ `pmap`  
使用 `pmap <PID>` 查看进程的内存分布，确定是哪个内存段占用最大（如 heap、stack）。

#### 分析 Go 程序的内存使用
+ **内存分配情况：**`pprof`  
使用<font style="color:#DF2A3F;"> Go 的 </font>`<font style="color:#DF2A3F;">pprof</font>`<font style="color:#DF2A3F;"> 工具生成内存快照</font>（heap profile）:

```bash
go tool pprof http://localhost:6060/debug/pprof/heap
```

分析结果可以<font style="color:#DF2A3F;">帮助你识别哪些对象在堆上占用最多的内存</font>。

+ **查看 Goroutine 状态**  
使用 `<font style="color:#DF2A3F;">pprof</font>`<font style="color:#DF2A3F;"> 中的 Goroutine 分析工具</font>：

```bash
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

查看<font style="color:#DF2A3F;">是否存在大量 Goroutine 导致内存占用</font>。

#### 检查内存泄漏
+ **是否有未释放的内存**  
检查代码中是否存在未释放的资源，如未关闭的文件描述符、数据库连接、未清理的缓存等。
+ **工具：**`leaktest`**, **`goleak`  
使用 `leaktest` 或 `goleak` 工具检测 Goroutine 泄漏，这些泄漏可能会导致内存无法被回收。

#### 优化内存使用
+ **<font style="color:#DF2A3F;">减少对象分配</font>**<font style="color:#DF2A3F;">  
</font><font style="color:#DF2A3F;">尽量复用内存，如使用 </font>`<font style="color:#DF2A3F;">sync.Pool</font>`<font style="color:#DF2A3F;"> 来管理重复使用的对象，避免频繁的内存分配和 GC 压力。</font>
+ **优化数据结构**  
检查是否使用了不必要的大型数据结构（如 map, slice），考虑更合适的替代方案。

通过以上步骤，通常可以较为全面地排查和解决 CPU 不高但内存高的问题。



### 怎么实时查看k8s内存占用的
要实时查看 Kubernetes 集群中 Pod 的内存占用情况，有几种常见的方法：

#### 使用 `kubectl top` 命令
`**<font style="color:#DF2A3F;">kubectl top</font>**` 是 Kubernetes 提供的一个工具，可以<font style="color:#DF2A3F;">实时查看 Pod 和节点的资源使用情况（包括 CPU 和内存）</font>。

```bash
# 查看某个命名空间下所有 Pod 的资源使用情况
kubectl top pod -n <namespace>

# 查看指定 Pod 的资源使用情况
kubectl top pod <pod-name> -n <namespace>

# 查看集群中所有节点的资源使用情况
kubectl top nodes
```

这个命令会显示每个 Pod 当前的 CPU 和内存使用量，以及节点的总资源消耗。

#### 使用 `kubectl describe pod` 命令
`kubectl describe` 命令可以查看单个 Pod 的详细信息，包括资源请求和限制：

```bash
kubectl describe pod <pod-name> -n <namespace>
```

这不会直接显示实时内存使用情况，但你可以看到 Pod 所请求和限制的内存资源。

#### 使用 Kubernetes Dashboard
<font style="color:#DF2A3F;">Kubernetes Dashboard 是一个 web 界面的 UI，可以查看集群中各种资源的使用情况，包括实时的内存消耗</font>。

+ 安装 Kubernetes Dashboard：

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0/aio/deploy/recommended.yaml
```

+ 访问 Dashboard：

```bash
kubectl proxy
```

然后在浏览器中打开 `http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/`。

在 Dashboard 中，你可以查看各个 Pod 的详细资源使用情况，包括实时的内存和 CPU 使用。

#### 使用 <font style="color:#DF2A3F;">Prometheus + Grafana</font> 监控
如果你的集群已经配置了 Prometheus 和 Grafana，你可以使用它们来实时监控内存使用情况：

+ **<font style="color:#DF2A3F;">Prometheus</font>**<font style="color:#DF2A3F;">：收集和存储 Kubernetes 中的指标数据。</font>
+ **<font style="color:#DF2A3F;">Grafana</font>**<font style="color:#DF2A3F;">：提供丰富的仪表盘，可以实时显示集群中各个资源的使用情况</font>。

在 Grafana 中，你可以创建或使用现有的仪表盘来监控 Pod 和节点的内存使用情况。

#### 直接查看容器内的内存使用
如果你想直接查看某个容器的内存使用情况，可以<font style="color:#DF2A3F;">进入容器内部，然后使用 </font>`<font style="color:#DF2A3F;">top</font>`<font style="color:#DF2A3F;"> 或 </font>`<font style="color:#DF2A3F;">free</font>`<font style="color:#DF2A3F;"> 等命令</font>：

```bash
kubectl exec -it <pod-name> -n <namespace> -- bash

# 在容器内使用 top 或 free 命令
top
free -m
```

### 6. **使用 **`kubectl get --raw`** 命令**
你可以直接通过 Kubernetes API 获取内存使用情况，返回结果为 JSON 格式：

```bash
kubectl get --raw /apis/metrics.k8s.io/v1beta1/namespaces/<namespace>/pods/<pod-name>
```

这个方法适合进行脚本化或编程访问资源使用数据。

通过以上这些方法，你可以实时查看 Kubernetes 中的内存使用情况，并及时了解资源的分配与消耗。

## 参考并致谢
1、可可酱 [可可酱：Golang常见面试题](https://zhuanlan.zhihu.com/p/434629143)  
2、Bel_Ami同学 [golang 面试题(从基础到高级)](https://link.zhihu.com/?target=https%3A//blog.csdn.net/Bel_Ami_n/article/details/123352478)

