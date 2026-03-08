# Dockerfile for b2b_portal
FROM node:22-alpine AS deps

# Create root directories to mimic monorepo structure
WORKDIR /
RUN mkdir -p /app /sdks/insuretech-typescript-sdk

WORKDIR /app
COPY b2b_portal/package.json ./

# Copy the local SDK exactly where package.json expects it (../sdks/...)
COPY sdks/insuretech-typescript-sdk/lifeplus-insuretech-sdk-0.1.0.tgz /sdks/insuretech-typescript-sdk/

# Now npm install will resolve ../sdks correctly from /app
RUN npm install

FROM node:22-alpine AS builder
# Recreate the structure in the builder stage
WORKDIR /
RUN mkdir -p /app /sdks/insuretech-typescript-sdk

WORKDIR /app
COPY b2b_portal/ ./
COPY --from=deps /app/node_modules ./node_modules
COPY sdks/insuretech-typescript-sdk/lifeplus-insuretech-sdk-0.1.0.tgz /sdks/insuretech-typescript-sdk/

# NEXT_PUBLIC_* vars are baked into the JS bundle at build time by Next.js.
# Accept the gateway URL as a build argument so it is correctly embedded.
# Usage: docker build --build-arg NEXT_PUBLIC_INSURETECH_API_BASE_URL=https://api.yourdomain.com ...
ARG NEXT_PUBLIC_INSURETECH_API_BASE_URL=http://localhost:8080
ENV NEXT_PUBLIC_INSURETECH_API_BASE_URL=$NEXT_PUBLIC_INSURETECH_API_BASE_URL

# INSURETECH_API_BASE_URL (non-public) is read at runtime by Node.js server.
# It is set via docker-compose environment / runtime env — not needed at build time.
RUN npm run build

FROM node:22-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

EXPOSE 3000
CMD ["node", "server.js"]
