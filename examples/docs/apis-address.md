# 地址管理

收货地址管理接口

## apis

### 添加地址接口

author: _alovn_

```text
POST /user/address/create
```

__Request__:

parameter|parameterType|dataType|required|validate|example|description
--|--|--|--|--|--|--
__address__|_form_|string|true|required|"abc"|地址
__city_id__|_form_|int64|true|required|123|城市ID

__Response__:

```json
// StatusCode: 200 
{  //object(common.Response), 通用返回结果
  "code": 0,  //int, 返回状态码
  "msg": "success",  //string, 返回消息
}
```

---

### 更新地址接口

author: _alovn_

```text
POST /user/address/update
```

__Request__:

parameter|parameterType|dataType|required|validate|example|description
--|--|--|--|--|--|--
__address__|_form_|string|true|required|"abc"|地址
__id__|_form_|int64|true|required|123|地址ID

__Response__:

```json
// StatusCode: 200 
{  //object(common.Response), 通用返回结果
  "code": 0,  //int, 返回状态码
  "msg": "success",  //string, 返回消息
}
```

---

### 删除地址接口

author: _alovn_

```text
POST /user/address/delete
```

__Request__:

parameter|parameterType|dataType|required|validate|example|description
--|--|--|--|--|--|--
__id__|_form_|int64|true|required|123|地址ID

__Response__:

```json
// StatusCode: 200 
{  //object(common.Response), 通用返回结果
  "code": 0,  //int, 返回状态码
  "msg": "success",  //string, 返回消息
}
```

---

### 获取地址信息

author: _alovn_

```text
GET /user/address/get/:id
```

__Request__:

parameter|parameterType|dataType|required|validate|example|description
--|--|--|--|--|--|--
__id__|_param_|int64|false||123|地址ID

__Response__:

```json
// StatusCode: 200 
{  //object(common.Response), 通用返回结果
  "code": 0,  //int, 返回状态码
  "data": {  //object(handler.AddressResponse), 返回地址信息
    "address": "abc",  //string, 地址信息
    "city_id": 123,  //int64, 城市ID
    "id": 123  //int64, 地址ID
  },
  "msg": "success"  //string, 返回消息
}
```

---

### 获取地址列表

___Deprecated___

获取收货地址列表

author: _alovn_

```text
GET /user/address/list
```

__Response__:

```json
// StatusCode: 200 
{  //object(common.Response), 通用返回结果
  "code": 0,  //int, 返回状态码
  "data": [  //array[handler.AddressResponse]
    {  //object(handler.AddressResponse), 返回地址信息
      "address": "abc",  //string, 地址信息
      "city_id": 123,  //int64, 城市ID
      "id": 123  //int64, 地址ID
    }
  ],
  "msg": "success"  //string, 返回消息
}
```

---