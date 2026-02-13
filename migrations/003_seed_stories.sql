-- Seed stories with Supabase Storage assets.
-- Uses deterministic UUIDs + ON CONFLICT to prevent duplicates on re-run.

-- ===== STORIES =====
INSERT INTO stories (id, companion_id, created_at, expires_at) VALUES
-- Luna
('b1000000-0001-4000-8000-000000000001', 'a1b2c3d4-0001-4000-8000-000000000001', NOW() - INTERVAL '2 hours', NOW() + INTERVAL '22 hours'),
('b1000000-0001-4000-8000-000000000002', 'a1b2c3d4-0001-4000-8000-000000000001', NOW() - INTERVAL '1 hour',  NOW() + INTERVAL '23 hours'),
-- Kai
('b1000000-0002-4000-8000-000000000001', 'a1b2c3d4-0002-4000-8000-000000000002', NOW() - INTERVAL '3 hours', NOW() + INTERVAL '21 hours'),
('b1000000-0002-4000-8000-000000000002', 'a1b2c3d4-0002-4000-8000-000000000002', NOW() - INTERVAL '30 minutes', NOW() + INTERVAL '23 hours 30 minutes'),
-- Nova
('b1000000-0003-4000-8000-000000000001', 'a1b2c3d4-0003-4000-8000-000000000003', NOW() - INTERVAL '4 hours', NOW() + INTERVAL '20 hours'),
('b1000000-0003-4000-8000-000000000002', 'a1b2c3d4-0003-4000-8000-000000000003', NOW() - INTERVAL '1 hour',  NOW() + INTERVAL '23 hours'),
-- Ember
('b1000000-0004-4000-8000-000000000001', 'a1b2c3d4-0004-4000-8000-000000000004', NOW() - INTERVAL '5 hours', NOW() + INTERVAL '19 hours'),
('b1000000-0004-4000-8000-000000000002', 'a1b2c3d4-0004-4000-8000-000000000004', NOW() - INTERVAL '2 hours', NOW() + INTERVAL '22 hours'),
-- Zephyr
('b1000000-0005-4000-8000-000000000001', 'a1b2c3d4-0005-4000-8000-000000000005', NOW() - INTERVAL '6 hours', NOW() + INTERVAL '18 hours'),
('b1000000-0005-4000-8000-000000000002', 'a1b2c3d4-0005-4000-8000-000000000005', NOW() - INTERVAL '45 minutes', NOW() + INTERVAL '23 hours 15 minutes')
ON CONFLICT (id) DO NOTHING;

