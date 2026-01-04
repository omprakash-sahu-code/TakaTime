local M = {}

-- Save secrets to ~/.local/share/nvim/taka_data.json
local data_path = vim.fn.stdpath("data") .. "/taka_data.json"

function M.save_secret(uri)
	local file = io.open(data_path, "w")
	if file then
		file:write(vim.json.encode({ mongo_uri = uri }))
		file:close()
		print("✅ TakaTime: Mongo URI saved successfully!")
	else
		print("❌ TakaTime Error: Could not save secret to " .. data_path)
	end
end

function M.load_secret()
	local file = io.open(data_path, "r")
	if file then
		local content = file:read("*a")
		file:close()
		-- Safely decode JSON
		local ok, data = pcall(vim.json.decode, content)
		if ok and data.mongo_uri then
			return data.mongo_uri
		end
	end
	return nil
end

return M
