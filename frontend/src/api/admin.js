// src/api/admin.js
import { fetchWithAuth } from './client';


export const getPendingAdmins = () => {
    return fetchWithAuth('/users/promote/queue', { method: 'GET' });
};

export const approvePromotion = (userId) => {
    return fetchWithAuth('/users/promote/approve', {
        method: 'POST',
        body: JSON.stringify({ user_id: parseInt(userId) })
    });
};


export const getPendingRequests = () => {
    return fetchWithAuth('/transactions/queue', { method: 'GET' });
};

export const approveCheckout = (txId, bookId) => {
    return fetchWithAuth('/transactions/approve', {
        method: 'POST',
        body: JSON.stringify({ transaction_id: parseInt(txId), book_id: parseInt(bookId) })
    });
};


export const approveCheckin = (txId) => {
    return fetchWithAuth('/transactions/checkin/approve', {
        method: 'POST',
        body: JSON.stringify({ transaction_id: parseInt(txId) })
    });
};


export const rejectCheckout = (txId) => {
    return fetchWithAuth('/transactions/reject/checkout', {
        method: 'POST',
        body: JSON.stringify({ transaction_id: parseInt(txId) })
    });
};

export const rejectCheckin = (txId) => {
    return fetchWithAuth('/transactions/reject/checkin', {
        method: 'POST',
        body: JSON.stringify({ transaction_id: parseInt(txId) })
    });
};