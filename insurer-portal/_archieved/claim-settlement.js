document.addEventListener('DOMContentLoaded', function() {
    // Claim Tabs
    const claimTabs = document.querySelectorAll('.claim-tab');
    claimTabs.forEach(tab => {
        tab.addEventListener('click', () => {
            claimTabs.forEach(t => t.classList.remove('active'));
            tab.classList.add('active');
        });
    });

    // Modal functionality
    const modal = document.getElementById('claimModal');
    const closeModalBtn = document.getElementById('closeModal');
    const viewBtns = document.querySelectorAll('.view-btn');

    // Open modal when clicking view button
    viewBtns.forEach(btn => {
        btn.addEventListener('click', () => {
            modal.classList.add('active');
            document.body.style.overflow = 'hidden';
        });
    });

    // Close modal
    closeModalBtn.addEventListener('click', () => {
        modal.classList.remove('active');
        document.body.style.overflow = '';
    });

    // Close modal when clicking outside
    modal.addEventListener('click', (e) => {
        if (e.target === modal) {
            modal.classList.remove('active');
            document.body.style.overflow = '';
        }
    });

    // Modal tabs
    const modalTabs = document.querySelectorAll('.modal-tab');
    const modalPanels = document.querySelectorAll('.modal-panel');

    modalTabs.forEach(tab => {
        tab.addEventListener('click', () => {
            const tabId = tab.dataset.tab;
            
            modalTabs.forEach(t => t.classList.remove('active'));
            modalPanels.forEach(p => p.classList.remove('active'));
            
            tab.classList.add('active');
            document.getElementById(`${tabId}-panel`).classList.add('active');
        });
    });

    // Status dropdown
    const statusDropdown = document.querySelector('.status-dropdown');
    const statusTrigger = document.querySelector('.status-trigger');
    const statusOptions = document.querySelectorAll('.status-option');

    statusTrigger.addEventListener('click', (e) => {
        e.stopPropagation();
        statusDropdown.classList.toggle('open');
    });

    statusOptions.forEach(option => {
        option.addEventListener('click', () => {
            const value = option.dataset.value;
            const text = option.textContent;
            
            statusTrigger.querySelector('span').textContent = text;
            statusDropdown.classList.remove('open');
        });
    });

    // Close dropdown when clicking outside
    document.addEventListener('click', () => {
        statusDropdown.classList.remove('open');
    });

    // Table sorting (basic implementation)
    const sortableHeaders = document.querySelectorAll('.claims-table th svg');
    sortableHeaders.forEach(header => {
        header.style.cursor = 'pointer';
        header.addEventListener('click', () => {
            // Add sorting logic here if needed
            console.log('Sort by:', header.closest('th').textContent.trim());
        });
    });

    // Search functionality
    const searchInput = document.querySelector('.table-search input');
    const tableRows = document.querySelectorAll('.claims-table tbody tr');

    searchInput.addEventListener('input', (e) => {
        const searchTerm = e.target.value.toLowerCase();
        
        tableRows.forEach(row => {
            const text = row.textContent.toLowerCase();
            row.style.display = text.includes(searchTerm) ? '' : 'none';
        });
    });

    // Pagination
    const pageBtns = document.querySelectorAll('.page-btn:not(.prev):not(.next)');
    pageBtns.forEach(btn => {
        btn.addEventListener('click', () => {
            pageBtns.forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
        });
    });

    // Show list dropdown
    const showListSelect = document.querySelector('.show-list select');
    showListSelect.addEventListener('change', (e) => {
        const count = parseInt(e.target.value);
        tableRows.forEach((row, index) => {
            row.style.display = index < count ? '' : 'none';
        });
    });

    // Edit button functionality
    const editBtns = document.querySelectorAll('.edit-btn');
    editBtns.forEach(btn => {
        btn.addEventListener('click', () => {
            const row = btn.closest('tr');
            const claimId = row.dataset.claimId;
            console.log('Edit claim:', claimId);
            // Add edit logic here
        });
    });

    // Stats cards click handler (for navigation or filtering)
    const statCards = document.querySelectorAll('.stat-card');
    statCards.forEach((card, index) => {
        card.style.cursor = 'pointer';
        card.addEventListener('click', () => {
            const labels = ['total', 'under-review', 'approved', 'requested', 'settled', 'rejected'];
            const filter = labels[index];
            
            if (filter === 'total') {
                tableRows.forEach(row => row.style.display = '');
            } else {
                tableRows.forEach(row => {
                    const status = row.querySelector('.status');
                    if (status) {
                        const statusClass = status.classList[1];
                        row.style.display = statusClass === filter ? '' : 'none';
                    }
                });
            }
        });
    });

    // Filter button
    const filterBtn = document.querySelector('.btn-filter');
    filterBtn.addEventListener('click', () => {
        // Show all rows (reset filter)
        tableRows.forEach(row => row.style.display = '');
    });

    // Export button dropdown (simplified)
    const exportBtn = document.querySelector('.btn-text');
    exportBtn.addEventListener('click', () => {
        console.log('Export options: Excel, PDF, CSV');
    });

    // Add Claim button
    const addClaimBtn = document.querySelector('.btn-primary');
    addClaimBtn.addEventListener('click', () => {
        console.log('Add new claim');
        // Add new claim form logic here
    });

    // Upload button
    const uploadBtn = document.querySelector('.btn-outline');
    uploadBtn.addEventListener('click', () => {
        const input = document.createElement('input');
        input.type = 'file';
        input.accept = '.xlsx,.xls,.csv';
        input.click();
        input.addEventListener('change', (e) => {
            const file = e.target.files[0];
            if (file) {
                console.log('File selected:', file.name);
            }
        });
    });
});
