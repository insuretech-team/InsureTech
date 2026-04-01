document.addEventListener('DOMContentLoaded', function() {
    // Tab switching
    const tabs = document.querySelectorAll('.tab');
    const tabPanels = document.querySelectorAll('.tab-panel');

    tabs.forEach(tab => {
        tab.addEventListener('click', () => {
            tabs.forEach(t => t.classList.remove('active'));
            tabPanels.forEach(p => p.classList.remove('active'));
            
            tab.classList.add('active');
            const panelId = tab.dataset.tab + '-panel';
            document.getElementById(panelId).classList.add('active');
        });
    });

    // Custom Select Dropdown functionality
    initCustomSelects();

    function initCustomSelects() {
        const customSelects = document.querySelectorAll('.custom-select-wrapper');
        
        customSelects.forEach(select => {
            const trigger = select.querySelector('.custom-select-trigger');
            const options = select.querySelectorAll('.custom-select-option');
            
            trigger.addEventListener('click', (e) => {
                e.stopPropagation();
                // Close all other dropdowns
                document.querySelectorAll('.custom-select-wrapper.open').forEach(s => {
                    if (s !== select) s.classList.remove('open');
                });
                select.classList.toggle('open');
            });
            
            options.forEach(option => {
                option.addEventListener('click', () => {
                    const value = option.dataset.value;
                    const text = option.querySelector('span:last-child').textContent;
                    const iconHtml = option.querySelector('.option-icon').innerHTML;
                    
                    trigger.querySelector('.select-icon').innerHTML = iconHtml;
                    trigger.querySelector('.select-text').textContent = text;
                    
                    select.classList.remove('open');
                    
                    // Update the options list based on selection
                    const questionCard = select.closest('.question-card');
                    if (questionCard) {
                        updateQuestionType(questionCard, value);
                    }
                });
            });
        });
        
        // Close dropdowns when clicking outside
        document.addEventListener('click', () => {
            document.querySelectorAll('.custom-select-wrapper.open').forEach(s => {
                s.classList.remove('open');
            });
        });
    }

    function updateQuestionType(questionCard, type) {
        const optionsList = questionCard.querySelector('.options-list');
        const addOptionBtn = questionCard.querySelector('.add-option-btn');
        
        if (type === 'date') {
            // Hide options for date type
            if (optionsList) optionsList.style.display = 'none';
            if (addOptionBtn) addOptionBtn.style.display = 'none';
        } else {
            // Show options for multiple choice and combo box
            if (optionsList) optionsList.style.display = 'flex';
            if (addOptionBtn) addOptionBtn.style.display = 'flex';
            
            // Update radio circles to checkboxes for combo box
            const radioCircles = optionsList.querySelectorAll('.radio-circle');
            radioCircles.forEach(circle => {
                if (type === 'combo-box') {
                    circle.style.borderRadius = '4px';
                } else {
                    circle.style.borderRadius = '50%';
                }
            });
        }
    }

    // Submenu toggle
    const planTemplateToggle = document.getElementById('planTemplateToggle');
    const planTemplateSubmenu = document.getElementById('planTemplateSubmenu');

    planTemplateToggle.addEventListener('click', (e) => {
        e.preventDefault();
        planTemplateSubmenu.classList.toggle('show');
    });

    // Submenu item click
    const submenuItems = document.querySelectorAll('.submenu-item');
    const coverageType = document.getElementById('coverageType');

    submenuItems.forEach(item => {
        item.addEventListener('click', (e) => {
            e.preventDefault();
            submenuItems.forEach(i => i.classList.remove('active'));
            item.classList.add('active');
            
            const type = item.dataset.type;
            if (coverageType) {
                coverageType.textContent = type.charAt(0).toUpperCase() + type.slice(1);
            }
        });
    });

    // Drag and Drop functionality
    const draggables = document.querySelectorAll('.element-item');
    const dropZone = document.getElementById('planInfoDropZone');

    draggables.forEach(draggable => {
        draggable.addEventListener('dragstart', (e) => {
            draggable.classList.add('dragging');
            e.dataTransfer.setData('text/plain', draggable.dataset.element);
        });

        draggable.addEventListener('dragend', () => {
            draggable.classList.remove('dragging');
        });
    });

    if (dropZone) {
        dropZone.addEventListener('dragover', (e) => {
            e.preventDefault();
            dropZone.classList.add('drag-over');
        });

        dropZone.addEventListener('dragleave', () => {
            dropZone.classList.remove('drag-over');
        });

        dropZone.addEventListener('drop', (e) => {
            e.preventDefault();
            dropZone.classList.remove('drag-over');
            
            const elementType = e.dataTransfer.getData('text/plain');
            createDroppedElement(elementType, dropZone);
        });
    }

    function createDroppedElement(type, container) {
        const droppedItem = document.createElement('div');
        droppedItem.className = 'dropped-item';
        
        const removeBtn = document.createElement('button');
        removeBtn.className = 'remove-btn';
        removeBtn.innerHTML = `
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
        `;
        removeBtn.addEventListener('click', () => droppedItem.remove());
        droppedItem.appendChild(removeBtn);

        switch(type) {
            case 'plan-card':
                droppedItem.innerHTML += `
                    <div class="dropped-plan-card">
                        <input type="text" placeholder="Enter plan name">
                        <div class="dropped-plan-info">
                            <div class="dropped-coverage">
                                <label>Coverage up to</label>
                                <span class="amount">৳ 65,000</span>
                            </div>
                            <div class="dropped-meta">
                                <div class="meta-item">
                                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#22c55e" stroke-width="2">
                                        <path d="M12 2v4M12 18v4M4.93 4.93l2.83 2.83M16.24 16.24l2.83 2.83M2 12h4M18 12h4M4.93 19.07l2.83-2.83M16.24 7.76l2.83-2.83"/>
                                    </svg>
                                    <span>Premium price:</span>
                                    <span class="value green">৳ 1,200</span>
                                </div>
                                <div class="meta-item">
                                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#22c55e" stroke-width="2">
                                        <circle cx="12" cy="12" r="10"/>
                                        <polyline points="12,6 12,12 16,14"/>
                                    </svg>
                                    <span>Policy duration:</span>
                                    <span class="value green">1 year</span>
                                </div>
                            </div>
                        </div>
                    </div>
                `;
                break;
            case 'header':
                droppedItem.innerHTML += `
                    <div class="dropped-header">
                        <input type="text" placeholder="Enter header text" style="font-weight: 600; font-size: 18px;">
                        <input type="text" placeholder="Enter paragraph text" style="color: #6b7280;">
                    </div>
                `;
                break;
            case 'description':
            case 'bullet':
                droppedItem.innerHTML += `
                    <div class="dropped-description">
                        <div class="icon-badge ${type === 'bullet' ? 'orange' : ''}">
                            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="${type === 'bullet' ? '#f97316' : '#22c55e'}" stroke-width="2">
                                <rect x="3" y="3" width="18" height="18" rx="2"/>
                                <path d="M9 12h6${type === 'description' ? 'M12 9v6' : ''}"/>
                            </svg>
                        </div>
                        <input type="text" placeholder="Enter label" value="${type === 'description' ? 'IPD Coverage' : 'Cabin room rent'}">
                    </div>
                `;
                break;
            case 'ipd-coverage':
                droppedItem.innerHTML += `
                    <div class="dropped-ipd">
                        <input type="text" placeholder="Enter header" style="font-weight: 600;">
                        <div class="point-item">
                            <span class="point-icon red">
                                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                                    <line x1="18" y1="6" x2="6" y2="18"/>
                                    <line x1="6" y1="6" x2="18" y2="18"/>
                                </svg>
                            </span>
                            <input type="text" placeholder="Enter point">
                        </div>
                        <button class="add-point-btn" onclick="addPoint(this)">+ Add Point</button>
                    </div>
                `;
                break;
        }

        container.appendChild(droppedItem);
    }

    // Add point functionality
    window.addPoint = function(btn) {
        const container = btn.parentElement;
        const newPoint = document.createElement('div');
        newPoint.className = 'point-item';
        newPoint.innerHTML = `
            <span class="point-icon red">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                    <line x1="18" y1="6" x2="6" y2="18"/>
                    <line x1="6" y1="6" x2="18" y2="18"/>
                </svg>
            </span>
            <input type="text" placeholder="Enter point">
        `;
        container.insertBefore(newPoint, btn);
    };

    // Add option functionality for questionnaire
    const addOptionBtns = document.querySelectorAll('.add-option-btn');
    addOptionBtns.forEach(btn => {
        btn.addEventListener('click', () => {
            const optionsList = btn.previousElementSibling;
            const optionCount = optionsList.children.length + 1;
            const newOption = document.createElement('div');
            newOption.className = 'option-item';
            newOption.innerHTML = `
                <div class="radio-circle"></div>
                <input type="text" placeholder="Option ${optionCount}">
            `;
            optionsList.appendChild(newOption);
        });
    });

    // Add question functionality
    const addQuestionBtn = document.querySelector('.add-question-btn');
    if (addQuestionBtn) {
        addQuestionBtn.addEventListener('click', () => {
            const questionsContainer = document.querySelector('.questionnaire-builder');
            const newQuestion = document.createElement('div');
            newQuestion.className = 'question-card';
            newQuestion.innerHTML = `
                <div class="question-header">
                    <input type="text" class="question-input" placeholder="Add Question">
                    <div class="question-type">
                        <label>Question Type</label>
                        <div class="custom-select-wrapper">
                            <button class="custom-select-trigger" type="button">
                                <span class="select-icon multiple-choice-icon">
                                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#22c55e" stroke-width="2">
                                        <circle cx="12" cy="12" r="10"/>
                                        <circle cx="12" cy="12" r="4" fill="#22c55e"/>
                                    </svg>
                                </span>
                                <span class="select-text">Multiple Choice</span>
                                <svg class="chevron" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <polyline points="6 9 12 15 18 9"/>
                                </svg>
                            </button>
                            <div class="custom-select-dropdown">
                                <div class="custom-select-option" data-value="multiple-choice">
                                    <span class="option-icon">
                                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#22c55e" stroke-width="2">
                                            <circle cx="12" cy="12" r="10"/>
                                            <circle cx="12" cy="12" r="4" fill="#22c55e"/>
                                        </svg>
                                    </span>
                                    <span>Multiple Choice</span>
                                </div>
                                <div class="custom-select-option" data-value="combo-box">
                                    <span class="option-icon">
                                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#22c55e" stroke-width="2.5">
                                            <polyline points="20 6 9 17 4 12"/>
                                        </svg>
                                    </span>
                                    <span>Combo Box</span>
                                </div>
                                <div class="custom-select-option" data-value="date">
                                    <span class="option-icon">
                                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#22c55e" stroke-width="2">
                                            <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
                                            <line x1="16" y1="2" x2="16" y2="6"/>
                                            <line x1="8" y1="2" x2="8" y2="6"/>
                                            <line x1="3" y1="10" x2="21" y2="10"/>
                                        </svg>
                                    </span>
                                    <span>Date</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="options-list">
                    <div class="option-item">
                        <div class="radio-circle"></div>
                        <input type="text" placeholder="Option 1">
                    </div>
                </div>
                <button class="add-option-btn">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <circle cx="12" cy="12" r="10"/>
                        <line x1="12" y1="8" x2="12" y2="16"/>
                        <line x1="8" y1="12" x2="16" y2="12"/>
                    </svg>
                    Add Options
                </button>
                <div class="question-actions">
                    <div class="toggle-group">
                        <span>Required</span>
                        <label class="toggle">
                            <input type="checkbox">
                            <span class="toggle-slider"></span>
                        </label>
                    </div>
                    <div class="delete-group">
                        <span>Delete Question</span>
                        <button class="delete-btn" onclick="this.closest('.question-card').remove()">
                            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <polyline points="3 6 5 6 21 6"/>
                                <path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/>
                            </svg>
                        </button>
                    </div>
                </div>
            `;
            questionsContainer.insertBefore(newQuestion, addQuestionBtn);

            // Re-initialize custom selects for the new question
            initCustomSelectsForElement(newQuestion);

            // Add event listener to new add option button
            const newAddOptionBtn = newQuestion.querySelector('.add-option-btn');
            newAddOptionBtn.addEventListener('click', () => {
                const optionsList = newAddOptionBtn.previousElementSibling;
                const optionCount = optionsList.children.length + 1;
                const newOption = document.createElement('div');
                newOption.className = 'option-item';
                newOption.innerHTML = `
                    <div class="radio-circle"></div>
                    <input type="text" placeholder="Option ${optionCount}">
                `;
                optionsList.appendChild(newOption);
            });
        });
    }

    // Initialize custom selects for dynamically added elements
    function initCustomSelectsForElement(element) {
        const customSelects = element.querySelectorAll('.custom-select-wrapper');
        
        customSelects.forEach(select => {
            const trigger = select.querySelector('.custom-select-trigger');
            const options = select.querySelectorAll('.custom-select-option');
            
            trigger.addEventListener('click', (e) => {
                e.stopPropagation();
                document.querySelectorAll('.custom-select-wrapper.open').forEach(s => {
                    if (s !== select) s.classList.remove('open');
                });
                select.classList.toggle('open');
            });
            
            options.forEach(option => {
                option.addEventListener('click', () => {
                    const value = option.dataset.value;
                    const text = option.querySelector('span:last-child').textContent;
                    const iconHtml = option.querySelector('.option-icon').innerHTML;
                    
                    trigger.querySelector('.select-icon').innerHTML = iconHtml;
                    trigger.querySelector('.select-text').textContent = text;
                    
                    select.classList.remove('open');
                    
                    const questionCard = select.closest('.question-card');
                    if (questionCard) {
                        updateQuestionType(questionCard, value);
                    }
                });
            });
        });
    }

    // Add document field functionality
    const addDocumentBtn = document.querySelector('.add-document-btn');
    if (addDocumentBtn) {
        addDocumentBtn.addEventListener('click', () => {
            const documentsContainer = document.querySelector('.documents-builder');
            const newDocument = document.createElement('div');
            newDocument.className = 'document-upload-card';
            newDocument.innerHTML = `
                <div class="upload-icon">
                    <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="#9ca3af" stroke-width="1.5">
                        <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/>
                        <polyline points="17 8 12 3 7 8"/>
                        <line x1="12" y1="3" x2="12" y2="15"/>
                    </svg>
                </div>
                <input type="text" class="document-name-input" placeholder="Document Name">
                <p class="file-hint">jpg, png file only</p>
            `;
            documentsContainer.insertBefore(newDocument, addDocumentBtn);
        });
    }

    // Search functionality
    const searchInput = document.getElementById('elementSearch');
    const elementsList = document.getElementById('elementsList');

    if (searchInput && elementsList) {
        searchInput.addEventListener('input', (e) => {
            const searchTerm = e.target.value.toLowerCase();
            const items = elementsList.querySelectorAll('.element-item, .element-label');
            
            items.forEach(item => {
                const text = item.textContent.toLowerCase();
                if (text.includes(searchTerm) || searchTerm === '') {
                    item.style.display = '';
                } else {
                    item.style.display = 'none';
                }
            });
        });
    }

    // Delete question buttons
    document.querySelectorAll('.delete-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            btn.closest('.question-card').remove();
        });
    });
});
