// script.js
function extractSubpath(url) {
    const anchor = document.createElement('a');
    anchor.href = url;
    return anchor.pathname;
}
const currentURL = window.location.href;
const subpath = extractSubpath(currentURL);
const withoutfirstletter = subpath.slice(1);
const cleanURL = withoutfirstletter.substring(0, withoutfirstletter.length-1);
console.log(cleanURL);

document.addEventListener("DOMContentLoaded", () => {
    const headlineElement = document.getElementById('headline');
    headlineElement.textContent = 'ShowMaster - V3 |  '+cleanURL;

    const ws = new WebSocket('/ws');

    ws.onopen = function(event) {
        console.log('WebSocket connection established.');
    };

    ws.onmessage = function(event) {
        console.log(event.data);
        const msgs = JSON.parse(event.data);
        if (msgs === cleanURL+":refresh") {
            fetchData();
        } else if (msgs === cleanURL+":reset") {
            resetStopwatch();
        } else if (msgs === cleanURL+":start") {
            startStopwatch();
        } else if (msgs === cleanURL+":stop") {
            stopStopwatch();
            onStopUpdate();
        }    
    };

    const stopwatchElement = document.getElementById('timer');
    const startPauseBtn = document.getElementById("startPauseBtn");
    const resetBtn = document.getElementById("resetBtn");
    let duration = 0;
    let running = false;
    let intervalId = null;

    // Function to format the duration into HH:MM:SS
    function formatDuration(duration) {
        const hours = String(Math.floor(duration / 3600)).padStart(2, '0');
        const minutes = String(Math.floor((duration % 3600) / 60)).padStart(2, '0');
        const seconds = String(duration % 60).padStart(2, '0');
        return `${hours}:${minutes}:${seconds}`;
    }

    // Function to update the stopwatch display
    function updateStopwatch() {
        stopwatchElement.textContent = formatDuration(duration);
    }

    // Function to start the stopwatch
    function startStopwatch() {
        if (!running) {
            running = true;
            intervalId = setInterval(() => {
                duration += 1;
                updateStopwatch();
            }, 1000);
            startPauseBtn.textContent = "Stop";
        }
    }

    // Function to stop the stopwatch
    function stopStopwatch() {
        if (running) {
            running = false;
            clearInterval(intervalId);
            intervalId = null;
            startPauseBtn.textContent = "Start";
        }
    }

    // Function to reset the stopwatch
    function resetStopwatch() {
        stopStopwatch();
        duration = 0;
        updateStopwatch();
        startPauseBtn.textContent = "Start";
    }

    // Initialize the stopwatch with data from the backend
    fetch('/api/stopwatch-status:'+cleanURL)
        .then(response => response.json())
        .then(data => {
            duration = data.Duration;
            updateStopwatch();
            if (data.Running) {
                startStopwatch();
                startPauseBtn.textContent = "Stop"
            } 
        })
        .catch(error => console.error('Error fetching stopwatch status:', error));

    function onStopUpdate() {
        fetch('/api/stopwatch-status:'+cleanURL)
        .then(response => response.json())
        .then(data => {
            duration = data.Duration;
            updateStopwatch();
            if (data.Running) {
                startStopwatch();
            }
        })
        .catch(error => console.error('Error fetching stopwatch status:', error));
    }

    startPauseBtn.addEventListener("click", () => {
        if (running) { // To stop
            const data = {
                RUNNING: false,
                RESET: false,
            };
    
            // Make the POST request
            fetch('/api/stopwatch-update:'+cleanURL, {
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
            startPauseBtn.textContent = "Start";
        } else { // To start
            const data = {
                RUNNING: true,
                RESET: false,
            };
    
            // Make the POST request
            fetch('/api/stopwatch-update:'+cleanURL, {
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
            startPauseBtn.textContent = "Pause";
        }
    });

    resetBtn.addEventListener("click", () => {
        const data = {
            RUNNING: false,
            RESET: true,
        };

        // Make the POST request
        fetch('/api/stopwatch-update:'+cleanURL, {
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
        startPauseBtn.textContent = "Start";
    });

    fetch('/api/getdata:'+cleanURL)
    .then(response => response.json())
    .then(data => {
        displayRows(data);
    })
    .catch(error => console.error('Error fetching data:', error));

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

    function sendDataToBackend(num) {
        // Construct the data to send
        const data = {
            number: num
        };

        // Make the POST request
        fetch('/api/updatehighlightedrow:'+cleanURL, {
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
        fetch('/api/getdata:'+cleanURL)
        .then(response => response.json())
        .then(data => {
            displayRows(data);
        })
        .catch(error => console.error('Error fetching data:', error));
    }
});