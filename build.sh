#!/bin/bash

echo "编译中..."
mkdir -p dist
wails build -clean -o image_proce.exe -upx -webview2 download
echo "复制DLL文件..."
cp -r dll/* build/bin/
echo "复制编译文件..."
cp -r build/bin/* dist/
echo "编译完成"