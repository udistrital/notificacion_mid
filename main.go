// Copyright 2013 Beego Samples authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// This sample is about using long polling and WebSocket to build a web-based chat room based on beego.
package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	_ "github.com/udistrital/notificacion_mid/routers"
	apistatus "github.com/udistrital/utils_oas/apiStatusLib"
	"github.com/udistrital/utils_oas/customerrorv2"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders: []string{"Origin", "x-requested-with",
			"content-type",
			"accept",
			"origin",
			"authorization",
			"x-csrftoken"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	beego.ErrorController(&customerrorv2.CustomErrorController{})
	apistatus.Init()
	beego.Run()
}

//go:generate sh -c "echo 'package routers; import \"github.com/astaxie/beego\"; func init() {beego.BConfig.RunMode = beego.DEV}' > routers/0.go"
//go:generate sh -c "echo 'package routers; import \"os\"; func init() {os.Exit(0)}' > routers/z.go"
//go:generate go run $GOFILE
//go:generate sh -c "rm routers/0.go routers/z.go"
