--
-- PostgreSQL database dump
--

-- Dumped from database version 15.2 (Debian 15.2-1.pgdg110+1)
-- Dumped by pg_dump version 15.2

-- Started on 2023-04-01 05:29:04 UTC

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

DROP DATABASE guest_book;
--
-- TOC entry 3370 (class 1262 OID 49214)
-- Name: guest_book; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE guest_book WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';


ALTER DATABASE guest_book OWNER TO postgres;

\connect guest_book

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
-- TOC entry 2 (class 3079 OID 49215)
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- TOC entry 3371 (class 0 OID 0)
-- Dependencies: 2
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 215 (class 1259 OID 49226)
-- Name: threads; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.threads (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    content text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    visit_id uuid NOT NULL,
    user_id uuid NOT NULL
);


ALTER TABLE public.threads OWNER TO postgres;

--
-- TOC entry 216 (class 1259 OID 49233)
-- Name: user_matrices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_matrices (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    endpoint text NOT NULL,
    is_admin boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.user_matrices OWNER TO postgres;

--
-- TOC entry 217 (class 1259 OID 49241)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    full_name text NOT NULL,
    email text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    password text NOT NULL,
    is_admin boolean DEFAULT false NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- TOC entry 218 (class 1259 OID 49249)
-- Name: visits; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.visits (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.visits OWNER TO postgres;

--
-- TOC entry 3361 (class 0 OID 49226)
-- Dependencies: 215
-- Data for Name: threads; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.threads (id, content, created_at, updated_at, deleted_at, visit_id, user_id) VALUES ('bf9db0bd-bfca-4d52-9014-e68efd49dfc4', 'tempat makannya adem, makanannya enak', '2023-04-01 05:20:13.558874+00', '2023-04-01 05:20:13.558874+00', '0001-01-01 00:00:00+00', '33dd9076-783a-4539-af61-96afd28bad46', '36182264-f4a4-4a89-8ec8-589bc6f0e519');
INSERT INTO public.threads (id, content, created_at, updated_at, deleted_at, visit_id, user_id) VALUES ('ca04c881-5153-4ec9-8060-96b88bb6c0a4', 'tempat makan ramah anak', '2023-04-01 05:20:24.723979+00', '2023-04-01 05:20:24.723979+00', '0001-01-01 00:00:00+00', '33dd9076-783a-4539-af61-96afd28bad46', '36182264-f4a4-4a89-8ec8-589bc6f0e519');


--
-- TOC entry 3362 (class 0 OID 49233)
-- Dependencies: 216
-- Data for Name: user_matrices; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.user_matrices (id, endpoint, is_admin, created_at, updated_at, deleted_at) VALUES ('c5315982-f432-4ed4-9b08-227bcc2a70b4', '/api/v1/users', true, '2023-04-01 00:32:53.700288+00', NULL, NULL);
INSERT INTO public.user_matrices (id, endpoint, is_admin, created_at, updated_at, deleted_at) VALUES ('182ea506-55bb-4c8e-a6c4-7c927cf9ecd6', '/api/v1/users', false, '2023-04-01 00:42:06.857335+00', NULL, NULL);
INSERT INTO public.user_matrices (id, endpoint, is_admin, created_at, updated_at, deleted_at) VALUES ('819dc907-9974-48cb-b03a-bada5a1753f4', '/api/v1/users/list', true, '2023-04-01 00:42:06.857335+00', NULL, NULL);
INSERT INTO public.user_matrices (id, endpoint, is_admin, created_at, updated_at, deleted_at) VALUES ('73f3ad28-40a0-40b4-be33-96abe3343c62', '/api/v1/users/[id]', true, '2023-04-01 00:42:06.857335+00', NULL, NULL);
INSERT INTO public.user_matrices (id, endpoint, is_admin, created_at, updated_at, deleted_at) VALUES ('b7a206a1-546d-4665-8b06-e43cb8f3c488', '/api/v1/users/[id]', false, '2023-04-01 00:42:06.857335+00', NULL, NULL);
INSERT INTO public.user_matrices (id, endpoint, is_admin, created_at, updated_at, deleted_at) VALUES ('d223d244-3a37-47e6-9262-baa3412e0b4a', '/api/v1/users/[id]/edit', true, '2023-04-01 00:42:06.857335+00', NULL, NULL);
INSERT INTO public.user_matrices (id, endpoint, is_admin, created_at, updated_at, deleted_at) VALUES ('8699542a-451f-41ac-8a8c-9747324c138b', '/api/v1/users/[id]/edit', false, '2023-04-01 00:42:06.857335+00', NULL, NULL);
INSERT INTO public.user_matrices (id, endpoint, is_admin, created_at, updated_at, deleted_at) VALUES ('c3cccd62-e5cc-4bbf-b8bf-03ac0c34e1ad', '/api/v1/users/[id]/delete', true, '2023-04-01 00:42:06.857335+00', NULL, NULL);
INSERT INTO public.user_matrices (id, endpoint, is_admin, created_at, updated_at, deleted_at) VALUES ('12e7aac3-88ae-41e9-8838-c0f0c222cc0c', '/api/v1/users/[id]/delete', false, '2023-04-01 00:42:06.857335+00', NULL, NULL);
INSERT INTO public.user_matrices (id, endpoint, is_admin, created_at, updated_at, deleted_at) VALUES ('cfd59aba-0616-4c5b-b071-2acf3b9d7e49', '/api/v1/visits', true, '2023-04-01 00:42:06.857335+00', NULL, NULL);
INSERT INTO public.user_matrices (id, endpoint, is_admin, created_at, updated_at, deleted_at) VALUES ('ccef5fbc-5843-46de-968a-9bbefb3f8aa9', '/api/v1/visits', false, '2023-04-01 00:42:06.857335+00', NULL, NULL);
INSERT INTO public.user_matrices (id, endpoint, is_admin, created_at, updated_at, deleted_at) VALUES ('e68e536f-d5dc-4d30-9d3e-c18e27009486', '/api/v1/visits/list', true, '2023-04-01 00:42:06.857335+00', NULL, NULL);
INSERT INTO public.user_matrices (id, endpoint, is_admin, created_at, updated_at, deleted_at) VALUES ('7ceb3854-d522-4a69-accb-8d53c1180e3c', '/api/v1/visits/list', false, '2023-04-01 00:42:06.857335+00', NULL, NULL);
INSERT INTO public.user_matrices (id, endpoint, is_admin, created_at, updated_at, deleted_at) VALUES ('a7aaf8b6-7632-43ab-8e87-3a8468bd8da6', '/api/v1/visits/[id]', true, '2023-04-01 00:42:06.857335+00', NULL, NULL);
INSERT INTO public.user_matrices (id, endpoint, is_admin, created_at, updated_at, deleted_at) VALUES ('b16ee009-5837-49ba-93f7-5124c5f41693', '/api/v1/visits/[id]', false, '2023-04-01 00:42:06.857335+00', NULL, NULL);
INSERT INTO public.user_matrices (id, endpoint, is_admin, created_at, updated_at, deleted_at) VALUES ('603eefe5-f00e-40cd-8e5b-2870e48e66d9', '/api/v1/visits/[id]/delete', true, '2023-04-01 00:42:06.857335+00', NULL, NULL);
INSERT INTO public.user_matrices (id, endpoint, is_admin, created_at, updated_at, deleted_at) VALUES ('c87d72cf-6220-4e56-8900-9518421ab145', '/api/v1/threads', true, '2023-04-01 00:42:06.857335+00', NULL, NULL);
INSERT INTO public.user_matrices (id, endpoint, is_admin, created_at, updated_at, deleted_at) VALUES ('60451b6d-c856-4ac0-b622-9ff7d95192f9', '/api/v1/threads', false, '2023-04-01 00:42:06.857335+00', NULL, NULL);
INSERT INTO public.user_matrices (id, endpoint, is_admin, created_at, updated_at, deleted_at) VALUES ('0aa777b4-3acf-4194-8a43-dd89f920b29a', '/api/v1/threads/list', true, '2023-04-01 00:42:06.857335+00', NULL, NULL);
INSERT INTO public.user_matrices (id, endpoint, is_admin, created_at, updated_at, deleted_at) VALUES ('a4537dcc-6c47-4f61-8a4d-710f79c89ed7', '/api/v1/threads/list', false, '2023-04-01 00:42:06.857335+00', NULL, NULL);
INSERT INTO public.user_matrices (id, endpoint, is_admin, created_at, updated_at, deleted_at) VALUES ('ea1cddaa-3578-4875-aee9-b7d57faa8175', '/api/v1/threads/[id]', true, '2023-04-01 00:42:06.857335+00', NULL, NULL);
INSERT INTO public.user_matrices (id, endpoint, is_admin, created_at, updated_at, deleted_at) VALUES ('61ddd01d-898a-40dd-a1a1-8e407d89fa84', '/api/v1/threads/[id]', false, '2023-04-01 00:42:06.857335+00', NULL, NULL);
INSERT INTO public.user_matrices (id, endpoint, is_admin, created_at, updated_at, deleted_at) VALUES ('8ceb2120-d5d0-430d-bb9f-24624872b12f', '/api/v1/threads/[id]/delete', true, '2023-04-01 00:42:06.857335+00', NULL, NULL);


--
-- TOC entry 3363 (class 0 OID 49241)
-- Dependencies: 217
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.users (id, full_name, email, created_at, updated_at, deleted_at, password, is_admin) VALUES ('355b4c07-bebc-4789-a8a0-f3f6b9c3b6d3', 'Admin', 'admin@admin.com', '2023-04-01 05:13:12.170498+00', '2023-04-01 05:13:12.170498+00', '0001-01-01 00:00:00+00', '2bb80d537b1da3e38bd30361aa855686bde0eacd7162fef6a25fe97bf527a25b', true);
INSERT INTO public.users (id, full_name, email, created_at, updated_at, deleted_at, password, is_admin) VALUES ('36182264-f4a4-4a89-8ec8-589bc6f0e519', 'John Doe', 'johndoe@mail.com', '2023-04-01 05:13:50.811607+00', '2023-04-01 05:13:50.811607+00', '0001-01-01 00:00:00+00', '2bb80d537b1da3e38bd30361aa855686bde0eacd7162fef6a25fe97bf527a25b', false);


--
-- TOC entry 3364 (class 0 OID 49249)
-- Dependencies: 218
-- Data for Name: visits; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.visits (id, user_id, created_at, updated_at, deleted_at) VALUES ('33dd9076-783a-4539-af61-96afd28bad46', '36182264-f4a4-4a89-8ec8-589bc6f0e519', '2023-04-01 05:19:20.645207+00', '2023-04-01 05:19:20.645207+00', '0001-01-01 00:00:00+00');


--
-- TOC entry 3213 (class 2606 OID 49255)
-- Name: users guests_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT guests_pkey PRIMARY KEY (id);


--
-- TOC entry 3209 (class 2606 OID 49257)
-- Name: threads threads_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.threads
    ADD CONSTRAINT threads_pkey PRIMARY KEY (id);


--
-- TOC entry 3211 (class 2606 OID 49259)
-- Name: user_matrices user_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_matrices
    ADD CONSTRAINT user_roles_pkey PRIMARY KEY (id);


--
-- TOC entry 3215 (class 2606 OID 49261)
-- Name: visits visits_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.visits
    ADD CONSTRAINT visits_pkey PRIMARY KEY (id);


--
-- TOC entry 3216 (class 2606 OID 49262)
-- Name: threads threads_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.threads
    ADD CONSTRAINT threads_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- TOC entry 3217 (class 2606 OID 49267)
-- Name: threads threads_visit_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.threads
    ADD CONSTRAINT threads_visit_id_fkey FOREIGN KEY (visit_id) REFERENCES public.visits(id);


--
-- TOC entry 3218 (class 2606 OID 49272)
-- Name: visits visits_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.visits
    ADD CONSTRAINT visits_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


-- Completed on 2023-04-01 05:29:04 UTC

--
-- PostgreSQL database dump complete
--

