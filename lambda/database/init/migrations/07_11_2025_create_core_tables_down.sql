ALTER TABLE `purchases` DROP FOREIGN KEY `FK_Purchases_ClubMembers`;

ALTER TABLE `point_histories` DROP FOREIGN KEY `FK_PointHistories_ClubMembers`;

ALTER TABLE `point_histories` DROP FOREIGN KEY `FK_PointHistories_PointSources`;

ALTER TABLE `events_to_clubs` DROP FOREIGN KEY `FK_EventsToClubs_Events`;

ALTER TABLE `events_to_clubs` DROP FOREIGN KEY `FK_EventsToClubs_Clubs`;

ALTER TABLE `event_versions` DROP FOREIGN KEY `FK_EventVersions_ClubMembers`;

ALTER TABLE `events` DROP FOREIGN KEY `FK_Events_CurrentVersion`;

ALTER TABLE `events` DROP FOREIGN KEY `FK_Events_OriginalVersion`;

ALTER TABLE `club_members` DROP FOREIGN KEY `FK_ClubMembers_Students`;

ALTER TABLE `club_members` DROP FOREIGN KEY `FK_ClubMembers_Clubs`;

ALTER TABLE `students_to_interests` DROP FOREIGN KEY `FK_StudentsToInterests_Students`;

ALTER TABLE `students_to_interests` DROP FOREIGN KEY `FK_StudentsToInterests_Interests`;

ALTER TABLE `student_info` DROP FOREIGN KEY `FK_StudentInfo_Students`;

DROP TABLE IF EXISTS `PURCHASES`;

DROP TABLE IF EXISTS `POINT_HISTORIES`;

DROP TABLE IF EXISTS `POINT_SOURCES`;

DROP TABLE IF EXISTS `EVENTS_TO_CLUBS`;

DROP TABLE IF EXISTS `EVENT_VERSIONS`;

DROP TABLE IF EXISTS `EVENTS`;

DROP TABLE IF EXISTS `CLUB_MEMBERS`;

DROP TABLE IF EXISTS `CLUBS`;

DROP TABLE IF EXISTS `INTERESTS`;

DROP TABLE IF EXISTS `STUDENTS_TO_INTERESTS`;

DROP TABLE IF EXISTS `STUDENT_INFO`;

DROP TABLE IF EXISTS `STUDENTS`;