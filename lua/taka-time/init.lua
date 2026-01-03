-- init.lua ok my firat plugin
-- redir > nvim_logs.txt | silent messages | redir END
-- lua/taka-time/init.lua
local M = {}
local config = require("taka-time.config")
local utils = require("taka-time.utils")
local core = require("taka-time.core")

function M.setup(opts)
	-- 1. Setup Config
	config.setup(opts)

	-- 2. Ensure Binary Exists (Download if needed)
	-- We run this in a pcall to prevent crashing if internet is down
	pcall(utils.ensure_binary)

	-- 3. Reset Timer
	-- Accessing a private state variable isn't ideal, but for v1 it's fine.
	-- Ideally 'core' would expose a reset function.

	-- 4. Setup Autocommands
	local group = vim.api.nvim_create_augroup("TakaTimeGroup", { clear = true })

	vim.api.nvim_create_autocmd("BufWritePost", {
		group = group,
		pattern = "*",
		callback = core.on_save,
	})

	vim.api.nvim_create_autocmd("VimLeavePre", {
		group = group,
		callback = core.on_exit,
	})
end

return M
