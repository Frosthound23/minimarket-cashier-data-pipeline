# superset/Dockerfile
from apache/superset:latest

user root

copy requirements.txt /app/requirements.txt
run pip install --no-cache-dir -r /app/requirements.txt

user superset