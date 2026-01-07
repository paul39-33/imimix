// Auth Check
const token = localStorage.getItem('token');
const user = JSON.parse(localStorage.getItem('user') || '{}');

if (!token) {
    window.location.href = 'index.html';
}

document.getElementById('user-display').textContent = user.username || 'User';

// Logout
document.getElementById('logoutBtn').addEventListener('click', () => {
    localStorage.clear();
    window.location.href = 'index.html';
});

// Helper Functions for Date
function formatDateToDDMMYYYY(isoStr) {
    if (!isoStr) return '';
    const d = new Date(isoStr);
    if (isNaN(d.getTime())) return '';
    const day = String(d.getDate()).padStart(2, '0');
    const month = String(d.getMonth() + 1).padStart(2, '0');
    const year = d.getFullYear();
    return `${day}/${month}/${year}`;
}

function parseDDMMYYYYToISO(dateStr) {
    if (!dateStr) return null;
    const parts = dateStr.split('/');
    if (parts.length === 3) {
        const day = parseInt(parts[0], 10);
        const month = parseInt(parts[1], 10) - 1;
        const year = parseInt(parts[2], 10);
        const dateObj = new Date(year, month, day);
        if (!isNaN(dateObj.getTime())) {
            return dateObj.toISOString();
        }
    }
    // Fallback: try standard parse
    const d = new Date(dateStr);
    if (!isNaN(d.getTime())) return d.toISOString();
    return null;
}

// State
let allObjects = [];
let currentPage = 1;
const itemsPerPage = 8; // Adjust as needed

// Elements
const searchInput = document.getElementById('searchInput');
const searchBtn = document.getElementById('searchBtn');
const tableBody = document.getElementById('tableBody');
const paginationControls = document.getElementById('paginationControls');

// Initial Load
fetchObjects();

// Search Event
searchBtn.addEventListener('click', () => {
    const query = searchInput.value.trim();
    fetchObjects(query);
});
// Allow Enter key
searchInput.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
        const query = searchInput.value.trim();
        fetchObjects(query);
    }
});

// Add New Object (Placeholder)
// Add New Request Modal Logic
const addReqBtn = document.getElementById('addReqBtn');
const modal = document.getElementById('add-req-modal');
const closeModal = document.querySelector('.close-modal');
const addReqForm = document.getElementById('add-req-form');

function toggleModal(show) {
    modal.style.display = show ? 'block' : 'none';
}

if (addReqBtn) {
    addReqBtn.addEventListener('click', () => toggleModal(true));
}

if (closeModal) {
    closeModal.addEventListener('click', () => toggleModal(false));
}

window.addEventListener('click', (event) => {
    if (event.target == modal) {
        toggleModal(false);
    }
});

// Handle Add Request Submit
if (addReqForm) {
    addReqForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const formData = new FormData(addReqForm);
        const data = Object.fromEntries(formData.entries());

        if (data.promote_date) {
            const parts = data.promote_date.split('/');
            if (parts.length === 3) {
                const day = parseInt(parts[0], 10);
                const month = parseInt(parts[1], 10) - 1;
                const year = parseInt(parts[2], 10);
                const dateObj = new Date(year, month, day);
                if (!isNaN(dateObj.getTime())) {
                    data.promote_date = dateObj.toISOString();
                } else {
                    alert('Invalid date format. Please use dd/mm/yyyy.');
                    return;
                }
            } else {
                // Try standard parse as fallback
                const dateObj = new Date(data.promote_date);
                if (!isNaN(dateObj.getTime())) {
                    data.promote_date = dateObj.toISOString();
                }
            }
        }

        try {
            const response = await fetch('/api/create_obj_req', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(data)
            });

            const result = await response.json();

            if (!response.ok) {
                throw new Error(result.error || 'Failed to create request');
            }

            alert('Request created successfully!');
            toggleModal(false);
            addReqForm.reset();
            // Switch to requests tab and refresh
            switchTab('requests');
        } catch (error) {
            console.error('Error creating request:', error);
            alert('Error: ' + error.message);
        }
    });
}

