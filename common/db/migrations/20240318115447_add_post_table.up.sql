CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    post_in_html TEXT NOT NULL,
)

CREATE TABLE IF NOT EXISTS post_tags (
    post_id INT REFERENCES posts(id),
    tag VARCHAR(50) NOT NULL,
);

CREATE INDEX IF NOT EXISTS post_id  ON post_tags(post_id);
