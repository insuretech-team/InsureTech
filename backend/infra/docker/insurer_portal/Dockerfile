FROM node:22-alpine AS deps

WORKDIR /
RUN mkdir -p /app /sdks/insuretech-typescript-sdk

WORKDIR /app
COPY insurer-portal/package.json ./
COPY sdks/insuretech-typescript-sdk/lifeplus-insuretech-sdk-0.1.0.tgz /sdks/insuretech-typescript-sdk/

RUN npm install

FROM node:22-alpine AS builder

WORKDIR /
RUN mkdir -p /app /sdks/insuretech-typescript-sdk

WORKDIR /app
COPY insurer-portal/ ./
COPY --from=deps /app/node_modules ./node_modules
COPY sdks/insuretech-typescript-sdk/lifeplus-insuretech-sdk-0.1.0.tgz /sdks/insuretech-typescript-sdk/

ARG NEXT_PUBLIC_INSURETECH_API_BASE_URL=http://localhost:8080
ENV NEXT_PUBLIC_INSURETECH_API_BASE_URL=$NEXT_PUBLIC_INSURETECH_API_BASE_URL

RUN NODE_OPTIONS="--max-old-space-size=4096" npm run build

FROM node:22-alpine AS runner

WORKDIR /app
ENV NODE_ENV=production
ENV PORT=3000

COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

EXPOSE 3000

CMD ["node", "server.js"]
