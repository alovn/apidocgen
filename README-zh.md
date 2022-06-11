# apidocgen

[English](./README.md) | 简体中文

apidocgen 是一个为Go项目生成api文档(makrdown格式)和mock的工具。

## 安装

```bash
go install github.com/alovn/apidocgen@latest
```

## 命令

```bash
$ apidocgen --help
apidocgen is a tool for Go to generate apis markdown docs and mocks. For example:

apidocgen --dir= --excludes= --output= --output-index= --template= --single --gen-mocks
apidocgen --mock --mock-listen=:8001
apidocgen mock --data= --listen=

Usage:
  apidocgen [flags]
  apidocgen [command]

Available Commands:
  help        Help about any command
  mock        
  version     print version

Flags:
      --dir string            Search apis dir, comma separated (default ".")
      --excludes string       Exclude directories and files when searching, comma separated
      --gen-mocks             Is generate the mock files
  -h, --help                  help for apidocgen
      --mock-listen string    Mock Server listen address (default "localhost:8001")
      --mock                  Serve a mock server
      --output string         Generate markdown files dir (default "./docs/")
      --output-index string   Generate index file name
      --single                Is generate a single doc
      --template string       Template name or custom template directory, built-in includes markdown and apidocs

Use "apidocgen [command] --help" for more information about a command.
```

## 模板

模板用于生成文档，内置的模板名称有 `markdown` 和 `apidocs`, 默认是 `markdown`.

其中 `apidocs` 模板是用于 [apidocs](https://github.com/alovn/apidocs) 这个项目中显示的。

也可以使用自己的模板文件:

```bash
apidocgen \
    --dir=svc-user,common \
    --template=/Users/xxx/workspace/apidocs/custom-template-direcoty \
    --output=./docs
```

## 怎样使用

apidocgen 支持Go语言下所有的web框架. 下面是一个 [bytego](https://github.com/gostack-labs/bytego) 框架的示例.

1. 在`main.go`文件中添加注解(注释):

    ```go
    //@title UserService
    //@service svc-user
    //@version 1.0.1
    //@desc the api about users
    //@baseurl /user
    func main() {
        r := bytego.New()
        c := controller.NewController()
        //@group account
        //@title Account
        //@desc account register and login
        account := r.Group("/account")
        {
            account.POST("/register", c.Register)
            account.POST("/login", c.Login)
        }
        _ = r.Run(":8000")
    }
    ```

2. 为API添加注解.

    ```go
    //@title AccountLogin
    //@api POST /account/login
    //@group account
    //@accept json
    //@format json
    //@request LoginRequest
    //@response 200 common.Response{code=0,msg="success",data=LoginResponse} "success description" //mock
    //@response 200 common.Response{code=10020,msg="password_error"} "error description"
    //@author alovn
    func (c *Controller) Login(c *bytego.Ctx) {
        //bind LoginRequest
        res := common.NewResponse(0, "success", &LoginResponse{
            WelcomeMsg: "welcome",
        })
        c.JSON(http.StatusOK, res)
    }
    ```

3. 执行 `apidocgen`.

    ```bash
    apidocgen
    ```

    执行后，markdown格式的文档将默认生成在 `./docs` 目录下.

## 示例

更完整的示例在这里: [apidocgen/examples](https://github.com/alovn/apidocgen/tree/main/examples).

这是一个用`apidocs`模板生成的在线浏览的网站模板示例: [https://apidocgen.netlify.app/](https://apidocgen.netlify.app/)

## 注解(注释)

目前支持的注解有下面这些：

annotation|description|example
--|--|--
service|Required, api服务标识|@service svc-user
baseurl|所有api的url前缀|@baseurl /
group|api的分组, 可以添加到api的注释上，或添加到文件头部的注释(该文件中所有的api都会添加到该分组下,api中可省略该注解).|@group account
title|service、group、api的标题|@title UserService
desc|service、group、api的详细描述信息|@desc xxx
api|api的标识|@api POST /account/register
order|group、api的排序(从小到大)|@order 1
author|api的作者(维护者)|@author alovn
version|service、api的版本号|@version 1.0.1
accept|reuest请求的数据格式, 支持 json、xml|@accept json
format|response输出的数据格式, 支持 json/xml|@format json
request|request请求对象|@request LoginRequest
response|response输出对象。 [http code] [data type]|@response 200 LoginResponse
success|同@response|@success 200 LoginResponse
failure|同@response|@failure 200 LoginResponse
param|URL中的path路由参数 `/user/:id`。参数由空格分割 [name] [type] [required] [comment],|@param id int true "user_id"
query|URL中的query请求参数, `/user?id=`,格式同 @param|@query id int true "user_id"
header|header请求参数, 格式同 @param|@header X-Request-ID string false "request id"
form|The form parameter of request, 格式同 @param|@form id int true "user_id"
deprecated|标记api过期.|@deprecated

## Response

1. 输出自定义的结构体.

    ```go
    type User struct {
        ID int `json:"id"`
        Name string `json:"name"`
    }

    //@response 200 User "description"
    ```

2. 输出结构体的数组.

    ```go
    //@response 200 []User
    ```

3. 输出一个结构体，并替换结构体的字段.

    ```go
    type Response struct {
        Code int `json:"code"`
        Msg string `json:"msg"`
        Data interface{} `json:"data"`
    }

    //@response 200 Response{code=10001,msg="some error"} "some error description"
    //@response 200 Response{code=0,data=User} "success description"
    //@response 200 Response{code=0,data=[]User} "success description"
    ```

    如果 Response 结构体是在common包中, 就要使用 `common.Response`:

    ```go
    import (
        "common"
    )

    //@response 200 common.Response{code=0,data=User} "success description"
    ```

4. 添加example tag, 为结构体输出示例值。

    ```go
    type User struct {
        ID int `json:"id" example:"100010"`
        Name string `json:"name" example:"user name"`
    } //User Infomation
    ```

## Mock

1. 为api生成mock文件. 在 api 的 `@response` 注解后添加 `//mock` 注释`，如果没有则默认使用第一个.

    ```go
    //@response 200 Response{code=0,data=[]User} "success description" //mock
    ```

    执行以下命令生成mock数据文件，默认生成的在 `./docs/mocks` 目录下.

    ```bash
    apidocgen --dir=common,svc-user --gen-mocks
    ```

2. 使用生成的mock数据，启动一个mock server.

    ```bash
    apidocgen mock --data=./mocks --listen=:8001
    ```

3. 也可以不生成mock数据文件，直接启动一个mock server.

    ```bash
    apidocgen --dir=common,svc-user --mock --mock-listen=:8001
    ```

4. 测试、使用mock server

    ```bash
    $ curl -X POST -d "" localhost:8001/user/account/login   
    {
      "code": 0,
      "data": {
        "welcome_msg": "string"
      },
      "msg": "success"
    }
    ```
