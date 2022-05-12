# 测试示例

其它一些示例演示

1. [struct数组](#1-struct数组)
2. [int数组](#2-int数组)
3. [int](#3-int)
4. [DemoMap](#4-DemoMap)

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

### 2. int数组

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

### 3. int

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

### 4. DemoMap

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
