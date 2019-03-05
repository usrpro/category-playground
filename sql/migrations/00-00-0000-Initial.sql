create table users (
        id serial primary key,
        email text unique not null,
        fname text,
        lname text,
        phone text unique,
        password text not null
);

create table images (
        id bigserial primary key,
        user_id int not null references users (id),
        file text unique not null
);

create table categories (
        id serial primary key,
        parent int not null,
        name text not null
);

create index ix_categories_parent on categories (parent);

create table articles (
        id bigserial primary key,
        title text,
        body text,
        user_id int references users (id)
);

create table article_images (
        article_id bigint references articles (id),
        image_id bigint references images(id),
        primary key (article_id, image_id)
);

create table article_categories (
        category_id int references categories (id),
        article_id int,
        primary key (category_id, article_id)
);

create index ix_article_category on article_categories (article_id, category_id);




