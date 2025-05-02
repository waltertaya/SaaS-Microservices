CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    plan VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL
);

CREATE TABLE transactions (
    transaction_id VARCHAR(100) PRIMARY KEY,
    tx_ref VARCHAR(100) NOT NULL,
    amount VARCHAR(50) NOT NULL,
    currency VARCHAR(30) NOT NULL,
    status VARCHAR(50) NOT NULL,
    customer_email VARCHAR(100) NOT NULL,
    transaction_date TIMESTAMP NOT NULL
)

CREATE TABLE failed_webhooks (
    id SERIAL PRIMARY KEY,
    event VARCHAR(100) NOT NULL,
    payload VARCHAR(100) NOT NULL,
    retry_count INT NOT NULL,
)
