CREATE TABLE
    `address` (
        `id` INT NOT NULL UNIQUE PRIMARY KEY AUTO_INCREMENT,
        `user_id` INT NOT NULL,
        `title` VARCHAR(255) NOT NULL,
        `province_id` INT NOT NULL,
        `city` VARCHAR(255) NOT NULL,
        `address` TEXT NOT NULL,
        `postal_code` VARCHAR(255) NOT NULL,
        `deleted_at` TIMESTAMP NULL,
        `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        FOREIGN KEY (`user_id`) REFERENCES `user` (`id`),
        FOREIGN KEY (`province_id`) REFERENCES `province` (`id`)
    )