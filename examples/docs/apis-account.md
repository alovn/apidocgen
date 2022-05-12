# 账户相关

账户相关的接口，含用户注册、登录等

1. [用户注册接口](#1-用户注册接口)
2. [用户登录接口](#2-用户登录接口)

## apis

### 1. 用户注册接口

用户注册接口说明

author: _alovn_

```text
POST /user/account/register
```

__Request__:

parameter|parameterType|dataType|required|validate|example|description
--|--|--|--|--|--|--
__from__|_query_|string|false|||test
__password__|_form_|string|true|required||密码
__username__|_form_|string|true|required||用户名
__x-request-id__|_header_|string|false|||request id

_body_:

```javascript
{  //object(handler.RegisterRequest), 注册请求参数
  "password": "abc",  //string, required, 密码
  "username": "abc"  //string, required, 用户名
}
```

__Response__:

```javascript
//StatusCode: 200 注册成功返回数据
{  //object(common.Response), 通用返回结果
  "code": 0,  //int, 返回状态码
  "data": {  //object(handler.RegisterResponse), 注册返回数据
    "user_id": 123,  //int64, 注册的用户ID
    "username": "abc",  //string, 注册的用户名
    "welcome_msg": "abc"  //string, 注册后的欢迎语
  },
  "msg": "success"  //string, 返回消息
}
```

```javascript
//StatusCode: 200 密码格式错误
{  //object(common.Response), 通用返回结果
  "code": 10011,  //int, 返回状态码
  "msg": "password format error",  //string, 返回消息
}
```

---

### 2. 用户登录接口

author: _alovn_

```text
POST /user/account/login
```

__Request__:

parameter|parameterType|dataType|required|validate|example|description
--|--|--|--|--|--|--
__password__|_form_|string|true|required||登录密码
__username__|_form_|string|true|required||登录用户名
__validate_code__|_form_|string|false|||验证码

__Response__:

```javascript
//StatusCode: 200 登录成功返回数据
{  //object(common.Response), 通用返回结果
  "code": 0,  //int, 返回状态码
  "data": {  //object(handler.LoginResponse), 登录返回数据
    "welcome_msg": "abc"  //string, 登录成功欢迎语
  },
  "msg": "success"  //string, 返回消息
}
```

```javascript
//StatusCode: 200 密码错误
{  //object(common.Response), 通用返回结果
  "code": 10020,  //int, 返回状态码
  "msg": "password_error",  //string, 返回消息
}
```

---
