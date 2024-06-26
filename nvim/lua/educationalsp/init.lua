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
	callback = function(event)
		vim.lsp.buf_attach_client(0, client)
		vim.keymap.set("n", "K", vim.lsp.buf.hover, { buffer = event.buf })
		vim.keymap.set("n", "gd", vim.lsp.buf.definition, { buffer = event.buf })
		vim.keymap.set("n", "<leader>ca", vim.lsp.buf.code_action, { buffer = event.buf })
	end,
})
