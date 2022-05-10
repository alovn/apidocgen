# 账户相关

账户相关的接口

## apis

### 用户注册接口

*author: alovn*

```text
POST /user/account/register
```

**Request**:

parameters|type|required|validate|example|description
--|--|--|--|--|--
**password**|_form_|true|required|"abc"|密码
**username**|_form_|true|required|"abc"|用户名

**Response**:

```json
// StatusCode: 200

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

### 用户登录接口

*author: alovn*

```text
POST /user/account/login
```

**Request**:

parameters|type|required|validate|example|description
--|--|--|--|--|--
**password**|_form_|true|required|"abc"|登录密码
**username**|_form_|true|required|"abc"|登录用户名
**validate_code**|_form_|false||"abc"|验证码

**Response**:

```json
// StatusCode: 200

// 登录成功返回数据
{  //object(common.Response), 通用返回结果
  "code": 0,  //int, 返回状态码
  "data": {  //object(handler.LoginResponse), 登录返回数据
    "welcome_msg": "abc"  //string, 登录成功欢迎语
  },
  "msg": "success"  //string, 返回消息
}
```

```json
// StatusCode: 200

// 密码错误
{  //object(common.Response), 通用返回结果
  "code": 10020,  //int, 返回状态码
  "msg": "password_error",  //string, 返回消息
}
```
