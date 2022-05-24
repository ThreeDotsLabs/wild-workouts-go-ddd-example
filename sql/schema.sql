CREATE TABLE `hours`
(
    hour         TIMESTAMP                                                 NOT NULL DEFAULT 0,
    availability ENUM ('available', 'not_available', 'training_scheduled') NOT NULL,
    PRIMARY KEY (hour)
);