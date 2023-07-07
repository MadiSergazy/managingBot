-- USE telegrambot;

CREATE TABLE if not exists admins (
                                      id INT PRIMARY KEY AUTO_INCREMENT,
                                      phone_number VARCHAR(20) NOT NULL UNIQUE
    );

CREATE TABLE if not exists workers (
                                       id INT PRIMARY KEY AUTO_INCREMENT,
                                       phone_number VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL
    );

CREATE TABLE  if not exists tasks (
                                      id INT PRIMARY KEY AUTO_INCREMENT,
                                      nameOfTask VARCHAR(255) NOT NULL,
    dateStart DATE,
    dateEnd DATE,
    isDone BOOLEAN,
    worker_id INT,
    FOREIGN KEY (worker_id) REFERENCES workers(id)
    );

CREATE TABLE if not exists lifts (
                                     id INT PRIMARY KEY AUTO_INCREMENT,
                                     name_lift VARCHAR(255) NOT NULL,
    task_id INT,
    worker_id INT,
    FOREIGN KEY (task_id) REFERENCES tasks(id),
    FOREIGN KEY (worker_id) REFERENCES workers(id)
    );

CREATE TABLE if not exists projects (
                                        id INT PRIMARY KEY AUTO_INCREMENT,
                                        name_resident VARCHAR(255) NOT NULL,
    name_lift VARCHAR(255) NOT NULL,
    worker_id INT,
    lift_id INT,
    FOREIGN KEY (worker_id) REFERENCES workers(id),
    FOREIGN KEY (lift_id) REFERENCES lifts(id)
    );