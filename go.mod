module github.com/tgruben/validator

go 1.14

require (
	github.com/pilosa/go-pilosa v1.3.1-0.20190715210601-8606626b90d6
	github.com/pilosa/pilosa/v2 v2.0.0-alpha.1
)

replace github.com/pilosa/go-pilosa => github.com/molecula/go-pilosa v0.0.0-20200415220604-ab86f987c746

replace github.com/pilosa/pilosa/v2 => github.com/molecula/pilosa/v2 v2.0.0-alpha.16.0.20200415212437-0691b6c015ff
