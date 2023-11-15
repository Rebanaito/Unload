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
<form action="/register_employer" method="post">
	<button>Register as an employer</button>
</form>
<form action="/register_worker" method="post">
	<button>Register as a worker</button>
</form>
</html>`

const registerEmployer = `<html>
<form action="/register">
	<label for="username">Username: </label>
	<input type="text" id="username" name="username" value="" pattern="[a-zA-Z0-9]+" required><br>
	<label for="password">Password: </label>
	<input type="password" id="password" name="password" value="" pattern="[a-zA-Z0-9]+" required><br>
	<label for="cash">Available funds: </label>
	<input type="text" id="cash" name="cash" value="" pattern="[0-9]+" required><br>
	<input type="hidden" id="role" name="role" value="employer">
	<input type="submit" value="Register">
</form>
</html>`

const registerWorker = `<html>
<form action="/register">
	<label for="username">Username: </label>
	<input type="text" id="username" name="username" value="" pattern="[a-zA-Z0-9]+" required><br>
	<label for="password">Password: </label>
	<input type="password" id="password" name="password" value="" pattern="[a-zA-Z0-9]+" required><br>
	<label for="wage">Desired wage: </label>
	<input type="text" id="wage" name="wage" value="" pattern="[0-9]+" required><br>
	<label for="weight">Max weight: </label>
	<input type="text" id="weight" name="weight" value="" pattern="[0-9]+" required><br>
	Alcohol? <input type="radio" id="true" name="drinks" value="true" required>
	<label for="true">HELL YEAH</label>
	<input type="radio" id="false" name="drinks" value="false" required>
	<label for="true">No.</label><br>
	<input type="hidden" id="role" name="role" value="worker">
	<input type="submit" value="Register">
</form>
</html>`

const homeWorker = `<html>
<h2>Logged in as '%s' (worker)</h2>
<form action="/me" method="post">
<input type="hidden" id="token" name="token" value="%s">
<input type="submit" value="My profile">
</form>
<form action="/tasks" method="post">
<input type="hidden" id="token" name="token" value="%s">
<input type="submit" value="My tasks">
</form>
</html>`

const homeEmployer = `<html>
<h2>Logged in as '%s' (employer)</h2>
<input type="hidden" id="token" name="token" value="%s">
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

const profileWorker = `<html>
<h2>Logged in as '%s' (worker)</h2>
<input type="hidden" id="token" name="token" value="%s">
<form action="/home" method="post">
	<button>Home</button>
</form>
<form action="/tasks" method="post">
	<button>My tasks</button>
</form>
<table>
	<tr>
		<th>ID</th>
		<th>Wage</th>
		<th>Fatigue</th>
		<th>Max weight</th>
		<th>Drinking</th>
	</tr>
	<tr>
		<td>%d</td>
		<td>%d</td>
		<td>%d</td>
		<td>%d</td>
		<td>%v</td>
	</tr>
</table>
</html>`

const profileEmployer = `<html>
<h2>Logged in as '%s' (employer)</h2>
<input type="hidden" id="token" name="token" value="%s">
<form action="/home" method="post">
	<button>Home</button>
</form>
<form action="/tasks" method="post">
	<button>My tasks</button>
</form>
<form action="/start" method="post">
	<button>Play the game</button>
</form>
<h2>Cash - %d</h2>%s</html>`

const tasksWorker = `<html>
<h2>Logged in as '%s' (tasks)</h2>
<input type="hidden" id="token" name="token" value="%s">
<form action="/home" method="post">
	<button>Home</button>
</form>
<form action="/me" method="post">
	<button>My profile</button>
</form>%s</html>`

const tasksEmployer = `<html>
<h2>Logged in as '%s' (employer)</h2>
<input type="hidden" id="token" name="token" value="%s">
<form action="/home" method="post">
	<button>Home</button>
</form>
<form action="/me" method="post">
	<button>My profile</button>
</form>%s</html>`

const unauthorizedAccess = `<html>Unauthorized access</html>`
