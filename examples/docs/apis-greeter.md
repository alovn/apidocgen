# greeter分组

greeter分组说明

## Apis

### @api 测试greeter

```text
GET /greeter
```

**Request**:
parameters|type|required|validate|example|description
--|--|--|--|--|--
**id**|*header,query*|true||12357|this id
**tid**|*param*|true|required|123|
**token**|*header*|false||"example"|
**Response**:

```json
// StatusCode: 200

// 输出对象 dd
{  //object(common.Response)
  "Code": 123,  //int
  "Common": 1.23,  //float64
  "Common2": "example",  //string
  "Data": null,  //any
  "Msg": "example",  //string
  "data": {  //object(main.TestData)
    "MyIntArray": [  //array[int]
      123
    ],
    "MyTestData2Array": [  //array[main.TestData2]
      {  //object(main.TestData2)
        "MyAge2": 123,  //int
        "MyTitle2": "example"  //string, 标题2
      }
    ],
    "data2": {  //object(main.TestData2)
      "MyAge2": 123,  //int
      "MyTitle2": "example"  //string, 标题2
    },
    "my_title": "example"  //string, 标题
  }
}
```

```json
// StatusCode: 500

// 出错了
{  //object(main.Response), 通用返回结果
  "code": 10010,  //int, 返回状态码
  "msg": "异常",  //string, 返回文本消息
}
```

### @api 测试greeter2

```text
GET /greeter2
```

**Response**:

```json
// StatusCode: 200

// 输出对象 dd
{  //object(main.TestData)
  "MyIntArray": [  //array[int]
    123
  ],
  "MyTestData2Array": [  //array[main.TestData2]
    {  //object(main.TestData2)
      "MyAge2": 123,  //int
      "MyTitle2": "example"  //string, 标题2
    }
  ],
  "data2": {  //object(main.TestData2)
    "MyAge2": 123,  //int
    "MyTitle2": "example"  //string, 标题2
  },
  "my_title": "example"  //string, 标题
}
```
