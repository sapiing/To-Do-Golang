-- --------------------------------------------------------
-- Host:                         127.0.0.1
-- Server version:               8.0.30 - MySQL Community Server - GPL
-- Server OS:                    Win64
-- HeidiSQL Version:             12.7.0.6850
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


-- Dumping database structure for todo_golang
CREATE DATABASE IF NOT EXISTS `todo_golang` /*!40100 DEFAULT CHARACTER SET latin1 */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `todo_golang`;

-- Dumping structure for table todo_golang.tasks
CREATE TABLE IF NOT EXISTS `tasks` (
  `id` int NOT NULL AUTO_INCREMENT,
  `title` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL DEFAULT '',
  `description` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `completed` tinyint(1) DEFAULT '0',
  `date` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=latin1;

-- Dumping data for table todo_golang.tasks: ~3 rows (approximately)
INSERT INTO `tasks` (`id`, `title`, `description`, `completed`, `date`) VALUES
	(1, 'Cuci Piring', 'Cuci piring yang ada di tempat cucian\n', 0, '2024-07-10 13:16:11'),
	(2, 'Sapu', 'Sapu rumah hingga bersih', 0, '2024-07-10 13:17:47'),
	(3, 'Tidur', 'zzz', 0, '2024-07-10 13:18:26'),
	(4, 'Aziek', 'Test', 0, '2024-07-12 13:25:15');

-- Dumping structure for table todo_golang.work_log
CREATE TABLE IF NOT EXISTS `work_log` (
  `id` int NOT NULL AUTO_INCREMENT,
  `task_id` int NOT NULL,
  `date` date NOT NULL,
  `completed` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `task_id` (`task_id`),
  CONSTRAINT `work_log_ibfk_1` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=116 DEFAULT CHARSET=latin1;

-- Dumping data for table todo_golang.work_log: ~4 rows (approximately)
INSERT INTO `work_log` (`id`, `task_id`, `date`, `completed`, `created_at`, `updated_at`) VALUES
	(112, 1, '2024-07-12', 0, '2024-07-12 13:49:36', '2024-07-12 13:49:36'),
	(113, 2, '2024-07-12', 0, '2024-07-12 13:49:36', '2024-07-12 13:49:36'),
	(114, 3, '2024-07-12', 0, '2024-07-12 13:49:36', '2024-07-12 13:49:36'),
	(115, 4, '2024-07-12', 0, '2024-07-12 13:49:36', '2024-07-12 13:49:36');

-- Dumping structure for table todo_golang.work_log_history
CREATE TABLE IF NOT EXISTS `work_log_history` (
  `id` int NOT NULL AUTO_INCREMENT,
  `task_id` int NOT NULL,
  `date` date NOT NULL,
  `completed` tinyint(1) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT (now()),
  PRIMARY KEY (`id`),
  KEY `task_id` (`task_id`),
  CONSTRAINT `work_log_history_ibfk_1` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=64 DEFAULT CHARSET=latin1;

-- Dumping data for table todo_golang.work_log_history: ~12 rows (approximately)
INSERT INTO `work_log_history` (`id`, `task_id`, `date`, `completed`, `created_at`, `updated_at`) VALUES
	(16, 1, '2024-07-14', 0, '2024-07-12 13:22:18', '2024-07-12 13:35:25'),
	(17, 2, '2024-07-12', 0, '2024-07-12 13:22:18', '2024-07-12 13:35:25'),
	(18, 3, '2024-07-12', 1, '2024-07-12 13:22:18', '2024-07-12 13:35:25'),
	(19, 1, '2024-07-12', 1, '2024-07-12 13:24:18', '2024-07-12 13:35:25'),
	(20, 2, '2024-07-09', 1, '2024-07-12 13:24:18', '2024-07-12 13:35:25'),
	(21, 3, '2024-07-12', 0, '2024-07-12 13:24:18', '2024-07-12 13:35:25'),
	(22, 1, '2024-07-12', 1, '2024-07-12 13:26:18', '2024-07-12 13:35:25'),
	(23, 2, '2024-07-06', 1, '2024-07-12 13:26:19', '2024-07-12 13:35:25'),
	(50, 1, '2024-07-12', 0, '2024-07-12 13:41:30', '2024-07-12 13:42:53'),
	(51, 2, '2024-07-12', 0, '2024-07-12 13:41:30', '2024-07-12 13:42:53'),
	(52, 3, '2024-07-12', 0, '2024-07-12 13:41:30', '2024-07-12 13:42:53'),
	(53, 4, '2024-07-12', 0, '2024-07-12 13:41:30', '2024-07-12 13:42:53'),
	(57, 1, '2024-07-12', 0, '2024-07-12 13:42:53', '2024-07-12 13:49:36'),
	(58, 2, '2024-07-12', 0, '2024-07-12 13:42:53', '2024-07-12 13:49:36'),
	(59, 3, '2024-07-12', 0, '2024-07-12 13:42:53', '2024-07-12 13:49:36'),
	(60, 4, '2024-07-12', 0, '2024-07-12 13:42:53', '2024-07-12 13:49:36');

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
