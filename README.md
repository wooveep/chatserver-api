<!--
 * @Author: cloudyi.li
 * @Date: 2023-05-10 09:15:49
 * @LastEditTime: 2023-05-11 14:50:34
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/README.md
-->


# 基于OPENAI的ChatGPT API开发的AI助手服务

本仓库为后端API服务，依赖的前端仓库为[wooveep](https://github.com/wooveep)/[chatserver-web](https://github.com/wooveep/chatserver-web)


## 系统演示

![操作演示](docs/操作演示.gif)
## 项目截图

#### 用户登录

![用户登录](docs/用户登录.png)

#### 用户注册

![用户注册](docs/用户注册.png)

#### 用户信息展示

![用户信息展示](docs/用户信息展示.png)

#### 会话上下文设置

![会话上下文设置](docs/会话上下文长度设置.png)

#### 会话角色控制

![会话角色控制](docs/会话角色控制.png)

#### 基于本地知识库的问答

![InceptorSQL调优](docs/inceptor调优.png)

![设备配置咨询](docs/设备配置咨询.png)

## 实现功能

- [x] 登录、注册、用户管理
- [x] 用户额度、会员有效期管理
- [x] 服务端保存用户会话和聊天记录
- [x] 会话角色管理
- [x] 长回复功能实现 <!--API返回消息因为TOKEN长度中断时自动处理-->
- [x] 支持结合本地知识库问答
- [x] 多会话储存和上下文逻辑

## 应用场景

我们 的AI助手适用于以下场景：

- 企业用户：结合公开外发客户的产品手册，解答产品问题，优化客服成本。
- 保险行业：结合保险条款内容，解答客户关于保险相关信息，推荐客户更合适保险产品。
- 教育行业：结合常见题库、文本，作为私人家教。

## 目标群体

- 个人部署
- 商业部署
- 企业部署
- 团队部署

## 待实现列表

- [ ] 基于卡密方式的用户额度充值
- [ ] 系统后台管理界面
- [ ] 用户系统设置模块
- [ ] 自定义敏感词
- [ ] 用户角色权限管理
- [ ] 自定义AI角色页面
- [ ] 本地文档向量化操作页面
- [ ] 语音问答

## 感谢

OpenAI API SDK：  [sashabaranov](https://github.com/sashabaranov)/**[go-openai](https://github.com/sashabaranov/go-openai)**

OpenAIt Token 计算：[pkoukk](https://github.com/pkoukk)/**[tiktoken-go](https://github.com/pkoukk/tiktoken-go)**

向量存储库：  [pgvector](https://github.com/pgvector)/**[pgvector](https://github.com/pgvector/pgvector)**

Gin项目框架： [xmgtony](https://github.com/xmgtony)/**[apiserver-gin](https://github.com/xmgtony/apiserver-gin)**

关键词提取：[yanyiwu](https://github.com/yanyiwu)/**[gojieba](https://github.com/yanyiwu/gojieba)**

## License

MIT © [Cloudyi Li](https://github.com/wooveep/chatserver-api/blob/main/LICENSE)
