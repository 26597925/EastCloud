set path=%path%;E:\GoWork\src\sapi\tools\protoc-3.13.0\protoc-win64\bin

// protoc --go_out=plugins=grpc:{输出目录}  {proto文件}
// protoc --go_out=plugins=grpc:./test/ ./test.proto