PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE `api_key` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `org_id` INTEGER NOT NULL
, `name` TEXT NOT NULL
, `key` TEXT NOT NULL
, `role` TEXT NOT NULL
, `created` DATETIME NOT NULL
, `updated` DATETIME NOT NULL
);
CREATE TABLE `dashboard_snapshot` (
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
);
CREATE TABLE `quota` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `org_id` INTEGER NULL
, `user_id` INTEGER NULL
, `target` TEXT NOT NULL
, `limit` INTEGER NOT NULL
, `created` DATETIME NOT NULL
, `updated` DATETIME NOT NULL
);
CREATE TABLE `plugin_setting` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `org_id` INTEGER NULL
, `plugin_id` TEXT NOT NULL
, `enabled` INTEGER NOT NULL
, `pinned` INTEGER NOT NULL
, `json_data` TEXT NULL
, `secure_json_data` TEXT NULL
, `created` DATETIME NOT NULL
, `updated` DATETIME NOT NULL
, `plugin_version` TEXT NULL);
CREATE TABLE `session` (
`key` TEXT PRIMARY KEY NOT NULL
, `data` BLOB NOT NULL
, `expiry` INTEGER NOT NULL
);
CREATE TABLE `playlist` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `name` TEXT NOT NULL
, `interval` TEXT NOT NULL
, `org_id` INTEGER NOT NULL
);
CREATE TABLE `playlist_item` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `playlist_id` INTEGER NOT NULL
, `type` TEXT NOT NULL
, `value` TEXT NOT NULL
, `title` TEXT NOT NULL
, `order` INTEGER NOT NULL
);
CREATE TABLE `preferences` (
`id` INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL
, `org_id` INTEGER NOT NULL
, `user_id` INTEGER NOT NULL
, `version` INTEGER NOT NULL
, `home_dashboard_id` INTEGER NOT NULL
, `timezone` TEXT NOT NULL
, `theme` TEXT NOT NULL
, `created` DATETIME NOT NULL
, `updated` DATETIME NOT NULL
);
INSERT INTO "preferences" VALUES(1,1,0,0,0,'','light','2016-05-16 17:21:04','2016-05-16 17:21:04');
DELETE FROM sqlite_sequence;
INSERT INTO "sqlite_sequence" VALUES('migration_log',93);
INSERT INTO "sqlite_sequence" VALUES('user',1);
INSERT INTO "sqlite_sequence" VALUES('dashboard',4);
INSERT INTO "sqlite_sequence" VALUES('data_source',2);
INSERT INTO "sqlite_sequence" VALUES('api_key',0);
INSERT INTO "sqlite_sequence" VALUES('org',1);
INSERT INTO "sqlite_sequence" VALUES('org_user',1);
INSERT INTO "sqlite_sequence" VALUES('preferences',1);
INSERT INTO "sqlite_sequence" VALUES('star',1);
CREATE UNIQUE INDEX `UQE_user_login` ON `user` (`login`);
CREATE UNIQUE INDEX `UQE_user_email` ON `user` (`email`);
CREATE INDEX `IDX_temp_user_email` ON `temp_user` (`email`);
CREATE INDEX `IDX_temp_user_org_id` ON `temp_user` (`org_id`);
CREATE INDEX `IDX_temp_user_code` ON `temp_user` (`code`);
CREATE INDEX `IDX_temp_user_status` ON `temp_user` (`status`);
CREATE UNIQUE INDEX `UQE_star_user_id_dashboard_id` ON `star` (`user_id`,`dashboard_id`);
CREATE UNIQUE INDEX `UQE_org_name` ON `org` (`name`);
CREATE INDEX `IDX_org_user_org_id` ON `org_user` (`org_id`);
CREATE UNIQUE INDEX `UQE_org_user_org_id_user_id` ON `org_user` (`org_id`,`user_id`);
CREATE INDEX `IDX_dashboard_org_id` ON `dashboard` (`org_id`);
CREATE UNIQUE INDEX `UQE_dashboard_org_id_slug` ON `dashboard` (`org_id`,`slug`);
CREATE INDEX `IDX_data_source_org_id` ON `data_source` (`org_id`);
CREATE UNIQUE INDEX `UQE_data_source_org_id_name` ON `data_source` (`org_id`,`name`);
CREATE INDEX `IDX_api_key_org_id` ON `api_key` (`org_id`);
CREATE UNIQUE INDEX `UQE_api_key_key` ON `api_key` (`key`);
CREATE UNIQUE INDEX `UQE_api_key_org_id_name` ON `api_key` (`org_id`,`name`);
CREATE UNIQUE INDEX `UQE_dashboard_snapshot_key` ON `dashboard_snapshot` (`key`);
CREATE UNIQUE INDEX `UQE_dashboard_snapshot_delete_key` ON `dashboard_snapshot` (`delete_key`);
CREATE INDEX `IDX_dashboard_snapshot_user_id` ON `dashboard_snapshot` (`user_id`);
CREATE UNIQUE INDEX `UQE_quota_org_id_user_id_target` ON `quota` (`org_id`,`user_id`,`target`);
CREATE UNIQUE INDEX `UQE_plugin_setting_org_id_plugin_id` ON `plugin_setting` (`org_id`,`plugin_id`);
CREATE INDEX `IDX_dashboard_gnet_id` ON `dashboard` (`gnet_id`);
CREATE INDEX `IDX_dashboard_org_id_plugin_id` ON `dashboard` (`org_id`,`plugin_id`);
COMMIT;
