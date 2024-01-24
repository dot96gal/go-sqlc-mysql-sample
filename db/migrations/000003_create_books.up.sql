CREATE TABLE `books` (
  `id`           BIGINT  NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `title`        TEXT    NOT NULL,
  `publisher_id` BIGINT  NOT NULL,
  FOREIGN KEY (`publisher_id`) REFERENCES `publishers` (`id`)
);
