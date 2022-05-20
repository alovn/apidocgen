# apidocgen

apidocgen is a tool for Go to generate apis markdown docs.

## Install

```bash
go install github.com/alovn/apidocgen@latest
```

## Cli

```bash
$ apidocgen --help
apidocgen is a tool for Go to generate apis markdown docs.

Usage:
  apidocgen --dir= --excludes= --output= --output-index= --template= --single

Flags:
        --dir:          Search apis dir, comma separated, default .
        --excludes:     Exclude directories and files when searching, comma separated
        --output:       Generate markdown files dir, default ./docs/
        --output-index: Generate index file name.
        --template:     Template name or custom template directory, built-in includes markdown and apidocs, default markdown.
        --single:       Generate a single markdown file.
```

built-in templates include `markdown` and `apidocs`, default is `markdown`.

run the command in the go module directory.

```bash
cd your-api-service-dir
apidocgen \
    --dir=svc-user,common \
    --output=./docs

apidocgen \
    --dir=svc-user,common \
    --template=apidocs \
    --output=./docs

apidocgen --output-index=index.md //generate index.md
apidocgen --output-index=@{service}.md //generate index file name by your @service comment, example: svc-user.md.
```

## Template

The built-in includes `markdown` and `apidocs`.

The built-in template `apidocs` is the template for generate website [apidocs](https://github.com/alovn/apidocs).

You can also use the custom template:

```bash
apidocgen \
    --dir=svc-user,common \
    --template=/Users/xxx/workspace/apidocs/custom-template-direcoty \
    --output=./docs
```

## How to use

apidocgen supported any web frameworks. here are an example using [bytego](https://github.com/gostack-labs/bytego).

1. Add API annotations in main.go code:

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

2. Add API annotations in httpHandler.

    ```go
    //@title AccountLogin
    //@api POST /account/login
    //@group account
    //@accept json
    //@format json
    //@request LoginRequest
    //@response 200 common.Response{code=0,msg="success",data=LoginResponse} "登录成功返回数据"
    //@response 200 common.Response{code=10020,msg="password_error"} "密码错误"
    //@author alovn
    func (c *Controller) Login(c *bytego.Ctx) {
        //bind LoginRequest
        res := common.NewResponse(0, "success", &LoginResponse{
            WelcomeMsg: "welcome",
        })
        c.JSON(http.StatusOK, res)
    }
    ```

3. Execute `apidocgen`.

    ```bash
    apidocgen
    ```

    The markdown files will be generated in `./docs`.

## Examples

Some examples and generated markdown docs are here: [apidocgen/examples](https://github.com/alovn/apidocgen/tree/main/examples).

An online generated api docs site: [https://apidocgen.netlify.app/](https://apidocgen.netlify.app/)

## Comments(annotations)

annotation|description|example
--|--|--
service|Required, the service's identification|@service svc-user
baseurl|The base url of api|@baseurl /
group|The group of service|@group account
title|The title of service, group and api|@title UserService
desc|The description of service, group and api|@desc xxx
api|The http api handler|@api POST /account/register
order|Sort groups and apis|@order 1
author|The author of api|@author alovn
version|The version of service or api|@version 1.0.1
accept|The request format, support json/xml|@accept json
format|The response format, support json/xml|@format json
request|The request body|@request LoginRequest
response|The response body, [http code] [data type]|@response 200 LoginResponse
success|As same as response|@success 200 LoginResponse
failure|As same as response|@failure 200 LoginResponse
param|The path parameter of router `/user/:id`, parameters separated by spaces [name] [type] [required] [comment],|@param id int true "user_id"
query|The query parameter of route, `/user?id=`, parameters same as @param|@query id int true "user_id"
header|The request HTTP header parameter, parameters same as @param|@header X-Request-ID string false "request id"
form|The form parameter of request, parameters same as @param|@form id int true "user_id"
deprecated|Mark api as deprecated.|@deprecated

## Response

1. response user defined struct.

    ```go
    type User struct {
        ID int `json:"id"`
        Name string `json:"name"`
    }

    //response 200 User "description"
    ```

2. response struct with array.

    ```go
    //response 200 []User
    ```

3. a composition common response.

    ```go
    type Response struct {
        Code int `json:"code"`
        Msg string `json:"msg"`
        Data interface{} `json:"data"`
    }

    //response 200 Response{code=10001,msg="some error"} "some error description"
    //response 200 Response{code=0,data=User} "success description"
    //response 200 Response{code=0,data=[]User} "success description"
    ```

    if import package of `common.Response`:

    ```go
    import (
        "common"
    )
    //response 200 common.Response{code=0,data=User} "success description"
    ```

4. example value of struct

    ```go
    type User struct {
        ID int `json:"id" example:"100010"`
        Name string `json:"name" example:"user name"`
    } //User Infomation
    ```
