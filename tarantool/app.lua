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

function by_last_ts(threrad_id, timestamp)
  local tuples = {}

  local events_space = prefix.."events"
  for _, tuple in box.space.chronograph_events.index.threadline_idx:pairs({threrad_id}, {iterator = box.index.Req}) do
    if tuple[3] >= timestamp then
      table.insert(tuples,1,tuple)
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

// создает записи в threadline для всех подписчиков трэда
function new_event_in_threadline(thread_id, created_at, event_id)
  for _, tuple in box.space.chronograph_subscriptions.index.primary:pairs({threrad_id}, {iterator = box.index.Req}) do
    log.info('Info insert threadline event to user %s (thread %s, event %s)', tuple[1], thread_id, event_id)
    box.space.chronograph_threadline:insert({tuple[1], thread_id, created_at, event_id})
  end
end

function threadline_enabled(thread_id)
  threads = box.space.chronograph_threads.index.primary:select({thread_id})
  if next(threads)== nil then
    return false
  end
  if threads[1][2] then
    return true
  end
  return false
end

// возвращает threadline
function threadline_by_last_ts(user_id, thread_id, timestamp)
  local tuples = {}

-- receive events ids
  for _, tuple in box.space.chronograph_threadline.index.threadline_real_idx:pairs({user_id, thread_id}, {iterator = box.index.Req}) do
    if tuple[3] >= timestamp then
      table.insert(tuples,1,tuple)
    end
  end
  if next(tuples) == nil then
    return
  end

-- recieve events
  local events = {}
  for _, tuple in pairs(tuples) do
    event = box.space.chronograph_events.index.primary:select({tuple[4]})
    table.insert(events, 1, unpack(event))
  end
  if next(events) == nil then
    return
  end
  return unpack(events)
end


// проверяет включен ли у трэда threadline и возвращает соответствующий результат
function get_threadline(user_id, thread_id, timestatmp)
  if threadline_enabled(thread_id) then
    return threadline_by_last_ts(user_id, thread_id, timestatmp)
  end
  return by_last_ts(thread_id, timestatmp)
end
