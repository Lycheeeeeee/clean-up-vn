-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE "users" (
  "id" SERIAL NOT NULL PRIMARY KEY,
  "created_at" timestamptz DEFAULT now(),
  "deleted_at" timestamptz,
  "updated_at" timestamptz,
  "social" varchar,
  "displayname" varchar(200),
  "email" varchar(200),
  "password" varchar,
  "token" varchar,
  "issubscribed" boolean
);

INSERT INTO users(id, displayname, email, issubscribed) VALUES (123456,'aws', 'khanhniii07@gmail.com',false);
INSERT INTO users(id, displayname, email, issubscribed) VALUES (456789,'gcloud', 'khanhniii07@gmail.com',false);


CREATE TABLE "projects"(
  "id" SERIAL NOT NULL PRIMARY KEY,
  "created_at" timestamptz DEFAULT now(),
  "deleted_at" timestamptz,
  "updated_at" timestamptz,
  "name" varchar(200),
  "description" varchar,
  "owner" integer REFERENCES users(id),
  "longtitude" float,
  "latitude" float,
  "status" varchar(20),
  "time" timestamptz,
  "result" integer 
);
INSERT INTO projects(name, description, owner, longtitude, latitude, status, time, result) VALUES ("RMIT main gate","each people will have one pair of gloves", 123456, 12.32, 34.15, "upcoming","2019-11-12T07:59:24.251337Z",0);


CREATE TABLE "user_projects" (
    "created_at" timestamptz DEFAULT now(),
    "deleted_at" timestamptz,
    "updated_at" timestamptz,
    "id" integer REFERENCES projects(id),
    "user_id" integer REFERENCES users(id),
    PRIMARY KEY(id,user_id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;

