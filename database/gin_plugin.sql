CREATE DATABASE  IF NOT EXISTS `gin` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin */;
USE `gin`;
-- MySQL dump 10.13  Distrib 8.0.15, for Win64 (x86_64)
--
-- Host: localhost    Database: gin
-- ------------------------------------------------------
-- Server version	8.0.15

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
 SET NAMES utf8 ;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `plugin`
--

DROP TABLE IF EXISTS `plugin`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `plugin` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(45) COLLATE utf8mb4_bin NOT NULL,
  `clean_name` varchar(45) COLLATE utf8mb4_bin NOT NULL,
  `uuid` varchar(45) COLLATE utf8mb4_bin NOT NULL,
  `description` longtext COLLATE utf8mb4_bin,
  `installation_date` datetime DEFAULT NULL,
  `status_id` int(11) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `config` longtext COLLATE utf8mb4_bin NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  KEY `status_id_status_idx` (`status_id`),
  CONSTRAINT `status_id_status` FOREIGN KEY (`status_id`) REFERENCES `plugin_status` (`id`) ON DELETE SET NULL ON UPDATE SET NULL
) ENGINE=InnoDB AUTO_INCREMENT=66 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `plugin`
--

LOCK TABLES `plugin` WRITE;
/*!40000 ALTER TABLE `plugin` DISABLE KEYS */;
INSERT INTO `plugin` VALUES (64,'plugin-1','plugin1','3fc198ac-1559-492b-ab2f-ff9be488be71','First plugin for testing purpouses','2020-05-19 12:43:26',3,'2020-05-19 12:43:26','2020-05-19 12:43:36','2020-05-19 12:43:59','{\r\n  \"name\": \"plugin-1\",\r\n  \"clean_name\":\"plugin1\",\r\n  \"uuid\": \"3fc198ac-1559-492b-ab2f-ff9be488be71\",\r\n  \"description\": \"First plugin for testing purpouses\",\r\n  \"install_elements\":[\"html\", \"css\", \"js\", \"controller\", \"routing\", \"configuration\"]\r\n}'),(65,'plugin-1','plugin1','3fc198ac-1559-492b-ab2f-ff9be488be71','First plugin for testing purpouses','2020-05-20 16:31:50',3,'2020-05-20 16:31:50','2020-05-20 16:32:05','2020-05-20 16:33:34','{\r\n  \"name\": \"plugin-1\",\r\n  \"clean_name\":\"plugin1\",\r\n  \"uuid\": \"3fc198ac-1559-492b-ab2f-ff9be488be71\",\r\n  \"description\": \"First plugin for testing purpouses\",\r\n  \"install_elements\":[\"html\", \"css\", \"js\", \"controller\", \"routing\", \"configuration\"]\r\n}');
/*!40000 ALTER TABLE `plugin` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-06-04 16:46:54
