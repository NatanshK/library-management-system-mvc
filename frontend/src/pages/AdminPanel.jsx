// src/pages/AdminPanel.jsx
import { useEffect, useState, useContext } from 'react';
import { AuthContext } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';
import { getPendingRequests, approveCheckout, approveCheckin, rejectCheckout, rejectCheckin, approvePromotion, getPendingAdmins } from '../api/admin';
import { getCatalog, addBook, deleteBook, updateBook } from '../api/books'; // Added updateBook

export default function AdminPanel() {
    const { logout } = useContext(AuthContext);
    const navigate = useNavigate();

    const [queue, setQueue] = useState([]);
    const [hrQueue, setHrQueue] = useState([]);
    const [catalog, setCatalog] = useState([]);

    // Smart Form State (Handles both Add and Edit)
    const [newBook, setNewBook] = useState({ title: '', author: '', copies: '' });
    const [isEditing, setIsEditing] = useState(false);
    const [editBookId, setEditBookId] = useState(null);

    const [message, setMessage] = useState('');
    const [error, setError] = useState('');

    const loadData = async () => {

        try {
            const queueData = await getPendingRequests();
            setQueue(queueData || []);
        } catch (err) {
            console.error("Queue load failed:", err);
        }


        try {
            const catalogData = await getCatalog();
            setCatalog(catalogData || []);
        } catch (err) {
            console.error("Catalog load failed:", err);
        }


        try {
            const hrData = await getPendingAdmins();
            setHrQueue(hrData || []);
        } catch (err) {
            console.error("HR Queue load failed:", err);

        }
    };

    useEffect(() => {
        loadData();
    }, []);

    const handleAction = async (actionFn, id, successMsg) => {
        try {
            await actionFn(id);
            setMessage(successMsg);
            loadData();
        } catch (err) {
            setError(err.message);
        }
    };

    const handleApproveCheckout = async (txId, bookId) => {
        try {
            await approveCheckout(txId, bookId);
            setMessage("Checkout Approved!");
            loadData();
        } catch (err) {
            setError(err.message);
        }
    };

    // The Smart Submit Handler
    const handleBookSubmit = async (e) => {
        e.preventDefault();
        setMessage('');
        setError('');
        try {
            if (isEditing) {
                await updateBook(editBookId, newBook.title, newBook.author, newBook.copies);
                setMessage("Book updated successfully!");
            } else {
                await addBook(newBook.title, newBook.author, newBook.copies);
                setMessage("Book added to catalog!");
            }
            // Reset Form
            setNewBook({ title: '', author: '', copies: '' });
            setIsEditing(false);
            setEditBookId(null);
            loadData();
        } catch (err) {
            setError(err.message);
        }
    };

    // Populates the form when "Edit" is clicked
    const handleEditClick = (book) => {
        setIsEditing(true);
        setEditBookId(book.id);
        setNewBook({ title: book.title, author: book.author, copies: book.total_copies });
        window.scrollTo({ top: 0, behavior: 'smooth' }); // Scroll up to the form
    };

    // Cancels the edit and clears the form
    const cancelEdit = () => {
        setIsEditing(false);
        setEditBookId(null);
        setNewBook({ title: '', author: '', copies: '' });
    };

    return (
        <div className="p-8 max-w-7xl mx-auto">
            <div className="flex justify-between items-center mb-8">
                <h1 className="text-3xl font-bold text-gray-800">Librarian Command Center</h1>
                <button onClick={() => { logout(); navigate('/login'); }} className="bg-gray-800 hover:bg-gray-700 text-white px-4 py-2 rounded-lg">Logout</button>
            </div>

            {message && <div className="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded mb-6 text-center font-semibold">{message}</div>}
            {error && <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6 text-center font-semibold">{error}</div>}

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mb-8">

                {/* HR Queue - SCROLLABLE */}
                <div className="bg-white rounded-xl shadow-sm border lg:col-span-1 p-6 h-[400px] flex flex-col">
                    <h2 className="font-bold text-xl mb-4">Pending HR Requests</h2>
                    <div className="overflow-y-auto flex-1 border rounded-lg">
                        <table className="w-full text-left text-sm">
                            <thead className="bg-gray-50 sticky top-0">
                                <tr>
                                    <th className="p-3">User</th>
                                    <th className="p-3">Action</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y">
                                {hrQueue.length > 0 ? hrQueue.map((user) => (
                                    <tr key={user.id} className="hover:bg-gray-50">
                                        <td className="p-3">
                                            <div className="font-medium">{user.username}</div>
                                            <div className="text-xs text-gray-500">{user.email}</div>
                                        </td>
                                        <td className="p-3">
                                            <button onClick={() => handleAction(approvePromotion, user.id, "Admin Promoted!")} className="bg-blue-600 text-white px-3 py-1 rounded hover:bg-blue-700">Approve</button>
                                        </td>
                                    </tr>
                                )) : <tr><td colSpan="2" className="p-4 text-center text-gray-500">No pending requests</td></tr>}
                            </tbody>
                        </table>
                    </div>
                </div>

                {/* Inventory Management - SCROLLABLE */}
                <div className="bg-white p-6 rounded-xl shadow-sm border lg:col-span-2 h-[400px] flex flex-col">
                    <div className="flex justify-between items-center mb-4">
                        <h2 className="font-bold text-xl">{isEditing ? "Edit Book" : "Add New Book"}</h2>
                        {isEditing && <button onClick={cancelEdit} className="text-sm text-gray-500 underline hover:text-gray-800">Cancel Edit</button>}
                    </div>

                    <form onSubmit={handleBookSubmit} className="flex gap-2 mb-4">
                        <input type="text" placeholder="Title" required value={newBook.title} onChange={(e) => setNewBook({ ...newBook, title: e.target.value })} className="flex-2 border p-2 rounded outline-none focus:ring-2 focus:ring-blue-500" />
                        <input type="text" placeholder="Author" required value={newBook.author} onChange={(e) => setNewBook({ ...newBook, author: e.target.value })} className="flex-2 border p-2 rounded outline-none focus:ring-2 focus:ring-blue-500" />
                        <input type="number" placeholder="Total Copies" required value={newBook.copies} onChange={(e) => setNewBook({ ...newBook, copies: e.target.value })} className="w-28 border p-2 rounded outline-none focus:ring-2 focus:ring-blue-500" />
                        <button type="submit" className={`px-4 py-2 rounded text-white font-medium ${isEditing ? 'bg-orange-500 hover:bg-orange-600' : 'bg-green-600 hover:bg-green-700'}`}>
                            {isEditing ? 'Update' : 'Add'}
                        </button>
                    </form>

                    <div className="overflow-y-auto flex-1 border rounded-lg">
                        <table className="w-full text-left text-sm">
                            <thead className="bg-gray-50 sticky top-0 shadow-sm">
                                <tr>
                                    <th className="p-3">Title</th>
                                    <th className="p-3">Copies</th>
                                    <th className="p-3 text-right">Actions</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y">
                                {catalog.map(book => (
                                    <tr key={book.id} className={`hover:bg-gray-50 ${editBookId === book.id ? 'bg-orange-50' : ''}`}>
                                        <td className="p-3">
                                            <div className="font-medium">{book.title}</div>
                                            <div className="text-xs text-gray-500">{book.author}</div>
                                        </td>
                                        <td className="p-3">{book.available_copies} / {book.total_copies}</td>
                                        <td className="p-3 text-right">
                                            <button onClick={() => handleEditClick(book)} className="text-blue-600 hover:text-blue-800 font-medium mr-4">Edit</button>
                                            <button onClick={() => handleAction(deleteBook, book.id, "Book deleted!")} className="text-red-600 hover:text-red-800 font-medium">Remove</button>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>

            {/* Book Transactions Queue - SCROLLABLE */}
            <div className="bg-white shadow rounded-lg p-6 border flex flex-col">
                <h2 className="text-xl font-bold mb-4">Book Transactions Queue</h2>
                <div className="overflow-y-auto max-h-96 border rounded-lg">
                    <table className="w-full text-left">
                        <thead className="bg-gray-50 sticky top-0 shadow-sm">
                            <tr>
                                <th className="p-4">User</th>
                                <th className="p-4">Book</th>
                                <th className="p-4">Status</th>
                                <th className="p-4 text-center">Actions</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y">
                            {queue.length > 0 ? queue.map((req) => (
                                <tr key={req.transaction_id} className="hover:bg-gray-50 transition-colors">
                                    <td className="p-4">
                                        <div className="font-medium">{req.user_name}</div>
                                        <div className="text-xs text-gray-500">{req.user_email}</div>
                                    </td>
                                    <td className="p-4 font-medium">{req.book_title}</td>
                                    <td className="p-4"><span className="bg-yellow-100 text-yellow-800 px-2 py-1 rounded text-sm font-semibold">{req.status}</span></td>
                                    <td className="p-4 flex justify-center gap-2">
                                        {req.status === 'checkout_requested' && (
                                            <>
                                                <button onClick={() => handleApproveCheckout(req.transaction_id, req.book_id)} className="bg-green-500 hover:bg-green-600 text-white px-3 py-1 rounded text-sm transition-colors">Approve Checkout</button>
                                                <button onClick={() => handleAction(rejectCheckout, req.transaction_id, "Checkout Rejected")} className="bg-red-100 hover:bg-red-200 text-red-700 px-3 py-1 rounded text-sm transition-colors">Deny</button>
                                            </>
                                        )}
                                        {req.status === 'checkin_requested' && (
                                            <>
                                                <button onClick={() => handleAction(approveCheckin, req.transaction_id, "Return Approved!")} className="bg-blue-500 hover:bg-blue-600 text-white px-3 py-1 rounded text-sm transition-colors">Approve Return</button>
                                                <button onClick={() => handleAction(rejectCheckin, req.transaction_id, "Return Rejected")} className="bg-orange-100 hover:bg-orange-200 text-orange-700 px-3 py-1 rounded text-sm transition-colors">Deny</button>
                                            </>
                                        )}
                                    </td>
                                </tr>
                            )) : <tr><td colSpan="4" className="p-8 text-center text-gray-500">No pending requests at the moment.</td></tr>}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
}