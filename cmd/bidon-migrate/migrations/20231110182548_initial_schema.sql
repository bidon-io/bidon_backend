-- +goose Up
-- +goose StatementBegin


--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: app_demand_profiles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.app_demand_profiles
(
    id               bigint                         NOT NULL,
    app_id           bigint                         NOT NULL,
    account_type     character varying              NOT NULL,
    account_id       bigint                         NOT NULL,
    demand_source_id bigint                         NOT NULL,
    data             jsonb DEFAULT '{}'::jsonb,
    created_at       timestamp(6) without time zone NOT NULL,
    updated_at       timestamp(6) without time zone NOT NULL,
    public_uid       bigint
);


--
-- Name: app_demand_profiles_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.app_demand_profiles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: app_demand_profiles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.app_demand_profiles_id_seq OWNED BY public.app_demand_profiles.id;


--
-- Name: app_mmp_profiles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.app_mmp_profiles
(
    id                                   bigint                         NOT NULL,
    app_id                               bigint                         NOT NULL,
    start_date                           date                           NOT NULL,
    mmp_platform                         integer DEFAULT 0,
    primary_mmp_account                  bigint,
    secondary_mmp_account                bigint,
    get_spend_from_secondary_mmp_account boolean DEFAULT false,
    primary_mmp_raw_data_source          integer,
    secondary_mmp_raw_data_source        integer,
    adjust_app_token                     character varying,
    adjust_s2s_token                     character varying,
    adjust_environment                   character varying,
    appsflyer_dev_key                    character varying,
    appsflyer_app_id                     character varying,
    appsflyer_conversion_keys            character varying,
    firebase_config_keys                 character varying,
    firebase_expiration_duration         integer,
    firebase_tracking                    boolean DEFAULT false,
    facebook_tracking                    boolean DEFAULT false,
    created_at                           timestamp(6) without time zone NOT NULL,
    updated_at                           timestamp(6) without time zone NOT NULL
);


--
-- Name: app_mmp_profiles_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.app_mmp_profiles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: app_mmp_profiles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.app_mmp_profiles_id_seq OWNED BY public.app_mmp_profiles.id;


--
-- Name: apps; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.apps
(
    id           bigint                         NOT NULL,
    user_id      bigint                         NOT NULL,
    platform_id  integer                        NOT NULL,
    human_name   character varying              NOT NULL,
    package_name character varying,
    app_key      character varying,
    settings     jsonb DEFAULT '{}'::jsonb,
    created_at   timestamp(6) without time zone NOT NULL,
    updated_at   timestamp(6) without time zone NOT NULL,
    public_uid   bigint
);

ALTER TABLE ONLY public.apps
    REPLICA IDENTITY FULL;


--
-- Name: apps_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.apps_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: apps_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.apps_id_seq OWNED BY public.apps.id;


