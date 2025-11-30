-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TYPE user_role AS ENUM ('USER', 'ADMIN', 'AGENT');
CREATE TYPE user_status AS ENUM ('pending_approval', 'approved', 'none');

CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  full_name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  age INT,
  phone TEXT,
  date_of_birth TEXT,
  country TEXT,
  state TEXT,
  school TEXT,
  profile_completed BOOLEAN,
  status user_status NOT NULL DEFAULT 'none',

  password TEXT NOT NULL,
  role user_role NOT NULL DEFAULT 'USER',

  -- Admin-only fields
  admin_code TEXT,
  department TEXT,
  


  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  last_login TIMESTAMP,
  updated_at TIMESTAMP
);


ALTER TABLE users
ADD CONSTRAINT admin_code_required_for_admin
CHECK (
  (role = 'ADMIN' AND admin_code IS NOT NULL)
  OR
  (role = 'USER')
);

ALTER TABLE users
ADD CONSTRAINT department_required_for_admin
CHECK (
  (role = 'ADMIN' AND department IS NOT NULL)
  OR
  (role = 'USER')
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
