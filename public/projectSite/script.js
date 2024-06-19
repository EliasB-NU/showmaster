// script.js

document.addEventListener('DOMContentLoaded', function() {
    fetch('/api/gettables')
        .then(response => response.json())
        .then(data => {
            const dropdown = document.getElementById('dropdownMenu');
            data.forEach(item => {
                const option = document.createElement('option');
                option.value = item.table;
                option.textContent = item.table;
                dropdown.appendChild(option);
            });
        })
        .catch((error) => console.error('Error fetching dropdown data:', error));
        
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
            
    document.getElementById('tableForm').addEventListener('submit', function(event) {
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
        
    document.getElementById('newInsertForm').addEventListener('submit', function(event) {
        event.preventDefault();
        const table = document.getElementById('dropdownMenu').value;
        const idString = document.getElementById('id').value;
        const id = Number(idString);
        const name = document.getElementById('name').value;
        const audio = document.getElementById('audio').value;
        const licht = document.getElementById('licht').value;
        const pptx = document.getElementById('pptx').value;
        const notes = document.getElementById('notes').value;

        console.log(id);
            
        fetch('/api/newinsert:'+table, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                id: id,
                name: name,
                audio: audio,
                licht: licht,
                pptx: pptx,
                notes: notes,
            })
        })
        .then(response => response.json())
            .then(data => console.log('Success:', data))
            .catch((error) => console.error('Error:', error));
    });
});

