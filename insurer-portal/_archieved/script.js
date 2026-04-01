// Claim Ratio Donut Chart
const claimRatioCtx = document.getElementById('claimRatioChart').getContext('2d');
new Chart(claimRatioCtx, {
    type: 'doughnut',
    data: {
        labels: ['Settled', 'Approved', 'Under Review', 'Document Requested', 'Rejected'],
        datasets: [{
            data: [47.6, 24.3, 10.2, 7.9, 7.2],
            backgroundColor: [
                '#1f2937',
                '#16a34a',
                '#3b82f6',
                '#f97316',
                '#ef4444'
            ],
            borderWidth: 0,
            cutout: '70%'
        }]
    },
    options: {
        responsive: true,
        maintainAspectRatio: true,
        plugins: {
            legend: {
                display: false
            },
            tooltip: {
                callbacks: {
                    label: function(context) {
                        return context.label + ': ' + context.raw + '%';
                    }
                }
            }
        }
    }
});

// Premium Summary Line Chart
const premiumCtx = document.getElementById('premiumChart').getContext('2d');
new Chart(premiumCtx, {
    type: 'line',
    data: {
        labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov'],
        datasets: [{
            label: 'Premium',
            data: [45000, 52000, 48000, 61000, 55000, 67000, 73997, 71000, 68000, 65000, 62000],
            borderColor: '#16a34a',
            backgroundColor: 'transparent',
            borderWidth: 2,
            tension: 0.4,
            pointBackgroundColor: '#16a34a',
            pointBorderColor: '#ffffff',
            pointBorderWidth: 2,
            pointRadius: 5,
            pointHoverRadius: 7
        }]
    },
    options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: {
                display: false
            },
            tooltip: {
                backgroundColor: '#ffffff',
                titleColor: '#6b7280',
                bodyColor: '#16a34a',
                bodyFont: {
                    size: 16,
                    weight: 'bold'
                },
                borderColor: '#e5e7eb',
                borderWidth: 1,
                padding: 12,
                displayColors: false,
                callbacks: {
                    title: function(context) {
                        return context[0].label + ' Premium';
                    },
                    label: function(context) {
                        return '৳ ' + context.raw.toLocaleString();
                    }
                }
            }
        },
        scales: {
            x: {
                grid: {
                    display: false
                },
                ticks: {
                    color: '#9ca3af',
                    font: {
                        size: 12
                    }
                }
            },
            y: {
                min: 0,
                max: 150000,
                ticks: {
                    stepSize: 35000,
                    color: '#9ca3af',
                    font: {
                        size: 12
                    },
                    callback: function(value) {
                        return value.toLocaleString();
                    }
                },
                grid: {
                    color: '#f3f4f6'
                }
            }
        }
    }
});

// Settlement Status Bar Chart
const settlementCtx = document.getElementById('settlementChart').getContext('2d');
new Chart(settlementCtx, {
    type: 'bar',
    data: {
        labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May'],
        datasets: [
            {
                label: 'Settled',
                data: [1200, 1400, 1100, 900, 800],
                backgroundColor: '#1f2937',
                borderRadius: 4,
                barPercentage: 0.7,
                categoryPercentage: 0.8
            },
            {
                label: 'Approved',
                data: [1500, 1300, 1400, 1200, 1100],
                backgroundColor: '#16a34a',
                borderRadius: 4,
                barPercentage: 0.7,
                categoryPercentage: 0.8
            },
            {
                label: 'Under Review',
                data: [800, 900, 700, 600, 500],
                backgroundColor: '#3b82f6',
                borderRadius: 4,
                barPercentage: 0.7,
                categoryPercentage: 0.8
            },
            {
                label: 'Document Requested',
                data: [600, 500, 550, 450, 400],
                backgroundColor: '#f97316',
                borderRadius: 4,
                barPercentage: 0.7,
                categoryPercentage: 0.8
            },
            {
                label: 'Rejected',
                data: [400, 350, 300, 250, 200],
                backgroundColor: '#ef4444',
                borderRadius: 4,
                barPercentage: 0.7,
                categoryPercentage: 0.8
            }
        ]
    },
    options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: {
                display: false
            },
            tooltip: {
                backgroundColor: '#ffffff',
                titleColor: '#1f2937',
                bodyColor: '#6b7280',
                borderColor: '#e5e7eb',
                borderWidth: 1,
                padding: 12
            }
        },
        scales: {
            x: {
                grid: {
                    display: false
                },
                ticks: {
                    color: '#9ca3af',
                    font: {
                        size: 12
                    }
                }
            },
            y: {
                min: 0,
                max: 1500,
                ticks: {
                    stepSize: 500,
                    color: '#9ca3af',
                    font: {
                        size: 12
                    },
                    callback: function(value) {
                        return value.toLocaleString();
                    }
                },
                grid: {
                    color: '#f3f4f6'
                }
            }
        }
    }
});

// Navigation active state
document.querySelectorAll('.nav-item').forEach(item => {
    item.addEventListener('click', function(e) {
        e.preventDefault();
        document.querySelectorAll('.nav-item').forEach(nav => nav.classList.remove('active'));
        this.classList.add('active');
    });
});

// Dropdown functionality
document.querySelectorAll('.dropdown-btn').forEach(btn => {
    btn.addEventListener('click', function() {
        // Toggle dropdown menu (placeholder for future implementation)
        console.log('Dropdown clicked');
    });
});
