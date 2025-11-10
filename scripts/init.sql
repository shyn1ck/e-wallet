-- Create database if not exists
SELECT 'CREATE DATABASE e_wallet_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'e_wallet_db')\gexec

-- Connect to the database
\c e_wallet_db
