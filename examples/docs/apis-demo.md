# 测试示例

其它一些示例演示

1. [struct数组](#1-struct数组)
2. [struct嵌套](#2-struct嵌套)
3. [int数组](#3-int数组)
4. [int](#4-int)
5. [DemoMap](#5-DemoMap)

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
    "DemoObject": {  //object(handler.DemoObject)
      "name": "abc"  //string
    },
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
    "DemoObject": {  //object(handler.DemoObject)
      "name": "abc"  //string
    },
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
    "title": "abc"  //string, 标题
  }
}
```

---
