# Dev image: Vite dev server with the source mounted at /app (see docker-compose.yml).
FROM node:26-alpine
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
EXPOSE 5173
CMD ["sh", "-c", "npm ci --no-audit --no-fund && npx vite --host 0.0.0.0"]
