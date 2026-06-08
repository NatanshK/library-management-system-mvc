// src/pages/StudentDashboard.jsx
import { useContext, useState, useEffect } from 'react';
import { AuthContext } from '../context/AuthContext';
import { requestCheckout, requestCheckin, getUserTransactions, requestAdminStatus } from '../api/transactions';
import { getCatalog } from '../api/books';


export default function StudentDashboard() {
    const { user, logout } = useContext(AuthContext);
    const [history, setHistory] = useState([]);
    const [catalog, setCatalog] = useState([]);
    const [searchQuery, setSearchQuery] = useState('');
    const [message, setMessage] = useState('');
    const [error, setError] = useState('');

    const loadData = async () => {
        try {
            const txData = await getUserTransactions(user.user_id);
            setHistory(txData || []);

            const bookData = await getCatalog();
            const safeBooks = (bookData || []).map(b => ({
                ...b,
                title: b.title || '',
                author: b.author || ''
            }));
            setCatalog(safeBooks);
        } catch (err) {
            console.error("Failed to load data", err);
        }
    };

    useEffect(() => {
        loadData();
    }, []);

    // Direct Action Handlers
    const handleCheckout = async (bookId) => {
        setMessage('');
        setError('');
        try {
            await requestCheckout(user.user_id, bookId);
            setMessage("Checkout requested successfully! Pending admin approval.");
            loadData();
        } catch (err) {
            setError('Error: ' + err.message);
        }
    };

    const handleReturn = async (txId) => {
        setMessage('');
        setError('');
        try {
            await requestCheckin(user.user_id, txId);
            setMessage("Return requested successfully! Pending admin verification.");
            loadData();
        } catch (err) {
            setError('Error: ' + err.message);
        }
    };

    const handleAdminRequest = async () => {
        setMessage('');
        setError('');
        try {
            await requestAdminStatus(user.user_id);
            setMessage("Admin status requested! Please wait for an existing admin to approve.");
        } catch (err) {
            setError('Error: ' + err.message);
        }
    };

    const filteredCatalog = catalog.filter(book =>
        book.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
        book.author.toLowerCase().includes(searchQuery.toLowerCase())
    );

    return (
        <div className="p-8 max-w-6xl mx-auto">
            <div className="flex justify-between items-center mb-8">
                <h1 className="text-3xl font-bold">Student Dashboard</h1>
                <div className="flex gap-4">
                    <button
                        onClick={handleAdminRequest}
                        className="bg-purple-600 text-white px-4 py-2 rounded-lg hover:bg-purple-700 transition-colors font-medium shadow-sm"
                    >
                        Request Admin Access
                    </button>
                    <button
                        onClick={logout}
                        className="bg-gray-800 text-white px-4 py-2 rounded-lg hover:bg-gray-700 transition-colors shadow-sm"
                    >
                        Logout
                    </button>
                </div>
            </div>

            {message && <div className="bg-green-100 border border-green-400 text-green-800 p-3 rounded-lg mb-6 text-center font-semibold">{message}</div>}
            {error && <div className="bg-red-100 border border-red-400 text-red-800 p-3 rounded-lg mb-6 text-center font-semibold">{error}</div>}

            {/* Catalog Panel (Now takes full width) */}
            <div className="bg-white shadow-sm rounded-xl border mb-8 p-6">
                <div className="flex justify-between items-center mb-4">
                    <h2 className="font-bold text-xl">Library Catalog</h2>
                    <input
                        type="text"
                        placeholder="Search by title or author..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="border p-2 rounded-lg w-64 focus:ring-2 focus:ring-blue-500 outline-none"
                    />
                </div>
                <div className="overflow-y-auto max-h-80 border rounded-lg">
                    <table className="w-full text-left text-sm">
                        <thead className="bg-gray-50 sticky top-0 shadow-sm">
                            <tr>
                                <th className="p-3">Title</th>
                                <th className="p-3">Author</th>
                                <th className="p-3">Availability</th>
                                <th className="p-3 text-right">Action</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y">
                            {filteredCatalog.map(book => (
                                <tr key={book.id} className="hover:bg-gray-50 transition-colors">
                                    <td className="p-3 font-medium">{book.title}</td>
                                    <td className="p-3 text-gray-600">{book.author}</td>
                                    <td className="p-3">
                                        <span className={`px-2 py-1 rounded-full text-xs font-semibold ${book.available_copies > 0 ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>
                                            {book.available_copies} / {book.total_copies}
                                        </span>
                                    </td>
                                    <td className="p-3 text-right">
                                        <button
                                            onClick={() => handleCheckout(book.id)}
                                            disabled={book.available_copies <= 0}
                                            className="bg-blue-600 text-white px-4 py-1.5 rounded-lg text-sm font-medium hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors"
                                        >
                                            {book.available_copies > 0 ? 'Request' : 'Out of Stock'}
                                        </button>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            </div>

            {/* History Table (Now features Return buttons) */}
            <div className="bg-white shadow-sm rounded-xl border overflow-hidden">
                <h2 className="p-6 text-xl font-bold border-b bg-gray-50">Your Borrowing History</h2>
                <div className="overflow-y-auto max-h-80 border rounded-b-xl">
                    <table className="w-full text-left text-sm">
                        <thead className="bg-gray-50 sticky top-0 shadow-sm">
                            <tr>
                                <th className="p-4">Tx ID</th>
                                <th className="p-4">Book</th>
                                <th className="p-4">Status</th>
                                <th className="p-4">Fine</th>
                                <th className="p-4 text-right">Action</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y">
                            {history.length > 0 ? history.map((tx) => (
                                <tr key={tx.transaction_id} className="hover:bg-gray-50 transition-colors">
                                    <td className="p-4 font-mono text-gray-500">{tx.transaction_id}</td>
                                    <td className="p-4 font-medium">{tx.book_title}</td>
                                    <td className="p-4 text-gray-600">
                                        <span className="bg-gray-100 text-gray-700 px-2 py-1 rounded text-xs font-semibold">
                                            {tx.status}
                                        </span>
                                    </td>
                                    <td className="p-4 text-red-600 font-medium">${tx.fine_amount?.toFixed(2) || "0.00"}</td>
                                    <td className="p-4 text-right">
                                        {/* Only show Return button if the book is actually checked out! */}
                                        {tx.status === 'checkout_accepted' ? (
                                            <button
                                                onClick={() => handleReturn(tx.transaction_id)}
                                                className="bg-green-600 hover:bg-green-700 text-white px-4 py-1.5 rounded-lg text-sm font-medium transition-colors"
                                            >
                                                Return Book
                                            </button>
                                        ) : (
                                            <span className="text-gray-400 text-sm">N/A</span>
                                        )}
                                    </td>
                                </tr>
                            )) : (
                                <tr>
                                    <td colSpan="5" className="p-6 text-center text-gray-500">No borrowing history found.</td>
                                </tr>
                            )}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
}