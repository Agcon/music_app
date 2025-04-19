const audio = document.getElementById('audio');
const playPause = document.getElementById('playPause');
const progressBar = document.getElementById('progressBar');
const currentTime = document.getElementById('currentTime');
const duration = document.getElementById('duration');
const volume = document.getElementById('volume');

// Обновить play/pause
playPause.addEventListener('click', () => {
    if (audio.paused) {
        audio.play();
        playPause.textContent = '⏸';
    } else {
        audio.pause();
        playPause.textContent = '▶️';
    }
});

// Обновлять прогресс и таймер
audio.addEventListener('timeupdate', () => {
    progressBar.value = Math.floor(audio.currentTime);
    currentTime.textContent = formatTime(audio.currentTime);
});

// Установить длительность при загрузке
audio.addEventListener('loadedmetadata', () => {
    progressBar.max = Math.floor(audio.duration);
    duration.textContent = formatTime(audio.duration);
});

// Перемотка
progressBar.addEventListener('input', () => {
    audio.currentTime = progressBar.value;
});

// Громкость
volume.addEventListener('input', () => {
    audio.volume = volume.value;
});

// Форматирование времени
function formatTime(time) {
    const minutes = Math.floor(time / 60);
    const seconds = Math.floor(time % 60).toString().padStart(2, '0');
    return `${minutes}:${seconds}`;
}
