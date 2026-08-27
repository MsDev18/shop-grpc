CREATE TABLE `session` (
    `id` INT PRIMARY KEY NOT NULL UNIQUE AUTO_INCREMENT,
    `user_id` INT NOT NULL ,
    `expires_at` TIMESTAMP NOT NULL ,
    `revoke_at` TIMESTAMP NULL,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`)
);