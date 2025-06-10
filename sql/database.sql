CREATE DATABASE rinha;
\c rinha

CREATE TABLE person(
  id uuid,
  nickname VARCHAR(255),
  birthdate date,
  name VARCHAR(255),

  PRIMARY KEY(id)
);

CREATE TABLE language (
  id uuid ,
  name VARCHAR(255),

  PRIMARY KEY(id)
);


CREATE TABLE stack (
  person_id uuid references person(id),
  language_id uuid references language(id),

  PRIMARY KEY(person_id, language_id)
);
