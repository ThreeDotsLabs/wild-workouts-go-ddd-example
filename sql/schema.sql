CREATE TABLE `hours`
(
    hour         TIMESTAMP                                                 NOT NULL,
    availability ENUM ('available', 'not_available', 'training_scheduled') NOT NULL,
    PRIMARY KEY (hour)
);