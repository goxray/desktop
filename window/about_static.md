# Go VPN client for XRay
[https://github.com/goxray](https://github.com/goxray)

This project brings fully functioning [XRay](https://github.com/XTLS/Xray-core) VPN client implementation in Go.

#### What is XRay?
Please visit [https://xtls.github.io/en](https://xtls.github.io/en) for more info.

## How it works
- Application sets up new TUN device;
- Adds additional routes to route all system traffic to this newly created TUN device
- Adds exception for XRay outbound address (basically your VPN server IP)
- Tunnel is created to process all incoming IP packets via TCP/IP stack
- All outbound traffic is routed through the XRay inbound proxy
- All incoming packets are routed back via TUN device