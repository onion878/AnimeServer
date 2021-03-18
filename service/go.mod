module service

go 1.16

require (
	github.com/gin-gonic/gin v1.6.3
	structs v0.0.0
	utils v0.0.0
)

replace service v0.0.0 => ../service

replace structs v0.0.0 => ../structs

replace utils v0.0.0 => ../utils