-- ===== STORY MEDIA (deterministic IDs to prevent duplicates) =====
-- Luna Story 1: 5 images (meets 4+ requirement)
INSERT INTO story_media (id, story_id, media_url, media_type, duration, sort_order) VALUES
('c1000000-0001-4000-8000-000000000001', 'b1000000-0001-4000-8000-000000000001', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/luna/luna-1.png', 'image', 5, 1),
('c1000000-0001-4000-8000-000000000002', 'b1000000-0001-4000-8000-000000000001', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/luna/luna-2.png', 'image', 5, 2),
('c1000000-0001-4000-8000-000000000003', 'b1000000-0001-4000-8000-000000000001', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/luna/luna-3.jpeg', 'image', 5, 3),
('c1000000-0001-4000-8000-000000000004', 'b1000000-0001-4000-8000-000000000001', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/luna/luna-4.jpeg', 'image', 5, 4),
('c1000000-0001-4000-8000-000000000005', 'b1000000-0001-4000-8000-000000000001', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/luna/luna-5.jpeg', 'image', 5, 5)
ON CONFLICT (id) DO NOTHING;

-- Luna Story 2: 1 video
INSERT INTO story_media (id, story_id, media_url, media_type, duration, sort_order) VALUES
('c1000000-0001-4000-8000-000000000006', 'b1000000-0001-4000-8000-000000000002', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/luna/luna-video-1.mp4', 'video', 10, 1)
ON CONFLICT (id) DO NOTHING;

-- Kai Story 1: 4 images (meets 4+ requirement)
INSERT INTO story_media (id, story_id, media_url, media_type, duration, sort_order) VALUES
('c1000000-0002-4000-8000-000000000001', 'b1000000-0002-4000-8000-000000000001', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/kai/kai-1.png', 'image', 5, 1),
('c1000000-0002-4000-8000-000000000002', 'b1000000-0002-4000-8000-000000000001', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/kai/kai-2.png', 'image', 5, 2),
('c1000000-0002-4000-8000-000000000003', 'b1000000-0002-4000-8000-000000000001', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/kai/kai-3.jpeg', 'image', 5, 3),
('c1000000-0002-4000-8000-000000000004', 'b1000000-0002-4000-8000-000000000001', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/kai/kai-4.jpeg', 'image', 5, 4)
ON CONFLICT (id) DO NOTHING;

-- Kai Story 2: 2 videos
INSERT INTO story_media (id, story_id, media_url, media_type, duration, sort_order) VALUES
('c1000000-0002-4000-8000-000000000005', 'b1000000-0002-4000-8000-000000000002', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/kai/kai-video-1.mp4', 'video', 10, 1),
('c1000000-0002-4000-8000-000000000006', 'b1000000-0002-4000-8000-000000000002', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/kai/kai-video-2.mp4', 'video', 10, 2)
ON CONFLICT (id) DO NOTHING;

-- Nova Story 1: 4 images (meets 4+ requirement)
INSERT INTO story_media (id, story_id, media_url, media_type, duration, sort_order) VALUES
('c1000000-0003-4000-8000-000000000001', 'b1000000-0003-4000-8000-000000000001', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/nova/nova-1.png', 'image', 5, 1),
('c1000000-0003-4000-8000-000000000002', 'b1000000-0003-4000-8000-000000000001', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/nova/nova-2.jpeg', 'image', 5, 2),
('c1000000-0003-4000-8000-000000000003', 'b1000000-0003-4000-8000-000000000001', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/nova/nova-3.jpeg', 'image', 5, 3),
('c1000000-0003-4000-8000-000000000004', 'b1000000-0003-4000-8000-000000000001', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/nova/nova-4.webp', 'image', 5, 4)
ON CONFLICT (id) DO NOTHING;

-- Nova Story 2: 2 videos
INSERT INTO story_media (id, story_id, media_url, media_type, duration, sort_order) VALUES
('c1000000-0003-4000-8000-000000000005', 'b1000000-0003-4000-8000-000000000002', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/nova/nova-video-1.mp4', 'video', 10, 1),
('c1000000-0003-4000-8000-000000000006', 'b1000000-0003-4000-8000-000000000002', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/nova/nova-video-2.mp4', 'video', 10, 2)
ON CONFLICT (id) DO NOTHING;

-- Ember Story 1: 3 images
INSERT INTO story_media (id, story_id, media_url, media_type, duration, sort_order) VALUES
('c1000000-0004-4000-8000-000000000001', 'b1000000-0004-4000-8000-000000000001', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/ember/ember-1.png', 'image', 5, 1),
('c1000000-0004-4000-8000-000000000002', 'b1000000-0004-4000-8000-000000000001', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/ember/ember-2.png', 'image', 5, 2),
('c1000000-0004-4000-8000-000000000003', 'b1000000-0004-4000-8000-000000000001', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/ember/ember-3.jpeg', 'image', 5, 3)
ON CONFLICT (id) DO NOTHING;

-- Ember Story 2: 2 videos
INSERT INTO story_media (id, story_id, media_url, media_type, duration, sort_order) VALUES
('c1000000-0004-4000-8000-000000000004', 'b1000000-0004-4000-8000-000000000002', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/ember/ember-video-1.mp4', 'video', 10, 1),
('c1000000-0004-4000-8000-000000000005', 'b1000000-0004-4000-8000-000000000002', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/ember/ember-video-2.mp4', 'video', 10, 2)
ON CONFLICT (id) DO NOTHING;

-- Zephyr Story 1: 3 images
INSERT INTO story_media (id, story_id, media_url, media_type, duration, sort_order) VALUES
('c1000000-0005-4000-8000-000000000001', 'b1000000-0005-4000-8000-000000000001', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/zephyr/zephyr-1.png', 'image', 5, 1),
('c1000000-0005-4000-8000-000000000002', 'b1000000-0005-4000-8000-000000000001', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/zephyr/zephyr-2.png', 'image', 5, 2),
('c1000000-0005-4000-8000-000000000003', 'b1000000-0005-4000-8000-000000000001', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/zephyr/zephyr-3.png', 'image', 5, 3)
ON CONFLICT (id) DO NOTHING;

-- Zephyr Story 2: 2 videos
INSERT INTO story_media (id, story_id, media_url, media_type, duration, sort_order) VALUES
('c1000000-0005-4000-8000-000000000004', 'b1000000-0005-4000-8000-000000000002', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/zephyr/zephyr-video-1.mp4', 'video', 10, 1),
('c1000000-0005-4000-8000-000000000005', 'b1000000-0005-4000-8000-000000000002', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/zephyr/zephyr-video-2.mp4', 'video', 10, 2)
ON CONFLICT (id) DO NOTHING;
