module github.com/TeaOSLab/EdgeAdmin

go 1.15

replace github.com/TeaOSLab/EdgeCommon => ../EdgeCommon

require (
	github.com/TeaOSLab/EdgeCommon v0.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/iwind/TeaGo v0.0.0-20201110043415-859f4b3b98f3
	golang.org/x/sys v0.0.0-20200724161237-0e2f3a69832c // indirect
	google.golang.org/grpc v1.32.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)
