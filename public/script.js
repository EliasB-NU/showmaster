// script.js

const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = function(event) {
    console.log('WebSocket connection established.');
};

ws.onmessage = function(event) {
    console.log(event.data);
    const rows = JSON.parse(event.data);
    displayRows(rows);
};

function displayRows(rows) {
    const container = document.getElementById('rows-container');
    container.innerHTML = ''; // Clear existing content

    const highlightedRowID = localStorage.getItem('highlightedRowID');
    

    rows.forEach(row => {
        const rowElement = document.createElement('div');
        rowElement.textContent = `${row.id}: ${row.name}`;
        rowElement.classList.add('row'); // Add CSS class to row
        rowElement.dataset.id = row.id; // Set data-id attribute
        container.appendChild(rowElement);

        if (row.id === highlightedRowID) {
            rowElement.classList.add('highlighted');
        } else {
            rowElement.classList.remove('highlighted');
        }

        // Add click event listener to each row
        rowElement.addEventListener('click', () => {
            selectRow(row.id);
        });
    });
}

function selectRow(rowID) {
    ws.send(rowID); // Send message to server to broadcast the new highlighted row
    console.log(highlightedRowID);
}