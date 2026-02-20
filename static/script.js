document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('prevWeekButton').addEventListener('click', function() {
        updateWeek('previous');
    });
    document.getElementById('nextWeekButton').addEventListener('click', function() {
        updateWeek('next');
    });
    document.getElementById('week').addEventListener('click', function() {
        resetToCurrentWeek();
    });

    document.getElementById('prevWeekButton').addEventListener('mousedown', function(e) {
        e.preventDefault();
    });
    document.getElementById('nextWeekButton').addEventListener('mousedown', function(e) {
        e.preventDefault();
    });
    document.getElementById('week').addEventListener('mousedown', function(e) {
        e.preventDefault();
    });

    document.addEventListener('keydown', function(event) {
        if (event.key === 'ArrowLeft') {
            event.preventDefault();
            updateWeek('previous');
        } else if (event.key === 'ArrowRight') {
            event.preventDefault();
            updateWeek('next');
        }
    });
});

function updateWeek(direction) {
    const currentWeek = parseInt(document.getElementById('week').innerText);
    const firstDateText = document.getElementById('firstDate').textContent;
    const firstDateYear = firstDateText.split('-')[0];
    const lastDateText = document.getElementById('lastDate').textContent;
    const lastDateYear = lastDateText.split('-')[0];

    const currentYear = parseInt(currentWeek) > 50 && (parseInt(firstDateYear) < parseInt(lastDateYear)) ? firstDateYear : lastDateYear;
    console.log(currentYear);
    let url = '';
    if (direction === 'previous') {
        url = `/api/previous/year/${currentYear}/week/${currentWeek}`;
    } else if (direction === 'next') {
        url = `/api/next/year/${currentYear}/week/${currentWeek}`;
    }

    fetch(url)
        .then(response => response.json())
        .then(data => updateWeekInfo(data))
        .catch(error => console.error('Error:', error));
}

function resetToCurrentWeek() {
    fetch(`/api/week/current`)
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