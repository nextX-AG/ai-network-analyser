// Formular absenden
document.getElementById('config-form').addEventListener('submit', function(e) {
    e.preventDefault();
    
    const formData = {
        server_url: document.getElementById('server-url').value,
        name: document.getElementById('agent-name').value,
        interface: document.getElementById('interface').value,
        api_key: document.getElementById('api-key').value
    };
    
    fetch('/admin/config', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData)
    })
    .then(response => response.json())
    .then(data => {
        const successMsg = document.getElementById('success-message');
        const errorMsg = document.getElementById('error-message');
        
        if (data.success) {
            successMsg.style.display = 'block';
            errorMsg.style.display = 'none';
            setTimeout(() => {
                successMsg.style.display = 'none';
            }, 3000);
        } else {
            errorMsg.textContent = data.error || 'Fehler beim Speichern der Konfiguration';
            errorMsg.style.display = 'block';
            successMsg.style.display = 'none';
        }
    })
    .catch(err => {
        const errorMsg = document.getElementById('error-message');
        errorMsg.textContent = 'Fehler bei der Kommunikation mit dem Server: ' + err.message;
        errorMsg.style.display = 'block';
    });
});

// Neustart-Button
document.getElementById('restart-button').addEventListener('click', function() {
    if (confirm('Möchten Sie den Agent wirklich neustarten?')) {
        fetch('/admin/restart', {
            method: 'POST'
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                alert('Agent wird neu gestartet...');
            } else {
                alert('Fehler beim Neustarten: ' + (data.error || 'Unbekannter Fehler'));
            }
        })
        .catch(err => {
            alert('Fehler bei der Kommunikation mit dem Server: ' + err.message);
        });
    }
});

// Registrierungs-Button
document.getElementById('register-button').addEventListener('click', function() {
    fetch('/admin/register', {
        method: 'POST'
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            alert('Registrierung erfolgreich!');
            document.getElementById('server-connection').textContent = 'Verbunden';
        } else {
            alert('Fehler bei der Registrierung: ' + (data.error || 'Unbekannter Fehler'));
        }
    })
    .catch(err => {
        alert('Fehler bei der Kommunikation mit dem Server: ' + err.message);
    });
});

// Status-Updates
function updateStatus() {
    fetch('/admin/status')
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            document.getElementById('agent-status').textContent = data.data.status;
            document.getElementById('packets-captured').textContent = data.data.packets_captured;
            document.getElementById('active-interface').textContent = data.data.interface;
            document.getElementById('server-connection').textContent = 
                data.data.connected ? 'Verbunden' : 'Nicht verbunden';
        }
    })
    .catch(err => console.error('Fehler beim Abrufen des Status:', err));
}

// Status alle 5 Sekunden aktualisieren
setInterval(updateStatus, 5000); 