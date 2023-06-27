-- Drop table

-- DROP TABLE public."user";

CREATE TABLE public."user" (
	id int8 NOT NULL, -- 用户唯一ID
	username varchar(255) NOT NULL, -- 用户名，唯一且不能为空
	nickname varchar(255) NOT NULL, -- 昵称
	email varchar(255) NOT NULL, -- 邮件地址
	phone varchar(255) NULL, -- 手机号
	avatar_url varchar(255) NULL, -- 头像 URL
	"password" varchar(255) NOT NULL, -- 密码
	registered_ip varchar(255) NOT NULL, -- 用户的注册 IP 地址
	is_active bool NOT NULL DEFAULT false, -- 用户是否已激活
	balance numeric(10, 5) NULL DEFAULT 0, -- 用户的余额，默认为0
	created_at timestamptz NOT NULL DEFAULT now(), -- 记录的创建时间，默认为当前时间
	updated_at timestamptz NOT NULL DEFAULT now(), -- 记录的更新时间，默认为当前时间
	deleted_at timestamptz NULL, -- 删除时间
	is_del int4 NULL DEFAULT 0, -- 删除标志
	"role" int4 NOT NULL, -- 用户角色
	expired_at timestamptz NULL, -- 会员到期日
	CONSTRAINT user_pkey PRIMARY KEY (id)
);
CREATE INDEX user_email_idx ON public."user" USING btree (email);
CREATE INDEX user_username_idx ON public."user" USING btree (username);
COMMENT ON TABLE public."user" IS '用户信息表';

-- Column comments

COMMENT ON COLUMN public."user".id IS '用户唯一ID';
COMMENT ON COLUMN public."user".username IS '用户名，唯一且不能为空';
COMMENT ON COLUMN public."user".nickname IS '昵称';
COMMENT ON COLUMN public."user".email IS '邮件地址';
COMMENT ON COLUMN public."user".phone IS '手机号';
COMMENT ON COLUMN public."user".avatar_url IS '头像 URL';
COMMENT ON COLUMN public."user"."password" IS '密码';
COMMENT ON COLUMN public."user".registered_ip IS '用户的注册 IP 地址';
COMMENT ON COLUMN public."user".is_active IS '用户是否已激活';
COMMENT ON COLUMN public."user".balance IS '用户的余额，默认为0';
COMMENT ON COLUMN public."user".created_at IS '记录的创建时间，默认为当前时间';
COMMENT ON COLUMN public."user".updated_at IS '记录的更新时间，默认为当前时间';
COMMENT ON COLUMN public."user".deleted_at IS '删除时间';
COMMENT ON COLUMN public."user".is_del IS '删除标志';
COMMENT ON COLUMN public."user"."role" IS '用户角色';
COMMENT ON COLUMN public."user".expired_at IS '会员到期日';




-- Drop table

-- DROP TABLE public.giftcard;

