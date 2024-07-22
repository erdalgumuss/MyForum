function fetchThreads() {
    fetch('/getpost', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to fetch threads');
        }
        return response.json();
    })
    .then(data => {
        const threadsDiv = document.getElementById('threads');
        threadsDiv.innerHTML = '';
        data.forEach(thread => {
            const threadDiv = document.createElement('div');
            threadDiv.innerHTML = `<h2><a href="/posts/${thread.id}">${thread.title}</a></h2><p>${thread.content}</p>`;
            threadsDiv.appendChild(threadDiv);
        });
    })
    .catch(error => {
        console.error('Error fetching threads:', error);
        alert('Error fetching threads. Please try again later.'); // Display error to user
    });
}

document.addEventListener('DOMContentLoaded', () => {
    const loginPopup = document.getElementById('login-popup');
    const registerPopup = document.getElementById('register-popup');
    const loginBtn = document.getElementById('login-btn');
    const registerBtn = document.getElementById('register-btn');
    const closeLogin = document.getElementById('close-login');
    const closeRegister = document.getElementById('close-register');
    const logoutBtn = document.getElementById('logout-btn');
    const userInfoContainer = document.getElementById('user-info');
    const userNameElement = document.getElementById('user-name');
    const userEmailElement = document.getElementById('user-email');

    const togglePopup = (popup, action) => {
        popup.style.display = action === 'open' ? 'block' : 'none';
    };

    loginBtn.addEventListener('click', (event) => {
        event.preventDefault();
        togglePopup(loginPopup, 'open');
    });

    registerBtn.addEventListener('click', (event) => {
        event.preventDefault();
        togglePopup(registerPopup, 'open');
    });

    closeLogin.addEventListener('click', () => {
        togglePopup(loginPopup, 'close');
    });

    closeRegister.addEventListener('click', () => {
        togglePopup(registerPopup, 'close');
    });

    window.addEventListener('click', (event) => {
        if (event.target === loginPopup) {
            togglePopup(loginPopup, 'close');
        }
        if (event.target === registerPopup) {
            togglePopup(registerPopup, 'close');
        }
    });

    const handleFormSubmit = (form, url) => {
        form.addEventListener('submit', async (event) => {
            event.preventDefault();
            const formData = new FormData(form);
            const data = Object.fromEntries(formData);

            try {
                const response = await fetch(url, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(data)
                });

                const responseData = await response.json();
                if (response.ok) {
                    if (form.id === 'login-form') {
                        localStorage.setItem('user', JSON.stringify(responseData));
                        loadUser(); // Reload user information
                        togglePopup(loginPopup, 'close');
                    } else if (form.id === 'register-form') {
                        togglePopup(registerPopup, 'close');
                    }
                } else {
                    alert(responseData.error);
                }
            } catch (error) {
                console.error(`${form.id} request failed:`, error);
            }
        });
    };

    handleFormSubmit(document.getElementById('login-form'), '/login');
    handleFormSubmit(document.getElementById('register-form'), '/register');

    const toggleUserUI = (isLoggedIn) => {
        loginBtn.style.display = isLoggedIn ? 'none' : 'inline';
        registerBtn.style.display = isLoggedIn ? 'none' : 'inline';
        logoutBtn.style.display = isLoggedIn ? 'inline' : 'none';
        userInfoContainer.style.display = isLoggedIn ? 'inline' : 'none';

        // Show or hide the profile link
        const profileLink = document.getElementById('profile-link');
        if (profileLink) {
            profileLink.style.display = isLoggedIn ? 'inline' : 'none';
        }
    };

    const loadUser = async () => {
        try {
            const response = await fetch('/models/user');
            const user = await response.json();
            if (response.ok) {
                toggleUserUI(true);
                userNameElement.textContent = `${user.name} ${user.surname}`;
                userEmailElement.textContent = user.email;
                // If on profile.html, update profile information
                if (window.location.pathname === '/profile.html') {
                    document.getElementById('profile-name').textContent = `${user.name} ${user.surname}`;
                    document.getElementById('profile-email').textContent = user.email;

                    // Load user-specific content
                    loadUserPosts(user.id);
                    loadUserLikes(user.id);
                    loadUserComments(user.id);
                }
            } else {
                toggleUserUI(false);
            }
        } catch (error) {
            console.error('Error loading user:', error);
            toggleUserUI(false);
        }
    };

    const loadUserPosts = async (userId) => {
        try {
            const response = await fetch(`/user/${userId}/posts`);
            const posts = await response.json();
            const postsContainer = document.getElementById('posts-container');
            postsContainer.innerHTML = ''; // Mevcut içeriği temizle
            posts.forEach(post => {
                const postElement = document.createElement('div');
                postElement.classList.add('post');
                postElement.innerHTML = `
                    <h4>${post.title}</h4>
                    <p>${post.content}</p>
                `;
                postsContainer.appendChild(postElement);
            });
        } catch (error) {
            console.error('Error loading user posts:', error);
        }
    };

    const loadUserLikes = async (userId) => {
        try {
            const response = await fetch(`/user/${userId}/likes`);
            const likes = await response.json();
            const likesList = document.getElementById('likes-list');
            likesList.innerHTML = ''; // Mevcut içeriği temizle
            likes.forEach(like => {
                const likeElement = document.createElement('li');
                likeElement.textContent = like.postTitle;
                likesList.appendChild(likeElement);
            });
        } catch (error) {
            console.error('Error loading user likes:', error);
        }
    };

    const loadUserComments = async (userId) => {
        try {
            const response = await fetch(`/user/${userId}/comments`);
            const comments = await response.json();
            const commentsList = document.getElementById('comments-list');
            commentsList.innerHTML = ''; // Mevcut içeriği temizle
            comments.forEach(comment => {
                const commentElement = document.createElement('li');
                commentElement.innerHTML = `
                    <strong>${comment.postTitle}</strong>: ${comment.content}
                `;
                commentsList.appendChild(commentElement);
            });
        } catch (error) {
            console.error('Error loading user comments:', error);
        }
    };

    loadUser();

    // Fetch threads only if on the forum page
    if (window.location.pathname === '/forum') {
        fetchThreads();
    }

    logoutBtn.addEventListener('click', async () => {
        try {
            const response = await fetch('/logout', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                }
            });

            if (response.ok) {
                const responseData = await response.json();
                alert(responseData.message);
                localStorage.removeItem('user');
                toggleUserUI(false);
                window.location.href = "/";  // Redirect to homepage after logout
            } else {
                const responseData = await response.json();
                alert(responseData.error);
            }
        } catch (error) {
            console.error('Logout request failed:', error);
        }
    });
});


