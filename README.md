

简易 Deeplx 免费api 中转项目，统一管理多个端点，请求失败则自动降级    

### Usage

重命名 `config.yaml.tmpl` 填入自己要管理的端点，`go run .` 启动程序即可

### Feature

- 支持设置超时时间
- 兼容原有 api


### 后续计划
- [ ] 自动排序模式。阶段性 ping 一次（可调），维护到内存，然后动态维护一个权重排名
- [ ] header 鉴权支持
- [ ] 重试实现（好像没必要？再观察是否方便实现）
- [ ] 并发模式。并发获取结果，选择第一个获取响应的
- [ ] 暂时基本考虑应该不兼容别家翻译 api 中转降级，但是自家的降级可以考虑兼容适配一下。（把各家翻译接口大一统，弱化具体到底调用了的考虑，谁只关心翻译可用性，即做一个翻译api界的oneapi，确实感觉是一个非常不错的需求。但是接口兼容性比想象中的重要，目前说白了，没法用到沉浸式翻译里面毫无价值，暂时不徒增难度了。后期要做的话基本也是往 deeplx-api-compatiable 方向去做）
- [ ] 增加前端（？在犹豫，前端有前端的好处，调整配置方便。但是也有坏处，配置复杂，依赖db维护。目前判断轻量级更重要。然后也再探索一下怎么动态插拔一种在线配置的方式，重启还是很麻烦的。）