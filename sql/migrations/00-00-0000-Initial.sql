CREATE TABLE categories (
       id INT NOT NULL,
       parent INT NOT NULL,
       lft INT NOT NULL,
       rgt INT NOT NULL,
       data VARCHAR(100) NOT NULL,
       stuffing VARCHAR(100) NOT NULL
);

ALTER TABLE categories ADD CONSTRAINT pk_categories_id PRIMARY KEY (id);
CREATE INDEX ix_categories_lft ON categories (lft);
CREATE INDEX ix_categories_rgt ON categories (rgt);
CREATE INDEX ix_categories_parent ON categories (parent);