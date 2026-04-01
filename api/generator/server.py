import http.server
import socketserver
import os
import sys
import threading

PORT = 8080
os.chdir(os.path.join(os.path.dirname(os.path.abspath(__file__)), '..'))

class CustomHandler(http.server.SimpleHTTPRequestHandler):
    def end_headers(self):
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET, POST, OPTIONS')
        self.send_header('Access-Control-Allow-Headers', 'Content-Type')
        self.send_header('Cache-Control', 'no-cache, no-store, must-revalidate')
        super().end_headers()

    def do_GET(self):
        if self.path == '/' or self.path == '':
            self.send_response(302)
            self.send_header('Location', '/docs/index.html')
            self.end_headers()
            return
        return http.server.SimpleHTTPRequestHandler.do_GET(self)

    def log_message(self, fmt, *args):
        pass

httpd = None
max_attempts = 5
for attempt in range(max_attempts):
    try:
        httpd = socketserver.TCPServer(('', PORT), CustomHandler)
        httpd.allow_reuse_address = True
        break
    except OSError:
        print(f'Port {PORT} in use, trying {PORT + 1}...')
        PORT += 1
else:
    print('Could not find available port after ' + str(max_attempts) + ' attempts')
    sys.exit(1)

print('')
print('='*60)
print('  InsureTech API Documentation Server')
print('='*60)
print('  Server running at: http://localhost:' + str(PORT) + '/')
print('  Swagger UI:        http://localhost:' + str(PORT) + '/docs/swagger.html')
print('  ReDoc:             http://localhost:' + str(PORT) + '/docs/redoc.html')
print('  Schema Visualizer: http://localhost:' + str(PORT) + '/docs/index.html')
print('  OpenAPI Spec:      http://localhost:' + str(PORT) + '/openapi.yaml')
print('='*60)
print('  Press Ctrl+C to stop')
print('='*60)
print('')

thread = threading.Thread(target=httpd.serve_forever, daemon=True)
thread.start()

try:
    thread.join()
except KeyboardInterrupt:
    pass
finally:
    print('\n  Shutting down documentation server...')
    httpd.shutdown()
    httpd.server_close()
    print('  Server stopped.')
    sys.exit(0)