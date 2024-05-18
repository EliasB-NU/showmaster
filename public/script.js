// script.js

const ws = new WebSocket('/ws/');

ws.onopen = function(event) {
    console.log('WebSocket connection established.');
};

ws.onmessage = function(event) {
    console.log(event.data);
    const msgs = JSON.parse(event.data);
    if (msgs === "refresh") {
        fetchData();
    } else {
        displayRows(msgs);
    }
    
};

var highlightedRow = -1
var rows = []

function displayRows(msgs) {
    msgs.sort((a, b) => a.row.id - b.row.id);

    const container = document.getElementById('tableBody');
    container.innerHTML = ''; // Clear existing content    

    msgs.forEach(msg => {
        rows.push(msg.row.id);
        const tr = document.createElement('tr');
        tr.dataset.id = msg.row.id;
        tr.innerHTML = `
            <td>${msg.row.id}</td>
            <td>${msg.row.name}</td>
            <td>${msg.row.audio}</td>
            <td>${msg.row.licht}</td>
            <td>${msg.row.pptx}</td>
            <td>${msg.row.notes}</td>
        `;
        tableBody.appendChild(tr);

        if (msg.highlighted === true) {
            const highlightedRowVar = document.querySelector(`tr[data-id="${msg.row.id}"]`);
            setTimeout(() => {
                highlightedRowVar.scrollIntoView({ behavior: 'smooth', block: 'center' });
            }, 100);
            highlightedRowVar.classList.add('highlighted');
            highlightedRow = rows.indexOf(msg.row.id);
        } else {
            tr.classList.remove('highlighted')
        }

        tr.addEventListener('click', () => {
            console.log(msg.row.id)
            sendDataToBackend(msg.row.id)
        })
    });
}

document.addEventListener('keydown', (event) => {
    if (event.key === 'ArrowDown') {
        event.preventDefault();
        console.log(rows[highlightedRow+1]);
        sendDataToBackend(rows[highlightedRow+1]);
    } else if (event.key === 'ArrowUp') {
        event.preventDefault();
        console.log(rows[highlightedRow-1]);
        sendDataToBackend(rows[highlightedRow-1]);
    }
})

function sendDataToBackend(number) {
    // Construct the data to send
    const data = {
        number: number
    };

    // Make the POST request
    fetch('/api/highlightedrow', {
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
    fetch('/api/data')
        .then(response => response.json())
        .then(data => {
            displayRows(data);
        })
        .catch(error => console.error('Error fetching data:', error));
}