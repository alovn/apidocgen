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
    "MyIntArray": [  // array(int), 
      0
    ],
    "MyTestData2Array": [  // array(main.TestData2), 
      {
        "MyAge2": 0,  // int
        "MyTitle2": "example"  // string, 标题2
      }
    ],
    "my_title": "example"  // string, 标题
  },
  "msg": "返回消息"  // string, 返回文本消息
}
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
  "MyIntArray": [  // array(int), 
    0
  ],
  "MyTestData2Array": [  // array(main.TestData2), 
    {
      "MyAge2": 0,  // int
      "MyTitle2": "example"  // string, 标题2
    }
  ],
  "my_title": "example"  // string, 标题
}
```
