-- USE telegrambot;
/*


CREATE TABLE if not exists  `admins` (  `id` int NOT NULL AUTO_INCREMENT, 
 `phone_number` varchar(20) NOT NULL,  `identifier` int, 
 PRIMARY KEY (`id`),  UNIQUE KEY `phone_number` (`phone_number`)) 
 ENGINE=InnoDB  AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE if not exists  `workers`  (  `id` int NOT NULL AUTO_INCREMENT, 
 `phone_number` varchar(20) DEFAULT NULL,  `name` varchar(255) NOT NULL, `identifier` int, 
 PRIMARY KEY (`id`),  UNIQUE KEY `phone_number` (`phone_number`))
 ENGINE=InnoDB  AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE if not exists  `tasks` (  `id` int NOT NULL AUTO_INCREMENT, 
 `nameOfTask` varchar(255) NOT NULL,  `dateStart` date, 
 `dateEnd` date,  
 `isDone` BOOL  DEFAULT FALSE, 
 `worker_id` int DEFAULT NULL,  PRIMARY KEY (`id`),  
 KEY `worker_id` (`worker_id`),  CONSTRAINT `tasks_ibfk_1` 
 FOREIGN KEY (`worker_id`) REFERENCES `workers` (`id`))
 ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

 CREATE TABLE  if not exists `lifts` (  `id` int NOT NULL AUTO_INCREMENT, 
 `name_lift` varchar(255) NOT NULL,  `worker_id` int DEFAULT NULL,
 PRIMARY KEY (`id`),  KEY `lifts_ibfk_1` (`worker_id`), 
 CONSTRAINT `lifts_ibfk_1` FOREIGN KEY (`worker_id`)
 REFERENCES `workers` (`id`)) ENGINE=InnoDB  
 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `projects` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name_resident` VARCHAR(255) NOT NULL,
  `name_lift` VARCHAR(255) NOT NULL,
  `worker_id` INT DEFAULT NULL,
  `lift_id` INT DEFAULT NULL,
  `lift_details_id` INT DEFAULT NULL,
  `task_id` INT DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `worker_id` (`worker_id`),
  KEY `lift_id` (`lift_id`),
  KEY `lift_details_id` (`lift_details_id`), -- Add this line
   KEY `task_id` (`task_id`),
  CONSTRAINT `projects_ibfk_1` FOREIGN KEY (`worker_id`) REFERENCES `workers` (`id`),
  CONSTRAINT `projects_ibfk_2` FOREIGN KEY (`lift_id`) REFERENCES `lifts` (`id`),
  CONSTRAINT `projects_ibfk_3` FOREIGN KEY (`lift_details_id`) REFERENCES `lift_details` (`id`) -- Add this line
  CONSTRAINT `projects_ibfk_4` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

 

CREATE TABLE if not exists  `lift_tasks` (  `lift_id` int NOT NULL,  `task_id` int NOT NULL,
  PRIMARY KEY (`lift_id`,`task_id`),  KEY `lift_tasks_ibfk_2` (`task_id`),
  CONSTRAINT `lift_tasks_ibfk_1` FOREIGN KEY (`lift_id`) REFERENCES `lifts`
  (`id`),  CONSTRAINT `lift_tasks_ibfk_2` FOREIGN KEY (`task_id`)
  REFERENCES `tasks` (`id`)) ENGINE=InnoDB  DEFAULT
  CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE if not exists  `task_of_lifts` (  `id` int NOT NULL AUTO_INCREMENT, 
 `task_name` varchar(255) NOT NULL,  PRIMARY KEY (`id`))
 ENGINE=InnoDB  AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;



  CREATE TABLE IF NOT EXISTS `lift_details` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name_resident` VARCHAR(255) NOT NULL,
  `name_lift` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;



 
 INSERT INTO admins (phone_number) VALUES ('77078566392');
 
INSERT INTO task_of_lifts (id, task_name) VALUES
(1, 'Приемка строительной части лифта'),
(2, 'Приемка оборудования и технической документации для монтажа, замены, модернизации'),
(3, 'Определение координат установки оборудования лифта'),
(4, 'Установка кронштейнов крепления направляющих кабины и противовеса'),
(5, 'Монтаж направляющих кабины и противовеса'),
(6, 'Монтаж дверей шахты'),
(7, 'Монтаж оборудования приямка'),
(8, 'Монтаж противовеса'),
(9, 'Монтаж кабины'),
(10, 'Монтаж лебедки главного привода'),
(11, 'Монтаж ограничителя скорости'),
(12, 'Навеска гибких тяговых элементов'),
(13, 'Монтаж электроаппаратуры, кабелей, электропроводки и цепей заземления'),
(14, 'Пусконаладочные работы');

*/



-- Create the admins table
CREATE TABLE IF NOT EXISTS admins (
  id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
  phone_number VARCHAR(20) NOT NULL,
  identifier INT
);

-- Create the workers table
CREATE TABLE IF NOT EXISTS workers (
  id INT NOT NULL PRIMARY KEY,
  phone_number VARCHAR(20) DEFAULT NULL,
  name VARCHAR(255) NOT NULL,
  identifier INT
);

-- Create the lifts table
CREATE TABLE IF NOT EXISTS lifts (
  id INT NOT NULL PRIMARY KEY,
  name_lift VARCHAR(255) NOT NULL,
  worker_id INT,
  FOREIGN KEY (worker_id) REFERENCES workers (id)
);


-- Create the projects table
CREATE TABLE IF NOT EXISTS projects (
  id INT NOT NULL PRIMARY KEY,
  name_resident VARCHAR(255) NOT NULL,
  worker_id INT,
  lift_id INT,
  FOREIGN KEY (worker_id) REFERENCES workers (id),
  FOREIGN KEY (lift_id) REFERENCES lifts (id)
);

