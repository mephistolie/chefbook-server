CREATE TABLE users
(
    user_id      SERIAL PRIMARY KEY       NOT NULL UNIQUE,
    email        VARCHAR(255)             NOT NULL UNIQUE,
    username     VARCHAR(255),
    password     bytea                    NOT NULL,

    registered   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT timezone('utc', now()),
    is_activated BOOLEAN                  NOT NULL DEFAULT false,

    avatar       VARCHAR(255)                      DEFAULT NULL,

    premium      TIMESTAMP WITH TIME ZONE          DEFAULT NULL,
    broccoins    INT                      NOT NULL DEFAULT 0,

    is_blocked   BOOLEAN                  NOT NULL DEFAULT false,
    key          VARCHAR(255)                      DEFAULT NULL
);

CREATE TABLE activation_links
(
    id              SERIAL PRIMARY KEY                               NOT NULL UNIQUE,
    activation_link uuid                                             NOT NULL,
    user_id         INT REFERENCES users (user_id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE shopping_list
(
    user_id       INT REFERENCES users (user_id) ON DELETE CASCADE NOT NULL UNIQUE,
    shopping_list JSONB                                            NOT NULL DEFAULT '[]'::jsonb
);

CREATE TABLE categories
(
    category_id SERIAL PRIMARY KEY                               NOT NULL UNIQUE,
    name        VARCHAR(255)                                     NOT NULL,
    cover       VARCHAR(20) DEFAULT '',
    user_id     INT REFERENCES users (user_id) ON DELETE CASCADE NOT NULL
);

CREATE TYPE visibility_type as ENUM ('private', 'shared', 'public');

CREATE TABLE recipes
(
    recipe_id          SERIAL PRIMARY KEY                               NOT NULL UNIQUE,
    name               VARCHAR(255)                                     NOT NULL,
    owner_id           INT REFERENCES users (user_id) ON DELETE CASCADE NOT NULL,

    visibility         visibility_type                                  NOT NULL DEFAULT 'private',
    language           VARCHAR(2)                                       NOT NULL DEFAULT 'en',
    description        TEXT                                                      DEFAULT NULL,
    likes              INT                                              NOT NULL DEFAULT 0,

    servings           SMALLINT                                                  DEFAULT NULL,
    time               SMALLINT                                                  DEFAULT NULL,
    calories           SMALLINT                                                  DEFAULT NULL,
    protein            SMALLINT                                                  DEFAULT NULL,
    fats               SMALLINT                                                  DEFAULT NULL,
    carbohydrates      SMALLINT                                                  DEFAULT NULL,

    ingredients        JSONB                                            NOT NULL,
    cooking            JSONB                                            NOT NULL,

    preview            VARCHAR(255)                                              DEFAULT NULL,
    encrypted          BOOLEAN                                          NOT NULL DEFAULT false,
    key                VARCHAR(255)                                              DEFAULT NULL,
    creation_timestamp TIMESTAMP WITH TIME ZONE                         NOT NULL DEFAULT timezone('utc', now()),
    update_timestamp   TIMESTAMP WITH TIME ZONE                         NOT NULL DEFAULT timezone('utc', now())
);

CREATE TABLE users_recipes
(
    user_id    INT REFERENCES users (user_id) ON DELETE CASCADE     NOT NULL,
    recipe_id  INT REFERENCES recipes (recipe_id) ON DELETE CASCADE NOT NULL,
    favourite  BOOLEAN                                              NOT NULL DEFAULT false,
    user_key   TEXT                                                          DEFAULT NULL,
    recipe_key TEXT                                                          DEFAULT NULL
);

CREATE TABLE encrypted_recipes_requests
(
    user_id              INT REFERENCES users (user_id) ON DELETE CASCADE     NOT NULL,
    recipe_id            INT REFERENCES recipes (recipe_id) ON DELETE CASCADE NOT NULL,
    encrypted_user_key   TEXT                                                 NOT NULL,
    encrypted_recipe_key TEXT DEFAULT NULL
);

CREATE TABLE likes
(
    user_id   INT REFERENCES users (user_id) ON DELETE CASCADE     NOT NULL,
    recipe_id INT REFERENCES recipes (recipe_id) ON DELETE CASCADE NOT NULL
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
    expires_at    TIMESTAMP WITH TIME ZONE                         NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE                         NOT NULL DEFAULT timezone('utc', now())
);

CREATE TYPE role as ENUM ('credentials', 'admin');

CREATE TABLE roles
(
    role_id SERIAL PRIMARY KEY                               NOT NULL UNIQUE,
    name    role                                             NOT NULL DEFAULT 'credentials',
    user_id INT REFERENCES users (user_id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE news
(
    id    SERIAL PRIMARY KEY NOT NULL UNIQUE,
    name  VARCHAR(255)       NOT NULL,
    text  TEXT               NOT NULL,
    cover VARCHAR(255) DEFAULT NULL
);