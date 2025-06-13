module benchmarks

go 1.24

toolchain go1.24.4

require (
	github.com/cloudresty/emit v0.0.0
	github.com/sirupsen/logrus v1.9.3
	go.uber.org/zap v1.27.0
)

require (
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)

replace github.com/cloudresty/emit => ../
