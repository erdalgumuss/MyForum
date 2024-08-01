alert("forum.js loaded");

// Function to fetch and display threads
function fetchThreads(category = '') {
    let url = '/getpost';
    if (category) {
        url += `?category=${encodeURIComponent(category)}`;
    }

    fetch(url, {
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

const handleFormSubmit = (form, url) => {
    if (form) {
        form.addEventListener('submit', async (event) => {
            event.preventDefault();
            const formData = new FormData(form);

            try {
                const response = await fetch(url, {
                    method: 'POST',
                    body: formData // Use formData directly for multipart/form-data
                });

                if (response.ok) {
                    const responseData = await response.json();
                    if (form.id === 'create-post-form') {
                        // Redirect to the new post's URL
                        window.location.href = `/posts/${responseData.postID}`;

                    } else if (form.id === 'create-comment-form') {
                        // Redirect to the post's URL with the new comment
                        window.location.href = `/posts/${responseData.postID}`;
                    }
                } else {
                    const responseData = await response.json();
                    if (responseData.error === "File size exceeds 20 MB") {
                        alert("File size exceeds 20 MB");
                    } else if (response.status === 401) {
                        alert("Please log in to create a post");
                    } else {
                        alert(responseData.error || "An error occurred. Please try again.");
                    }
                }
            } catch (error) {
                console.error(`${form.id} request failed:`, error);
            }
        });
    }
};

document.addEventListener('DOMContentLoaded', () => {
    handleFormSubmit(document.getElementById('create-post-form'), '/create-post');
    handleFormSubmit(document.getElementById('create-comment-form'), '/create-comment');
});

// CHECK IF THE USER IS THE SAME USER WHO WANTS TO EDIT THE POST
document.addEventListener('DOMContentLoaded', () => {
    const editPostLinkElement = document.getElementById('edit-post-link');
    const postID = document.getElementById('post-container').getAttribute('data-post-id');

    const fetchCurrentUser = async () => {
        try {
            const response = await fetch('/models/user');
            if (!response.ok) {
                throw new Error('Failed to fetch user data');
            }
            const user = await response.json();
            console.log("Current user fetched:", user); // Log current user data
            return user;
        } catch (error) {
            console.error('Error loading user:', error);
            return null;
        }
    };

    const checkIfUserCanEdit = async () => {
        const currentUser = await fetchCurrentUser();
        if (!currentUser) {
            return;
        }

        const postUsername = document.getElementById('post-author-username').textContent.trim();
        console.log("Post username:", postUsername); // Log post username
        console.log("Current username:", currentUser.username); // Log current user username

        if (currentUser.username === postUsername) {
            editPostLinkElement.innerHTML = `<a href="/post/edit/${postID}">Edit Post</a>`;
        }
    };

    checkIfUserCanEdit();
});
