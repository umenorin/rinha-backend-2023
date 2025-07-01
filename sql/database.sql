DO
$do$
BEGIN
   IF EXISTS (SELECT FROM pg_database WHERE datname = 'rinha') THEN
      RAISE NOTICE 'Database already exists';
   ELSE
      PERFORM dblink_exec('dbname=' || current_database()
                        , 'CREATE DATABASE rinha');
   END IF;
END
$do$;

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
