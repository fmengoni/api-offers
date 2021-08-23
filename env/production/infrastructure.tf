terraform {
  backend "s3" {
    bucket = "basset-terraform-state-storage"
    key    = "production/api-geo/terraform.tfstate"
    region = "us-east-1"
  }
}

variable "app_version" {}
variable "priority" {}

module "api-geo" {
  source = "git::https://9599414b9b6f5431d7372d504e8e9f5d76440329@github.com/basset-la/infrastructure//modules/ec2?ref=5.4.0"

  application_name            = "api-geo"
  application_version         = "${var.app_version}"
  application_path            = "/geo"
  application_health_check    = "/health-check"
  application_memory_assigned = 256
  application_min_capacity    = 2
  application_max_capacity    = 6
  application_desired_count   = 2
  environment                 = "production"
  ami                         = "ami-0b9a214f40c38d5eb"
  instance_type               = "t3.nano"
  key_pair_name               = "basset-api-kp"
  host_port                   = 80
  container_port              = 8080
  listener_priority           = "${var.priority}"
  listener_arn                = "arn:aws:elasticloadbalancing:us-east-1:215088504831:listener/app/basset-api-private-lb/843d06b8585745cd/da0e58ce3b1b093a"
  vpc_id                      = "vpc-2b2b3e52"
  subnet_ids                  = ["subnet-da3960d6", "subnet-6b39fe54", "subnet-b339c9f8", "subnet-c1e4cfed", "subnet-6262b806", "subnet-740a262e"]
  security_groups             = ["sg-021ccd71"]
  associate_public_ip_address = false
  newrelic_key                = "badc3d500b4fb4b0cb607962c545ef67ae9c182c"
}
