package main

import (
	runtime "github.com/tuxounet/k2-sdk"

	app "github.com/tuxounet/k2-sdk/samples/with_plugins/plug_host/app"
)

// @title			PlugHostApp
// @version		0.0
// @description	This is the Sample api
// @termsOfService	http://swagger.io/terms/
// @contact.name	API Support
// @contact.url	http://www.swagger.io/support
// @contact.email	support@swagger.io
// @license.name	GPL-3.0
// @license.url	http://www.gnu.org/licenses/gpl-3.0.html
func main() {
	runtime.HostApp(
		app.NewApp(),
	)
}
