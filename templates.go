package main

const defaultPage = `<html>
<form action="/login_page" method="post">
	<button>Login</button>
</form>
<form action="/register_page" method="post">
	<button>Register</button>
</form>
</html>`

const loginPage = `<html>
<form action="/login">
	<label for="username">Username: </label>
	<input type="text" id="username" name="username" value="" pattern="[a-zA-Z0-9]+" required><br>
	<label for="password">Password: </label>
	<input type="password" id="password" name="password" value="" pattern="[a-zA-Z0-9]+" required><br>
	<input type="submit" value="Login">
</form>
</html>`

const loginError = `<html>%s</html>`

const registerPage = `<html>
<form action="/register">
	<label for="username">Username: </label>
	<input type="text" id="username" name="username" value="" pattern="[a-zA-Z0-9]+" required><br>
	<label for="password">Password: </label>
	<input type="password" id="password" name="password" value="" pattern="[a-zA-Z0-9]+" required><br>
	<input type="radio" id="employer" name="role" value="employer" required>
	<label for="employer">Employer</label>
	<input type="radio" id="worker" name="role" value="worker">
	<label for="employer">Worker</label>
	<input type="submit" value="Login">
</form>
</html>`

const homeWorker = `<html>
<form action="/me" method="post">
	<button>My profile</button>
</form>
<form action="/tasks" method="post">
	<button>My tasks</button>
</form>
</html>`

const homeEmployer = `<html>
<form action="/me" method="post">
	<button>My profile</button>
</form>
<form action="/tasks" method="post">
	<button>My tasks</button>
</form>
<form action="/start" method="post">
	<button>Play the game</button>
</form>
</html>`
