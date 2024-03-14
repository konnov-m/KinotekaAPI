

-- Добавление актера
INSERT INTO actors (name, surname, patronymic, birthday, sex, information)
VALUES ($1, $2, $3, $4, $5, $6);

-- изменение информации актера
UPDATE actors SET name=$1, surname=$2, patronymic=$3, birthday=$4, sex=$5, information=$6
WHERE id=$7;

-- удаление информациии актера
DELETE from actors WHERE id=$1;

-- добавление информации о фильме
INSERT INTO films (title, year, information, rating)
VALUES ($1, $2, $3, $4);

-- изменение информации о фильме
UPDATE films SET title=$1, year=$2, information=$3, rating=$4
WHERE id=$5;

-- удаление информации о фильме
DELETE from films WHERE id=$1;

-- получение списка фильмов с возможностью сортировки по названию, рейтингу, по дате выпуска
SELECT id, title, year, information, rating
FROM films
ORDER BY $1 DESC;

-- получение списка актёров, для каждого актера выдаётся также список фильмов с его участием
-- сомнительная реализация
SELECT
    a.id AS actor_id,
    a.name AS actor_name,
    a.surname AS actor_surname,
    a.patronymic AS actor_patronymic,
    a.birthday AS actor_birthday,
    a.sex AS actor_sex,
    a.information AS actor_information,
    f.id AS film_id,
    f.title AS film_title,
    f.year AS film_year,
    f.information AS film_information,
    f.rating AS film_rating
FROM
    actors a
        JOIN
    films_actors fa ON a.id = fa.actor_id
        JOIN
    films f ON fa.film_id = f.id;

--получение фильма по фрагменту названия
SELECT * FROM films WHERE LOWER(title) LIKE '%'+ $1 + '%';

-- добавление актеров к фильмам
INSERT INTO films_actors (film_id, actor_id)
VALUES ($1, $2);