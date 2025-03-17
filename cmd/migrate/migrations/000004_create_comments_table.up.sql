CREATE TABLE IF NOT EXISTS comments (
  id BIGSERIAL PRIMARY KEY,
  content TEXT NOT NULL,
  likes INTEGER NOT NULL DEFAULT 0,
  dislikes INTEGER NOT NULL DEFAULT 0,
  manga_id BIGINT NOT NULL,
  chapter_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE comments ADD CONSTRAINT fk_comments_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
CREATE INDEX idx_comments_manga_id ON comments (manga_id);
CREATE INDEX idx_comments_chapter_id ON comments (chapter_id);
