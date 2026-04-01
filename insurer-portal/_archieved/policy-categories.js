document.addEventListener('DOMContentLoaded', function() {
    // Form state tracking
    let formState = 'initial'; // 'initial', 'planAdded', 'subPlanAdded', 'subPlanLevel1'
    
    // DOM Elements
    const iconUploadArea = document.getElementById('iconUploadArea');
    const iconInput = document.getElementById('iconInput');
    const iconPreview = document.getElementById('iconPreview');
    const insuranceNameInput = document.getElementById('insuranceName');
    const planNameField = document.getElementById('planNameField');
    const subPlanNameField = document.getElementById('subPlanNameField');
    const mainActionBtn = document.getElementById('mainActionBtn');
    const rightColumn = document.getElementById('rightColumn');
    const addSubPlanLevel1Btn = document.getElementById('addSubPlanLevel1Btn');

    // Icon Upload
    iconUploadArea.addEventListener('click', () => {
        iconInput.click();
    });

    iconInput.addEventListener('change', (e) => {
        const file = e.target.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = (e) => {
                iconPreview.src = e.target.result;
                iconUploadArea.classList.add('has-image');
            };
            reader.readAsDataURL(file);
        }
    });

    // Drag and drop for icon
    iconUploadArea.addEventListener('dragover', (e) => {
        e.preventDefault();
        iconUploadArea.style.background = 'var(--primary-light)';
    });

    iconUploadArea.addEventListener('dragleave', (e) => {
        e.preventDefault();
        iconUploadArea.style.background = '';
    });

    iconUploadArea.addEventListener('drop', (e) => {
        e.preventDefault();
        iconUploadArea.style.background = '';
        
        const file = e.dataTransfer.files[0];
        if (file && (file.type === 'image/jpeg' || file.type === 'image/png')) {
            const reader = new FileReader();
            reader.onload = (e) => {
                iconPreview.src = e.target.result;
                iconUploadArea.classList.add('has-image');
            };
            reader.readAsDataURL(file);
        }
    });

    // Main Action Button Click
    mainActionBtn.addEventListener('click', () => {
        switch (formState) {
            case 'initial':
                // Validate Insurance Name
                if (!insuranceNameInput.value.trim()) {
                    insuranceNameInput.focus();
                    insuranceNameInput.style.borderColor = '#ef4444';
                    setTimeout(() => {
                        insuranceNameInput.style.borderColor = '';
                    }, 2000);
                    return;
                }
                
                // Show Plan Name field
                planNameField.classList.remove('hidden');
                mainActionBtn.textContent = 'Add Sub Plan';
                formState = 'planAdded';
                break;

            case 'planAdded':
                // Validate Plan Name
                const planNameInput = document.getElementById('planName');
                if (!planNameInput.value.trim()) {
                    planNameInput.focus();
                    planNameInput.style.borderColor = '#ef4444';
                    setTimeout(() => {
                        planNameInput.style.borderColor = '';
                    }, 2000);
                    return;
                }
                
                // Show Sub Plan Name field
                subPlanNameField.classList.remove('hidden');
                mainActionBtn.textContent = 'Add Sub Plan Level 1';
                formState = 'subPlanAdded';
                break;

            case 'subPlanAdded':
                // Validate Sub Plan Name
                const subPlanNameInput = document.getElementById('subPlanName');
                if (!subPlanNameInput.value.trim()) {
                    subPlanNameInput.focus();
                    subPlanNameInput.style.borderColor = '#ef4444';
                    setTimeout(() => {
                        subPlanNameInput.style.borderColor = '';
                    }, 2000);
                    return;
                }
                
                // Show right column
                rightColumn.classList.remove('hidden');
                formState = 'subPlanLevel1';
                break;

            case 'subPlanLevel1':
                // Submit the form / Add to table
                addPolicyToTable();
                resetForm();
                break;
        }
    });

    // Add Sub Plan Level 1 Button
    addSubPlanLevel1Btn.addEventListener('click', () => {
        const subPlanLevel1Input = document.getElementById('subPlanLevel1Name');
        if (!subPlanLevel1Input.value.trim()) {
            subPlanLevel1Input.focus();
            subPlanLevel1Input.style.borderColor = '#ef4444';
            setTimeout(() => {
                subPlanLevel1Input.style.borderColor = '';
            }, 2000);
            return;
        }
        
        // Add logic to handle sub plan level 1
        alert('Sub Plan Level 1 added: ' + subPlanLevel1Input.value);
        subPlanLevel1Input.value = '';
    });

    // Add Policy to Table
    function addPolicyToTable() {
        const table = document.querySelector('.data-table tbody');
        const rowCount = table.querySelectorAll('tr').length + 1;
        const insuranceName = insuranceNameInput.value;
        const planName = document.getElementById('planName').value;

        const newRow = document.createElement('tr');
        newRow.innerHTML = `
            <td>${String(rowCount).padStart(2, '0')}</td>
            <td>${planName}</td>
            <td>${insuranceName}</td>
            <td class="actions-cell">
                <button class="action-btn edit-btn">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/>
                        <path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"/>
                    </svg>
                </button>
                <button class="action-btn delete-btn" onclick="this.closest('tr').remove()">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <polyline points="3 6 5 6 21 6"/>
                        <path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/>
                        <line x1="10" y1="11" x2="10" y2="17"/>
                        <line x1="14" y1="11" x2="14" y2="17"/>
                    </svg>
                </button>
            </td>
        `;
        table.appendChild(newRow);
    }

    // Reset Form
    function resetForm() {
        formState = 'initial';
        insuranceNameInput.value = '';
        document.getElementById('planName').value = '';
        document.getElementById('subPlanName').value = '';
        document.getElementById('subPlanLevel1Name').value = '';
        planNameField.classList.add('hidden');
        subPlanNameField.classList.add('hidden');
        rightColumn.classList.add('hidden');
        mainActionBtn.textContent = 'Add Plan';
        iconUploadArea.classList.remove('has-image');
        iconPreview.src = '';
    }

    // Search functionality
    const searchInput = document.querySelector('.search-box input');
    searchInput.addEventListener('input', (e) => {
        const searchTerm = e.target.value.toLowerCase();
        const rows = document.querySelectorAll('.data-table tbody tr');
        
        rows.forEach(row => {
            const planName = row.cells[1].textContent.toLowerCase();
            const category = row.cells[2].textContent.toLowerCase();
            
            if (planName.includes(searchTerm) || category.includes(searchTerm)) {
                row.style.display = '';
            } else {
                row.style.display = 'none';
            }
        });
    });

    // Delete row functionality
    document.querySelectorAll('.delete-btn').forEach(btn => {
        btn.addEventListener('click', function() {
            if (confirm('Are you sure you want to delete this policy?')) {
                this.closest('tr').remove();
                // Renumber rows
                const rows = document.querySelectorAll('.data-table tbody tr');
                rows.forEach((row, index) => {
                    row.cells[0].textContent = String(index + 1).padStart(2, '0');
                });
            }
        });
    });

    // Show list select
    const showListSelect = document.getElementById('showListSelect');
    showListSelect.addEventListener('change', (e) => {
        // In a real app, this would paginate the table
        console.log('Show', e.target.value, 'items');
    });

    // Add Department button
    const addDepartmentBtn = document.getElementById('addDepartmentBtn');
    addDepartmentBtn.addEventListener('click', () => {
        // Scroll to add policy section
        document.querySelector('.add-policy-section').scrollIntoView({ behavior: 'smooth' });
        insuranceNameInput.focus();
    });
});
