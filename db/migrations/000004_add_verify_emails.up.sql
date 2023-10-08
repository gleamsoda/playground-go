CREATE TABLE `verify_emails` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `username` varchar(255) NOT NULL,
  `email` varchar(255) UNIQUE NOT NULL,
  `secret_code` varchar(32) NOT NULL,
  `is_used` tinyint(1) NOT NULL DEFAULT 0,
  `expired_at` datetime NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE `users` ADD COLUMN `is_email_verified` tinyint(1) NOT NULL DEFAULT 0;
