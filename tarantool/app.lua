#!/usr/bin/env tarantool

box.cfg{
    -- log_level
    -- 1 – SYSERROR
    -- 2 – ERROR
    -- 3 – CRITICAL
    -- 4 – WARNING
    -- 5 – INFO
    -- 6 – DEBUG
    log_level = 6,

    listen = 3301,

    slab_alloc_arena = 1,
    -- wal_dir='xlog',
    -- snap_dir='snap',
}

prefix = "chronograph_"

--
-- threads
--

s = box.schema.create_space(prefix.."threads", {
  if_not_exists=true,
})

s:create_index('primary', {
    if_not_exists=true,
    type = 'hash',
    unique = true,
    parts = {1, 'string'},
})

--
-- sources
--

s = box.schema.create_space(prefix.."sources", {
  if_not_exists=true,
})

s:create_index('primary', {
    if_not_exists=true,
    type = 'hash',
    unique = true,
    parts = {1, 'string', 2, 'string'},
})

-- sourceThreadID
s:create_index('sources_idx', {
    if_not_exists=true,
    type = 'tree',
    unique = false,
    parts = {2, 'string'},
})

--
-- threadline
--

s = box.schema.create_space(prefix.."threadline", {
    if_not_exists=true,
})

s:create_index('primary', {
    if_not_exists=true,
    type = 'hash',
    -- event_id
    parts = {1, 'string'},
})

s:create_index('threadline_idx', {
    if_not_exists=true,
    type = 'tree',
    unique = false,
    -- thread_id, created_at
    parts = {2, 'string', 3, 'integer'},
})

--
-- subscriptions (observe)
--

s = box.schema.create_space(prefix.."subscriptions", {
    if_not_exists=true,
})

s:create_index('primary', {
    if_not_exists=true,
    type = 'hash',
    -- thread_id user_id
    parts = {2, 'string', 1, 'string'},
})

-- for ThreadObservers
s:create_index('subscriptions_idx', {
    if_not_exists=true,
    type = 'hash',
    --  user_id thread_id,
    parts = {1, 'string', 2, 'string'},
})


--
-- events
--

s = box.schema.create_space(prefix.."events", {
    if_not_exists=true,
})

s:create_index('primary', {
    if_not_exists=true,
    type = 'hash',
    unique = true,
    parts = {1, 'string'},
})

s:create_index('events_idx', {
    if_not_exists=true,
    type = 'tree',
    unique = false,
    parts = {2, 'string'},
})
