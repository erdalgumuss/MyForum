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
