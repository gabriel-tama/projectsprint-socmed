CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    post_in_html TEXT NOT NULL,
    user_id INT REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS post_tags (
    post_id INT REFERENCES posts(id),
    tag VARCHAR(50) NOT NULL
);

CREATE INDEX IF NOT EXISTS post_id_idx  ON post_tags(post_id);
CREATE INDEX IF NOT EXISTS user_id_idx  ON posts(user_id);

