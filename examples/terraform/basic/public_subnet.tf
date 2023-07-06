# Public subnet
resource "aws_subnet" "public_subnet" {
  count = length(var.azs)

  vpc_id     = aws_vpc.ecos_vpc.id
  cidr_block = "10.0.${count.index}.0/24"

  map_public_ip_on_launch = true

  tags = {
    Name = "ecos-public-${var.azs[count.index]}"
  }
}

# Route table
resource "aws_route_table" "public_rte" {
  vpc_id = aws_vpc.ecos_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.ecos_igw.id
  }

  tags = {
    Name = "ecos-public-rte"
  }
}

# Associate the public route table with the public subnet
resource "aws_route_table_association" "public_assoc" {
  count = length(var.azs)

  subnet_id      = aws_subnet.public_subnet[count.index].id
  route_table_id = aws_route_table.public_rte.id
}

# Associate the public route table as the VPC's main route table
resource "aws_main_route_table_association" "public_assoc" {
  vpc_id         = aws_vpc.ecos_vpc.id
  route_table_id = aws_route_table.public_rte.id
}
