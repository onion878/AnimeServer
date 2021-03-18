module utils

go 1.16

require (
	github.com/axgle/mahonia v0.0.0-20180208002826-3358181d7394
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-xorm/xorm v0.7.9
	github.com/kr/pretty v0.2.1 // indirect
	github.com/satori/go.uuid v1.2.0
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	structs v0.0.0
	xorm.io/core v0.7.2-0.20190928055935-90aeac8d08eb
)

replace structs v0.0.0 => ../structs