terraform {
  backend "s3" {
    bucket = "basset-terraform-state-storage"
    key    = "development/api-geo/terraform.tfstate"
    region = "us-east-1"
  }
}

variable "app_version" {}
variable "priority" {}

module "api-geo" {
  source = "git::https://9599414b9b6f5431d7372d504e8e9f5d76440329@github.com/basset-la/infrastructure//modules/ec2?ref=5.4.0"

  application_name = "api-geo"
  application_version = "${var.app_version}"
  application_path = "/geo"
  application_health_check = "/health-check"
  application_memory_assigned = 256
  application_min_capacity = 1
  application_max_capacity = 4
  application_desired_count = 1
  environment = "development"
  ami = "ami-0b9a214f40c38d5eb"
  instance_type = "t3.nano"
  key_pair_name = "basset-dev-api-kp"
  host_port = 80
  container_port = 8080
  listener_priority = "${var.priority}"
  listener_arn = "arn:aws:elasticloadbalancing:us-east-1:215088504831:listener/app/basset-dev-api-private-lb/ba068a12b68150c4/312b54da43befbd3"
  vpc_id = "vpc-9d15a5e5"
  subnet_ids = [ "subnet-edccfeb0", "subnet-95786bf1", "subnet-08fc5f07", "subnet-2b84a814", "subnet-a6f533ec", "subnet-f45a66db", ]
  security_groups = [ "sg-14696866", ]
  associate_public_ip_address = false
  newrelic_key = "1bb55c167a9cd56851acc0e1225fb9a92a43dd7c"
}
