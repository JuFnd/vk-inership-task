DROP TABLE IF EXISTS film_actor;
CREATE TABLE film_actor (
                            id SERIAL PRIMARY KEY,
                            film_id INTEGER,
                            actor_id INTEGER
);

DROP TABLE IF EXISTS film;
CREATE TABLE film (
                      id SERIAL PRIMARY KEY,
                      name TEXT,
                      description TEXT,
                      rating float,
                      releaseDate DATE
);

DROP TABLE IF EXISTS actor;
CREATE TABLE actor (
                       id SERIAL PRIMARY KEY,
                       name TEXT,
                       gender TEXT,
                       birthdate DATE
);

ALTER TABLE film_actor
    ADD FOREIGN KEY (film_id) REFERENCES film (id),
    ADD FOREIGN KEY (actor_id) REFERENCES actor (id);

INSERT INTO film (name, description, rating, releaseDate) VALUES
                                                              ('Film 1', 'Description 1', 4.5, '2022-01-01'),
                                                              ('Film 2', 'Description 2', 3.8, '2022-02-15'),
                                                              ('Film 3', 'Description 3', 4.2, '2022-03-10'),
                                                              ('Film 4', 'Description 4', 3.5, '2022-04-20'),
                                                              ('Film 5', 'Description 5', 4.0, '2022-05-05'),
                                                              ('Film 6', 'Description 6', 4.7, '2022-06-30'),
                                                              ('Film 7', 'Description 7', 3.9, '2022-07-15'),
                                                              ('Film 8', 'Description 8', 4.1, '2022-08-25'),
                                                              ('Film 9', 'Description 9', 3.6, '2022-09-10'),
                                                              ('Film 10', 'Description 10', 4.3, '2022-10-31');

INSERT INTO actor (name, gender, birthdate) VALUES
                                                ('Actor 1', 'Male', '1990-01-01'),
                                                ('Actor 2', 'Female', '1992-02-15'),
                                                ('Actor 3', 'Male', '1985-03-10'),
                                                ('Actor 4', 'Female', '1988-04-20'),
                                                ('Actor 5', 'Male', '1995-05-05'),
                                                ('Actor 6', 'Female', '1993-06-30'),
                                                ('Actor 7', 'Male', '1991-07-15'),
                                                ('Actor 8', 'Female', '1987-08-25'),
                                                ('Actor 9', 'Male', '1994-09-10'),
                                                ('Actor 10', 'Female', '1989-10-31');

INSERT INTO film_actor (film_id, actor_id) VALUES
                                               (1, 1),
                                               (1, 2),
                                               (2, 3),
                                               (2, 8),
                                               (2, 4),
                                               (3, 5),
                                               (3, 6),
                                               (4, 7),
                                               (4, 8),
                                               (4, 6),
                                               (5, 9),
                                               (5, 1),
                                               (5, 10);
