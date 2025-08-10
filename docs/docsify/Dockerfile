FROM node:18-alpine

# Install docsify-cli globally
RUN npm install -g docsify-cli@4.4.3

# Create app directory
WORKDIR /docs

# Copy documentation files
COPY . .

# Expose the port the app runs on
EXPOSE 3000

# Start docsify
CMD ["docsify", "serve", "-p", "3000", "."]
