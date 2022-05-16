# 测试示例

其它一些示例演示

1. [struct数组](#1-struct数组)
2. [struct嵌套](#2-struct嵌套)
3. [int数组](#3-int数组)
4. [int](#4-int)
5. [map](#5-map)
6. [xml](#6-xml)
7. [time](#7-time)

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
    "map": {  //object(map[string]int)
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

### 5. map

```text
GET /user/demo/map
```

__Response__:

```javascript
//StatusCode: 200 demo map
{  //object(map[string]handler.DemoData)
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
    "map": {  //object(map[string]int)
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

### 6. xml

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

```javascript
//StatusCode: 200 
<response> //object(common.Response), 通用返回结果
  <code>0</code> //int, 返回状态码
  <data> //object(handler.DemoXMLResponse2), XML测试返回对象2
    <address>abc</address> //string, 地址信息
    <city_id>123</city_id> //int64, 城市ID
    <id>123</id> //int64, 地址ID
  </data>
  <msg>success</msg> //string, 返回消息
</response>
```

```javascript
//StatusCode: 200 
<response> //object(common.Response), 通用返回结果
  <code>10010</code> //int, 返回状态码
  <msg>sme error</msg> //string, 返回消息
</response>
```

---

### 7. time

author: _alovn_

```text
GET /user/demo/time
```

__Response__:

```javascript
//StatusCode: 200 
{  //object(common.Response), 通用返回结果
  "code": 0,  //int, 返回状态码
  "data": {  //object(handler.DemoTime)
    "time_1": "2022-05-16T02:47:13.923883+08:00",  //object(time.Time), example1
    "time_2": "2022-05-14 15:04:05",  //object(time.Time), example2
    "time_3": "2022-05-16T02:47:13.924032+08:00"  //object(time.Time)
  },
  "msg": "success"  //string, 返回消息
}
```

---
