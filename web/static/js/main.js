    document.addEventListener("DOMContentLoaded", function () {
    // Регистрация
    const registerForm = document.getElementById("registerForm");
    if (registerForm) {
    registerForm.addEventListener("submit", async function(e) {
    e.preventDefault();
    const formData = new FormData(registerForm);
    try {
    const response = await fetch("/register", {
    method: "POST",
    body: formData
});

    if (response.redirected) {
    window.location.href = response.url;
} else {
    const html = await response.text();
    document.open();
    document.write(html);
    document.close();
}
} catch (err) {
    console.error(err);
}
});
}

    // Авторизация
    const loginForm = document.getElementById("loginForm");
    if (loginForm) {
    loginForm.addEventListener("submit", async function(e) {
    e.preventDefault();
    const formData = new FormData(loginForm);
    try {
    const response = await fetch("/login", {
    method: "POST",
    body: formData
});

    if (response.redirected) {
    window.location.href = response.url;
} else {
    const html = await response.text();
    document.open();
    document.write(html);
    document.close();
}
} catch (err) {
    console.error(err);
}
});
}

    // Загрузка трека
    const uploadForm = document.getElementById("uploadForm");
    if (uploadForm) {
    uploadForm.addEventListener("submit", async function(e) {
    e.preventDefault();
    const formData = new FormData(uploadForm);
    try {
    const response = await fetch("/tracks", {
    method: "POST",
    body: formData
});

    if (response.redirected) {
    window.location.href = response.url;
} else {
    const result = await response.json();
    document.getElementById("uploadResult").innerText = result.error || "Ошибка загрузки";
}
} catch (err) {
    console.error(err);
}
});
}
});
