CREATE TABLE `items` (
    `id` INT(8) unsigned NOT NULL AUTO_INCREMENT,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `hash` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`)
);

ALTER TABLE `items` ADD INDEX idx_users_username(name);
