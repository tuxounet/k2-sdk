package main

import (
	runtime "github.com/tuxounet/k2-sdk"

	app "{{ .module }}/app"
)

// @title			{{ .name }}
// @description		{{ .description }}
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
