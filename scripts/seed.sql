-- Seed data for E-Wallet API

-- API Clients
INSERT INTO api_clients (user_id, secret_key, is_active, created_at, updated_at)
VALUES 
    ('alif_partner', 'alif_secret_2025', true, NOW(), NOW()),
    ('megafon_api', 'megafon_key_secure', true, NOW(), NOW()),
    ('tcell_integration', 'tcell_hmac_key', true, NOW(), NOW())
ON CONFLICT (user_id) DO NOTHING;

-- Unidentified wallets (max 10,000 TJS)
INSERT INTO wallets (account_id, type, balance, created_at, updated_at)
VALUES 
    ('992900123456', 'unidentified', 250000, NOW(), NOW()),
    ('992935789012', 'unidentified', 500000, NOW(), NOW()),
    ('992918765432', 'unidentified', 750000, NOW(), NOW()),
    ('992987654321', 'unidentified', 1000000, NOW(), NOW()),
    ('992901234567', 'unidentified', 0, NOW(), NOW())
ON CONFLICT (account_id) DO NOTHING;

-- Identified wallets (max 100,000 TJS)
INSERT INTO wallets (account_id, type, balance, created_at, updated_at)
VALUES 
    ('992900111222', 'identified', 2500000, NOW(), NOW()),
    ('992935333444', 'identified', 5000000, NOW(), NOW()),
    ('992918555666', 'identified', 7500000, NOW(), NOW()),
    ('992987777888', 'identified', 10000000, NOW(), NOW()),
    ('992901999000', 'identified', 0, NOW(), NOW())
ON CONFLICT (account_id) DO NOTHING;

-- Sample transactions for testing monthly stats
DO $$
DECLARE
    w_id BIGINT;
BEGIN
    SELECT id INTO w_id FROM wallets WHERE account_id = '992900123456' LIMIT 1;
    IF w_id IS NOT NULL THEN
        INSERT INTO transactions (wallet_id, type, amount, created_at)
        VALUES 
            (w_id, 'deposit', 50000, NOW() - INTERVAL '2 days'),
            (w_id, 'deposit', 100000, NOW() - INTERVAL '7 days'),
            (w_id, 'deposit', 75000, NOW() - INTERVAL '12 days');
    END IF;
    
    SELECT id INTO w_id FROM wallets WHERE account_id = '992900111222' LIMIT 1;
    IF w_id IS NOT NULL THEN
        INSERT INTO transactions (wallet_id, type, amount, created_at)
        VALUES 
            (w_id, 'deposit', 250000, NOW() - INTERVAL '1 day'),
            (w_id, 'deposit', 500000, NOW() - INTERVAL '5 days'),
            (w_id, 'deposit', 1000000, NOW() - INTERVAL '10 days');
    END IF;
END $$;
