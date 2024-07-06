// DOM yüklendiğinde çağrılan işlevler
document.addEventListener('DOMContentLoaded', function() {
    // Beğenme butonu
    document.getElementById('like-btn').addEventListener('click', likePost);

    // Beğenmeme butonu
    document.getElementById('dislike-btn').addEventListener('click', dislikePost);
});

// Beğeni işlevi
function likePost() {
    fetch('/like', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ action: 'like' })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            likeCount++;
            updateLikes();
        } else {
            alert("Beğeni işlemi başarısız oldu");
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert("Sunucu hatası oluştu");
    });
}

// Beğenmeme işlevi
function dislikePost() {
    fetch('/like', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ action: 'dislike' })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            dislikeCount++;
            updateLikes();
        } else {
            alert("Beğenmeme işlemi başarısız oldu");
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert("Sunucu hatası oluştu");
    });
}

// Beğeni sayılarını güncelleme işlevi
function updateLikes() {
    document.getElementById('like-count').textContent = likeCount;
    document.getElementById('dislike-count').textContent = dislikeCount;
}
