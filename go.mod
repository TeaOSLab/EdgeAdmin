module github.com/TeaOSLab/EdgeAdmin

go 1.16

replace github.com/TeaOSLab/EdgeCommon => ../EdgeCommon

require (
	github.com/TeaOSLab/EdgeCommon v0.0.0-00010101000000-000000000000
	github.com/cespare/xxhash v1.1.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/iwind/TeaGo v0.0.0-20220304042834-d8b5103e9027
	github.com/iwind/gosock v0.0.0-20211103081026-ee4652210ca4
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/miekg/dns v1.1.35
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/tealeg/xlsx/v3 v3.2.3
	github.com/xlzd/gotp v0.0.0-20181030022105-c8557ba2c119
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220227234510-4e6760a101f9
	google.golang.org/grpc v1.44.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)
