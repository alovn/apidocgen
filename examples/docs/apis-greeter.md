# greeter分组

greeter分组说明

## Apis

### 测试greeter

```text
GET /greeter
```
**Response**:

```json

// 输出对象 dd
// StatusCode: 200
// object main.Response 
{
  "code": 0,  // int, 返回状态码
  "data": {
    "data2": {
      "MyTitle2": "example"
    },
    "my_title": "example"  // string, 标题
  },
  "msg": "返回消息"  // string, 返回文本消息
}
```

```json

// 出错了
// StatusCode: 500
// object main.Response 
{
  "code": 10010,  // int, 返回状态码
  "msg": "异常",  // string, 返回文本消息
}
```

```json

// 错误
// StatusCode: 500
// integer integer 

```


### 测试greeter2

```text
GET /greeter2
```
**Response**:

```json

// 输出对象 dd
// StatusCode: 200
// object main.TestData 
{
  "data2": {
    "MyTitle2": "example"
  },
  "my_title": "example"  // string, 标题
}
```