-- Create the tasks table
CREATE TABLE IF NOT EXISTS tasks (
  id INT NOT NULL PRIMARY KEY,
  nameOfTask VARCHAR(255) NOT NULL,
  dateStart DATE,
  dateEnd DATE,
  isDone BOOL DEFAULT FALSE,
  lift_id INT,
  is_validate BOOL DEFAULT FALSE,  --todo add this
  file_id VARCHAR(255),            --todo add this 
  date_requested_to_validate DATE, --todo ad dthis and check after 48hour with is_validate if is_val false and 48 hour send to hr_manager and also send to worker notification 
  is_rejected BOOL DEFAULT FALSE,  -- todo add this if is reject is true send notification to worker with text description
  reject_description TEXT, 
  FOREIGN KEY (lift_id) REFERENCES lifts (id)
);


CREATE TABLE force_majeure (
    id INT PRIMARY KEY AUTO_INCREMENT,
    task_id INT,
    residential_complex VARCHAR(255) NOT NULL,
    elevator_name VARCHAR(255) NOT NULL,
    employee_phone_number VARCHAR(15) NOT NULL,
    incident_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    description TEXT NOT NULL
);

CREATE TABLE `change_requests` 
( `id` int NOT NULL AUTO_INCREMENT,
  `task_id` int DEFAULT NULL,
  `residential_complex` varchar(255) NOT NULL,
  `elevator_name` varchar(255) NOT NULL,
  `employee_phone_number` varchar(15) NOT NULL,
  `employee_identifier` int NOT NULL,
  `incident_time` timestamp DEFAULT CURRENT_TIMESTAMP,
  `description` text NOT NULL,
  `description_of_what_done` text,
  `is_done` bool DEFAULT false,
PRIMARY KEY (`id`)) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE overdue_task (
    id INT PRIMARY KEY AUTO_INCREMENT,
    task_id INT NOT NULL,
    name_resident VARCHAR(255) NOT NULL,
    name_lift VARCHAR(255) NOT NULL,
    phone_number VARCHAR(15) NOT NULL,
    name_of_task VARCHAR(255) NOT NULL,
    date_start DATE NOT NULL,
    date_end DATE NOT NULL,
    is_done_by_worker BOOL NOT NULL, 
    is_done_by_admin BOOL default false, 
    is_done_by_hr_manager BOOL default false,
    description_by_admin text,
    description_by_hr_manager text
);


CREATE TABLE IF NOT EXISTS recommendations (
  id INT NOT NULL AUTO_INCREMENT,
  date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  phone_number VARCHAR(15) NOT NULL,
  description TEXT,
  PRIMARY KEY (id)
);


-- Insert sample adminsta into the tables
INSERT INTO admins (phone_number) VALUES
  ('77078566392');

-- INSERT INTO workers (id, phone_number, name, identifier) VALUES
--   (1, '1111111111', 'John Doe', 100),
--   (2, '2222222222', 'Jane Smith', 200);

-- INSERT INTO lifts (id, name_lift, worker_id) VALUES
--   (1, 'Lift1', 1),
--   (2, 'Lift2', 1);

-- INSERT INTO projects (id, name_resident, worker_id, lift_id) VALUES
--   (1, 'Resident1', 1, 1),
--   (2, 'Resident2', 1, 2);

-- INSERT INTO tasks (id, nameOfTask, dateStart, dateEnd, isDone, lift_id) VALUES
--   (1, 'Task1', '2023-07-03', '2023-07-10', 1, 1),
--   (2, 'Task2', '2023-07-03', '2023-07-10', 0, 1),
--   (3, 'Task3', '2023-07-03', '2023-07-10', 1, 1),
--   (4, 'Task4', '2023-07-03', '2023-07-15', 0, 2),
--   (5, 'Task5', '2023-07-03', '2023-07-10', 1, 2),
--   (6, 'Task6', '2023-07-03', '2023-07-15', 0, 2);


CREATE TABLE if not exists  `task_of_lifts` ( 
 `id` int NOT NULL AUTO_INCREMENT, 
 `task_name` varchar(255) NOT NULL,  PRIMARY KEY (`id`))
 ENGINE=InnoDB  AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


INSERT INTO task_of_lifts (id, task_name) VALUES
(1, 'Приемка строительной части лифта'),
(2, 'Приемка оборудования и технической документации для монтажа, замены, модернизации'),
(3, 'Определение координат установки оборудования лифта'),
(4, 'Установка кронштейнов крепления направляющих кабины и противовеса'),
(5, 'Монтаж направляющих кабины и противовеса'),
(6, 'Монтаж дверей шахты'),
(7, 'Монтаж оборудования приямка'),
(8, 'Монтаж противовеса'),
(9, 'Монтаж кабины'),
(10, 'Монтаж лебедки главного привода'),
(11, 'Монтаж ограничителя скорости'),
(12, 'Навеска гибких тяговых элементов'),
(13, 'Монтаж электроаппаратуры, кабелей, электропроводки и цепей заземления'),
(14, 'Пусконаладочные работы');

-- SELECT
--   t.id,
--   p.name_resident,
--   l.name_lift,
--   w.phone_number,
--   t.nameOfTask,
--   t.dateStart,
--   t.dateEnd,
--   t.isDone
-- FROM projects p
-- JOIN lifts l ON p.lift_id = l.id
-- JOIN workers w ON p.worker_id = w.id
-- JOIN tasks t ON t.lift_id = l.id
-- WHERE w.id = 1
-- ORDER BY p.id, t.id;