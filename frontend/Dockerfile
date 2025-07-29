# 构建阶段
FROM node:18-alpine as build-stage

# 安装yarn
RUN corepack enable && corepack prepare yarn@stable --activate

WORKDIR /app

# 复制package.json和yarn.lock
COPY package.json yarn.lock* ./

# 安装依赖
RUN yarn install --frozen-lockfile

# 复制源代码
COPY . .

# 构建生产版本
RUN yarn build

# 生产阶段
FROM nginx:alpine

# 复制构建产物到nginx目录
COPY --from=build-stage /app/dist /usr/share/nginx/html

# 复制nginx配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 暴露80端口
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]