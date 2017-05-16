#!/usr/bin/env tarantool

box.cfg{
    -- log_level
    -- 1 – SYSERROR
    -- 2 – ERROR
    -- 3 – CRITICAL
    -- 4 – WARNING
    -- 5 – INFO
    -- 6 – DEBUG
    log_level = 5,

    --logger = 'tarantool.txt',

    listen = 3301,

    slab_alloc_arena = 1,
    -- wal_dir='xlog',
    -- snap_dir='snap',
}

log = require('log')

log.info('Info log enabled')

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
-- subscriptions (observe)
--

s = box.schema.create_space(prefix.."subscriptions", {
    if_not_exists=true,
})

s:create_index('primary', {
    if_not_exists=true,
    type = 'tree',
    -- thread_id user_id
    parts = {2, 'string', 1, 'string'},
})

-- for ThreadObservers
s:create_index('subs_thread_id_idx', {
    if_not_exists=true,
    type = 'tree',
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

s:create_index('threadline_idx', {
    if_not_exists=true,
    type = 'tree',
    unique = false,
    -- thread_id, created_at
    parts = {2, 'string', 3, 'integer'},
})

function by_last_ts(threrad_id, timestamp, limit)
  local tuples = {}
  local count = 0

  local events_space = prefix.."events"
  for _, tuple in box.space.chronograph_events.index.threadline_idx:pairs({threrad_id}, {iterator = box.index.REQ}) do
    if tuple[3] >= timestamp then
      count = count + 1
      table.insert(tuples, tuple)
      if count == limit then break end
    end

  end
  if next(tuples) == nil then
    return
  end
  return unpack(tuples)
end

--
-- users
--

s = box.schema.create_space(prefix.."users", {
    if_not_exists=true,
})

s:create_index('primary', {
    if_not_exists=true,
    type = 'hash',
    unique = true,
    parts = {1, 'string'},
})

--
-- related data
--

s = box.schema.create_space(prefix.."related", {
    if_not_exists=true,
})

s:create_index('primary', {
    if_not_exists=true,
    type = 'tree',
    unique = true,
    -- user_id, event_id
    parts = {1, 'string', 2, 'string'},
})

--
-- threadline
--

s = box.schema.create_space(prefix.."threadline", {
    if_not_exists=true,
})

-- unused index for primary
s:create_index('primary', {
    if_not_exists=true,
    type = 'hash',
    unique = true,
    -- user_id, thread_id
    parts = {1, 'string', 2, 'string', 4, 'string'},
})

s:create_index('threadline_real_idx', {
    if_not_exists=true,
    type = 'tree',
    unique = false,
    -- user_id, thread_id, created_at
    parts = {1, 'string', 2, 'string', 3, 'integer'},
})

-- создает записи в threadline для всех подписчиков трэда
function new_event_in_threadline(thread_id, created_at, event_id)
  for _, tuple in box.space.chronograph_subscriptions.index.primary:pairs({thread_id}, {iterator = box.index.REQ}) do
    box.space.chronograph_threadline:insert({tuple[1], thread_id, created_at, event_id})
  end
end


-- возвращает threadline
function threadline(user_id, thread_id, limit, offset)
  local events = {}
  local count = 0

-- receive events ids
  for _, tuple in box.space.chronograph_threadline.index.threadline_real_idx:pairs({user_id, thread_id}, {iterator = box.index.REQ, offset = offset, limit = limit}) do
      count = count + 1
      event = box.space.chronograph_events.index.primary:select({tuple[4]})
      table.insert(events, unpack(event))
      if count == limit then break end
  end
  -- возвращаем nil, если результат пустой
  if next(events) == nil then
    return
  end

  return unpack(events)
end
