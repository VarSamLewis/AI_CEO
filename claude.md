# MealPrep Buddy - Building in Public

## About the Founder

I'm a non-technical startup founder who's passionate about solving everyday problems that people actually face. For months, I've been watching friends and family struggle with the same issue: they want to eat healthier and save money, but existing meal planning apps are either too complicated, require ingredients nobody has, or feel like homework.

I'm not a developer, but I know a real problem when I see one. That's why I brought on my founding engineer to turn this vision into reality.

## The Problem

Everyone I know wants to eat healthier and save money, but meal planning is broken:
- Apps are too complicated with calorie tracking and nutrition analytics nobody asked for
- Recipe suggestions require 47 obscure ingredients you'll use once
- They don't account for what you already have at home
- People download them, use them once, and forget about them

## The Solution: MealPrep Buddy

An AI-powered meal planning app that feels like having a really organized friend who's good at cooking help you out.

**Core concept:**
1. Tell the app what ingredients you have at home
2. Set simple preferences (dietary restrictions, cooking time)
3. Get practical meal suggestions that use what you own
4. Save favorites for later

No complicated tracking. No nutrition guilt-tripping. Just helpful meal ideas.

## What We're Building (MVP)

After discussion with my founding engineer, here's what we decided to build first:

### Week 1-2: Core Features
- **User accounts** - Login/signup system so preferences persist
- **Preference settings** - Dietary restrictions, cooking time preferences, ingredient dislikes
- **Ingredient input** - Manual text entry (skipping photo recognition for MVP)
- **AI meal generation** - Feed user inputs to LLM with custom prompts to generate 3 meal ideas with cooking instructions
- **Saved favorites** - Simple database to store meals users love

### Week 3-4: Testing Phase
- Get 10-15 real people using it
- Each person generates meals 2-3 times per week
- Measure: Do they come back? That's the key question.

## Decision Log

### Decision 1: Manual Input vs Photo Recognition
**Date:** November 11, 2025
**Decision:** Start with manual ingredient input
**Reasoning:** Photo recognition sounds cool but adds complexity. We need to validate if people even want ingredient-based meal suggestions before building fancy features. We can add photos later if the core idea works.

### Decision 2: Function Over Form
**Date:** November 11, 2025
**Decision:** Build functional but not beautiful for MVP
**Reasoning:** We have 2 weeks. Pretty designs can wait. We need to test the core hypothesis: will people use an app that turns their pantry into meal ideas? Aesthetics are a Week 5 problem.

### Decision 3: AI Cost Structure
**Date:** November 11, 2025
**Research:** Consulted with technical consultant Gloria Sonnet
**Finding:** $0.0027 - $0.008 per meal generation request
**Impact:** Even at 100 meals per user per month, that's less than $1 in costs. This is totally viable for free testing and won't break our budget.

### Decision 4: Core Validation Metric
**Date:** November 11, 2025
**Decision:** Measure repeat usage, not signups
**Reasoning:** Anyone will try a free app once. The real question is: do they come back? If 10 people test it and 8 of them use it multiple times per week, we have something. If they use it once and forget, the idea doesn't work.

### Decision 5: Monetization Strategy
**Date:** November 27, 2025
**Decision:** Free testing with usage caps (Week 3-4), then immediate paid launch (Week 5+)

**Testing Phase (Week 3-4):**
- Each tester gets 20 free meal generations
- Total cost to founder: ~$5 (acceptable for validation)
- Tell testers upfront: "This will cost $10/month after testing. Would you pay that?"
- Tests willingness to pay WITHOUT payment friction affecting honest feedback

**Launch Phase (Week 5+):**
- **Pricing:** $10/month for unlimited meal generations
- **Payment:** Stripe (2.9% + $0.30 per transaction)
- **Unit Economics:** Power users (100 meals/month) cost <$1 in LLM fees, giving 90% margin
- **Positioning:** Comparable to Netflix/Spotify pricing people understand

**Reasoning:** No funding means we can't subsidize usage long-term. But $5 to validate with 15 real users is cheap market research. If they say "yes, I'd pay $10/month" and demonstrate repeat usage, we have product-market fit. Then we charge from day one of public launch. Usage caps during testing prevent cost explosion.

**Technical Requirements:** Need usage tracking per user and rate limiting system to enforce 20-meal cap during testing.

## Technical Stack

- Web application with user authentication
- Database for user preferences and saved meals
- LLM integration for meal generation (custom system prompts)
- Simple, functional UI

## Timeline

**Week 1-2:** Build MVP
- Engineer handles all development
- Founder writes copy for the app and prepares test user recruitment

**Week 3:** Initial Testing
- Recruit 10-15 friends/family who expressed interest
- Goal: Each person uses it 2-3 times minimum
- Collect feedback informally

**Week 4:** Evaluation & Decision
- Analyze usage patterns
- Talk to test users
- Decide: Is this worth building further?
- Key question: Would they pay $5-10/month for this?

## What We're NOT Building Yet

These are ideas for later, after we validate the core concept:
- Photo recognition of fridge/pantry
- Weekly meal planning calendar
- Shopping list optimization
- Nutrition tracking
- Recipe ratings and reviews
- Social features
- Mobile apps
- Integration with grocery delivery

## The Bet We're Making

We believe that people don't need another complicated meal planning app. They need something dead simple that just helps them answer the question: "What should I cook with what I have?"

If we're right, this could help thousands of people eat better, waste less food, and save money. If we're wrong, we'll learn that in 4 weeks and move on.

## Progress Updates

### November 11, 2025 - Day 1
- Finalized MVP scope with founding engineer
- Decided on 2-week development timeline
- Confirmed technical feasibility and cost structure
- Engineer beginning development

### November 27, 2025 - Week 1 Progress
- Backend development mostly complete
- Finalized monetization strategy: Free testing with usage caps, then $10/month subscription
- Next: Implement usage tracking and rate limiting (20 meals/user during testing)
- Next: Stripe payment integration for post-testing launch

---

**Current Status:** In Development (Week 1)
**Last Updated:** November 27, 2025

*This is a living document. Check back for updates as we build.*
