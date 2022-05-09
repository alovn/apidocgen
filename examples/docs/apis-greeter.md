# greeter分组

greeter分组说明

## Apis

### 测试greeter

```text
GET /greeter
```

**Response**:

```json
// StatusCode: 200

// 输出对象 dd
{  //object(main.Response), 通用返回结果
  "code": 0,  //int, 返回状态码
  "data": {  //object(main.TestData)
    "MyInt": 123,  //int
    "MyInts": [  //array[int]
      123
    ]
  },
  "msg": "返回消息"  //string, 返回文本消息
}
```

### 测试greeter2

```text
GET /greeter2
```

**Response**:
