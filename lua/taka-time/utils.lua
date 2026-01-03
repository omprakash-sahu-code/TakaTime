local M = {}
local config = require("taka-time.config")

-- Helper: Get path to the Go binary
function M.get_binary_path()
	local plugin_root = vim.fn.fnamemodify(debug.getinfo(1).source:sub(2), ":h:h:h")
	return plugin_root .. "/taka-upload"
end

-- Helper: Detect OS and Architecture
local function get_os_info()
	local uname = vim.loop.os_uname()
	local os_name = uname.sysname:lower()
	local arch = uname.machine:lower()

	if os_name == "linux" then
		os_name = "linux"
	elseif os_name == "darwin" then
		os_name = "darwin"
	else
		return nil, nil
	end -- Windows requires extra logic later

	if arch == "x86_64" then
		arch = "amd64"
	end
	-- arm64 usually stays arm64, but sometimes aarch64
	if arch == "aarch64" then
		arch = "arm64"
	end

	return os_name, arch
end

-- Public: Check if binary exists, if not, download it
function M.ensure_binary()
	local bin_path = M.get_binary_path()

	-- If file exists, we are good
	if vim.fn.filereadable(bin_path) == 1 then
		return
	end

	local os_name, arch = get_os_info()
	if not os_name then
		print("[Taka] Auto-install not supported for this OS. Please build manually.")
		return
	end

	-- In lua/taka-time/utils.lua
	local version = config.options.binary_version
	local url = string.format(
		"https://github.com/Rtarun3606k/TakaTime/releases/download/%s/taka-upload-%s-%s",
		version,
		os_name,
		arch
	)

	print("[Taka] Installing binary (" .. version .. ")...")

	-- Download
	vim.fn.system({ "curl", "-L", "-o", bin_path, url })
	-- Make executable
	vim.fn.system({ "chmod", "+x", bin_path })

	print("[Taka] Installation complete.")
end

return M
