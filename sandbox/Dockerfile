# Sandbox runner page, built and served as static files with its CSP (build context: repo root).
FROM node:26-alpine
WORKDIR /repo/sandbox
COPY sandbox/package.json sandbox/package-lock.json ./
RUN npm ci --no-audit --no-fund
COPY sandbox/ ./
COPY runners/python/harness/cq/__init__.py /repo/runners/python/harness/cq/__init__.py
ARG VITE_APP_ORIGINS=http://localhost:5173
ARG SANDBOX_ORIGIN=http://localhost:5174
ENV VITE_APP_ORIGINS=$VITE_APP_ORIGINS SANDBOX_ORIGIN=$SANDBOX_ORIGIN
# vite preview writes its compiled config here at start-up, as the node user.
RUN npm run build && mkdir -p node_modules/.vite-temp && chown -R node:node node_modules/.vite-temp
USER node
EXPOSE 5174
CMD ["npm", "run", "serve"]
