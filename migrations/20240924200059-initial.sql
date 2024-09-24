-- +migrate Up
CREATE TABLE users (
    id INT PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    refresh_token TEXT NOT NULL
);

INSERT INTO users (id, email, username, refresh_token) VALUES
(1, 'email1@mail.com', 'user1', '$2y$10$a6ECj5MfyDQZ45inPf38weHvX0PaXP25RkcQ.1NkQ67CN86v2YmdK'),
(2, 'email2@mail.com', 'user2', '$2y$10$j0QR17FYUpMugT8aAjkp3u9LvuHbtIh70wkrbJFKpB51wOlNRicQ2'),
(3, 'email3@mail.com', 'user3', '$2y$10$TfI8/QNwnGE4REauwfiz/Om1KX7iputeoilflJPtJkJPXDCZ9JZDu');

-- +migrate Down
DROP TABLE IF EXISTS users;