// Add Object Modal Logic
const addObjBtn = document.getElementById('addObjBtn');
const objModal = document.getElementById('add-obj-modal');
// We need to handle multiple modals. Let's make toggleModal accept an element ID or generic logic.
// Simpler: Just replicate logic or generalize.
// Since we have specific IDs, let's replicate for clarity.

function toggleObjModal(show) {
    if (objModal) objModal.style.display = show ? 'block' : 'none';
}

if (addObjBtn) {
    addObjBtn.addEventListener('click', () => toggleObjModal(true));
}

// Close button for Obj Modal (assuming same class .close-modal but inside specific modal)
// Actually querySelector('.close-modal') only gets the first one. We need to iterate or be specific.
document.querySelectorAll('.close-modal').forEach(btn => {
    btn.addEventListener('click', () => {
        toggleModal(false);
        toggleObjModal(false);
    });
});

window.addEventListener('click', (event) => {
    if (event.target == objModal) {
        toggleObjModal(false);
    }
});

const addObjForm = document.getElementById('add-obj-form');
if (addObjForm) {
    addObjForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const formData = new FormData(addObjForm);
        const data = Object.fromEntries(formData.entries());

        if (data.promote_date) {
            const parts = data.promote_date.split('/');
            if (parts.length === 3) {
                const day = parseInt(parts[0], 10);
                const month = parseInt(parts[1], 10) - 1;
                const year = parseInt(parts[2], 10);
                const dateObj = new Date(year, month, day);
                if (!isNaN(dateObj.getTime())) {
                    data.promote_date = dateObj.toISOString();
                } else {
                    alert('Invalid date format. Please use dd/mm/yyyy.');
                    return;
                }
            } else {
                // Try standard parse as fallback
                const dateObj = new Date(data.promote_date);
                if (!isNaN(dateObj.getTime())) {
                    data.promote_date = dateObj.toISOString();
                }
            }
        }

        try {
            const response = await fetch('/api/add_mimix_obj', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(data)
            });

            const result = await response.json();

            if (!response.ok) {
                throw new Error(result.error || 'Failed to create object');
            }

            alert('Object created successfully!');
            toggleObjModal(false);
            addObjForm.reset();
            switchTab('objects');
            fetchObjects();
        } catch (error) {
            console.error('Error creating object:', error);
            alert('Error: ' + error.message);
        }
    });
}

