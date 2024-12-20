name: Build TGExporter
run-name: ${{ gitea.actor }} Build TGExporter
on: [push]

jobs:
  Build:
    runs-on: ubuntu-latest
    steps:
      - name: Check out repository code
        uses: actions/checkout@v4

      - name: Set up Node.js
        uses: actions/setup-node@v3
        with:
          node-version: "20"

      - name: Install dependencies and build frontend
        run: |
          cd frontend
          npm config set registry https://npm.hub.ipao.vip
          npm install
          npm run build

      - name: Set up Go
        uses: actions/setup-go@v3
        with:
          go-version: "1.22"

      - name: Build Go application
        run: |
          cd backend
          mkdir -p build
          go env -w GOPROXY=https://go.hub.ipao.vip,direct
          go env -w GONOPROXY='git.ipao.vip'
          go env -w GONOSUMDB='git.ipao.vip'
          go mod tidy
          CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/app .

      - name: Build final Docker image
        run: |
          docker login -u ${{ secrets.DOCKER_AF_USERNAME }} -p ${{ secrets.DOCKER_AF_PASSWORD }} docker-af.hub.ipao.vip
          docker build --push -t docker-af.hub.ipao.vip/rogeecn/qvyun:latest  .
