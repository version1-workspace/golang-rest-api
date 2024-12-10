INSERT INTO users (id, username, email) VALUES
  (1, 'user 1', 'This is the first post.'),
  (2, 'user 2', 'This is the second post.'),
  (3, 'user 3', 'This is the third post.');

INSERT INTO posts (user_id, title, content) VALUES
  (1, 'Post 1', 'This is the first post.'),
  (1, 'Post 2', 'This is the second post.'),
  (1, 'Post 3', 'This is the third post.');

INSERT INTO tags (slug, name) VALUES
  ('tag-1', 'Tag 1'),
  ('tag-2', 'Tag 2'),
  ('tag-3', 'Tag 3');

INSERT INTO tag_posts (post_id, tag_id) VALUES
  (1, 1),
  (1, 2),
  (2, 2),
  (2, 3),
  (3, 1),
  (3, 3);
