# 菜单管理

菜单管理接口

1 [获取菜单节点](#1-获取菜单节点)

## apis

### 1. 获取菜单节点

测试数组、递归结构体

author: _alovn_

```text
GET /user/menu/nodes
```

__Response__:

```json
// StatusCode: 200 
{  //object(common.Response), 通用返回结果
  "code": 0,  //int, 返回状态码
  "data": [  //array[handler.Node]
    {  //object(handler.Node)
      "id": 123,  //int64
      "name": "abc",  //string
      "nodes": [  //array[handler.Node]

      ]
    }
  ],
  "msg": "success"  //string, 返回消息
}
```

---
