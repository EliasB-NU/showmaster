const ws = new WebSocket('ws://localhost:8080/ws');

ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    displayRows(data);
};

function displayRows(rows) {
    const container = document.getElementById('rows-container');
    container.innerHTML = ''; // Clear existing content

    rows.forEach(row => {
        const rowElement = document.createElement('div');
        rowElement.textContent = `${row.id}: ${row.name} | ${row.audio} | ${row.licht} | ${row.pptx} | ${row.notes}`;
        container.appendChild(rowElement);
    });
}
