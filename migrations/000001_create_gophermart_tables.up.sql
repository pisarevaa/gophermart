CREATE TABLE IF NOT EXISTS users (
    "login" 		VARCHAR(250) PRIMARY KEY,
	"password" 		VARCHAR(250) NOT NULL,
	"balance"       DECIMAL NOT NULL DEFAULT 0,
	"withdrawn"     DECIMAL NOT NULL DEFAULT 0,
	"created_at" 	TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS orders (
    "number" 		VARCHAR(50) PRIMARY KEY,
	"status" 		VARCHAR(15) NOT NULL,
	"accrual" 		DECIMAL NOT NULL DEFAULT 0,
	"withdrawn" 	DECIMAL NOT NULL DEFAULT 0,
	"login"			VARCHAR(250) NOT NULL,
	"uploaded_at"	TIMESTAMPTZ NOT NULL,
	"processed_at" 	TIMESTAMPTZ NULL
);
ALTER TABLE orders ADD CONSTRAINT users_login_fkey FOREIGN KEY ("login") REFERENCES users("login");
