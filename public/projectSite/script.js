// script.js

fetchTables();

function fetchTables() {
    const container = document.getElementById('mytables');
    container.innerHTML = '';
    fetch('/api/gettables')
        .then(response => response.json())
        .then(data => {
            if (data != null) {
                console.log(data);
                data.forEach(table => {
                    let rowElement = document.createElement('a');
                    rowElement.href = `/${table.table}/`;
                    rowElement.textContent = table.table;
                    rowElement.classList.add('row');
                    container.appendChild(rowElement);
                });
            } else {
                let element = document.createElement('h1');
                element.textContent = 'No Tables found';
                container.appendChild(element);
            }
        })  
        .catch(error => console.error('Error fetching data:', error));
}

document.getElementById('dataForm').addEventListener('submit', function(event) {
    event.preventDefault();
    const data = document.getElementById('dataInput').value;
    fetch('/api/newtable', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ name: data })
    })
    .then(response => response.json())
    .then(data => console.log('Success:', data))
    .catch((error) => console.error('Error:', error));
    fetchTables();
});