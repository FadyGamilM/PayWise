
CREATE TABLE accounts (
  id BIGSERIAL PRIMARY KEY,
  owner_name VARCHAR(100) NOT NULL,
  balance DECIMAL NOT NULL CHECK(balance > 0.0),
  currency VARCHAR(10) NOT NULL,
  removed BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS entries(
	id BIGSERIAL PRIMARY KEY,
  	account_id BIGINT REFERENCES accounts(id) ON DELETE CASCADE,
  	amount DECIMAL Not NULL,
    removed BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
	updated_at TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS transfers(
	id BIGSERIAL PRIMARY KEY,
  	to_account BIGINT REFERENCES accounts(id) ON DELETE CASCADE,
  	from_account BIGINT REFERENCES accounts(id) ON DELETE CASCADE,
  	amount DECIMAL Not NULL CHECK(Amount >= 0),
  	removed BOOLEAN Not NULL DEFAULT FALSE,
	created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
	updated_at TIMESTAMPTZ DEFAULT now() NOT NULL
);



CREATE OR REPLACE FUNCTION update_timestamp_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now();
   RETURN NEW;
END;
$$ language 'plpgsql';


CREATE TRIGGER update_timestamp_accounts
BEFORE UPDATE ON accounts
FOR EACH ROW
EXECUTE PROCEDURE update_timestamp_column();

CREATE TRIGGER update_timestamp_entries
BEFORE UPDATE ON entries
FOR EACH ROW
EXECUTE PROCEDURE update_timestamp_column();

CREATE TRIGGER update_timestamp_transfers
BEFORE UPDATE ON transfers
FOR EACH ROW
EXECUTE PROCEDURE update_timestamp_column();

create INDEX on accounts(id);
create INDEX on entries(account_id);
create INDEX on transfers(to_account);
create INDEX on transfers(from_account);
create INDEX on transfers(to_account, from_account);