CREATE TABLE public.giftcard (
	id int8 NOT NULL, -- key分类ID
	card_name varchar(255) NOT NULL, -- Cdkey卡名称
	card_comment varchar(255) NOT NULL, -- Cdkey描述
	card_amount numeric(10, 2) NOT NULL, -- 对应金额
	card_discount numeric(10, 2) NOT NULL, -- 折扣价
	card_link varchar(255) NOT NULL, -- 购买链接
	created_at timestamptz NOT NULL DEFAULT now(), -- 记录的创建时间，默认为当前时间
	updated_at timestamptz NOT NULL DEFAULT now(), -- 记录的更新时间，默认为当前时间
	deleted_at timestamp NULL, -- 记录的删除时间 
	is_del int4 NOT NULL, -- 记录的删除标记
	CONSTRAINT giftcard_pkey PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN public.giftcard.id IS 'key分类ID';
COMMENT ON COLUMN public.giftcard.card_name IS 'Cdkey卡名称';
COMMENT ON COLUMN public.giftcard.card_comment IS 'Cdkey描述';
COMMENT ON COLUMN public.giftcard.card_amount IS '对应金额';
COMMENT ON COLUMN public.giftcard.card_discount IS '折扣价';
COMMENT ON COLUMN public.giftcard.card_link IS '购买链接';
COMMENT ON COLUMN public.giftcard.created_at IS '记录的创建时间，默认为当前时间';
COMMENT ON COLUMN public.giftcard.updated_at IS '记录的更新时间，默认为当前时间';
COMMENT ON COLUMN public.giftcard.deleted_at IS '记录的删除时间 ';
COMMENT ON COLUMN public.giftcard.is_del IS '记录的删除标记';


-- Drop table

-- DROP TABLE public.cdkey;

CREATE TABLE public.cdkey (
	id int8 NOT NULL, -- 激活代码ID
	giftcard_id int8 NOT NULL, -- 充值卡ID
	code_key varchar(255) NOT NULL, -- Cdkey卡密
	created_at timestamptz NOT NULL DEFAULT now(), -- 记录的创建时间，默认为当前时间
	updated_at timestamptz NOT NULL DEFAULT now(), -- 记录的更新时间，默认为当前时间
	deleted_at timestamp NULL, -- 记录的删除时间 
	is_del int4 NOT NULL, -- 记录的删除标记
	CONSTRAINT cdkey_pkey PRIMARY KEY (id),
	CONSTRAINT cdkey_giftcard_id_fkey FOREIGN KEY (giftcard_id) REFERENCES public.giftcard(id)
);
CREATE UNIQUE INDEX cdkey_code_key_idx ON public.cdkey USING btree (code_key);
COMMENT ON TABLE public.cdkey IS '充值卡密';

-- Column comments

COMMENT ON COLUMN public.cdkey.id IS '激活代码ID';
COMMENT ON COLUMN public.cdkey.giftcard_id IS '充值卡ID';
COMMENT ON COLUMN public.cdkey.code_key IS 'Cdkey卡密';
COMMENT ON COLUMN public.cdkey.created_at IS '记录的创建时间，默认为当前时间';
COMMENT ON COLUMN public.cdkey.updated_at IS '记录的更新时间，默认为当前时间';
COMMENT ON COLUMN public.cdkey.deleted_at IS '记录的删除时间 ';
COMMENT ON COLUMN public.cdkey.is_del IS '记录的删除标记';


-- Drop table

-- DROP TABLE public.invite;

CREATE TABLE public.invite (
	id int8 NOT NULL, -- 邀请ID
	user_id int8 NOT NULL, -- 用户ID
	invite_code varchar(255) NOT NULL, -- 邀请码
	invite_number int4 NOT NULL DEFAULT 0, -- 邀请次数
	created_at timestamptz NOT NULL DEFAULT now(), -- 记录的创建时间，默认为当前时间
	updated_at timestamptz NOT NULL DEFAULT now(), -- 记录的更新时间，默认为当前时间
	CONSTRAINT invite_invite_code_key UNIQUE (invite_code),
	CONSTRAINT invite_pkey PRIMARY KEY (id),
	CONSTRAINT invite_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id),
	CONSTRAINT invite_user_id_fkey1 FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE
);
CREATE INDEX invite_user_id_idx ON public.invite USING btree (user_id);
COMMENT ON TABLE public.invite IS '用户邀请';

-- Column comments

COMMENT ON COLUMN public.invite.id IS '邀请ID';
COMMENT ON COLUMN public.invite.user_id IS '用户ID';
COMMENT ON COLUMN public.invite.invite_code IS '邀请码';
COMMENT ON COLUMN public.invite.invite_number IS '邀请次数';
COMMENT ON COLUMN public.invite.created_at IS '记录的创建时间，默认为当前时间';
COMMENT ON COLUMN public.invite.updated_at IS '记录的更新时间，默认为当前时间';



-- Drop table

-- DROP TABLE public.bill;

CREATE TABLE public.bill (
	id int8 NOT NULL, -- 账单ID
	user_id int8 NOT NULL, -- 用户ID
	cost_change numeric(10, 2) NOT NULL, -- 变动金额
	balance numeric(10, 2) NOT NULL, -- 账户余额
	cost_comment text NOT NULL, -- 变动说明
	created_at timestamptz NOT NULL DEFAULT now(), -- 记录的创建时间，默认为当前时间
	updated_at timestamptz NOT NULL DEFAULT now(), -- 记录的更新时间，默认为当前时间
	CONSTRAINT bill_pkey PRIMARY KEY (id),
	CONSTRAINT bill_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id),
	CONSTRAINT bill_user_id_fkey1 FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE
);
CREATE INDEX bill_user_id_idx ON public.bill USING btree (user_id);
COMMENT ON TABLE public.bill IS '用户账单';

