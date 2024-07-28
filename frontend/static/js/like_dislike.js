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
