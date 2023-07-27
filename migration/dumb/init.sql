SET NAMES utf8mb4;

DROP TABLE IF EXISTS `admins`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admins` (
  `id` int NOT NULL AUTO_INCREMENT,
  `phone_number` varchar(20) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `identifier` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `admins`
--

LOCK TABLES `admins` WRITE;
/*!40000 ALTER TABLE `admins` DISABLE KEYS */;
INSERT INTO `admins` VALUES (1,'77078566392','JurtBala20O21',590541456);
/*!40000 ALTER TABLE `admins` ENABLE KEYS */;
UNLOCK TABLES;



DROP TABLE IF EXISTS `workers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `workers` (
  `id` int NOT NULL AUTO_INCREMENT,
  `phone_number` varchar(20) DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `identifier` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `workers`
--

LOCK TABLES `workers` WRITE;
/*!40000 ALTER TABLE `workers` DISABLE KEYS */;
INSERT INTO `workers` VALUES (1,'77009654486','UzuiNU',2088886948),(2,'77078566392','JurtBala20O21',590541456),(3,'77078566392','JurtBala20O21',590541456),(4,'77078566392','JurtBala20O21',590541456);
/*!40000 ALTER TABLE `workers` ENABLE KEYS */;
UNLOCK TABLES;



DROP TABLE IF EXISTS `hr_manager`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `hr_manager` (
  `id` int NOT NULL AUTO_INCREMENT,
  `phone_number` varchar(20) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `identifier` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `hr_manager`
--

LOCK TABLES `hr_manager` WRITE;
/*!40000 ALTER TABLE `hr_manager` DISABLE KEYS */;
INSERT INTO `hr_manager` VALUES (1,'77078566392','JurtBala20O21',590541456);
/*!40000 ALTER TABLE `hr_manager` ENABLE KEYS */;
UNLOCK TABLES;




DROP TABLE IF EXISTS `force_majeure`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `force_majeure` (
  `id` int NOT NULL AUTO_INCREMENT,
  `task_id` int DEFAULT NULL,
  `residential_complex` varchar(255) NOT NULL,
  `elevator_name` varchar(255) NOT NULL,
  `employee_phone_number` varchar(15) NOT NULL,
  `employee_identifier` int NOT NULL,
  `incident_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `description` text NOT NULL,
  `description_of_what_done` text,
  `is_done` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `force_majeure`
--

LOCK TABLES `force_majeure` WRITE;
/*!40000 ALTER TABLE `force_majeure` DISABLE KEYS */;
INSERT INTO `force_majeure` VALUES (1,1,'Daryn','Lift1','77009654486',2088886948,'2023-07-03 09:19:17','Some details','Something',1),(2,1,'Daryn','Lift1','77009654486',2088886948,'2023-05-05 10:49:12','Some details2',NULL,0),(3,1,'Daryn','Lift1','77009654486',2088886948,'2023-07-05 10:51:22','Some 2',NULL,0),(4,2,'Daryn','Lift1','77009654486',2088886948,'2023-07-05 11:03:29','Taskid2',NULL,0),(5,1,'Daryn','Lift1','77009654486',2088886948,'2023-07-05 11:26:33','De',NULL,0),(6,2,'Daryn','Lift1','77009654486',2088886948,'2023-07-05 11:27:34','De',NULL,0),(7,3,'Daryn','Lift1','77009654486',2088886948,'2023-07-05 11:39:36','ForcemaJor By taskId3',NULL,1),(8,3,'Daryn','Lift1','77009654486',2088886948,'2023-07-05 11:46:34','ForcemaJor By taskId3','I did something',1);
/*!40000 ALTER TABLE `force_majeure` ENABLE KEYS */;
UNLOCK TABLES;




DROP TABLE IF EXISTS `change_requests`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `change_requests` (
  `id` int NOT NULL AUTO_INCREMENT,
  `task_id` int DEFAULT NULL,
  `residential_complex` varchar(255) NOT NULL,
  `elevator_name` varchar(255) NOT NULL,
  `employee_phone_number` varchar(15) NOT NULL,
  `employee_identifier` int NOT NULL,
  `incident_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `description` text NOT NULL,
  `description_of_what_done` text,
  `is_done` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `change_requests`
--

LOCK TABLES `change_requests` WRITE;
/*!40000 ALTER TABLE `change_requests` DISABLE KEYS */;
INSERT INTO `change_requests` VALUES (9,5,'Rakhmet','Lift2','77009654486',2088886948,'2023-07-10 07:32:35','Bla bla bla','Somechit',1);
/*!40000 ALTER TABLE `change_requests` ENABLE KEYS */;
UNLOCK TABLES;





DROP TABLE IF EXISTS `overdue_task`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `overdue_task` (
  `id` int NOT NULL AUTO_INCREMENT,
  `task_id` int NOT NULL,
  `name_resident` varchar(255) NOT NULL,
  `name_lift` varchar(255) NOT NULL,
  `phone_number` varchar(15) NOT NULL,
  `name_of_task` varchar(255) NOT NULL,
  `date_start` date NOT NULL,
  `date_end` date NOT NULL,
  `is_done_by_worker` tinyint(1) NOT NULL,
  `is_done_by_admin` tinyint(1) DEFAULT '0',
  `is_done_by_hr_manager` tinyint(1) DEFAULT '0',
  `description_by_admin` text,
  `description_by_hr_manager` text,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=40 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `overdue_task`
--

LOCK TABLES `overdue_task` WRITE;
/*!40000 ALTER TABLE `overdue_task` DISABLE KEYS */;
INSERT INTO `overdue_task` VALUES (38,3,'Daryn','Lift1','77009654486','Определение координат установки оборудования лифта','2023-07-04','2023-07-10',0,1,0,'Hdhdh',NULL),(39,3,'Daryn','Lift1','77009654486','Определение координат установки оборудования лифта','2023-07-04','2023-07-10',0,1,0,'Hdhdb',NULL);
/*!40000 ALTER TABLE `overdue_task` ENABLE KEYS */;
UNLOCK TABLES;




DROP TABLE IF EXISTS `recommendations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `recommendations` (
  `id` int NOT NULL AUTO_INCREMENT,
  `date_created` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `phone_number` varchar(15) NOT NULL,
  `description` text,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `recommendations`
--

LOCK TABLES `recommendations` WRITE;
/*!40000 ALTER TABLE `recommendations` DISABLE KEYS */;
INSERT INTO `recommendations` VALUES (1,'2023-07-12 07:17:54','UzuiNU','I want you to do something'),(2,'2023-07-12 11:29:22','UzuiNU','Xcf');
/*!40000 ALTER TABLE `recommendations` ENABLE KEYS */;
UNLOCK TABLES;


DROP TABLE IF EXISTS `task_of_lifts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `task_of_lifts` (
  `id` int NOT NULL AUTO_INCREMENT,
  `task_name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `task_of_lifts`
--

LOCK TABLES `task_of_lifts` WRITE;
/*!40000 ALTER TABLE `task_of_lifts` DISABLE KEYS */;
UPDATE task_of_lifts SET task_name=CONVERT(CONVERT(task_name USING binary) USING utf8);
INSERT INTO `task_of_lifts` VALUES (1,'Приемка строительной части лифта'),(2,'Приемка оборудования и технической документации для монтажа, замены, модернизации'),(3,'Определение координат установки оборудования лифта'),(4,'Установка кронштейнов крепления направляющих кабины и противовеса'),(5,'Монтаж направляющих кабины и противовеса'),(6,'Монтаж дверей шахты'),(7,'Монтаж оборудования приямка'),(8,'Монтаж противовеса'),(9,'Монтаж кабины'),(10,'Монтаж лебедки главного привода'),(11,'Монтаж ограничителя скорости'),(12,'Навеска гибких тяговых элементов'),(13,'Монтаж электроаппаратуры, кабелей, электропроводки и цепей заземления'),(14,'Пусконаладочные работы');
/*!40000 ALTER TABLE `task_of_lifts` ENABLE KEYS */;
UNLOCK TABLES;




DROP TABLE IF EXISTS `lifts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `lifts` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name_lift` varchar(255) NOT NULL,
  `worker_id` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `worker_id` (`worker_id`),
  CONSTRAINT `lifts_ibfk_1` FOREIGN KEY (`worker_id`) REFERENCES `workers` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `lifts`
--

LOCK TABLES `lifts` WRITE;
/*!40000 ALTER TABLE `lifts` DISABLE KEYS */;
INSERT INTO `lifts` VALUES (1,'Lift1',1),(2,'Lift1',1),(3,'Lift2',1);
/*!40000 ALTER TABLE `lifts` ENABLE KEYS */;
UNLOCK TABLES;




DROP TABLE IF EXISTS `tasks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tasks` (
  `id` int NOT NULL AUTO_INCREMENT,
  `nameOfTask` varchar(255) NOT NULL,
  `dateStart` date DEFAULT NULL,
  `dateEnd` date DEFAULT NULL,
  `isDone` tinyint(1) DEFAULT '0',
  `lift_id` int DEFAULT NULL,
  `is_validate` tinyint(1) DEFAULT '0',
  `file_id` varchar(255) DEFAULT NULL,
  `date_requested_to_validate` date DEFAULT NULL,
  `is_rejected` tinyint(1) DEFAULT '0',
  `reject_description` text,
  PRIMARY KEY (`id`),
  KEY `lift_id` (`lift_id`),
  CONSTRAINT `tasks_ibfk_1` FOREIGN KEY (`lift_id`) REFERENCES `lifts` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tasks`
--

LOCK TABLES `tasks` WRITE;
/*!40000 ALTER TABLE `tasks` DISABLE KEYS */;
UPDATE tasks SET nameOfTask=CONVERT(CONVERT(nameOfTask USING binary) USING utf8); 
INSERT INTO `tasks` VALUES (1,'Приемка строительной части лифта','2023-07-04','2023-07-13',1,2,0,NULL,'2023-07-10',0,NULL),(2,'Приемка оборудования и технической документации для монтажа, замены, модернизации','2023-07-04','2023-07-13',1,2,0,NULL,'2023-07-09',0,NULL),(3,'Определение координат установки оборудования лифта','2023-07-04','2023-07-13',1,2,0,NULL,'2023-07-11',0,NULL),(4,'Приемка строительной части лифта','2023-07-04','2023-07-15',1,3,0,NULL,NULL,0,NULL),(5,'Приемка оборудования и технической документации для монтажа, замены, модернизации','2023-07-04','2023-07-15',1,3,1,'BAACAgIAAxkBAAII92SuPpxWvVBulPnLOQtYXdVISAu-AAJtOAACKH5wSbd9rLM_Wqg6LwQ','2023-07-12',1,'Фото офиса'),(6,'Определение координат установки оборудования лифта','2023-07-04','2023-07-15',0,3,0,NULL,NULL,0,NULL);
/*!40000 ALTER TABLE `tasks` ENABLE KEYS */;
UNLOCK TABLES;



DROP TABLE IF EXISTS `projects`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `projects` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name_resident` varchar(255) NOT NULL,
  `worker_id` int DEFAULT NULL,
  `lift_id` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `worker_id` (`worker_id`),
  KEY `lift_id` (`lift_id`),
  CONSTRAINT `projects_ibfk_1` FOREIGN KEY (`worker_id`) REFERENCES `workers` (`id`),
  CONSTRAINT `projects_ibfk_2` FOREIGN KEY (`lift_id`) REFERENCES `lifts` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `projects`
--

LOCK TABLES `projects` WRITE;
/*!40000 ALTER TABLE `projects` DISABLE KEYS */;
INSERT INTO `projects` VALUES (1,'Daryn',1,2),(2,'Rakhmet',1,3);
/*!40000 ALTER TABLE `projects` ENABLE KEYS */;
UNLOCK TABLES;