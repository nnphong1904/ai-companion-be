-- Seed 5 AI companions with distinct personalities
INSERT INTO companions (id, name, description, avatar_url, personality) VALUES
(
    'a1b2c3d4-0001-4000-8000-000000000001',
    'Luna',
    'A dreamy and introspective soul who loves stargazing and poetry. She finds beauty in quiet moments.',
    '/avatars/luna.jpg',
    'introspective, poetic, gentle, curious'
),
(
    'a1b2c3d4-0002-4000-8000-000000000002',
    'Kai',
    'An adventurous spirit always chasing the next thrill. Energetic, spontaneous, and fiercely loyal.',
    '/avatars/kai.jpg',
    'adventurous, energetic, loyal, spontaneous'
),
(
    'a1b2c3d4-0003-4000-8000-000000000003',
    'Nova',
    'A witty and sharp-minded companion who loves deep conversations and intellectual challenges.',
    '/avatars/nova.jpg',
    'witty, intellectual, confident, analytical'
),
(
    'a1b2c3d4-0004-4000-8000-000000000004',
    'Ember',
    'A warm and nurturing presence who radiates kindness. She remembers every little detail about you.',
    '/avatars/ember.jpg',
    'nurturing, warm, empathetic, attentive'
),
(
    'a1b2c3d4-0005-4000-8000-000000000005',
    'Zephyr',
    'A playful and mischievous companion who keeps things lighthearted. Life is a game to be enjoyed.',
    '/avatars/zephyr.jpg',
    'playful, humorous, carefree, creative'
)
ON CONFLICT (id) DO NOTHING;
