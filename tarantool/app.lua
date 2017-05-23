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

-- подписчики
s:create_index('primary', {
    if_not_exists=true,
    type = 'tree',
    -- thread_id user_id
    parts = {2, 'string', 1, 'string'},
})

-- подписки
s:create_index('subs_thread_id_idx', {
    if_not_exists=true,
    type = 'tree',
    --  user_id thread_id,
    parts = {1, 'string', 4, 'integer', 2, 'string'},
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

s:create_index('_updated_idx', {
    if_not_exists=true,
    type = 'tree',
    unique = true,
    -- user_id, event_id
    parts = {1, 'string', 4, 'string'},
})

-- создает записи в threadline для всех подписчиков трэда
function new_event_in_threadline(thread_id, created_at, event_id)
  for _, tuple in box.space.chronograph_subscriptions.index.primary:pairs({thread_id}, {iterator = box.index.REQ}) do
    box.space.chronograph_threadline:insert({tuple[1], thread_id, created_at, event_id})
  end
end


-- возвращает threadline
function threadline(user_id, thread_id, lim, off)
  local events = {}

-- receive events ids
  for _, tuple in pairs(box.space.chronograph_threadline.index.threadline_real_idx:select({user_id, thread_id}, {iterator = box.index.REQ, offset = off, limit = lim})) do
      event = box.space.chronograph_events.index.primary:get({tuple[4]})
      table.insert(events, event)
  end
  -- возвращаем nil, если результат пустой
  if next(events) == nil then
    return
  end

  return unpack(events)
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

function count_events(user_id, thread_id, time, limit, offset)
  cnt = 0
  if threadline_enabled(thread_id) then
    for _, tuple in pairs(box.space.chronograph_threadline.index.threadline_real_idx:select({user_id, thread_id}, {iterator = box.index.REQ, limit = limit, offset = offset})) do
        if tuple[3] > time then
          cnt = cnt + 1
        else
          break
        end
    end
  else
    for _, tuple in pairs(box.space.chronograph_events.index.threadline_idx:select({threrad_id}, {iterator = box.index.REQ, limit = limit, offset = offset})) do
      if tuple[3] > time then
        cnt = cnt + 1
      else
        break
      end
    end
  end
  return cnt
end
