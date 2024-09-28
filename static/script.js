document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('prevWeekButton').addEventListener('click', function() {
        updateWeek(-1);
    });
    document.getElementById('nextWeekButton').addEventListener('click', function() {
        updateWeek(1);
    });
    document.getElementById('week').addEventListener('click', function() {
        resetToCurrentWeek();
    });
});

function updateWeek(weekChange) {
    const currentWeek = parseInt(document.getElementById('week').innerText);
    const newWeek = currentWeek + weekChange;

    fetch(`/week/${newWeek}`)
        .then(response => response.json())
        .then(data => updateWeekInfo(data))
        .catch(error => console.error('Error:', error));
}

function resetToCurrentWeek() {
    fetch(`/week/current`)
        .then(response => response.json())
        .then(data => updateWeekInfo(data))
        .catch(error => console.error('Error:', error));
}

function updateWeekInfo(data) {
    document.getElementById('week').innerText = data.Week;
    document.getElementById('firstDate').innerText = data.FirstDate;
    document.getElementById('lastDate').innerText = data.LastDate;
    document.getElementById('pageTitle').innerText = `Current week is ${data.Week} | ${data.FirstDate} - ${data.LastDate}`;
}