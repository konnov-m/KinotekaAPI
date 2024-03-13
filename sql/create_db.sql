DROP TABLE IF EXISTS users_roles;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS films_actors;
DROP TABLE IF EXISTS actors;
DROP TABLE IF EXISTS films;

CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    login varchar(256) not null,
    password varchar(1024) not null
);

CREATE TABLE roles(
    id SERIAL PRIMARY KEY,
    name varchar(256) not null
);

CREATE TABLE users_roles(
    user_id INTEGER NOT NULL REFERENCES users(id),
    permission_id INTEGER NOT NULL REFERENCES roles(id),
    PRIMARY KEY(user_id, permission_id)
);


CREATE TABLE actors(
    id SERIAL PRIMARY KEY,
    name varchar(256) not null,
    surname varchar(256) not null,
    patronymic varchar(256),
    birthday DATE not null,
    sex CHAR(1) not null,
    information varchar(2048)
);

CREATE TABLE films(
    id SERIAL PRIMARY KEY,
    title varchar(150) not null,
    year INT not null,
    information varchar(1000),
    rating DECIMAL(3,1) CHECK (rating BETWEEN 0 AND 10)
);


CREATE TABLE films_actors(
    film_id INTEGER NOT NULL REFERENCES films(id),
    actor_id INTEGER NOT NULL REFERENCES actors(id),
    PRIMARY KEY(film_id, actor_id)
);





