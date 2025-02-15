document.addEventListener('DOMContentLoaded', () => {
    const modal = document.getElementById('player-modal');
    const videoPlayer = document.getElementById('video-player');
    const closeBtn = document.querySelector('.close');

    // Fetch movies
    fetch('/api/movies')
        .then(response => response.json())
        .then(movies => {
            const grid = document.getElementById('movies-grid');
            movies.forEach(movie => {
                const card = document.createElement('div');
                card.className = 'movie-card';
                card.innerHTML = `
                    <h3>${movie}</h3>
                    <div class="play-button">▶ Play</div>
                `;
                
                card.addEventListener('click', () => {
                    const encodedPath = encodeURIComponent(movie).replace(/%2F/g, '/');
                    videoPlayer.src = `/media/${encodedPath}`;
                    modal.style.display = 'block';
                    videoPlayer.play();
                });
                
                grid.appendChild(card);
            });
        });

    // Close modal
    closeBtn.onclick = () => {
        modal.style.display = 'none';
        videoPlayer.pause();
    };

    window.onclick = (event) => {
        if (event.target === modal) {
            modal.style.display = 'none';
            videoPlayer.pause();
        }
    };
});