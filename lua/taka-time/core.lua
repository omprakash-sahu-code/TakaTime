local M = {}
local config = require("taka-time.config")
local utils = require("taka-time.utils")

-- STATE
local state = {
    last_active = os.time(),
    pending_duration = 0,
    job_id = 0
}

-- Internal: The actual upload logic
local function attempt_upload()
    if state.job_id ~= 0 then return end -- Busy
    if state.pending_duration < config.options.debounce_seconds then return end -- Not enough data

    -- Snapshot data to send
    local time_to_send = state.pending_duration
    state.pending_duration = 0

    -- Prepare args
    local file_path = vim.fn.expand('%:p')
    local project = vim.fn.fnamemodify(vim.fn.getcwd(), ':t')
    local ext = vim.fn.fnamemodify(file_path, ':e')
    if ext == "" then ext = "text" end

    local cmd = {
        utils.get_binary_path(),
        "-uri", config.options.mongo_uri,
        "-project", project,
        "-language", ext,
        "-file", file_path,
        "-duration", tostring(time_to_send)
    }

    if config.options.debug then print("[Taka] Syncing " .. time_to_send .. "s...") end

    state.job_id = vim.fn.jobstart(cmd, {
        on_exit = function(_, code)
            state.job_id = 0
            if code ~= 0 then
                -- Failure: Put time back in bucket
                state.pending_duration = state.pending_duration + time_to_send
                if config.options.debug then print("[Taka] Failed. Retrying.") end
            elseif config.options.debug then
                print("[Taka] Success.")
            end
        end,
    })
end

-- Public: Called on :w
function M.on_save()
    local now = os.time()
    state.pending_duration = state.pending_duration + (now - state.last_active)
    state.last_active = now
    
    attempt_upload()
end

-- Public: Called on Exit
function M.on_exit()
    local now = os.time()
    state.pending_duration = state.pending_duration + (now - state.last_active)

    if state.pending_duration > 0 then
        print("[Taka] Uploading final data...")
        -- Blocking call for exit
        vim.fn.system({
            utils.get_binary_path(),
            "-uri", config.options.mongo_uri,
            "-project", "unknown",
            "-language", "text",
            "-file", "closing_session",
            "-duration", tostring(state.pending_duration)
        })
    end
end

return M
