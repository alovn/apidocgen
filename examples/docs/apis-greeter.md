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
    "Map2": {  //object(main.TestData2)
      "abc": null  //object
    },
    "Map3": {  //object(main.TestData2)
      "abc": {  //object(main.TestData2)
        "MyAge2": 123,  //int
        "MyTitle2": "abc"  //string, 标题2
      }
    },
    "MyInts": [  //array[int]
      123
    ],
    "Nodes": {  //object(main.Node)
      "abc": {  //object(main.Node)
        "Name": "abc",  //string
        "Nodes": {  //object(main.Node)
          "abc": null  //object
        }
      }
    },
    "amap": ,  //object
    "data2": {  //object(main.TestData2)
      "MyAge2": 123,  //int
      "MyTitle2": "abc"  //string, 标题2
    }
  },
  "msg": "返回消息"  //string, 返回文本消息
}
```

### 测试greeter2

```text
GET /greeter2
```

**Response**:
