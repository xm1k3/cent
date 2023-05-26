GOOS=linux GOARCH=amd64 go build -o ./build/cent_linux_amd64
GOOS=linux GOARCH=arm go build -o ./build/cent_linux_arm
GOOS=windows GOARCH=amd64 go build -o ./build/cent_windows_amd64.exe
GOOS=windows GOARCH=386 go build -o ./build/cent_windows_386.exe
GOOS=darwin GOARCH=amd64 go build -o ./build/cent_macos_amd64
