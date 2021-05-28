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
* 引入了加密token的生成（代码为照搬）

### 遇到的问题
* 生成token的时候，一直进行token不匹配报错。最终发现为unix时间生成的int64直接转换成string不对
* int64转string需要使用strconv.FormatInt()解决


