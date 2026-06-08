// src/api/books.js
import { fetchWithAuth } from './client';

export const getCatalog = () => fetchWithAuth('/books', { method: 'GET' });

export const addBook = (title, author, totalCopies) => {
    return fetchWithAuth('/books/add', {
        method: 'POST',
        body: JSON.stringify({
            title,
            author,
            total_copies: parseInt(totalCopies),
            available_copies: parseInt(totalCopies)
        })
    });
};

// FIX 1: Attach ID directly to the URL query string
export const deleteBook = (bookId) => {
    return fetchWithAuth(`/books/delete?id=${bookId}`, {
        method: 'POST'
    });
};


export const updateBook = (bookId, title, author, totalCopies) => {
    return fetchWithAuth(`/books/update?id=${bookId}`, {
        method: 'POST',
        body: JSON.stringify({
            title,
            author,
            total_copies: parseInt(totalCopies)
        })
    });
};