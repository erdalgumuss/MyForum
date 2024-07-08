document.addEventListener('DOMContentLoaded', function () {
    const loginPopup = document.getElementById('login-popup');
    const registerPopup = document.getElementById('register-popup');
    const loginBtn = document.getElementById('login-btn');
    const registerBtn = document.getElementById('register-btn');
    const closeLogin = document.getElementById('close-login');
    const closeRegister = document.getElementById('close-register');

    loginBtn.addEventListener('click', function (event) {
        event.preventDefault();
        loginPopup.style.display = 'block';
    });

    registerBtn.addEventListener('click', function (event) {
        event.preventDefault();
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
});

document.addEventListener('DOMContentLoaded', function () {
    const loginForm = document.getElementById('login-form');
    loginForm.addEventListener('submit', function (event) {
        event.preventDefault();
        const formData = new FormData(loginForm);

        // FormData'yı JSON formatına çevirme
        const data = {};
        formData.forEach((value, key) => {
            data[key] = value;
        });

        fetch('/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                console.error('Login error:', data.error);
            } else {
                console.log('Login successful:', data.message);
                window.location.href = "/"; // Redirect to homepage
            }
        })
        .catch(error => console.error('Login request failed:', error));
    });
});


/*
function loadProfileInfo() {
    fetch('/api/profile')
        .then(response => response.json())
        .then(data => {
            document.getElementById('profile-picture').innerHTML = `<img src="${data.profilePicture}" alt="Profil Resmi">`;
            document.getElementById('profile-details').innerHTML = `
                <h3>${data.username}</h3>
                <p>${data.email}</p>
                <!-- Diğer profil bilgileri buraya eklenebilir -->
            `;
        })
        .catch(error => console.error('Profil bilgileri yüklenirken hata oluştu:', error));
}

// Sayfa yüklendiğinde profil bilgilerini yükle
document.addEventListener('DOMContentLoaded', function() {
    loadProfileInfo();
});
function loadUserLikes() {
    fetch('/api/likes')
        .then(response => response.json())
        .then(data => {
            const likesList = document.getElementById('user-likes');
            likesList.innerHTML = ''; // Önceki içeriği temizle
            data.likes.forEach(like => {
                const li = document.createElement('li');
                li.textContent = like.title;
                likesList.appendChild(li);
            });
        })
        .catch(error => console.error('Beğeniler yüklenirken hata oluştu:', error));
}

// Sayfa yüklendiğinde beğenileri yükle
document.addEventListener('DOMContentLoaded', function() {
    loadUserLikes();
});
function loadUserComments() {
    fetch('/api/comments')
        .then(response => response.json())
        .then(data => {
            const commentsList = document.getElementById('user-comments');
            commentsList.innerHTML = ''; // Önceki içeriği temizle
            data.comments.forEach(comment => {
                const li = document.createElement('li');
                li.textContent = comment.content;
                commentsList.appendChild(li);
            });
        })
        .catch(error => console.error('Yorumlar yüklenirken hata oluştu:', error));
}

// Sayfa yüklendiğinde yorumları yükle
document.addEventListener('DOMContentLoaded', function() {
    loadUserComments();
});
function loadUserPosts() {
    fetch('/api/posts')
        .then(response => response.json())
        .then(data => {
            const postsContainer = document.getElementById('user-posts');
            postsContainer.innerHTML = ''; // Önceki içeriği temizle
            data.posts.forEach(post => {
                const div = document.createElement('div');
                div.classList.add('post');
                div.innerHTML = `
                    <h4>${post.title}</h4>
                    <p>${post.content}</p>
                `;
                postsContainer.appendChild(div);
            });
        })
        .catch(error => console.error('Gönderiler yüklenirken hata oluştu:', error));
}

// Sayfa yüklendiğinde gönderileri yükle
document.addEventListener('DOMContentLoaded', function() {
    loadUserPosts();
});
*/