-- This is the SQL script that will be used to initialize the database schema.
-- We will evaluate you based on how well you design your database.
-- 1. How you design the tables.
-- 2. How you choose the data types and keys.
-- 3. How you name the fields.
-- In this assignment we will use PostgreSQL as the database.

-- Tables: estates
CREATE TABLE IF NOT EXISTS estates (
    id UUID PRIMARY KEY,
    width INT NOT NULL CHECK (width > 0 AND width <= 50000),
    length INT NOT NULL CHECK (length > 0 AND length <= 50000),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

-- Table: trees
CREATE TABLE IF NOT EXISTS trees (
    id UUID PRIMARY KEY,
    estate_id UUID NOT NULL REFERENCES estates(id) ON DELETE CASCADE,
    x INT NOT NULL CHECK (x > 0),
    y INT NOT NULL CHECK (y > 0),
    height INT NOT NULL CHECK (height >= 1 AND height <= 30),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

-- Table: estate_stats
CREATE TABLE IF NOT EXISTS estate_stats (
    id UUID PRIMARY KEY,
    estate_id UUID UNIQUE NOT NULL REFERENCES estates(id) ON DELETE CASCADE,
	tree_count BIGINT NOT NULL DEFAULT 0,
    max_height INT NOT NULL DEFAULT 0,
    min_height INT NOT NULL DEFAULT 0,
    median_height INT NOT NULL DEFAULT 0,
	drone_distance BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

-- Index for trees table
CREATE INDEX IF NOT EXISTS idx_trees_estate_id ON trees(estate_id);