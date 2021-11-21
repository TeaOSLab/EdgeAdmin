module github.com/TeaOSLab/EdgeAdmin

go 1.16

replace github.com/TeaOSLab/EdgeCommon => ../EdgeCommon

require (
	github.com/TeaOSLab/EdgeCommon v0.0.0-00010101000000-000000000000
	github.com/cespare/xxhash v1.1.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/iwind/TeaGo v0.0.0-20211026123858-7de7a21cad24
	github.com/iwind/gosock v0.0.0-20211103081026-ee4652210ca4
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/miekg/dns v1.1.35
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/tealeg/xlsx/v3 v3.2.3
	github.com/xlzd/gotp v0.0.0-20181030022105-c8557ba2c119
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22
	google.golang.org/genproto v0.0.0-20210617175327-b9e0b3197ced // indirect
	google.golang.org/grpc v1.38.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)