--
-- Name: ar_internal_metadata; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ar_internal_metadata
(
    key        character varying              NOT NULL,
    value      character varying,
    created_at timestamp(6) without time zone NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: auction_configurations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.auction_configurations
(
    id                         bigint                         NOT NULL,
    name                       character varying,
    app_id                     bigint                         NOT NULL,
    ad_type                    integer                        NOT NULL,
    rounds                     jsonb   DEFAULT '[]'::jsonb,
    status                     integer,
    settings                   jsonb   DEFAULT '{}'::jsonb,
    pricefloor                 double precision               NOT NULL,
    created_at                 timestamp(6) without time zone NOT NULL,
    updated_at                 timestamp(6) without time zone NOT NULL,
    segment_id                 bigint,
    external_win_notifications boolean DEFAULT false          NOT NULL,
    public_uid                 bigint
);


--
-- Name: auction_configurations_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.auction_configurations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: auction_configurations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.auction_configurations_id_seq OWNED BY public.auction_configurations.id;


--
-- Name: countries; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.countries
(
    id           bigint                         NOT NULL,
    alpha_2_code character varying              NOT NULL,
    alpha_3_code character varying              NOT NULL,
    human_name   character varying,
    created_at   timestamp(6) without time zone NOT NULL,
    updated_at   timestamp(6) without time zone NOT NULL
);


--
-- Name: countries_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.countries_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: countries_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.countries_id_seq OWNED BY public.countries.id;


--
-- Name: demand_source_accounts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.demand_source_accounts
(
    id               bigint                         NOT NULL,
    demand_source_id bigint                         NOT NULL,
    user_id          bigint,
    type             character varying              NOT NULL,
    extra            jsonb   DEFAULT '{}'::jsonb,
    bidding          boolean DEFAULT false,
    is_default       boolean,
    created_at       timestamp(6) without time zone NOT NULL,
    updated_at       timestamp(6) without time zone NOT NULL,
    label            character varying,
    public_uid       bigint
);


--
-- Name: demand_source_accounts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.demand_source_accounts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: demand_source_accounts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.demand_source_accounts_id_seq OWNED BY public.demand_source_accounts.id;


--
-- Name: demand_sources; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.demand_sources
(
    id         bigint                         NOT NULL,
    api_key    character varying              NOT NULL,
    human_name character varying              NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL
);


--
-- Name: demand_sources_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.demand_sources_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: demand_sources_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.demand_sources_id_seq OWNED BY public.demand_sources.id;


--
-- Name: line_items; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.line_items
(
    id           bigint                         NOT NULL,
    app_id       bigint                         NOT NULL,
    account_type character varying              NOT NULL,
    account_id   bigint                         NOT NULL,
    human_name   character varying              NOT NULL,
    code         character varying              NOT NULL,
    bid_floor    numeric,
    ad_type      integer                        NOT NULL,
    extra        jsonb   DEFAULT '{}'::jsonb,
    created_at   timestamp(6) without time zone NOT NULL,
    updated_at   timestamp(6) without time zone NOT NULL,
    width        integer DEFAULT 0              NOT NULL,
    height       integer DEFAULT 0              NOT NULL,
    format       character varying,
    public_uid   bigint,
    bidding      boolean
);


--
-- Name: line_items_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.line_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: line_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.line_items_id_seq OWNED BY public.line_items.id;


--
-- Name: mmp_accounts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.mmp_accounts
(
    id                   bigint                         NOT NULL,
    user_id              bigint                         NOT NULL,
    human_name           character varying              NOT NULL,
    account_type         integer                        NOT NULL,
    use_s3               boolean DEFAULT false,
    s3_access_key_id     character varying,
    s3_secret_access_key character varying,
    s3_bucket_name       character varying,
    s3_region            character varying,
    s3_home_folder       character varying,
    master_api_token     character varying,
    user_token           character varying,
    is_global_account    boolean DEFAULT false,
    created_at           timestamp(6) without time zone NOT NULL,
    updated_at           timestamp(6) without time zone NOT NULL
);


--
-- Name: mmp_accounts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.mmp_accounts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: mmp_accounts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.mmp_accounts_id_seq OWNED BY public.mmp_accounts.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations
(
    version character varying NOT NULL
);


--
-- Name: segments; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.segments
(
    id          bigint                         NOT NULL,
    name        character varying              NOT NULL,
    description text                           NOT NULL,
    filters     jsonb   DEFAULT '[]'::jsonb    NOT NULL,
    enabled     boolean DEFAULT true           NOT NULL,
    app_id      bigint                         NOT NULL,
    created_at  timestamp(6) without time zone NOT NULL,
    updated_at  timestamp(6) without time zone NOT NULL,
    priority    integer DEFAULT 0              NOT NULL,
    public_uid  bigint
);


--
-- Name: segments_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.segments_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: segments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.segments_id_seq OWNED BY public.segments.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users
(
    id            bigint                         NOT NULL,
    email         character varying              NOT NULL,
    created_at    timestamp(6) without time zone NOT NULL,
    updated_at    timestamp(6) without time zone NOT NULL,
    is_admin      boolean DEFAULT false          NOT NULL,
    password_hash character varying              NOT NULL,
    public_uid    bigint
);


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: app_demand_profiles id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.app_demand_profiles
    ALTER COLUMN id SET DEFAULT nextval('public.app_demand_profiles_id_seq'::regclass);


--
-- Name: app_mmp_profiles id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.app_mmp_profiles
    ALTER COLUMN id SET DEFAULT nextval('public.app_mmp_profiles_id_seq'::regclass);


--
-- Name: apps id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.apps
    ALTER COLUMN id SET DEFAULT nextval('public.apps_id_seq'::regclass);


--
-- Name: auction_configurations id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.auction_configurations
    ALTER COLUMN id SET DEFAULT nextval('public.auction_configurations_id_seq'::regclass);


--
-- Name: countries id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.countries
    ALTER COLUMN id SET DEFAULT nextval('public.countries_id_seq'::regclass);


--
-- Name: demand_source_accounts id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.demand_source_accounts
    ALTER COLUMN id SET DEFAULT nextval('public.demand_source_accounts_id_seq'::regclass);


--
-- Name: demand_sources id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.demand_sources
    ALTER COLUMN id SET DEFAULT nextval('public.demand_sources_id_seq'::regclass);


--
-- Name: line_items id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.line_items
    ALTER COLUMN id SET DEFAULT nextval('public.line_items_id_seq'::regclass);


--
-- Name: mmp_accounts id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.mmp_accounts
    ALTER COLUMN id SET DEFAULT nextval('public.mmp_accounts_id_seq'::regclass);


--
-- Name: segments id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.segments
    ALTER COLUMN id SET DEFAULT nextval('public.segments_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: app_demand_profiles app_demand_profiles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.app_demand_profiles
    ADD CONSTRAINT app_demand_profiles_pkey PRIMARY KEY (id);


--
-- Name: app_mmp_profiles app_mmp_profiles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.app_mmp_profiles
    ADD CONSTRAINT app_mmp_profiles_pkey PRIMARY KEY (id);


--
-- Name: apps apps_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.apps
    ADD CONSTRAINT apps_pkey PRIMARY KEY (id);


--
-- Name: ar_internal_metadata ar_internal_metadata_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ar_internal_metadata
    ADD CONSTRAINT ar_internal_metadata_pkey PRIMARY KEY (key);


--
-- Name: auction_configurations auction_configurations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.auction_configurations
    ADD CONSTRAINT auction_configurations_pkey PRIMARY KEY (id);


--
-- Name: countries countries_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.countries
    ADD CONSTRAINT countries_pkey PRIMARY KEY (id);


--
-- Name: demand_source_accounts demand_source_accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.demand_source_accounts
    ADD CONSTRAINT demand_source_accounts_pkey PRIMARY KEY (id);


--
-- Name: demand_sources demand_sources_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.demand_sources
    ADD CONSTRAINT demand_sources_pkey PRIMARY KEY (id);


--
-- Name: line_items line_items_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.line_items
    ADD CONSTRAINT line_items_pkey PRIMARY KEY (id);


--
-- Name: mmp_accounts mmp_accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.mmp_accounts
    ADD CONSTRAINT mmp_accounts_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: segments segments_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.segments
    ADD CONSTRAINT segments_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: index_app_demand_profiles_on_account; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_app_demand_profiles_on_account ON public.app_demand_profiles USING btree (account_type, account_id);


--
-- Name: index_app_demand_profiles_on_app_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_app_demand_profiles_on_app_id ON public.app_demand_profiles USING btree (app_id);


--
-- Name: index_app_demand_profiles_on_app_id_and_demand_source_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_app_demand_profiles_on_app_id_and_demand_source_id ON public.app_demand_profiles USING btree (app_id, demand_source_id);


--
-- Name: index_app_demand_profiles_on_demand_source_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_app_demand_profiles_on_demand_source_id ON public.app_demand_profiles USING btree (demand_source_id);


--
-- Name: index_app_demand_profiles_on_public_uid; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_app_demand_profiles_on_public_uid ON public.app_demand_profiles USING btree (public_uid);


--
-- Name: index_app_mmp_profiles_on_app_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_app_mmp_profiles_on_app_id ON public.app_mmp_profiles USING btree (app_id);


--
-- Name: index_apps_on_app_key; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_apps_on_app_key ON public.apps USING btree (app_key);


--
-- Name: index_apps_on_package_name_and_platform_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_apps_on_package_name_and_platform_id ON public.apps USING btree (package_name, platform_id);


--
-- Name: index_apps_on_public_uid; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_apps_on_public_uid ON public.apps USING btree (public_uid);


--
-- Name: index_apps_on_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_apps_on_user_id ON public.apps USING btree (user_id);


--
-- Name: index_auction_configurations_on_app_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_auction_configurations_on_app_id ON public.auction_configurations USING btree (app_id);


--
-- Name: index_auction_configurations_on_public_uid; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_auction_configurations_on_public_uid ON public.auction_configurations USING btree (public_uid);


--
-- Name: index_auction_configurations_on_segment_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_auction_configurations_on_segment_id ON public.auction_configurations USING btree (segment_id);


--
-- Name: index_countries_on_alpha_2_code; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_countries_on_alpha_2_code ON public.countries USING btree (alpha_2_code);


--
-- Name: index_countries_on_alpha_3_code; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_countries_on_alpha_3_code ON public.countries USING btree (alpha_3_code);


--
-- Name: index_demand_source_accounts_on_demand_source_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_demand_source_accounts_on_demand_source_id ON public.demand_source_accounts USING btree (demand_source_id);


--
-- Name: index_demand_source_accounts_on_public_uid; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_demand_source_accounts_on_public_uid ON public.demand_source_accounts USING btree (public_uid);


--
-- Name: index_demand_sources_on_api_key; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_demand_sources_on_api_key ON public.demand_sources USING btree (api_key);


--
-- Name: index_line_items_on_account; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_line_items_on_account ON public.line_items USING btree (account_type, account_id);


--
-- Name: index_line_items_on_app_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_line_items_on_app_id ON public.line_items USING btree (app_id);


--
-- Name: index_line_items_on_public_uid; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_line_items_on_public_uid ON public.line_items USING btree (public_uid);


--
-- Name: index_mmp_accounts_on_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_mmp_accounts_on_user_id ON public.mmp_accounts USING btree (user_id);


--
-- Name: index_segments_on_app_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX index_segments_on_app_id ON public.segments USING btree (app_id);


--
-- Name: index_segments_on_public_uid; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_segments_on_public_uid ON public.segments USING btree (public_uid);


--
-- Name: index_users_on_email; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_users_on_email ON public.users USING btree (email);


--
-- Name: index_users_on_public_uid; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX index_users_on_public_uid ON public.users USING btree (public_uid);


--
-- Name: line_items fk_rails_0dfd52f3ee; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.line_items
    ADD CONSTRAINT fk_rails_0dfd52f3ee FOREIGN KEY (app_id) REFERENCES public.apps (id);


--
-- Name: mmp_accounts fk_rails_2325069c61; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.mmp_accounts
    ADD CONSTRAINT fk_rails_2325069c61 FOREIGN KEY (user_id) REFERENCES public.users (id);


--
-- Name: app_demand_profiles fk_rails_30ca2bbbfc; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.app_demand_profiles
    ADD CONSTRAINT fk_rails_30ca2bbbfc FOREIGN KEY (demand_source_id) REFERENCES public.demand_sources (id);


--
-- Name: auction_configurations fk_rails_350c02af85; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.auction_configurations
    ADD CONSTRAINT fk_rails_350c02af85 FOREIGN KEY (segment_id) REFERENCES public.segments (id);


--
-- Name: app_mmp_profiles fk_rails_398e0a7e82; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.app_mmp_profiles
    ADD CONSTRAINT fk_rails_398e0a7e82 FOREIGN KEY (app_id) REFERENCES public.apps (id);


--
-- Name: segments fk_rails_42e22d5b29; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.segments
    ADD CONSTRAINT fk_rails_42e22d5b29 FOREIGN KEY (app_id) REFERENCES public.apps (id);


--
-- Name: demand_source_accounts fk_rails_464e6b9aaa; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.demand_source_accounts
    ADD CONSTRAINT fk_rails_464e6b9aaa FOREIGN KEY (demand_source_id) REFERENCES public.demand_sources (id);


--
-- Name: apps fk_rails_995ae3be76; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.apps
    ADD CONSTRAINT fk_rails_995ae3be76 FOREIGN KEY (user_id) REFERENCES public.users (id);


--
-- Name: auction_configurations fk_rails_9a69f14fed; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.auction_configurations
    ADD CONSTRAINT fk_rails_9a69f14fed FOREIGN KEY (app_id) REFERENCES public.apps (id);


--
-- Name: app_demand_profiles fk_rails_dbb170baa2; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.app_demand_profiles
    ADD CONSTRAINT fk_rails_dbb170baa2 FOREIGN KEY (account_id) REFERENCES public.demand_source_accounts (id);


--
-- Name: app_demand_profiles fk_rails_e5a2299135; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.app_demand_profiles
    ADD CONSTRAINT fk_rails_e5a2299135 FOREIGN KEY (app_id) REFERENCES public.apps (id);


--
-- Name: line_items fk_rails_f1f73a548d; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.line_items
    ADD CONSTRAINT fk_rails_f1f73a548d FOREIGN KEY (account_id) REFERENCES public.demand_source_accounts (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP EXTENSION IF EXISTS pgcrypto CASCADE;

DROP TABLE public.app_demand_profiles CASCADE;
DROP TABLE public.app_mmp_profiles CASCADE;
DROP TABLE public.apps CASCADE;
DROP TABLE public.auction_configurations CASCADE;
DROP TABLE public.countries CASCADE;
DROP TABLE public.demand_source_accounts CASCADE;
DROP TABLE public.demand_sources CASCADE;
DROP TABLE public.line_items CASCADE;
DROP TABLE public.mmp_accounts CASCADE;
DROP TABLE public.segments CASCADE;
DROP TABLE public.users CASCADE;

-- Rails leftovers
DROP TABLE public.ar_internal_metadata CASCADE;
DROP TABLE public.schema_migrations CASCADE;
-- +goose StatementEnd
