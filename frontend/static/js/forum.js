alert("forum.js loaded");

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
                    alert(responseData.error);
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
