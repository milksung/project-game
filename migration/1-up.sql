CREATE TABLE
    `Admins` (
        `id` bigint PRIMARY KEY AUTO_INCREMENT,
        `username` varchar(255) DEFAULT NULL,
        `password` varchar(255) DEFAULT NULL,
        `role` ENUM ('SUPER_ADMIN', 'ADMIN') DEFAULT NULL,
        `status` ENUM ('ACTIVE', 'DEACTIVE') DEFAULT NULL,
        `firstname` varchar(255) DEFAULT NULL,
        `lastname` varchar(255) DEFAULT NULL,
        `fullname` varchar(255) DEFAULT NULL,
        `email` varchar(255) DEFAULT NULL,
        `phone` varchar(255) DEFAULT NULL,
        `ip` varchar(255) DEFAULT NULL,
        `admin_group_id` int DEFAULT NULL,
        `logedin_at` datetime DEFAULT NULL,
        `created_at` datetime DEFAULT (now()),
        `updated_at` datetime DEFAULT NULL,
        `deleted_at` datetime DEFAULT NULL
    );

INSERT INTO
    Admins (
        `username`,
        `password`,
        `role`,
        `status`,
        `firstname`,
        `lastname`
    )
VALUES (
        'superadmin',
        '$2a$12$aioRRrLa0bPY9IzlFd2kTeRKZo/..mgqi3BMcDCnsO3UevYVbxEbe',
        'SUPER_ADMIN',
        'active',
        'admin',
        'admin'
    );

CREATE TABLE
    `Admin_groups` (
        `id` bigint PRIMARY KEY AUTO_INCREMENT,
        `name` varchar(255) DEFAULT NULL,
        `admin_count` int DEFAULT 0,
        `created_at` datetime DEFAULT (now()),
        `updated_at` datetime DEFAULT NULL,
        `deleted_at` datetime DEFAULT NULL
    );

CREATE TABLE
    `Permissions` (
        `id` bigint PRIMARY KEY AUTO_INCREMENT,
        `permission_key` varchar(255) DEFAULT NULL,
        `name` varchar(255) DEFAULT NULL,
        `created_at` datetime DEFAULT (now()),
        `updated_at` datetime DEFAULT NULL,
        `deleted_at` datetime DEFAULT NULL
    );

CREATE TABLE
    `Admin_group_permissions` (
        `id` bigint PRIMARY KEY AUTO_INCREMENT,
        `group_id` bigint DEFAULT NULL,
        `permission_id` bigint DEFAULT NULL,
        `created_at` datetime DEFAULT (now()),
        `updated_at` datetime DEFAULT NULL,
        `deleted_at` datetime DEFAULT NULL
    );

CREATE TABLE
    `Admin_permissions` (
        `id` bigint PRIMARY KEY AUTO_INCREMENT,
        `admin_id` bigint DEFAULT NULL,
        `permission_id` bigint DEFAULT NULL,
        `created_at` datetime DEFAULT (now()),
        `updated_at` datetime DEFAULT NULL,
        `deleted_at` datetime DEFAULT NULL
    );

CREATE TABLE
    `Users` (
        `id` bigint PRIMARY KEY AUTO_INCREMENT,
        `member_code` varchar(255) NOT NULL,
        `username` varchar(255) NOT NULL,
        `phone` varchar(255) NOT NULL,
        `password` varchar(255) NOT NULL,
        `status` ENUM ('ACTIVE', 'DEACTIVE') DEFAULT NULL,
        `firstname` varchar(255) DEFAULT NULL,
        `lastname` varchar(255) DEFAULT NULL,
        `fullname` varchar(255) DEFAULT NULL,
        `credit` int DEFAULT 0,
        `ip` varchar(255) DEFAULT NULL,
        `logedin_at` datetime DEFAULT NULL,
        `created_at` datetime DEFAULT (now()),
        `updated_at` datetime DEFAULT NULL,
        `deleted_at` datetime DEFAULT NULL
    );

CREATE TABLE
    `User_login_logs` (
        `id` bigint PRIMARY KEY AUTO_INCREMENT,
        `user_id` bigint DEFAULT NULL,
        `ip` varchar(255) DEFAULT NULL,
        `created_at` datetime DEFAULT (now())
    );

CREATE TABLE
    `User_update_logs` (
        `id` bigint PRIMARY KEY AUTO_INCREMENT,
        `user_id` bigint DEFAULT NULL,
        `description` varchar(255) DEFAULT NULL,
        `created_by_username` varchar(255) DEFAULT NULL,
        `ip` varchar(255) DEFAULT NULL,
        `created_at` datetime DEFAULT (now())
    );

CREATE UNIQUE INDEX `uni_username` ON `Admins` (`username`);

CREATE INDEX `idx_status` ON `Admins` (`status`);

CREATE UNIQUE INDEX `uni_email` ON `Admins` (`email`);

CREATE INDEX `idx_admin_group_id` ON `Admins` (`admin_group_id`);

CREATE INDEX `idx_created_at` ON `Admins` (`created_at`);

CREATE INDEX `idx_deleted_at` ON `Admins` (`deleted_at`);

CREATE UNIQUE INDEX `uni_name` ON `Admin_groups` (`name`);

CREATE INDEX `idx_created_at` ON `Admin_groups` (`created_at`);

CREATE INDEX `idx_deleted_at` ON `Admin_groups` (`deleted_at`);

CREATE INDEX `idx_created_at` ON `Permissions` (`created_at`);

CREATE INDEX `idx_deleted_at` ON `Permissions` (`deleted_at`);

CREATE INDEX
    `idx_group_id` ON `Admin_group_permissions` (`group_id`);

CREATE INDEX
    `idx_permission_id` ON `Admin_group_permissions` (`permission_id`);

CREATE INDEX
    `idx_created_at` ON `Admin_group_permissions` (`created_at`);

CREATE INDEX
    `idx_deleted_at` ON `Admin_group_permissions` (`deleted_at`);

CREATE INDEX `idx_admin_id` ON `Admin_permissions` (`admin_id`);

CREATE INDEX `idx_created_at` ON `Admin_permissions` (`created_at`);

CREATE INDEX `idx_deleted_at` ON `Admin_permissions` (`deleted_at`);

CREATE UNIQUE INDEX `uni_username` ON `Users` (`username`);

CREATE INDEX `idx_status` ON `Users` (`status`);

CREATE INDEX `idx_created_at` ON `Users` (`created_at`);

CREATE INDEX `idx_deleted_at` ON `Users` (`deleted_at`);

CREATE INDEX `idx_user_id` ON `User_login_logs` (`user_id`);

CREATE INDEX `idx_created_at` ON `User_login_logs` (`created_at`);

CREATE INDEX `idx_user_id` ON `User_update_logs` (`user_id`);

CREATE INDEX `idx_created_at` ON `User_update_logs` (`created_at`);