// Fetch Logic
async function fetchObjects(query = '') {
    try {
        // Construct URL: if query is empty, hit /api/obj/search, else /api/obj/search/query
        // Actually the backend might handle /api/obj/search/ (trailing slash) or we use the new route
        let url = '/api/obj/search';
        if (query) {
            url += `/${encodeURIComponent(query)}`;
        }

        const response = await fetch(url, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (response.status === 401) {
            // Token expired or invalid
            localStorage.clear();
            window.location.href = 'index.html';
            return;
        }

        if (!response.ok) throw new Error('Failed to fetch objects');

        const data = await response.json();
        // API might return null if empty, default to array
        allObjects = data || [];

        // Reset to page 1 on new search
        currentPage = 1;
        renderTable();

    } catch (error) {
        console.error('Error fetching objects:', error);
        tableBody.innerHTML = `<tr><td colspan="8" style="text-align:center; color: #ef4444;">Error loading data: ${error.message}</td></tr>`;
    }
}

// Render Table
function renderTable() {
    tableBody.innerHTML = '';

    if (allObjects.length === 0) {
        tableBody.innerHTML = `<tr><td colspan="8" style="text-align:center; color: var(--text-muted);">No objects found.</td></tr>`;
        renderPagination(0);
        return;
    }

    // Pagination Logic
    const startIndex = (currentPage - 1) * itemsPerPage;
    const endIndex = startIndex + itemsPerPage;
    const pageItems = allObjects.slice(startIndex, endIndex);

    pageItems.forEach(obj => {
        const row = document.createElement('tr');
        row.dataset.id = obj.id; // Store ID for lookup

        // Format Date (simplified)
        const dateStr = obj.promote_date ? new Date(obj.promote_date).toLocaleDateString() : '-';

        row.innerHTML = `
            <td style="font-weight: 500; color: white;">${obj.obj}</td>
            <td>${obj.obj_type}</td>
            <td>${dateStr}</td>
            <td>${obj.lib}</td>
            <td>${obj.obj_ver}</td>
            <td><span class="status-badge status-temp">${obj.mimix_status}</span></td>
            <td>${obj.developer}</td>
            <td style="color: var(--text-muted); font-size: 0.875rem;">${obj.keterangan || '-'}</td>
            <td>
                <button class="action-btn" onclick="addToRequest('${obj.id}')" title="Add to Mimix Request">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="22" y1="2" x2="11" y2="13"></line><polygon points="22 2 15 22 11 13 2 9 22 2"></polygon></svg>
                </button>
                <button class="action-btn" onclick="editObject('${obj.id}')" title="Edit">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
                </button>
                <button class="action-btn delete" onclick="deleteObject('${obj.id}')" title="Delete">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><line x1="10" y1="11" x2="10" y2="17"></line><line x1="14" y1="11" x2="14" y2="17"></line></svg>
                </button>
            </td>
        `;
        tableBody.appendChild(row);
    });

    // Clean up old listeners (simple re-render handles usage here, but being safe)
    // Actually, simply replacing innerHTML clears old listeners on those elements.

    const totalPages = Math.ceil(allObjects.length / itemsPerPage);
    renderPagination(totalPages);
}

// Render Pagination Controls
function renderPagination(totalPages) {
    paginationControls.innerHTML = '';

    if (totalPages <= 1) return;

    // Previous
    const prevBtn = document.createElement('button');
    prevBtn.className = 'page-btn';
    prevBtn.textContent = '<';
    prevBtn.disabled = currentPage === 1;
    prevBtn.onclick = () => { if (currentPage > 1) { currentPage--; renderTable(); } };
    paginationControls.appendChild(prevBtn);

    // Page Numbers
    for (let i = 1; i <= totalPages; i++) {
        const btn = document.createElement('button');
        btn.className = `page-btn ${i === currentPage ? 'active' : ''}`;
        btn.textContent = i;
        btn.onclick = () => { currentPage = i; renderTable(); };
        paginationControls.appendChild(btn);
    }

    // Next
    const nextBtn = document.createElement('button');
    nextBtn.className = 'page-btn';
    nextBtn.textContent = '>';
    nextBtn.disabled = currentPage === totalPages;
    nextBtn.onclick = () => { if (currentPage < totalPages) { currentPage++; renderTable(); } };
    paginationControls.appendChild(nextBtn);
}


// Requests State
let allRequests = [];
let currentReqPage = 1;

// Request Elements
const reqSearchInput = document.getElementById('reqSearchInput');
const reqSearchBtn = document.getElementById('reqSearchBtn');
const reqTableBody = document.getElementById('reqTableBody');
const reqPaginationControls = document.getElementById('reqPaginationControls');

// Fetch Requests
async function fetchRequests(query = '') {
    try {
        let url = '/api/obj_req/search';
        if (query) {
            url += `/${encodeURIComponent(query)}`;
        }

        const response = await fetch(url, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (response.status === 401) {
            localStorage.removeItem('token');
            window.location.href = 'index.html';
            return;
        }

        if (!response.ok) {
            throw new Error('Failed to fetch requests');
        }

        allRequests = await response.json();
        if (!allRequests) allRequests = []; // Handle null/empty response
        currentReqPage = 1;
        renderRequestsTable();

    } catch (error) {
        console.error('Error fetching requests:', error);
        reqTableBody.innerHTML = `<tr><td colspan="10" style="text-align:center; color: var(--text-muted);">Error loading requests: ${error.message}</td></tr>`;
    }
}

// Render Requests Table
function renderRequestsTable() {
    reqTableBody.innerHTML = '';

    if (allRequests.length === 0) {
        reqTableBody.innerHTML = `<tr><td colspan="10" style="text-align:center; color: var(--text-muted);">No requests found.</td></tr>`;
        renderReqPagination(0);
        return;
    }

    const startIndex = (currentReqPage - 1) * itemsPerPage;
    const endIndex = startIndex + itemsPerPage;
    const pageItems = allRequests.slice(startIndex, endIndex);

    pageItems.forEach(req => {
        const row = document.createElement('tr');

        // Format Dates
        const updateDate = req.updated_at ? new Date(req.updated_at).toLocaleString() : '-';
        const promoteDate = req.promote_date ? new Date(req.promote_date).toLocaleDateString() : '-';

        row.innerHTML = `
            <td style="font-weight: 500; color: white;">${req.obj_name}</td>
            <td>${req.requester}</td>
            <td>${updateDate}</td>
            <td>${req.lib}</td>
            <td>${req.obj_ver}</td>
            <td>${req.obj_type}</td>
            <td>${promoteDate}</td>
            <td>${req.developer}</td>
            <td><span class="status-badge status-temp">${req.promote_status || '-'}</span></td>
            <td><span class="status-badge status-temp">${req.req_status}</span></td>
            <td>
                <button class="action-btn" onclick="convertRequest('${req.id}')" title="Convert to Object">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="12" y1="18" x2="12" y2="12"></line><line x1="9" y1="15" x2="15" y2="15"></line></svg>
                </button>
                <button class="action-btn" onclick="editRequest('${req.id}')" title="Edit">
                    ${icons.edit}
                </button>
                <button class="action-btn delete" onclick="deleteRequest('${req.id}')" title="Delete">
                    ${icons.delete}
                </button>
            </td>
        `;
        reqTableBody.appendChild(row);
    });

    const totalPages = Math.ceil(allRequests.length / itemsPerPage);
    renderReqPagination(totalPages);
}

// Render Request Pagination
function renderReqPagination(totalPages) {
    reqPaginationControls.innerHTML = '';
    if (totalPages <= 1) return;

    const prevBtn = document.createElement('button');
    prevBtn.className = 'page-btn';
    prevBtn.textContent = '<';
    prevBtn.disabled = currentReqPage === 1;
    prevBtn.onclick = () => { if (currentReqPage > 1) { currentReqPage--; renderRequestsTable(); } };
    reqPaginationControls.appendChild(prevBtn);

    for (let i = 1; i <= totalPages; i++) {
        const btn = document.createElement('button');
        btn.className = `page-btn ${i === currentReqPage ? 'active' : ''}`;
        btn.textContent = i;
        btn.onclick = () => { currentReqPage = i; renderRequestsTable(); };
        reqPaginationControls.appendChild(btn);
    }

    const nextBtn = document.createElement('button');
    nextBtn.className = 'page-btn';
    nextBtn.textContent = '>';
    nextBtn.disabled = currentReqPage === totalPages;
    nextBtn.onclick = () => { if (currentReqPage < totalPages) { currentReqPage++; renderRequestsTable(); } };
    reqPaginationControls.appendChild(nextBtn);
}


// Tab Switching Update
// Tab Switching Update
function switchTab(tabName) {
    // Content
    document.querySelectorAll('.tab-content').forEach(content => content.classList.remove('active'));

    // Buttons
    document.querySelectorAll('.tab-btn').forEach(btn => btn.classList.remove('active'));

    if (tabName === 'objects') {
        document.getElementById('objects-tab').classList.add('active');
        const btn = document.querySelector(`button[onclick="switchTab('objects')"]`);
        if (btn) btn.classList.add('active');
    } else {
        document.getElementById('requests-tab').classList.add('active');
        const btn = document.querySelector(`button[onclick="switchTab('requests')"]`);
        if (btn) btn.classList.add('active');
        fetchRequests(); // Fetch data when switching to requests
    }
}

// Event Listeners for Requests
reqSearchBtn.addEventListener('click', () => {
    fetchRequests(reqSearchInput.value);
});
reqSearchInput.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') fetchRequests(reqSearchInput.value);
});

