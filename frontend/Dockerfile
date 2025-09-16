FROM node:20-bullseye

# Setujemo radni folder
WORKDIR /app

# Kopiramo samo package fajlove za caching
COPY package*.json ./

# Instaliramo dependencies
RUN npm install

# Kopiramo ostatak projekta
COPY . .

# Expose Angular dev server port
EXPOSE 4200

# Pokreni Angular dev server sa live reload i poll da Docker vidi promene fajlova
CMD ["npx", "ng", "serve", "--host", "0.0.0.0", "--poll=2000"]
