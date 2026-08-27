CREATE TABLE `otp` (
    `id` INT PRIMARY KEY NOT NULL UNIQUE AUTO_INCREMENT,
    `user_id` INT NOT NULL UNIQUE ,
    `code` VARCHAR(10) NOT NULL ,
    `expires_at` TIMESTAMP NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`)
);