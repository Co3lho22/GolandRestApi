CREATE DATABASE IF NOT EXISTS RestApi;

USE RestApi;

CREATE TABLE USERS (
                    id INT AUTO_INCREMENT PRIMARY KEY,
                    username VARCHAR(255) NOT NULL,
                    hashed_password VARCHAR(255) NOT NULL,
                    email VARCHAR(255),
                    country VARCHAR(255),
                    phone VARCHAR(255),
                    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE USER_AUTH (
                    user_id INT,
                    refresh_token VARCHAR(255),
                    PRIMARY KEY (user_id),
                    FOREIGN KEY (user_id) REFERENCES USERS(id)
);

CREATE TABLE ROLE (
                    id INT AUTO_INCREMENT PRIMARY KEY,
                    name VARCHAR(255) NOT NULL
);

CREATE TABLE PERMISSION (
                     id INT AUTO_INCREMENT PRIMARY KEY,
                     name VARCHAR(255) NOT NULL
);

CREATE TABLE ROLE_PERMISSION (
                    role_id INT,
                    permission_id INT,
                    PRIMARY KEY (role_id, permission_id),
                    FOREIGN KEY (role_id) REFERENCES ROLE(id),
                    FOREIGN KEY (permission_id) REFERENCES PERMISSION(id)
);

CREATE TABLE USER_ROLE (
                    user_id INT,
                    role_id INT,
                    PRIMARY KEY (user_id, role_id),
                    FOREIGN KEY (user_id) REFERENCES USERS(id),
                    FOREIGN KEY (role_id) REFERENCES ROLE(id)
);

INSERT INTO ROLE (name) VALUES ('admin'), ('user');

INSERT INTO PERMISSION (name) VALUES ('read'), ('write'), ('delete');

INSERT INTO ROLE_PERMISSION (role_id, permission_id) VALUES (1, 1), (1, 2), (1, 3);

INSERT INTO ROLE_PERMISSION (role_id, permission_id) VALUES (2, 1);