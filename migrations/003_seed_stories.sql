-- Seed stories and story media for each companion.
-- Each companion gets 2–3 stories with 2–4 media slides.
-- Stories expire 30 days from migration run so they stay active.

-- ===================
-- Luna — stargazing & poetry
-- ===================
INSERT INTO stories (id, companion_id, created_at, expires_at) VALUES
('b1000000-0001-4000-8000-000000000001', 'a1b2c3d4-0001-4000-8000-000000000001', now() - interval '2 hours',  now() + interval '30 days'),
('b1000000-0001-4000-8000-000000000002', 'a1b2c3d4-0001-4000-8000-000000000001', now() - interval '30 minutes', now() + interval '30 days')
ON CONFLICT (id) DO NOTHING;

INSERT INTO story_media (id, story_id, media_url, media_type, duration, sort_order) VALUES
-- Luna story 1: stargazing night
('c1000000-0001-4000-8000-000000000001', 'b1000000-0001-4000-8000-000000000001', 'https://images.unsplash.com/photo-1519681393784-d120267933ba?w=800', 'image', 5, 0),
('c1000000-0001-4000-8000-000000000002', 'b1000000-0001-4000-8000-000000000001', 'https://images.unsplash.com/photo-1507400492013-162706c8c05e?w=800', 'image', 5, 1),
('c1000000-0001-4000-8000-000000000003', 'b1000000-0001-4000-8000-000000000001', 'https://images.unsplash.com/photo-1444703686981-a3abbc4d4fe3?w=800', 'image', 5, 2),
-- Luna story 2: poetry & books
('c1000000-0001-4000-8000-000000000004', 'b1000000-0001-4000-8000-000000000002', 'https://images.unsplash.com/photo-1512820790803-83ca734da794?w=800', 'image', 5, 0),
('c1000000-0001-4000-8000-000000000005', 'b1000000-0001-4000-8000-000000000002', 'https://images.unsplash.com/photo-1474366521946-c3b8e5c7e68a?w=800', 'image', 5, 1)
ON CONFLICT (id) DO NOTHING;

-- ===================
-- Kai — adventure & surfing
-- ===================
INSERT INTO stories (id, companion_id, created_at, expires_at) VALUES
('b1000000-0002-4000-8000-000000000001', 'a1b2c3d4-0002-4000-8000-000000000002', now() - interval '3 hours',  now() + interval '30 days'),
('b1000000-0002-4000-8000-000000000002', 'a1b2c3d4-0002-4000-8000-000000000002', now() - interval '1 hour',   now() + interval '30 days'),
('b1000000-0002-4000-8000-000000000003', 'a1b2c3d4-0002-4000-8000-000000000002', now() - interval '15 minutes', now() + interval '30 days')
ON CONFLICT (id) DO NOTHING;

INSERT INTO story_media (id, story_id, media_url, media_type, duration, sort_order) VALUES
-- Kai story 1: mountain hiking
('c1000000-0002-4000-8000-000000000001', 'b1000000-0002-4000-8000-000000000001', 'https://images.unsplash.com/photo-1464822759023-fed622ff2c3b?w=800', 'image', 5, 0),
('c1000000-0002-4000-8000-000000000002', 'b1000000-0002-4000-8000-000000000001', 'https://images.unsplash.com/photo-1551632811-561732d1e306?w=800', 'image', 5, 1),
('c1000000-0002-4000-8000-000000000003', 'b1000000-0002-4000-8000-000000000001', 'https://images.unsplash.com/photo-1483728642387-6c3bdd6c93e5?w=800', 'image', 5, 2),
-- Kai story 2: surfing
('c1000000-0002-4000-8000-000000000004', 'b1000000-0002-4000-8000-000000000002', 'https://images.unsplash.com/photo-1502680390548-bdbac40b3e1a?w=800', 'image', 5, 0),
('c1000000-0002-4000-8000-000000000005', 'b1000000-0002-4000-8000-000000000002', 'https://images.unsplash.com/photo-1455729552457-5c322b382635?w=800', 'image', 5, 1),
-- Kai story 3: campfire at night
('c1000000-0002-4000-8000-000000000006', 'b1000000-0002-4000-8000-000000000003', 'https://images.unsplash.com/photo-1475483768296-6163e08872a1?w=800', 'image', 5, 0),
('c1000000-0002-4000-8000-000000000007', 'b1000000-0002-4000-8000-000000000003', 'https://images.unsplash.com/photo-1504851149312-7a075b496cc7?w=800', 'image', 5, 1)
ON CONFLICT (id) DO NOTHING;

-- ===================
-- Nova — city & science
-- ===================
INSERT INTO stories (id, companion_id, created_at, expires_at) VALUES
('b1000000-0003-4000-8000-000000000001', 'a1b2c3d4-0003-4000-8000-000000000003', now() - interval '4 hours',  now() + interval '30 days'),
('b1000000-0003-4000-8000-000000000002', 'a1b2c3d4-0003-4000-8000-000000000003', now() - interval '45 minutes', now() + interval '30 days')
ON CONFLICT (id) DO NOTHING;

