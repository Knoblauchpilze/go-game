
CREATE TABLE users (
  id uuid NOT NULL DEFAULT uuid_generate_v4(),
  mail text NOT NULL,
  name text NOT NULL,
  password text NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE (mail)
);
