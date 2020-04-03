var params = (new URL(window.location.href)).searchParams;

if (params.get('unauthorized') != null) {
    document.querySelector('#error_unauthorized').style.display = 'block';
} else if (params.get('logout') != null) {
    document.querySelector('#logged_out').style.display = 'block';
}

var session;

document.querySelector('#login').addEventListener("submit", function(e) {
    e.preventDefault();

    var loginButton = document.querySelector('#login_button'),
        usernameInput = document.querySelector('#username'),
        passwordInput = document.querySelector('#password'),
        passwordError = document.querySelector('#error_password');

    loginButton.disabled = true;
    loginButton.textContent = 'Loading...';

    var payload = {
        Username: usernameInput.value,
        Password: passwordInput.value
    };

    fetch("/api/login", {
        method: 'POST',
        credentials: 'include',
        body: JSON.stringify(payload)
    }).then(function (res) {
        return res.json();
    }).then(function (data) {
        if (data.code == 200) {
            passwordError.style.display = 'none';
            window.location.href = '/';
        } else {
            usernameInput.value = '';
            passwordInput.value = '';
            usernameInput.focus();
            passwordError.style.display = 'block';
            loginButton.disabled = false;
            loginButton.textContent = 'Login';
        }
    });
    return false;
});