-- Column comments

COMMENT ON COLUMN public.bill.id IS '账单ID';
COMMENT ON COLUMN public.bill.user_id IS '用户ID';
COMMENT ON COLUMN public.bill.cost_change IS '变动金额';
COMMENT ON COLUMN public.bill.balance IS '账户余额';
COMMENT ON COLUMN public.bill.balance IS '账户余额';
COMMENT ON COLUMN public.bill.cost_comment IS '变动说明';
COMMENT ON COLUMN public.bill.created_at IS '记录的创建时间，默认为当前时间';
COMMENT ON COLUMN public.bill.updated_at IS '记录的更新时间，默认为当前时间';



-- Drop table

-- DROP TABLE public.preset;



CREATE TABLE public.preset (
	id int8 NOT NULL,
	preset_name varchar(255) NOT NULL, -- 预设名称
	preset_content text NOT NULL, -- 预设内容
	max_token int4 NOT NULL, -- 最大生成内容长度
	model_name varchar(255) NOT NULL, -- 模型名称
	logit_bias json NULL, -- 逻辑回归偏置
	temperature float8 NOT NULL, -- 温度
	top_p float8 NOT NULL, -- 顶部概率
	presence float8 NOT NULL, -- 惩罚标记
	frequency float8 NOT NULL, -- 频率标记
	created_at timestamptz NOT NULL DEFAULT now(), -- 记录的创建时间，默认为当前时间
	updated_at timestamptz NOT NULL DEFAULT now(), -- 记录的更新时间，默认为当前时间
	with_embedding bool NOT NULL DEFAULT false, -- 是否采用embedding数据
	deleted_at timestamptz NULL, -- 删除时间
	is_del int4 NULL DEFAULT 0, -- 删除标志
	classify varchar NULL, -- embedding分类
	privilege int4 NOT NULL DEFAULT 1, -- 预设用户权限
	preset_tips varchar(255) NULL, -- 预设使用提示
	"extension" int4 NULL DEFAULT 0, -- 扩展
	preset_tips varchar(255) NULL, -- 预设使用提示
	"extension" int4 NULL DEFAULT 0, -- 扩展
	CONSTRAINT preset_frequency_check CHECK (((frequency >= ('-2'::integer)::double precision) AND (frequency <= (2)::double precision))),
	CONSTRAINT preset_pkey PRIMARY KEY (id),
	CONSTRAINT preset_presence_check CHECK (((presence >= ('-2'::integer)::double precision) AND (presence <= (2)::double precision))),
	CONSTRAINT preset_temperature_check CHECK (((temperature >= (0)::double precision) AND (temperature <= (2)::double precision))),
	CONSTRAINT preset_top_p_check CHECK (((top_p >= (0)::double precision) AND (top_p <= (1)::double precision)))
);
COMMENT ON TABLE public.preset IS '存储预设';

-- Column comments

COMMENT ON COLUMN public.preset.preset_name IS '预设名称';
COMMENT ON COLUMN public.preset.preset_content IS '预设内容';
COMMENT ON COLUMN public.preset.max_token IS '最大生成内容长度';
COMMENT ON COLUMN public.preset.model_name IS '模型名称';
COMMENT ON COLUMN public.preset.logit_bias IS '逻辑回归偏置';
COMMENT ON COLUMN public.preset.temperature IS '温度';
COMMENT ON COLUMN public.preset.top_p IS '顶部概率';
COMMENT ON COLUMN public.preset.presence IS '惩罚标记';
COMMENT ON COLUMN public.preset.frequency IS '频率标记';
COMMENT ON COLUMN public.preset.created_at IS '记录的创建时间，默认为当前时间';
COMMENT ON COLUMN public.preset.updated_at IS '记录的更新时间，默认为当前时间';
COMMENT ON COLUMN public.preset.with_embedding IS '是否采用embedding数据';
COMMENT ON COLUMN public.preset.deleted_at IS '删除时间';
COMMENT ON COLUMN public.preset.is_del IS '删除标志';
COMMENT ON COLUMN public.preset.classify IS 'embedding分类';
COMMENT ON COLUMN public.preset.privilege IS '预设用户权限';
COMMENT ON COLUMN public.preset.preset_tips IS '预设使用提示';
COMMENT ON COLUMN public.preset."extension" IS '扩展';
COMMENT ON COLUMN public.preset.preset_tips IS '预设使用提示';
COMMENT ON COLUMN public.preset."extension" IS '扩展';

