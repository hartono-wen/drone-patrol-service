-- This is the SQL script that will be used to initialize the database schema.
-- We will evaluate you based on how well you design your database.
-- 1. How you design the tables.
-- 2. How you choose the data types and keys.
-- 3. How you name the fields.
-- In this assignment we will use PostgreSQL as the database.

CREATE SCHEMA plantation_management_service;

CREATE TABLE IF NOT EXISTS plantation_management_service.estates (
	id UUID NOT NULL,
	length INTEGER NOT NULL CHECK (length BETWEEN 1 AND 50000),
	width INTEGER NOT NULL CHECK (width BETWEEN 1 AND 50000),
	created_at TIMESTAMPTZ NOT NULL,
    
	CONSTRAINT estate_pk PRIMARY KEY (id),
	CONSTRAINT estates_unique_keys UNIQUE (length, width)
);

CREATE TABLE IF NOT EXISTS plantation_management_service.trees (
	id UUID NOT NULL,
	estate_id UUID NOT NULL,
	x INTEGER NOT NULL,
	y INTEGER NOT NULL,
	height SMALLINT NOT NULL CHECK (height BETWEEN 1 AND 30),
	created_at TIMESTAMPTZ NOT NULL,
    
	CONSTRAINT tree_pk PRIMARY KEY (id),
	CONSTRAINT trees_estate_id_fk_estates_estate_id FOREIGN KEY(estate_id) REFERENCES plantation_management_service.estates(id)
);

