CREATE TABLE IF NOT EXISTS videos (
    id SERIAL PRIMARY KEY,
    center_id INT NOT NULL,
    title VARCHAR(255) NOT NULL UNIQUE,
    short_text TEXT NOT NULL,
    video_url TEXT NOT NULL,

    createAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_center FOREIGN KEY (center_id) REFERENCES centers(id) ON DELETE CASCADE
);
