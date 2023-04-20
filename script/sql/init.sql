-- DROP SCHEMA public;

CREATE SCHEMA public AUTHORIZATION postgres;

COMMENT ON SCHEMA public IS 'standard public schema';
-- public.preset definition

-- Drop table

-- DROP TABLE public.preset;

CREATE TABLE public.preset (
	id int8 NOT NULL,
	preset_name varchar(255) NOT NULL, -- 预设名称
	preset_content text NULL, -- 预设内容
	max_token int4 NULL, -- 最大生成内容长度
	model_name varchar(255) NULL, -- 模型名称
	logit_bias json NULL, -- 逻辑回归偏置
	temperature float8 NULL, -- 温度
	top_p float8 NULL, -- 顶部概率
	presence float8 NULL, -- 惩罚标记
	frequency float8 NULL, -- 频率标记
	created_at timestamptz NULL DEFAULT now(), -- 记录的创建时间，默认为当前时间
	updated_at timestamptz NULL DEFAULT now(), -- 记录的更新时间，默认为当前时间
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

-- Permissions

ALTER TABLE public.preset OWNER TO whatserver;
GRANT ALL ON TABLE public.preset TO whatserver;


-- public."user" definition

-- Drop table

-- DROP TABLE public."user";

CREATE TABLE public."user" (
	id int8 NOT NULL, -- 用户唯一ID
	username varchar(255) NOT NULL, -- 用户名，唯一且不能为空
	nickname varchar(255) NULL, -- 昵称
	email varchar(255) NOT NULL, -- 邮件地址
	phone varchar(255) NULL, -- 手机号
	avatar_url varchar(255) NULL, -- 头像 URL
	"password" varchar(255) NULL, -- 密码
	expired_at timestamp NULL DEFAULT now(), -- 用户失效时间
	registered_ip varchar(255) NULL, -- 用户的注册 IP 地址
	is_active bool NULL DEFAULT false, -- 用户是否已激活
	balance numeric(10, 5) NULL DEFAULT 0, -- 用户的余额，默认为0
	created_at timestamptz NULL DEFAULT now(), -- 记录的创建时间，默认为当前时间
	updated_at timestamptz NULL DEFAULT now(), -- 记录的更新时间，默认为当前时间
	CONSTRAINT user_email_key UNIQUE (email),
	CONSTRAINT user_pkey PRIMARY KEY (id),
	CONSTRAINT user_username_key UNIQUE (username)
);
COMMENT ON TABLE public."user" IS '用户信息表';

-- Column comments

COMMENT ON COLUMN public."user".id IS '用户唯一ID';
COMMENT ON COLUMN public."user".username IS '用户名，唯一且不能为空';
COMMENT ON COLUMN public."user".nickname IS '昵称';
COMMENT ON COLUMN public."user".email IS '邮件地址';
COMMENT ON COLUMN public."user".phone IS '手机号';
COMMENT ON COLUMN public."user".avatar_url IS '头像 URL';
COMMENT ON COLUMN public."user"."password" IS '密码';
COMMENT ON COLUMN public."user".expired_at IS '用户失效时间';
COMMENT ON COLUMN public."user".registered_ip IS '用户的注册 IP 地址';
COMMENT ON COLUMN public."user".is_active IS '用户是否已激活';
COMMENT ON COLUMN public."user".balance IS '用户的余额，默认为0';
COMMENT ON COLUMN public."user".created_at IS '记录的创建时间，默认为当前时间';
COMMENT ON COLUMN public."user".updated_at IS '记录的更新时间，默认为当前时间';

-- Permissions

ALTER TABLE public."user" OWNER TO whatserver;
GRANT ALL ON TABLE public."user" TO whatserver;


-- public.activecode definition

-- Drop table

-- DROP TABLE public.activecode;

CREATE TABLE public.activecode (
	id int8 NOT NULL, -- 激活代码ID
	user_id int8 NULL, -- 用户主键ID
	active_code text NOT NULL, -- 激活代码
	expire_date timestamptz NULL, -- 过期日期
	created_at timestamptz NULL DEFAULT now(), -- 记录的创建时间，默认为当前时间
	updated_at timestamptz NULL DEFAULT now(), -- 记录的更新时间，默认为当前时间
	CONSTRAINT activecode_pkey PRIMARY KEY (id),
	CONSTRAINT activecode_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id),
	CONSTRAINT activecode_user_id_fkey1 FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE
);
COMMENT ON TABLE public.activecode IS '用户激活代码表';

-- Column comments

COMMENT ON COLUMN public.activecode.id IS '激活代码ID';
COMMENT ON COLUMN public.activecode.user_id IS '用户主键ID';
COMMENT ON COLUMN public.activecode.active_code IS '激活代码';
COMMENT ON COLUMN public.activecode.expire_date IS '过期日期';
COMMENT ON COLUMN public.activecode.created_at IS '记录的创建时间，默认为当前时间';
COMMENT ON COLUMN public.activecode.updated_at IS '记录的更新时间，默认为当前时间';

