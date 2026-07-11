CREATE TABLE IF NOT EXISTS posts (
    id           INT AUTO_INCREMENT PRIMARY KEY,
    title        VARCHAR(200)  NOT NULL,
    content      TEXT          NOT NULL,
    category     VARCHAR(100)  NOT NULL,
    created_date TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_date TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    status       VARCHAR(100)  NOT NULL COMMENT 'publish | draft | thrash'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;
