CREATE TABLE users_data (
  user_id PRIMARY KEY AUTO_INCREMENT,
  first_name VARCHAR(50) NOT NULL,
  second_name VARCHAR(50) NOT NULL,
  email VARCHAR(100) NOT NULL,
  phone VARCHAR(20),
  current_password VARCHAR(255) NOT NULL,
  role_in_network ENUM('admin', 'user') NOT NULL,
  account_status ENUM('active', 'deleted') NOT NULL,
  avatar_ref MEDIUMTEXT
);

CREATE TABLE users_passwords (
  password_id PRIMARY KEY AUTO_INCREMENT,
  password VARCHAR(200) NOT NULL,
  password_created TIMESTAMP NOT NULL,
  user_id INTEGER REFERENCES users_data(user_id)
);

CREATE TABLE friend_list (
  friend_id INTEGER REFERENCES users_data(user_id),
  user_id INTEGER REFERENCES users_data(user_id),
  user_id_equals_friend_id CHAR(0) AS (CASE WHEN friend_id NOT IN (user_id) THEN '' END) VIRTUAL NOT NULL
);
