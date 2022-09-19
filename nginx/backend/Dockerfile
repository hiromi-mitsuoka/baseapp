# https://qiita.com/yoship1639/items/92405ab31779c8527c08#dockerfile-nginx%E7%94%A8
FROM nginx:latest
# https://qiita.com/YumaInaura/items/3432cc3f8a8553e05a6e
# -y : 問い合わせがあった場合は，全て「y」と答える
# RUN apt-get update && apt-get install -y vim
COPY ./nginx.conf /etc/nginx/nginx.conf
COPY ./default.conf /etc/nginx/conf.d/default.conf