// LIKE DISLIKE 4 POSTS N THREADS STARTED 

function updateLikesDislikes(type, id, likes, dislikes) {
    const elementId = type === 'post' ? 'post-likes-dislikes' : `comment-likes-dislikes-${id}`;
    const element = document.getElementById(elementId);
    element.textContent = `Likes: ${likes} | Dislikes: ${dislikes}`;
}

async function likePost(postId) {
    try {
        const response = await fetch(`/posts/${postId}/like`, { method: 'POST' });
        if (response.ok) {
            const data = await response.json();
            updateLikesDislikes('post', postId, data.likes, data.dislikes);
        } else {
            console.error('Failed to like post:', await response.text());
        }
    } catch (error) {
        console.error('Error liking post:', error);
    }
}

async function dislikePost(postId) {
    try {
        const response = await fetch(`/posts/${postId}/dislike`, { method: 'POST' });
        if (response.ok) {
            const data = await response.json();
            updateLikesDislikes('post', postId, data.likes, data.dislikes);
        } else {
            console.error('Failed to dislike post:', await response.text());
        }
    } catch (error) {
        console.error('Error disliking post:', error);
    }
}

async function likeComment(commentId) {
    try {
        const response = await fetch(`/comments/${commentId}/like`, { method: 'POST' });
        if (response.ok) {
            const data = await response.json();
            updateLikesDislikes('comment', commentId, data.likes, data.dislikes);
        } else {
            console.error('Failed to like comment:', await response.text());
        }
    } catch (error) {
        console.error('Error liking comment:', error);
    }
}

async function dislikeComment(commentId) {
    try {
        const response = await fetch(`/comments/${commentId}/dislike`, { method: 'POST' });
        if (response.ok) {
            const data = await response.json();
            updateLikesDislikes('comment', commentId, data.likes, data.dislikes);
        } else {
            console.error('Failed to dislike comment:', await response.text());
        }
    } catch (error) {
        console.error('Error disliking comment:', error);
    }
}



// LIKE DISLIKE ENDED