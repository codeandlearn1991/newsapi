CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS news (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  author TEXT NOT NULL,
  title TEXT NOT NULL,
  summary TEXT NOT NULL,
  content TEXT NOT NULL,
  source TEXT NOT NULL,
  tags TEXT[] NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

INSERT INTO news (id, author, title, summary, content, source, tags, created_at, updated_at)
VALUES (
  '17628bea-9d11-47f9-986e-16703a87e451',
  'Batman', 
  'Breaking News', 
  'A brief summary of the news', 
  'Full content of the news article', 
  'https://www.example.com', 
  ARRAY ['tag1', 'tag2'], 
  NOW(), 
  NOW());
INSERT INTO news (id, author, title, summary, content, source, tags, created_at, updated_at)
VALUES (
  'bde0c593-0df6-4eba-9326-3f00be67aade',
  'Superman', 
  'Breaking News', 
  'A brief summary of the news', 
  'Full content of the news article', 
  'https://www.example.com', 
  ARRAY ['tag1', 'tag2'], 
  NOW(), 
  NOW());
-- Deleted News --
INSERT INTO news (id, author, title, summary, content, source, tags, created_at, updated_at, deleted_at)
VALUES (
  'f710bc79-9ad3-4e0f-8dab-e43d94b42fbb',
  'Spiderman', 
  'Breaking News', 
  'A brief summary of the news', 
  'Full content of the news article', 
  'https://www.batman.com', 
  ARRAY ['tag1', 'Superhero'], 
  NOW(), 
  NOW(), 
  NOW());