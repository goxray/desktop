# <img width="25px" src="icon/assets/app.png" alt="app logo"> GoXRay: Desktop VPN client for XRay (macOS)
![Static Badge](https://img.shields.io/badge/OS-macOS-blue?style=flat&logo=apple&logoColor=white&logoSize=auto&color=blue)
![Static Badge](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go&logoColor=white)

This project brings fully functional [XRay](https://github.com/XTLS/Xray-core) desktop client for macOS, powered by [fyne](https://github.com/fyne-io/fyne) and written in Go.

> [!NOTE]
> The program will not damage your routing rules, default route is intact and only additional rules are added for the lifetime of application's TUN device. There are also additional complementary clean up procedures in place.

<img alt="" src=".github/images/overview_screenshot.png">

#### What is XRay?
Please visit https://xtls.github.io/en for more info.

#### Tested and supported on:
- macOS (tested on Sequoia 15.1.1)

> Feel free to test this on your system and let me know in the issues :)

## ✨ Features
- Stupidly easy to use
- Adding and editing XRay URL configurations
- Supports all [Xray-core](https://github.com/XTLS/Xray-core) protocols (vless, vmess e.t.c.) using link notation (`vless://` e.t.c.)
- Real-time network statistics for each configuration
- Responsive, lightweight and dynamic UI, focusing on tray menu for quick and easy interactions
- Only soft routing rules are applied, no changes made to default routes

## ⚡️ Usage

### 🌐 VPN Client

> [!TIP]
> You can easily set the application to start at login.
> Go to `System Settings` > `General` > `Login Items & Extensions` > `Open at Login`, then press `+` and browse for GoXRay.app

Get the latest release app bundle from [Releases](https://github.com/goxray/desktop/releases) and... just run it.
You will be prompted for privileged access, and your GoXRay VPN is ready.
Don't forget to add the app to your `Applications` and `Open at Login` items!

### 🛠️ Building the project

> [!IMPORTANT]
> - `sudo` is required, the application will prompt you for privileged access when you run it.
> - CGO_ENABLED=1 is required.
> - Visit https://docs.fyne.io/started/ to get required dependencies for your platform.

If you have all the required dependencies in place for building the project. You can simply run:
```bash
CGO_ENABLED=1 go build . -o goxray
```

If you want to get the macOS application bundle (.app structured file), you can run:
```bash
# go install fyne.io/fyne/v2/cmd/fyne@latest
fyne package -os darwin # Check out the newly created GoXRay.app file
```

### 📋 TODO
- [ ] Add logs tab with all std logs and export button
- [ ] Add rearranging items functionality
