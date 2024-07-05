// Simüle edilmiş beğeni sayısı
let likeCount = 0;
let dislikeCount = 0;

// Beğeni işlevi
function likePost() {
    likeCount++;
    updateLikes();
}

// Beğenmeme işlevi
function dislikePost() {
    dislikeCount++;
    updateLikes();
}

// Beğeni sayılarını güncelleme işlevi
function updateLikes() {
    document.getElementById('like-count').textContent = likeCount;
    document.getElementById('dislike-count').textContent = dislikeCount;
}

// DOM yüklendiğinde çağrılan işlevler
document.addEventListener('DOMContentLoaded', function() {
    // Beğenme butonu
    document.getElementById('like-btn').addEventListener('click', likePost);

    // Beğenmeme butonu
    document.getElementById('dislike-btn').addEventListener('click', dislikePost);
});
