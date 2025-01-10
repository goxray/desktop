fyne-cross darwin -arch=arm64
xattr -c ./fyne-cross/dist/darwin-arm64/GoXRay.app
echo "If you encounter 'damaged' error, run this 'xattr -c \"GoXRay.app\"' to remove macOS quarantine flag for externally downloaded apps." > ./fyne-cross/dist/darwin-arm64/README.txt
tar -cJf ./fyne-cross/GoXRay_darwin_arm64.tar.xz -C ./fyne-cross/dist/darwin-arm64 .

fyne-cross darwin -arch=amd64
xattr -c ./fyne-cross/dist/darwin-amd64/GoXRay.app
echo "If you encounter 'damaged' error, run this 'xattr -c \"GoXRay.app\"' to remove macOS quarantine flag for externally downloaded apps." > ./fyne-cross/dist/darwin-amd64/README.txt
tar -cJf ./fyne-cross/GoXRay_darwin_amd64.tar.xz -C ./fyne-cross/dist/darwin-amd64 .

fyne-cross linux -arch=arm64
echo "After installation, set the binary privileges for network access with: 'sudo setcap cap_net_raw,cap_net_admin,cap_net_bind_service+eip goxray_binary_path'" > ./fyne-cross/dist/linux-arm64/README.txt
unxz ./fyne-cross/dist/linux-arm64/GoXRay.tar.xz
tar -rvf ./fyne-cross/dist/linux-arm64/GoXRay.tar -C ./fyne-cross/dist/linux-arm64 README.txt
cp ./fyne-cross/dist/linux-arm64/GoXRay.tar ./fyne-cross/GoXRay_linux_arm64.tar

fyne-cross linux -arch=amd64
echo "After installation, set the binary privileges for network access with: 'sudo setcap cap_net_raw,cap_net_admin,cap_net_bind_service+eip goxray_binary_path'" > ./fyne-cross/dist/linux-amd64/README.txt
unxz ./fyne-cross/dist/linux-amd64/GoXRay.tar.xz
tar -rvf ./fyne-cross/dist/linux-amd64/GoXRay.tar -C ./fyne-cross/dist/linux-amd64 README.txt
cp ./fyne-cross/dist/linux-amd64/GoXRay.tar ./fyne-cross/GoXRay_linux_amd64.tar
