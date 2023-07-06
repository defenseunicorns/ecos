# VPC
resource "aws_vpc" "ecos_vpc" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "crossplane-secure-storage"
  }
}

# Internet Gateway
resource "aws_internet_gateway" "ecos_igw" {
  vpc_id = aws_vpc.ecos_vpc.id

  tags = {
    Name = "ecos-igw"
  }
}
