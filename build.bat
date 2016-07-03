@pushd %~dp0
@cd ../../../../
@set GOPATH=%CD%
@cd src/github.com/sryanyuan/gocodesite
@go build
@popd