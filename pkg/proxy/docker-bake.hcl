group "default" {
  targets = ["build1", "build2", "build3", "build4", "build5"]
}

target "build1" {
  context = "./pkg/proxy"
  dockerfile = "Dockerfile.large"
  tags = ["test:1"]
  args = {
    UNIQUE_ID = "1"
  }
  cache-from = ["type=gha,url=http://127.0.0.1:49160/,version=1,scope=build1"]
  cache-to = ["type=gha,url=http://127.0.0.1:49160/,mode=max,version=1,scope=build1"]
}

target "build2" {
  context = "./pkg/proxy"
  dockerfile = "Dockerfile.large"
  tags = ["test:2"]
  args = {
    UNIQUE_ID = "2"
  }
  cache-from = ["type=gha,url=http://127.0.0.1:49160/,version=1,scope=build2"]
  cache-to = ["type=gha,url=http://127.0.0.1:49160/,mode=max,version=1,scope=build2"]
}

target "build3" {
  context = "./pkg/proxy"
  dockerfile = "Dockerfile.large"
  tags = ["test:3"]
  args = {
    UNIQUE_ID = "3"
  }
  cache-from = ["type=gha,url=http://127.0.0.1:49160/,version=1,scope=build3"]
  cache-to = ["type=gha,url=http://127.0.0.1:49160/,mode=max,version=1,scope=build3"]
}

target "build4" {
  context = "./pkg/proxy"
  dockerfile = "Dockerfile.large"
  tags = ["test:4"]
  args = {
    UNIQUE_ID = "4"
  }
  cache-from = ["type=gha,url=http://127.0.0.1:49160/,version=1,scope=build4"]
  cache-to = ["type=gha,url=http://127.0.0.1:49160/,mode=max,version=1,scope=build4"]
}

target "build5" {
  context = "./pkg/proxy"
  dockerfile = "Dockerfile.large"
  tags = ["test:5"]
  args = {
    UNIQUE_ID = "5"
  }
  cache-from = ["type=gha,url=http://127.0.0.1:49160/,version=1,scope=build5"]
  cache-to = ["type=gha,url=http://127.0.0.1:49160/,mode=max,version=1,scope=build5"]
}

