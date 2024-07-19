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

function submitForm() {
    const form = document.getElementById('createPostForm');
    if (!form.checkValidity()) {
        form.reportValidity();
        return;
    }
    const formData = new FormData(form);

    // Convert selected categories into an array of strings
    const selectedCategories = Array.from(formData.getAll('categories'));

    // Add each selected category individually to formData
    formData.delete('categories');
    selectedCategories.forEach(category => {
        formData.append('categories', category);
    });

    fetch('/create-post', {
        method: 'POST',
        body: formData
    })
    .then(response => response.json())
    .then(data => {
        console.log('Success:', data);
        if (data.message === "Post başarıyla oluşturuldu") {
            if (window.location.pathname === '/forum') {
                fetchThreads();
            }
            form.reset(); // Optionally reset the form after successful submission
        }
    })
    .catch((error) => {
        console.error('Error:', error);
    });
}

// COMMENT SECTION STARTED //

function submitComment(event) {
    event.preventDefault();

    const form = document.getElementById('commentForm');
    const formData = new FormData(form);
    const commentData = Object.fromEntries(formData);

    fetch('/comments', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            content: commentData.content,
            post_id: parseInt(commentData.post_id, 10),
            user_id: parseInt(commentData.user_id, 10)
        }),
    })
    .then(response => response.json())
    .then(data => {
        if (data.message === "Comment created successfully") {
            alert("Comment submitted successfully!");
            form.reset();
            loadComments(commentData.post_id);  // Refresh comments section
        } else {
            alert("Failed to submit comment: " + data.error);
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert('Error submitting comment. Please try again later.');
    });
}

function loadComments(postId) {
    fetch(`/posts/${postId}/comments`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => response.json())
    .then(data => {
        const commentsDiv = document.getElementById('comments');
        commentsDiv.innerHTML = ''; // Clear previous comments

        if (!Array.isArray(data) || data.length === 0) {
            commentsDiv.textContent = 'No comments yet.';
            return;
        }

        data.forEach(comment => {
            const commentContainer = document.createElement('div');
            commentContainer.classList.add('comment');

            const contentElement = document.createElement('p');
            contentElement.textContent = comment.content;
            commentContainer.appendChild(contentElement);

            const authorElement = document.createElement('p');
            authorElement.textContent = `Posted by User ID: ${comment.user_id} on ${new Date(comment.created_at).toLocaleString()}`;
            commentContainer.appendChild(authorElement);

            commentsDiv.appendChild(commentContainer);
            const hrElement = document.createElement('hr');
            commentsDiv.appendChild(hrElement);
        });
    })
    .catch(error => {
        console.error('Error fetching comments:', error);
    });
}

document.addEventListener('DOMContentLoaded', () => {
    const commentForm = document.getElementById('commentForm');
    if (commentForm) {
        commentForm.addEventListener('submit', submitComment);
    }

    // Load comments for the current post
    const postId = document.getElementById('post_id').value;
    loadComments(postId);
});


// COMMENT SECTION ENDED //

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
