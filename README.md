# <img width="25px" src="icon/assets/app.png" alt="app logo"> GoXRay: Desktop VPN client for XRay [![Static Badge](https://img.shields.io/badge/featured_on_apps.fyne.io-00ADD8?style=flat&logo=graphite&logoColor=lightblue&color=blue&link=https%3A%2F%2Fapps.fyne.io%2Fapps%2Fcom.github.goxray.client.app.html)](https://apps.fyne.io/apps/com.github.goxray.client.app.html)
![Static Badge](https://img.shields.io/badge/OS-macOS%20%7C%20Linux-blue?style=flat&logo=linux&logoColor=white&logoSize=auto&color=blue)
![Static Badge](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go&logoColor=white)
[![Go Report Card](https://goreportcard.com/badge/github.com/goxray/desktop)](https://goreportcard.com/report/github.com/goxray/desktop)
![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/goxray/desktop/total?color=blue)



This project brings fully functional [XRay](https://github.com/XTLS/Xray-core) desktop client for macOS and Linux, powered by [fyne](https://github.com/fyne-io/fyne) and written in Go.

> [!NOTE]
> The program will not damage your routing rules, default route is intact and only additional rules are added for the lifetime of application's TUN device. There are also additional complementary clean up procedures in place.

> For CLI version see https://github.com/goxray/tun.

| macOS (15.1) Light  | Linux (Ubuntu) Dark |
| ------------- | ------------- |
| <img alt="" src=".github/images/overview_screenshot.png"> | <img alt="" src=".github/images/overview_screenshot_linux.jpg">  |

#### What is XRay?
Please visit https://xtls.github.io/en for more info.

#### Tested and supported on:
- macOS (tested on Sequoia 15.1.1 M3)
- Linux:
    - Ubuntu 24.10 ARM
    - Mint 22 amd64

## ✨ Features
- Stupidly easy to use
- Adding and editing XRay URL configurations
- Supports all [Xray-core](https://github.com/XTLS/Xray-core) protocols (vless, vmess e.t.c.) using link notation (`vless://` e.t.c.)
- Real-time network statistics for each configuration
- Responsive, lightweight and dynamic UI, focusing on tray menu for quick and easy interactions
- Only soft routing rules are applied, no changes made to default routes

## ⚡️ Installation

### 🌐 VPN Client

#### macOS

> [!TIP]
> You can easily set the application to start at login.
> Go to `System Settings` > `General` > `Login Items & Extensions` > `Open at Login`, then press `+` and browse for GoXRay.app

Get the latest release app bundle from [Releases](https://github.com/goxray/desktop/releases) and... just run it. If you get "damaged" error then run this command `xattr -c "GoXRay.app"`, apple marks externaly downloaded files, this command removes this mark.
You will be prompted for admin password, and your GoXRay VPN is ready.
Don't forget to add the app to your `Applications` and `Open at Login` items!

#### Linux

Get the latest release from [Releases](https://github.com/goxray/desktop/releases).
After unpacking the archive - you can install the app using `make user-install` or run the binary in `usr/local/bin` folder.

> [!IMPORTANT]
> After installing the application, don't forget to set networking privileges to the installed binary (see `scripts` directory for more info):
> ```bash
> sudo setcap cap_net_raw,cap_net_admin,cap_net_bind_service+eip goxray_binary_path
> ```

##### 📦 Using [twdragon](https://github.com/twdragon) debian package

The client is available from the PPA repository `ppa:twdragon/xray`, maintained by [twdragon](https://github.com/twdragon). The package is in sync with this repo releases. You can check the pipeline at the [dedicated repository](https://github.com/twdragon/xray-debian-pkg). To install use:

```bash
sudo add-apt-repository ppa:twdragon/xray
sudo apt update
sudo apt install goxray-gui
```

After the installation, the package might be updated automatically as it is done in Ubuntu. All the packages in this PPA are cryptographically signed by [twdragon](https://github.com/twdragon). The PPA is also available on Ubuntu Launchpad: https://launchpad.net/~twdragon/+archive/ubuntu/xray

The experimental unsigned artifact packages of the program could be also found in the Github Action Artifacts section of the [pipeline repository](https://github.com/twdragon/xray-debian-pkg/actions).

### 🛠️ Building from source

> [!IMPORTANT]
> - `sudo` is required, the application will prompt you for privileged access when you run it.
> - CGO_ENABLED=1 is required.
> - Visit https://docs.fyne.io/started/ to get required dependencies for your platform.

If you have all the required dependencies in place for building the project. You can simply run:
```bash
CGO_ENABLED=1 go build . -o goxray
```

For Linux, you should also set additional privileges to run the binary, see [#Linux](#linux) info above.

If you want to get the macOS application bundle (.app structured file), you can run:
```bash
# go install fyne.io/fyne/v2/cmd/fyne@latest
fyne package -os darwin # Check out the newly created GoXRay.app file
```

> [!TIP]
> For headache-free cross-compilation please see https://github.com/fyne-io/fyne-cross.

### 📋 Credits

- https://github.com/xtls/xray-core
- https://github.com/fyne-io/fyne
- https://github.com/lilendian0x00/xray-knife
- https://github.com/jackpal/gateway
