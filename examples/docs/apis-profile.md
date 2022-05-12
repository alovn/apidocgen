# 资料管理

用户资料管理接口

1. [获取用户资料](#1-获取用户资料)

## apis

### 1. 获取用户资料

author: _alovn_

```text
GET /user/profile/get
```

__Response__:

```javascript
//StatusCode: 200 
{  //object(common.Response), 通用返回结果
  "code": 0,  //int, 返回状态码
  "data": {  //object(handler.ProfileResponse)
    "extends": {  //object(Extends)
      "abc": "abc"  //string, 扩展信息
    },
    "gender": 1,  //uint8
    "username": "abc"  //string
  },
  "msg": "success"  //string, 返回消息
}
```

---
