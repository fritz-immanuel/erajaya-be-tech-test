CREATE TABLE user_actions (
  id VARCHAR(255) NOT NULL,
  user_id VARCHAR(255) DEFAULT '',
  table_name VARCHAR(200) DEFAULT NULL,
  action VARCHAR(100) DEFAULT NULL,
  action_value INT DEFAULT 0,
  ref_id VARCHAR(255) DEFAULT '0',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  INDEX idx_user_id (user_id),
  INDEX idx_ref_id (ref_id),
  INDEX idx_table_name (table_name)
);