CREATE TABLE `sessions` (
  `id` char(36) PRIMARY KEY,
  `username` varchar(255) NOT NULL,
  `refresh_token` text NOT NULL,
  `user_agent` varchar(255) NOT NULL,
  `client_ip` varchar(255) NOT NULL,
  `is_blocked` tinyint(1) NOT NULL DEFAULT 0,
  `expires_at` datetime NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);
