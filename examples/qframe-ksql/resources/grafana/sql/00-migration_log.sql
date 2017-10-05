PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE `migration_log` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `migration_id` TEXT NOT NULL
, `sql` TEXT NOT NULL
, `success` INTEGER NOT NULL
, `error` TEXT NOT NULL
, `timestamp` DATETIME NOT NULL
);
INSERT INTO "migration_log" VALUES(1,'create migration_log table','CREATE TABLE IF NOT EXISTS `migration_log` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `migration_id` TEXT NOT NULL
, `sql` TEXT NOT NULL
, `success` INTEGER NOT NULL
, `error` TEXT NOT NULL
, `timestamp` DATETIME NOT NULL
);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(2,'create user table','CREATE TABLE IF NOT EXISTS `user` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `version` INTEGER NOT NULL
, `login` TEXT NOT NULL
, `email` TEXT NOT NULL
, `name` TEXT NULL
, `password` TEXT NULL
, `salt` TEXT NULL
, `rands` TEXT NULL
, `company` TEXT NULL
, `account_id` INTEGER NOT NULL
, `is_admin` INTEGER NOT NULL
, `created` DATETIME NOT NULL
, `updated` DATETIME NOT NULL
);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(3,'add unique index user.login','CREATE UNIQUE INDEX `UQE_user_login` ON `user` (`login`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(4,'add unique index user.email','CREATE UNIQUE INDEX `UQE_user_email` ON `user` (`email`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(5,'drop index UQE_user_login - v1','DROP INDEX `UQE_user_login`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(6,'drop index UQE_user_email - v1','DROP INDEX `UQE_user_email`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(7,'Rename table user to user_v1 - v1','ALTER TABLE `user` RENAME TO `user_v1`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(8,'create user table v2','CREATE TABLE IF NOT EXISTS `user` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `version` INTEGER NOT NULL
, `login` TEXT NOT NULL
, `email` TEXT NOT NULL
, `name` TEXT NULL
, `password` TEXT NULL
, `salt` TEXT NULL
, `rands` TEXT NULL
, `company` TEXT NULL
, `org_id` INTEGER NOT NULL
, `is_admin` INTEGER NOT NULL
, `email_verified` INTEGER NULL
, `theme` TEXT NULL
, `created` DATETIME NOT NULL
, `updated` DATETIME NOT NULL
);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(9,'create index UQE_user_login - v2','CREATE UNIQUE INDEX `UQE_user_login` ON `user` (`login`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(10,'create index UQE_user_email - v2','CREATE UNIQUE INDEX `UQE_user_email` ON `user` (`email`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(11,'copy data_source v1 to v2','INSERT INTO `user` (`login`
, `password`
, `salt`
, `org_id`
, `created`
, `id`
, `email`
, `name`
, `rands`
, `company`
, `is_admin`
, `updated`
, `version`) SELECT `login`
, `password`
, `salt`
, `account_id`
, `created`
, `id`
, `email`
, `name`
, `rands`
, `company`
, `is_admin`
, `updated`
, `version` FROM `user_v1`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(12,'Drop old table user_v1','DROP TABLE IF EXISTS `user_v1`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(13,'create temp user table v1-7','CREATE TABLE IF NOT EXISTS `temp_user` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `org_id` INTEGER NOT NULL
, `version` INTEGER NOT NULL
, `email` TEXT NOT NULL
, `name` TEXT NULL
, `role` TEXT NULL
, `code` TEXT NOT NULL
, `status` TEXT NOT NULL
, `invited_by_user_id` INTEGER NULL
, `email_sent` INTEGER NOT NULL
, `email_sent_on` DATETIME NULL
, `remote_addr` TEXT NULL
, `created` DATETIME NOT NULL
, `updated` DATETIME NOT NULL
);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(14,'create index IDX_temp_user_email - v1-7','CREATE INDEX `IDX_temp_user_email` ON `temp_user` (`email`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(15,'create index IDX_temp_user_org_id - v1-7','CREATE INDEX `IDX_temp_user_org_id` ON `temp_user` (`org_id`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(16,'create index IDX_temp_user_code - v1-7','CREATE INDEX `IDX_temp_user_code` ON `temp_user` (`code`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(17,'create index IDX_temp_user_status - v1-7','CREATE INDEX `IDX_temp_user_status` ON `temp_user` (`status`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(18,'create star table','CREATE TABLE IF NOT EXISTS `star` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `user_id` INTEGER NOT NULL
, `dashboard_id` INTEGER NOT NULL
);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(19,'add unique index star.user_id_dashboard_id','CREATE UNIQUE INDEX `UQE_star_user_id_dashboard_id` ON `star` (`user_id`,`dashboard_id`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(20,'create org table v1','CREATE TABLE IF NOT EXISTS `org` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `version` INTEGER NOT NULL
, `name` TEXT NOT NULL
, `address1` TEXT NULL
, `address2` TEXT NULL
, `city` TEXT NULL
, `state` TEXT NULL
, `zip_code` TEXT NULL
, `country` TEXT NULL
, `billing_email` TEXT NULL
, `created` DATETIME NOT NULL
, `updated` DATETIME NOT NULL
);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(21,'create index UQE_org_name - v1','CREATE UNIQUE INDEX `UQE_org_name` ON `org` (`name`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(22,'create org_user table v1','CREATE TABLE IF NOT EXISTS `org_user` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `org_id` INTEGER NOT NULL
, `user_id` INTEGER NOT NULL
, `role` TEXT NOT NULL
, `created` DATETIME NOT NULL
, `updated` DATETIME NOT NULL
);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(23,'create index IDX_org_user_org_id - v1','CREATE INDEX `IDX_org_user_org_id` ON `org_user` (`org_id`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(24,'create index UQE_org_user_org_id_user_id - v1','CREATE UNIQUE INDEX `UQE_org_user_org_id_user_id` ON `org_user` (`org_id`,`user_id`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(25,'copy data account to org','INSERT INTO `org` (`id`
, `version`
, `name`
, `created`
, `updated`) SELECT `id`
, `version`
, `name`
, `created`
, `updated` FROM `account`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(26,'copy data account_user to org_user','INSERT INTO `org_user` (`id`
, `org_id`
, `user_id`
, `role`
, `created`
, `updated`) SELECT `id`
, `account_id`
, `user_id`
, `role`
, `created`
, `updated` FROM `account_user`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(27,'Drop old table account','DROP TABLE IF EXISTS `account`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(28,'Drop old table account_user','DROP TABLE IF EXISTS `account_user`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(29,'create dashboard table','CREATE TABLE IF NOT EXISTS `dashboard` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `version` INTEGER NOT NULL
, `slug` TEXT NOT NULL
, `title` TEXT NOT NULL
, `data` TEXT NOT NULL
, `account_id` INTEGER NOT NULL
, `created` DATETIME NOT NULL
, `updated` DATETIME NOT NULL
);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(30,'add index dashboard.account_id','CREATE INDEX `IDX_dashboard_account_id` ON `dashboard` (`account_id`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(31,'add unique index dashboard_account_id_slug','CREATE UNIQUE INDEX `UQE_dashboard_account_id_slug` ON `dashboard` (`account_id`,`slug`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(32,'create dashboard_tag table','CREATE TABLE IF NOT EXISTS `dashboard_tag` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `dashboard_id` INTEGER NOT NULL
, `term` TEXT NOT NULL
);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(33,'add unique index dashboard_tag.dasboard_id_term','CREATE UNIQUE INDEX `UQE_dashboard_tag_dashboard_id_term` ON `dashboard_tag` (`dashboard_id`,`term`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(34,'drop index UQE_dashboard_tag_dashboard_id_term - v1','DROP INDEX `UQE_dashboard_tag_dashboard_id_term`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(35,'Rename table dashboard to dashboard_v1 - v1','ALTER TABLE `dashboard` RENAME TO `dashboard_v1`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(36,'create dashboard v2','CREATE TABLE IF NOT EXISTS `dashboard` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `version` INTEGER NOT NULL
, `slug` TEXT NOT NULL
, `title` TEXT NOT NULL
, `data` TEXT NOT NULL
, `org_id` INTEGER NOT NULL
, `created` DATETIME NOT NULL
, `updated` DATETIME NOT NULL
);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(37,'create index IDX_dashboard_org_id - v2','CREATE INDEX `IDX_dashboard_org_id` ON `dashboard` (`org_id`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(38,'create index UQE_dashboard_org_id_slug - v2','CREATE UNIQUE INDEX `UQE_dashboard_org_id_slug` ON `dashboard` (`org_id`,`slug`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(39,'copy dashboard v1 to v2','INSERT INTO `dashboard` (`id`
, `version`
, `slug`
, `title`
, `data`
, `org_id`
, `created`
, `updated`) SELECT `id`
, `version`
, `slug`
, `title`
, `data`
, `account_id`
, `created`
, `updated` FROM `dashboard_v1`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(40,'drop table dashboard_v1','DROP TABLE IF EXISTS `dashboard_v1`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(41,'alter dashboard.data to mediumtext v1','SELECT 0 WHERE 0;',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(42,'Add column updated_by in dashboard - v2','alter table `dashboard` ADD COLUMN `updated_by` INTEGER NULL ',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(43,'Add column created_by in dashboard - v2','alter table `dashboard` ADD COLUMN `created_by` INTEGER NULL ',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(44,'create data_source table','CREATE TABLE IF NOT EXISTS `data_source` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `account_id` INTEGER NOT NULL
, `version` INTEGER NOT NULL
, `type` TEXT NOT NULL
, `name` TEXT NOT NULL
, `access` TEXT NOT NULL
, `url` TEXT NOT NULL
, `password` TEXT NULL
, `user` TEXT NULL
, `database` TEXT NULL
, `basic_auth` INTEGER NOT NULL
, `basic_auth_user` TEXT NULL
, `basic_auth_password` TEXT NULL
, `is_default` INTEGER NOT NULL
, `created` DATETIME NOT NULL
, `updated` DATETIME NOT NULL
);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(45,'add index data_source.account_id','CREATE INDEX `IDX_data_source_account_id` ON `data_source` (`account_id`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(46,'add unique index data_source.account_id_name','CREATE UNIQUE INDEX `UQE_data_source_account_id_name` ON `data_source` (`account_id`,`name`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(47,'drop index IDX_data_source_account_id - v1','DROP INDEX `IDX_data_source_account_id`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(48,'drop index UQE_data_source_account_id_name - v1','DROP INDEX `UQE_data_source_account_id_name`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(49,'Rename table data_source to data_source_v1 - v1','ALTER TABLE `data_source` RENAME TO `data_source_v1`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(50,'create data_source table v2','CREATE TABLE IF NOT EXISTS `data_source` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `org_id` INTEGER NOT NULL
, `version` INTEGER NOT NULL
, `type` TEXT NOT NULL
, `name` TEXT NOT NULL
, `access` TEXT NOT NULL
, `url` TEXT NOT NULL
, `password` TEXT NULL
, `user` TEXT NULL
, `database` TEXT NULL
, `basic_auth` INTEGER NOT NULL
, `basic_auth_user` TEXT NULL
, `basic_auth_password` TEXT NULL
, `is_default` INTEGER NOT NULL
, `json_data` TEXT NULL
, `created` DATETIME NOT NULL
, `updated` DATETIME NOT NULL
);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(51,'create index IDX_data_source_org_id - v2','CREATE INDEX `IDX_data_source_org_id` ON `data_source` (`org_id`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(52,'create index UQE_data_source_org_id_name - v2','CREATE UNIQUE INDEX `UQE_data_source_org_id_name` ON `data_source` (`org_id`,`name`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(53,'copy data_source v1 to v2','INSERT INTO `data_source` (`password`
, `basic_auth`
, `type`
, `name`
, `access`
, `user`
, `org_id`
, `basic_auth_user`
, `basic_auth_password`
, `database`
, `is_default`
, `updated`
, `id`
, `version`
, `url`
, `created`) SELECT `password`
, `basic_auth`
, `type`
, `name`
, `access`
, `user`
, `account_id`
, `basic_auth_user`
, `basic_auth_password`
, `database`
, `is_default`
, `updated`
, `id`
, `version`
, `url`
, `created` FROM `data_source_v1`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(54,'Drop old table data_source_v1 #2','DROP TABLE IF EXISTS `data_source_v1`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(55,'Add column with_credentials','alter table `data_source` ADD COLUMN `with_credentials` INTEGER NOT NULL DEFAULT 0 ',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(56,'create api_key table','CREATE TABLE IF NOT EXISTS `api_key` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `account_id` INTEGER NOT NULL
, `name` TEXT NOT NULL
, `key` TEXT NOT NULL
, `role` TEXT NOT NULL
, `created` DATETIME NOT NULL
, `updated` DATETIME NOT NULL
);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(57,'add index api_key.account_id','CREATE INDEX `IDX_api_key_account_id` ON `api_key` (`account_id`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(58,'add index api_key.key','CREATE UNIQUE INDEX `UQE_api_key_key` ON `api_key` (`key`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(59,'add index api_key.account_id_name','CREATE UNIQUE INDEX `UQE_api_key_account_id_name` ON `api_key` (`account_id`,`name`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(60,'drop index IDX_api_key_account_id - v1','DROP INDEX `IDX_api_key_account_id`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(61,'drop index UQE_api_key_key - v1','DROP INDEX `UQE_api_key_key`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(62,'drop index UQE_api_key_account_id_name - v1','DROP INDEX `UQE_api_key_account_id_name`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(63,'Rename table api_key to api_key_v1 - v1','ALTER TABLE `api_key` RENAME TO `api_key_v1`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(64,'create api_key table v2','CREATE TABLE IF NOT EXISTS `api_key` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `org_id` INTEGER NOT NULL
, `name` TEXT NOT NULL
, `key` TEXT NOT NULL
, `role` TEXT NOT NULL
, `created` DATETIME NOT NULL
, `updated` DATETIME NOT NULL
);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(65,'create index IDX_api_key_org_id - v2','CREATE INDEX `IDX_api_key_org_id` ON `api_key` (`org_id`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(66,'create index UQE_api_key_key - v2','CREATE UNIQUE INDEX `UQE_api_key_key` ON `api_key` (`key`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(67,'create index UQE_api_key_org_id_name - v2','CREATE UNIQUE INDEX `UQE_api_key_org_id_name` ON `api_key` (`org_id`,`name`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(68,'copy api_key v1 to v2','INSERT INTO `api_key` (`name`
, `key`
, `role`
, `created`
, `updated`
, `id`
, `org_id`) SELECT `name`
, `key`
, `role`
, `created`
, `updated`
, `id`
, `account_id` FROM `api_key_v1`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(69,'Drop old table api_key_v1','DROP TABLE IF EXISTS `api_key_v1`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(70,'create dashboard_snapshot table v4','CREATE TABLE IF NOT EXISTS `dashboard_snapshot` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `name` TEXT NOT NULL
, `key` TEXT NOT NULL
, `dashboard` TEXT NOT NULL
, `expires` DATETIME NOT NULL
, `created` DATETIME NOT NULL
, `updated` DATETIME NOT NULL
);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(71,'drop table dashboard_snapshot_v4 #1','DROP TABLE IF EXISTS `dashboard_snapshot`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(72,'create dashboard_snapshot table v5 #2','CREATE TABLE IF NOT EXISTS `dashboard_snapshot` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `name` TEXT NOT NULL
, `key` TEXT NOT NULL
, `delete_key` TEXT NOT NULL
, `org_id` INTEGER NOT NULL
, `user_id` INTEGER NOT NULL
, `external` INTEGER NOT NULL
, `external_url` TEXT NOT NULL
, `dashboard` TEXT NOT NULL
, `expires` DATETIME NOT NULL
, `created` DATETIME NOT NULL
, `updated` DATETIME NOT NULL
);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(73,'create index UQE_dashboard_snapshot_key - v5','CREATE UNIQUE INDEX `UQE_dashboard_snapshot_key` ON `dashboard_snapshot` (`key`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(74,'create index UQE_dashboard_snapshot_delete_key - v5','CREATE UNIQUE INDEX `UQE_dashboard_snapshot_delete_key` ON `dashboard_snapshot` (`delete_key`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(75,'create index IDX_dashboard_snapshot_user_id - v5','CREATE INDEX `IDX_dashboard_snapshot_user_id` ON `dashboard_snapshot` (`user_id`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(76,'alter dashboard_snapshot to mediumtext v2','SELECT 0 WHERE 0;',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(77,'create quota table v1','CREATE TABLE IF NOT EXISTS `quota` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `org_id` INTEGER NULL
, `user_id` INTEGER NULL
, `target` TEXT NOT NULL
, `limit` INTEGER NOT NULL
, `created` DATETIME NOT NULL
, `updated` DATETIME NOT NULL
);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(78,'create index UQE_quota_org_id_user_id_target - v1','CREATE UNIQUE INDEX `UQE_quota_org_id_user_id_target` ON `quota` (`org_id`,`user_id`,`target`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(79,'create plugin_setting table','CREATE TABLE IF NOT EXISTS `plugin_setting` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `org_id` INTEGER NULL
, `plugin_id` TEXT NOT NULL
, `enabled` INTEGER NOT NULL
, `pinned` INTEGER NOT NULL
, `json_data` TEXT NULL
, `secure_json_data` TEXT NULL
, `created` DATETIME NOT NULL
, `updated` DATETIME NOT NULL
);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(80,'create index UQE_plugin_setting_org_id_plugin_id - v1','CREATE UNIQUE INDEX `UQE_plugin_setting_org_id_plugin_id` ON `plugin_setting` (`org_id`,`plugin_id`);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(81,'create session table','CREATE TABLE IF NOT EXISTS `session` (
`key` TEXT PRIMARY KEY NOT NULL
, `data` BLOB NOT NULL
, `expiry` INTEGER NOT NULL
);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(82,'Drop old table playlist table','DROP TABLE IF EXISTS `playlist`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(83,'Drop old table playlist_item table','DROP TABLE IF EXISTS `playlist_item`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(84,'create playlist table v2','CREATE TABLE IF NOT EXISTS `playlist` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `name` TEXT NOT NULL
, `interval` TEXT NOT NULL
, `org_id` INTEGER NOT NULL
);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(85,'create playlist item table v2','CREATE TABLE IF NOT EXISTS `playlist_item` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `playlist_id` INTEGER NOT NULL
, `type` TEXT NOT NULL
, `value` TEXT NOT NULL
, `title` TEXT NOT NULL
, `order` INTEGER NOT NULL
);',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(86,'drop preferences table v2','DROP TABLE IF EXISTS `preferences`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(87,'drop preferences table v3','DROP TABLE IF EXISTS `preferences`',1,'','2016-05-14 11:07:42');
INSERT INTO "migration_log" VALUES(88,'create preferences table v3','CREATE TABLE IF NOT EXISTS `preferences` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `org_id` INTEGER NOT NULL
, `user_id` INTEGER NOT NULL
, `version` INTEGER NOT NULL
, `home_dashboard_id` INTEGER NOT NULL
, `timezone` TEXT NOT NULL
, `theme` TEXT NOT NULL
, `created` DATETIME NOT NULL
, `updated` DATETIME NOT NULL
);',1,'','2016-05-14 11:07:43');
INSERT INTO "migration_log" VALUES(89,'Add column gnetId in dashboard','alter table `dashboard` ADD COLUMN `gnet_id` INTEGER NULL ',1,'','2016-09-13 23:22:09');
INSERT INTO "migration_log" VALUES(90,'Add index for gnetId in dashboard','CREATE INDEX `IDX_dashboard_gnet_id` ON `dashboard` (`gnet_id`);',1,'','2016-09-13 23:22:09');
INSERT INTO "migration_log" VALUES(91,'Add column plugin_id in dashboard','alter table `dashboard` ADD COLUMN `plugin_id` TEXT NULL ',1,'','2016-09-13 23:22:09');
INSERT INTO "migration_log" VALUES(92,'Add index for plugin_id in dashboard','CREATE INDEX `IDX_dashboard_org_id_plugin_id` ON `dashboard` (`org_id`,`plugin_id`);',1,'','2016-09-13 23:22:09');
INSERT INTO "migration_log" VALUES(93,'Add column plugin_version to plugin_settings','alter table `plugin_setting` ADD COLUMN `plugin_version` TEXT NULL ',1,'','2016-09-13 23:22:09');
COMMIT;
