CREATE TABLE `author_books` (
  `author_uuid` VARBINARY(36) NOT NULL,
  `book_uuid` VARBINARY(36) NOT NULL,
  PRIMARY KEY (`author_uuid`, `book_uuid`),
  FOREIGN KEY (`author_uuid`) REFERENCES `authors` (`uuid`),
  FOREIGN KEY (`book_uuid`) REFERENCES `books` (`uuid`)
);
