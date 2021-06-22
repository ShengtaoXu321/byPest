# Key

## 说明
这是根据`daqiuyin`接口，进行`POST`方式访问，获取的数据

## 特性
* 涉及到头部部分字段加密解密，此处代码为@huanting撰写。直接搬用过来
* 进行POST方式访问`daqiuyin`这个提供的接口，需要填充头部信息
* 头部信息协议如下
  
  *  `Content-Type`字段为`application/json`
  *  `X-DAQIUYIN-ID`字段为`5f45d17204da596300*****02`
  *  `SIGN`字段为需要解密的信息
  *  `DATE`字段为UNIX时间戳
 
* 需要解密的思路为：将时间与某参数进行sha256运算
  



