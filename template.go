package authgo

import (
	"net/http"
)

func writeString(w http.ResponseWriter, v string) error {
	_, err := w.Write([]byte(v))

	if err != nil {
		return err
	}

	return nil
}

const (
	templateLogin = `<!doctype html>
<html lang="en">

<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
	<title>authgo</title>
	<link rel="stylesheet"
		href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-beta.3/css/bootstrap.min.css"
		integrity="sha384-Zug+QiDoJOrZ5t4lssLdxGhVrurbmBWopoEl+M6BdEfwnCJZtKxi1KgxUyJq13dy"
		crossorigin="anonymous"
	/>
</head>

<body class="text-white bg-dark">
	<main class="container-fluid">
		<section class="row justify-content-center">
			<form class="col-auto mt-5" action="/login" method="post">
				<fieldset>
					<legend>authgo</legend>
					<div class="form-group">
						<label for="email" class="sr-only">Email:</label>
						<input id="email" class="form-control" type="text" name="email" placeholder="Email">
					</div>
					<div class="form-group">
						<label for="password" class="sr-only">Password:</label>
						<input id="password" class="form-control" type="password" name="password" placeholder="Password">
					</div>
					<button type="submit" class="btn btn-primary">Submit</button>
				</fieldset>
			</form>
		</section>
	</main>
	<script src="https://code.jquery.com/jquery-3.2.1.slim.min.js"
		integrity="sha384-KJ3o2DKtIkvYIK3UENzmM7KCkRr/rE9/Qpg6aAZGJwFDMVNA/GpGFF93hXpG5KkN"
		crossorigin="anonymous">
	</script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.12.9/umd/popper.min.js"
		integrity="sha384-ApNbgh9B+Y1QKtv3Rn7W3mgPxhU9K/ScQsAP7hUibX39j7fakFPskvXusvfa0b4Q"
		crossorigin="anonymous">
	</script>
	<script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-beta.3/js/bootstrap.min.js"
		integrity="sha384-a5N7Y/aK3qNeh15eJKGWxsqtnX/wWdSZSKp+81YjTmS15nvnvxKHuzaWwXHDli+4"
		crossorigin="anonymous">
	</script>
</body>

</html>`
	templateGraphiQL = `<!doctype html>
<html lang="en">

<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
	<title>authgo</title>
	<link rel="stylesheet"
		href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.css"
		crossorigin="anonymous"
	/>
</head>

<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
	<main id="graphiql" style="height: 100vh;">
		Loading...
	</main>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/2.0.3/fetch.min.js"
		crossorigin="anonymous">
	</script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/react/16.2.0/umd/react.production.min.js" 
		crossorigin="anonymous">
	</script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/react-dom/16.2.0/umd/react-dom.production.min.js"
		crossorigin="anonymous">
	</script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.js"
		crossorigin="anonymous">
	</script>
	<script>
		function graphQLFetcher(graphQLParams) {
			return fetch('/graphql', {
				method: "post",
				body: JSON.stringify(graphQLParams),
				credentials: 'include',
			}).then(function (response) {
				return response.text();
			}).then(function (responseBody) {
				try {
					return JSON.parse(responseBody);
				} catch (error) {
					return responseBody;
				}
			});
		}
		ReactDOM.render(
			React.createElement(GraphiQL, { fetcher: graphQLFetcher }),
			document.getElementById('graphiql')
		);
	</script>
</body>

</html>`
)
