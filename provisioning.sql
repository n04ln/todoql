CREATE TABLE todos (
    id varchar(128) UNIQUE,
    text varchar(128),
    done boolean,
    user_id varchar(128) -- referenced users.id
);

CREATE TABLE users (
    id varchar(128) UNIQUE,
    name varchar(128)
);

-- predata
INSERT INTO users(id, name) VALUES ('qwer', 'noah');
INSERT INTO todos(id, text, done, user_id) VALUES ('asdf', 'check it out', false, 'qwer');
INSERT INTO todos(id, text, done, user_id) VALUES ('zxcv', 'disscuss gqlgen', true, 'qwer');
