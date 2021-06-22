# 北邮虫害服务器

## 第一版本
### 特性
* 基础的数据透传，未作解析
* 给硬件平台提供了上报接口
* docker部署，端口为18088

### 完成时间
2021年4月1号

## 第二版本
### 特性
* 完成数据解析，设置筛选规则
* 给硬件平台提供了上报接口
* 数据解析后存入数据库
* 给网页前端提供查询历史数据和最新数据的接口

### 完成时间
2021年5月21号

## 第三版
### 特性
* 在跟前端网页进行对接时，出现了`跨域`问题；第三版本已解决
* 请求最新数据时，新增了Get方法

    ####  补充问题：跨域
    1. 什么是跨域？
    
    跨域是由浏览器同源政策引起的，即为页面请求的接口地址，必须与页面url处于同域上（域名、端口、协议相同）。
       这是为了防止某域名下的接口被其他域名下的网页非法调用。
   
   2. 举例说明： 同源、跨域  以http://store.company.com/dir/page.html为例子
    
    |  URL   | 结果  |原因|
    |  :----:  | :----:  |:---:|
    |  http://store.company.com/dir2/other.html | 同源 |只有路径不同|
    | http://store.company.com/dir/inner/another.html  | 同源 |只有路径不同|
    |https://store.company.com/secure.html|失败|http和https协议不同|
    |http://store.company.com:81/dir/etc.html|失败|端口不同，80和81|
    |http://news.company.com/dir/other.html|失败|主机不同|
    
  3. 同源政策
        
    同源政策的目的是为了保证用户的信息安全，防止恶意网站窃取数据。
    非同源，会有三种行为受到限制
    
     （1）Cookie、LocalStorage和IndexDB无法读取
      
     （2）DOM无法获得

     （3）AJAX请求无法发送

  ####  如何解决跨域
  * 跨域资源共享(Cross-origin resource sharing)---CORS
    
  CROS允许浏览器向跨源服务器，发出`XMLHttpRequest`请求（一个浏览器接口，从而是的js可以进行HTTP(S)通信），从而克服了AJAX只能同源使用的限制

  对于后端开发来说，浏览器一旦发现AJAX请求跨源，就会自动添加一些附加的头部信息，要实现CORS通信的关键就是服务器实现CORS接口。
  
  * CORS两种请求
    
   **简单请求(simple request)**和**非简单请求(not-so-simple request)**
  
  

### 完成时间
2021年5月25日


## 第四版
### 特性
* 新增霉变数据的获取功能
* 引入了加密`token`的生成（代码为照搬）

### 遇到的问题
* 生成`token`的时候，一直进行`token`不匹配报错。最终发现为unix时间生成的int64直接转换成string不对
* int64转string需要使用`strconv.FormatInt()`解决

### 完成时间
2021年5月28日

## 第五版
### 特性
* 对获取到的霉变数据，使用`gjson`进行解析
* 将解析后的数据，进行插入Germs集合的操作
* 使用`time.sleep`函数作为定时，设置24小时获取一次霉变数据

#### 遇到的问题
* 数据格式繁多，在进行数据库操作时候，需要定义好结构体
* 如何使用gjosn取出类似json数据和不是json的数据

### 完成时间
2021年5月29日

## 第六版
### 特性
* 完成前端请求霉变数据的请求接口
* 修正文件格式，使其规范化

#### 遇到的问题
* windows本机运行无问题，部署时在Linux上获取不到POST请求
* 改用go的交叉编译时，编译不成功
* 交叉编译成功后，利用dockerfile部署的时候，报错x506：证书错误

### 解决方法
* Docker上面POST请求报错，是因为x506:证书不通过。解决方法：忽略http请求证书
* 解决参考：https://stackoverflow.com/questions/12122159/how-to-do-a-https-request-with-bad-certificate
* 交叉编译：windows上面一步一步执行

### 完成时间
2021年6月3日

## 第七版
### 特性
* 霉变数据新增一个level字段
* 更改docker部署方式，交叉编译后，直接linux运行

### 部署时注意
```shell
# 先查看端口占用情况
lsof -i:18088
# 杀死进程
kill pid
# 赋予权限
chmod xxx pest
# 运行文件,并将结果格式化输出到
nohup ./pest > out.file 2>&1 &

```
### 解决docker部署访问不到霉变数据的问题
* 在dockerfile中加入了时间配置，安装tzdata软件包，并生成软连接
```dockerfile
RUN apt-get update
RUN apt-get install -y tzdata && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
```
* 解决代码中对ymd非标准时间转UNIX时间戳出现的问题
* 解决docker中POST请求，出现证书不通过的问题

### 存在的问题
* 代码臃肿，复用程度不够
* 在进行数据库的增删改查时，处理逻辑需要改进
* 在解决docker中时区上海定位不到的问题，除了上述的解决思路，是否还有其他的解决方案

### 完成时间
2021年6月17日


## 第八版
### 特性
* 注释掉获取硬件数据的API，改为模拟数据生成
* 新增一个perm1接口，该接口功能：前端发送时间段内的所有时间戳、`iddev`种类数据的上报

### 完成时间
2021年6月22日


