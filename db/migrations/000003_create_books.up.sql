CREATE TABLE `books` (
  `uuid` VARBINARY(36) NOT NULL,
  `title` TEXT NOT NULL,
  `publisher_uuid` VARBINARY(36) NOT NULL,
  PRIMARY KEY (`uuid`),
  FOREIGN KEY (`publisher_uuid`) REFERENCES `publishers` (`uuid`)
);
