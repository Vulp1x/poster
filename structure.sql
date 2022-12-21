--
-- PostgreSQL database dump
--

-- Dumped from database version 14.5 (Homebrew)
-- Dumped by pg_dump version 14.5

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- Name: medias_kind; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.medias_kind AS ENUM (
    'photo',
    'reels'
);


--
-- Name: pgqueue_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.pgqueue_status AS ENUM (
    'new',
    'must_retry',
    'no_attempts_left',
    'cancelled',
    'succeeded'
);


--
-- Name: targets_interaction; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.targets_interaction AS ENUM (
    'none',
    'post_description',
    'photo_tag'
);


--
-- Name: targets_status; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.targets_status AS ENUM (
    'new',
    'in_progress',
    'failed',
    'notified'
);


SET default_table_access_method = heap;

--
-- Name: bot_accounts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bot_accounts (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    task_id uuid NOT NULL,
    username text NOT NULL,
    password text NOT NULL,
    user_agent text NOT NULL,
    device_data jsonb NOT NULL,
    session jsonb NOT NULL,
    headers jsonb NOT NULL,
    res_proxy jsonb,
    work_proxy jsonb,
    status smallint NOT NULL,
    posts_count smallint DEFAULT '-1'::integer NOT NULL,
    started_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone,
    file_order integer NOT NULL,
    inst_id bigint NOT NULL,
    CONSTRAINT not_empty_device CHECK (((device_data <> '[]'::jsonb) AND (device_data <> '{}'::jsonb))),
    CONSTRAINT not_empty_headers CHECK (((headers <> '[]'::jsonb) AND (headers <> '{}'::jsonb))),
    CONSTRAINT not_empty_session CHECK (((session <> '[]'::jsonb) AND (session <> '{}'::jsonb)))
);


--
-- Name: goose_db_version_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.goose_db_version_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: logs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.logs (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    bot_id uuid NOT NULL,
    operation text NOT NULL,
    request jsonb NOT NULL,
    response jsonb NOT NULL,
    response_code integer NOT NULL,
    request_time timestamp without time zone NOT NULL,
    proxy_url text
);


--
-- Name: medias; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.medias (
    id bigint NOT NULL,
    kind public.medias_kind NOT NULL,
    inst_id text NOT NULL,
    bot_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: medias_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.medias_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: medias_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.medias_id_seq OWNED BY public.medias.id;


--
-- Name: pgqueue; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.pgqueue (
    id bigint NOT NULL,
    kind smallint NOT NULL,
    payload bytea NOT NULL,
    external_key text,
    status public.pgqueue_status DEFAULT 'new'::public.pgqueue_status NOT NULL,
    messages text[] DEFAULT ARRAY[]::text[] NOT NULL,
    attempts_left smallint NOT NULL,
    attempts_elapsed smallint DEFAULT 0 NOT NULL,
    delayed_till timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
)
WITH (fillfactor='80');


--
-- Name: pgqueue_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.pgqueue_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: pgqueue_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.pgqueue_id_seq OWNED BY public.pgqueue.id;


--
-- Name: proxies; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.proxies (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    task_id uuid NOT NULL,
    assigned_to uuid,
    host text NOT NULL,
    port integer NOT NULL,
    login text NOT NULL,
    pass text NOT NULL,
    type smallint NOT NULL
);


--
-- Name: python_bots; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.python_bots (
    session_id text NOT NULL,
    settings jsonb NOT NULL
);


--
-- Name: target_users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.target_users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    task_id uuid NOT NULL,
    username text NOT NULL,
    user_id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone,
    media_fk bigint,
    status public.targets_status DEFAULT 'new'::public.targets_status NOT NULL,
    interaction_type public.targets_interaction DEFAULT 'none'::public.targets_interaction NOT NULL
);


--
-- Name: target_users_to_tasks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.target_users_to_tasks (
    target_id uuid NOT NULL,
    task_id uuid NOT NULL,
    notified_at timestamp without time zone
);


--
-- Name: tasks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tasks (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    manager_id uuid NOT NULL,
    text_template text NOT NULL,
    landing_accounts text[] NOT NULL,
    account_profile_images bytea[],
    account_names text[],
    account_urls text[],
    images bytea[] NOT NULL,
    status smallint NOT NULL,
    title text NOT NULL,
    bots_filename text,
    cheap_proxies_filename text,
    res_proxies_filename text,
    targets_filename text,
    created_at timestamp with time zone NOT NULL,
    started_at timestamp with time zone,
    stopped_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    account_last_names text[],
    follow_targets boolean DEFAULT false NOT NULL,
    need_photo_tags boolean DEFAULT false NOT NULL,
    per_post_sleep_seconds integer DEFAULT 0 NOT NULL,
    photo_tags_delay_seconds integer DEFAULT 0 NOT NULL,
    type smallint DEFAULT 0 NOT NULL,
    video_filename text,
    posts_per_bot smallint DEFAULT 0 NOT NULL,
    targets_per_post smallint DEFAULT 0 NOT NULL,
    photo_tags_posts_per_bot smallint DEFAULT 0 NOT NULL,
    photo_targets_per_post smallint DEFAULT 0 NOT NULL
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    login text NOT NULL,
    password_hash text NOT NULL,
    role smallint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    CONSTRAINT valid_role CHECK ((role = ANY (ARRAY[0, 1])))
);


