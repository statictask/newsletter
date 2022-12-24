## Example Form HTML page

This page includes a very simple form for creating new subscriptions.
The code will send requests to the local backend, so before using it
make sure you have the server running.

	make run

Then, after creating your first project manually ([check this for more
information](../examples/README.md)), edit the `newsletter.js` file
to include the correct project id.

Now open the `index.html` file on any browser. If you're using linux
and have `xdg` available, you can run

	xdg-open ./form/index.html
