document.addEventListener('DOMContentLoaded', function () {
    // Giriş ve Kayıt Pop-up'ları ile ilgili kodlar
    const loginPopup = document.getElementById('login-popup');
    const registerPopup = document.getElementById('register-popup');
    const loginBtn = document.getElementById('login-btn');
    const registerBtn = document.getElementById('register-btn');
    const closeLogin = document.getElementById('close-login');
    const closeRegister = document.getElementById('close-register');

    loginBtn.addEventListener('click', function () {
        loginPopup.style.display = 'block';
    });

    registerBtn.addEventListener('click', function () {
        registerPopup.style.display = 'block';
    });

    closeLogin.addEventListener('click', function () {
        loginPopup.style.display = 'none';
    });

    closeRegister.addEventListener('click', function () {
        registerPopup.style.display = 'none';
    });

    window.addEventListener('click', function (event) {
        if (event.target === loginPopup) {
            loginPopup.style.display = 'none';
        }
        if (event.target === registerPopup) {
            registerPopup.style.display = 'none';
        }
    });

    // Arama işlevi ile ilgili kodlar
    const searchBtn = document.getElementById('search-btn');
    const searchInput = document.getElementById('search');

    searchBtn.addEventListener('click', function () {
        const query = searchInput.value.toLowerCase();
        const posts = document.querySelectorAll('#posts .post');

        posts.forEach(function (post) {
            const title = post.querySelector('.post-title').textContent.toLowerCase();
            const content = post.querySelector('p').textContent.toLowerCase();
            if (title.includes(query) || content.includes(query)) {
                post.style.display = 'block';
            } else {
                post.style.display = 'none';
            }
        });
    });
});
