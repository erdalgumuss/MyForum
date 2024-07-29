document.addEventListener('DOMContentLoaded', () => {
    const logoutBtn = document.getElementById('logout-btn');
    const userNameElement = document.getElementById('user-name');
    const userEmailElement = document.getElementById('user-email');
    const profilePictureElement = document.querySelector('#profile-picture img');
    const postsContainer = document.getElementById('posts-container');
    const likesList = document.getElementById('topics-likes-list');
    const commentsList = document.getElementById('comments-likes-list');
    const requestModeratorBtn = document.getElementById('request-moderator-btn');
    const moderatorRequestDiv = document.getElementById('moderator-request');
    const userRoleElement = document.getElementById('user-role');

    const toggleUserUI = (isLoggedIn) => {
        logoutBtn.style.display = isLoggedIn ? 'inline' : 'none';
    };

    const loadUser = async () => {
        try {
            const response = await fetch('/models/user');
            const user = await response.json();
            console.log('User data:', user); // Debug log
            if (response.ok) {
                toggleUserUI(true);
                userNameElement.textContent = `${user.name} ${user.surname}`;
                userEmailElement.textContent = user.email;
                profilePictureElement.src = user.profilePicture || '/static/images/default-profile.png';
                userRoleElement.textContent = user.role; // Set user role in hidden div

                // Check if the user is a moderator and redirect if necessary
                if (user.role === 'moderator') {
                    requestModeratorBtn.textContent = "Go to Moderator Dashboard";
                    requestModeratorBtn.onclick = () => {
                        window.location.href = '/moderator/dashboard';
                    };
                } else {
                    requestModeratorBtn.textContent = "Request Moderator";
                    requestModeratorBtn.onclick = requestModerator;
                }

                // Load user's posts, likes, and comments
                loadUserPosts(user.id);
                loadUserLikes(user.id);
                loadUserComments(user.id);
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
            postsContainer.innerHTML = ''; // Clear existing content
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
            likesList.innerHTML = ''; // Clear existing content
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
            commentsList.innerHTML = ''; // Clear existing content
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

    const requestModerator = async () => {
        try {
            const response = await fetch('/user/request_moderator', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                }
            });

            if (response.ok) {
                alert('Moderator request submitted successfully.');
                document.getElementById('moderator-request').style.display = 'none';
            } else {
                alert('Failed to submit moderator request.');
            }
        } catch (error) {
            console.error('Error submitting moderator request:', error);
        }
    };

    loadUser();

    logoutBtn.addEventListener('click', async (event) => {
        event.preventDefault();
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
                window.location.href = "/";  // Logout işleminden sonra anasayfaya yönlendir
            } else {
                const responseData = await response.json();
                alert(responseData.error);
            }
        } catch (error) {
            console.error('Logout request failed:', error);
        }
    });
});
