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
    started_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone,
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
    request jsonb NOT NULL,
    response jsonb NOT NULL,
    response_code integer NOT NULL,
    request_time timestamp without time zone NOT NULL,
    proxy_url text
);


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
-- Name: target_users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.target_users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    task_id uuid NOT NULL,
    username text NOT NULL,
    user_id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone
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
    image bytea NOT NULL,
    status smallint NOT NULL,
    title text NOT NULL,
    bots_filename text,
    proxies_filename text,
    targets_filename text,
    created_at timestamp with time zone NOT NULL,
    started_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
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
-- Name: proxies proxies_host_port_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.proxies
    ADD CONSTRAINT proxies_host_port_key UNIQUE (host, port);


--
-- Name: proxies proxies_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.proxies
    ADD CONSTRAINT proxies_pkey PRIMARY KEY (id);


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
-- Name: target_users_uniq_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX target_users_uniq_idx ON public.target_users USING btree (task_id, username, user_id);


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

