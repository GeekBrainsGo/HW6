/*
 * HomeWork-6: Mongo in BeeGo
 * Created on 28.09.19 22:19
 * Copyright (c) 2019 - Eugene Klimov
 */

package conf

import (
	"fmt"
	"github.com/astaxie/beego"
)

// GetURI returns DSN from config or defaults
func GetURI() string {
	host := beego.AppConfig.String("DBHOST")
	if host == "" {
		host = "localhost"
	}
	port := beego.AppConfig.String("DBPORT")
	if port == "" {
		port = "27017"
	}
	user := beego.AppConfig.String("DBUSER")
	pass := beego.AppConfig.String("DBPASS")

	if user != "" || pass != "" {
		user += ":"
	}

	return fmt.Sprintf("mongodb://%s%s@%s:%s", user, pass, host, port)
}
