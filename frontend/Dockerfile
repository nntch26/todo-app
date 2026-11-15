# ---------- Stage 1: install dependencies ----------
FROM node:22-alpine AS deps
WORKDIR /app

# Install from the lockfile only: this layer is cached until package*.json change
COPY package.json package-lock.json ./
RUN npm ci

# ---------- Stage 2: build ----------
FROM node:22-alpine AS build
WORKDIR /app
ENV NEXT_TELEMETRY_DISABLED=1

# API_URL is read by next.config.ts at BUILD time and baked into the output,
# so it must be passed as a build arg (setting it when the container runs has no effect)
ARG API_URL
RUN test -n "$API_URL" || (echo "build arg API_URL is required" && exit 1)
ENV API_URL=$API_URL

COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build

# ---------- Stage 3: run (only the files the server needs) ----------
FROM node:22-alpine AS run
WORKDIR /app
ENV NODE_ENV=production \
    NEXT_TELEMETRY_DISABLED=1 \
    PORT=3000 \
    HOSTNAME=0.0.0.0

# standalone = server.js + the traced node_modules; static files must be copied next to it
COPY --from=build --chown=node:node /app/.next/standalone ./
COPY --from=build --chown=node:node /app/.next/static ./.next/static

# "node" is a non-root user that comes with the node image
USER node

EXPOSE 3000
CMD ["node", "server.js"]