--
-- Name: medias id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.medias ALTER COLUMN id SET DEFAULT nextval('public.medias_id_seq'::regclass);


--
-- Name: pgqueue id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pgqueue ALTER COLUMN id SET DEFAULT nextval('public.pgqueue_id_seq'::regclass);


--
-- Name: bot_accounts bot_accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bot_accounts
    ADD CONSTRAINT bot_accounts_pkey PRIMARY KEY (id);


--
-- Name: bot_accounts bot_accounts_username_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bot_accounts
    ADD CONSTRAINT bot_accounts_username_key UNIQUE (username);


--
-- Name: logs logs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.logs
    ADD CONSTRAINT logs_pkey PRIMARY KEY (id);


--
-- Name: medias medias_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.medias
    ADD CONSTRAINT medias_pkey PRIMARY KEY (id);


--
-- Name: pgqueue pgqueue_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.pgqueue
    ADD CONSTRAINT pgqueue_pkey PRIMARY KEY (id);


--
-- Name: proxies proxies_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.proxies
    ADD CONSTRAINT proxies_pkey PRIMARY KEY (id);


--
-- Name: python_bots python_bots_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.python_bots
    ADD CONSTRAINT python_bots_pkey PRIMARY KEY (session_id);


--
-- Name: target_users target_users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.target_users
    ADD CONSTRAINT target_users_pkey PRIMARY KEY (id);


--
-- Name: tasks tasks_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT tasks_pkey PRIMARY KEY (id);


--
-- Name: bot_accounts uniq_file_order; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bot_accounts
    ADD CONSTRAINT uniq_file_order UNIQUE (task_id, file_order);


--
-- Name: users users_login_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_login_key UNIQUE (login);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: pgqueue_broken_tasks_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX pgqueue_broken_tasks_idx ON public.pgqueue USING btree (kind, created_at) WHERE (status = 'no_attempts_left'::public.pgqueue_status);


--
-- Name: pgqueue_idempotency_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX pgqueue_idempotency_idx ON public.pgqueue USING btree (kind, external_key);


--
-- Name: pgqueue_open_tasks_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX pgqueue_open_tasks_idx ON public.pgqueue USING btree (kind, delayed_till) WHERE (status = ANY (ARRAY['new'::public.pgqueue_status, 'must_retry'::public.pgqueue_status]));


--
-- Name: pgqueue_terminal_tasks_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX pgqueue_terminal_tasks_idx ON public.pgqueue USING btree (kind, updated_at) WHERE (status = ANY (ARRAY['cancelled'::public.pgqueue_status, 'succeeded'::public.pgqueue_status]));


--
-- Name: target_users_uniq_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX target_users_uniq_idx ON public.target_users USING btree (task_id, username, user_id);


--
-- Name: targets_task_with_status_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX targets_task_with_status_idx ON public.target_users USING btree (task_id, status);


--
-- Name: targets_user_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX targets_user_id_idx ON public.target_users USING btree (task_id, user_id);


--
-- Name: bot_accounts bot_accounts_task_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bot_accounts
    ADD CONSTRAINT bot_accounts_task_id_fkey FOREIGN KEY (task_id) REFERENCES public.tasks(id);


--
-- Name: logs logs_bot_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.logs
    ADD CONSTRAINT logs_bot_id_fkey FOREIGN KEY (bot_id) REFERENCES public.bot_accounts(id);


--
-- Name: medias medias_bot_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.medias
    ADD CONSTRAINT medias_bot_id_fkey FOREIGN KEY (bot_id) REFERENCES public.bot_accounts(id);


--
-- Name: proxies proxies_assigned_to_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.proxies
    ADD CONSTRAINT proxies_assigned_to_fkey FOREIGN KEY (assigned_to) REFERENCES public.bot_accounts(id);


--
-- Name: proxies proxies_task_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.proxies
    ADD CONSTRAINT proxies_task_id_fkey FOREIGN KEY (task_id) REFERENCES public.tasks(id);


--
-- Name: target_users target_users_media_fk_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.target_users
    ADD CONSTRAINT target_users_media_fk_fkey FOREIGN KEY (media_fk) REFERENCES public.medias(id);


--
-- Name: target_users target_users_task_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.target_users
    ADD CONSTRAINT target_users_task_id_fkey FOREIGN KEY (task_id) REFERENCES public.tasks(id);


--
-- Name: target_users_to_tasks target_users_to_tasks_target_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.target_users_to_tasks
    ADD CONSTRAINT target_users_to_tasks_target_id_fkey FOREIGN KEY (target_id) REFERENCES public.target_users(id);


--
-- Name: target_users_to_tasks target_users_to_tasks_task_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.target_users_to_tasks
    ADD CONSTRAINT target_users_to_tasks_task_id_fkey FOREIGN KEY (task_id) REFERENCES public.tasks(id);


--
-- Name: tasks tasks_manager_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT tasks_manager_id_fkey FOREIGN KEY (manager_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

