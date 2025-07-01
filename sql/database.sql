CREATE DATABASE rinha;
\c rinha

CREATE TABLE IF NOT EXISTS  person(
  id uuid DEFAULT gen_random_uuid(),
  nickname varchar(255),
  birthdate date,
  name VARCHAR(255),

  PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS  language (
  id uuid DEFAULT gen_random_uuid(),
  name VARCHAR(255),

  PRIMARY KEY(id)
);


CREATE TABLE IF NOT EXISTS  stack (
  person_id uuid references person(id),
  language_id uuid references language(id),

  PRIMARY KEY(person_id, language_id)
);
