# greeter分组

greeter分组说明

## Apis

### 测试greeter

```text
GET /greeter
```


Response:
```json

// 输出对象 dd
// HTTP StatusCode: 200
// object main 
{
  "code": 0,  // int, 返回状态码
  "msg": "返回消息",  // string, 返回文本消息
  "data": nil  // any, 返回的具体数据
}
```


### 测试greeter2

```text
GET /greeter2
```


Response:
```json

// 输出对象 dd
// HTTP StatusCode: 200
// object main 
{
  "my_title": ""
}
```

