terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 2.22.0"
    }
  }
}

provider "docker" {
  host = "unix:///var/run/docker.sock"
}

resource "docker_network" "infrastructure" {
  name   = "infrastructure"
  driver = "bridge"
}

resource "docker_volume" "mysql_data" {
  name = "mysql_data"
}

resource "docker_container" "db" {
  name  = "db-go-modular"
  image = "mysql:8.0"
  ports {
    internal = 3306
    external = 3307
  }
  env = [
    "MYSQL_ROOT_HOST=%",
    "MYSQL_ROOT_PASSWORD=root_admin",
    "MYSQL_DATABASE=backend_modules",
    "MYSQL_USER=user",
    "MYSQL_PASSWORD=password",
    "TZ=Asia/Jakarta"
  ]
  networks_advanced {
    name = docker_network.infrastructure.name
  }
  volumes {
    volume_name    = docker_volume.mysql_data.name
    container_path = "/var/lib/mysql"
  }
  healthcheck {
    test         = ["CMD-SHELL", "mysqladmin ping -h localhost -uroot -proot_admin || exit 1"]
    interval     = "10s"
    timeout      = "5s"
    retries      = 3
    start_period = "30s"
  }
}

resource "docker_image" "backend_modules_image" {
  name = "backend_modules:latest"
  build {
    path       = "."
    dockerfile = "Dockerfile"
  }
}

resource "docker_container" "backend_modules" {
  name  = "go-modular-backend"
  image = docker_image.backend_modules_image.name
  ports {
    internal = 9000
    external = 9000
  }
  restart = "always"
  depends_on = [
    docker_container.db
  ]
  networks_advanced {
    name = docker_network.infrastructure.name
  }
}
