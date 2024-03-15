

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





INSERT INTO roles(name) VALUES ('admin'), ('user');


INSERT INTO actors (id, name, surname, patronymic, birthday, sex, information) VALUES
(1, 'Райан', 'Томас Гослинг', '', DATE('12/11/1980'), 'm', 'Томас Гослинг Райан'),
(2, 'Киллиан', 'Мерфи', '', DATE('05/25/1976'), 'm', null),
(3, 'Питер', 'Хейден Динклэйдж', '', DATE('11/06/1969'), 'm', null),
(4, 'Кристиан', 'Бэйл', '', DATE('01/30/1974'), 'm', null),
(5,'Леонардо', 'ДиКаприо', '', DATE('11/11/1974'), 'm', null),
(6,'Брэд', ' Питт', '', DATE('12/18/1963'), 'm', null),
(7,'Эдвард', ' Нортон', '', DATE('08/18/1969'), 'm', null),
(8,'Рами', 'Малек', '', DATE('12/05/1981'), 'm', null),
(9,'Сэмюэл', 'Л. Джексон','', DATE('12/21/1948'),'m', null),
(10,'Брайан', 'Крэнстон','', DATE('7/03/1956'), 'm', null),
(11, 'Мэтью', 'Макконахи', '', DATE('4/11/1969'), 'm', null),
(12, 'Сергей', 'Бодров', '', DATE('12/27/1971'), 'm', null),
(13, 'Тоби',  'Магуайр', '', DATE('06/27/1975'), 'm', null);


INSERT INTO films (id, title, year, information, rating) VALUES
(1, 'Оппенгеймер', 2023, '03:00', 9.0),
(2, 'Бойцовский клуб', 1999, '02:19', 9.1),
(3, 'Начало', 2010, '02:19', 8.9),
(4, 'Криминальное чтиво', 1994, '02:34', 8.5),
(5, 'Джентльмены', 2019 , '1:53', 8.6),
(6, 'Брат', 1997, '1:40', 8.6),
(7, 'Брат 2', 2000,  '2:07', 8.6),
(8, 'Человек-паук', 2002, '2:01', 8.3),
(9, 'Человек-паук 2', 2004, '2:07', 8.2),
(10, 'Человек-паук 3', 2007, '2:19', 8.1),
(11, 'Ходячий замок', 2004, '1:59', 8.1),
(12, 'Бесславные ублюдки', 2009, '2:34', 8.9),
(13, 'Однажды в... Голливуде', 2019, '2:41', 8.7);

INSERT INTO films_actors (film_id, actor_id) VALUES
(1, 2),
(2, 6),
(3, 5),
(4 ,9),
(5, 11),
(6, 12),
(7, 12),
(8, 13),
(9, 13),
(10, 13),
(12, 6),
(13, 6),
(13, 5);