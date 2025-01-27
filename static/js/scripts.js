function checkoutBook(id) {
    fetch(`/books/${id}/checkout`, { method: 'PATCH' })
        .then(response => response.json())
        .then(data => {
            console.log('Book checked out:', data);
            alert('Book checked out successfully!');
        })
        .catch(error => console.error('Error:', error));
}

function returnBook(id) {
    fetch(`/books/${id}/return`, { method: 'PATCH' })
        .then(response => response.json())
        .then(data => {
            console.log('Book returned:', data);
            alert('Book returned successfully!');
        })
        .catch(error => console.error('Error:', error));
}
