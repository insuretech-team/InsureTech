import http.server
import socketserver
import os
import sys
from urllib.parse import urlparse

PORT = 8080
os.chdir(r'E:\Projects\InsureTech\api')

class CustomHandler(http.server.SimpleHTTPRequestHandler):
    def end_headers(self):
        # Add CORS headers for Swagger/ReDoc
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET, POST, OPTIONS')
        self.send_header('Access-Control-Allow-Headers', 'Content-Type')
        self.send_header('Cache-Control', 'no-cache, no-store, must-revalidate')
        super().end_headers()
    
    def do_GET(self):
        # Redirect root to docs/index.html
        if self.path == '/' or self.path == '':
            self.send_response(302)
            self.send_header('Location', '/docs/index.html')
            self.end_headers()
            return
        # Serve other files normally
        return http.server.SimpleHTTPRequestHandler.do_GET(self)

# Try to bind to port, retry with next port if occupied
max_attempts = 5
for attempt in range(max_attempts):
    try:
        with socketserver.TCPServer(('', PORT), CustomHandler) as httpd:
            print('')
            print('='*60)
            print('  InsureTech API Documentation Server')
            print('='*60)
            print('  Server running at: http://localhost:' + str(PORT) + '/')
            print('  Documentation:     http://localhost:' + str(PORT) + '/docs/')
            print('  Swagger UI:        http://localhost:' + str(PORT) + '/docs/swagger.html')
            print('  ReDoc:             http://localhost:' + str(PORT) + '/docs/redoc.html')
            print('  Schema Visualizer: http://localhost:' + str(PORT) + '/docs/index.html (🎨 tab)')
            print('  OpenAPI Spec:      http://localhost:' + str(PORT) + '/openapi.yaml')
            print('='*60)
            print('  Press Ctrl+C to stop the server')
            print('='*60)
            print('')
            httpd.serve_forever()
        break
    except OSError as e:
        if e.winerror == 10048:  # Port in use
            print('Port ' + str(PORT) + ' is in use, trying ' + str(PORT + 1) + '...')
            PORT += 1
        else:
            raise
else:
    print('Could not find available port after ' + str(max_attempts) + ' attempts')
    sys.exit(1)
