// src/api/transactions.js
import { fetchWithAuth } from './client';

export const requestCheckout = (userId, bookId) => {
    return fetchWithAuth('/transactions/request', {
        method: 'POST',
        body: JSON.stringify({
            user_id: parseInt(userId),
            book_id: parseInt(bookId)
        })
    });
};

export const requestCheckin = (userId, txId) => {
    return fetchWithAuth('/transactions/return', {
        method: 'POST',
        body: JSON.stringify({
            user_id: parseInt(userId),
            transaction_id: parseInt(txId)
        })
    });
};

export const getUserTransactions = (userId) => {
    return fetchWithAuth(`/transactions/history?user_id=${userId}`, {
        method: 'GET'
    });
};

export const requestAdminStatus = (userId) => {
    return fetchWithAuth('/users/promote/request', {
        method: 'POST',
        body: JSON.stringify({ user_id: parseInt(userId) })
    });
};