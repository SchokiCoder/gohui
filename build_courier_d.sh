BIN_NAME="courier_d"
VERSION="v0.2.0"

go build -o "$BIN_NAME" -gcflags "-l -N" -ldflags "-X 'main.Version=$VERSION'" ./courier