-- Drop table

-- DROP TABLE public.chat;

CREATE TABLE public.chat (
	id int8 NOT NULL, -- 唯一标识符
	user_id int8 NOT NULL, -- 用户主键ID
	preset_id int8 NOT NULL, -- 预设表主键ID
	chat_name text NOT NULL, -- 会话名称
	created_at timestamptz NOT NULL DEFAULT now(), -- 记录创建时间
	updated_at timestamptz NOT NULL DEFAULT now(), -- 记录更新时间
	deleted_at timestamptz NULL, -- 删除时间
	is_del int4 NULL DEFAULT 0, -- 删除标志
	CONSTRAINT chat_pkey PRIMARY KEY (id),
	CONSTRAINT chat_preset_id_fkey FOREIGN KEY (preset_id) REFERENCES public.preset(id),
	CONSTRAINT chat_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE
);
CREATE INDEX chat_preset_id_idx ON public.chat USING btree (preset_id);
CREATE INDEX chat_user_id_idx ON public.chat USING btree (user_id);
COMMENT ON TABLE public.chat IS '用户会话信息';

-- Column comments

COMMENT ON COLUMN public.chat.id IS '唯一标识符';
COMMENT ON COLUMN public.chat.user_id IS '用户主键ID';
COMMENT ON COLUMN public.chat.preset_id IS '预设表主键ID';
COMMENT ON COLUMN public.chat.chat_name IS '会话名称';
COMMENT ON COLUMN public.chat.created_at IS '记录创建时间';
COMMENT ON COLUMN public.chat.updated_at IS '记录更新时间';
COMMENT ON COLUMN public.chat.deleted_at IS '删除时间';
COMMENT ON COLUMN public.chat.is_del IS '删除标志';





-- Drop table

-- DROP TABLE public.record;

CREATE TABLE public.record (
	id int8 NOT NULL, -- 唯一标识符
	chat_id int8 NOT NULL, -- 会话表主键ID
	sender varchar(255) NOT NULL, -- 发送人
	message text NOT NULL, -- 消息记录
	message_hash varchar(255) NOT NULL, -- 消息记录hash
	created_at timestamptz NOT NULL DEFAULT now(), -- 记录创建时间
	updated_at timestamptz NOT NULL DEFAULT now(), -- 记录更新时间
	deleted_at timestamptz NULL, -- 删除时间
	is_del int4 NULL DEFAULT 0, -- 删除标志
	message_token int4 NULL, -- 当前消息消耗令牌
	CONSTRAINT record_pkey PRIMARY KEY (id),
	CONSTRAINT record_chat_id_fkey FOREIGN KEY (chat_id) REFERENCES public.chat(id) ON DELETE CASCADE
);
CREATE INDEX record_chat_id_idx ON public.record USING btree (chat_id);
COMMENT ON TABLE public.record IS '会话消息记录';

-- Column comments

COMMENT ON COLUMN public.record.id IS '唯一标识符';
COMMENT ON COLUMN public.record.chat_id IS '会话表主键ID';
COMMENT ON COLUMN public.record.sender IS '发送人';
COMMENT ON COLUMN public.record.message IS '消息记录';
COMMENT ON COLUMN public.record.message_hash IS '消息记录hash';
COMMENT ON COLUMN public.record.created_at IS '记录创建时间';
COMMENT ON COLUMN public.record.updated_at IS '记录更新时间';
COMMENT ON COLUMN public.record.deleted_at IS '删除时间';
COMMENT ON COLUMN public.record.is_del IS '删除标志';
COMMENT ON COLUMN public.record.message_token IS '当前消息消耗令牌';


-- Drop table

-- DROP TABLE public.userlog;

