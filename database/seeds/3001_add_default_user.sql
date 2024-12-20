USE sb_db;

INSERT INTO users (email, password, created_at, updated_at) VALUES ('admin@example.com', '$2a$10$h0GW3t/zIv2dv5HcMeZubOa9K2c4HtECtG9nGN6R8EOPjMnjSWPKW', NOW(), NOW());  -- password: admin123
