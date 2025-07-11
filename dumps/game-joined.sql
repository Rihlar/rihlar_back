-- MySQL dump 10.13  Distrib 9.3.0, for Linux (aarch64)
--
-- Host: localhost    Database: maindb
-- ------------------------------------------------------
-- Server version	9.3.0

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `Circles`
--

DROP TABLE IF EXISTS `Circles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Circles` (
  `circle_id` varchar(191) NOT NULL,
  `game_id` longtext,
  `team_id` longtext,
  `user_id` longtext,
  `size` bigint NOT NULL,
  `level` bigint NOT NULL,
  `latitude` double DEFAULT NULL,
  `longitude` double DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `image_id` longtext,
  `steps` bigint DEFAULT NULL,
  PRIMARY KEY (`circle_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Circles`
--

LOCK TABLES `Circles` WRITE;
/*!40000 ALTER TABLE `Circles` DISABLE KEYS */;
/*!40000 ALTER TABLE `Circles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Members`
--

DROP TABLE IF EXISTS `Members`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Members` (
  `game_id` varchar(50) NOT NULL,
  `team_id` varchar(50) NOT NULL,
  `user_id` varchar(191) NOT NULL,
  `points` bigint NOT NULL,
  PRIMARY KEY (`game_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Members`
--

LOCK TABLES `Members` WRITE;
/*!40000 ALTER TABLE `Members` DISABLE KEYS */;
INSERT INTO `Members` VALUES ('f36eb7ce-4e24-4805-99a5-b3ae3468708a','b5fef636-b22e-4057-b1fe-acc7bde6add0','e3abf90d-4bcf-4c3b-bbde-37694b1611b3',0),('gameid-413a287b-213c-414f-a287-c1397db8f9bf','teamid-c6799da3-6d9c-449c-98c1-86108faf3b00','userid-79541130-3275-4b90-8677-01323045aca5',0),('gameid-413a287b-213c-414f-a287-c1397db8f9bf','teamid-06adac36-33ef-42ae-af2e-6209d321302f','userid-a1b2c3d4-e5f6-7890-1234-567890abcdef',0),('gameid-413a287b-213c-414f-a287-c1397db8f9bf','teamid-73004fb5-f235-492b-987a-843794783652','userid-a7b8c9d0-e1f2-3456-7890-abcdef012345',0),('gameid-413a287b-213c-414f-a287-c1397db8f9bf','teamid-b8fa6afd-8f9f-4219-8502-db746005ce54','userid-b2c3d4e5-f6a7-8901-2345-67890abcdef0',0),('gameid-413a287b-213c-414f-a287-c1397db8f9bf','teamid-d798b901-455a-471f-8a9d-cef467896cbe','userid-b8c9d0e1-f2a3-4567-890a-bcdef0123456',0),('gameid-413a287b-213c-414f-a287-c1397db8f9bf','teamid-c09427ac-0598-43b5-8dbd-dbbbb8d43044','userid-c3d4e5f6-a7b8-9012-3456-7890abcdef01',0),('gameid-413a287b-213c-414f-a287-c1397db8f9bf','teamid-00895409-ee75-43c9-a973-69552293c350','userid-c9d0e1f2-a3b4-5678-90ab-cdef01234567',0),('gameid-413a287b-213c-414f-a287-c1397db8f9bf','teamid-ed9befa3-f443-4e82-9ddb-dfa69ef6fd86','userid-d0e1f2a3-b4c5-6789-0abc-def012345678',0),('gameid-413a287b-213c-414f-a287-c1397db8f9bf','teamid-9faa04db-8a3e-4117-90ec-63bc276be6d4','userid-d4e5f6a7-b8c9-0123-4567-890abcdef012',0),('gameid-413a287b-213c-414f-a287-c1397db8f9bf','teamid-4ed878c4-3be2-4e93-a2cd-5954eaebfc54','userid-e5f6a7b8-c9d0-1234-5678-90abcdef0123',0),('gameid-413a287b-213c-414f-a287-c1397db8f9bf','teamid-57a4e8b1-fa30-413e-bcd1-158f243f7a5f','userid-f224c85f-2ac4-4f50-a9dc-af80a659b671',0),('gameid-413a287b-213c-414f-a287-c1397db8f9bf','teamid-35c743c1-cc4c-4c6a-9012-dee65cf70301','userid-f6a7b8c9-d0e1-2345-6789-0abcdef01234',0),('gameid-9fcb784b-04a8-49c3-9ed9-ca9588eb86a8','teamid-957c39a5-3573-488c-9ba5-4027eda85bca','userid-79541130-3275-4b90-8677-01323045aca5',0),('gameid-a4f3e2d1-c0b9-8765-4321-0fedcba98765','teamid-2fc371c9-22d7-4a22-a098-3dd6cc782265','userid-f6a7b8c9-d0e1-2345-6789-0abcdef01234',0),('gameid-b385def4-5833-4153-b54e-bdb134ab9fc8','teamid-c660f82d-878e-4a5a-baae-dbe40e0c3edc','userid-f224c85f-2ac4-4f50-a9dc-af80a659b671',0),('gameid-b5a4f3e2-d1c0-9876-5432-10fedcba9876','teamid-bba4c683-6fe9-4aac-9777-fca946843fd2','userid-e5f6a7b8-c9d0-1234-5678-90abcdef0123',0),('gameid-c0b9a8f7-e6d5-4321-0987-fedcba987651','teamid-907f391a-233f-46e0-9f68-660fe35c0a5a','userid-d0e1f2a3-b4c5-6789-0abc-def012345678',0),('gameid-c6b5a4f3-e2d1-0987-6543-210fedcba987','teamid-5f10a41b-0b6c-4c78-9849-c6d278eca161','userid-d4e5f6a7-b8c9-0123-4567-890abcdef012',0),('gameid-d1c0b9a8-f7e6-5432-1098-fedcba987652','teamid-be7b4a89-1e9f-4494-b9e2-f0ed92d27504','userid-c9d0e1f2-a3b4-5678-90ab-cdef01234567',0),('gameid-d7c6b5a4-f3e2-1098-7654-3210fedcba98','teamid-835f0493-40e2-4248-8c17-bb21ad889af5','userid-c3d4e5f6-a7b8-9012-3456-7890abcdef01',0),('gameid-e2d1c0b9-a8f7-6543-2109-fedcba987653','teamid-681b3089-8324-461e-84e3-e8dde99c3efb','userid-b8c9d0e1-f2a3-4567-890a-bcdef0123456',0),('gameid-e8d7c6b5-a4f3-2109-8765-43210fedcba9','teamid-6423fbed-4e1f-45f0-9d15-4241b28aaea4','userid-b2c3d4e5-f6a7-8901-2345-67890abcdef0',0),('gameid-f3e2d1c0-b9a8-7654-3210-fedcba987654','teamid-8d0fefc3-a2fb-42f1-88bd-3d878d499394','userid-a7b8c9d0-e1f2-3456-7890-abcdef012345',0),('gameid-f9e8d7c6-b5a4-3210-fedc-ba9876543210','teamid-6c9d6cde-de7f-4f10-8495-de0e334dfc27','userid-a1b2c3d4-e5f6-7890-1234-567890abcdef',0);
/*!40000 ALTER TABLE `Members` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Profile`
--

DROP TABLE IF EXISTS `Profile`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Profile` (
  `user_id` varchar(50) NOT NULL,
  `record_id` varchar(50) DEFAULT NULL,
  `comment` varchar(255) DEFAULT '',
  `latitude` double DEFAULT '0',
  `longitude` double DEFAULT '0',
  `size` bigint DEFAULT '0',
  `region_id` varchar(50) DEFAULT '',
  `sys_game` varchar(50) DEFAULT '',
  `adm_game` varchar(50) DEFAULT '',
  `name` varchar(100) DEFAULT '',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Profile`
--

LOCK TABLES `Profile` WRITE;
/*!40000 ALTER TABLE `Profile` DISABLE KEYS */;
INSERT INTO `Profile` VALUES ('userid-79541130-3275-4b90-8677-01323045aca5','','',0,0,0,'','gameid-9fcb784b-04a8-49c3-9ed9-ca9588eb86a8','gameid-413a287b-213c-414f-a287-c1397db8f9bf','Alice'),('userid-a1b2c3d4-e5f6-7890-1234-567890abcdef','','',0,0,0,'','gameid-f9e8d7c6-b5a4-3210-fedc-ba9876543210','gameid-413a287b-213c-414f-a287-c1397db8f9bf','Charlie'),('userid-a7b8c9d0-e1f2-3456-7890-abcdef012345','','',0,0,0,'','gameid-f3e2d1c0-b9a8-7654-3210-fedcba987654','gameid-413a287b-213c-414f-a287-c1397db8f9bf','Ivan'),('userid-b2c3d4e5-f6a7-8901-2345-67890abcdef0','','',0,0,0,'','gameid-e8d7c6b5-a4f3-2109-8765-43210fedcba9','gameid-413a287b-213c-414f-a287-c1397db8f9bf','David'),('userid-b8c9d0e1-f2a3-4567-890a-bcdef0123456','','',0,0,0,'','gameid-e2d1c0b9-a8f7-6543-2109-fedcba987653','gameid-413a287b-213c-414f-a287-c1397db8f9bf','Judy'),('userid-c3d4e5f6-a7b8-9012-3456-7890abcdef01','','',0,0,0,'','gameid-d7c6b5a4-f3e2-1098-7654-3210fedcba98','gameid-413a287b-213c-414f-a287-c1397db8f9bf','Eve'),('userid-c9d0e1f2-a3b4-5678-90ab-cdef01234567','','',0,0,0,'','gameid-d1c0b9a8-f7e6-5432-1098-fedcba987652','gameid-413a287b-213c-414f-a287-c1397db8f9bf','Kevin'),('userid-d0e1f2a3-b4c5-6789-0abc-def012345678','','',0,0,0,'','gameid-c0b9a8f7-e6d5-4321-0987-fedcba987651','gameid-413a287b-213c-414f-a287-c1397db8f9bf','Liam'),('userid-d4e5f6a7-b8c9-0123-4567-890abcdef012','','',0,0,0,'','gameid-c6b5a4f3-e2d1-0987-6543-210fedcba987','gameid-413a287b-213c-414f-a287-c1397db8f9bf','Frank'),('userid-e5f6a7b8-c9d0-1234-5678-90abcdef0123','','',0,0,0,'','gameid-b5a4f3e2-d1c0-9876-5432-10fedcba9876','gameid-413a287b-213c-414f-a287-c1397db8f9bf','Grace'),('userid-f224c85f-2ac4-4f50-a9dc-af80a659b671','','',0,0,0,'','gameid-b385def4-5833-4153-b54e-bdb134ab9fc8','gameid-413a287b-213c-414f-a287-c1397db8f9bf','Bob'),('userid-f6a7b8c9-d0e1-2345-6789-0abcdef01234','','',0,0,0,'','gameid-a4f3e2d1-c0b9-8765-4321-0fedcba98765','gameid-413a287b-213c-414f-a287-c1397db8f9bf','Heidi');
/*!40000 ALTER TABLE `Profile` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Teams`
--

DROP TABLE IF EXISTS `Teams`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Teams` (
  `team_id` varchar(50) NOT NULL,
  `game_id` varchar(50) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `points` bigint NOT NULL,
  PRIMARY KEY (`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Teams`
--

LOCK TABLES `Teams` WRITE;
/*!40000 ALTER TABLE `Teams` DISABLE KEYS */;
INSERT INTO `Teams` VALUES ('b5fef636-b22e-4057-b1fe-acc7bde6add0','f36eb7ce-4e24-4805-99a5-b3ae3468708a','2025-07-11 04:59:44.904',0),('teamid-00895409-ee75-43c9-a973-69552293c350','gameid-413a287b-213c-414f-a287-c1397db8f9bf','2025-07-11 04:59:44.999',0),('teamid-06adac36-33ef-42ae-af2e-6209d321302f','gameid-413a287b-213c-414f-a287-c1397db8f9bf','2025-07-11 04:59:44.954',0),('teamid-2fc371c9-22d7-4a22-a098-3dd6cc782265','gameid-a4f3e2d1-c0b9-8765-4321-0fedcba98765','2025-07-11 04:59:44.980',0),('teamid-35c743c1-cc4c-4c6a-9012-dee65cf70301','gameid-413a287b-213c-414f-a287-c1397db8f9bf','2025-07-11 04:59:44.983',0),('teamid-4ed878c4-3be2-4e93-a2cd-5954eaebfc54','gameid-413a287b-213c-414f-a287-c1397db8f9bf','2025-07-11 04:59:44.978',0),('teamid-57a4e8b1-fa30-413e-bcd1-158f243f7a5f','gameid-413a287b-213c-414f-a287-c1397db8f9bf','2025-07-11 04:59:44.949',0),('teamid-5f10a41b-0b6c-4c78-9849-c6d278eca161','gameid-c6b5a4f3-e2d1-0987-6543-210fedcba987','2025-07-11 04:59:44.970',0),('teamid-6423fbed-4e1f-45f0-9d15-4241b28aaea4','gameid-e8d7c6b5-a4f3-2109-8765-43210fedcba9','2025-07-11 04:59:44.957',0),('teamid-681b3089-8324-461e-84e3-e8dde99c3efb','gameid-e2d1c0b9-a8f7-6543-2109-fedcba987653','2025-07-11 04:59:44.992',0),('teamid-6c9d6cde-de7f-4f10-8495-de0e334dfc27','gameid-f9e8d7c6-b5a4-3210-fedc-ba9876543210','2025-07-11 04:59:44.952',0),('teamid-73004fb5-f235-492b-987a-843794783652','gameid-413a287b-213c-414f-a287-c1397db8f9bf','2025-07-11 04:59:44.990',0),('teamid-835f0493-40e2-4248-8c17-bb21ad889af5','gameid-d7c6b5a4-f3e2-1098-7654-3210fedcba98','2025-07-11 04:59:44.962',0),('teamid-8d0fefc3-a2fb-42f1-88bd-3d878d499394','gameid-f3e2d1c0-b9a8-7654-3210-fedcba987654','2025-07-11 04:59:44.987',0),('teamid-907f391a-233f-46e0-9f68-660fe35c0a5a','gameid-c0b9a8f7-e6d5-4321-0987-fedcba987651','2025-07-11 04:59:45.001',0),('teamid-957c39a5-3573-488c-9ba5-4027eda85bca','gameid-9fcb784b-04a8-49c3-9ed9-ca9588eb86a8','2025-07-11 04:59:44.938',0),('teamid-9faa04db-8a3e-4117-90ec-63bc276be6d4','gameid-413a287b-213c-414f-a287-c1397db8f9bf','2025-07-11 04:59:44.973',0),('teamid-b8fa6afd-8f9f-4219-8502-db746005ce54','gameid-413a287b-213c-414f-a287-c1397db8f9bf','2025-07-11 04:59:44.959',0),('teamid-bba4c683-6fe9-4aac-9777-fca946843fd2','gameid-b5a4f3e2-d1c0-9876-5432-10fedcba9876','2025-07-11 04:59:44.975',0),('teamid-be7b4a89-1e9f-4494-b9e2-f0ed92d27504','gameid-d1c0b9a8-f7e6-5432-1098-fedcba987652','2025-07-11 04:59:44.997',0),('teamid-c09427ac-0598-43b5-8dbd-dbbbb8d43044','gameid-413a287b-213c-414f-a287-c1397db8f9bf','2025-07-11 04:59:44.968',0),('teamid-c660f82d-878e-4a5a-baae-dbe40e0c3edc','gameid-b385def4-5833-4153-b54e-bdb134ab9fc8','2025-07-11 04:59:44.947',0),('teamid-c6799da3-6d9c-449c-98c1-86108faf3b00','gameid-413a287b-213c-414f-a287-c1397db8f9bf','2025-07-11 04:59:44.941',0),('teamid-d798b901-455a-471f-8a9d-cef467896cbe','gameid-413a287b-213c-414f-a287-c1397db8f9bf','2025-07-11 04:59:44.995',0),('teamid-ed9befa3-f443-4e82-9ddb-dfa69ef6fd86','gameid-413a287b-213c-414f-a287-c1397db8f9bf','2025-07-11 04:59:45.005',0);
/*!40000 ALTER TABLE `Teams` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `game_chunks`
--

DROP TABLE IF EXISTS `game_chunks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `game_chunks` (
  `chunk_id` varchar(191) NOT NULL,
  `game_id` longtext,
  `image_id` longtext,
  `owner_id` longtext,
  `start_lat` double NOT NULL,
  `start_lon` double NOT NULL,
  `end_lat` double NOT NULL,
  `end_lon` double NOT NULL,
  `level` bigint NOT NULL,
  PRIMARY KEY (`chunk_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `game_chunks`
--

LOCK TABLES `game_chunks` WRITE;
/*!40000 ALTER TABLE `game_chunks` DISABLE KEYS */;
INSERT INTO `game_chunks` VALUES ('3325d4ee-ef32-42a3-91d1-33d3582dffc2','f36eb7ce-4e24-4805-99a5-b3ae3468708a','76bd1e16-3105-4916-ad6b-7da9554c9601','e9178c88-3b64-4e61-b823-fd874d177d3c',0,0,0,0,2);
/*!40000 ALTER TABLE `game_chunks` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `games`
--

DROP TABLE IF EXISTS `games`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `games` (
  `game_id` varchar(50) NOT NULL,
  `start_time` datetime(3) NOT NULL,
  `end_time` datetime(3) NOT NULL,
  `flag` bigint NOT NULL,
  `type` bigint NOT NULL,
  `status` bigint NOT NULL,
  `region_id` longtext,
  PRIMARY KEY (`game_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `games`
--

LOCK TABLES `games` WRITE;
/*!40000 ALTER TABLE `games` DISABLE KEYS */;
INSERT INTO `games` VALUES ('gameid-413a287b-213c-414f-a287-c1397db8f9bf','2025-07-11 04:59:44.937','2025-07-31 04:59:44.937',0,1,1,'regionId-c161edb9-6aff-4244-8749-707bff2fa3be'),('gameid-9fcb784b-04a8-49c3-9ed9-ca9588eb86a8','2025-07-11 04:59:44.916','2025-07-31 04:59:44.916',0,0,1,'regionId-c161edb9-6aff-4244-8749-707bff2fa3be'),('gameid-a4f3e2d1-c0b9-8765-4321-0fedcba98765','2025-07-11 04:59:44.932','2025-07-31 04:59:44.932',0,0,1,'regionId-c161edb9-6aff-4244-8749-707bff2fa3be'),('gameid-b385def4-5833-4153-b54e-bdb134ab9fc8','2025-07-11 04:59:44.918','2025-07-31 04:59:44.918',0,0,1,'regionId-c161edb9-6aff-4244-8749-707bff2fa3be'),('gameid-b5a4f3e2-d1c0-9876-5432-10fedcba9876','2025-07-11 04:59:44.930','2025-07-31 04:59:44.930',0,0,1,'regionId-c161edb9-6aff-4244-8749-707bff2fa3be'),('gameid-c0b9a8f7-e6d5-4321-0987-fedcba987651','2025-07-11 04:59:44.936','2025-07-31 04:59:44.936',0,0,1,'regionId-c161edb9-6aff-4244-8749-707bff2fa3be'),('gameid-c6b5a4f3-e2d1-0987-6543-210fedcba987','2025-07-11 04:59:44.929','2025-07-31 04:59:44.929',0,0,1,'regionId-c161edb9-6aff-4244-8749-707bff2fa3be'),('gameid-d1c0b9a8-f7e6-5432-1098-fedcba987652','2025-07-11 04:59:44.935','2025-07-31 04:59:44.935',0,0,1,'regionId-c161edb9-6aff-4244-8749-707bff2fa3be'),('gameid-d7c6b5a4-f3e2-1098-7654-3210fedcba98','2025-07-11 04:59:44.928','2025-07-31 04:59:44.928',0,0,1,'regionId-c161edb9-6aff-4244-8749-707bff2fa3be'),('gameid-e2d1c0b9-a8f7-6543-2109-fedcba987653','2025-07-11 04:59:44.934','2025-07-31 04:59:44.934',0,0,1,'regionId-c161edb9-6aff-4244-8749-707bff2fa3be'),('gameid-e8d7c6b5-a4f3-2109-8765-43210fedcba9','2025-07-11 04:59:44.927','2025-07-31 04:59:44.927',0,0,1,'regionId-c161edb9-6aff-4244-8749-707bff2fa3be'),('gameid-f3e2d1c0-b9a8-7654-3210-fedcba987654','2025-07-11 04:59:44.933','2025-07-31 04:59:44.933',0,0,1,'regionId-c161edb9-6aff-4244-8749-707bff2fa3be'),('gameid-f9e8d7c6-b5a4-3210-fedc-ba9876543210','2025-07-11 04:59:44.925','2025-07-31 04:59:44.925',0,0,1,'regionId-c161edb9-6aff-4244-8749-707bff2fa3be');
/*!40000 ALTER TABLE `games` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `movementLog`
--

DROP TABLE IF EXISTS `movementLog`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `movementLog` (
  `movement_id` varchar(191) NOT NULL,
  `user_id` varchar(191) NOT NULL,
  `latitude` double DEFAULT NULL,
  `longitude` double DEFAULT NULL,
  `steps` bigint DEFAULT NULL,
  `game_id` varchar(191) NOT NULL,
  `time_stamp` bigint NOT NULL,
  PRIMARY KEY (`movement_id`,`user_id`,`game_id`,`time_stamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `movementLog`
--

LOCK TABLES `movementLog` WRITE;
/*!40000 ALTER TABLE `movementLog` DISABLE KEYS */;
/*!40000 ALTER TABLE `movementLog` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `regions`
--

DROP TABLE IF EXISTS `regions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `regions` (
  `region_id` varchar(191) NOT NULL,
  `region_name` longtext NOT NULL,
  `start_lat` double NOT NULL,
  `start_lon` double NOT NULL,
  `end_lat` double NOT NULL,
  `end_lon` double NOT NULL,
  PRIMARY KEY (`region_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `regions`
--

LOCK TABLES `regions` WRITE;
/*!40000 ALTER TABLE `regions` DISABLE KEYS */;
INSERT INTO `regions` VALUES ('regionId-005344a4-523d-4aab-9d1e-d232322cf54e','中国',36,130.5,33.5,134.5),('regionId-16d687a3-8eab-4b36-8563-b62514823fe8','中部',38,135.5,34,139.5),('regionId-3501902d-8bab-40cd-926f-30c53d80efc5','四国',34.5,132,32.5,135.5),('regionId-65051f6a-9e94-439d-a7f6-5c127ad0c885','東北',41.6,138,37,142.5),('regionId-a3d3dcd1-7a73-4a31-9908-b9cab944280d','北海道',45.55,139.5,41.3,148.5),('regionId-c161edb9-6aff-4244-8749-707bff2fa3be','関西',35.8,134,33.5,137.5),('regionId-ef5aa179-53e0-481d-b64d-ae7654049a88','九州',34.5,128,30.5,132.5),('regionId-fb145c05-e0e5-4f22-86e1-9f40326faf31','関東',37,138.5,34.5,141);
/*!40000 ALTER TABLE `regions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `samples`
--

DROP TABLE IF EXISTS `samples`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `samples` (
  `user_id` varchar(191) NOT NULL,
  `name` varchar(100) NOT NULL,
  `email` varchar(191) NOT NULL,
  `age` bigint DEFAULT '0',
  `is_active` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `uni_samples_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `samples`
--

LOCK TABLES `samples` WRITE;
/*!40000 ALTER TABLE `samples` DISABLE KEYS */;
INSERT INTO `samples` VALUES ('e3abf90d-4bcf-4c3b-bbde-37694b1611b3','aiueo','test@mattuu.com',20,1);
/*!40000 ALTER TABLE `samples` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-07-11 14:01:01