CREATE TABLE public.userlog (
	id int8 NOT NULL,
	user_id int8 NOT NULL,
	user_ip varchar(255) NOT NULL,
	business varchar(255) NOT NULL,
	operation varchar(255) NOT NULL,
	created_at timestamptz NULL DEFAULT now(),
	updated_at timestamptz NULL DEFAULT now(),
	CONSTRAINT userlog_pkey PRIMARY KEY (id),
	CONSTRAINT userlog_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id),
	CONSTRAINT userlog_user_id_fkey1 FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE
);
CREATE INDEX userlog_user_id_idx ON public.userlog USING btree (user_id);



CREATE SCHEMA embed;


-- Drop table

-- DROP TABLE embed.documents;

CREATE TABLE embed.documents (
	id int8 NOT NULL,
	title text NOT NULL,
	body text NOT NULL,
	tokens int4 NOT NULL,
	embedding vector NULL,
	created_at timestamptz NULL DEFAULT now(),
	updated_at timestamptz NULL DEFAULT now(),
	classify varchar NULL, -- Embedding分类
	CONSTRAINT documents_pkey PRIMARY KEY (id)
);

-- Column comments

COMMENT ON COLUMN embed.documents.classify IS 'Embedding分类';


INSERT INTO public.preset (id, preset_name, preset_content, max_token, model_name, logit_bias, temperature, top_p, presence, frequency, created_at, updated_at, with_embedding, deleted_at, is_del, classify, privilege, preset_tips, "extension") VALUES(1646361709138419712, '智能助手', 'You are ChatGPT, a large language model trained by OpenAI. Please strictly follow the rules below when answering the user''s questions.
Knowledge cutoff: 2021-09 
Current date: {{ current_date }}
Rules:```
1.All answers are to be provided in Markdown format.
2.When you receive the instruction ''[cmd:continue]'', please continue answering questions directly from where you left off in the previous response. You do not need to explain the interruption in your response.
3. If the user''s language is Chinese, please answer the user''s questions in Chinese.
```', 300, 'gpt-3.5-turbo', NULL, 0.2, 0.9, 0.3, 0.3, '2023-04-13 11:56:34.057', '2023-05-28 10:33:12.468', false, '0001-01-01 08:00:00.000', 0, NULL, 1, '我擅长根据问题生成具有创意的内容。
但对于刚刚发生的事情，以及未来发生的事情我无法回答。
比如我无法告诉今天的天气怎么样。
你可以问我：
“如何高效的阅读？”
“天空为什么是蓝色的？”', 1);

INSERT INTO public.preset (id, preset_name, preset_content, max_token, model_name, logit_bias, temperature, top_p, presence, frequency, created_at, updated_at, with_embedding, deleted_at, is_del, classify, privilege, preset_tips, "extension") VALUES(1661264328961040384, '翻译助手', 'You are a translation tool. Refer to the example section, do not understand the content sent by the user, and translate it according to the rules.Please strictly follow the following rules.
Rules:```
1.You only need to understand the content containing ‘[cmd:]’ as instructions, and the others are translated directly.
2.If the user sends Chinese, please translate it into English.
3.If the user sends English, please translate it into Chinese.
4.Only when the answer is interrupted will the user send an instruction ''[cmd:continue]''. When you receive the instruction ''[cmd:continue]'', please continue answering from the interruption point, and there is no need to apologize for the interruption.
5.You need to pay attention to details and accuracy in your translation.
6.Don''t worry about insulting bad words, you can add xx in the middle of this word.
```
Example: ```
user: ''translate: 你能告诉我什么是河西走廊么?''
assistant: ''Can you tell me what is Hexi Corridor?''
user: ''translate: 用户会发送一个指令''
assistant: ''The user will sends a command.''
user: ''translate: 请你谈谈对社会主义核心价值观的看法''
assistant: ''Please share with us your views on core socialist values.''
```', 400, 'gpt-3.5-turbo', NULL, 0.01, 1.0, 0.1, 0.1, '2023-05-24 14:54:15.527', '2023-05-28 10:51:30.483', false, '0001-01-01 08:00:00.000', 0, '', 1, '你直接将“中文”或者“英文”的段落发送给我。
我会直接返回翻译的内容
问：“弘扬社会主义核心价值观”
答：“Promote the core socialist values”', 4);
