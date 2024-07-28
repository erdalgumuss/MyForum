document.addEventListener('DOMContentLoaded', () => {
    console.log("DOM fully loaded and parsed - profile.js");
    const userNameElement = document.getElementById('user-name');
    const userEmailElement = document.getElementById('user-email');
    const postsContainer = document.getElementById('posts-container');
    const topicsLikesList = document.getElementById('topics-likes-list');
    const commentsLikesList = document.getElementById('comments-likes-list');
    const commentsContainer = document.getElementById('comments-container');
    const requestModeratorBtn = document.getElementById('request-moderator-btn');

    const loadUser = async () => {
        console.log("Loading user profile - profile.js");
        try {
            const response = await fetch('/models/user');
            const user = await response.json();
            if (response.ok) {
                console.log("User profile loaded - profile.js", user);

                userNameElement.textContent = `${user.name} ${user.surname}`;
                userEmailElement.textContent = user.email;

                loadUserPosts(user.id);
                loadUserLikes(user.id);
                loadUserComments(user.id);
            } else {
                console.error('Failed to load user profile - profile.js');
            }
        } catch (error) {
            console.error('Error loading user - profile.js:', error);
        }
    };

    const loadUserPosts = async (userId) => {
        console.log("Loading user posts - profile.js");
        try {
            const response = await fetch(`/user/${userId}/posts`);
            const posts = await response.json();
            postsContainer.innerHTML = '';
            posts.forEach(post => {
                const postElement = document.createElement('div');
                postElement.classList.add('post');
                postElement.innerHTML = `<h4><a href="/posts/${post.id}">${post.title}</a></h4><p>${post.content}</p>`;
                postsContainer.appendChild(postElement);
            });
        } catch (error) {
            console.error('Error loading user posts - profile.js:', error);
        }
    };

    const loadUserLikes = async (userId) => {
        console.log("Loading user likes - profile.js");
        try {
            const response = await fetch(`/user/${userId}/likes`);
            const likes = await response.json();
            topicsLikesList.innerHTML = '';
            commentsLikesList.innerHTML = '';
            likes.forEach(like => {
                const likeElement = document.createElement('li');
                if (like.post_id && like.post_id.Valid) {
                    likeElement.innerHTML = `<a href="/posts/${like.post_id.Int64}">${like.post_title}</a>`;
                    topicsLikesList.appendChild(likeElement);
                } else if (like.comment_id && like.comment_id.Valid) {
                    likeElement.innerHTML = `<a href="/posts/${like.comment_id.Int64}">${like.post_title}</a>`;
                    commentsLikesList.appendChild(likeElement);
                }
            });
        } catch (error) {
            console.error('Error loading user likes - profile.js:', error);
        }
    };

    const loadUserComments = async (userId) => {
        console.log("Loading user comments - profile.js");
        try {
            const response = await fetch(`/user/${userId}/comments`);
            const comments = await response.json();
            commentsContainer.innerHTML = '';
            comments.forEach(comment => {
                const commentElement = document.createElement('div');
                commentElement.classList.add('comment');
                commentElement.innerHTML = `<h4><a href="/posts/${comment.post_id}">${comment.PostTitle}</a></h4><p>${comment.content}</p>`;
                commentsContainer.appendChild(commentElement);
            });
        } catch (error) {
            console.error('Error loading user comments - profile.js:', error);
        }
    };

    const requestModerator = async () => {
        console.log("Requesting moderator - profile.js");
        try {
            const response = await fetch('/user/request_moderator', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                }
            });

            if (response.ok) {
                alert('Moderator request submitted successfully.');
                requestModeratorBtn.style.display = 'none';
            } else {
                alert('Failed to submit moderator request.');
            }
        } catch (error) {
            console.error('Error submitting moderator request - profile.js:', error);
        }
    };

    requestModeratorBtn.addEventListener('click', requestModerator);

    loadUser();
});
