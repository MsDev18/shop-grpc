CREATE TABLE
    `product` (
        `id` INT NOT NULL UNIQUE PRIMARY KEY AUTO_INCREMENT,
        `name` VARCHAR(255) NOT NULL,
        `slug` VARCHAR(255) NOT NULL UNIQUE,
        `description` TEXT NOT NULL,
        `price` INT NOT NULL,
        `stock` INT NOT NULL DEFAULT 0,
        `main_image` VARCHAR(255) NOT NULL,
        `category_id` INT NOT NULL,
        `deleted_at` TIMESTAMP NULL,
        `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        FOREIGN KEY (`category_id`) REFERENCES `category` (`id`)
    );