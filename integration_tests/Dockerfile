FROM python:3.11-slim

WORKDIR /tests

# Copy only requirements first for better layer caching
COPY requirements.txt .

# Install dependencies (cached as separate layer)
RUN pip install -r requirements.txt

# Copy test files
COPY . .
