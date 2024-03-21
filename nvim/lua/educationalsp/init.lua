local client = vim.lsp.start_client({
	name = "educationalsp",
	cmd = { "/home/zi/personal/educationalsp/build/educationalsp" },
})

if not client then
	vim.notify("lsp-dev client configuration error")
	return
end

vim.api.nvim_create_autocmd("FileType", {
	pattern = "markdown",
	callback = function()
		vim.lsp.buf_attach_client(0, client)
	end,
})
