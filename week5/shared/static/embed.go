// Package static содержит встроенные статические файлы для HTTP сервера
package static

import "embed"

//go:embed swagger-ui.html generated/order.swagger.yaml generated/inventory/v1/inventory.swagger.json generated/payment/v1/payment.swagger.json
var FS embed.FS
