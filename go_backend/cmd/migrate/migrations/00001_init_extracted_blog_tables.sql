-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS blogs (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(64) NOT NULL,
    description TEXT NOT NULL,
    publication_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    type TEXT CHECK (type IN ('extracted', 'user')) NOT NULL
);

CREATE TABLE IF NOT EXISTS companies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS authors (
    id SERIAL PRIMARY KEY,
    company_id SERIAL NOT NULL,
    name VARCHAR(255) NOT NULL,

    CONSTRAINT fk_authors_company_id_companies
        FOREIGN KEY (company_id)
        REFERENCES companies (id)
);

CREATE TABLE IF NOT EXISTS extracted_blog_authors (
    blog_id uuid NOT NULL,
    author_id SERIAL NOT NULL,

    CONSTRAINT fk_extracted_blog_authors_blog_id_blogs
        FOREIGN KEY (blog_id)
        REFERENCES blogs(id),
    
    CONSTRAINT fk_extracted_blog_authors_author_id_authors
        FOREIGN KEY (author_id)
        REFERENCES authors(id) 
);

CREATE TABLE IF NOT EXISTS extracted_blogs (
    blog_id uuid,
    company_id SERIAL,
    content TEXT,
    link TEXT,

    CONSTRAINT fk_extracted_blogs_blog_id_blogs
        FOREIGN KEY (blog_id)
        REFERENCES blogs(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_extracted_blogs_company_id_companies
        FOREIGN KEY (company_id)
        REFERENCES companies (id)
);

CREATE TABLE IF NOT EXISTS tags (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(64)
);

CREATE TABLE IF NOT EXISTS blog_tags (
    tag_id uuid NOT NULL,
    blog_id uuid NOT NULL,

    CONSTRAINT fk_blog_tags_tag_id_tags
        FOREIGN KEY (tag_id)
        REFERENCES tags(id)
        ON DELETE CASCADE,
    
    CONSTRAINT fk_blog_tags_blog_id_blogs
        FOREIGN KEY (blog_id)
        REFERENCES blogs(id)
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS blogs CASCADE;
DROP TABLE IF EXISTS companies CASCADE;
DROP TABLE IF EXISTS authors CASCADE;
DROP TABLE IF EXISTS extracted_blog_authors CASCADE;
DROP TABLE IF EXISTS extracted_blogs CASCADE;
DROP TABLE IF EXISTS tags CASCADE;
DROP TABLE IF EXISTS blog_tags CASCADE;
-- +goose StatementEnd
