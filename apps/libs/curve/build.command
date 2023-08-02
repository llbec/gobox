set GOARCH=386
set CGO_ENABLED=1
go build -ldflags "-s -w"  -o main.dll -buildmode=c-shared