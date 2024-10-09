# build
# GOOS=linux GOARCH=amd64 go build -o frp-filter-linux-amd64

set -e

cd web
npm run build
cd ../
cp -r web/dist app/

echo ''

targetname=frp-filter-linux-amd64

# build by docker
if docker images | grep -qw "frp-filter-builder"; then
    echo "The image frp-filter-builder exists."
else
    docker build --build-arg targetname=$targetname -t frp-filter-builder .
fi

docker run --rm -e HTTPS_PROXY='' -e HTTP_PROXY='' -e http_proxy='' -e https_proxy='' -v $(pwd)/app:/app frp-filter-builder sh -c "cd /app && go env -w  GOPROXY=https://goproxy.io,direct && go build -o $targetname ."
mv app/$targetname ./
#rm app/dist -r

# docker create --name frp-filter-builder frp-filter-builder
# docker cp frp-filter-builder:/app/$targetname .
# docker rm frp-filter-builder

sudo chmod +x $targetname
sudo chown $(id -u):$(id -g) $targetname
