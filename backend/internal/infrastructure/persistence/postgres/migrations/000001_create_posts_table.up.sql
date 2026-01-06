-- Migration: Create posts and cities tables
-- Version: 000001
-- Description: Create posts table for storing exposure content and cities table for city configuration

-- Create posts table
CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_name VARCHAR(100) NOT NULL,
    city_code VARCHAR(50) NOT NULL,
    city_name VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    occurred_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for posts table
-- Index for city code (used in FindByCity)
CREATE INDEX IF NOT EXISTS idx_posts_city_code ON posts(city_code);

-- Index for created_at (used for ordering, DESC for latest first)
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC);

-- Index for company name (used for filtering and search)
CREATE INDEX IF NOT EXISTS idx_posts_company_name ON posts(company_name);

-- Full-text search index (using PostgreSQL built-in, not pg_jieba for now)
-- Note: pg_jieba requires additional extension installation
-- For now, use simple_tsvector which works out of the box
-- Future: Can be upgraded to use pg_jieba for better Chinese support
CREATE INDEX IF NOT EXISTS idx_posts_search ON posts USING GIN(
    to_tsvector('simple', company_name || ' ' || content)
);

-- Create cities table (configuration table for city data)
CREATE TABLE IF NOT EXISTS cities (
    code VARCHAR(50) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    pinyin VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create index for cities name (used for searching cities)
CREATE INDEX IF NOT EXISTS idx_cities_name ON cities(name);

-- Add comment to tables
COMMENT ON TABLE posts IS 'Stores exposure content posts';
COMMENT ON TABLE cities IS 'Stores city configuration data';

COMMENT ON COLUMN posts.id IS 'Unique identifier (UUID)';
COMMENT ON COLUMN posts.company_name IS 'Company name (1-100 characters)';
COMMENT ON COLUMN posts.city_code IS 'City code (e.g., beijing, shanghai)';
COMMENT ON COLUMN posts.city_name IS 'City name (e.g., 北京, 上海)';
COMMENT ON COLUMN posts.content IS 'Post content (10-5000 characters)';
COMMENT ON COLUMN posts.occurred_at IS 'When the incident occurred (optional)';
COMMENT ON COLUMN posts.created_at IS 'When the post was created';
COMMENT ON COLUMN posts.updated_at IS 'When the post was last updated';

COMMENT ON COLUMN cities.code IS 'City code (primary key)';
COMMENT ON COLUMN cities.name IS 'City name';
COMMENT ON COLUMN cities.pinyin IS 'City name in pinyin (optional)';

