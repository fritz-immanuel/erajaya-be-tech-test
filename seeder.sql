CREATE DATABASE IF NOT EXISTS erajaya_be_tech_test;
USE erajaya_be_tech_test;

CREATE TABLE IF NOT EXISTS schema_migrations (
  version BIGINT NOT NULL PRIMARY KEY,
  dirty TINYINT(1) NOT NULL
);
INSERT INTO schema_migrations (version, dirty)
VALUES (202410251400, 0);

CREATE TABLE status (
  id VARCHAR(5) NOT NULL,
  name VARCHAR(255) DEFAULT "",
  PRIMARY KEY (id)
);