go build
# rm -rf log/*.log
pkill -f packet-verify
./packet-verify --mode=TCPServer &
./packet-verify --mode=RemoteProxy &
./packet-verify --mode=LocalProxy &