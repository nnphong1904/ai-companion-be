# Full Stack Engineer Take-Home Assessment – Project Specification

---

# 1️⃣ Assessment Requirements (From PDF)

## Context

Build a fully functional **AI Companion app** that includes:

- Stories feature (similar behavior to Instagram Stories)
- 2 additional feature-complete items
- Optimized for both mobile and desktop
- Focus on product depth and feature creativity
- Clean architecture and scalability

The goal is to evaluate:

- Product-level thinking
- UX polish
- Scalability decisions
- Code quality
- Clean separation between frontend and backend
- Creativity in custom features

---

# 2️⃣ Technical Requirements

## Frontend

- NextJS (App Router)
- TypeScript expected
- Mobile + Desktop optimized
- Clean UI component separation

## Backend

- Go
- REST API
- Clean separation of handlers / services / repositories

## Database

- Any (PostgreSQL recommended)
- Must justify choice in README
- Proper indexing required

## Storage

- Any cloud storage solution
- Used for images/videos

## Deployment

- Vercel (frontend) or equivalent
- Backend hosted and live
- Public link required

---

# 3️⃣ Asset Requirements (Very Important)

- Minimum 5 companions
- Each companion must have:
  - At least 1 photo OR video
  - At least 1 Story
- At least one Story must contain 4 slides
- Stories must include BOTH photos and videos
- Avoid generic stock images
- Assets must feel believable and emotionally aligned
- 9:16 vertical format for Stories

---

# 4️⃣ Core App Concept

This is NOT a chatbot.

It is a **Reactive AI Companion Social App** where:

- AI companions post Stories
- They react emotionally to user behavior
- Relationship evolves over time
- Important emotional moments are stored and resurfaced

The app must create a “Social Loop” that increases daily engagement.

---

# 5️⃣ Required Base Feature

## Stories System

- Horizontal story feed
- Full-screen story modal
- Animated progress bars
- Tap left/right navigation
- React buttons (❤️ 😢 😍 😡)
- Preloaded media
- At least one 4-slide story sequence

---

# 6️⃣ Feature Complete Item #1

## Dynamic Mood & Relationship System

### Purpose

Make companions emotionally reactive.

### Includes

- mood_score (numeric)
- relationship_score (numeric)
- Mood labels:
  - <20 → Distant
  - 20–50 → Neutral
  - 50–80 → Happy
  - 80+ → Attached
- Relationship progression display
- Mood changes based on:
  - Chat interactions
  - Story reactions
  - Time decay (inactivity reduces mood)

### Backend Model

RelationshipState:

- user_id
- companion_id
- mood_score
- relationship_score
- last_interaction
- updated_at

Composite unique index on:
(user_id, companion_id)

Mood recalculated on read (no background jobs).

This is a feature-complete system:

- Dedicated UI
- Backend logic
- Affects whole product experience

---

# 7️⃣ Feature Complete Item #2

## Emotional Memory Timeline

### Purpose

Provide emotional continuity.

This is NOT chat history.
It is curated meaningful shared moments.

### Includes

- Dedicated Memory page
- Timeline UI
- Memory card:
  - content
  - date
  - optional tag
- "Mark as Memory" action in chat
- Delete memory
- Optional pin feature

### Backend Model

Memories:

- id
- user_id
- companion_id
- content
- tag (optional)
- created_at

### Optional Resurfacing

When fetching Stories:
If memory older than X days exists,
Insert dynamic slide:

"I was thinking about when you said..."

No complex AI logic required.

This qualifies as a feature-complete system because:

- Own UI
- Own data model
- Own interaction logic
- Provides independent user value

---

# 8️⃣ Authentication (Expected but Not Counted as Feature)

- Email + Password
- bcrypt hashing
- JWT
- Protected routes
- Required for user-specific relationship + memories

Keep simple and minimal.

---

# 9️⃣ Required Frontend Routes

- /login
- /signup
- /onboarding
- /dashboard
- /companions/[id]
- /chat/[id]
- /memories/[id]

---

# 🔟 Onboarding Flow

1. Welcome screen
2. Select one of 5 companions
3. Create initial RelationshipState
4. Redirect to dashboard

Must feel intentional and polished.

---

# 1️⃣1️⃣ Core Database Entities

- Users
- Companions
- Stories
- StoryMedia
- Messages
- RelationshipState
- Memories

Indexes required:

- user_id
- companion_id
- created_at

Avoid N+1 queries.
Use pagination for messages.

---

# 1️⃣2️⃣ Architecture Constraints

## Backend

- Clean folder separation:
  - handlers
  - services
  - repositories
  - models
  - middleware
- REST only
- No WebSockets
- No cron jobs
- Proper error handling
- Cursor or offset pagination
- CORS enabled

## Frontend

- App Router
- Component modularization
- Mobile-first design
- Smooth story transitions
- Animated progress bars
- Dark emotional theme
- Separation of UI and data logic

---

# 1️⃣3️⃣ Design Philosophy

App should feel:

- Emotional
- Reactive
- Social
- Persistent
- Immersive

Not:

- Generic CRUD
- Plain chatbot
- Corporate admin tool

---

# 1️⃣4️⃣ Evaluation Criteria

The submission should demonstrate:

- Product thinking
- Creativity
- Feature depth
- Scalable schema design
- Query optimization
- UX polish
- Clean architecture
- Proper TypeScript usage
- API/UI separation

README must justify:

- Feature choices
- Design decisions
- Architecture choices
- Database selection
- Competitor research

