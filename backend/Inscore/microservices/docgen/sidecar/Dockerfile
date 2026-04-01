FROM python:3.12-slim-bookworm

WORKDIR /app

# Install WeasyPrint system dependencies (Pango, Cairo, GDK-PixBuf, fonts)
RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list.d/debian.sources 2>/dev/null || true && \
    apt-get update && apt-get install -y --no-install-recommends \
    libpango-1.0-0 \
    libpangoft2-1.0-0 \
    libpangocairo-1.0-0 \
    libcairo2 \
    libgdk-pixbuf-2.0-0 \
    libffi8 \
    shared-mime-info \
    fonts-liberation \
    fonts-dejavu-core \
    wget \
    && rm -rf /var/lib/apt/lists/*

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

EXPOSE 8500

HEALTHCHECK --interval=15s --timeout=5s --retries=5 --start-period=30s \
    CMD wget -qO- http://localhost:8500/health >/dev/null 2>&1 || exit 1

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8500", "--workers", "2"]
