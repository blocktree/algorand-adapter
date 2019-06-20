# algorand-adapter

algorand-adapter适配了openwallet.AssetsAdapter接口，给应用提供了底层的区块链协议支持。

## 如何测试

openwtester包下的测试用例已经集成了openwallet钱包体系，创建conf文件，新建ALGO.ini文件，编辑如下内容：

```ini

# algod service
ServerAPI = "http://127.0.0.1:8080"
ServerToken = ""

```

## 资料介绍

### 官网

https://developer.algorand.org

### 区块浏览器

https://algoexplorer.io

### github

https://github.com/algorand
