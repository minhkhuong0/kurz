CREATE TABLE IF NOT EXISTS shorturls(
   code VARCHAR(10) PRIMARY KEY,
   url TEXT NOT NULL,
   created_at TIMESTAMP NOT NULL,
   expires_at TIMESTAMP
);
