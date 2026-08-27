CREATE TABLE
    `product_image` (
        `id` INT NOT NULL UNIQUE PRIMARY KEY AUTO_INCREMENT,
        `image` VARCHAR(255) NOT NULL,
        `product_id` INT NOT NULL,
        `deleted_at` TIMESTAMP NULL,
        `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (`product_id`) REFERENCES `product` (`id`)
    );