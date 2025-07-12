CREATE TABLE IF NOT EXISTS `STUDENTS` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `first_name` VARCHAR(50),
  `last_name` VARCHAR(50),
  `email` VARCHAR(255),
  `password` VARCHAR(255) COMMENT 'hashed, to be added eventually'
);

CREATE TABLE IF NOT EXISTS `STUDENT_INFO` (
  `student_id` int UNIQUE,
  `major` VARCHAR(255),
  `emplid` VARCHAR(255),
  `grad_year` int,
  `interests_id` int COMMENT 'FK',
  `dietary_restrictions` VARCHAR(255),
  `comments` VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS `STUDENTS_TO_INTERESTS` (
  `student_id` int COMMENT 'FK',
  `interest_id` int COMMENT 'FK'
);

CREATE TABLE IF NOT EXISTS `INTERESTS` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `label` VARCHAR(255) COMMENT 'career, community,, and volunteer when starting'
);

CREATE TABLE IF NOT EXISTS `CLUBS` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `club_name` VARCHAR(255),
  `club_icon` VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS `CLUB_MEMBERS` (
  `student_id` int,
  `club_id` int,
  `role` ENUM ('member', 'eboard'),
  PRIMARY KEY (`student_id`, `club_id`)
);

CREATE TABLE IF NOT EXISTS `EVENTS` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `current_version_id` int,
  `original_version_id` int
);

CREATE TABLE IF NOT EXISTS `EVENT_VERSIONS` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `author_id` int COMMENT 'FK',
  `event_name` VARCHAR(255),
  `event_img` VARCHAR(255),
  `event_status` ENUM ('drafted', 'posted'),
  `event_date` datetime,
  `type` ENUM ('create', 'edit', 'delete'),
  `timestamp` datetime
);

CREATE TABLE IF NOT EXISTS `EVENTS_TO_CLUBS` (
  `event_id` int COMMENT 'FK',
  `club_id` int COMMENT 'FK',
  PRIMARY KEY (`event_id`, `club_id`)
);

CREATE TABLE IF NOT EXISTS `POINT_SOURCES` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `title` VARCHAR(255),
  `points` int
);

CREATE TABLE IF NOT EXISTS `POINT_HISTORIES` (
  `member_id` int COMMENT 'FK',
  `point_source_id` int COMMENT 'FK',
  `points_earned` int,
  `points_after_gain` int,
  `timestamp` datetime,
  PRIMARY KEY (`member_id`, `point_source_id`)
);

CREATE TABLE IF NOT EXISTS `PURCHASES` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `member_id` int COMMENT 'FK',
  `item_title` VARCHAR(255),
  `points_cost` int,
  `points_after_purchase` int,
  `timestamp` datetime
);

ALTER TABLE `student_info` ADD CONSTRAINT `FK_StudentInfo_Students` FOREIGN KEY (`student_id`) REFERENCES `students` (`id`);

ALTER TABLE `students_to_interests` ADD CONSTRAINT `FK_StudentsToInterests_Students` FOREIGN KEY (`student_id`) REFERENCES `students` (`id`);

ALTER TABLE `students_to_interests` ADD CONSTRAINT `FK_StudentsToInterests_Interests` FOREIGN KEY (`interest_id`) REFERENCES `interests` (`id`);

ALTER TABLE `club_members` ADD CONSTRAINT `FK_ClubMembers_Students` FOREIGN KEY (`student_id`) REFERENCES `students` (`id`);

ALTER TABLE `club_members` ADD CONSTRAINT `FK_ClubMembers_Clubs` FOREIGN KEY (`club_id`) REFERENCES `clubs` (`id`);

ALTER TABLE `events` ADD CONSTRAINT `FK_Events_CurrentVersion` FOREIGN KEY (`current_version_id`) REFERENCES `event_versions` (`id`);

ALTER TABLE `events` ADD CONSTRAINT `FK_Events_OriginalVersion` FOREIGN KEY (`original_version_id`) REFERENCES `event_versions` (`id`);

ALTER TABLE `event_versions` ADD CONSTRAINT `FK_EventVersions_ClubMembers` FOREIGN KEY (`author_id`) REFERENCES `students` (`id`);

ALTER TABLE `events_to_clubs` ADD CONSTRAINT `FK_EventsToClubs_Events` FOREIGN KEY (`event_id`) REFERENCES `events` (`id`);

ALTER TABLE `events_to_clubs` ADD CONSTRAINT `FK_EventsToClubs_Clubs` FOREIGN KEY (`club_id`) REFERENCES `clubs` (`id`);

ALTER TABLE `point_histories` ADD CONSTRAINT `FK_PointHistories_ClubMembers` FOREIGN KEY (`member_id`) REFERENCES `students` (`id`);

ALTER TABLE `point_histories` ADD CONSTRAINT `FK_PointHistories_PointSources` FOREIGN KEY (`point_source_id`) REFERENCES `point_sources` (`id`);

ALTER TABLE `purchases` ADD CONSTRAINT `FK_Purchases_ClubMembers` FOREIGN KEY (`member_id`) REFERENCES `students` (`id`);