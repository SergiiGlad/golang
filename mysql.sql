CREATE TABLE IF NOT EXISTS users_data (
  user_id SERIAL PRIMARY KEY,
  first_name VARCHAR(50) NOT NULL,
  second_name VARCHAR(50) NOT NULL,
  email VARCHAR(100) NOT NULL,
  phone VARCHAR(20),
  current_password VARCHAR(255) NOT NULL,
  role_in_network ENUM('admin', 'user') NOT NULL,
  account_status ENUM('inactive','active', 'deleted') NOT NULL,
  avatar_ref MEDIUMTEXT
);

CREATE TABLE IF NOT EXISTS users_passwords (
  password_id SERIAL PRIMARY KEY,
  password VARCHAR(200) NOT NULL,
  password_created TIMESTAMP NOT NULL,
  user_id INTEGER REFERENCES users_data(user_id)
);

CREATE TABLE IF NOT EXISTS friend_list (
  friend_id INTEGER REFERENCES users_data(user_id),
  user_id INTEGER REFERENCES users_data(user_id),
  user_id_equals_friend_id CHAR(0) AS (CASE WHEN friend_id NOT IN (user_id) THEN '' END) VIRTUAL NOT NULL
);

CREATE TABLE IF NOT EXISTS user_tokens (
    token_id SERIAL PRIMARY KEY,
    token VARCHAR(128) NOT NULL,
    email VARCHAR(100) NOT NULL,
    is_active BOOLEAN,
    user_id INTEGER REFERENCES users_data(user_id)
  );
