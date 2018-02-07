CREATE TABLE users_data (
  user_id SERIAL PRIMARY KEY,
  email VARCHAR(100) NOT NULL,
  first_name VARCHAR(50) NOT NULL,
  last_name VARCHAR(50) NOT NULL,
  phone VARCHAR(20),
  role_in_network ENUM('admin', 'user') NOT NULL,
  account_status ENUM('active', 'deleted') NOT NULL,
  avatar_ref MEDIUMTEXT
);

CREATE TABLE users_passwords (
  password_id SERIAL PRIMARY KEY,
  password VARCHAR(200) NOT NULL,
  password_created TIMESTAMP NOT NULL,
  user_id INTEGER REFERENCES users_data(user_id)
);

CREATE TABLE friend_list (
  friend_user_id INTEGER REFERENCES users_data(user_id),
  user_id INTEGER REFERENCES users_data(user_id),
  connection_status ENUM('approved', 'rejected', 'waiting') NOT NULL,
  user_id_equals_friend_id CHAR(0) AS (CASE WHEN friend_id NOT IN (user_id) THEN '' END) VIRTUAL NOT NULL
);