-- Permissions

ALTER TABLE public.activecode OWNER TO whatserver;
GRANT ALL ON TABLE public.activecode TO whatserver;


-- public.chat definition

-- Drop table

-- DROP TABLE public.chat;

CREATE TABLE public.chat (
	id int8 NOT NULL, -- 唯一标识符
	user_id int8 NULL, -- 用户主键ID
	preset_id int8 NULL, -- 预设表主键ID
	chat_name text NOT NULL, -- 会话名称
	created_at timestamptz NULL DEFAULT now(), -- 记录创建时间
	updated_at timestamptz NULL DEFAULT now(), -- 记录更新时间
	CONSTRAINT chat_pkey PRIMARY KEY (id),
	CONSTRAINT chat_preset_id_fkey FOREIGN KEY (preset_id) REFERENCES public.preset(id),
	CONSTRAINT chat_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE
);
COMMENT ON TABLE public.chat IS '用户会话信息';

-- Column comments

COMMENT ON COLUMN public.chat.id IS '唯一标识符';
COMMENT ON COLUMN public.chat.user_id IS '用户主键ID';
COMMENT ON COLUMN public.chat.preset_id IS '预设表主键ID';
COMMENT ON COLUMN public.chat.chat_name IS '会话名称';
COMMENT ON COLUMN public.chat.created_at IS '记录创建时间';
COMMENT ON COLUMN public.chat.updated_at IS '记录更新时间';

-- Permissions

ALTER TABLE public.chat OWNER TO whatserver;
GRANT ALL ON TABLE public.chat TO whatserver;


-- public.record definition

-- Drop table

-- DROP TABLE public.record;

CREATE TABLE public.record (
	id int8 NOT NULL, -- 唯一标识符
	chat_id int8 NULL, -- 会话表主键ID
	sender text NOT NULL, -- 发送人
	message text NOT NULL, -- 消息记录
	message_hash text NOT NULL, -- 消息记录hash
	created_at timestamptz NULL DEFAULT now(), -- 记录创建时间
	updated_at timestamptz NULL DEFAULT now(), -- 记录更新时间
	CONSTRAINT record_pkey PRIMARY KEY (id),
	CONSTRAINT record_chat_id_fkey FOREIGN KEY (chat_id) REFERENCES public.chat(id) ON DELETE CASCADE
);
COMMENT ON TABLE public.record IS '会话消息记录';

-- Column comments

COMMENT ON COLUMN public.record.id IS '唯一标识符';
COMMENT ON COLUMN public.record.chat_id IS '会话表主键ID';
COMMENT ON COLUMN public.record.sender IS '发送人';
COMMENT ON COLUMN public.record.message IS '消息记录';
COMMENT ON COLUMN public.record.message_hash IS '消息记录hash';
COMMENT ON COLUMN public.record.created_at IS '记录创建时间';
COMMENT ON COLUMN public.record.updated_at IS '记录更新时间';

-- Permissions

ALTER TABLE public.record OWNER TO whatserver;
GRANT ALL ON TABLE public.record TO whatserver;


-- public.settings definition

-- Drop table

-- DROP TABLE public.settings;

CREATE TABLE public.settings (
	id int8 NOT NULL, -- 唯一标识符
	user_id int8 NULL, -- 用户主键ID
	api_key text NULL, -- API密钥
	created_at timestamptz NULL DEFAULT now(), -- 记录创建时间
	updated_at timestamptz NULL DEFAULT now(), -- 记录更新时间
	CONSTRAINT settings_pkey PRIMARY KEY (id),
	CONSTRAINT settings_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id) ON DELETE CASCADE
);
COMMENT ON TABLE public.settings IS '用户设置信息';

-- Column comments

COMMENT ON COLUMN public.settings.id IS '唯一标识符';
COMMENT ON COLUMN public.settings.user_id IS '用户主键ID';
COMMENT ON COLUMN public.settings.api_key IS 'API密钥';
COMMENT ON COLUMN public.settings.created_at IS '记录创建时间';
COMMENT ON COLUMN public.settings.updated_at IS '记录更新时间';

-- Permissions

ALTER TABLE public.settings OWNER TO whatserver;
GRANT ALL ON TABLE public.settings TO whatserver;




-- Permissions

GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO public;


INSERT INTO public.preset (id,preset_name,preset_content,max_token,model_name,logit_bias,temperature,top_p,presence,frequency,created_at,updated_at) VALUES
	 (1646361709138419712,'人工智能AI助手','You are ChatGPT, a large language model trained by OpenAI. Follow the user''s instructions carefully. Respond using markdown.If you do not know the answer to a question that is based on common sense, truthfulness or has specific information, please answer "I don''t know" instead of providing uncertain or potentially incorrect information.',400,'gpt-3.5-turbo',NULL,0.3,1.0,0.1,0.1,'2023-04-13 11:56:34.05797+08','2023-04-13 11:56:34.05797+08');