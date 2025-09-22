FROM node:20-bullseye

WORKDIR /app

COPY package*.json ./

#Build arg sa default npm install --legacy-peer-deps
ARG NPM_INSTALL_CMD="npm install --legacy-peer-deps"
RUN $NPM_INSTALL_CMD

COPY . .

EXPOSE 4200

CMD ["npx", "ng", "serve", "--host", "0.0.0.0", "--poll=2000"]
