<p align="center">
  <a href="https://github.com/jurycu/umi-dva-antd-admin">
    <img width="100" src="https://github.com/jurycu/probix/blob/main/probix-monitor.png">
  </a>
</p>

<h1 align="center">Probix</h1>

<div align="center">
基于prometheus-operator的通用黑盒监控方案。
</div>

## 🌿 Probix是什么？
probix是一个云原生通用的黑盒监控解决方案，用户不用关心如何集成prometheus，只需根据业务的监控指标，写自己的相关脚本，吐出相关json数据，probix就可将其转换成metrics指标，可无缝接入prometheus，来达到对业务的黑盒探测。

## 🗣名词解释
`黑盒监控`：主要关注的现象，一般都是正在发生的东西，例如出现一个告警，业务接口不正常，那么这种监控就是站在用户的角度能看到的监控，重点在于能对正在发生的故障进行告警。

`白盒监控`：主要关注的是原因，也就是系统内部暴露的一些指标，例如redis的info中显示redis slave down，这个就是redis info显示的一个内部的指标，重点在于原因，可能是在黑盒监控中看到redis down，而查看内部信息的时候，显示redis port is refused connection。

## ✨ 特性

- 🛡 **无缝接入**: 完全基于原生prometheus-operator，依赖少
- 💎 **使用简单**：创建一个黑盒监控，只是简单配置CRD
- 🚀 **最新技术栈**：使用CRD,参考blackbox的理念设计
- 🔢 **黑盒扩展**：填补blackbox不能自定义监控的缺点



## 📜 架构

<img src="https://github.com/jurycu/probix/blob/main/probix-arch.png" /> 

## 🎉 部署
````
kubectl apply -f config/deployment/probix.yaml
````

## 📚  效果
用户吐出的json指标
````
{
	"values": [{
		"labels": {
			"aas": "as",
			"tt": "ss"
		},
		"value": 12
	}, {
		"labels": {
			"aas": "aas",
			"tt": "sff"
		},
		"value": 15
	}, {
		"labels": {
			"aas": "asga",
			"tt": "asfa"
		},
		"value": 11
	}]
}
````

probix转换后的效果
````
# HELP probix_duration_seconds Returns how long the probe took to complete in seconds
# TYPE probix_duration_seconds gauge
probix_duration_seconds{method="GET"} 0.008774518
# HELP probix_success Displays whether or not the probix was a success
# TYPE probix_success gauge
probix_success{method="GET"} 1
# HELP tt aa
# TYPE tt gauge
tt{aas="aas",tt="sff"} 15
tt{aas="as",tt="ss"} 12
tt{aas="asga",tt="asfa"} 11
````

## 👁 CRD示例
````
apiVersion: ferulax.jurycu.io/v1alpha1
kind: Probix
metadata:
  name: probix-sample
  namespace: probix
spec:
  interval: 30s #拉取时间间隔
  scrapeTimeout: 30s  #超时时间
  targets:
    - metricsName: test_ss  #指标名
      metricsHelp: test  指标说明
      target: http://xxx.com #接口定义
      method: POST #请求方法
      body: '{"foo": "bar"}' #请求body
````

## ⌨️ 本地开发

```sh
# 克隆项目到本地
git clone https://github.com/jurycu/probix

# 安装依赖
go mod tidy

# 启动服务
make run
```


## 👥 社区互助

| Github Issue                                      | 钉钉                                                                                     | 微信                                                                                   |
| ------------------------------------------------- | ------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------------------- |
| [issues](https://github.com/jurycu/probix/issues) | <img src="https://github.com/jurycu/umi-dva-antd-admin/blob/main/src/assets/dingtalk.jpg" width="100" /> | <img src="https://github.com/jurycu/umi-dva-antd-admin/blob/main/src/assets/wechat.png" width="100" /> |
