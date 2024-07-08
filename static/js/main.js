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

    const loginForm = document.getElementById('login-form');
    loginForm.addEventListener('submit', function (event) {
        event.preventDefault();
        const formData = new FormData(this);
    
        const data = {};
        formData.forEach((value, key) => {
            data[key] = value;
        });
    
        console.log('Login form data:', data); // Form verilerini konsola yazdır
    
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
                window.location.href = "/profile"; // Redirect to profile page
            }
        })
        .catch(error => console.error('Login request failed:', error));
    });

    const registerForm = document.getElementById('register-form');
    registerForm.addEventListener('submit', function (event) {
        event.preventDefault();
        const formData = new FormData(this);

        // FormData'yı JSON formatına çevirme
        const data = {};
        formData.forEach((value, key) => {
            data[key] = value;
        });

        console.log('Register form data:', data); // Bu satırı ekleyin

        fetch('/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                console.error('Register error:', data.error);
            } else {
                console.log('Register successful:', data.message);
                window.location.href = "/"; // Anasayfaya veya giriş sayfasına yönlendirme
            }
        })
        .catch(error => console.error('Register request failed:', error));
    });

    function loadProfileInfo() {
        fetch('/models/profile')
            .then(response => response.json())
            .then(data => {
                document.getElementById('profile-picture').innerHTML = `<img src="${data.profilePicture}" alt="Profil Resmi">`;
                document.getElementById('profile-details').innerHTML = `
                    <h3>${data.username}</h3>
                    <p>${data.email}</p>
                `;
            })
            .catch(error => console.error('Profil bilgileri yüklenirken hata oluştu:', error));
    }

    function loadUserLikes() {
        fetch('/models/topic')
            .then(response => response.json())
            .then(data => {
                const likesList = document.getElementById('likes-list');
                likesList.innerHTML = '';
                data.likes.forEach(like => {
                    const li = document.createElement('li');
                    li.textContent = like.title;
                    likesList.appendChild(li);
                });
            })
            .catch(error => console.error('Beğeniler yüklenirken hata oluştu:', error));
    }

    function loadUserComments() {
        fetch('/models/comment')
            .then(response => response.json())
            .then(data => {
                const commentsList = document.getElementById('comments-list');
                commentsList.innerHTML = '';
                data.comments.forEach(comment => {
                    const li = document.createElement('li');
                    li.textContent = comment.content;
                    commentsList.appendChild(li);
                });
            })
            .catch(error => console.error('Yorumlar yüklenirken hata oluştu:', error));
    }

    function loadUserPosts() {
        fetch('/models/post')
            .then(response => response.json())
            .then(data => {
                const postsContainer = document.getElementById('posts-container');
                postsContainer.innerHTML = '';
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

    // Sayfa yüklendiğinde profil bilgilerini, beğenileri, yorumları ve gönderileri yükle
    loadProfileInfo();
    loadUserLikes();
    loadUserComments();
    loadUserPosts();
});
document.addEventListener('DOMContentLoaded', function() {
    const likeBtn = document.getElementById('like-btn');
    const dislikeBtn = document.getElementById('dislike-btn');
    const likeCount = document.getElementById('like-count');
    const dislikeCount = document.getElementById('dislike-count');

    likeBtn.addEventListener('click', function() {
        let count = parseInt(likeCount.innerText);
        likeCount.innerText = count + 1;
    });

    dislikeBtn.addEventListener('click', function() {
        let count = parseInt(dislikeCount.innerText);
        dislikeCount.innerText = count + 1;
    });
});
