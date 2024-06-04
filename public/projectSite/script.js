// script.js

fetchTables();

function fetchTables() {
    fetch('/api/gettables')
        .then(response => response.json())
        .then(data => {
            console.log(data);
            const container = document.getElementById('mytables');
            container.innerHTML = '';

            data.forEach(table => {
                let rowElement = document.createElement('a');
                rowElement.href = `/${table.table}/`;
                rowElement.textContent = table.table;
                rowElement.classList.add('row');
                container.appendChild(rowElement);
            });
        })
        .catch(error => console.error('Error fetching data:', error));
}