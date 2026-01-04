local M = {}
local config = require("taka-time.config")
local utils = require("taka-time.utils")
local core = require("taka-time.core")
local storage = require("taka-time.storage")

function M.setup(opts)
	-- 1. Setup Config
	config.setup(opts)

	-- 2. AUTOMATICALLY LOAD SECRET (Crucial Step!)
	-- If user didn't provide URI in config, try to load from storage
	if config.options.mongo_uri == "" then
		local saved_uri = storage.load_secret()
		if saved_uri then
			config.options.mongo_uri = saved_uri
		end
	end

	-- 3. Create Commands
	vim.api.nvim_create_user_command("TakaInit", function()
		-- Use inputsecret to hide characters (******)
		local uri = vim.fn.inputsecret("Enter your Mongo URI: ")
		if uri and uri ~= "" then
			storage.save_secret(uri)
			config.options.mongo_uri = uri -- Update runtime config immediately
		else
			print("❌ TakaTime: No URI Entered")
		end
	end, {})

	vim.api.nvim_create_user_command("TakaStatus", function()
		if config.options.mongo_uri and config.options.mongo_uri ~= "" then
			print("✅ TakaTime is configured and running.")
		else
			print("⚠️ TakaTime is NOT configured. Run :TakaInit")
		end
	end, {})

	-- 4. Start Timer
	core.start_timer()

	-- 5. Ensure Binary Exists (Download if needed)
	-- Using pcall ensures we don't crash if internet is down
	pcall(utils.ensure_binary)

	-- 6. Setup Autocommands
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
