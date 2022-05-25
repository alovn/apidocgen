# apidocgen

apidocgen is a tool for Go to generate apis markdown docs and mocks.

## Install

```bash
go install github.com/alovn/apidocgen@latest
```

## Cli

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

## Template

The built-in templates include `markdown` and `apidocs`, default is `markdown`.

The template `apidocs` is the template for generate website [apidocs](https://github.com/alovn/apidocs).

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
group|The group of service, add it at api comment or at the head of file comment.|@group account
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

    //@response 200 User "description"
    ```

2. response struct with array.

    ```go
    //@response 200 []User
    ```

3. a composition common response.

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

    if import package of `common.Response`:

    ```go
    import (
        "common"
    )

    //@response 200 common.Response{code=0,data=User} "success description"
    ```

4. example value of struct

    ```go
    type User struct {
        ID int `json:"id" example:"100010"`
        Name string `json:"name" example:"user name"`
    } //User Infomation
    ```

## Mock

1. generate apis mocks files. add `//mock` at the end of `response`, default use first.

    ```go
    //@response 200 Response{code=0,data=[]User} "success description" //mock
    ```

     generate apis mocks files, default generated in the directory `./docs/mocks`.

    ```bash
    apidocgen --dir=common,svc-user --gen-mocks
    ```

2. serve the mock server from source code.

    ```bash
    apidocgen --dir=common,svc-user --mock --mock-listen=:8001
    ```

3. serve the mock server from generated mock files.

    ```bash
    apidocgen mock --data=./mocks --listen=:8001
    ```

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
