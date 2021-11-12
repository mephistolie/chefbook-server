CREATE TABLE users
(
    user_id         SERIAL PRIMARY KEY NOT NULL UNIQUE,
    email           VARCHAR(255)       NOT NULL UNIQUE,
    username        VARCHAR(255),
    password        bytea              NOT NULL,

    is_activated    BOOLEAN            NOT NULL DEFAULT false,
    activation_link uuid               NOT NULL,

    avatar          VARCHAR(255) UNIQUE         DEFAULT NULL,
    vk_id           VARCHAR(255) UNIQUE         DEFAULT NULL,
    premium         TIMESTAMP                   DEFAULT null,

    is_blocked      BOOLEAN            NOT NULL DEFAULT false
);

CREATE TABLE shopping_list
(
    item_id SERIAL PRIMARY KEY             NOT NULL UNIQUE,
    name    VARCHAR(255)                   NOT NULL,
    user_id INT REFERENCES users (user_id) NOT NULL
);

CREATE TYPE category_type as ENUM ('red', 'blue');

CREATE TABLE categories
(
    category_id SERIAL PRIMARY KEY             NOT NULL UNIQUE,
    name        VARCHAR(255)                   NOT NULL,
    type        category_type DEFAULT NULL,
    user_id     INT REFERENCES users (user_id) NOT NULL
);

CREATE TYPE visibility_type as ENUM ('private', 'shared', 'public');

CREATE TABLE recipes
(
    recipe_id          SERIAL PRIMARY KEY                               NOT NULL UNIQUE,
    name               VARCHAR(255)                                     NOT NULL,
    owner_id           INT REFERENCES users (user_id) ON DELETE CASCADE NOT NULL,

    servings           SMALLINT                                         NOT NULL DEFAULT 1,
    time               INT                                              NOT NULL,
    calories           SMALLINT                                                  DEFAULT NULL,

    ingredients        JSONB                                            NOT NULL,
    cooking            JSONB                                            NOT NULL,
    preview            VARCHAR(255),
    visibility         visibility_type                                  NOT NULL DEFAULT 'private',
    encrypted          BOOLEAN                                          NOT NULL DEFAULT false,
    creation_timestamp TIMESTAMP                                        NOT NULL DEFAULT now(),
    update_timestamp   TIMESTAMP                                        NOT NULL DEFAULT now()
);

CREATE TABLE users_recipes
(
    user_id   INT REFERENCES users (user_id) ON DELETE CASCADE     NOT NULL,
    recipe_id INT REFERENCES recipes (recipe_id) ON DELETE CASCADE NOT NULL,
    favourite BOOLEAN                                              NOT NULL DEFAULT false,
    liked     BOOLEAN                                              NOT NULL DEFAULT false
);

CREATE TABLE recipes_categories
(
    recipe_id   INT REFERENCES recipes (recipe_id) ON DELETE CASCADE      NOT NULL,
    category_id INT REFERENCES categories (category_id) ON DELETE CASCADE NOT NULL,
    user_id     INT REFERENCES users (user_id) ON DELETE CASCADE          NOT NULL
);

CREATE TABLE sessions
(
    session_id    SERIAL PRIMARY KEY                               NOT NULL UNIQUE,
    user_id       INT REFERENCES users (user_id) ON DELETE CASCADE NOT NULL,
    refresh_token VARCHAR(255)                                     NOT NULL UNIQUE,
    ip            VARCHAR(255)                                     NOT NULL,
    expires_at    TIMESTAMP                                        NOT NULL,
    created_at    TIMESTAMP                                        NOT NULL DEFAULT now()
);

CREATE TYPE role as ENUM ('user', 'admin');

CREATE TABLE roles
(
    role_id SERIAL PRIMARY KEY                               NOT NULL UNIQUE,
    name    role                                             NOT NULL DEFAULT 'user',
    user_id INT REFERENCES users (user_id) ON DELETE CASCADE NOT NULL
)