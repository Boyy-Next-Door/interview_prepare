# 经历
2021年本科毕业于湖南大学计算机科学与技术专业，校招进入腾讯CSIG医疗健康事业部，从事toB、toG的医疗、公卫相关系统的后台开发。

# 技术栈
第一年主要使用nodejs进行开发，近半年开始转golang。

经手的业务有部署在公有云上的SaaS版本，也有基于k8s部署的私有化版本，深度参与过项目私有化部署迁移的工程改造与基础建设。

- 熟悉MySQL，有一定的优化经验。

- 熟悉redis，对其中大多数数据结构有过实际应用。

- 了解Kafka、RocketMQ、RabbitMQ等消息队列，有过基本的使用。

- 有过k8s的基本使用经验，参与过CI/CD流程的搭建。

- 了解一些基本的前端技术栈，如HTML、js、CSS、Vue。利用低代码平台+自定义模块，开发过业务后管系统。

# 项目经验
1. 深圳市流行病学调查处置系统

    是一个针对流行病阳性病例的调查处置系统，核心功能包括：阳性病人的人工和自动上报、流调任务的部署分发、在线电话流调、报告生成、轨迹信息分析几大块功能。

    我在其中深度参与了几乎所有功能模块，比较完整的有：
    - 统一登录系统，负责流调系统以及上下游多个相关系统的账号托管、API鉴权、角色权限管理、单点登录。
    - 在线语音电话ASR转录推送，涉及到websocket服务从公有云到私有云的迁移改造、基于k8s的横向扩容改造；流调核心信息提取。
    - 流调报告生成器。基于docx文件底层的OOXML结构，对目标内容进行替换和逻辑处理，支持任意样式、内容的docx文件模板快速生成流调报告。
    - 对外API网关，托管系统所对接上下游系统的主调、被调API。简单实现了接口导入、验签鉴权、参数校验、限流等功能。
    - 深度参与业务接口的优化。包括业务逻辑、缓存设计、SQL优化等。
    - 参与CI/CD流程的搭建，将系统发版时间从四个小时提升到二十分钟。
2. 家庭医生智能助手

    目前属于孵化demo阶段，利用openai的API，基于家庭医生与患者的聊天对话，提取患者主诉，生成诊疗记录、给处合理的诊断、治疗和其他健康建议。

    个人负责了整个demo后端的所有逻辑，其中包含了：
    - 利用消息队列，实现对话内容实时拉取与分发。
    - 实现API主调服务，集成了会话上下文管理与消息持久化相关功能。
    - 实现了http代理服务器，内部支持流式传输和关键词黑白名单功能。
