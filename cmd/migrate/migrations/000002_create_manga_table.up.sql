CREATE TABLE IF NOT EXISTS mangas (
  id BIGSERIAL PRIMARY KEY,
  title varchar(255) NOT NULL,
  description TEXT,
  slug VARCHAR(255) NOT NULL,
  cover VARCHAR(255),
  tags VARCHAR(100) [],
  type SMALLINT DEFAULT 0,
  status SMALLINT DEFAULT 0,
  painter VARCHAR(255),
  user_id BIGSERIAL NOT NULL,
  version INT DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT unique_slug UNIQUE(slug),
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)
);