// Original Tab Switching (Removing old function to replace with updated one above)
// Function above effectively overwrites it if placed correctly or we replace it.

// --- Action Functions ---

const icons = {
    save: `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"></path><polyline points="17 21 17 13 7 13 7 21"></polyline><polyline points="7 3 7 8 15 8"></polyline></svg>`,
    cancel: `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>`,
    edit: `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>`,
    delete: `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><line x1="10" y1="11" x2="10" y2="17"></line><line x1="14" y1="11" x2="14" y2="17"></line></svg>`,
    send: `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="22" y1="2" x2="11" y2="13"></line><polygon points="22 2 15 22 11 13 2 9 22 2"></polygon></svg>`
};

async function deleteObject(id) {
    if (!confirm('Are you sure you want to delete this object?')) return;

    try {
        const response = await fetch(`/api/delete_mimix_obj/${id}`, {
            method: 'DELETE',
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (!response.ok) throw new Error('Failed to delete object');

        fetchObjects(searchInput.value);
    } catch (error) {
        console.error('Delete error:', error);
        alert('Error deleting object');
    }
}

async function addToRequest(id) {
    if (!confirm('Add this object to Mimix Request?')) return;

    try {
        const response = await fetch(`/api/add_obj_to_obj_req/${id}`, {
            method: 'POST',
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (!response.ok) {
            const data = await response.json();
            throw new Error(data.error || 'Failed to add to request');
        }

        alert('Object added to request successfully');
        // Optionally switch tab: switchTab('requests');
    } catch (error) {
        console.error('Add Request error:', error);
        alert('Error: ' + error.message);
    }
}

// Global variable to store original row HTML
let editingRowId = null;
let originalRowHTML = '';

function editObject(id) {
    if (editingRowId) {
        cancelEdit(editingRowId);
    }

    const row = document.querySelector(`tr[data-id="${id}"]`);
    if (!row) return;

    editingRowId = id;
    originalRowHTML = row.innerHTML;

    // Find the object data from allObjects array
    const objData = allObjects.find(o => o.id === id);
    if (!objData) return;

    // Replace cells with inputs
    // Adjusted to allow editing Name (obj), Type (obj_type), and Promote Date
    row.innerHTML = `
        <td><input class="edit-input" id="edit-obj-${id}" value="${objData.obj}"></td>
        <td><input class="edit-input" id="edit-type-${id}" value="${objData.obj_type}"></td>
        <td><input class="edit-input" id="edit-date-${id}" value="${formatDateToDDMMYYYY(objData.promote_date)}" placeholder="dd/mm/yyyy"></td>
        
        <td><input class="edit-input" id="edit-lib-${id}" value="${objData.lib}"></td>
        <td><input class="edit-input" id="edit-ver-${id}" value="${objData.obj_ver}"></td>
        
        <td>
             <select class="edit-input" id="edit-status-${id}">
                <option value="pending" ${objData.mimix_status === 'pending' ? 'selected' : ''}>Pending</option>
                <option value="on progress" ${objData.mimix_status === 'on progress' ? 'selected' : ''}>On Progress</option>
                <option value="done" ${objData.mimix_status === 'done' ? 'selected' : ''}>Done</option>
                <option value="error" ${objData.mimix_status === 'error' ? 'selected' : ''}>Error</option>
            </select>
        </td>
        
        <td><input class="edit-input" id="edit-dev-${id}" value="${objData.developer}"></td>
        <td><input class="edit-input" id="edit-ket-${id}" value="${objData.keterangan || ''}"></td>
        
        <td>
            <button class="action-btn save" onclick="saveObject('${id}')" title="Save">${icons.save}</button>
            <button class="action-btn" onclick="cancelEdit('${id}')" title="Cancel">${icons.cancel}</button>
        </td>
    `;
}

function cancelEdit(id) {
    const row = document.querySelector(`tr[data-id="${id}"]`);
    if (row && editingRowId === id) {
        row.innerHTML = originalRowHTML;
        editingRowId = null;
        originalRowHTML = '';
    }
}

async function saveObject(id) {
    const objName = document.getElementById(`edit-obj-${id}`).value;
    const objType = document.getElementById(`edit-type-${id}`).value;
    const dateStr = document.getElementById(`edit-date-${id}`).value;
    const lib = document.getElementById(`edit-lib-${id}`).value;
    const ver = document.getElementById(`edit-ver-${id}`).value;
    const status = document.getElementById(`edit-status-${id}`).value;
    const dev = document.getElementById(`edit-dev-${id}`).value;
    const ket = document.getElementById(`edit-ket-${id}`).value;

    let promoteDate = null;
    if (dateStr) {
        promoteDate = parseDDMMYYYYToISO(dateStr);
        if (!promoteDate) {
            alert('Invalid date format. Please use dd/mm/yyyy');
            return;
        }
    }

    try {
        // Update Info including Name and Type
        const responseInfo = await fetch(`/api/update_mimix_obj_info/${id}`, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                obj: objName,
                obj_type: objType,
                promote_date: promoteDate,
                lib: lib,
                obj_ver: ver,
                developer: dev,
                keterangan: ket,
                mimix_status: status // Include status in info update as well if supported/needed
            })
        });

        if (!responseInfo.ok) {
            const res = await responseInfo.json();
            throw new Error(res.error || 'Failed to update info');
        }

        editingRowId = null;
        fetchObjects(searchInput.value);

    } catch (error) {
        console.error('Save error:', error);
        alert('Error saving object: ' + error.message);
    }
}


async function convertRequest(id) {
    if (!confirm('Are you sure you want to convert this request to an object?')) return;

    try {
        const response = await fetch(`/api/convert_obj_req/${id}`, {
            method: 'POST',
            headers: { 'Authorization': `Bearer ${token}` }
        });

        const result = await response.json();

        if (!response.ok) {
            throw new Error(result.error || 'Failed to convert request');
        }

        alert('Request converted successfully!');
        fetchRequests(reqSearchInput.value);
    } catch (error) {
        console.error('Convert error:', error);
        alert('Error: ' + error.message);
    }
}

async function deleteRequest(id) {
    if (!confirm('Are you sure you want to delete this request?')) return;

    try {
        const response = await fetch(`/api/delete_obj_req/${id}`, {
            method: 'DELETE',
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (!response.ok) throw new Error('Failed to delete request');
        fetchRequests(reqSearchInput.value);
    } catch (error) {
        console.error('Delete request error:', error);
        alert('Error deleting request');
    }
}

let editingReqId = null;
let originalReqRowHTML = '';

function editRequest(id) {
    if (editingReqId) cancelEditRequest(editingReqId);

    // For now, finding row by button context if id not on tr?
    const btn = document.querySelector(`button[onclick="editRequest('${id}')"]`);
    if (!btn) return;
    const row = btn.closest('tr'); // Safer

    editingReqId = id;
    originalReqRowHTML = row.innerHTML;

    const reqData = allRequests.find(r => r.id === id);
    if (!reqData) return;

    // Editable fields: ObjName, Lib, ObjVer, ObjType, PromoteDate, Developer, Statuses?
    // User didn't specify, but UpdateObjReqInfo allows: ObjName, Lib, ObjVer, ObjType, PromoteDate, Developer, PromoteStatus, ReqStatus.
    // Let's allow editing common fields.

    row.innerHTML = `
        <td><input class="edit-input" id="req-obj-${id}" value="${reqData.obj_name}"></td>
        <td>${reqData.requester}</td>
        <td>${reqData.updated_at ? new Date(reqData.updated_at).toLocaleString() : '-'}</td>
        <td><input class="edit-input" id="req-lib-${id}" value="${reqData.lib}"></td>
        <td><input class="edit-input" id="req-ver-${id}" value="${reqData.obj_ver}"></td>
        <td><input class="edit-input" id="req-type-${id}" value="${reqData.obj_type}"></td> 
        <td><input type="text" class="edit-input" id="req-promotedate-${id}" value="${formatDateToDDMMYYYY(reqData.promote_date)}" placeholder="dd/mm/yyyy"></td>
        <td><input class="edit-input" id="req-dev-${id}" value="${reqData.developer || ''}"></td>
        <td>
             <select class="edit-input" id="req-promotestatus-${id}">
                <option value="">-</option> 
                <option value="in_progress" ${reqData.promote_status === 'in_progress' ? 'selected' : ''}>In Progress</option>
                <option value="deployed" ${reqData.promote_status === 'deployed' ? 'selected' : ''}>Deployed</option>
            </select>
        </td>
        <td>
             <select class="edit-input" id="req-status-${id}">
                <option value="pending" ${reqData.req_status === 'pending' ? 'selected' : ''}>Pending</option>
                <option value="completed" ${reqData.req_status === 'completed' ? 'selected' : ''}>Completed</option>
            </select>
        </td>
        <td>
            <button class="action-btn save" onclick="saveRequest('${id}')" title="Save">${icons.save}</button>
            <button class="action-btn" onclick="cancelEditRequest('${id}')" title="Cancel">${icons.cancel}</button>
        </td>
    `;
}

function cancelEditRequest(id) {
    // Find row by editingReqId if id matches
    if (editingReqId === id) {
        // We need to find the row again. Since we replaced innerHTML, the button ref is gone.
        // But we can find it by ... wait we replaced the row content. 
        // We can find it by being the row that has the input with this ID?
        const input = document.getElementById(`req-obj-${id}`);
        if (input) {
            const row = input.closest('tr');
            row.innerHTML = originalReqRowHTML;
            editingReqId = null;
            originalReqRowHTML = '';
        }
    }
}

async function saveRequest(id) {
    const objName = document.getElementById(`req-obj-${id}`).value;
    const lib = document.getElementById(`req-lib-${id}`).value;
    const ver = document.getElementById(`req-ver-${id}`).value;
    const type = document.getElementById(`req-type-${id}`).value;
    const pDate = document.getElementById(`req-promotedate-${id}`).value;
    const dev = document.getElementById(`req-dev-${id}`).value;
    const pStatus = document.getElementById(`req-promotestatus-${id}`).value;
    const rStatus = document.getElementById(`req-status-${id}`).value;

    const data = {
        obj_name: objName,
        lib: lib,
        obj_ver: ver,
        obj_type: type,
        developer: dev,
        req_status: rStatus,
        promote_status: pStatus || ""
    };

    if (pDate) {
        const iso = parseDDMMYYYYToISO(pDate);
        if (iso) {
            data.promote_date = iso;
        } else {
            alert('Invalid Promte Date format. Use dd/mm/yyyy');
            return;
        }
    }

    try {
        const response = await fetch(`/api/update_obj_req_info/${id}`, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(data)
        });

        if (!response.ok) {
            const json = await response.json();
            throw new Error(json.error || 'Failed to update request');
        }

        editingReqId = null;
        fetchRequests();

    } catch (error) {
        console.error('Save request error:', error);
        alert('Error saving request: ' + error.message);
    }
}
