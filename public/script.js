// script.js

const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = function(event) {
    console.log('WebSocket connection established.');
};

ws.onmessage = function(event) {
    const msgs = JSON.parse(event.data);
    displayRows(msgs);
};

function displayRows(msgs) {
    msgs.sort((a, b) => a.row.id - b.row.id);

    const container = document.getElementById('rows-container');
    container.innerHTML = ''; // Clear existing content    

    msgs.forEach(msg => {
        const rowElement = document.createElement('div');
        rowElement.textContent = `${msg.row.id}: ${msg.row.name} | ${msg.row.audio} | ${msg.row.licht} | ${msg.row.pptx} | ${msg.row.notes}`;
        rowElement.classList.add('row'); // Add CSS class to row
        rowElement.dataset.id = msg.row.id; // Set data-id attribute
        container.appendChild(rowElement);
        if (msg.highlighted === true) {
            rowElement.classList.add('highlighted');
        } else {
            rowElement.classList.remove('highlighted')
        }

        rowElement.addEventListener('click', () => {
            console.log(msg.row.id)
            sendDataToBackend(msg.row.id)
        })
    });
}

function sendDataToBackend(number) {
    // Construct the data to send
    const data = {
        number: number
    };

    // Make the POST request
    fetch('http://localhost:8080/api/highlightedrow', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(data => {
        console.log('Response from backend:', data);
    })
    .catch(error => {
        console.error('Error sending data to backend:', error);
    });
}

function fetchData() {
    fetch('localhost:8080/api/data')
        .then(response => response.json())
        .then(data => {
            displayRows(data);
        })
        .catch(error => console.error('Error fetching data:', error));
}


setInterval(function() {
    fetchData();
}, 1000);
