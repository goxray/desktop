go build -o app .
sudo setcap cap_net_raw,cap_net_admin,cap_net_bind_service+eip app
./app