INSERT INTO story_media (id, story_id, media_url, media_type, duration, sort_order) VALUES
-- Nova story 1: city exploration
('c1000000-0003-4000-8000-000000000001', 'b1000000-0003-4000-8000-000000000001', 'https://images.unsplash.com/photo-1480714378408-67cf0d13bc1b?w=800', 'image', 5, 0),
('c1000000-0003-4000-8000-000000000002', 'b1000000-0003-4000-8000-000000000001', 'https://images.unsplash.com/photo-1514565131-fce0801e5785?w=800', 'image', 5, 1),
('c1000000-0003-4000-8000-000000000003', 'b1000000-0003-4000-8000-000000000001', 'https://images.unsplash.com/photo-1517248135467-4c7edcad34c4?w=800', 'image', 5, 2),
-- Nova story 2: tech & science
('c1000000-0003-4000-8000-000000000004', 'b1000000-0003-4000-8000-000000000002', 'https://images.unsplash.com/photo-1518770660439-4636190af475?w=800', 'image', 5, 0),
('c1000000-0003-4000-8000-000000000005', 'b1000000-0003-4000-8000-000000000002', 'https://images.unsplash.com/photo-1507413245164-6160d8298b31?w=800', 'image', 5, 1),
('c1000000-0003-4000-8000-000000000006', 'b1000000-0003-4000-8000-000000000002', 'https://images.unsplash.com/photo-1451187580459-43490279c0fa?w=800', 'image', 5, 2)
ON CONFLICT (id) DO NOTHING;

-- ===================
-- Ember — cooking & garden
-- ===================
INSERT INTO stories (id, companion_id, created_at, expires_at) VALUES
('b1000000-0004-4000-8000-000000000001', 'a1b2c3d4-0004-4000-8000-000000000004', now() - interval '5 hours',  now() + interval '30 days'),
('b1000000-0004-4000-8000-000000000002', 'a1b2c3d4-0004-4000-8000-000000000004', now() - interval '1 hour',   now() + interval '30 days')
ON CONFLICT (id) DO NOTHING;

INSERT INTO story_media (id, story_id, media_url, media_type, duration, sort_order) VALUES
-- Ember story 1: baking & cooking
('c1000000-0004-4000-8000-000000000001', 'b1000000-0004-4000-8000-000000000001', 'https://images.unsplash.com/photo-1556909114-f6e7ad7d3136?w=800', 'image', 5, 0),
('c1000000-0004-4000-8000-000000000002', 'b1000000-0004-4000-8000-000000000001', 'https://images.unsplash.com/photo-1486427944544-d2c246c4df16?w=800', 'image', 5, 1),
('c1000000-0004-4000-8000-000000000003', 'b1000000-0004-4000-8000-000000000001', 'https://images.unsplash.com/photo-1464305795204-6f5bbfc7fb81?w=800', 'image', 5, 2),
-- Ember story 2: garden & flowers
('c1000000-0004-4000-8000-000000000004', 'b1000000-0004-4000-8000-000000000002', 'https://images.unsplash.com/photo-1416879595882-3373a0480b5b?w=800', 'image', 5, 0),
('c1000000-0004-4000-8000-000000000005', 'b1000000-0004-4000-8000-000000000002', 'https://images.unsplash.com/photo-1490750967868-88aa4f44baee?w=800', 'image', 5, 1)
ON CONFLICT (id) DO NOTHING;

-- ===================
-- Zephyr — street art & festival
-- ===================
INSERT INTO stories (id, companion_id, created_at, expires_at) VALUES
('b1000000-0005-4000-8000-000000000001', 'a1b2c3d4-0005-4000-8000-000000000005', now() - interval '6 hours',  now() + interval '30 days'),
('b1000000-0005-4000-8000-000000000002', 'a1b2c3d4-0005-4000-8000-000000000005', now() - interval '20 minutes', now() + interval '30 days'),
('b1000000-0005-4000-8000-000000000003', 'a1b2c3d4-0005-4000-8000-000000000005', now() - interval '5 minutes',  now() + interval '30 days')
ON CONFLICT (id) DO NOTHING;

INSERT INTO story_media (id, story_id, media_url, media_type, duration, sort_order) VALUES
-- Zephyr story 1: street art
('c1000000-0005-4000-8000-000000000001', 'b1000000-0005-4000-8000-000000000001', 'https://images.unsplash.com/photo-1499781350541-7783f6c6a0c8?w=800', 'image', 5, 0),
('c1000000-0005-4000-8000-000000000002', 'b1000000-0005-4000-8000-000000000001', 'https://images.unsplash.com/photo-1460661419201-fd4cecdf8a8b?w=800', 'image', 5, 1),
('c1000000-0005-4000-8000-000000000003', 'b1000000-0005-4000-8000-000000000001', 'https://images.unsplash.com/photo-1561059488-916d69792237?w=800', 'image', 5, 2),
-- Zephyr story 2: music festival
('c1000000-0005-4000-8000-000000000004', 'b1000000-0005-4000-8000-000000000002', 'https://images.unsplash.com/photo-1459749411175-04bf5292ceea?w=800', 'image', 5, 0),
('c1000000-0005-4000-8000-000000000005', 'b1000000-0005-4000-8000-000000000002', 'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=800', 'image', 5, 1),
('c1000000-0005-4000-8000-000000000006', 'b1000000-0005-4000-8000-000000000002', 'https://images.unsplash.com/photo-1470229722913-7c0e2dbbafd3?w=800', 'image', 5, 2),
-- Zephyr story 3: skateboarding
('c1000000-0005-4000-8000-000000000007', 'b1000000-0005-4000-8000-000000000003', 'https://images.unsplash.com/photo-1564429238961-bf8f8be6b7cc?w=800', 'image', 5, 0),
('c1000000-0005-4000-8000-000000000008', 'b1000000-0005-4000-8000-000000000003', 'https://images.unsplash.com/photo-1547447134-cd3f5c716030?w=800', 'image', 5, 1)
ON CONFLICT (id) DO NOTHING;
