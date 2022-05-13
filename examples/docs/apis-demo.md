# 测试示例

其它一些示例演示

1. [struct数组](#1-struct数组)
2. [struct嵌套](#2-struct嵌套)
3. [int数组](#3-int数组)
4. [int](#4-int)
5. [DemoMap](#5-DemoMap)
6. [XML测试](#6-XML测试)

## apis

### 1. struct数组

```text
GET /user/demo/struct_array
```

__Response__:

```javascript
//StatusCode: 200 demo struct array
[  //array[handler.DemoData]
  {  //object(handler.DemoData)
    "count": 123,  //int
    "description": "abc",  //string
    "float_array": [  //array[float64]
      1.23
    ],
    "int_array": [  //array[int]
      123
    ],
    "int_pointer": 123,  //int
    "map": {  //object(Map)
      "abc": 123  //int
    },
    "object_1": {  //object(handler.DemoObject)
      "name": "abc"  //string
    },
    "object_2": {  //object(handler.DemoObject)
      "name": "abc"  //string
    },
    "title": "abc"  //string, 标题
  }
]
```

---

### 2. struct嵌套

```text
GET /user/demo/struct_nested
```

__Response__:

```javascript
//StatusCode: 200 nested struct
{  //object(handler.Struct1)
  "Name": "abc",  //string
  "Name2": "abc"  //string
}
```

---

### 3. int数组

```text
GET /user/demo/int_array
```

__Response__:

```javascript
//StatusCode: 200 demo int array
[  //array[int]
  123
]
```

---

### 4. int

```text
GET /user/demo/int
```

__Response__:

```javascript
//StatusCode: 200 demo int
//int
123
```

---

### 5. DemoMap

```text
GET /user/demo/map
```

__Response__:

```javascript
//StatusCode: 200 demo map
{  //object(handler.DemoData)
  "abc": {  //object(handler.DemoData)
    "count": 123,  //int
    "description": "abc",  //string
    "float_array": [  //array[float64]
      1.23
    ],
    "int_array": [  //array[int]
      123
    ],
    "int_pointer": 123,  //int
    "map": {  //object(Map)
      "abc": 123  //int
    },
    "object_1": {  //object(handler.DemoObject)
      "name": "abc"  //string
    },
    "object_2": {  //object(handler.DemoObject)
      "name": "abc"  //string
    },
    "title": "abc"  //string, 标题
  }
}
```

---

### 6. XML测试

author: _alovn_

```text
GET /user/demo/xml
```

__Request__:

parameter|parameterType|dataType|required|validate|example|description
--|--|--|--|--|--|--
__id__|_param_|int64|false|||DemoID

_body_:

```javascript
<request> //object(handler.DemoXMLRequest), XML测试请求对象
  <id>123</id> //int64, DemoID
</request>
```

__Response__:

```javascript
//StatusCode: 200 
<response> //object(common.Response), 通用返回结果
  <code>0</code> //int, 返回状态码
  <demo> //object(handler.DemoXMLResponse), XML测试返回对象
    <address>abc</address> //string, 地址信息
    <city_id>123</city_id> //int64, 城市ID
    <id>123</id> //int64, 地址ID
  </demo>
  <msg>success</msg> //string, 返回消息
</response>
```

---
