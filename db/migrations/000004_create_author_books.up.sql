CREATE TABLE `author_books` (
  `author_id` BIGINT  NOT NULL,
  `book_id`   BIGINT  NOT NULL,
  PRIMARY KEY (`author_id`, `book_id`),
  FOREIGN KEY (`author_id`) REFERENCES `authors` (`id`),
  FOREIGN KEY (`book_id`)   REFERENCES `books` (`id`)
);
