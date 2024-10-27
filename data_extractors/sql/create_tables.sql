BEGIN;


CREATE TABLE IF NOT EXISTS public.company_blog_site
(
    id uuid NOT NULL,
    blog_name text NOT NULL,
    feed_link text NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.company_blog_posts
(
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 ),
    company_id uuid NOT NULL,
    title text NOT NULL,
    description text,
    publication_date timestamp with time zone,
    link text NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.company_blog_post_content
(
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 ),
    company_blog_post_id integer NOT NULL,
    tag_type text NOT NULL,
    tag_content text NOT NULL,
    "order" integer NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.tags
(
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 ),
    name text NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.blog_tags
(
    blog_post_id integer NOT NULL,
    tag_id integer NOT NULL,
    PRIMARY KEY (blog_post_id, tag_id)
);

CREATE TABLE IF NOT EXISTS public.authors
(
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 ),
    name text NOT NULL,
    company_id uuid NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.blog_post_authors
(
    blog_id integer NOT NULL,
    author_id integer NOT NULL,
    PRIMARY KEY (blog_id, author_id)
);

CREATE TABLE IF NOT EXISTS public.accounts
(
    id uuid NOT NULL,
    username text NOT NULL,
    password text NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.followings
(
    id uuid NOT NULL,
    account_id uuid NOT NULL,
    following_type text NOT NULL,
    following_id uuid NOT NULL,
    PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public.company_blog_posts
    ADD CONSTRAINT "FK_company_id" FOREIGN KEY (company_id)
    REFERENCES public.company_blog_site (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS public.company_blog_post_content
    ADD CONSTRAINT "FK_company_blog_post_id" FOREIGN KEY (company_blog_post_id)
    REFERENCES public.company_blog_posts (id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE CASCADE
    NOT VALID;


ALTER TABLE IF EXISTS public.blog_tags
    ADD CONSTRAINT "FK_company_blog_posts_id" FOREIGN KEY (blog_post_id)
    REFERENCES public.company_blog_posts (id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE CASCADE
    NOT VALID;


ALTER TABLE IF EXISTS public.blog_tags
    ADD CONSTRAINT "FK_tags_id" FOREIGN KEY (tag_id)
    REFERENCES public.tags (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS public.authors
    ADD CONSTRAINT "FK_company_blog_site_id" FOREIGN KEY (company_id)
    REFERENCES public.company_blog_site (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS public.blog_post_authors
    ADD CONSTRAINT "FK_authors_id" FOREIGN KEY (author_id)
    REFERENCES public.authors (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS public.blog_post_authors
    ADD CONSTRAINT "FK_company_blog_post_id" FOREIGN KEY (blog_id)
    REFERENCES public.company_blog_posts (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS public.followings
    ADD CONSTRAINT "FK_accounts_id" FOREIGN KEY (account_id)
    REFERENCES public.accounts (id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE CASCADE
    NOT VALID;


ALTER TABLE IF EXISTS public.followings
    ADD CONSTRAINT "FK_accounts_following_id" FOREIGN KEY (following_id)
    REFERENCES public.accounts (id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE CASCADE
    NOT VALID;


ALTER TABLE IF EXISTS public.followings
    ADD CONSTRAINT "FK_company_blog_site_id" FOREIGN KEY (following_id)
    REFERENCES public.company_blog_site (id) MATCH SIMPLE
    ON UPDATE CASCADE
    ON DELETE CASCADE
    NOT VALID;

END;