go test
export GOOS=linux
go build -o main main.go
$GOPATH\\bin\\build-lambda-zip.exe -output main.zip main