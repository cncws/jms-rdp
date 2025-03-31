# JumpServer RDP Launcher

加速 JumpServer 连接 RDP 资产，无需登陆网页再连接资产。

在 MacOS(arm64) 和 Windows(amd64) 上测试通过。

## 使用方法

- MacOS 需安装 Windows App；windows 已经内置 mstsc，无需安装
  - `brew install --cask windows-app`
- 在 JumpServer 网页端获取 API Key
- 使用命令打开远程桌面

```bash
# 连接资产
jms-rdp -url URL -ak AK -sk SK -account ACCOUNT IP

# 查看使用帮助
jms-rdp
```
