function updateWeek(weekChange) {
    const currentWeek = parseInt(document.getElementById('week').innerText);
    const newWeek = currentWeek + weekChange;

    fetch(`/week/${newWeek}`)
        .then(response => response.json())
        .then(data => {
            document.getElementById('week').innerText = data.Week;
            document.getElementById('dateRange').innerText = `${data.FirstDate} - ${data.LastDate}`;
            document.getElementById('pageTitle').innerText = `Current week is ${data.Week} | ${data.FirstDate} - ${data.LastDate}`;
        })
        .catch(error => console.error('Error:', error));
}