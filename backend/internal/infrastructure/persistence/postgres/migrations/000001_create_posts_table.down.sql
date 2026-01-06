-- Migration: Drop posts and cities tables
-- Version: 000001
-- Description: Rollback migration - drop all tables and indexes created in up migration

-- Drop indexes first (they depend on tables)
DROP INDEX IF EXISTS idx_cities_name;
DROP INDEX IF EXISTS idx_posts_search;
DROP INDEX IF EXISTS idx_posts_company_name;
DROP INDEX IF EXISTS idx_posts_created_at;
DROP INDEX IF EXISTS idx_posts_city_code;

-- Drop tables
DROP TABLE IF EXISTS cities;
DROP TABLE IF EXISTS posts;

