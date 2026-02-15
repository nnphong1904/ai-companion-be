-- Seed stories with Supabase Storage assets.
-- Each story has exactly 1 media asset. Uses ON CONFLICT to prevent duplicates.

-- Base URL: https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories

-- ===== STORIES (1 per asset) =====
INSERT INTO stories (id, companion_id, created_at, expires_at) VALUES
-- Luna (6 assets)
('b1000000-0001-4000-8000-000000000001', 'a1b2c3d4-0001-4000-8000-000000000001', NOW() - INTERVAL '5 hours',      NOW() + INTERVAL '19 hours'),
('b1000000-0001-4000-8000-000000000002', 'a1b2c3d4-0001-4000-8000-000000000001', NOW() - INTERVAL '4 hours',      NOW() + INTERVAL '20 hours'),
('b1000000-0001-4000-8000-000000000003', 'a1b2c3d4-0001-4000-8000-000000000001', NOW() - INTERVAL '3 hours',      NOW() + INTERVAL '21 hours'),
('b1000000-0001-4000-8000-000000000004', 'a1b2c3d4-0001-4000-8000-000000000001', NOW() - INTERVAL '2 hours',      NOW() + INTERVAL '22 hours'),
('b1000000-0001-4000-8000-000000000005', 'a1b2c3d4-0001-4000-8000-000000000001', NOW() - INTERVAL '1 hour',       NOW() + INTERVAL '23 hours'),
('b1000000-0001-4000-8000-000000000006', 'a1b2c3d4-0001-4000-8000-000000000001', NOW() - INTERVAL '30 minutes',   NOW() + INTERVAL '23 hours 30 minutes'),
-- Kai (6 assets)
('b1000000-0002-4000-8000-000000000001', 'a1b2c3d4-0002-4000-8000-000000000002', NOW() - INTERVAL '5 hours',      NOW() + INTERVAL '19 hours'),
('b1000000-0002-4000-8000-000000000002', 'a1b2c3d4-0002-4000-8000-000000000002', NOW() - INTERVAL '4 hours',      NOW() + INTERVAL '20 hours'),
('b1000000-0002-4000-8000-000000000003', 'a1b2c3d4-0002-4000-8000-000000000002', NOW() - INTERVAL '3 hours',      NOW() + INTERVAL '21 hours'),
('b1000000-0002-4000-8000-000000000004', 'a1b2c3d4-0002-4000-8000-000000000002', NOW() - INTERVAL '2 hours',      NOW() + INTERVAL '22 hours'),
('b1000000-0002-4000-8000-000000000005', 'a1b2c3d4-0002-4000-8000-000000000002', NOW() - INTERVAL '1 hour',       NOW() + INTERVAL '23 hours'),
('b1000000-0002-4000-8000-000000000006', 'a1b2c3d4-0002-4000-8000-000000000002', NOW() - INTERVAL '30 minutes',   NOW() + INTERVAL '23 hours 30 minutes'),
-- Nova (6 assets)
('b1000000-0003-4000-8000-000000000001', 'a1b2c3d4-0003-4000-8000-000000000003', NOW() - INTERVAL '5 hours',      NOW() + INTERVAL '19 hours'),
('b1000000-0003-4000-8000-000000000002', 'a1b2c3d4-0003-4000-8000-000000000003', NOW() - INTERVAL '4 hours',      NOW() + INTERVAL '20 hours'),
('b1000000-0003-4000-8000-000000000003', 'a1b2c3d4-0003-4000-8000-000000000003', NOW() - INTERVAL '3 hours',      NOW() + INTERVAL '21 hours'),
('b1000000-0003-4000-8000-000000000004', 'a1b2c3d4-0003-4000-8000-000000000003', NOW() - INTERVAL '2 hours',      NOW() + INTERVAL '22 hours'),
('b1000000-0003-4000-8000-000000000005', 'a1b2c3d4-0003-4000-8000-000000000003', NOW() - INTERVAL '1 hour',       NOW() + INTERVAL '23 hours'),
('b1000000-0003-4000-8000-000000000006', 'a1b2c3d4-0003-4000-8000-000000000003', NOW() - INTERVAL '30 minutes',   NOW() + INTERVAL '23 hours 30 minutes'),
-- Ember (5 assets)
('b1000000-0004-4000-8000-000000000001', 'a1b2c3d4-0004-4000-8000-000000000004', NOW() - INTERVAL '4 hours',      NOW() + INTERVAL '20 hours'),
('b1000000-0004-4000-8000-000000000002', 'a1b2c3d4-0004-4000-8000-000000000004', NOW() - INTERVAL '3 hours',      NOW() + INTERVAL '21 hours'),
('b1000000-0004-4000-8000-000000000003', 'a1b2c3d4-0004-4000-8000-000000000004', NOW() - INTERVAL '2 hours',      NOW() + INTERVAL '22 hours'),
('b1000000-0004-4000-8000-000000000004', 'a1b2c3d4-0004-4000-8000-000000000004', NOW() - INTERVAL '1 hour',       NOW() + INTERVAL '23 hours'),
('b1000000-0004-4000-8000-000000000005', 'a1b2c3d4-0004-4000-8000-000000000004', NOW() - INTERVAL '30 minutes',   NOW() + INTERVAL '23 hours 30 minutes'),
-- Zephyr (5 assets)
('b1000000-0005-4000-8000-000000000001', 'a1b2c3d4-0005-4000-8000-000000000005', NOW() - INTERVAL '4 hours',      NOW() + INTERVAL '20 hours'),
('b1000000-0005-4000-8000-000000000002', 'a1b2c3d4-0005-4000-8000-000000000005', NOW() - INTERVAL '3 hours',      NOW() + INTERVAL '21 hours'),
('b1000000-0005-4000-8000-000000000003', 'a1b2c3d4-0005-4000-8000-000000000005', NOW() - INTERVAL '2 hours',      NOW() + INTERVAL '22 hours'),
('b1000000-0005-4000-8000-000000000004', 'a1b2c3d4-0005-4000-8000-000000000005', NOW() - INTERVAL '1 hour',       NOW() + INTERVAL '23 hours'),
('b1000000-0005-4000-8000-000000000005', 'a1b2c3d4-0005-4000-8000-000000000005', NOW() - INTERVAL '30 minutes',   NOW() + INTERVAL '23 hours 30 minutes')
ON CONFLICT (id) DO UPDATE SET
  created_at = EXCLUDED.created_at,
  expires_at = EXCLUDED.expires_at;

