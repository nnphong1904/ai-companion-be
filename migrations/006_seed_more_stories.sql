-- Additional stories to meet assessment requirements:
--   - Each companion has a mix of photos and videos
--   - At least one story has 4+ media slides (Kai's road trip has 5, Luna's aurora night has 4)
--   - Video URLs use verified public test videos (test-videos.co.uk, w3schools)

-- ===================
-- Luna — aurora night (4 slides, mix of photo + video)
-- ===================
INSERT INTO stories (id, companion_id, created_at, expires_at) VALUES
('b1000000-0001-4000-8000-000000000003', 'a1b2c3d4-0001-4000-8000-000000000001', now() - interval '10 minutes', now() + interval '30 days')
ON CONFLICT (id) DO NOTHING;

INSERT INTO story_media (id, story_id, media_url, media_type, duration, sort_order) VALUES
('c1000000-0001-4000-8000-000000000006', 'b1000000-0001-4000-8000-000000000003', 'https://images.unsplash.com/photo-1531366936337-7c912a4589a7?w=800', 'image', 5, 0),
('c1000000-0001-4000-8000-000000000007', 'b1000000-0001-4000-8000-000000000003', 'https://images.unsplash.com/photo-1483347756197-71ef80e95f73?w=800', 'image', 5, 1),
('c1000000-0001-4000-8000-000000000008', 'b1000000-0001-4000-8000-000000000003', 'https://test-videos.co.uk/vids/sintel/mp4/h264/720/Sintel_720_10s_2MB.mp4', 'video', 10, 2),
('c1000000-0001-4000-8000-000000000009', 'b1000000-0001-4000-8000-000000000003', 'https://images.unsplash.com/photo-1507400492013-162706c8c05e?w=800', 'image', 5, 3)
ON CONFLICT (id) DO NOTHING;

-- ===================
-- Kai — epic road trip (5 slides — the longest story, mix of photo + video)
-- ===================
INSERT INTO stories (id, companion_id, created_at, expires_at) VALUES
('b1000000-0002-4000-8000-000000000004', 'a1b2c3d4-0002-4000-8000-000000000002', now() - interval '5 minutes', now() + interval '30 days')
ON CONFLICT (id) DO NOTHING;

INSERT INTO story_media (id, story_id, media_url, media_type, duration, sort_order) VALUES
('c1000000-0002-4000-8000-000000000008', 'b1000000-0002-4000-8000-000000000004', 'https://images.unsplash.com/photo-1469854523086-cc02fe5d8800?w=800', 'image', 5, 0),
('c1000000-0002-4000-8000-000000000009', 'b1000000-0002-4000-8000-000000000004', 'https://test-videos.co.uk/vids/bigbuckbunny/mp4/h264/720/Big_Buck_Bunny_720_10s_2MB.mp4', 'video', 10, 1),
('c1000000-0002-4000-8000-000000000010', 'b1000000-0002-4000-8000-000000000004', 'https://images.unsplash.com/photo-1500534314209-a25ddb2bd429?w=800', 'image', 5, 2),
('c1000000-0002-4000-8000-000000000011', 'b1000000-0002-4000-8000-000000000004', 'https://images.unsplash.com/photo-1501785888041-af3ef285b470?w=800', 'image', 5, 3),
('c1000000-0002-4000-8000-000000000012', 'b1000000-0002-4000-8000-000000000004', 'https://test-videos.co.uk/vids/jellyfish/mp4/h264/720/Jellyfish_720_10s_2MB.mp4', 'video', 10, 4)
ON CONFLICT (id) DO NOTHING;

-- ===================
-- Nova — rooftop timelapse (video story)
-- ===================
INSERT INTO stories (id, companion_id, created_at, expires_at) VALUES
('b1000000-0003-4000-8000-000000000003', 'a1b2c3d4-0003-4000-8000-000000000003', now() - interval '8 minutes', now() + interval '30 days')
ON CONFLICT (id) DO NOTHING;

INSERT INTO story_media (id, story_id, media_url, media_type, duration, sort_order) VALUES
('c1000000-0003-4000-8000-000000000007', 'b1000000-0003-4000-8000-000000000003', 'https://test-videos.co.uk/vids/bigbuckbunny/mp4/h264/360/Big_Buck_Bunny_360_10s_1MB.mp4', 'video', 10, 0),
('c1000000-0003-4000-8000-000000000008', 'b1000000-0003-4000-8000-000000000003', 'https://images.unsplash.com/photo-1519501025264-65ba15a82390?w=800', 'image', 5, 1),
('c1000000-0003-4000-8000-000000000009', 'b1000000-0003-4000-8000-000000000003', 'https://images.unsplash.com/photo-1477959858617-67f85cf4f1df?w=800', 'image', 5, 2)
ON CONFLICT (id) DO NOTHING;

-- ===================
-- Ember — cozy morning routine (video + photos)
-- ===================
INSERT INTO stories (id, companion_id, created_at, expires_at) VALUES
('b1000000-0004-4000-8000-000000000003', 'a1b2c3d4-0004-4000-8000-000000000004', now() - interval '12 minutes', now() + interval '30 days')
ON CONFLICT (id) DO NOTHING;

INSERT INTO story_media (id, story_id, media_url, media_type, duration, sort_order) VALUES
('c1000000-0004-4000-8000-000000000006', 'b1000000-0004-4000-8000-000000000003', 'https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?w=800', 'image', 5, 0),
('c1000000-0004-4000-8000-000000000007', 'b1000000-0004-4000-8000-000000000003', 'https://www.w3schools.com/html/mov_bbb.mp4', 'video', 10, 1),
('c1000000-0004-4000-8000-000000000008', 'b1000000-0004-4000-8000-000000000003', 'https://images.unsplash.com/photo-1504754524776-8f4f37790ca0?w=800', 'image', 5, 2)
ON CONFLICT (id) DO NOTHING;

-- ===================
-- Zephyr — late night jam session (video-heavy)
-- ===================
INSERT INTO stories (id, companion_id, created_at, expires_at) VALUES
('b1000000-0005-4000-8000-000000000004', 'a1b2c3d4-0005-4000-8000-000000000005', now() - interval '3 minutes', now() + interval '30 days')
ON CONFLICT (id) DO NOTHING;

INSERT INTO story_media (id, story_id, media_url, media_type, duration, sort_order) VALUES
('c1000000-0005-4000-8000-000000000009', 'b1000000-0005-4000-8000-000000000004', 'https://test-videos.co.uk/vids/jellyfish/mp4/h264/360/Jellyfish_360_10s_1MB.mp4', 'video', 10, 0),
('c1000000-0005-4000-8000-000000000010', 'b1000000-0005-4000-8000-000000000004', 'https://images.unsplash.com/photo-1511671782779-c97d3d27a1d4?w=800', 'image', 5, 1),
('c1000000-0005-4000-8000-000000000011', 'b1000000-0005-4000-8000-000000000004', 'https://test-videos.co.uk/vids/sintel/mp4/h264/360/Sintel_360_10s_1MB.mp4', 'video', 10, 2)
ON CONFLICT (id) DO NOTHING;
