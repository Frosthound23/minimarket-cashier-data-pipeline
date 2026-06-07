# superset/Dockerfile
FROM apache/superset:latest

USER root

COPY requirements.txt /app/requirements.txt
RUN pip install --no-cache-dir -r /app/requirements.txt

USER superset