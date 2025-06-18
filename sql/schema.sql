CREATE TABLE scraped_data (
    id SERIAL PRIMARY KEY,
    platform VARCHAR(50) NOT NULL,
    data JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_platform ON scraped_data(platform);
CREATE INDEX idx_created_at ON scraped_data(created_at);