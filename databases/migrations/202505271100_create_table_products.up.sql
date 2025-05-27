CREATE TABLE products (
  id VARCHAR(255) NOT NULL,
  name VARCHAR(255) DEFAULT '',
  price DECIMAL(25,2) DEFAULT 0,
  description LONGTEXT NOT NULL,
  quantity INT DEFAULT 0,

  status_id VARCHAR(5) DEFAULT "1",
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  created_by VARCHAR(255) DEFAULT '',
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_by VARCHAR(255) DEFAULT '',
  PRIMARY KEY (id),
  INDEX idx_products_name (name)
);