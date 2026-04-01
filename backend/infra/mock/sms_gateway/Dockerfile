FROM python:3.12-alpine

WORKDIR /app

RUN pip install --no-cache-dir flask

COPY mock_sms.py .

EXPOSE 8600

HEALTHCHECK --interval=5s --timeout=3s --retries=5 \
    CMD python -c "import urllib.request; urllib.request.urlopen('http://127.0.0.1:8600/health')" >/dev/null 2>&1 || exit 1

CMD ["python", "mock_sms.py"]