-- ===== STORY MEDIA (1 media per story) =====
INSERT INTO story_media (id, story_id, media_url, media_type, duration, sort_order) VALUES
-- Luna
('c1000000-0001-4000-8000-000000000001', 'b1000000-0001-4000-8000-000000000001', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/luna/luna-1.webp',      'image', 5,  0),
('c1000000-0001-4000-8000-000000000002', 'b1000000-0001-4000-8000-000000000002', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/luna/luna-2.webp',      'image', 5,  0),
('c1000000-0001-4000-8000-000000000003', 'b1000000-0001-4000-8000-000000000003', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/luna/luna-3.webp',      'image', 5,  0),
('c1000000-0001-4000-8000-000000000004', 'b1000000-0001-4000-8000-000000000004', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/luna/luna-4.webp',      'image', 5,  0),
('c1000000-0001-4000-8000-000000000005', 'b1000000-0001-4000-8000-000000000005', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/luna/luna-5.webp',      'image', 5,  0),
('c1000000-0001-4000-8000-000000000006', 'b1000000-0001-4000-8000-000000000006', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/luna/luna-video-1.mp4', 'video', 10, 0),
-- Kai
('c1000000-0002-4000-8000-000000000001', 'b1000000-0002-4000-8000-000000000001', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/kai/kai-1.webp',        'image', 5,  0),
('c1000000-0002-4000-8000-000000000002', 'b1000000-0002-4000-8000-000000000002', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/kai/kai-2.webp',        'image', 5,  0),
('c1000000-0002-4000-8000-000000000003', 'b1000000-0002-4000-8000-000000000003', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/kai/kai-3.webp',        'image', 5,  0),
('c1000000-0002-4000-8000-000000000004', 'b1000000-0002-4000-8000-000000000004', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/kai/kai-4.webp',        'image', 5,  0),
('c1000000-0002-4000-8000-000000000005', 'b1000000-0002-4000-8000-000000000005', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/kai/kai-video-1.mp4',   'video', 10, 0),
('c1000000-0002-4000-8000-000000000006', 'b1000000-0002-4000-8000-000000000006', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/kai/kai-video-2.mp4',   'video', 10, 0),
-- Nova
('c1000000-0003-4000-8000-000000000001', 'b1000000-0003-4000-8000-000000000001', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/nova/nova-1.webp',      'image', 5,  0),
('c1000000-0003-4000-8000-000000000002', 'b1000000-0003-4000-8000-000000000002', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/nova/nova-2.webp',      'image', 5,  0),
('c1000000-0003-4000-8000-000000000003', 'b1000000-0003-4000-8000-000000000003', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/nova/nova-3.webp',      'image', 5,  0),
('c1000000-0003-4000-8000-000000000004', 'b1000000-0003-4000-8000-000000000004', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/nova/nova-4.webp',      'image', 5,  0),
('c1000000-0003-4000-8000-000000000005', 'b1000000-0003-4000-8000-000000000005', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/nova/nova-video-1.mp4', 'video', 10, 0),
('c1000000-0003-4000-8000-000000000006', 'b1000000-0003-4000-8000-000000000006', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/nova/nova-video-2.mp4', 'video', 10, 0),
-- Ember
('c1000000-0004-4000-8000-000000000001', 'b1000000-0004-4000-8000-000000000001', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/ember/ember-1.webp',      'image', 5,  0),
('c1000000-0004-4000-8000-000000000002', 'b1000000-0004-4000-8000-000000000002', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/ember/ember-2.webp',      'image', 5,  0),
('c1000000-0004-4000-8000-000000000003', 'b1000000-0004-4000-8000-000000000003', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/ember/ember-3.webp',      'image', 5,  0),
('c1000000-0004-4000-8000-000000000004', 'b1000000-0004-4000-8000-000000000004', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/ember/ember-video-1.mp4', 'video', 10, 0),
('c1000000-0004-4000-8000-000000000005', 'b1000000-0004-4000-8000-000000000005', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/ember/ember-video-2.mp4', 'video', 10, 0),
-- Zephyr
('c1000000-0005-4000-8000-000000000001', 'b1000000-0005-4000-8000-000000000001', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/zephyr/zephyr-1.webp',      'image', 5,  0),
('c1000000-0005-4000-8000-000000000002', 'b1000000-0005-4000-8000-000000000002', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/zephyr/zephyr-2.webp',      'image', 5,  0),
('c1000000-0005-4000-8000-000000000003', 'b1000000-0005-4000-8000-000000000003', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/zephyr/zephyr-3.webp',      'image', 5,  0),
('c1000000-0005-4000-8000-000000000004', 'b1000000-0005-4000-8000-000000000004', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/zephyr/zephyr-video-1.mp4', 'video', 10, 0),
('c1000000-0005-4000-8000-000000000005', 'b1000000-0005-4000-8000-000000000005', 'https://wiyiltwuiplfenhbpgsg.supabase.co/storage/v1/object/public/stories/zephyr/zephyr-video-2.mp4', 'video', 10, 0)
ON CONFLICT (id) DO NOTHING;
