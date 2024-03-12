

-- Добавление актера
INSERT INTO actors (name, surname, patronymic, birthday, sex, information)
VALUES ($1, $2, $3, $4, $5, $6);

-- изменение информации актера
UPDATE actors SET name=$1, surname=$2, patronymic=$3, birthday=$4, sex=$5, information=$6
WHERE id=$7;

-- удаление информациии актера
DELETE from actors WHERE id=$1;




CREATE TABLE films(
                      id SERIAL PRIMARY KEY,
                      title varchar(150) not null,
                      year INT not null,
                      information varchar(1000),
                      rating DECIMAL(3,1) CHECK (rating BETWEEN 0 AND 10)
);
-- добавление информации о фильме

INSERT INTO films (title, year, information, rating)
VALUES ($1, $2, $3, $4);

-- удаление информации о фильме

-- получение списка фильмов с возможностью сортировки по названию, рейтингу, по дате выпуска

-- получение списка актёров, для каждого актера выдаётся также список ффильмов с его участием


-- добавление актеров к фильмам
