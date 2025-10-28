FROM python:3.11-slim

WORKDIR /tests

# Copy only requirements first for better layer caching
COPY requirements.txt .

# Install dependencies in a single layer
RUN pip install --no-cache-dir -r requirements.txt

# Copy test files
COPY . .

# Set environment variable for base URL
ENV BASE_URL=http://server:8080

# Run tests by default
CMD ["pytest", "-v", "--tb=short", "-n", "auto"]
