var i = {
    id: Number,
    name: String,
    audio: String,
    licht: String,
    pptx: String,
    notes: String
}

document.getElementById('newInsert').addEventListener('submit', function(event) {
    event.preventDefault();
    const formData = new FormData(event.data);
    console.log(formData);
    newInsert(formData);
})

function newInsert(data) {
    fetch('/api/newinsert/', {
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
    })
}