variable "aws_region" {
  default = "us-east-1"
}

variable "azs" {
  default = ["us-east-1a", "us-east-1b"]
}

variable "bastion_type" {
  default = "t2.small" # 2 x 2
}

variable "k8s_type" {
  default = "t2.xlarge" # 4 x 16
}

variable "amis" {
  default = {
    us-east-1 = "ami-0cff7528ff583bf9a" # AMD 64-bit AWS Linux 2
  }
}

